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

type PaymentValidator struct {
	log          adaptors.Logger
	repositories struct {
		payment repositories.PaymentRepository
	}
}

func (p *PaymentValidator) Init(container container.Container) error {
	p.log = container.Resolve(application.ModuleLogger).(adaptors.Logger).NewLog(adaptors.LoggerPrefixed(
		`validators.payment-id`))
	p.repositories.payment = container.Resolve(application.ModulePaymentRepo).(repositories.PaymentRepository)

	return nil
}

func (p PaymentValidator) Validate(ctx context.Context, v interface{}, requestParams krouter.RequestParams) error {
	paymentID := requestParams.Param(request.PathParamPaymentID.String())
	exists, err := p.repositories.payment.Exists(ctx, paymentID)
	if err != nil {
		return errors.Wrap(err, "failed to validate payment ID")
	}

	if !exists {
		return ValidatorError{errors.New(ErrMsgInvalidPaymentID), ErrMsgInvalidPaymentID}
	}

	return nil
}
