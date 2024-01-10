package firebase

import (
	fbAuth "firebase.google.com/go/auth"
	"github.com/Puena/auction-house-auth/internal/core/domain"
)

func mapDomainFirebaseUserToUpdate(toUpdate *domain.FirebaseUserToUpdate) fbAuth.UserToUpdate {
	var update fbAuth.UserToUpdate

	if toUpdate.DisplayName != nil {
		update.DisplayName(*toUpdate.DisplayName)
	}
	if toUpdate.Email != nil {
		update.Email(*toUpdate.Email)
	}
	if toUpdate.PhoneNumber != nil {
		update.PhoneNumber(*toUpdate.PhoneNumber)
	}
	if toUpdate.PhotoURL != nil {
		update.PhotoURL(*toUpdate.PhotoURL)
	}
	if toUpdate.Disabled != nil {
		update.Disabled(*toUpdate.Disabled)
	}

	return update
}

func mapUserRecordToDomainFirebaseUser(record *fbAuth.UserRecord) domain.FirebaseUser {
	return domain.FirebaseUser{
		UID:                  record.UID,
		DisplayName:          record.DisplayName,
		Email:                record.Email,
		PhoneNumber:          record.PhoneNumber,
		PhotoURL:             record.PhotoURL,
		ProviderID:           record.ProviderID,
		CreationTimestamp:    record.UserMetadata.CreationTimestamp,
		LastLogInTimestamp:   record.UserMetadata.LastLogInTimestamp,
		LastRefreshTimestamp: record.UserMetadata.LastRefreshTimestamp,
		Disabled:             record.Disabled,
	}
}
