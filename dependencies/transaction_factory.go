package dependencies

import (
	databasepkg "authserver/database"
	postgresadapter "authserver/database/postgres_adapter"
	sqladapter "authserver/database/sql_adapter"
	"sync"
)

var transactionFactoryOnce sync.Once
var transactionFactory databasepkg.TransactionFactory

// ResolveTransactionFactory resolves the TransactionFactory dependency.
// Only the first call to this function will create a new TransactionFactory, after which it will be retrieved from the cache.
func ResolveTransactionFactory() databasepkg.TransactionFactory {
	transactionFactoryOnce.Do(func() {
		transactionFactory = sqladapter.SQLTransactionFactory{
			DB: &ResolveDatabase().(*postgresadapter.PostgresDB).SQLDB,
		}
	})
	return transactionFactory
}
