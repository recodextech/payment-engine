package validator

import (
	"context"

	"payment-engine/internal/domain/adaptors"
	"payment-engine/internal/domain/adaptors/repositories"
	"payment-engine/internal/domain/application"
	"payment-engine/internal/http/request"
	"payment-engine/pkg/errors"

	"github.com/recodextech/container"
	"github.com/recodextech/krouter"
)

type JobValidator struct {
	log     adaptors.Logger
	jobRepo repositories.JobRepository
}

func (j *JobValidator) Init(container container.Container) error {
	j.log = container.Resolve(application.ModuleLogger).(adaptors.Logger).NewLog(adaptors.LoggerPrefixed(
		`validators.job-id`))
	j.jobRepo = container.Resolve(application.ModuleJobRepo).(repositories.JobRepository)

	return nil
}

func (j JobValidator) Validate(ctx context.Context, v interface{}, requestParams krouter.RequestParams) error {
	jobID := requestParams.Param(request.PathParamJobID.String())

	exists, err := j.jobRepo.Exists(ctx, jobID)
	if err != nil {
		return errors.Wrap(err, "failed to validate job ID")
	}

	if !exists {
		return errors.New(ErrMsgInvalidJobID)
	}

	return nil
}
