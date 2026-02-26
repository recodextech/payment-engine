package services

import (
	"context"

	"payment-engine/internal/domain/adaptors/repositories"
	"payment-engine/internal/domain/application"
	"payment-engine/internal/domain/models"
	"payment-engine/internal/http/request"
	"payment-engine/pkg/errors"

	"github.com/recodextech/api-definitions/events"
	"github.com/recodextech/api-definitions/pkg/uuid"
	"github.com/recodextech/container"
)

type PaymentService struct {
	paymentRepo       repositories.PaymentRepository
	accountWalletRepo repositories.AccountWalletRepository
}

func (s *PaymentService) Init(c container.Container) error {
	s.paymentRepo = c.Resolve(application.ModulePaymentRepo).(repositories.PaymentRepository)
	s.accountWalletRepo = c.Resolve(application.ModuleAccountWalletRepo).(repositories.AccountWalletRepository)
	return nil
}

// CreatePayment creates a new payment from the request
func (s *PaymentService) CreatePayment(ctx context.Context, req request.CreatePaymentRequest) (string, error) {
	// Validate required fields
	if req.Amount <= 0 {
		return "", errors.New("amount must be greater than 0")
	}

	// Convert request transaction entries to domain model
	transactionEntries := make([]events.TransactionPayload, 0, len(req.TransactionEntries))
	for _, entry := range req.TransactionEntries {
		wallet, err := s.accountWalletRepo.GetWalletByID(ctx, entry.PayeeID)
		if err != nil {
			return "", errors.Wrap(err, "failed to get wallet")
		}
		transaction := events.TransactionPayload{
			Amount: entry.Amount,
			ID:     uuid.New().String(),
		}
		transaction.Payee.ID = entry.PayeeID
		transaction.Payee.CurrentBalance = wallet.Payload.Balance
		transaction.Payee.SequenceNumber = wallet.Payload.SequenceNumber

		transaction.PaymentInfromation.Method = req.PaymentInfromation.Method
		transaction.PaymentInfromation.Reference = req.PaymentInfromation.Reference
		transactionEntries = append(transactionEntries, transaction)
	}

	// Build payment model
	payment := events.PaymentEvent{}
	payment.EventMeta = events.NewMetaContext(ctx)
	payment.Payload.ID = uuid.New().String()
	payment.Payload.Type = req.Type
	payment.Payload.JobID = req.JobID
	payment.Payload.TransactionEntries = transactionEntries
	payment.Payload.Amount = req.Amount
	payment.Payload.PaymentInfromation.Method = req.PaymentInfromation.Method
	payment.Payload.PaymentInfromation.Reference = req.PaymentInfromation.Reference
	payment.Payload.PaymentInfromation.Version = 1

	// Create payment in repository
	paymentID, err := s.paymentRepo.CreatePayment(ctx, payment)
	if err != nil {
		return "", errors.Wrap(err, "failed to create payment")
	}

	return paymentID, nil
}

// GetPayment retrieves a payment by ID
func (s *PaymentService) GetPayment(ctx context.Context, paymentID string) (payment models.Payment, err error) {
	if paymentID == "" {
		return payment, errors.New("payment_id is required")
	}

	paymentEvent, err := s.paymentRepo.GetPaymentByID(ctx, paymentID)
	if err != nil {
		return payment, errors.Wrap(err, "failed to get payment")
	}

	payment = models.Payment{
		ID: paymentEvent.Payload.ID,
	}

	return payment, nil
}

// UpdatePayment updates an existing payment
func (s *PaymentService) UpdatePayment(ctx context.Context, payment events.PaymentEvent) error {
	if payment.Payload.ID == "" {
		return errors.New("payment_id is required")
	}

	err := s.paymentRepo.UpdatePayment(ctx, payment)
	if err != nil {
		return errors.Wrap(err, "failed to update payment")
	}

	return nil
}

// CancelPayment cancels an in-progress payment
func (s *PaymentService) CancelPayment(ctx context.Context, req request.CancelPaymentRequest) error {
	err := s.paymentRepo.UpdateInProgressPaymentToCancelled(ctx, req.PaymentID)
	if err != nil {
		return errors.Wrap(err, "failed to cancel payment")
	}

	return nil
}

// SuccessPayment marks an in-progress payment as successful
func (s *PaymentService) SuccessPayment(ctx context.Context, req request.SuccessPaymentRequest) error {
	err := s.paymentRepo.UpdateInProgressPaymentToSuccess(ctx, req.PaymentID)
	if err != nil {
		return errors.Wrap(err, "failed to update payment to success")
	}

	return nil
}
