package repositories

import (
	"context"
	"encoding/json"

	"payment-engine/internal/domain"
	"payment-engine/internal/domain/adaptors"
	"payment-engine/internal/domain/adaptors/database"
	"payment-engine/internal/domain/application"
	"payment-engine/internal/domain/models"
	"payment-engine/pkg/errors"

	gopostgres "github.com/HADLakmal/go-postgres"
	"github.com/recodextech/api-definitions/events"
	"github.com/recodextech/container"
)

type AccountWalletRepository struct {
	dbAdaptor database.FixFlowDB
	log       adaptors.Logger
	gopostgres.DatabaseReporter
}

func (r *AccountWalletRepository) Init(c container.Container) error {
	r.dbAdaptor = c.Resolve(application.MoudleDBConector).(database.FixFlowDB)
	r.log = c.Resolve(application.ModuleLogger).(adaptors.Logger).NewLog(adaptors.LoggerPrefixed(`repo.account-wallet`))
	r.DatabaseReporter = c.Resolve(application.ModuleSQL).(gopostgres.DatabaseReporter)
	return nil
}

// CreateWallet creates a new wallet record for an account and returns the wallet ID
func (r *AccountWalletRepository) CreateWallet(ctx context.Context, accountID string, wallet events.AccountWalletEvent) (string, error) {
	metaJSON, err := json.Marshal(wallet.EventMeta)
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal meta")
	}

	valueJSON, err := json.Marshal(wallet.Payload)
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal value")
	}

	status := domain.WalletStatusActive

	columns := []string{columnParamKey, columnParamAccountID, columnParamType, columnParamMeta, columnParamPayload, columnParamBalance, columnParamStatus}
	values := []any{wallet.Payload.ID, accountID, wallet.Payload.Type, string(metaJSON), string(valueJSON), wallet.Payload.Balance, status}

	walletID, err := r.dbAdaptor.InsertDataRow(ctx, domain.AccountWalletTable, columns, values)
	if err != nil {
		return "", errors.Wrap(err, "failed to create wallet")
	}

	return walletID, nil
}

// CreateInternalWallets creates a new internal wallet record for an account and returns the wallet ID
func (r *AccountWalletRepository) CreateInternalWallets(ctx context.Context, accountID string, wallets []events.AccountWalletEvent) ([]string, error) {
	var walletIDs []string
	err := r.DatabaseReporter.InTransaction(ctx, func(ctx context.Context) error {
		for _, wallet := range wallets {
			metaJSON, err := json.Marshal(wallet.EventMeta)
			if err != nil {
				return errors.Wrap(err, "failed to marshal meta")
			}
			valueJSON, err := json.Marshal(wallet.Payload)
			if err != nil {
				return errors.Wrap(err, "failed to marshal value")
			}

			status := domain.WalletStatusActive

			columns := []string{columnParamKey, columnParamAccountID, columnParamType, columnParamMeta, columnParamPayload, columnParamBalance, columnParamStatus}
			values := []any{wallet.Payload.ID, accountID, wallet.Payload.Type, string(metaJSON), string(valueJSON), wallet.Payload.Balance, status}

			walletID, err := r.dbAdaptor.InsertDataRow(ctx, domain.AccountWalletTable, columns, values)
			if err != nil {
				return errors.Wrap(err, "failed to create wallet")
			}
			walletIDs = append(walletIDs, walletID)
		}

		return nil
	}, nil)
	if err != nil {
		return []string{}, errors.Wrap(err, "failed to create internal wallet")
	}

	return walletIDs, nil
}

// GetWalletByID retrieves a wallet by its ID
func (r *AccountWalletRepository) GetWalletByID(ctx context.Context, walletID string) (walletRes events.AccountWalletEvent, err error) {
	columns := []string{columnParamAccountID, columnParamType, columnParamMeta, columnParamPayload, columnParamBalance}
	whereClause := `key = $1 AND deleted = false`

	result, err := r.dbAdaptor.GetDataRowWithResult(ctx, domain.AccountWalletTable, columns, whereClause, []interface{}{walletID})
	if err != nil {
		return walletRes, errors.Wrap(err, "failed to get wallet")
	}
	var walletPayloadBytes, walletMeta []byte
	var accountID, walletType string
	var walletBalance float64

	exist, err := result.Scan(
		&accountID,
		&walletType,
		&walletMeta,
		&walletPayloadBytes,
		&walletBalance)
	if err != nil {
		return walletRes, errors.Wrap(err, "failed to scan wallet row")
	}
	if !exist {
		return walletRes, errors.New("wallet not found")
	}
	err = json.Unmarshal(walletPayloadBytes, &walletRes.Payload)
	if err != nil {
		return walletRes, errors.Wrap(err, "failed to unmarshal wallet payload")
	}
	err = json.Unmarshal(walletMeta, &walletRes.EventMeta)
	if err != nil {
		return walletRes, errors.Wrap(err, "failed to unmarshal wallet meta")
	}

	walletRes.Payload.AccountID = accountID
	walletRes.Payload.Type = walletType
	walletRes.Payload.Balance = walletBalance

	return walletRes, nil
}

