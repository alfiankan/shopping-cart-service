package tests

import (
	"context"
	"testing"

	"github.com/alfiankan/haioo-shoping-cart/application/cart/repositories"
	"github.com/alfiankan/haioo-shoping-cart/config"
	"github.com/alfiankan/haioo-shoping-cart/infrastructure"
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

	t.Run("get cart", func(t *testing.T) {
		t.Error(nil)
	})

	t.Run("add item to cart", func(t *testing.T) {
		t.Error(nil)

	})

	t.Run("get cart items must be more than 0", func(t *testing.T) {
		t.Error(nil)

	})

	t.Run("delete cart item", func(t *testing.T) {
		t.Error(nil)

	})

}
