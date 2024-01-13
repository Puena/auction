package tests

import (
	"context"
	"fmt"
	"os"
	"syscall"
	"testing"

	"github.com/Puena/auction/product/internal/app"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	// streams
	natsProductStreamName          = "product"
	natsProductStreamSubjectFilter = "product.>"
	natsProductRetentionPolicy     = jetstream.WorkQueuePolicy
	natsErrorStreamName            = "error"
	natsErrorStreamSubjectFilter   = "error.>"
	natsErrorStreamRetentionPolicy = jetstream.LimitsPolicy
	// headers
	natsHeaderAuthUserID = "Auth-User-ID"
	natsHeaderOccuredAt  = "Occured-At"
	NatsHeaderMsgID      = "Msg-Id"
	// subjects
	natsSubjectEventProductCreated  = "product.event.product_created"
	natsSubjectEventProductUpdated  = "product.event.product_updated"
	natsSubjectEventProductDeleted  = "product.event.product_deleted"
	natsSubjectEventProductFound    = "product.event.product_found"
	natsSubjectEventProductsFound   = "product.event.products_found"
	natsSubjectCommandCreateProduct = "product.command.create_product"
	natsSubjectCommandUpdateProduct = "product.command.update_product"
	natsSubjectCommandDeleteProduct = "product.command.delete_product"
	natsSubjectQueryFindProduct     = "product.query.find_product"
	natsSubjectQueryFindProducts    = "product.query.find_products"
)

type ProductTestSuite struct {
	baseTestSuite
	natsHelpers
	productStream jetstream.Stream
	natsConn      *nats.Conn
	appPid        int
}

func (s *ProductTestSuite) SetupSuite() {
	s.baseTestSuite.SetupSuite()

	var err error
	s.natsConn, err = nats.Connect(s.natsContainer.URI)
	s.Require().NoError(err, "error connecting to nats")

	s.productStream, err = s.CreateStream(context.Background(), s.natsConn, jetstream.StreamConfig{
		Name:      natsProductStreamName,
		Subjects:  []string{natsProductStreamSubjectFilter},
		Retention: natsProductRetentionPolicy,
	})
	s.Require().NoError(err, "error creating nats stream")

	appConfig, err := app.NewConfig()
	s.Require().NoError(err, "error creating app config")
	appConfig.AppName = "product-service-test"
	appConfig.PostgresDSN = s.postgresContainer.URI
	appConfig.NatsURL = s.natsContainer.URI
	appConfig.NameProductStream = natsProductStreamName
	appConfig.NameErrorStream = natsErrorStreamName
	appConfig.NatsHeaderAuthUserID = natsHeaderAuthUserID
	appConfig.NatsHeaderOccuredAt = natsHeaderOccuredAt
	appConfig.NatsHeaderMsgID = NatsHeaderMsgID
	appConfig.SubjectEventProductCreated = natsSubjectEventProductCreated
	appConfig.SubjectEventProductUpdated = natsSubjectEventProductUpdated
	appConfig.SubjectEventProductDeleted = natsSubjectEventProductDeleted
	appConfig.SubjectEventProductFound = natsSubjectEventProductFound
	appConfig.SubjectEventProductsFound = natsSubjectEventProductsFound
	appConfig.SubjectCommandCreateProduct = natsSubjectCommandCreateProduct
	appConfig.SubjectCommandUpdateProduct = natsSubjectCommandUpdateProduct
	appConfig.SubjectCommandDeleteProduct = natsSubjectCommandDeleteProduct
	appConfig.SubjectQueryFindProduct = natsSubjectQueryFindProduct
	appConfig.SubjectQueryFindProducts = natsSubjectQueryFindProducts
	appConfig.HttpPort = 6040

	app := app.NewApp(appConfig)
	err = app.RunMigration()
	require.NoError(s.T(), err, "error running migration")

	go func() {
		s.appPid = os.Getpid()
		s.T().Log(fmt.Sprintf("running app pid: %d", s.appPid))
		err := app.Run(context.Background())
		require.NoError(s.T(), err, "error running app")
	}()
}

func (s *ProductTestSuite) TearDownSuite() {
	p, err := os.FindProcess(s.appPid)
	require.NoError(s.T(), err, "error finding process")
	p.Signal(syscall.SIGTERM)

	err = s.natsConn.Drain()
	s.Require().NoError(err, "error draining nats connection")

	s.baseTestSuite.TearDownSuite()
}

func (s *ProductTestSuite) TestCreateProductCommand() {
	// Arrange
	s.Assert().Equal(0, 0)
}

func TestProductSuite(t *testing.T) {
	suite.Run(t, new(ProductTestSuite))
}
