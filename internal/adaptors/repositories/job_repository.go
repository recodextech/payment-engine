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

type JobRepository struct {
	dbAdaptor database.FixFlowDB
	log       adaptors.Logger
}

func (r *JobRepository) Init(c container.Container) error {
	r.dbAdaptor = c.Resolve(application.MoudleDBConector).(database.FixFlowDB)
	r.log = c.Resolve(application.ModuleLogger).(adaptors.Logger).NewLog(adaptors.LoggerPrefixed(`repo.payment`))
	return nil
}

// job exists
func (r *JobRepository) Exists(ctx context.Context, jobID string) (bool, error) {
	exist, err := r.dbAdaptor.GetDataRowByID(ctx, domain.JobTable, jobID)
	if err != nil {
		return false, errors.Wrap(err, "failed to check job existence")
	}

	return exist, nil
}
