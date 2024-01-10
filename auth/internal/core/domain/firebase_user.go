package domain

// FirebaseUser represents user data from firebase.
type FirebaseUser struct {
	UID                  string
	DisplayName          string
	Email                string
	PhoneNumber          string
	PhotoURL             string
	ProviderID           string
	CreationTimestamp    int64
	LastLogInTimestamp   int64
	LastRefreshTimestamp int64
	Disabled             bool
}

type FirebaseUserToUpdate struct {
	DisplayName *string
	Email       *string
	PhoneNumber *string
	PhotoURL    *string
	Disabled    *bool
}

type FirebaseUserID struct {
	UID string
}
