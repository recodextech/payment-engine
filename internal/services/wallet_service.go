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
	walletRes, existingWallet, err := s.accountWalletRepo.GetWalletByType(ctx, req.AccountID, domain.WalletType(req.WalletType))
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
			Type:           domain.WalletType(req.WalletType).String(),
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

func (s *WalletService) CreateInternalWallet(ctx context.Context) ([]string, error) {
	accountID := ctx.Value(domain.ContextKeyAccountID.String()).(uuid.UUID).String()
	// check if wallet already exists for the account
	walletPoint, existingWallet, err := s.accountWalletRepo.GetWalletByType(ctx, accountID, domain.WalletPoints)
	if !errors.Is(err, errors.RepositoryDataNotExistError{}) {
		return nil, errors.Wrap(err, "failed to check existing wallet")
	}
	if existingWallet {
		return []string{walletPoint.Payload.ID}, nil
	}
	wallets := make([]events.AccountWalletEvent, 2)
	for i := 0; i < 2; i++ {
		wallets[i] = events.AccountWalletEvent{
			Payload: events.AccountWalletPayload{
				AccountID:      accountID,
				Balance:        0,
				SequenceNumber: 1,
				Status:         domain.WalletStatusActive,
				ID:             uuid.New().String(),
			},
			EventMeta: events.NewMetaContext(ctx),
		}
	}
	wallets[0].Payload.Type = domain.WalletPoints.String()
	wallets[1].Payload.Type = domain.WalletCash.String()
	// Save internal wallets to repository
	walletIDs, err := s.accountWalletRepo.CreateInternalWallets(ctx, accountID, wallets)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create internal wallets")
	}

	return walletIDs, nil
}

// GetWallets retrieves all wallets for an account
func (s *WalletService) GetWallets(ctx context.Context, accountID string) ([]events.AccountWalletEvent, error) {
	wallets, exist, err := s.accountWalletRepo.GetWalletsByAccountID(ctx, accountID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get wallets")
	}
	if !exist {
		return []events.AccountWalletEvent{}, nil
	}
	return wallets, nil
}
