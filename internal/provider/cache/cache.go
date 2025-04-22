package cache

import (
	"context"
	"time"
)

type Cache interface {
	Scan(ctx context.Context, key Key, value any) (ok bool, err error)

	Set(
		ctx context.Context,
		key Key,
		value any,
		expiration time.Duration,
	) error

	Delete(
		ctx context.Context,
		keys ...Key,
	) error
}

type Key = string

const (
	KeySyncTransactionsOffset Key = "sync_transactions_offset"
	KeySyncBalancesOffset     Key = "sync_balances_offset"
)
