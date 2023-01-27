package cart

import (
	"context"

	"github.com/google/uuid"
)

// ICartUseCase main application business logic hold cart usecases
type ICartUseCase interface {
	CreateCart(ctx context.Context) uuid.UUID
	GetCarts(ctx context.Context) (carts []Cart, err error)
	GetCartItems(ctx context.Context, cartID uuid.UUID, filter ItemFilter) (cart Cart, err error)

	AddToCart(ctx context.Context, cart Cart) (err error)
	DeleteCartItem(ctx context.Context, cartID uuid.UUID, productCode uuid.UUID) (err error)
}

// ICartRepository persistence db storage store cart
type ICartRepository interface {
	GetCarts(ctx context.Context) (carts []Cart, err error)
	GetItems(ctx context.Context, cartID uuid.UUID, filter ItemFilter) (carts Cart, err error)

	NewCart(ctx context.Context) (err error)
	AddCartItem(ctx context.Context, cartID uuid.UUID, item CartItem) (err error)
	DeleteCartItem(ctx context.Context, cartID uuid.UUID, productCode uuid.UUID) (err error)
}

// ICartCacheRepository cache cart items
type ICartCacheRepository interface {
	Save(ctx context.Context, cart Cart) (err error)
	Get(ctx context.Context, cartID string) (cart Cart, err error)
	InvalidateByCartID(ctx context.Context, cartID string) (err error)
}
