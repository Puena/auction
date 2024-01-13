package nats

import (
	"github.com/Puena/auction/pbgo/auction"
	"github.com/Puena/auction/product/internal/core/domain"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func eventProductToProtoProduct(eventProduct domain.Product) *auction.Product {
	return &auction.Product{
		Id:          eventProduct.ID,
		Name:        eventProduct.Name,
		Description: eventProduct.Description,
		Media:       eventProduct.Media,
		CreatedAt:   timestamppb.New(eventProduct.CreatedAt),
		CreatedBy:   eventProduct.CreatedBy,
	}
}
