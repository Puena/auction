package ports

import (
	"context"

	"github.com/Puena/auction-house-auth/internal/core/domain"
)

// FirebaseAuth represents the firebase auth driven port.
type FirebaseAuth interface {
	VerifyIDToken(ctx context.Context, token string) (string, error)
	VerifyIDTokenAndCheckRevoked(ctx context.Context, token string) (string, error)
	UpdateUser(ctx context.Context, uid string, user *domain.FirebaseUserToUpdate) (domain.FirebaseUser, error)
	GetUser(ctx context.Context, uid string) (domain.FirebaseUser, error)
	RevokedRefreshTokens(ctx context.Context, uid string) error
}
