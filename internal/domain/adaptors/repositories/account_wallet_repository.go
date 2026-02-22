package repositories

import (
	"context"
	"payment-engine/internal/domain"

	"github.com/recodextech/api-definitions/events"
)

// AccountWalletRepository defines the interface for account wallet data operations
type AccountWalletRepository interface {
	CreateWallet(ctx context.Context, accountID string, wallet events.AccountWalletEvent) (string, error)
	CreateInternalWallets(ctx context.Context, accountID string, wallets []events.AccountWalletEvent) ([]string, error)
	GetWalletByID(ctx context.Context, walletID string) (walletRes events.AccountWalletEvent, err error)
	GetWalletsByAccountID(ctx context.Context, accountID string) (walletRes []events.AccountWalletEvent, exist bool, err error)
	GetWalletByType(ctx context.Context, accountID string, walletType domain.WalletType) (walletRes events.AccountWalletEvent, exist bool, err error)
	UpdateWalletBalance(ctx context.Context, walletID string, balance float64) error
}
