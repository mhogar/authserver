package sqladapter

import (
	"context"
	"time"
)

// SQLAdapter is a SQL implementation of the Database interface.
type SQLAdapter struct {
	context    context.Context
	cancelFunc context.CancelFunc
	timeout    int

	// DriverName is the name of the sql driver.
	DriverName string

	// DbKey is the key that will be used to resolve the database's connection string.
	DbKey string

	SQLScriptRepository SQLScriptRepository
	SQLExecuter         SQLExecuter
}

// CreateStandardTimeoutContext creates a context with the timeout loaded from the database config.
// It is a child of the adapter's context and can be canceled by the adapter's cancel function.
// Returns the created context and cancel function.
func (adapter *SQLAdapter) CreateStandardTimeoutContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(adapter.context, time.Duration(adapter.timeout)*time.Millisecond)
}
