package repositories

import (
	"context"

	"github.com/recodextech/api-definitions/events"
)

// AccountWalletRepository defines the interface for account wallet data operations
type AccountWalletRepository interface {
	CreateWallet(ctx context.Context, accountID string, wallet events.AccountWalletEvent) (string, error)
	GetWalletByID(ctx context.Context, walletID string) (walletRes events.AccountWalletEvent, err error)
	GetWalletByAccountID(ctx context.Context, accountID string) (walletRes events.AccountWalletEvent, exist bool, err error)
	UpdateWalletBalance(ctx context.Context, walletID string, balance float64) error
}
