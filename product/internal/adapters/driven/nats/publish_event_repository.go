package nats

import (
	"context"
	"fmt"
	"time"

	logger "github.com/Puena/auction/logger"
	"github.com/Puena/auction/pbgo/auction"
	"github.com/Puena/auction/product/internal/core/domain"
	natsJS "github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/oklog/ulid/v2"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	natsMessageRetries = 3
)

type Config struct {
	ProductStreamHeaderAuthUserID string
	ProductStreamHeaderOccuredAt  string
	SubjectEventProductCreated    string
	SubjectEventProductUpdated    string
	SubjectEventProductDeleted    string
	SubjectEventProductFound      string
	SubjectEventProductsFound     string
	SubjectEventProductError      string
}

func (c *Config) Validate() error {
	if c.ProductStreamHeaderAuthUserID == "" {
		return fmt.Errorf("ProductStreamHeaderAuthUserID can't be empty")
	}
	if c.ProductStreamHeaderOccuredAt == "" {
		return fmt.Errorf("ProductStreamHeaderOccuredAt can't be empty")
	}
	if c.SubjectEventProductCreated == "" {
		return fmt.Errorf("SubjectProductEventProductCreated can't be empty")
	}
	if c.SubjectEventProductUpdated == "" {
		return fmt.Errorf("SubjectProductEventProductUpdated can't be empty")
	}
	if c.SubjectEventProductDeleted == "" {
		return fmt.Errorf("SubjectProductEventProductDeleted can't be empty")
	}
	if c.SubjectEventProductFound == "" {
		return fmt.Errorf("SubjectProductEventProductFound can't be empty")
	}
	if c.SubjectEventProductsFound == "" {
		return fmt.Errorf("SubjectProductEventProductsFound can't be empty")
	}

	return nil
}

type publishEventRepository struct {
	config Config
	broker jetstream.JetStream
}

// NewPublishEventRepository creates a new publish event repository.
func NewPublishEventRepository(nats jetstream.JetStream, config Config) *publishEventRepository {
	err := config.Validate()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed while validating config")
	}
	return &publishEventRepository{broker: nats, config: config}
}

func (p *publishEventRepository) setAdditionalHeaders(authUserID string, occuredAt time.Time) natsJS.Header {
	headers := natsJS.Header{}
	headers.Set(p.config.ProductStreamHeaderAuthUserID, authUserID)
	headers.Set(p.config.ProductStreamHeaderOccuredAt, occuredAt.Format(time.RFC3339))

	return headers
}

func (p *publishEventRepository) publishToNatsJetStream(subject string, msgID string, headers natsJS.Header, protoMsg protoreflect.ProtoMessage, opts ...jetstream.PublishOpt) (*jetstream.PubAck, error) {
	natsMsg := natsJS.NewMsg(subject)
	natsMsg.Header = headers
	var err error
	natsMsg.Data, err = proto.Marshal(protoMsg)
	if err != nil {
		return nil, fmt.Errorf("failed when marshaling proto message: %w", err)
	}

	configOpts := []jetstream.PublishOpt{
		jetstream.WithRetryAttempts(natsMessageRetries),
	}
	if msgID != "" {
		configOpts = append(configOpts, jetstream.WithMsgID(msgID))
	}
	opts = append(configOpts, opts...)

	return p.broker.PublishMsg(context.Background(), natsMsg, opts...)
}

// Publish event product created, userID is who initiated this action.
func (p *publishEventRepository) ProductCreated(ctx context.Context, userID string, event domain.EventProductCreated) error {
	protoMsg := &auction.EventProductCreated{
		Key: event.Value.ID,
		Value: &auction.Product{
			Id:          event.Value.ID,
			Name:        event.Value.Name,
			Media:       event.Value.Media,
			Description: event.Value.Description,
			CreatedBy:   event.Value.CreatedBy,
			CreatedAt:   timestamppb.New(event.Value.CreatedAt),
		},
	}

	_, err := p.publishToNatsJetStream(p.config.SubjectEventProductCreated, event.Value.ID, p.setAdditionalHeaders(userID, time.Now()), protoMsg)
	if err != nil {
		return err
	}

	return nil
}

