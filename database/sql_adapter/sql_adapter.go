package sqladapter

import (
	"context"
	"time"
)

// SQLAdapter contains methods and members common to the sql db and transaction structs.
type SQLAdapter struct {
	context    context.Context
	cancelFunc context.CancelFunc
	timeout    int

	// DbKey is the key that will be used to resolve the database's connection string.
	DbKey string

	// SQLExecuter is a dependency for executing sql scripts.
	SQLExecuter SQLExecuter

	// SQLDriver is a dependency for fetching the sql scripts and resolving the driver name.
	SQLDriver SQLDriver
}

// CreateStandardTimeoutContext creates a context with the timeout loaded from the database config.
// It is a child of the adapter's context and can be canceled by the adapter's cancel function.
// Returns the created context and cancel function.
func (adapter *SQLAdapter) CreateStandardTimeoutContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(adapter.context, time.Duration(adapter.timeout)*time.Millisecond)
}
