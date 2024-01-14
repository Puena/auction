package tests

import (
	"context"
	"fmt"
	"os"
	"sync/atomic"
	"syscall"
	"testing"
	"time"

	"github.com/Puena/auction/pbgo/auction"
	"github.com/Puena/auction/product/internal/app"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	// streams
	natsProductStreamName          = "product"
	natsProductStreamSubjectFilter = "product.>"
	natsProductRetentionPolicy     = jetstream.LimitsPolicy // for testing purposes in prod it will be workqueue
	natsErrorStreamName            = "error"
	natsErrorStreamSubjectFilter   = "error.>"
	natsErrorStreamRetentionPolicy = jetstream.LimitsPolicy
	// headers
	natsHeaderAuthUserID   = "Auth-User-ID"
	natsHeaderOccuredAt    = "Occured-At"
	natsHeaderMsgID        = "Nats-Msg-Id"
	natsHeaderReplyToMsgID = "Reply-To-Msg-Id"
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
	appPid        atomic.Uint32
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
	appConfig.NatsHeaderMsgID = natsHeaderMsgID
	appConfig.NatsHeaderReplyToMsgID = natsHeaderReplyToMsgID
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
		pid := os.Getpid()
		s.appPid.Store(uint32(pid))
		s.T().Log(fmt.Sprintf("running app pid: %d", pid))
		err := app.Run(context.Background())
		require.NoError(s.T(), err, "error running app")
	}()
}

func (s *ProductTestSuite) TearDownSuite() {
	pid := s.appPid.Load()
	p, err := os.FindProcess(int(pid))
	require.NoError(s.T(), err, "error finding process")
	p.Signal(syscall.SIGTERM)

	err = s.natsConn.Drain()
	s.Require().NoError(err, "error draining nats connection")

	s.baseTestSuite.TearDownSuite()
}

func (s *ProductTestSuite) publishToProductStream(ctx context.Context, subject string, userID string, msgID string, data protoreflect.ProtoMessage) (*jetstream.PubAck, error) {
	js, err := jetstream.New(s.natsConn)
	s.Require().NoError(err, "error creating nats jetstream")
	marshaled, err := proto.Marshal(data)
	s.Require().NoError(err, "error marshaling proto message")
	msg := nats.NewMsg(subject)
	msg.Header.Add(natsHeaderAuthUserID, userID)
	msg.Header.Add(natsHeaderOccuredAt, time.Now().Format(time.RFC3339))
	msg.Header.Add(natsHeaderMsgID, msgID)
	msg.Data = marshaled

	return js.PublishMsg(ctx, msg)
}

func (s *ProductTestSuite) listenProductStreamEvent(ctx context.Context, subject string) (jetstream.Consumer, error) {
	c, err := s.productStream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		DeliverPolicy: jetstream.DeliverNewPolicy,
		AckPolicy:     jetstream.AckExplicitPolicy,
		FilterSubject: subject,
	})
	s.Require().NoError(err, "error creating nats consumer")
	return c, err
}

