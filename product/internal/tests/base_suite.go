package tests

import (
	"context"
	"time"

	"github.com/stretchr/testify/suite"
	"golang.org/x/sync/errgroup"
)

type baseTestSuite struct {
	suite.Suite
	postgresContainer *postgresContainer
	natsContainer     *natsContainer
}

func (s *baseTestSuite) SetupSuite() {
	var err error

	ctx := context.Background()
	// create postgres container
	s.postgresContainer, err = newPostgresContainer(ctx)
	s.Require().NoError(err)

	// create nats container
	s.natsContainer, err = newNatsContainer(ctx)
	s.Require().NoError(err)
}

func (s *baseTestSuite) TearDownSuite() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	errGroup, groupCtx := errgroup.WithContext(ctx)
	errGroup.Go(func() error {
		return s.postgresContainer.Container.Terminate(groupCtx)
	})
	errGroup.Go(func() error {
		return s.natsContainer.Container.Terminate(groupCtx)
	})

	s.Require().NoError(errGroup.Wait())
}
