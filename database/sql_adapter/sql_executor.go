package sqladapter

import (
	"context"
	"database/sql"
)

// SQLExecuter is an interface for execting sql scripts.
type SQLExecuter interface {
	// ExecContext executes the sql statement. Returns its result and any errors.
	ExecContext(ctx context.Context, stmt string, args ...interface{}) (sql.Result, error)

	// QueryContext executes the sql query. Returns the resulting rows and any errors.
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)

	// QueryRowContext executes the sql query. Should be used when extactly one resulting row is expected. Returns the resulting row.
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}
