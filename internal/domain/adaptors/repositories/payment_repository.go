package repositories

import (
	"context"

	"github.com/recodextech/api-definitions/events"
)

// PaymentRepository defines the interface for payment data operations
type PaymentRepository interface {
	CreatePayment(ctx context.Context, payment events.PaymentEvent) (string, error)
	GetPaymentByID(ctx context.Context, paymentID string) (events.PaymentEvent, error)
	UpdatePayment(ctx context.Context, payment events.PaymentEvent) error
}
