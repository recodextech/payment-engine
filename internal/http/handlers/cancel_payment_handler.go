package handlers

import (
	"context"
	"payment-engine/internal/domain/application"
	"payment-engine/internal/domain/service"
	"payment-engine/internal/http/request"
	"payment-engine/pkg/errors"

	"github.com/recodextech/container"
	"github.com/recodextech/krouter"
)

type CancelPaymentHandler struct {
	paymentService service.PaymentService
}

func (h *CancelPaymentHandler) Init(c container.Container) error {
	h.paymentService = c.Resolve(application.ModulePaymentService).(service.PaymentService)
	return nil
}

func (h *CancelPaymentHandler) Handle(ctx context.Context, payload krouter.HttpPayload) (interface{}, error) {
	// Get account_id from header
	accountID := payload.Header(request.HeaderAccountID.String())

	if accountID == "" {
		return nil, errors.New("account-id header is required")
	}

	jobID := payload.Param(request.PathParamJobID.String())
	paymentID := payload.Param(request.PathParamPaymentID.String())

	// Call service to cancel payment
	err := h.paymentService.CancelPayment(ctx, request.CancelPaymentRequest{
		AccountID: accountID,
		JobID:     jobID,
		PaymentID: paymentID,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to cancel payment")
	}

	// Return response
	return nil, nil
}
