package sqladapter

import (
	"authserver/database"
	"authserver/helpers"
	"database/sql"
)

// SQLTransaction is a SQL implementation of the Transaction interface.
type SQLTransaction struct {
	SQLAdapter

	// TX is the sql transaction instance.
	Tx *sql.Tx
}

// SQLTransactionFactory is a SQL implementation of the TransactionFactory interface.
type SQLTransactionFactory struct {
	// DB is the sql db instance to create transactions for.
	DB *SQLDB
}

// CommitTransaction commits the sql transaction's transaction instance.
// Returns any errors.
func (tx *SQLTransaction) CommitTransaction() error {
	return tx.Tx.Commit()
}

// RollbackTransaction rollbacks the sql transaction's transaction instance.
// Returns any errors.
func (tx *SQLTransaction) RollbackTransaction() error {
	return tx.Tx.Rollback()
}

// CreateTransaction creates a new sql transaction. Returns any errors.
func (f SQLTransactionFactory) CreateTransaction() (database.Transaction, error) {
	tx, err := f.DB.DB.Begin()
	if err != nil {
		return nil, helpers.ChainError("error beginning transaction", err)
	}

	//copy the adapter and set its executor to the transaction
	adapter := f.DB.SQLAdapter
	adapter.SQLExecuter = tx

	transaction := &SQLTransaction{
		SQLAdapter: adapter,
		Tx:         tx,
	}

	return transaction, nil
}
