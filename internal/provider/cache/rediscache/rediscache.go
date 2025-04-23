package rediscache

import (
	"context"
	"encoding/json"
	"strconv"
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
	key string,
	value any,
) (bool, error) {
	strCmd := r.c.Get(ctx, key)
	if strCmd.Err() == redis.Nil {
		return false, nil
	}
	if err := strCmd.Err(); err != nil {
		return false, err
	}

	raw, err := strCmd.Result()
	if err != nil {
		return false, err
	}

	switch v := value.(type) {
	case *string:
		*v = raw
	case *[]byte:
		*v = []byte(raw)
	case *int:
		i, err := strconv.Atoi(raw)
		if err != nil {
			return false, err
		}
		*v = i
	case *int64:
		i, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			return false, err
		}
		*v = i
	case *float64:
		f, err := strconv.ParseFloat(raw, 64)
		if err != nil {
			return false, err
		}
		*v = f
	case *bool:
		b, err := strconv.ParseBool(raw)
		if err != nil {
			return false, err
		}
		*v = b
	default:
		if err := json.Unmarshal([]byte(raw), value); err != nil {
			return false, err
		}
	}

	return true, nil
}

func (r *RedisCache) Set(
	ctx context.Context,
	key string,
	value any,
	expiration time.Duration,
) error {
	var data any

	switch v := value.(type) {
	case string, []byte, int, int64, float64, bool:
		data = v
	default:
		b, err := json.Marshal(value)
		if err != nil {
			return err
		}
		data = b
	}

	return r.c.Set(ctx, key, data, expiration).Err()
}

func (r *RedisCache) Delete(
	ctx context.Context,
	keys ...string,
) error {
	ks := make([]string, len(keys))
	copy(ks, keys)
	return r.c.Del(ctx, ks...).Err()
}

var _ cache.Cache = (*RedisCache)(nil)
