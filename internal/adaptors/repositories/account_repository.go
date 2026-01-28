package repositories

import (
	"context"
	"payment-engine/internal/domain"
	"payment-engine/internal/domain/adaptors"
	"payment-engine/internal/domain/adaptors/database"
	"payment-engine/internal/domain/application"
	"payment-engine/pkg/errors"

	"github.com/recodextech/container"
)

type AccountRepository struct {
	log       adaptors.Logger
	dbAdaptor database.MOSDB
}

func (repo *AccountRepository) Init(c container.Container) error {
	repo.log = c.Resolve(application.ModuleLogger).(adaptors.Logger).
		NewLog(adaptors.LoggerPrefixed("repositories.account-repository"))
	repo.dbAdaptor = c.Resolve(application.MoudleDBConector).(database.MOSDB)

	return nil
}

func (repo *AccountRepository) Exists(ctx context.Context, key string) (exists bool, err error) {
	exist, err := repo.dbAdaptor.GetDataRowByID(ctx, domain.ProcessTable, key)
	if err != nil {
		return false, errors.Wrap(err, `process get failed`)
	}

	return exist, nil
}
