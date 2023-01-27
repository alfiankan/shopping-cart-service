package repositories

import (
	"context"
	"encoding/json"
	"time"

	"github.com/alfiankan/haioo-shoping-cart/application/cart"
	"github.com/go-redis/redis/v9"
)

type CartRepositoryCacheRedis struct {
	cache redis.UniversalClient
	ttl   time.Duration
}

func NewCartRepositoryCacheRedis(cache redis.UniversalClient, ttl time.Duration) cart.ICartCacheRepository {
	return &CartRepositoryCacheRedis{cache, ttl}
}

func (repo *CartRepositoryCacheRedis) Save(ctx context.Context, cart cart.Cart, cacheKey string) (err error) {

	cacheData, err := json.Marshal(cart)

	// save as string
	stat := repo.cache.Set(ctx, cacheKey, string(cacheData), repo.ttl)

	if stat.Err() != nil {
		err = stat.Err()
		return
	}

	return
}
func (repo *CartRepositoryCacheRedis) Get(ctx context.Context, cartID string) (cart cart.Cart, err error) {

	stat := repo.cache.Get(ctx, cartID)
	if stat.Err() != nil {
		err = stat.Err()
		return
	}

	var cachedData string

	if err = stat.Scan(&cachedData); err != nil {
		return
	}

	// unmarshal
	err = json.Unmarshal([]byte(cachedData), &cart)

	return
}
func (repo *CartRepositoryCacheRedis) InvalidateByCartID(ctx context.Context, cartID string) (err error) {

	if err = repo.cache.Del(ctx, cartID).Err(); err != nil {
		return
	}

	return
}
