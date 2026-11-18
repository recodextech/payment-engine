package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"strconv"
	"strings"

	"payment-engine/pkg/errors"

	gopostgres "github.com/HADLakmal/go-postgres"
	"github.com/recodextech/api-definitions/pkg/uuid"
)

// sqlRunner runs raw statements on the transaction in ctx when there is one.
type sqlRunner struct {
	db gopostgres.DatabaseReporter
}

func (r sqlRunner) exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	stmt, err := r.db.Prepare(ctx, query)
	if err != nil {
		return nil, errors.Wrap(err, "prepare failed")
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, args...)
	if err != nil {
		return nil, errors.Wrap(err, "exec failed")
	}
	return res, nil
}

// queryRow scans the first row into dest and reports whether a row existed.
func (r sqlRunner) queryRow(ctx context.Context, query string, args []any, dest ...any) (bool, error) {
	stmt, err := r.db.Prepare(ctx, query)
	if err != nil {
		return false, errors.Wrap(err, "prepare failed")
	}
	defer stmt.Close()
	row, err := stmt.QueryRowContext(ctx, args...)
	if err != nil {
		return false, errors.Wrap(err, "query failed")
	}
	exist, err := row.Scan(dest...)
	if err != nil {
		return false, errors.Wrap(err, "scan failed")
	}
	return exist, nil
}

// queryRows calls scan for every row; rows are closed before it returns.
func (r sqlRunner) queryRows(ctx context.Context, query string, args []any, scan func(gopostgres.RowsInterface) error) error {
	stmt, err := r.db.Prepare(ctx, query)
	if err != nil {
		return errors.Wrap(err, "prepare failed")
	}
	defer stmt.Close()
	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return errors.Wrap(err, "query failed")
	}
	defer rows.Close()
	for rows.Next() {
		if err := scan(rows); err != nil {
			return errors.Wrap(err, "failed to scan row")
		}
	}
	return nil
}

func toJSON(v any) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", errors.Wrap(err, "marshal failed")
	}
	return string(b), nil
}

func parseUUID(s string) uuid.UUID {
	id, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil
	}
	return id
}

func itoa(i int) string {
	return strconv.Itoa(i)
}

// columnList joins column names for a SELECT or INSERT column list.
func columnList(cols []string) string {
	return strings.Join(cols, ", ")
}

// placeholders returns "$1, $2, ..., $n".
func placeholders(n int) string {
	p := make([]string, n)
	for i := range p {
		p[i] = "$" + itoa(i+1)
	}
	return strings.Join(p, ", ")
}

// insertQuery builds "INSERT INTO table (cols...) VALUES ($1, ...)".
func insertQuery(table string, cols []string) string {
	return `INSERT INTO ` + table + ` (` + columnList(cols) + `) VALUES (` + placeholders(len(cols)) + `)`
}

// setClause builds "a = $start, b = $start+1, ..." for an UPDATE.
func setClause(start int, cols ...string) string {
	p := make([]string, len(cols))
	for i, c := range cols {
		p[i] = c + " = $" + itoa(start+i)
	}
	return strings.Join(p, ", ")
}

// excludedSet builds "a = EXCLUDED.a, ..." for an ON CONFLICT DO UPDATE.
func excludedSet(cols ...string) string {
	p := make([]string, len(cols))
	for i, c := range cols {
		p[i] = c + " = EXCLUDED." + c
	}
	return strings.Join(p, ", ")
}
