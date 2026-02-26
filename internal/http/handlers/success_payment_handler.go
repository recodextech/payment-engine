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

type SuccessPaymentHandler struct {
	paymentService service.PaymentService
}

func (h *SuccessPaymentHandler) Init(c container.Container) error {
	h.paymentService = c.Resolve(application.ModulePaymentService).(service.PaymentService)
	return nil
}

func (h *SuccessPaymentHandler) Handle(ctx context.Context, payload krouter.HttpPayload) (interface{}, error) {
	// Get account_id from header
	accountID := payload.Header(request.HeaderAccountID.String())

	if accountID == "" {
		return nil, errors.New("account-id header is required")
	}

	jobID := payload.Param(request.PathParamJobID.String())
	paymentID := payload.Param(request.PathParamPaymentID.String())

	// Call service to mark payment as successful
	err := h.paymentService.SuccessPayment(ctx, request.SuccessPaymentRequest{
		AccountID: accountID,
		JobID:     jobID,
		PaymentID: paymentID,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to update payment to success")
	}

	// Return response
	return nil, nil
}
