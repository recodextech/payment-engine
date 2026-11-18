package repositories

import (
	"context"
	"encoding/json"

	"payment-engine/internal/domain"
	"payment-engine/internal/domain/application"
	"payment-engine/pkg/errors"

	gopostgres "github.com/HADLakmal/go-postgres"
	"github.com/recodextech/api-definitions/events"
	"github.com/recodextech/container"
)

// SubscriptionRepository stores subscription rows (monthly plans and add-on packs), payments and the
// job quota ledger.
type SubscriptionRepository struct {
	sqlRunner
}

func (r *SubscriptionRepository) Init(c container.Container) error {
	r.db = c.Resolve(application.ModuleSQL).(gopostgres.DatabaseReporter)
	return nil
}

func (r *SubscriptionRepository) InTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return r.db.InTransaction(ctx, fn, nil)
}

func (r *SubscriptionRepository) LockContractor(ctx context.Context, contractorID string) error {
	var key string
	exist, err := r.queryRow(ctx, `SELECT key FROM `+domain.ContractorTable+` WHERE key = $1::uuid FOR UPDATE`,
		[]any{contractorID}, &key)
	if err != nil {
		return errors.Wrap(err, "failed to lock contractor")
	}
	if !exist {
		return errors.NotFound("contractor")
	}
	return nil
}

func unmarshalRow(meta, payload []byte, metaDst, payloadDst any, what string) error {
	if err := json.Unmarshal(meta, metaDst); err != nil {
		return errors.Wrap(err, "failed to unmarshal "+what+" meta")
	}
	if len(payload) > 0 {
		if err := json.Unmarshal(payload, payloadDst); err != nil {
			return errors.Wrap(err, "failed to unmarshal "+what+" payload")
		}
	}
	return nil
}

// ---- subscription rows ----

// rowColumns is the column order scanRow reads and CreateRow writes.
var rowColumns = []string{
	columnParamKey, columnParamMeta, columnParamContractorID, columnParamPlanID, columnParamKind,
	columnParamStatus, columnParamPeriodEnd, columnParamRemainingJobs, columnParamPayload,
}

func scanRow(scan func(dest ...any) error) (events.SubscriptionEvent, error) {
	var row events.SubscriptionEvent
	var key, contractorID, planID, kind, status string
	var end int64
	var remaining int
	var meta, payload []byte
	if err := scan(&key, &meta, &contractorID, &planID, &kind, &status, &end, &remaining, &payload); err != nil {
		return row, errors.Wrap(err, "failed to scan subscription")
	}
	if err := unmarshalRow(meta, payload, &row.EventMeta, &row.Payload, "subscription"); err != nil {
		return row, errors.Wrap(err, "failed to decode subscription row")
	}
	p := &row.Payload
	p.ID, p.ContractorID, p.PlanID = parseUUID(key), parseUUID(contractorID), parseUUID(planID)
	p.Kind, p.Status, p.PeriodEnd, p.RemainingJobs = kind, status, end, remaining
	return row, nil
}

func (r *SubscriptionRepository) listRows(ctx context.Context, where string, args []any) ([]events.SubscriptionEvent, error) {
	var rows []events.SubscriptionEvent
	err := r.queryRows(ctx, `SELECT `+columnList(rowColumns)+` FROM `+domain.SubscriptionTable+` WHERE deleted = FALSE AND `+where,
		args, func(res gopostgres.RowsInterface) error {
			row, err := scanRow(res.Scan)
			if err != nil {
				return errors.Wrap(err, "failed to read subscription row")
			}
			rows = append(rows, row)
			return nil
		})
	if err != nil {
		return nil, errors.Wrap(err, "failed to query subscriptions")
	}
	return rows, nil
}

func (r *SubscriptionRepository) ListActiveRows(ctx context.Context, contractorID string) ([]events.SubscriptionEvent, error) {
	return r.listRows(ctx, `contractor_id = $1::uuid AND status = 'ACTIVE' ORDER BY period_end, created_at`, []any{contractorID})
}

func (r *SubscriptionRepository) GetRowByID(ctx context.Context, id string) (events.SubscriptionEvent, bool, error) {
	rows, err := r.listRows(ctx, `key = $1::uuid`, []any{id})
	if err != nil {
		return events.SubscriptionEvent{}, false, errors.Wrap(err, "failed to get subscription row")
	}
	if len(rows) == 0 {
		return events.SubscriptionEvent{}, false, nil
	}
	return rows[0], true, nil
}

func (r *SubscriptionRepository) ListRows(ctx context.Context, status, kind string, limit int) ([]events.SubscriptionEvent, error) {
	return r.listRows(ctx, `($1 = '' OR status = $1) AND ($2 = '' OR kind = $2) ORDER BY created_at DESC LIMIT $3`,
		[]any{status, kind, limit})
}

