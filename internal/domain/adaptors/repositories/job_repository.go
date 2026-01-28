package repositories

import (
	"context"
)

type JobRepository interface {
	GetJobsByWorkerAndContractor(ctx context.Context, workerID, contractorID string) ([]JobWithDetails, error)
	GetJobsByWorker(ctx context.Context, workerID string) ([]JobWithDetails, error)
}

type JobWithDetails struct {
	WorkerID           string
	WorkerStatus       string
	WorkerCategories   string
	JobID              string
	JobStatus          string
	JobLatitude        float64
	JobLongitude       float64
	JobStartTime       string
	JobDurationHours   float64
	JobCategories      string
	MatchingCategories string
	ContractorID       string
	ContractorCompany  string
	ProcessID          string
	ProcessStatus      string
	DistanceMeters     float64
}
