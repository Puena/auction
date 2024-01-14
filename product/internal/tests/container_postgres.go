package tests

import (
	"context"
	"fmt"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	postgresContainerImage = "postgres:16"
	postgresContainerPort  = "5432"
	postgresPassword       = "some-strong-password"
	postgresUser           = "testcontainers"
	postgresDatabase       = "test"
	sslMode                = "disable"
)

type postgresContainer struct {
	Container testcontainers.Container
	URI       string
}

func newPostgresContainer(ctx context.Context) (*postgresContainer, error) {
	req := testcontainers.ContainerRequest{
		Image:        postgresContainerImage,
		ExposedPorts: []string{postgresContainerPort},
		Env: map[string]string{
			"POSTGRES_PASSWORD": postgresPassword,
			"POSTGRES_USER":     postgresUser,
			"POSTGRES_DB":       postgresDatabase,
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").WithOccurrence(2).WithStartupTimeout(10 * time.Second),
	}

	pgContainer, err := testcontainers.GenericContainer(
		ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("error creating postgres container: %w", err)
	}

	host, err := pgContainer.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting postgres container host: %w", err)
	}
	port, err := pgContainer.MappedPort(ctx, postgresContainerPort)
	if err != nil {
		return nil, fmt.Errorf("error getting postgres container port: %w", err)
	}

	return &postgresContainer{
		Container: pgContainer,
		URI:       fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", postgresUser, postgresPassword, host, port.Port(), postgresDatabase, sslMode),
	}, nil
}