func (r *SubscriptionRepository) ListDueContractors(ctx context.Context, nowMillis int64, limit int) ([]string, error) {
	var ids []string
	err := r.queryRows(ctx, `SELECT contractor_id FROM `+domain.SubscriptionTable+`
		WHERE deleted = FALSE AND status = 'ACTIVE' AND period_end <= $1
		GROUP BY contractor_id ORDER BY MIN(period_end) LIMIT $2`,
		[]any{nowMillis, limit}, func(res gopostgres.RowsInterface) error {
			var id string
			if err := res.Scan(&id); err != nil {
				return errors.Wrap(err, "failed to scan contractor id")
			}
			ids = append(ids, id)
			return nil
		})
	if err != nil {
		return nil, errors.Wrap(err, "failed to query due subscriptions")
	}
	return ids, nil
}

func (r *SubscriptionRepository) CreateRow(ctx context.Context, row events.SubscriptionEvent) error {
	meta, err := toJSON(row.EventMeta)
	if err != nil {
		return errors.Wrap(err, "failed to encode subscription meta")
	}
	payload, err := toJSON(row.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to encode subscription payload")
	}
	p := row.Payload
	_, err = r.exec(ctx, insertQuery(domain.SubscriptionTable, rowColumns),
		p.ID.String(), meta, p.ContractorID.String(), p.PlanID.String(), p.Kind, p.Status, p.PeriodEnd, p.RemainingJobs, payload)
	if err != nil {
		return errors.Wrap(err, "failed to create subscription")
	}
	return nil
}

// EndRow marks the row ENDED; remaining_jobs is left as it was.
func (r *SubscriptionRepository) EndRow(ctx context.Context, row events.SubscriptionEvent) error {
	meta, err := toJSON(row.EventMeta)
	if err != nil {
		return errors.Wrap(err, "failed to encode subscription meta")
	}
	payload, err := toJSON(row.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to encode subscription payload")
	}
	_, err = r.exec(ctx, `UPDATE `+domain.SubscriptionTable+` SET `+setClause(2, columnParamMeta, columnParamStatus, columnParamPayload)+
		` WHERE `+columnParamKey+` = $1::uuid`,
		row.Payload.ID.String(), meta, domain.SubscriptionEnded, payload)
	if err != nil {
		return errors.Wrap(err, "failed to end subscription")
	}
	return nil
}

func (r *SubscriptionRepository) TakeJob(ctx context.Context, rowID string) (int, bool, error) {
	var remaining int
	ok, err := r.queryRow(ctx, `UPDATE `+domain.SubscriptionTable+
		` SET remaining_jobs = remaining_jobs - 1 WHERE key = $1::uuid AND remaining_jobs > 0 RETURNING remaining_jobs`,
		[]any{rowID}, &remaining)
	if err != nil {
		return 0, false, errors.Wrap(err, "failed to take a job")
	}
	return remaining, ok, nil
}

func (r *SubscriptionRepository) ReturnJob(ctx context.Context, rowID string, nowMillis int64) (bool, error) {
	var remaining int
	ok, err := r.queryRow(ctx, `UPDATE `+domain.SubscriptionTable+
		` SET remaining_jobs = remaining_jobs + 1 WHERE key = $1::uuid AND status = 'ACTIVE' AND period_end > $2 RETURNING remaining_jobs`,
		[]any{rowID, nowMillis}, &remaining)
	if err != nil {
		return false, errors.Wrap(err, "failed to return a job")
	}
	return ok, nil
}

// ---- payments ----

// paymentColumns is the column order scanSubscriptionPayment reads and CreatePayment writes.
var paymentColumns = []string{
	columnParamKey, columnParamMeta, columnParamContractorID, columnParamProvider,
	columnParamProviderRef, columnParamStatus, columnParamPayload,
}

func scanSubscriptionPayment(scan func(dest ...any) error) (events.SubscriptionPaymentEvent, error) {
	var pay events.SubscriptionPaymentEvent
	var key, contractorID, provider, ref, status string
	var meta, payload []byte
	if err := scan(&key, &meta, &contractorID, &provider, &ref, &status, &payload); err != nil {
		return pay, errors.Wrap(err, "failed to scan payment")
	}
	if err := unmarshalRow(meta, payload, &pay.EventMeta, &pay.Payload, "payment"); err != nil {
		return pay, errors.Wrap(err, "failed to decode payment row")
	}
	p := &pay.Payload
	p.ID, p.ContractorID, p.Provider, p.ProviderRef, p.Status = parseUUID(key), parseUUID(contractorID), provider, ref, status
	return pay, nil
}