// Publish event product updated, userID is who initiated this action.
func (p *publishEventRepository) ProductUpdated(ctx context.Context, userID string, event domain.EventProductUpdated) error {
	protoMsg := &auction.EventProductUpdated{
		Key: event.Value.ID,
		Value: &auction.Product{
			Id:          event.Value.ID,
			Name:        event.Value.Name,
			Media:       event.Value.Media,
			Description: event.Value.Description,
			CreatedBy:   event.Value.CreatedBy,
			CreatedAt:   timestamppb.New(event.Value.CreatedAt),
		},
	}

	_, err := p.publishToNatsJetStream(p.config.SubjectEventProductUpdated, event.Value.ID, p.setAdditionalHeaders(userID, time.Now()), protoMsg)
	if err != nil {
		return err
	}

	return nil
}

// Publish event product deleted, userID is who initiated this action.
func (p *publishEventRepository) ProductDeleted(ctx context.Context, userID string, event domain.EventProductDeleted) error {
	protoMsg := &auction.EventProductDeleted{
		Key: event.Value.ID,
		Value: &auction.Product{
			Id:          event.Value.ID,
			Name:        event.Value.Name,
			Media:       event.Value.Media,
			Description: event.Value.Description,
			CreatedBy:   event.Value.CreatedBy,
			CreatedAt:   timestamppb.New(event.Value.CreatedAt),
		},
	}

	_, err := p.publishToNatsJetStream(p.config.SubjectEventProductDeleted, event.Value.ID, p.setAdditionalHeaders(userID, time.Now()), protoMsg)
	if err != nil {
		return err
	}

	return nil
}

// Publish event product found, userID is who initiated this action, if found nothing empty array returns.
func (p *publishEventRepository) ProductFound(ctx context.Context, userID string, event domain.EventProductFound) error {
	msgID := ulid.Make().String()
	protoMsg := &auction.EventProductFound{
		Key: event.Key,
		Value: &auction.Product{
			Id:          event.Value.ID,
			Name:        event.Value.Name,
			Media:       event.Value.Media,
			Description: event.Value.Description,
			CreatedBy:   event.Value.CreatedBy,
			CreatedAt:   timestamppb.New(event.Value.CreatedAt),
		},
	}

	_, err := p.publishToNatsJetStream(p.config.SubjectEventProductsFound, msgID, p.setAdditionalHeaders(userID, time.Now()), protoMsg)
	if err != nil {
		return err
	}

	return nil
}

// Publish event products found, userID is who initiated this action, if found nothing empty array returns.
func (p *publishEventRepository) ProductsFound(ctx context.Context, userID string, event domain.EventProductsFound) error {
	msgID := ulid.Make().String()
	protoMsg := &auction.EventProductsFound{
		Key:   msgID,
		Value: []*auction.Product{},
	}
	for _, ep := range event.Value {
		protoMsg.Value = append(protoMsg.Value, eventProductToProtoProduct(ep))
	}

	_, err := p.publishToNatsJetStream(p.config.SubjectEventProductsFound, msgID, p.setAdditionalHeaders(userID, time.Now()), protoMsg)
	if err != nil {
		return err
	}

	return nil
}

// ProductError publish event error.
func (p *publishEventRepository) ProductError(ctx context.Context, userID string, event domain.EventProductError) error {
	protoMsg := &auction.EventErrorOccurred{
		Key: event.Key,
		Value: &auction.EventError{
			StreamName:        event.Value.StreamName,
			ConsumerName:      event.Value.ConsumerName,
			Subject:           event.Value.Subject,
			ReferenceEventKey: event.Value.Reference_event_key,
			Message:           event.Value.Message,
			Code:              int32(event.Value.Code),
			Data:              event.Value.Data,
			Headers:           event.Value.Headers,
			Time:              timestamppb.New(event.Value.Time),
		},
	}

	_, err := p.publishToNatsJetStream(p.config.SubjectEventProductError, event.Key, p.setAdditionalHeaders(userID, time.Now()), protoMsg)
	if err != nil {
		return err
	}

	return nil
}
