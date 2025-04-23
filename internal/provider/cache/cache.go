package cache

import (
	"context"
	"time"
)

type Cache interface {
	Scan(ctx context.Context, key string, value any) (ok bool, err error)

	Set(
		ctx context.Context,
		key string,
		value any,
		expiration time.Duration,
	) error

	Delete(
		ctx context.Context,
		keys ...string,
	) error
}