// GetWalletByIDForUpdate retrieves a wallet by ID with a pessimistic row lock
// (SELECT ... FOR UPDATE).  It must be called inside an InTransaction callback.
func (r *AccountWalletRepository) GetWalletByIDForUpdate(ctx context.Context, walletID string) (walletRes events.AccountWalletEvent, err error) {
	columns := []string{columnParamAccountID, columnParamType, columnParamMeta, columnParamPayload, columnParamBalance}
	whereClause := `key = $1 AND deleted = false`

	result, err := r.dbAdaptor.GetDataRowWithResultForUpdate(ctx, domain.AccountWalletTable, columns, whereClause, []interface{}{walletID})
	if err != nil {
		return walletRes, errors.Wrap(err, "failed to lock wallet row")
	}
	var walletPayloadBytes, walletMeta []byte
	var accountID, walletType string
	var walletBalance float64

	exist, err := result.Scan(
		&accountID,
		&walletType,
		&walletMeta,
		&walletPayloadBytes,
		&walletBalance)
	if err != nil {
		return walletRes, errors.Wrap(err, "failed to scan wallet row")
	}
	if !exist {
		return walletRes, errors.New("wallet not found")
	}
	err = json.Unmarshal(walletPayloadBytes, &walletRes.Payload)
	if err != nil {
		return walletRes, errors.Wrap(err, "failed to unmarshal wallet payload")
	}
	err = json.Unmarshal(walletMeta, &walletRes.EventMeta)
	if err != nil {
		return walletRes, errors.Wrap(err, "failed to unmarshal wallet meta")
	}

	walletRes.Payload.AccountID = accountID
	walletRes.Payload.Type = walletType
	walletRes.Payload.Balance = walletBalance

	return walletRes, nil
}

// GetWalletByAccountID retrieves wallet(s) for an account
func (r *AccountWalletRepository) GetWalletsByAccountID(ctx context.Context, accountID string) (walletRes []events.AccountWalletEvent, exist bool, err error) {
	columns := []string{columnParamAccountID, columnParamType, columnParamMeta, columnParamPayload, columnParamBalance}
	whereClause := `account_id = $1 AND deleted = false`

	rows, err := r.dbAdaptor.GetDataRowsWithResult(ctx, domain.AccountWalletTable, columns, whereClause, []interface{}{accountID})
	if err != nil {
		return walletRes, false, errors.Wrap(err, "failed to get wallet")
	}

	defer rows.Close()

	for rows.Next() {
		walletData := events.AccountWalletEvent{}
		var walletPayloadBytes, walletMeta []byte
		var accountIDVar, walletType string
		var walletBalance float64
		err := rows.Scan(
			&accountIDVar,
			&walletType,
			&walletMeta,
			&walletPayloadBytes,
			&walletBalance)
		if err != nil {
			return walletRes, false, errors.Wrap(err, "failed to scan wallet row")
		}
		err = json.Unmarshal(walletPayloadBytes, &walletData.Payload)
		if err != nil {
			return walletRes, false, errors.Wrap(err, "failed to unmarshal wallet payload")
		}

		err = json.Unmarshal(walletMeta, &walletData.EventMeta)
		if err != nil {
			return walletRes, false, errors.Wrap(err, "failed to unmarshal wallet meta")
		}
		walletData.Payload.AccountID = accountIDVar
		walletData.Payload.Type = walletType
		walletData.Payload.Balance = walletBalance
		walletRes = append(walletRes, walletData)
	}
	if len(walletRes) == 0 {
		return walletRes, false, nil
	}

	return walletRes, true, nil
}

// GetWalletByType retrieves wallet type for an account
func (r *AccountWalletRepository) GetWalletByType(ctx context.Context, accountID string, walletType domain.WalletType) (walletRes events.AccountWalletEvent, exist bool, err error) {
	columns := []string{columnParamAccountID, columnParamMeta, columnParamPayload, columnParamBalance}
	whereClause := `account_id = $1 AND type = $2 AND deleted = false`

	result, err := r.dbAdaptor.GetDataRowWithResult(ctx, domain.AccountWalletTable, columns, whereClause, []interface{}{accountID, walletType.String()})
	if err != nil {
		return walletRes, false, errors.Wrap(err, "failed to get wallet")
	}

	var walletPayloadBytes, walletMeta []byte
	var accountIDVar string
	var walletBalance float64
	exist, err = result.Scan(
		&accountIDVar,
		&walletMeta,
		&walletPayloadBytes,
		&walletBalance)
	if err != nil {
		return walletRes, false, errors.Wrap(err, "failed to scan wallet row")
	}
	if !exist {
		return walletRes, false, errors.RepositoryDataNotExistError{}
	}
	err = json.Unmarshal(walletPayloadBytes, &walletRes.Payload)
	if err != nil {
		return walletRes, false, errors.Wrap(err, "failed to unmarshal wallet payload")
	}

	err = json.Unmarshal(walletMeta, &walletRes.EventMeta)
	if err != nil {
		return walletRes, false, errors.Wrap(err, "failed to unmarshal wallet meta")
	}
	walletRes.Payload.AccountID = accountIDVar
	walletRes.Payload.Type = walletType.String()
	walletRes.Payload.Balance = walletBalance

	return walletRes, true, nil
}

