package sqladapter

import "authserver/database"

type SQLTransactionFactory struct {
	DB *SQLDB
}

func (f SQLTransactionFactory) CreateTransaction() (database.Transaction, error) {
	tx, _ := f.DB.DB.Begin()

	//copy the adapter and set its executor to the transaction
	adapter := f.DB.SQLAdapter
	adapter.SQLExecuter = tx

	transaction := &SQLTransaction{
		SQLAdapter: adapter,
		Tx:         tx,
	}

	return transaction, nil
}
