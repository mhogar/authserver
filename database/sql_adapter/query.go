package sqladapter

import (
	"authserver/helpers"
	"context"
	"database/sql"
)

// ExecStatement executes the provided statement on its DB instance with a standard timeout context.
// Returns the result and any errors.
func (adapter *SQLAdapter) ExecStatement(stmt string, args ...interface{}) (sql.Result, error) {
	//use the active transaction if it exists
	var exec interface {
		ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	}
	if adapter.Tx != nil {
		exec = adapter.Tx
	} else {
		exec = adapter.DB
	}

	ctx, cancel := adapter.CreateStandardTimeoutContext()
	result, err := exec.ExecContext(ctx, stmt, args...)
	cancel()

	if err != nil {
		return nil, helpers.ChainError("error executing statment", err)
	}

	return result, nil
}

// ExecQuery executes the provided query on its DB instance with a standard timeout context.
// Returns the resulting rows and any errors.
func (adapter *SQLAdapter) ExecQuery(query string, args ...interface{}) (*sql.Rows, error) {
	//use the active transaction if it exists
	var q interface {
		QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	}
	if adapter.Tx != nil {
		q = adapter.Tx
	} else {
		q = adapter.DB
	}

	ctx, cancel := adapter.CreateStandardTimeoutContext()
	rows, err := q.QueryContext(ctx, query, args...)
	cancel()

	if err != nil {
		return nil, helpers.ChainError("error executing query", err)
	}

	return rows, nil
}

// ExecQueryRow executes the provided query on its DB instance with a standard timeout context.
// Returns the resulting row and any errors.
func (adapter *SQLAdapter) ExecQueryRow(query string, args ...interface{}) *sql.Row {
	//use the active transaction if it exists
	var q interface {
		QueryRowContext(context.Context, string, ...interface{}) *sql.Row
	}
	if adapter.Tx != nil {
		q = adapter.Tx
	} else {
		q = adapter.DB
	}

	ctx, cancel := adapter.CreateStandardTimeoutContext()
	row := q.QueryRowContext(ctx, query, args...)
	cancel()

	return row
}
