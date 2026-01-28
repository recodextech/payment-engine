package database

import (
	"context"

	gopostgres "github.com/HADLakmal/go-postgres"
)

type MOSDB interface {
	GetDataRowByID(ctx context.Context, table, key string) (bool, error)
	GetDataRowByAccountID(ctx context.Context, table, account string) (bool, error)
	GetDataRow(ctx context.Context, table, columns, whereClause string, args ...any) (bool, error)
	InsertResult(ctx context.Context, table string, columns []string, args []any) error
	GetDataRowWithResult(ctx context.Context, table string, columns []string, whereClause string, args []any) (gopostgres.RowInterface, error)
	GetDataRowsWithResult(ctx context.Context, table string, columns []string, whereClause string, args []any) (gopostgres.RowsInterface, error)
}
