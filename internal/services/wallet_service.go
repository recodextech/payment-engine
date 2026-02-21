package services

import (
	"context"

	"payment-engine/internal/domain"
	"payment-engine/internal/domain/adaptors/repositories"
	"payment-engine/internal/domain/application"
	"payment-engine/internal/http/request"
	"payment-engine/pkg/errors"

	"github.com/recodextech/api-definitions/events"
	"github.com/recodextech/api-definitions/pkg/uuid"
	"github.com/recodextech/container"
)

type WalletService struct {
	accountWalletRepo repositories.AccountWalletRepository
}

func (s *WalletService) Init(c container.Container) error {
	s.accountWalletRepo = c.Resolve(application.ModuleAccountWalletRepo).(repositories.AccountWalletRepository)
	return nil
}

// CreateWallet creates a new wallet for the account
func (s *WalletService) CreateWallet(ctx context.Context, req request.CreateWalletRequest) (string, error) {
	// check if wallet already exists for the account
	walletRes, existingWallet, err := s.accountWalletRepo.GetWalletByAccountID(ctx, req.AccountID)
	if err != nil {
		return "", errors.Wrap(err, "failed to check existing wallet")
	}
	if existingWallet {
		return walletRes.Payload.ID, nil
	}
	// Create wallet event
	walletEvent := events.AccountWalletEvent{
		Payload: events.AccountWalletPayload{
			AccountID:      req.AccountID,
			Balance:        0,
			SequenceNumber: 1,
			Status:         domain.WalletStatusActive,
			ID:             uuid.New().String(),
		},
	}
	walletEvent.EventMeta = events.NewMetaContext(ctx)

	// Save wallet to repository
	walletID, err := s.accountWalletRepo.CreateWallet(ctx, req.AccountID, walletEvent)
	if err != nil {
		return "", errors.Wrap(err, "failed to create wallet")
	}

	return walletID, nil
}