// UpdateWalletBalance updates the balance of a wallet.  The read is done with
// FOR UPDATE so the row is locked for the duration of the surrounding
// transaction, preventing concurrent double-spend.
func (r *AccountWalletRepository) UpdateWalletBalance(ctx context.Context, walletID string, balance float64) error {
	walletRes, err := r.GetWalletByIDForUpdate(ctx, walletID)
	if err != nil {
		return errors.Wrap(err, "failed to lock wallet for balance update")
	}

	walletRes.Payload.Balance = balance
	walletRes.Payload.SequenceNumber += 1
	walletRes.EventMeta = events.MetaUpdate(ctx, walletRes.EventMeta)

	metaJSON, err := json.Marshal(walletRes.EventMeta)
	if err != nil {
		return errors.Wrap(err, "failed to marshal meta")
	}

	valueJSON, err := json.Marshal(walletRes.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to marshal value")
	}

	err = r.dbAdaptor.UpdateDataRow(ctx, domain.AccountWalletTable, walletID, map[string]interface{}{
		columnParamMeta:    string(metaJSON),
		columnParamPayload: string(valueJSON),
		columnParamBalance: balance,
	})
	if err != nil {
		return errors.Wrap(err, "failed to update wallet balance")
	}

	return nil
}

// GetWalletTotals sums the wallet's ledger entries: credits (amount > 0) and
// debits (amount < 0, returned as a positive number).
func (r *AccountWalletRepository) GetWalletTotals(ctx context.Context, walletID string) (models.WalletTotals, error) {
	columns := []string{
		`COALESCE(SUM(amount) FILTER (WHERE amount > 0), 0)`,
		`COALESCE(-SUM(amount) FILTER (WHERE amount < 0), 0)`,
	}
	whereClause := `payee_id = $1 AND deleted = false`

	result, err := r.dbAdaptor.GetDataRowWithResult(ctx, domain.TransactionTable, columns, whereClause, []interface{}{walletID})
	if err != nil {
		return models.WalletTotals{}, errors.Wrap(err, "failed to get wallet totals")
	}

	var totals models.WalletTotals
	if _, err = result.Scan(&totals.Credited, &totals.Debited); err != nil {
		return models.WalletTotals{}, errors.Wrap(err, "failed to scan wallet totals")
	}

	totals.Held, err = r.getWalletHeld(ctx, walletID)
	if err != nil {
		return models.WalletTotals{}, err
	}
	return totals, nil
}

// getWalletHeld sums the debit entries of the wallet's in-progress payments:
// the amount on hold until those jobs complete or are cancelled.
func (r *AccountWalletRepository) getWalletHeld(ctx context.Context, walletID string) (float64, error) {
	table := domain.PaymentTable + ` p CROSS JOIN LATERAL jsonb_array_elements(p.payload::jsonb->'transaction_entries') e`
	columns := []string{`COALESCE(-SUM((e->>'amount')::float8), 0)`}
	whereClause := `p.status = $1 AND p.deleted = false AND e->'payee'->>'id' = $2 AND (e->>'amount')::float8 < 0`

	result, err := r.dbAdaptor.GetDataRowWithResult(ctx, table, columns, whereClause,
		[]any{domain.PaymentPending.String(), walletID})
	if err != nil {
		return 0, errors.Wrap(err, "failed to get wallet held amount")
	}

	var held float64
	if _, err = result.Scan(&held); err != nil {
		return 0, errors.Wrap(err, "failed to scan wallet held amount")
	}
	return held, nil
}

// GetWalletTransactions returns the wallet's latest ledger entries, newest first.
func (r *AccountWalletRepository) GetWalletTransactions(ctx context.Context, walletID string, limit int) ([]models.WalletTransaction, error) {
	table := domain.TransactionTable + ` t JOIN ` + domain.PaymentTable + ` p ON p.key = t.payment_id`
	columns := []string{
		`t.key`,
		`t.payment_id`,
		`COALESCE(p.job_id::text, '')`,
		`t.type`,
		`t.amount`,
		`COALESCE((t.payload->'payee'->>'current_balance')::float8, 0)`,
		`t.created_at`,
	}
	whereClause := `t.payee_id = $1 AND t.deleted = false ORDER BY t.created_at DESC LIMIT $2`

	rows, err := r.dbAdaptor.GetDataRowsWithResult(ctx, table, columns, whereClause, []interface{}{walletID, limit})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get wallet transactions")
	}
	defer rows.Close()

	transactions := make([]models.WalletTransaction, 0)
	for rows.Next() {
		var t models.WalletTransaction
		if err := rows.Scan(&t.ID, &t.PaymentID, &t.JobID, &t.WalletType, &t.Amount, &t.BalanceAfter, &t.CreatedAt); err != nil {
			return nil, errors.Wrap(err, "failed to scan wallet transaction")
		}
		transactions = append(transactions, t)
	}
	return transactions, nil
}
