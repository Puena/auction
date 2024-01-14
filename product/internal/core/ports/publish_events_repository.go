package ports

import (
	"context"

	"github.com/Puena/auction/product/internal/core/domain"
)

// PublishEventsRepository represent publish events repository.
type PublishEventsRepository interface {
	// Publish event product created.
	ProductCreated(ctx context.Context, userID string, event domain.EventProductCreated, replyToMsgID string) error
	// Publish event product updated.
	ProductUpdated(ctx context.Context, userID string, event domain.EventProductUpdated, replyToMsgID string) error
	// Publish event product deleted.
	ProductDeleted(ctx context.Context, userID string, event domain.EventProductDeleted, replyToMsgID string) error
	// Publish event product found.
	ProductFound(ctx context.Context, userID string, event domain.EventProductFound, replyToMsgID string) error
	// Publish event products found.
	ProductsFound(ctx context.Context, userID string, event domain.EventProductsFound, replyToMSgID string) error
	// ProductError publish product error.
	ProductError(ctx context.Context, userID string, event domain.EventProductError, replyToMsgID string) error
}
