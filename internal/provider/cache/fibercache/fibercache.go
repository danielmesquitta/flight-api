package fibercache

import (
	"context"
	"maps"
	"slices"
	"time"

	"github.com/danielmesquitta/flight-api/internal/provider/cache"
	"github.com/gofiber/fiber/v2"
)

type FiberCache struct {
	c          cache.Cache
	storedKeys map[string]struct{}
}

func NewFiberCache(c cache.Cache) *FiberCache {
	return &FiberCache{
		c:          c,
		storedKeys: map[string]struct{}{},
	}
}

func (f *FiberCache) Get(key string) ([]byte, error) {
	var val []byte
	ok, err := f.c.Scan(context.Background(), key, &val)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}
	return val, nil
}

func (f *FiberCache) Set(key string, val []byte, exp time.Duration) error {
	err := f.c.Set(context.Background(), key, val, exp)
	if err != nil {
		return err
	}
	f.storedKeys[key] = struct{}{}
	return nil
}

func (f *FiberCache) Delete(key string) error {
	err := f.c.Delete(context.Background(), key)
	if err != nil {
		return err
	}
	delete(f.storedKeys, key)
	return nil
}

func (f *FiberCache) Reset() error {
	keys := slices.Collect(maps.Keys(f.storedKeys))
	err := f.c.Delete(context.Background(), keys...)
	if err != nil {
		return err
	}
	f.storedKeys = map[string]struct{}{}
	return nil
}

func (f *FiberCache) Close() error {
	return nil
}

var _ fiber.Storage = (*FiberCache)(nil)
