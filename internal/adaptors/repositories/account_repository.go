package repositories

import (
	"context"
	"payment-engine/internal/domain"
	"payment-engine/internal/domain/adaptors"
	"payment-engine/internal/domain/adaptors/database"
	"payment-engine/internal/domain/application"
	"payment-engine/pkg/errors"

	gopostgres "github.com/HADLakmal/go-postgres"
	"github.com/recodextech/container"
)

type AccountRepository struct {
	log       adaptors.Logger
	dbAdaptor database.FixFlowDB
	sqlRunner
}

func (repo *AccountRepository) Init(c container.Container) error {
	repo.log = c.Resolve(application.ModuleLogger).(adaptors.Logger).
		NewLog(adaptors.LoggerPrefixed("repositories.account-repository"))
	repo.dbAdaptor = c.Resolve(application.MoudleDBConector).(database.FixFlowDB)
	repo.db = c.Resolve(application.ModuleSQL).(gopostgres.DatabaseReporter)

	return nil
}

func (repo *AccountRepository) Exists(ctx context.Context, key string) (exists bool, err error) {
	exist, err := repo.dbAdaptor.GetDataRowByID(ctx, domain.AccountTable, key)
	if err != nil {
		return false, errors.Wrap(err, `account get failed`)
	}

	return exist, nil
}

// GetAccountCurrency returns the account's preferred currency, empty when not set.
func (repo *AccountRepository) GetAccountCurrency(ctx context.Context, accountID string) (string, error) {
	var currency string
	exist, err := repo.queryRow(ctx, `SELECT COALESCE(`+columnParamPayload+`->'preference'->>'currency', '') FROM `+domain.AccountTable+
		` WHERE `+columnParamKey+` = $1::uuid`, []any{accountID}, &currency)
	if err != nil {
		return "", errors.Wrap(err, "failed to get account currency")
	}
	if !exist {
		return "", errors.NotFound("account")
	}
	return currency, nil
}
