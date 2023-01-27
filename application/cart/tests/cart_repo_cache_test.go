package tests

import (
	"context"
	"testing"
	"time"

	"github.com/alfiankan/haioo-shoping-cart/application/cart"
	"github.com/alfiankan/haioo-shoping-cart/application/cart/repositories"
	"github.com/alfiankan/haioo-shoping-cart/config"
	"github.com/alfiankan/haioo-shoping-cart/infrastructure"
	"github.com/google/uuid"
	"github.com/jaswdr/faker"
	"github.com/stretchr/testify/assert"
)

func TestCartRepoCahche(t *testing.T) {

	cfg := config.Load("../../../.env")

	redConn, err := infrastructure.NewRedisConnection(cfg)
	if err != nil {
		t.Error(err)
	}
	repo := repositories.NewCartRepositoryCacheRedis(redConn, time.Second*300)
	ctx := context.Background()
	newCart := cart.Cart{ID: uuid.New(), Items: []cart.CartItem{}}

	for i := 0; i < 10; i++ {
		fake := faker.New()
		newCart.Items = append(newCart.Items, cart.CartItem{
			ProductCode: uuid.New().String(),
			ProductName: fake.App().Name(),
			Quantity:    fake.IntBetween(1, 10),
		})
	}

	t.Run("add item", func(t *testing.T) {
		err := repo.Save(ctx, newCart)
		assert.NoError(t, err)
	})

	t.Run("get item", func(t *testing.T) {
		_, err := repo.Get(ctx, newCart.ID.String())
		assert.NoError(t, err)
	})

	t.Run("invalidate item", func(t *testing.T) {

		err := repo.InvalidateByCartID(ctx, newCart.ID.String())
		assert.NoError(t, err)
	})

}
