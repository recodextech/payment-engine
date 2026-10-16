package services

import (
	"context"
	"math"

	"payment-engine/internal/domain"
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
	if req.JobID == "" {
		return "", errors.New("job_id is required: every payment must be tied to a job")
	}
	if req.PayerID == "" {
		return "", errors.New("payer_id is required: the wallet the amount is debited from")
	}

	payer, err := s.accountWalletRepo.GetWalletByID(ctx, req.PayerID)
	if err != nil {
		return "", errors.Wrap(err, "failed to get payer wallet")
	}
	walletType := payer.Payload.Type

	// Double entry: the payer is debited the full amount and each payee is
	// credited its share, so worker earnings and contractor spend both stay traceable.
	transactionEntries := make([]events.TransactionPayload, 0, len(req.TransactionEntries)+1)
	transactionEntries = append(transactionEntries, s.newTransaction(req, req.PayerID, payer, -req.Amount))

	var credited float64
	for _, entry := range req.TransactionEntries {
		if entry.Amount <= 0 {
			return "", errors.New("transaction entry amount must be greater than 0")
		}
		if entry.WalletType != "" && entry.WalletType != walletType {
			return "", errors.Errorf("transaction entry wallet type %s does not match payer wallet type %s",
				entry.WalletType, walletType)
		}
		payeeID, payee, err := s.resolvePayeeWallet(ctx, entry, walletType)
		if err != nil {
			return "", err
		}
		if payeeID == req.PayerID {
			return "", errors.New("payer and payee wallets must be different")
		}
		if payee.Payload.Type != walletType {
			return "", errors.Errorf("payee wallet type %s does not match payer wallet type %s",
				payee.Payload.Type, walletType)
		}
		credited += entry.Amount
		transactionEntries = append(transactionEntries, s.newTransaction(req, payeeID, payee, entry.Amount))
	}
	if math.Abs(credited-req.Amount) > amountTolerance {
		return "", errors.Errorf("transaction entries add up to %.2f but the payment amount is %.2f", credited, req.Amount)
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

// amountTolerance absorbs float rounding when comparing amounts
const amountTolerance = 0.000001

// resolvePayeeWallet returns the wallet to credit: the given wallet, or the
// payee account's wallet of the payer's type.
func (s *PaymentService) resolvePayeeWallet(ctx context.Context, entry request.TransactionEntry,
	walletType string) (string, events.AccountWalletEvent, error) {
	if entry.PayeeID != "" {
		wallet, err := s.accountWalletRepo.GetWalletByID(ctx, entry.PayeeID)
		if err != nil {
			return "", wallet, errors.Wrap(err, "failed to get payee wallet")
		}
		return entry.PayeeID, wallet, nil
	}
	if entry.PayeeAccountID == "" {
		return "", events.AccountWalletEvent{}, errors.New("transaction entry needs payee_id or payee_account_id")
	}
	wallet, exist, err := s.accountWalletRepo.GetWalletByType(ctx, entry.PayeeAccountID, domain.WalletType(walletType))
	if err != nil && !errors.Is(err, errors.RepositoryDataNotExistError{}) {
		return "", wallet, errors.Wrap(err, "failed to get payee wallet")
	}
	if !exist {
		return "", wallet, errors.Errorf("payee account %s has no %s wallet", entry.PayeeAccountID, walletType)
	}
	return wallet.Payload.ID, wallet, nil
}

func (s *PaymentService) newTransaction(req request.CreatePaymentRequest, walletID string,
	wallet events.AccountWalletEvent, amount float64) events.TransactionPayload {
	transaction := events.TransactionPayload{
		ID:     uuid.New().String(),
		Type:   wallet.Payload.Type,
		Amount: amount,
	}
	transaction.Payee.ID = walletID
	transaction.Payee.CurrentBalance = wallet.Payload.Balance
	transaction.Payee.SequenceNumber = wallet.Payload.SequenceNumber
	transaction.PaymentInfromation.Method = req.PaymentInfromation.Method
	transaction.PaymentInfromation.Reference = req.PaymentInfromation.Reference
	return transaction
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

// CancelPayment cancels an in-progress payment
func (s *PaymentService) CancelPayment(ctx context.Context, req request.CancelPaymentRequest) error {
	err := s.paymentRepo.UpdateInProgressPaymentToCancelled(ctx, req.PaymentID, req.JobID)
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
