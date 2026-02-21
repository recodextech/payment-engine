package repositories

import (
	"context"
	"encoding/json"

	"payment-engine/internal/domain"
	"payment-engine/internal/domain/adaptors"
	"payment-engine/internal/domain/adaptors/database"
	"payment-engine/internal/domain/application"
	"payment-engine/pkg/errors"

	gopostgres "github.com/HADLakmal/go-postgres"
	"github.com/recodextech/api-definitions/events"
	"github.com/recodextech/container"
)

const (
	paymentTableName = `"fix.payments"`
)

type PaymentRepository struct {
	dbAdaptor database.FixFlowDB
	log       adaptors.Logger
	gopostgres.DatabaseReporter
}

func (r *PaymentRepository) Init(c container.Container) error {
	r.dbAdaptor = c.Resolve(application.MoudleDBConector).(database.FixFlowDB)
	r.log = c.Resolve(application.ModuleLogger).(adaptors.Logger).NewLog(adaptors.LoggerPrefixed(`repo.payment`))
	return nil
}

// CreatePayment creates a new payment record and returns the payment ID
func (r *PaymentRepository) CreatePayment(ctx context.Context, payment events.PaymentEvent) (string, error) {
	paymentID, err := r.createPaymentChange(ctx, payment)
	if err != nil {
		return "", errors.Wrap(err, "failed to create payment change")
	}
	return paymentID, nil
}

func (r *PaymentRepository) createPaymentChange(ctx context.Context, payment events.PaymentEvent) (workerID string, err error) {
	metaJSON, err := json.Marshal(payment.EventMeta)
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal meta")
	}

	payloadJSON, err := json.Marshal(payment.Payload)
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal payload")
	}

	status := domain.PaymentPending

	columns := []string{columnParamKey, columnParamMeta, columnParamJobID, columnParamType, columnParamStatus, columnParamAmount, columnParamPayload}
	values := []interface{}{payment.Payload.ID, string(metaJSON), payment.Payload.JobID, payment.Payload.Type, status, payment.Payload.Amount, string(payloadJSON)}

	paymentIDRet, err := r.dbAdaptor.InsertDataRow(ctx, domain.PaymentTable, columns, values)
	if err != nil {
		return "", errors.Wrap(err, "failed to create payment")
	}
	return paymentIDRet, nil
}

func (r *PaymentRepository) createTransaction(ctx context.Context, transaction events.TransactionPayload, meta events.Meta) (workerID string, err error) {
	metaJSON, err := json.Marshal(meta)
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal meta")
	}

	payloadJSON, err := json.Marshal(transaction)
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal payload")
	}

	columns := []string{columnParamKey, columnParamMeta, columnParamPaymentID, columnParamType, columnParamPayeeID, columnParamAmount, columnParamPayload}
	values := []interface{}{transaction.ID, string(metaJSON), transaction.PaymentID, transaction.Type, transaction.Payee.ID, transaction.Amount, string(payloadJSON)}

	paymentIDRet, err := r.dbAdaptor.InsertDataRow(ctx, domain.TransactionTable, columns, values)
	if err != nil {
		return "", errors.Wrap(err, "failed to create transaction")
	}
	return paymentIDRet, nil
}

// GetPaymentByID retrieves a payment by its ID
func (r *PaymentRepository) GetPaymentByID(ctx context.Context, paymentID string) (events.PaymentEvent, error) {
	return events.PaymentEvent{}, errors.New("not implemented")
}

// UpdatePayment updates an existing payment record
func (r *PaymentRepository) UpdatePayment(ctx context.Context, payment events.PaymentEvent) error {
	return errors.New("not implemented")
}
