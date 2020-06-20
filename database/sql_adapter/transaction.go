package sqladapter

import (
	"authserver/helpers"
)

// BeginTransaction creates a new database transaction.
func (adapter *SQLAdapter) BeginTransaction() error {
	tx, err := adapter.DB.Begin()
	if err != nil {
		return helpers.ChainError("error starting transaction", err)
	}

	adapter.Tx = tx
	return nil
}

// CommitTransaction commits the active transaction.
func (adapter *SQLAdapter) CommitTransaction() error {
	err := adapter.Tx.Commit()
	if err != nil {
		return helpers.ChainError("error committing transaction", err)
	}

	//reset the transaction
	adapter.Tx = nil

	return nil
}

// RollbackTransaction rollbacks the active transaction.
func (adapter *SQLAdapter) RollbackTransaction() error {
	err := adapter.Tx.Rollback()
	if err != nil {
		return helpers.ChainError("error rolling back transaction", err)
	}

	//reset the transaction
	adapter.Tx = nil

	return nil
}
