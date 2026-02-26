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

const CreatePayment = "payment.create"

type CreatePaymentHandler struct {
	paymentService service.PaymentService
}

func (h *CreatePaymentHandler) Init(c container.Container) error {
	h.paymentService = c.Resolve(application.ModulePaymentService).(service.PaymentService)
	return nil
}

func (h *CreatePaymentHandler) Handle(ctx context.Context, payload krouter.HttpPayload) (interface{}, error) {
	// Get account_id from header
	accountID := payload.Header(request.HeaderAccountID.String())

	if accountID == "" {
		return nil, errors.New("account-id header is required")
	}

	// Get request body
	reqBody, ok := payload.Body.(request.CreatePaymentRequest)
	if !ok {
		return nil, errors.New("invalid request body")
	}
	reqBody.AccountID = accountID
	reqBody.JobID = payload.Param(request.PathParamJobID.String())

	// Call service to create payment
	paymentID, err := h.paymentService.CreatePayment(ctx, reqBody)
	if err != nil {
		return nil, err
	}

	// Return response
	return responses.CreatePaymentResponse{
		PaymentID: paymentID,
	}, nil
}