func (s *ProductTestSuite) TestCRUDCommands() {
	// Arrange
	userID := ulid.Make().String()
	createProductCommand := &auction.CommandCreateProduct{
		Key: ulid.Make().String(),
		Value: &auction.CreateProduct{
			Name:        "Steering wheel",
			Media:       []string{},
			Description: "Ancient steering wheel",
			CreatedBy:   userID,
		},
	}

	// create consumer from product created events
	c, err := s.listenProductStreamEvent(context.Background(), natsSubjectEventProductCreated)
	s.Assert().NoError(err, "error creating nats consumer")
	found := make(chan protoreflect.ProtoMessage, 1)
	productCreatedConsumer, err := c.Consume(func(msg jetstream.Msg) {
		err := msg.Ack()
		s.Assert().NoError(err, "error acking nats message")

		s.assertMsgHeaders(msg)

		replyMsgID := msg.Headers().Get(natsHeaderReplyToMsgID)
		if replyMsgID != createProductCommand.Key {
			return // this can be possible if test run in parallel
		}
		// get data
		created := &auction.EventProductCreated{}
		err = proto.Unmarshal(msg.Data(), created)
		s.Assert().NoError(err, "error unmarshaling data")

		// Assert
		s.Assert().NotEmpty(created.Value.Id, "id should not be empty")
		s.Assert().Equal(createProductCommand.Value.Name, created.Value.Name, "name should be the same")
		s.Assert().Equal(createProductCommand.Value.Description, created.Value.Description, "description should be the same")
		s.Assert().Equal(len(createProductCommand.Value.Media), len(created.Value.Media), "media should be the same")
		s.Assert().Equal(createProductCommand.Value.CreatedBy, created.Value.CreatedBy, "created by should be the same")
		s.Assert().Equal(userID, created.Value.CreatedBy, "created should be the same as userID")
		s.Assert().NotEmpty(created.Value.CreatedAt, "created at should not be empty")

		found <- created
	})
	s.Assert().NoError(err, "error consuming nats stream")
	defer productCreatedConsumer.Stop()

	_, err = s.publishToProductStream(context.Background(), natsSubjectCommandCreateProduct, userID, createProductCommand.Key, createProductCommand)
	s.Assert().NoError(err, "error publishing to nats")

	entity := <-found
	createdProduct, ok := entity.(*auction.EventProductCreated)
	s.Assert().True(ok, "entity should be of type EventProductCreated")

	// update product
	commandUpdateProduct := &auction.CommandUpdateProduct{
		Key: ulid.Make().String(),
		Value: &auction.UpdateProduct{
			Id:        createdProduct.Value.Id,
			Media:     []string{"some/path"},
			CreatedBy: createdProduct.Value.CreatedBy,
		},
	}

	epu, err := s.listenProductStreamEvent(context.Background(), natsSubjectEventProductUpdated)
	s.Assert().NoError(err, "error creating nats consumer")
	productUpdatedConsumer, err := epu.Consume(func(msg jetstream.Msg) {
		err := msg.Ack()
		s.Assert().NoError(err, "error acking nats message")

		s.assertMsgHeaders(msg)
		replyMsgID := msg.Headers().Get(natsHeaderReplyToMsgID)
		if replyMsgID != commandUpdateProduct.Key {
			return // this can be possible if test run in parallel
		}

		// get data
		updated := &auction.EventProductUpdated{}
		err = proto.Unmarshal(msg.Data(), updated)
		s.Assert().NoError(err, "error unmarshaling data")

		// Assert
		s.Assert().Equal(commandUpdateProduct.Value.Id, updated.Value.Id, "id should be the same")
		s.Assert().Equal(commandUpdateProduct.Value.Media, updated.Value.Media, "media should be the same")
		s.Assert().Equal(commandUpdateProduct.Value.CreatedBy, updated.Value.CreatedBy, "created by should be the same")
		s.Assert().Equal(createdProduct.Value.Name, updated.Value.Name, "name should be the same")
		s.Assert().Equal(createdProduct.Value.CreatedBy, updated.Value.CreatedBy, "created should be the same as userID")
		s.Assert().Equal(createdProduct.Value.CreatedAt, updated.Value.CreatedAt, "created at should be the same")
		s.Assert().NotEmpty(updated.Value.UpdatedAt, "updated at should not be empty")

		found <- updated
	})
	s.Assert().NoError(err, "error consuming nats stream")
	defer productUpdatedConsumer.Stop()

	_, err = s.publishToProductStream(context.Background(), natsSubjectCommandUpdateProduct, userID, commandUpdateProduct.Key, commandUpdateProduct)
	s.Assert().NoError(err, "error publishing update command to nats")

	entity = <-found
	updatedProduct, ok := entity.(*auction.EventProductUpdated)
	s.Assert().True(ok, "entity should be of type EventProductUpdated")

	// found query
	query := &auction.QueryFindProduct{
		Key: ulid.Make().String(),
		Value: &auction.FindProduct{
			Id: createdProduct.Value.Id,
		},
	}

	epf, err := s.listenProductStreamEvent(context.Background(), natsSubjectEventProductFound)
	s.Assert().NoError(err, "error creating nats consumer")

	productFoundConsumer, err := epf.Consume(func(msg jetstream.Msg) {
		err := msg.Ack()
		s.Assert().NoError(err, "error acking nats message")

		s.assertMsgHeaders(msg)
		replyMsgID := msg.Headers().Get(natsHeaderReplyToMsgID)
		if replyMsgID != query.Key {
			return // this can be possible if test run in parallel
		}

		// get data
		productFounded := &auction.EventProductFound{}
		err = proto.Unmarshal(msg.Data(), productFounded)
		s.Assert().NoError(err, "error unmarshaling data")

		// Assert
		s.Assert().Equal(query.Value.Id, productFounded.Value.Id, "id should be the same")
		s.Assert().Equal(updatedProduct.Value.Name, productFounded.Value.Name, "name should be the same")
		s.Assert().Equal(updatedProduct.Value.Description, productFounded.Value.Description, "description should be the same")
		s.Assert().Equal(updatedProduct.Value.Media, productFounded.Value.Media, "media should be the same")
		s.Assert().Equal(updatedProduct.Value.CreatedBy, productFounded.Value.CreatedBy, "created by should be the same")
		s.Assert().Equal(updatedProduct.Value.CreatedAt, productFounded.Value.CreatedAt, "created at should be the same")
		s.Assert().Equal(updatedProduct.Value.UpdatedAt, productFounded.Value.UpdatedAt, "updated at should be the same")

		found <- productFounded
	})
	s.Assert().NoError(err, "error consuming nats stream")
	defer productFoundConsumer.Stop()

	_, err = s.publishToProductStream(context.Background(), natsSubjectQueryFindProduct, userID, query.Key, query)
	s.Assert().NoError(err, "error publishing query to nats")

	entity = <-found
	foundedProduct, ok := entity.(*auction.EventProductFound)
	s.Assert().True(ok, "entity should be of type EventProductFound")

	// delete product
	commandDeleteProduct := &auction.CommandDeleteProduct{
		Key: ulid.Make().String(),
		Value: &auction.DeleteProduct{
			Id:        foundedProduct.Value.Id,
			CreatedBy: foundedProduct.Value.CreatedBy,
		},
	}

	epd, err := s.listenProductStreamEvent(context.Background(), natsSubjectEventProductDeleted)
	s.Assert().NoError(err, "error creating nats consumer")
	productDeletedConsumer, err := epd.Consume(func(msg jetstream.Msg) {
		err := msg.Ack()
		s.Assert().NoError(err, "error acking nats message")

		s.assertMsgHeaders(msg)
		replyMsgID := msg.Headers().Get(natsHeaderReplyToMsgID)
		if replyMsgID != commandDeleteProduct.Key {
			return // this can be possible if test run in parallel
		}

		// get data
		deleted := &auction.EventProductDeleted{}
		err = proto.Unmarshal(msg.Data(), deleted)
		s.Assert().NoError(err, "error unmarshaling data")

		// Assert
		s.Assert().Equal(commandDeleteProduct.Value.Id, deleted.Value.Id, "id should be the same")
		s.Assert().Equal(commandDeleteProduct.Value.CreatedBy, deleted.Value.CreatedBy, "created by should be the same")
		s.Assert().Equal(foundedProduct.Value.Name, deleted.Value.Name, "name should be the same")
		s.Assert().Equal(foundedProduct.Value.Description, deleted.Value.Description, "description should be the same")
		s.Assert().Equal(foundedProduct.Value.Media, deleted.Value.Media, "media should be the same")
		s.Assert().Equal(foundedProduct.Value.CreatedBy, deleted.Value.CreatedBy, "created by should be the same")
		s.Assert().Equal(foundedProduct.Value.CreatedAt, deleted.Value.CreatedAt, "created at should be the same")
		s.Assert().Equal(foundedProduct.Value.UpdatedAt, deleted.Value.UpdatedAt, "updated at should be the same")

		found <- deleted
	})
	s.Assert().NoError(err, "error consuming nats stream")
	defer productDeletedConsumer.Stop()

	_, err = s.publishToProductStream(context.Background(), natsSubjectCommandDeleteProduct, userID, commandDeleteProduct.Key, commandDeleteProduct)
	s.Assert().NoError(err, "error publishing delete command to nats")

	entity = <-found
	deletedProduct, ok := entity.(*auction.EventProductDeleted)
	s.Assert().True(ok, "entity should be of type EventProductDeleted")
	fmt.Println(deletedProduct)
}

func (s *ProductTestSuite) assertMsgHeaders(msg jetstream.Msg) {
	msgID := msg.Headers().Get(natsHeaderMsgID)
	s.Assert().NotEmpty(msgID, "msg id header should not be empty")
	// assert reply msg id header
	replyMsgID := msg.Headers().Get(natsHeaderReplyToMsgID)
	s.Assert().NotEmpty(replyMsgID, "reply msg id header should not be empty")
	// assert occured at header
	occuredAtHeader := msg.Headers().Get(natsHeaderOccuredAt)
	s.Assert().NotEmpty(occuredAtHeader, "occured at header should not be empty")
	_, err := time.Parse(time.RFC3339, occuredAtHeader)
	s.Assert().NoError(err, "error when trying to parse occuredAt time from RFC3339")
	// assert auth user id header
	s.Assert().NotEmpty(msg.Headers().Get(natsHeaderAuthUserID), "auth user id header should not be empty")
}

func TestUpdateProductCommand(t *testing.T) {
}

func TestProductSuite(t *testing.T) {
	suite.Run(t, new(ProductTestSuite))
}
