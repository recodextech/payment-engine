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
	r.DatabaseReporter = c.Resolve(application.ModuleSQL).(gopostgres.DatabaseReporter)
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
	columns := []string{columnParamMeta, columnParamJobID, columnParamType, columnParamAmount, columnParamStatus, columnParamPayload}
	whereClause := `key = $1 AND deleted = false`

	result, err := r.dbAdaptor.GetDataRowWithResult(ctx, domain.PaymentTable, columns, whereClause, []interface{}{paymentID})
	if err != nil {
		return events.PaymentEvent{}, errors.Wrap(err, "failed to get payment")
	}

	var paymentPayloadBytes, paymentMeta []byte
	var jobIDVar, typeVar string
	var amountVar float64
	var statusVar string
	exist, err := result.Scan(
		&paymentMeta,
		&jobIDVar,
		&typeVar,
		&amountVar,
		&statusVar,
		&paymentPayloadBytes)
	if err != nil {
		return events.PaymentEvent{}, errors.Wrap(err, "failed to scan payment row")
	}
	if !exist {
		return events.PaymentEvent{}, errors.RepositoryDataNotExistError{}
	}
	paymentRes := events.PaymentEvent{}
	err = json.Unmarshal(paymentPayloadBytes, &paymentRes.Payload)
	if err != nil {
		return events.PaymentEvent{}, errors.Wrap(err, "failed to unmarshal payment payload")
	}

	err = json.Unmarshal(paymentMeta, &paymentRes.EventMeta)
	if err != nil {
		return events.PaymentEvent{}, errors.Wrap(err, "failed to unmarshal payment meta")
	}
	paymentRes.Payload.JobID = jobIDVar
	paymentRes.Payload.Type = typeVar
	paymentRes.Payload.Amount = amountVar
	paymentRes.Payload.Status = statusVar

	return paymentRes, nil
}

// UpdatePayment updates an existing payment record
func (r *PaymentRepository) UpdatePayment(ctx context.Context, payment events.PaymentEvent) error {
	return errors.New("not implemented")
}

// UpdatePayment updates an existing payment record
func (r *PaymentRepository) UpdateInProgressPaymentToCancelled(ctx context.Context, key string) error {
	err := r.DatabaseReporter.InTransaction(ctx, func(ctx context.Context) error {
		payment, errPayment := r.GetPaymentByID(ctx, key)
		if errPayment != nil {
			return errors.Wrap(errPayment, "failed to get payment by ID")
		}

		if payment.Payload.Status != domain.PaymentPending.String() {
			return errors.New("only in-progress payments can be cancelled")
		}
		payment.Payload.Status = domain.PaymentCancelled.String()
		payloadJSON, err := json.Marshal(payment.Payload)
		if err != nil {
			return errors.Wrap(err, "failed to marshal payload")
		}
		payment.EventMeta = events.MetaUpdate(ctx, payment.EventMeta)
		metaJSON, err := json.Marshal(payment.EventMeta)
		if err != nil {
			return errors.Wrap(err, "failed to marshal meta")
		}
		processUpdates := map[string]interface{}{
			columnParamStatus:    domain.PaymentCancelled.String(),
			columnParamUpdatedAt: "NOW()",
			columnParamPayload:   string(payloadJSON),
			columnParamMeta:      string(metaJSON),
		}
		err = r.dbAdaptor.UpdateDataRow(ctx, domain.PaymentTable, key, processUpdates)
		if err != nil {
			return errors.Wrap(err, "failed to update process status")
		}
		return nil
	}, nil)
	if err != nil {
		return errors.Wrap(err, "failed to update in progress payment to cancelled")
	}
	return nil
}

// UpdatePayment updates an existing payment record
func (r *PaymentRepository) UpdateInProgressPaymentToSuccess(ctx context.Context, key string) error {
	err := r.DatabaseReporter.InTransaction(ctx, func(ctx context.Context) error {
		payment, errPayment := r.GetPaymentByID(ctx, key)
		if errPayment != nil {
			return errors.Wrap(errPayment, "failed to get payment by ID")
		}
		if payment.Payload.Status != domain.PaymentPending.String() {
			return errors.New("only in-progress payments can be marked as successful")
		}
		payment.Payload.Status = domain.PaymentSuccess.String()
		payloadJSON, err := json.Marshal(payment.Payload)
		if err != nil {
			return errors.Wrap(err, "failed to marshal payload")
		}
		payment.EventMeta = events.MetaUpdate(ctx, payment.EventMeta)
		metaJSON, err := json.Marshal(payment.EventMeta)
		if err != nil {
			return errors.Wrap(err, "failed to marshal meta")
		}
		processUpdates := map[string]interface{}{
			columnParamStatus:    domain.PaymentSuccess.String(),
			columnParamUpdatedAt: "NOW()",
			columnParamPayload:   string(payloadJSON),
			columnParamMeta:      string(metaJSON),
		}
		err = r.dbAdaptor.UpdateDataRow(ctx, domain.PaymentTable, key, processUpdates)
		if err != nil {
			return errors.Wrap(err, "failed to update process status")
		}
		for _, transac := range payment.Payload.TransactionEntries {
			_, err = r.createTransaction(ctx, transac, payment.EventMeta)
			if err != nil {
				return errors.Wrap(err, "failed to create transaction for successful payment")
			}
		}
		return nil
	}, nil)
	if err != nil {
		return errors.Wrap(err, "failed to update in progress payment to success")
	}
	return nil
}

// Exists checks if a payment exists by its ID
func (r *PaymentRepository) Exists(ctx context.Context, paymentID string) (bool, error) {
	exist, err := r.dbAdaptor.GetDataRowByID(ctx, domain.PaymentTable, paymentID)
	if err != nil {
		return false, errors.Wrap(err, "failed to check payment existence")
	}

	return exist, nil
}
