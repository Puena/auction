package app

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	// import http pprof
	"net/http"
	_ "net/http/pprof"

	"github.com/Puena/auction/logger"
	natsDrivenAdapter "github.com/Puena/auction/product/internal/adapters/driven/nats"
	pgDrivenAdapter "github.com/Puena/auction/product/internal/adapters/driven/postgres"
	natsDriverAdapter "github.com/Puena/auction/product/internal/adapters/driver/nats"
	"github.com/Puena/auction/product/internal/core/service"
	"github.com/Puena/auction/product/internal/database/postgres"
	natsPkg "github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"golang.org/x/sync/errgroup"
)

type app struct {
	Config *config
}

// NewApp creates a new app.
func NewApp(cfg *config) *app {
	return &app{
		Config: cfg,
	}
}

// Run runs the app, blocking, return error if something happened.
func (a *app) Run(ctx context.Context) error {
	sigs := make(chan os.Signal, 1)
	done := make(chan error, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// init pprof server
	httpServer := a.initHttpServer()
	// init nats
	nats, natsJS, err := a.initNatsJS()
	if err != nil {
		return err
	}
	// init database
	db, err := a.initPostgres()
	if err != nil {
		return err
	}

	// init consumers err group and start consuming messages with err group
	consumers := a.initConsumers(natsJS, db)
	errGroup, groupCtx := a.initErrGroup(ctx)
	errGroup.Go(func() error {
		return consumers.ConsumeCreateProductCommand(groupCtx)
	})
	errGroup.Go(func() error {
		return consumers.ConsumeUpdateProductCommand(groupCtx)
	})
	errGroup.Go(func() error {
		return consumers.ConsumeDeleteProductCommand(groupCtx)
	})
	errGroup.Go(func() error {
		return consumers.ConsumeFindProductQuery(groupCtx)
	})
	errGroup.Go(func() error {
		return consumers.ConsumeFindProductsQuery(groupCtx)
	})
	errGroup.Go(func() error {
		return a.listenHttpServer(httpServer)
	})

	// handle signals
	go func() {
		select {
		case <-sigs:
			logger.Info().Msg("received signal to stop, start terminating...")
			err := a.gracefulStop(nats, db, httpServer)
			done <- err
		case <-groupCtx.Done():
			logger.Info().Msg("context done, start terminating...")
			err := a.gracefulStop(nats, db, httpServer)
			done <- err
		}
	}()

	// listen err group
	go a.listenErrGroup(errGroup)

	// wait until done
	err = <-done
	if err != nil {
		logger.Error().Err(err).Msg("app terminated with error")
	} else {
		logger.Info().Msg("app terminated successfully")
	}
	return err
}

// RunMigration runs the app database migration, return error if something happened.
func (a *app) RunMigration() error {
	db, err := a.initPostgres()
	if err != nil {
		return err
	}

	if err := postgres.UpMigration(db); err != nil {
		return fmt.Errorf("error running postgres up migration: %w", err)
	}
	err = db.Close()
	if err != nil {
		return fmt.Errorf("error closing postgres connection: %w", err)
	}
	return nil
}

func (a *app) initNatsJS() (*natsPkg.Conn, jetstream.JetStream, error) {
	nats, err := natsPkg.Connect(a.Config.NatsURL)
	if err != nil {
		return nil, nil, fmt.Errorf("error connecting to nats: %w", err)
	}

	js, err := jetstream.New(nats)
	if err != nil {
		return nil, nil, fmt.Errorf("error connecting to nats jetstream: %w", err)
	}

	return nats, js, nil
}

func (a *app) initPostgres() (*sql.DB, error) {
	db, err := postgres.Connect(postgres.Config{
		DSN: a.Config.PostgresDSN,
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (a *app) initHttpServer() *http.Server {
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
	return &http.Server{
		Addr: fmt.Sprintf(":%d", a.Config.HttpPort),
	}
}

func (a *app) listenHttpServer(httpServer *http.Server) error {
	logger.Info().Str("port", fmt.Sprintf("%d", a.Config.HttpPort)).Msg("start listening pprof...")
	err := httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		logger.Error().Err(err).Msg("error listening pprof")
		return err
	}

	return nil
}

func (a *app) initConsumers(natsJS jetstream.JetStream, db *sql.DB) *natsDriverAdapter.ProductStreamConsumer {
	productRepo := pgDrivenAdapter.NewProductRepository(db)
	natsPublishRepo := natsDrivenAdapter.NewPublishEventRepository(natsJS, natsDrivenAdapter.Config{
		ProductStreamHeaderAuthUserID: a.Config.NatsHeaderAuthUserID,
		ProductStreamHeaderOccuredAt:  a.Config.NatsHeaderOccuredAt,
		SubjectEventProductCreated:    a.Config.SubjectEventProductCreated,
		SubjectEventProductUpdated:    a.Config.SubjectEventProductUpdated,
		SubjectEventProductDeleted:    a.Config.SubjectEventProductDeleted,
		SubjectEventProductFound:      a.Config.SubjectEventProductFound,
		SubjectEventProductsFound:     a.Config.SubjectEventProductFound,
		SubjectEventProductError:      a.Config.SubjectEventProductError,
	})
	services := service.NewProductService(productRepo, natsPublishRepo)
	consumers := natsDriverAdapter.NewProductStreamConsumer(natsJS, services, natsDriverAdapter.Config{
		AppName:                            a.Config.AppName,
		ProductStreamName:                  a.Config.NameProductStream,
		NatsHeaderAuthUserID:               a.Config.NatsHeaderAuthUserID,
		NatsHeaderOccuredAt:                a.Config.NatsHeaderOccuredAt,
		NatsHeaderMsgID:                    a.Config.NatsHeaderMsgID,
		SubjectProductCommandCreateProduct: a.Config.SubjectCommandCreateProduct,
		SubjectProductCommandUpdateProduct: a.Config.SubjectCommandUpdateProduct,
		SubjectProductCommandDeleteProduct: a.Config.SubjectCommandDeleteProduct,
		SubjectProductQueryFindProduct:     a.Config.SubjectQueryFindProduct,
		SubjectProductQueryFindProducts:    a.Config.SubjectQueryFindProducts,
	})
	return consumers
}

func (a *app) initErrGroup(ctx context.Context) (*errgroup.Group, context.Context) {
	errGroup, gCtx := errgroup.WithContext(ctx)

	return errGroup, gCtx
}

func (a *app) listenConsumers(errGroup *errgroup.Group, consumers *natsDriverAdapter.ProductStreamConsumer) error {
	err := errGroup.Wait()
	if err != nil {
		return fmt.Errorf("error listening err group: %w", err)
	}
	return nil
}

func (a *app) listenErrGroup(errGroup *errgroup.Group) error {
	err := errGroup.Wait()
	if err != nil {
		return fmt.Errorf("error listening err group: %w", err)
	}
	return nil
}

func (a *app) gracefulStop(nats *natsPkg.Conn, db *sql.DB, pprofServer *http.Server) error {
	if err := nats.Drain(); err != nil {
		return fmt.Errorf("drain nats connection error: %w", err)
	}

	if err := db.Close(); err != nil {
		return fmt.Errorf("close db connection error: %w", err)
	}

	if err := pprofServer.Close(); err != nil {
		return fmt.Errorf("close pprof server error: %w", err)
	}

	return nil
}
