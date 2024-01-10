package domain

// EventErrorData represent [EventError] data when error occurs, oneof [FirebaseUser, FirebaseUserID].
type EventErrorData struct {
	User   *FirebaseUser
	UserID *FirebaseUserID
}

// EventError represent error struct for [EventUserError].
type EventError struct {
	EventID        string
	EventData      EventErrorData
	EventError     string
	EventCreatedAt int64
}

// EventUserRevoked represent event when user have been revoked.
type EventUserRevoked struct {
	Key   string
	Value FirebaseUserID
}

// EventUserError represent event when error occurs.
type EventUserError struct {
	Key   string
	Value EventError
}

// EventUserVerified represent event when user have been verified.
type EventUserVerified struct {
	Key   string
	Value FirebaseUserID
}

// EventUserUpdated represent event when user have been updated.
type EventUserUpdated struct {
	Key   string
	Value FirebaseUser
}

// EventUserData represent event when user have been requested.
type EventUserData struct {
	Key   string
	Value FirebaseUser
}

// NewEventUserRevoked create a new [EventUserRevoked].
func NewEventUserRevoked(key string, data *FirebaseUserID) EventUserRevoked {
	return EventUserRevoked{
		Key:   key,
		Value: *data,
	}
}

// NewEventUserVerified create a new [EventUserVerified].
func NewEventUserVerified(key string, data *FirebaseUserID) EventUserVerified {
	return EventUserVerified{
		Key:   key,
		Value: *data,
	}
}

// NewEventUserUpdated create a new [EventUserUpdated].
func NewEventUserUpdated(key string, data *FirebaseUser) EventUserUpdated {
	return EventUserUpdated{
		Key:   key,
		Value: *data,
	}
}

// NewEventUserData create a new [EventUserData].
func NewEventUserData(key string, data *FirebaseUser) EventUserData {
	return EventUserData{
		Key:   key,
		Value: *data,
	}
}

// NewEventUserError create a new [EventUserError].
func NewEventUserError(key string, eventID string, data *EventErrorData, eventError string, eventCreatedAt int64) EventUserError {
	return EventUserError{
		Key: key,
		Value: EventError{
			EventID:        eventID,
			EventData:      *data,
			EventError:     eventError,
			EventCreatedAt: eventCreatedAt,
		},
	}
}
