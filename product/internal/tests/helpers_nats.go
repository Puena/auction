package tests

import (
	"context"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type natsHelpers struct {
}

func (h *natsHelpers) CreateStream(ctx context.Context, connection *nats.Conn, conf jetstream.StreamConfig) (jetstream.Stream, error) {
	js, err := jetstream.New(connection)
	if err != nil {
		return nil, err
	}

	stream, err := js.CreateOrUpdateStream(ctx, conf)
	if err != nil {
		return nil, err
	}

	return stream, err
}
