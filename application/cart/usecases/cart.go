package usecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/alfiankan/haioo-shoping-cart/application/cart"
	"github.com/alfiankan/haioo-shoping-cart/common"
	"github.com/google/uuid"
)

type CartApplication struct {
	persistRepo cart.ICartRepository
	cacheRepo   cart.ICartCacheRepository
}

func NewCartApplication(persistRepo cart.ICartRepository, cacheRepo cart.ICartCacheRepository) cart.ICartUseCase {
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

	cacheKey := fmt.Sprintf("%s#%s#%d", cartID, filter.Name, filter.Qty)

	cart, err = uc.cacheRepo.Get(ctx, cacheKey)
	if err != nil {
		// get from persistence
		cart, err = uc.persistRepo.GetItems(ctx, cartID, filter)

		// save cache
		if errCaching := uc.cacheRepo.Save(ctx, cart, cacheKey); errCaching != nil {
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

		lastItem, errExist := uc.persistRepo.GetItemByProductID(ctx, cart.ID, cart.Items[0].ProductCode)
		if errExist != nil {
			if err = uc.persistRepo.AddCartItem(ctx, cart.ID, cart.Items[0]); err != nil {
				return
			}
		}

		// update

		if !uc.persistRepo.UpdateQtyByProductID(ctx, cart.ID, cart.Items[0].ProductCode, lastItem.Quantity+cart.Items[0].Quantity) {
			err = errors.New("Cant update cart")
			return
		}

	}
	return

}

func (uc *CartApplication) DeleteCartItem(ctx context.Context, cartID uuid.UUID, productCode string) (err error) {

	if err = uc.persistRepo.DeleteCartItem(ctx, cartID, productCode); err != nil {
		return
	}

	// invalidate cacheRepo

	err = uc.cacheRepo.InvalidateAll(ctx)

	return
}
