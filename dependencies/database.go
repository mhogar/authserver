package dependencies

import (
	databasepkg "authserver/database"
	sqladapter "authserver/database/sql_adapter"
	"sync"

	"github.com/spf13/viper"
)

var createDatabaseOnce sync.Once
var database databasepkg.Database

// ResolveDatabase resolves the Database dependency.
// Only the first call to this function will create a new Database, after which it will be retrieved from memory.
func ResolveDatabase() databasepkg.Database {
	createDatabaseOnce.Do(func() {
		database = sqladapter.CreateSQLDB(viper.GetString("db_key"), ResolveSQLDriver())
	})
	return database
}
