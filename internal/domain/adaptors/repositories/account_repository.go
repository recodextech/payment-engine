package repositories

import (
	"context"
)

type AccountRepository interface {
	Exists(ctx context.Context, key string) (exists bool, err error)
	// GetAccountCurrency returns the account's preferred currency, empty when not set.
	GetAccountCurrency(ctx context.Context, accountID string) (string, error)
}
