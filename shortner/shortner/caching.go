package shortner

import (
	"context"

	"github.com/go-redis/redis/v9"
)

type ShortnerCache interface {
	SetUrl(ctx context.Context, shortkey, url string) error
	GetUrl(ctx context.Context, shortkey string) (string, error)
}

type shortnerCache struct {
	redis *redis.Client
}

func (sc *shortnerCache) SetUrl(ctx context.Context, shortkey, url string) error {
	key := "USSO:" + shortkey
	return sc.redis.Set(ctx, key, url, 0).Err()
}

func (sc *shortnerCache) GetUrl(ctx context.Context, shortkey string) (string, error) {
	key := "USSO:" + shortkey
	return sc.redis.Get(ctx, key).Result()
}

func NewShortnerCache(redis *redis.Client) ShortnerCache {
	return &shortnerCache{
		redis: redis,
	}
}
