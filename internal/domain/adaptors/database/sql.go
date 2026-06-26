package database

import (
	"context"

	gopostgres "github.com/HADLakmal/go-postgres"
)

type FixFlowDB interface {
	GetDataRowByID(ctx context.Context, table, key string) (bool, error)
	GetDataRowByAccountID(ctx context.Context, table, account string) (bool, error)
	GetDataRow(ctx context.Context, table, columns, whereClause string, args ...any) (bool, error)
	InsertResult(ctx context.Context, table string, columns []string, args []any) error
	InsertDataRow(ctx context.Context, table string, columns []string, args []any) (string, error)
	UpdateDataRow(ctx context.Context, table string, key string, updates map[string]interface{}) error
	GetDataRowWithResult(ctx context.Context, table string, columns []string, whereClause string, args []any) (gopostgres.RowInterface, error)
	GetDataRowsWithResult(ctx context.Context, table string, columns []string, whereClause string, args []any) (gopostgres.RowsInterface, error)
	// GetDataRowWithResultForUpdate is identical to GetDataRowWithResult but
	// appends FOR UPDATE so the returned row is exclusively locked within the
	// current transaction.  Must only be called from inside an InTransaction
	// callback; otherwise the lock is released immediately on statement close.
	GetDataRowWithResultForUpdate(ctx context.Context, table string, columns []string, whereClause string, args []any) (gopostgres.RowInterface, error)
}
