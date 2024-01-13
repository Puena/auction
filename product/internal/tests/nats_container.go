package tests

import (
	"context"
	"fmt"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	natsContainerImage = "nats:2"
	natsJSFlag         = "--js"
	natsContainerPort  = "4222"
)

type natsContainer struct {
	Container testcontainers.Container
	URI       string
}

func newNatsContainer(ctx context.Context) (*natsContainer, error) {
	req := testcontainers.ContainerRequest{
		Image:        natsContainerImage,
		Cmd:          []string{natsJSFlag},
		ExposedPorts: []string{natsContainerPort},
		WaitingFor:   wait.ForLog("Server is ready").WithStartupTimeout(10 * time.Second),
	}

	ntsContainer, err := testcontainers.GenericContainer(
		ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("error creating nats container: %w", err)
	}

	host, err := ntsContainer.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting nats container host: %w", err)
	}
	port, err := ntsContainer.MappedPort(ctx, natsContainerPort)
	if err != nil {
		return nil, fmt.Errorf("error getting nats container port: %w", err)
	}

	return &natsContainer{
		Container: ntsContainer,
		URI:       fmt.Sprintf("nats://%s:%s", host, port.Port()),
	}, nil
}
