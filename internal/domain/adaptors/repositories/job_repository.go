package repositories

import (
	"context"

	"github.com/recodextech/api-definitions/events"
)

type JobRepository interface {
	Exists(ctx context.Context, jobID string) (bool, error)
	GetJobByID(ctx context.Context, jobID string) (events.JobPayload, error)
}
