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

type AccountValidator struct {
	log          adaptors.Logger
	repositories struct {
		account repositories.AccountRepository
	}
}

func (a *AccountValidator) Init(container container.Container) error {
	a.log = container.Resolve(application.ModuleLogger).(adaptors.Logger).NewLog(adaptors.LoggerPrefixed(
		`validators.account-id`))
	a.repositories.account = container.Resolve(application.ModuleAccountRepo).(repositories.AccountRepository)

	return nil
}

func (a AccountValidator) Validate(ctx context.Context, v interface{}, requestParams krouter.RequestParams) error {
	accountID := requestParams.Header(request.HeaderAccountID.String())
	exists, err := a.repositories.account.Exists(ctx, accountID)
	if err != nil {
		return errors.Wrap(err, "failed to validate account ID")
	}

	if !exists {
		return ValidatorError{errors.New(ErrMsgInvalidAccountID), ErrMsgInvalidAccountID}
	}

	return nil
}
