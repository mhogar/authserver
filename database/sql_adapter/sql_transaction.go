package sqladapter

import "database/sql"

type SQLTransaction struct {
	SQLAdapter
	Tx *sql.Tx
}

func (tx *SQLTransaction) CommitTransaction() error {
	return tx.Tx.Commit()
}

func (tx *SQLTransaction) RollbackTransaction() error {
	return tx.Tx.Rollback()
}
