package firebase

import (
	"context"

	fb "firebase.google.com/go"
	fbAuth "firebase.google.com/go/auth"
	"github.com/Puena/auction-house-auth/internal/core/domain"
)

type firebaseAuth struct {
	auth *fbAuth.Client
}

// NewFirebaseAuth creates a new firebase auth client.
func NewFirebaseAuth(ctx context.Context, app *fb.App) (*firebaseAuth, error) {
	if app == nil {
		// TODO: Replace to logger Fatal
		panic("firebase app is required")
	}

	cli, err := app.Auth(ctx)
	if err != nil {
		return nil, err
	}

	return &firebaseAuth{
		auth: cli,
	}, nil
}

// VerifyIDToken verifies the signature and data for the given ID token and return UID.
func (a *firebaseAuth) VerifyIDToken(ctx context.Context, token string) (string, error) {
	userToken, err := a.auth.VerifyIDToken(ctx, token)
	if err != nil {
		return "", err
	}

	return userToken.UID, nil
}

// VerifyIDTokenAndCheckRevoked verifies the signature and data for the given ID token and check if the token has been revoked, then return UID.
func (a *firebaseAuth) VerifyIDTokenAndCheckRevoked(ctx context.Context, token string) (string, error) {
	userToken, err := a.auth.VerifyIDTokenAndCheckRevoked(ctx, token)
	if err != nil {
		return "", err
	}

	return userToken.UID, nil
}

// UpdateUser updates an existing user account with the given parameters.
func (a *firebaseAuth) UpdateUser(ctx context.Context, uid string, user *domain.FirebaseUserToUpdate) (domain.FirebaseUser, error) {
	update := mapDomainFirebaseUserToUpdate(user)
	ur, err := a.auth.UpdateUser(ctx, uid, &update)
	if err != nil {
		return domain.FirebaseUser{}, err
	}

	return mapUserRecordToDomainFirebaseUser(ur), nil
}

// GetUser gets the user data corresponding to the specified user ID.
func (a *firebaseAuth) GetUser(ctx context.Context, uid string) (domain.FirebaseUser, error) {
	ur, err := a.auth.GetUser(ctx, uid)
	if err != nil {
		return domain.FirebaseUser{}, err
	}

	return mapUserRecordToDomainFirebaseUser(ur), nil
}

// RevokedRefreshTokens revokes all refresh tokens for a specified user identified by the given uid.
func (a *firebaseAuth) RevokedRefreshTokens(ctx context.Context, uid string) error {
	return a.auth.RevokeRefreshTokens(ctx, uid)
}
