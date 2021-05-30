package dependencies

import (
	databasepkg "authserver/database"
	sqladapter "authserver/database/sql_adapter"
	"sync"
)

var createTransactionFactoryOnce sync.Once
var transactionFactory databasepkg.TransactionFactory

// ResolveTransactionFactory resolves the TransactionFactory dependency.
// Only the first call to this function will create a new TransactionFactory, after which it will be retrieved from memory.
func ResolveTransactionFactory() databasepkg.TransactionFactory {
	createTransactionFactoryOnce.Do(func() {
		transactionFactory = sqladapter.SQLTransactionFactory{
			DB: ResolveDatabase().(*sqladapter.SQLDB),
		}
	})
	return transactionFactory
}
