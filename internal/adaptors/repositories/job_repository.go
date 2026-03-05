package repositories

import (
	"context"
	"encoding/json"
	"payment-engine/internal/domain"
	"payment-engine/internal/domain/adaptors"
	"payment-engine/internal/domain/adaptors/database"
	"payment-engine/internal/domain/application"
	"payment-engine/pkg/errors"

	"github.com/recodextech/api-definitions/events"
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

// GetJobByID retrieves a job by its ID
func (r *JobRepository) GetJobByID(ctx context.Context, jobID string) (events.JobPayload, error) {
	columns := []string{columnParamPayload}
	whereClause := `key = $1 AND deleted = false`

	result, err := r.dbAdaptor.GetDataRowWithResult(ctx, domain.JobTable, columns, whereClause, []interface{}{jobID})
	if err != nil {
		return events.JobPayload{}, errors.Wrap(err, "failed to get job by ID")
	}

	var jobPayloadBytes []byte
	exist, err := result.Scan(
		&jobPayloadBytes)
	if err != nil {
		return events.JobPayload{}, errors.Wrap(err, "failed to scan job row")
	}
	if !exist {
		return events.JobPayload{}, errors.RepositoryDataNotExistError{}
	}
	jobRes := events.JobPayload{}
	err = json.Unmarshal(jobPayloadBytes, &jobRes)
	if err != nil {
		return events.JobPayload{}, errors.Wrap(err, "failed to unmarshal job payload")
	}

	return jobRes, nil
}
