package sqladapter

import (
	"authserver/helpers"
	"database/sql"
)

// ExecStatement executes the provided statement on its DB instance with a standard timeout context.
// Returns the result and any errors.
func (adapter *SQLAdapter) ExecStatement(stmt string, args ...interface{}) (sql.Result, error) {
	ctx, cancel := adapter.CreateStandardTimeoutContext()
	result, err := adapter.SQLExecuter.ExecContext(ctx, stmt, args...)
	cancel()

	if err != nil {
		return nil, helpers.ChainError("error executing statment", err)
	}

	return result, nil
}

// ExecQuery executes the provided query on its DB instance with a standard timeout context.
// Returns the resulting rows and any errors.
func (adapter *SQLAdapter) ExecQuery(query string, args ...interface{}) (*sql.Rows, error) {
	ctx, cancel := adapter.CreateStandardTimeoutContext()
	rows, err := adapter.SQLExecuter.QueryContext(ctx, query, args...)
	cancel()

	if err != nil {
		return nil, helpers.ChainError("error executing query", err)
	}

	return rows, nil
}

// ExecQueryRow executes the provided query on its DB instance with a standard timeout context.
// Returns the resulting row and any errors.
func (adapter *SQLAdapter) ExecQueryRow(query string, args ...interface{}) *sql.Row {
	ctx, cancel := adapter.CreateStandardTimeoutContext()
	row := adapter.SQLExecuter.QueryRowContext(ctx, query, args...)
	cancel()

	return row
}