func (r *SubscriptionRepository) listPayments(ctx context.Context, where string, args []any) ([]events.SubscriptionPaymentEvent, error) {
	var out []events.SubscriptionPaymentEvent
	err := r.queryRows(ctx, `SELECT `+columnList(paymentColumns)+` FROM `+domain.SubscriptionPaymentTable+` WHERE deleted = FALSE AND `+where,
		args, func(res gopostgres.RowsInterface) error {
			pay, err := scanSubscriptionPayment(res.Scan)
			if err != nil {
				return errors.Wrap(err, "failed to read payment row")
			}
			out = append(out, pay)
			return nil
		})
	if err != nil {
		return nil, errors.Wrap(err, "failed to query payments")
	}
	return out, nil
}

func (r *SubscriptionRepository) CreatePayment(ctx context.Context, pay events.SubscriptionPaymentEvent) error {
	meta, err := toJSON(pay.EventMeta)
	if err != nil {
		return errors.Wrap(err, "failed to encode payment meta")
	}
	payload, err := toJSON(pay.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to encode payment payload")
	}
	p := pay.Payload
	res, err := r.exec(ctx, insertQuery(domain.SubscriptionPaymentTable, paymentColumns)+
		` ON CONFLICT (`+columnList([]string{columnParamProvider, columnParamProviderRef})+`) DO NOTHING`,
		p.ID.String(), meta, p.ContractorID.String(), p.Provider, p.ProviderRef, p.Status, payload)
	if err != nil {
		return errors.Wrap(err, "failed to create payment")
	}
	n, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "failed to read rows affected")
	}
	if n == 0 {
		return errors.Duplicate("payment " + p.ProviderRef)
	}
	return nil
}

func (r *SubscriptionRepository) SetPaymentStatus(ctx context.Context, pay events.SubscriptionPaymentEvent) error {
	meta, err := toJSON(pay.EventMeta)
	if err != nil {
		return errors.Wrap(err, "failed to encode payment meta")
	}
	_, err = r.exec(ctx, `UPDATE `+domain.SubscriptionPaymentTable+` SET `+setClause(2, columnParamMeta, columnParamStatus)+
		` WHERE `+columnParamKey+` = $1::uuid`,
		pay.Payload.ID.String(), meta, pay.Payload.Status)
	if err != nil {
		return errors.Wrap(err, "failed to update payment")
	}
	return nil
}

func (r *SubscriptionRepository) ListPaidPayments(ctx context.Context, contractorID string) ([]events.SubscriptionPaymentEvent, error) {
	return r.listPayments(ctx, `contractor_id = $1::uuid AND status = 'PAID' ORDER BY created_at`, []any{contractorID})
}

func (r *SubscriptionRepository) ListPayments(ctx context.Context, contractorID string, limit int) ([]events.SubscriptionPaymentEvent, error) {
	return r.listPayments(ctx, `contractor_id = $1::uuid ORDER BY created_at DESC LIMIT $2`, []any{contractorID, limit})
}

// ---- quota ledger ----

// ledgerColumns is the column order GetLedgerByProcessForUpdate reads and SaveLedger writes.
var ledgerColumns = []string{columnParamKey, columnParamMeta, columnParamProcessId, columnParamStatus, columnParamPayload}

func (r *SubscriptionRepository) GetLedgerByProcessForUpdate(ctx context.Context, processID string) (events.QuotaLedgerEvent, bool, error) {
	var entry events.QuotaLedgerEvent
	var key, process, status string
	var meta, payload []byte
	exist, err := r.queryRow(ctx, `SELECT `+columnList(ledgerColumns)+` FROM `+domain.QuotaLedgerTable+
		` WHERE `+columnParamProcessId+` = $1::uuid FOR UPDATE`, []any{processID}, &key, &meta, &process, &status, &payload)
	if err != nil {
		return entry, false, errors.Wrap(err, "failed to get quota ledger")
	}
	if !exist {
		return entry, false, nil
	}
	if err := unmarshalRow(meta, payload, &entry.EventMeta, &entry.Payload, "ledger"); err != nil {
		return entry, false, errors.Wrap(err, "failed to decode quota ledger row")
	}
	entry.Payload.ID, entry.Payload.ProcessID, entry.Payload.Status = parseUUID(key), parseUUID(process), status
	return entry, true, nil
}

func (r *SubscriptionRepository) SaveLedger(ctx context.Context, entry events.QuotaLedgerEvent) error {
	meta, err := toJSON(entry.EventMeta)
	if err != nil {
		return errors.Wrap(err, "failed to encode quota ledger meta")
	}
	payload, err := toJSON(entry.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to encode quota ledger payload")
	}
	_, err = r.exec(ctx, insertQuery(domain.QuotaLedgerTable, ledgerColumns)+
		` ON CONFLICT (`+columnParamProcessId+`) DO UPDATE SET `+excludedSet(columnParamStatus, columnParamMeta, columnParamPayload),
		entry.Payload.ID.String(), meta, entry.Payload.ProcessID.String(), entry.Payload.Status, payload)
	if err != nil {
		return errors.Wrap(err, "failed to save quota ledger")
	}
	return nil
}
