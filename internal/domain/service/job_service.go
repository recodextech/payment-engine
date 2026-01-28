package service

import (
	"context"

	"payment-engine/internal/domain/adaptors/repositories"
	"payment-engine/internal/http/request"
)

type JobService interface {
	GetJobsByWorkerAndContractor(ctx context.Context, req request.GetJobsByWorker) ([]repositories.JobWithDetails, error)
}
