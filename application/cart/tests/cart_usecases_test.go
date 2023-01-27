package tests

import (
	"context"
	"testing"
	"time"

	"github.com/alfiankan/haioo-shoping-cart/application/cart"
	"github.com/alfiankan/haioo-shoping-cart/application/cart/repositories"
	"github.com/alfiankan/haioo-shoping-cart/application/cart/usecases"
	"github.com/alfiankan/haioo-shoping-cart/config"
	"github.com/alfiankan/haioo-shoping-cart/infrastructure"
	"github.com/google/uuid"
	"github.com/jaswdr/faker"
	"github.com/stretchr/testify/assert"
)

func TestCartUseCases(t *testing.T) {

	cfg := config.Load("../../../.env")

	redisConn, _ := infrastructure.NewRedisConnection(cfg)
	pgConn, _ := infrastructure.NewPgConnection(cfg)

	persistRepo := repositories.NewCartRepositoryPostgree(pgConn)
	cacheRepo := repositories.NewCartRepositoryCacheRedis(redisConn, 5*time.Minute)

	uc := usecases.NewCartApplication(persistRepo, cacheRepo)
	ctx := context.Background()

	t.Run("create cart", func(t *testing.T) {
		err := uc.CreateCart(ctx)
		assert.NoError(t, err)
	})
	var cartID uuid.UUID
	t.Run("get carts", func(t *testing.T) {

		res, err := uc.GetCarts(ctx)
		cartID = res[0].ID
		assert.True(t, len(res) > 0)
		assert.NoError(t, err)

	})

	t.Run("add to cart", func(t *testing.T) {

		for i := 0; i < 10; i++ {
			fake := faker.New()
			errExec := uc.AddToCart(ctx, cart.Cart{
				ID: cartID,
				Items: []cart.CartItem{{
					ProductCode: uuid.New().String(),
					ProductName: fake.App().Name(),
					Quantity:    fake.IntBetween(1, 10),
				}},
			})
			assert.NoError(t, errExec)
		}

	})

	var productCode string
	t.Run("get cart items", func(t *testing.T) {

		cart, err := uc.GetCartItems(ctx, cartID, cart.ItemFilter{})

		assert.True(t, len(cart.Items) > 0)
		productCode = cart.Items[0].ProductCode
		assert.NoError(t, err)

	})

	t.Run("delete cart item", func(t *testing.T) {

		err := uc.DeleteCartItem(ctx, cartID, productCode)

		assert.NoError(t, err)
	})

}
