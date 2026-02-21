package service

import (
	"context"
	"payment-engine/internal/http/request"

	"github.com/recodextech/api-definitions/events"
)

type PaymentService interface {
	CreatePayment(ctx context.Context, req request.CreatePaymentRequest) (string, error)
	UpdatePayment(ctx context.Context, payment events.PaymentEvent) error
}
