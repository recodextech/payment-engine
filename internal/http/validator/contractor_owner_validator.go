package validator

import (
	"context"
	"net/http"

	"payment-engine/internal/domain/adaptors/repositories"
	"payment-engine/internal/domain/application"
	"payment-engine/internal/http/request"
	"payment-engine/internal/services"
	"payment-engine/pkg/errors"

	"github.com/recodextech/container"
	"github.com/recodextech/krouter"
)

// ContractorOwnerValidator checks that the contractor in the path belongs to the account-id header,
// so one account can't read or change another contractor's subscription.
type ContractorOwnerValidator struct {
	accounts repositories.ContractorAccountRepository
}

func (v *ContractorOwnerValidator) Init(c container.Container) error {
	v.accounts = c.Resolve(application.ModuleContractorAccountRepo).(repositories.ContractorAccountRepository)
	return nil
}

func (v *ContractorOwnerValidator) Validate(ctx context.Context, _ interface{}, params krouter.RequestParams) error {
	owns, err := v.accounts.ContractorBelongsToAccount(ctx, params.Param(request.PathParamContractorID.String()),
		params.Header(request.HeaderAccountID.String()))
	if err != nil {
		return errors.Wrap(err, "failed to validate contractor")
	}
	if !owns {
		return services.NewServiceError(http.StatusForbidden, services.CodeNotContractorOwner, "contractor does not belong to this account")
	}
	return nil
}
