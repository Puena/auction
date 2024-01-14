package nats

import (
	"github.com/Puena/auction/pbgo/auction"
	"github.com/Puena/auction/product/internal/core/domain"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func mapProductToProto(eventProduct domain.Product) *auction.Product {
	p := &auction.Product{
		Id:          eventProduct.ID,
		Name:        eventProduct.Name,
		Description: eventProduct.Description,
		Media:       eventProduct.Media,
		CreatedAt:   timestamppb.New(eventProduct.CreatedAt),
		CreatedBy:   eventProduct.CreatedBy,
	}
	if eventProduct.UpdatedAt != nil {
		p.UpdatedAt = timestamppb.New(*eventProduct.UpdatedAt)
	}

	return p
}
