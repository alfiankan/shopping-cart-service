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

	sql := "SELECT id, created_at FROM cart"
	row, err := repo.db.QueryContext(ctx, sql)

	for row.Next() {

		var cart cart.Cart

		if err = row.Scan(&cart.ID, &cart.CreatedAt); err != nil {
			return
		}
		carts = append(carts, cart)

	}

	return
}
func (repo *CartRepositoryPostgree) GetItems(ctx context.Context, cartID uuid.UUID, filter cart.ItemFilter) (res cart.Cart, err error) {

	args := []any{}
	sql := "SELECT * FROM cart_item"

	if filter.Name != "" || filter.Qty != 0 {
		sql += " WHERE "
	}

	if filter.Name != "" {
		args = append(args, filter.Name)
		sql += " product_name ILIKE %?% "
	}

	if filter.Name != "" {
		args = append(args, filter.Qty)
		sql += " qty = %?% "
	}

	row, err := repo.db.QueryContext(ctx, sql, args...)

	for row.Next() {

		var items cart.CartItem

		if err = row.Scan(&items.ItemID, &items.ProductCode, &items.ProductName, &items.Quantity, &res.ID); err != nil {
			return
		}
		res.Items = append(res.Items, items)

	}

	return
}

// NewCart create new cart
func (repo *CartRepositoryPostgree) NewCart(ctx context.Context) (err error) {

	sql := "INSERT INTO cart (id) VALUES (uuid_generate_v4());"
	_, err = repo.db.ExecContext(ctx, sql)

	return
}

// AddCartItem add to cart by cartid
func (repo *CartRepositoryPostgree) AddCartItem(ctx context.Context, cartID uuid.UUID, item cart.CartItem) (err error) {

	sql := "INSERT INTO cart_item (product_code, product_name, qty, cart_id) VALUES ($1, $2, $3, $4)"
	_, err = repo.db.ExecContext(ctx, sql, item.ProductCode, item.ProductName, item.Quantity, cartID)

	return
}

func (repo *CartRepositoryPostgree) DeleteCartItem(ctx context.Context, cartID uuid.UUID, productCode uuid.UUID) (err error) {

	sql := "DELETE FROM cart_item WHERE cart_id = $1 AND product_code = $2"
	_, err = repo.db.ExecContext(ctx, sql, cartID, productCode)

	return
}
