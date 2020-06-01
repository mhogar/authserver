package dependencies

import (
	databasepkg "authserver/database"
	postgresadapter "authserver/database/postgres_adapter"
	"sync"
)

var createDatabaseOnce sync.Once
var database databasepkg.Database

// ResolveDatabase resolves the Database dependency.
// Only the first call to this function will create a new Database, after which it will be retrieved from the cache.
func ResolveDatabase() databasepkg.Database {
	createDatabaseOnce.Do(func() {
		database = &postgresadapter.PostgresAdapter{
			DbKey: "core",
		}
	})
	return database
}
