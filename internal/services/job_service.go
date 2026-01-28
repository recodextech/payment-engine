package services

import (
	"context"

	"payment-engine/internal/domain/adaptors/repositories"
	"payment-engine/internal/domain/application"
	"payment-engine/internal/http/request"
	"payment-engine/pkg/errors"

	"github.com/recodextech/container"
)

type JobService struct {
	jobRepo repositories.JobRepository
}

func (s *JobService) Init(c container.Container) error {
	s.jobRepo = c.Resolve(application.ModuleJobRepo).(repositories.JobRepository)
	return nil
}

func (s *JobService) GetJobsByWorkerAndContractor(ctx context.Context, req request.GetJobsByWorker) ([]repositories.JobWithDetails, error) {
	if req.WorkerID == "" {
		return nil, errors.New("worker_id is required")
	}
	if req.AccountID == "" {
		return nil, errors.New("account_id is required")
	}

	jobs, err := s.jobRepo.GetJobsByWorker(ctx, req.WorkerID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get jobs by worker")
	}

	return jobs, nil
}
