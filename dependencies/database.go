package dependencies

import (
	databasepkg "authserver/database"
	"sync"

	postgresadpater "authserver/database/sql_adapter/postgres_adapter"
)

var createDatabaseOnce sync.Once
var database databasepkg.Database

// ResolveDatabase resolves the Database dependency.
// Only the first call to this function will create a new Database, after which it will be retrieved from the cache.
func ResolveDatabase() databasepkg.Database {
	createDatabaseOnce.Do(func() {
		database = postgresadpater.CreatePostgresAdapter("core")
	})
	return database
}
