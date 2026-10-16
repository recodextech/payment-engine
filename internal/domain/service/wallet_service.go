package service

import (
	"context"
	"payment-engine/internal/domain/models"
	"payment-engine/internal/http/request"

	"github.com/recodextech/api-definitions/events"
)

type WalletService interface {
	CreateWallet(ctx context.Context, req request.CreateWalletRequest) (string, error)
	CreateInternalWallet(ctx context.Context) ([]string, error)
	GetWallets(ctx context.Context, accountID string) ([]events.AccountWalletEvent, error)
	GetWalletTotals(ctx context.Context, walletID string) (models.WalletTotals, error)
	GetWalletTransactions(ctx context.Context, accountID, walletID string) ([]models.WalletTransaction, error)
}
