package usecases

import (
	"context"

	"github.com/alfiankan/haioo-shoping-cart/application/cart"
	"github.com/alfiankan/haioo-shoping-cart/common"
	"github.com/google/uuid"
)

type CartApplication struct {
	persistRepo cart.ICartRepository
	cacheRepo   cart.ICartCacheRepository
}

func NewCartApplicationI(persistRepo cart.ICartRepository, cacheRepo cart.ICartCacheRepository) cart.ICartUseCase {
	return &CartApplication{persistRepo, cacheRepo}
}

func (uc *CartApplication) CreateCart(ctx context.Context) (err error) {
	err = uc.persistRepo.NewCart(ctx)
	return
}

func (uc *CartApplication) GetCarts(ctx context.Context) (carts []cart.Cart, err error) {

	carts, err = uc.persistRepo.GetCarts(ctx)

	return
}

func (uc *CartApplication) GetCartItems(ctx context.Context, cartID uuid.UUID, filter cart.ItemFilter) (cart cart.Cart, err error) {

	cart, err = uc.cacheRepo.Get(ctx, cartID.String())
	if err != nil {
		// get from persistence
		cart, err = uc.persistRepo.GetItems(ctx, cartID, filter)

		// save cache
		if errCaching := uc.cacheRepo.Save(ctx, cart); errCaching != nil {
			common.Log(common.LOG_LEVEL_ERROR, errCaching.Error())
		}
	}

	return
}

func (uc *CartApplication) AddToCart(ctx context.Context, cart cart.Cart) (err error) {
	// invalidate cache
	if err = uc.cacheRepo.InvalidateByCartID(ctx, cart.ID.String()); err != nil {
		return
	}

	// add to AddToCart
	if len(cart.Items) > 0 {
		if err = uc.persistRepo.AddCartItem(ctx, cart.ID, cart.Items[0]); err != nil {
			return
		}
	}
	return

}

func (uc *CartApplication) DeleteCartItem(ctx context.Context, cartID uuid.UUID, productCode uuid.UUID) (err error) {

	return
}
