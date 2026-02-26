package repositories

import "context"

type JobRepository interface {
	Exists(ctx context.Context, jobID string) (bool, error)
}
