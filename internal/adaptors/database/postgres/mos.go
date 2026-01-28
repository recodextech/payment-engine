package postgres

import (
	"context"
	"fmt"
	"payment-engine/internal/domain/adaptors"
	"payment-engine/internal/domain/application"
	"payment-engine/pkg/errors"
	"strconv"
	"time"

	gopostgres "github.com/HADLakmal/go-postgres"
	"github.com/recodextech/container"
)

type DBConnector struct {
	gopostgres.DatabaseReporter
	log adaptors.Logger
}

func (d *DBConnector) Init(di container.Container) error {
	d.DatabaseReporter = di.Resolve(application.ModuleSQL).(gopostgres.DatabaseReporter)
	d.log = di.Resolve(application.ModuleLogger).(adaptors.Logger).NewLog(adaptors.LoggerPrefixed(`postgres.das`))
	return nil
}

func (d *DBConnector) GetDataRowByID(ctx context.Context, table, key string) (bool, error) {
	return d.GetDataRow(ctx, table, "key", "key = $1", key)
}

// GetDataRowByAccountID checks if a row exists by account_id
func (d *DBConnector) GetDataRowByAccountID(ctx context.Context, table, account string) (bool, error) {
	return d.GetDataRow(ctx, table, "key", "account_id = $1", account)
}

// GetDataRow is a generic function to query data with custom columns and WHERE clause
func (d *DBConnector) GetDataRow(ctx context.Context, table, columns, whereClause string, args ...any) (bool, error) {
	query := fmt.Sprintf(`select %s from %s where %s`, columns, table, whereClause)
	statement, err := d.DatabaseReporter.Prepare(ctx, query)
	if err != nil {
		return false, errors.Wrap(err, table+` : query failed`)
	}
	row, err := statement.QueryRowContext(ctx, args...)
	if err != nil {
		return false, errors.Wrap(err, `query execution failed`)
	}
	var dest any
	exist, err := row.Scan(&dest)
	if err != nil {
		return false, errors.Wrap(err, `query scan failed`)
	}
	return exist, nil
}

// InsertResult is a generic function that insert data into table
func (d *DBConnector) InsertResult(ctx context.Context, table string, columns []string, args []any) error {
	columns = append(columns, `created_at`)
	args = append(args, time.Now().UTC())
	var values, columnValues string
	for index, cl := range columns {
		values += `$` + strconv.FormatInt(int64(index)+1, 10)
		columnValues += cl
		if index != len(columns)-1 {
			values += `, `
			columnValues += `, `
		}
	}
	// Use PostgreSQL INSERT ... ON CONFLICT syntax instead of UPSERT
	updateClause := d.buildUpdateClause(columns)
	query := fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s) ON CONFLICT (key) DO UPDATE SET %s`, table, columnValues, values, updateClause)
	statement, err := d.DatabaseReporter.Prepare(ctx, query)
	if err != nil {
		return errors.Wrap(err, table+` : query failed`)
	}
	_, err = statement.ExecContext(ctx, args...)
	if err != nil {
		return errors.Wrap(err, `query execution failed`)
	}
	return nil
}

// buildUpdateClause builds the UPDATE SET clause for ON CONFLICT
func (d *DBConnector) buildUpdateClause(columns []string) string {
	var updateClause string
	for _, column := range columns {
		if column == "key" {
			continue // Skip the key column as it's the conflict target
		}
		if updateClause != "" {
			updateClause += ", "
		}
		updateClause += fmt.Sprintf("%s = EXCLUDED.%s", column, column)
	}
	return updateClause
}

// GetDataRowWithResult is a generic function that returns the scanned result
func (d *DBConnector) GetDataRowWithResult(ctx context.Context, table string, columns []string, whereClause string, args []any) (gopostgres.RowInterface, error) {
	var columnValues string
	for index, cl := range columns {
		columnValues += cl
		if index != len(columns)-1 {
			columnValues += `, `
		}
	}

	query := fmt.Sprintf(`select %s from %s where %s`, columnValues, table, whereClause)
	statement, err := d.DatabaseReporter.Prepare(ctx, query)
	if err != nil {
		return nil, errors.Wrap(err, table+` : query failed`)
	}
	row, err := statement.QueryRowContext(ctx, args...)
	if err != nil {
		return nil, errors.Wrap(err, `query execution failed`)
	}
	return row, nil
}

// GetDataRowWithResult is a generic function that returns the scanned result
func (d *DBConnector) GetDataRowsWithResult(ctx context.Context, table string, columns []string, whereClause string, args []any) (gopostgres.RowsInterface, error) {
	var columnValues string
	for index, cl := range columns {
		columnValues += cl
		if index != len(columns)-1 {
			columnValues += `, `
		}
	}

	query := fmt.Sprintf(`select %s from %s where %s`, columnValues, table, whereClause)
	statement, err := d.DatabaseReporter.Prepare(ctx, query)
	if err != nil {
		return nil, errors.Wrap(err, table+` : query failed`)
	}
	row, err := statement.QueryContext(ctx, args...)
	if err != nil {
		return nil, errors.Wrap(err, `query execution failed`)
	}
	return row, nil
}
