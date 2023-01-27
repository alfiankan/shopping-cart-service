package cart

import (
	"context"

	"github.com/google/uuid"
)

// ICartUseCase main application business logic hold cart usecases
type ICartUseCase interface {
	CreateCart(ctx context.Context) (err error)
	GetCarts(ctx context.Context) (carts []Cart, err error)
	GetCartItems(ctx context.Context, cartID uuid.UUID, filter ItemFilter) (cart Cart, err error)

	AddToCart(ctx context.Context, cart Cart) (err error)
	DeleteCartItem(ctx context.Context, cartID uuid.UUID, productCode string) (err error)
}

// ICartRepository persistence db storage store cart
type ICartRepository interface {
	GetCarts(ctx context.Context) (carts []Cart, err error)
	GetItems(ctx context.Context, cartID uuid.UUID, filter ItemFilter) (carts Cart, err error)
	GetItemByProductID(ctx context.Context, cartID uuid.UUID, productCode string) (item CartItem, err error)

	UpdateQtyByProductID(ctx context.Context, cartID uuid.UUID, productCode string, qty int) (ok bool)
	NewCart(ctx context.Context) (err error)
	AddCartItem(ctx context.Context, cartID uuid.UUID, item CartItem) (err error)
	DeleteCartItem(ctx context.Context, cartID uuid.UUID, productCode string) (err error)
}

// ICartCacheRepository cache cart items
type ICartCacheRepository interface {
	Save(ctx context.Context, cart Cart, cacheKey string) (err error)
	Get(ctx context.Context, cartID string) (cart Cart, err error)
	InvalidateByCartID(ctx context.Context, cartID string) (err error)
}
