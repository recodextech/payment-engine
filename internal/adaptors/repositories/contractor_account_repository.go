package repositories

import (
	"context"

	"payment-engine/internal/domain"
	"payment-engine/internal/domain/application"
	"payment-engine/internal/domain/models"
	"payment-engine/pkg/errors"

	gopostgres "github.com/HADLakmal/go-postgres"
	"github.com/recodextech/container"
)

// ContractorAccountRepository resolves the account, user and preferences behind a contractor.
type ContractorAccountRepository struct {
	sqlRunner
}

func (r *ContractorAccountRepository) Init(c container.Container) error {
	r.db = c.Resolve(application.ModuleSQL).(gopostgres.DatabaseReporter)
	return nil
}

func (r *ContractorAccountRepository) GetContractorAccount(ctx context.Context, contractorID string) (models.ContractorAccount, error) {
	return r.getContractorAccount(ctx, `c.key = $1::uuid`, contractorID)
}

// GetContractorAccountByNumber finds a contractor by the account number in its bank transfer reference.
func (r *ContractorAccountRepository) GetContractorAccountByNumber(ctx context.Context, accountNumber int64) (models.ContractorAccount, error) {
	return r.getContractorAccount(ctx, `a.account_number = $1`, accountNumber)
}

func (r *ContractorAccountRepository) getContractorAccount(ctx context.Context, where string, arg any) (models.ContractorAccount, error) {
	var acc models.ContractorAccount
	exist, err := r.queryRow(ctx, `SELECT c.key, a.key, a.account_number, a.user_id,
			COALESCE(a.payload->'preference'->>'currency', ''),
			COALESCE(c.payload->>'contractor_name', ''), COALESCE(c.payload->>'email', ''), COALESCE(c.payload->>'phone_number', '')
		FROM `+domain.ContractorTable+` c JOIN `+domain.AccountTable+` a ON a.key = c.account_id
		WHERE `+where+` AND c.deleted = FALSE`,
		[]any{arg}, &acc.ContractorID, &acc.AccountID, &acc.AccountNumber, &acc.UserID, &acc.Currency, &acc.Name, &acc.Email, &acc.Phone)
	if err != nil {
		return acc, errors.Wrap(err, "failed to get contractor account")
	}
	if !exist {
		return acc, errors.NotFound("contractor")
	}
	return acc, nil
}

func (r *ContractorAccountRepository) ContractorBelongsToAccount(ctx context.Context, contractorID, accountID string) (bool, error) {
	var key string
	exist, err := r.queryRow(ctx, `SELECT key FROM `+domain.ContractorTable+` WHERE key = $1::uuid AND account_id = $2::uuid AND deleted = FALSE`,
		[]any{contractorID, accountID}, &key)
	if err != nil {
		return false, errors.Wrap(err, "failed to check contractor ownership")
	}
	return exist, nil
}
