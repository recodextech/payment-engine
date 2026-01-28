package handlers

import (
	"context"

	"payment-engine/internal/domain/application"
	"payment-engine/internal/domain/service"
	"payment-engine/internal/http/request"
	"payment-engine/internal/http/responses"
	"payment-engine/pkg/errors"

	"github.com/recodextech/container"
	"github.com/recodextech/krouter"
)

const GetJobsByWorker = "job.get-by-worker"

type GetJobsByWorkerHandler struct {
	jobService service.JobService
}

func (h *GetJobsByWorkerHandler) Init(c container.Container) error {
	h.jobService = c.Resolve(application.ModuleJobService).(service.JobService)
	return nil
}

func (h *GetJobsByWorkerHandler) Handle(ctx context.Context, payload krouter.HttpPayload) (interface{}, error) {
	// Get worker_id from path parameter
	workerID := payload.Param(request.PathParamWorkerID.String())
	accountId := payload.Header(request.HeaderAccountID.String())

	if workerID == "" {
		return nil, errors.New("worker_id path parameter is required")
	}

	req := request.GetJobsByWorker{
		WorkerID:  workerID,
		AccountID: accountId,
	}

	jobs, err := h.jobService.GetJobsByWorkerAndContractor(ctx, req)
	if err != nil {
		return nil, err
	}

	// Convert to response format
	jobDetails := make([]responses.JobDetails, 0, len(jobs))
	for _, job := range jobs {
		jobDetails = append(jobDetails, responses.JobDetails{
			WorkerID:           job.WorkerID,
			WorkerStatus:       job.WorkerStatus,
			WorkerCategories:   job.WorkerCategories,
			JobID:              job.JobID,
			JobStatus:          job.JobStatus,
			JobLatitude:        job.JobLatitude,
			JobLongitude:       job.JobLongitude,
			JobStartTime:       job.JobStartTime,
			JobDurationHours:   job.JobDurationHours,
			JobCategories:      job.JobCategories,
			MatchingCategories: job.MatchingCategories,
			ContractorID:       job.ContractorID,
			ContractorCompany:  job.ContractorCompany,
			ProcessID:          job.ProcessID,
			ProcessStatus:      job.ProcessStatus,
			DistanceMeters:     job.DistanceMeters,
		})
	}

	return responses.GetJobsByWorkerResponse{
		Jobs: jobDetails,
	}, nil
}
