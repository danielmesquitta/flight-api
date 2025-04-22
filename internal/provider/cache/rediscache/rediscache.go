package rediscache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/danielmesquitta/flight-api/internal/config/env"
	"github.com/danielmesquitta/flight-api/internal/provider/cache"
)

type RedisCache struct {
	c *redis.Client
}

func NewRedisCache(
	e *env.Env,
) *RedisCache {
	opts, err := redis.ParseURL(e.RedisDatabaseURL)
	if err != nil {
		panic(err)
	}

	client := redis.NewClient(opts)

	status := client.Ping(context.Background())
	if status.Err() != nil {
		panic(status.Err())
	}

	return &RedisCache{
		c: client,
	}
}

func (r *RedisCache) Scan(
	ctx context.Context,
	key cache.Key,
	value any,
) (bool, error) {
	strCmd := r.c.Get(ctx, key)
	if strCmd.Err() == redis.Nil {
		return false, nil
	}
	if strCmd.Err() != nil {
		return false, strCmd.Err()
	}

	if err := strCmd.Scan(value); err != nil {
		return false, err
	}

	return true, nil
}

func (r *RedisCache) Set(
	ctx context.Context,
	key cache.Key,
	value any,
	expiration time.Duration,
) error {
	return r.c.Set(ctx, key, value, expiration).Err()
}

func (r *RedisCache) Delete(
	ctx context.Context,
	keys ...cache.Key,
) error {
	ks := make([]string, len(keys))
	copy(ks, keys)
	return r.c.Del(ctx, ks...).Err()
}

var _ cache.Cache = (*RedisCache)(nil)
