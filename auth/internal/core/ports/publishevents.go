package ports

import (
	"context"

	"github.com/Puena/auction-house-auth/internal/core/domain"
)

type PublishEvents interface {
	PublishEventUserRevoked(ctx context.Context, event *domain.EventUserRevoked) error
	PublishEventUserUpdated(ctx context.Context, event *domain.EventUserUpdated) error
	PublishEventUserVerified(ctx context.Context, event *domain.EventUserVerified) error
	PublishEventUserData(ctx context.Context, event *domain.EventUserData) error
	PublishEventUserError(ctx context.Context, event *domain.EventUserError) error
}
