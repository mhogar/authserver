package dependencies

import (
	databasepkg "authserver/database"
	postgresadapter "authserver/database/postgres_adapter"
	"sync"

	"github.com/spf13/viper"
)

var createDatabaseOnce sync.Once
var database databasepkg.Database

// ResolveDatabase resolves the Database dependency.
// Only the first call to this function will create a new Database, after which it will be retrieved from the cache.
func ResolveDatabase() databasepkg.Database {
	createDatabaseOnce.Do(func() {
		database = postgresadapter.CreatePostgresDB(viper.GetString("dbkey"))
	})
	return database
}
