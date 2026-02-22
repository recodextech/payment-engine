package handlers

import (
	"context"
	"payment-engine/internal/domain/application"
	"payment-engine/internal/domain/service"
	"payment-engine/internal/http/request"
	"payment-engine/internal/http/responses"
	"payment-engine/pkg/errors"

	"github.com/recodextech/container"
	"github.com/recodextech/krouter"
)

type GetWalletsHandler struct {
	walletService service.WalletService
}

func (h *GetWalletsHandler) Init(c container.Container) error {
	h.walletService = c.Resolve(application.ModuleWalletService).(service.WalletService)
	return nil
}

func (h *GetWalletsHandler) Handle(ctx context.Context, payload krouter.HttpPayload) (interface{}, error) {
	// Get account_id from header
	accountID := payload.Header(request.HeaderAccountID.String())

	// Call wallet service to get wallets
	wallets, err := h.walletService.GetWallets(ctx, accountID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get wallets")
	}

	// Return response with wallets
	return responses.NewGetWalletsResponse(wallets), nil
}
