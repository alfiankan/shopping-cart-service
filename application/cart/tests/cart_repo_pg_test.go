package tests

import (
	"context"
	"testing"

	"github.com/alfiankan/haioo-shoping-cart/application/cart"
	"github.com/alfiankan/haioo-shoping-cart/application/cart/repositories"
	"github.com/alfiankan/haioo-shoping-cart/config"
	"github.com/alfiankan/haioo-shoping-cart/infrastructure"
	"github.com/google/uuid"
	"github.com/jaswdr/faker"
	"github.com/stretchr/testify/assert"
)

func TestNewCart(t *testing.T) {

	/*
		ordered tests
		1. create new cart
		2. get carts
		3. add item to cart
		4. get cart items must be more than 0
		5. delete cart item
	*/
	cfg := config.Load("../../../.env")

	pgConn, err := infrastructure.NewPgConnection(cfg)
	if err != nil {
		t.Error(err)
	}
	repo := repositories.NewCartRepositoryPostgree(pgConn)
	ctx := context.Background()

	t.Run("create new cart", func(t *testing.T) {

		err := repo.NewCart(ctx)

		assert.NoError(t, err)

	})

	var currentCartID string

	t.Run("get carts", func(t *testing.T) {

		carts, err := repo.GetCarts(ctx)

		assert.NoError(t, err)

		currentCartID = carts[0].ID.String()
		assert.True(t, len(carts) > 0)
	})

	cartId, err := uuid.ParseBytes([]byte(currentCartID))
	if err != nil {
		t.Error(err)
	}

	t.Run("add item to cart", func(t *testing.T) {

		for i := 0; i < 10; i++ {
			fake := faker.New()
			errExec := repo.AddCartItem(ctx, cartId, cart.CartItem{
				ProductCode: uuid.New().String(),
				ProductName: fake.App().Name(),
				Quantity:    fake.IntBetween(1, 10),
			})
			assert.NoError(t, errExec)
		}
	})

	var currentProductCode string

	t.Run("get cart items must be more than 0", func(t *testing.T) {

		cartItems, err := repo.GetItems(ctx, cartId, cart.ItemFilter{})

		assert.NoError(t, err)
		currentProductCode = cartItems.Items[0].ProductCode
		assert.True(t, len(cartItems.Items) > 0)
	})

	if err != nil {
		t.Error(err)
	}

	t.Run("delete cart item", func(t *testing.T) {

		err := repo.DeleteCartItem(ctx, cartId, currentProductCode)
		assert.NoError(t, err)
	})

}
