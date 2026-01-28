package repositories

import (
	"context"

	"payment-engine/internal/domain"
	"payment-engine/internal/domain/adaptors/database"
	"payment-engine/internal/domain/adaptors/repositories"
	"payment-engine/internal/domain/application"
	"payment-engine/pkg/errors"

	gopostgres "github.com/HADLakmal/go-postgres"
	"github.com/recodextech/container"
)

type JobRepository struct {
	dbAdaptor database.MOSDB
	gopostgres.DatabaseReporter
}

func (repo *JobRepository) Init(c container.Container) error {
	repo.dbAdaptor = c.Resolve(application.MoudleDBConector).(database.MOSDB)
	repo.DatabaseReporter = c.Resolve(application.ModuleSQL).(gopostgres.DatabaseReporter)

	return nil
}

func (repo *JobRepository) GetJobsByWorkerAndContractor(ctx context.Context, workerID, contractorID string) ([]repositories.JobWithDetails, error) {
	whereClause := "worker_id = $1 AND contractor_id = $2"

	columns := []string{
		"worker_id",
		"worker_status",
		"worker_categories",
		"job_id",
		"job_status",
		"job_latitude",
		"job_longitude",
		"job_start_time::text",
		"job_duration_hours",
		"job_categories",
		"matching_categories",
		"contractor_id",
		"contractor_company",
		"process_id",
		"process_status",
		"distance_meters",
		"assignment_status",
	}

	rows, err := repo.dbAdaptor.GetDataRowsWithResult(ctx, domain.WorkerNearbyJobsView, columns, whereClause, []any{workerID, contractorID})
	if err != nil {
		return nil, errors.Wrap(err, `query jobs by worker and contractor failed`)
	}
	defer rows.Close()

	var jobs []repositories.JobWithDetails
	for rows.Next() {
		var job repositories.JobWithDetails
		err := rows.Scan(
			&job.WorkerID,
			&job.WorkerStatus,
			&job.WorkerCategories,
			&job.JobID,
			&job.JobStatus,
			&job.JobLatitude,
			&job.JobLongitude,
			&job.JobStartTime,
			&job.JobDurationHours,
			&job.JobCategories,
			&job.MatchingCategories,
			&job.ContractorID,
			&job.ContractorCompany,
			&job.ProcessID,
			&job.ProcessStatus,
			&job.DistanceMeters,
		)
		if err != nil {
			return nil, errors.Wrap(err, `scan job details failed`)
		}
		jobs = append(jobs, job)
	}

	return jobs, nil
}

func (repo *JobRepository) GetJobsByWorker(ctx context.Context, workerID string) ([]repositories.JobWithDetails, error) {
	whereClause := "worker_id = $1"

	columns := []string{
		"worker_id",
		"worker_status",
		"worker_categories",
		"job_id",
		"job_status",
		"job_latitude",
		"job_longitude",
		"job_start_time::text",
		"job_duration_hours",
		"job_categories",
		"matching_categories",
		"contractor_id",
		"contractor_company",
		"process_id",
		"process_status",
		"distance_meters",
	}

	rows, err := repo.dbAdaptor.GetDataRowsWithResult(ctx, domain.WorkerNearbyJobsView, columns, whereClause, []any{workerID})
	if err != nil {
		return nil, errors.Wrap(err, `query jobs by worker and contractor failed`)
	}
	defer rows.Close()

	var jobs []repositories.JobWithDetails
	for rows.Next() {
		var job repositories.JobWithDetails
		err := rows.Scan(
			&job.WorkerID,
			&job.WorkerStatus,
			&job.WorkerCategories,
			&job.JobID,
			&job.JobStatus,
			&job.JobLatitude,
			&job.JobLongitude,
			&job.JobStartTime,
			&job.JobDurationHours,
			&job.JobCategories,
			&job.MatchingCategories,
			&job.ContractorID,
			&job.ContractorCompany,
			&job.ProcessID,
			&job.ProcessStatus,
			&job.DistanceMeters,
		)
		if err != nil {
			return nil, errors.Wrap(err, `scan job details failed`)
		}
		jobs = append(jobs, job)
	}

	return jobs, nil
}
