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

type CreateWalletHandler struct {
	walletService service.WalletService
}

func (h *CreateWalletHandler) Init(c container.Container) error {
	h.walletService = c.Resolve(application.ModuleWalletService).(service.WalletService)
	return nil
}

func (h *CreateWalletHandler) Handle(ctx context.Context, payload krouter.HttpPayload) (interface{}, error) {
	// Get account_id from header
	accountID := payload.Header(request.HeaderAccountID.String())

	if accountID == "" {
		return nil, errors.New("account-id header is required")
	}

	// Create wallet request
	req := request.CreateWalletRequest{
		AccountID: accountID,
	}
	// Call wallet service to create wallet
	walletID, err := h.walletService.CreateWallet(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create wallet")
	}

	// Return response with wallet ID
	return responses.CreateWalletResponse{
		ID: walletID,
	}, nil
}
