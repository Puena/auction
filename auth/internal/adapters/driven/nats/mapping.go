package nats

import (
	"github.com/Puena/auction-house-auth/internal/core/domain"
	auctionProto "github.com/Puena/auction-messages-golang"
)

func mapEventUserDataToProto(event *domain.EventUserData) auctionProto.EventUserData {
	return auctionProto.EventUserData{
		Data: &auctionProto.User{
			Id:        event.Value.UID,
			Name:      event.Value.DisplayName,
			Email:     event.Value.Email,
			Phone:     event.Value.PhoneNumber,
			Photo:     event.Value.PhotoURL,
			CreatedAt: event.Value.CreationTimestamp,
			LastLogin: event.Value.LastLogInTimestamp,
			Disabled:  false,
		},
	}
}

func mapEventUserUpdatedToProto(event *domain.EventUserUpdated) auctionProto.EventUserUpdated {
	return auctionProto.EventUserUpdated{
		Data: &auctionProto.User{
			Id:        event.Value.UID,
			Name:      event.Value.DisplayName,
			Email:     event.Value.Email,
			Phone:     event.Value.PhoneNumber,
			Photo:     event.Value.PhotoURL,
			CreatedAt: event.Value.CreationTimestamp,
			LastLogin: event.Value.LastLogInTimestamp,
			Disabled:  false,
		},
	}
}

func mapEventUserRevokedToProto(event *domain.EventUserRevoked) auctionProto.EventUserRevoked {
	return auctionProto.EventUserRevoked{
		Data: &auctionProto.UserId{
			Id: event.Value.UID,
		},
	}
}

func mapEventUserVerifiedToProto(event *domain.EventUserVerified) auctionProto.EventUserVerified {
	return auctionProto.EventUserVerified{
		Data: &auctionProto.UserId{
			Id: event.Value.UID,
		},
	}
}
