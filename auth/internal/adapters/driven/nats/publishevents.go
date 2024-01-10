package nats

import (
	"context"

	"github.com/Puena/auction-house-auth/internal/core/domain"
	"github.com/Puena/auction-house-nats/publish"
	"google.golang.org/protobuf/proto"
)

const (
	eventSubjectUserRevoked  = "response.user.revoked"
	eventSubjectUserUpdated  = "response.user.updated"
	eventSubjectUserVerified = "response.user.verified"
	eventSubjectUserData     = "response.user.data"
	eventSubjectUserError    = "error.user"
)

type publishEvents struct {
	jetstream publish.JetStreamPublisher
}

func NewPublishEvents(jetstream publish.JetStreamPublisher) *publishEvents {
	return &publishEvents{
		jetstream: jetstream,
	}
}

func (p *publishEvents) PublishEventUserRevoked(ctx context.Context, event *domain.EventUserRevoked) error {
	msg := mapEventUserRevokedToProto(event)
	data, err := proto.Marshal(&msg)
	if err != nil {
		return err
	}

	_, err = p.jetstream.Publish(ctx, eventSubjectUserRevoked, data, event.Key)
	if err != nil {
		return err
	}

	return nil
}

func (p *publishEvents) PublishEventUserUpdated(ctx context.Context, event *domain.EventUserUpdated) error {
	msg := mapEventUserUpdatedToProto(event)
	data, err := proto.Marshal(&msg)
	if err != nil {
		return err
	}

	_, err = p.jetstream.Publish(ctx, eventSubjectUserUpdated, data, event.Key)
	if err != nil {
		return err
	}

	return nil
}

func (p *publishEvents) PublishEventUserVerified(ctx context.Context, event *domain.EventUserVerified) error {
	msg := mapEventUserVerifiedToProto(event)
	data, err := proto.Marshal(&msg)
	if err != nil {
		return err
	}

	_, err = p.jetstream.Publish(ctx, eventSubjectUserVerified, data, event.Key)
	if err != nil {
		return err
	}

	return nil
}

func (p *publishEvents) PublishEventUserData(ctx context.Context, event *domain.EventUserData) error {
	msg := mapEventUserDataToProto(event)
	data, err := proto.Marshal(&msg)
	if err != nil {
		return err
	}

	_, err = p.jetstream.Publish(ctx, eventSubjectUserData, data, event.Key)
	if err != nil {
		return err
	}

	return nil
}

func (p *publishEvents) PublishEventUserError(ctx context.Context, event *domain.EventUserError) error {
	data, err := proto.Marshal(nil)
	if err != nil {
		return err
	}

	_, err = p.jetstream.Publish(ctx, eventSubjectUserError, data, event.Key)
	if err != nil {
		return err
	}

	return nil
}
