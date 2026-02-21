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
	accountWalletTableName = `"fix.account.wallets"`
)

type AccountWalletRepository struct {
	dbAdaptor database.FixFlowDB
	log       adaptors.Logger
	gopostgres.DatabaseReporter
}

func (r *AccountWalletRepository) Init(c container.Container) error {
	r.dbAdaptor = c.Resolve(application.MoudleDBConector).(database.FixFlowDB)
	r.log = c.Resolve(application.ModuleLogger).(adaptors.Logger).NewLog(adaptors.LoggerPrefixed(`repo.account-wallet`))
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
	values := []interface{}{wallet.Payload.ID, accountID, wallet.Payload.Type, string(metaJSON), string(valueJSON), wallet.Payload.Balance, status}

	walletID, err := r.dbAdaptor.InsertDataRow(ctx, domain.AccountWalletTable, columns, values)
	if err != nil {
		return "", errors.Wrap(err, "failed to create wallet")
	}

	return walletID, nil
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

// GetWalletByAccountID retrieves wallet(s) for an account
func (r *AccountWalletRepository) GetWalletByAccountID(ctx context.Context, accountID string) (walletRes events.AccountWalletEvent, exist bool, err error) {
	columns := []string{columnParamAccountID, columnParamType, columnParamMeta, columnParamPayload, columnParamBalance}
	whereClause := `account_id = $1 AND deleted = false`

	result, err := r.dbAdaptor.GetDataRowWithResult(ctx, domain.AccountWalletTable, columns, whereClause, []interface{}{accountID})
	if err != nil {
		return walletRes, false, errors.Wrap(err, "failed to get wallet")
	}
	var walletPayloadBytes, walletMeta []byte
	var accountIDVar, walletType string
	var walletBalance float64

	existWallet, err := result.Scan(
		&accountIDVar,
		&walletType,
		&walletMeta,
		&walletPayloadBytes,
		&walletBalance)
	if err != nil {
		return walletRes, false, errors.Wrap(err, "failed to scan wallet row")
	}
	if !existWallet {
		return walletRes, false, nil
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
	walletRes.Payload.Type = walletType
	walletRes.Payload.Balance = walletBalance

	return walletRes, true, nil
}

// UpdateWalletBalance updates the balance of a wallet
func (r *AccountWalletRepository) UpdateWalletBalance(ctx context.Context, walletID string, balance float64) error {
	walletRes, err := r.GetWalletByID(ctx, walletID)
	if err != nil {
		return errors.Wrap(err, "failed to get wallet for balance update")
	}

	walletRes.Payload.Balance = balance
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
