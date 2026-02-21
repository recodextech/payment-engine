package service

import (
	"context"
	"payment-engine/internal/http/request"
)

type WalletService interface {
	CreateWallet(ctx context.Context, req request.CreateWalletRequest) (string, error)
}
