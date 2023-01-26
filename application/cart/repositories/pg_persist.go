package repositories

import (
	"context"
	"database/sql"

	"github.com/alfiankan/haioo-shoping-cart/application/cart"
	"github.com/google/uuid"
)

type CartRepositoryPostgree struct {
	db *sql.DB
}

func NewCartRepositoryPostgree(db *sql.DB) cart.ICartRepository {
	return &CartRepositoryPostgree{db}
}

func (repo *CartRepositoryPostgree) GetCarts(ctx context.Context) (carts []cart.Cart, err error) {
	return
}
func (repo *CartRepositoryPostgree) GetItems(ctx context.Context, filter cart.ItemFilter) (carts []cart.Cart, err error) {
	return
}
func (repo *CartRepositoryPostgree) NewCart(ctx context.Context, cart cart.Cart) (err error) {
	return
}
func (repo *CartRepositoryPostgree) AddCartItem(ctx context.Context, item cart.CartItem) (err error) {
	return
}
func (repo *CartRepositoryPostgree) DeleteCartItem(ctx context.Context, cartID uuid.UUID, productCode uuid.UUID) (err error) {
	return
}
