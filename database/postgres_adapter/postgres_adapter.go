package postgresadapter

import (
	sqladapter "authserver/database/sql_adapter"
	//import the postgres driver
	_ "github.com/lib/pq"
)

// PostgresAdapter is a SQL implementation of the Database interface.
type PostgresAdapter struct {
	sqladapter.SQLAdapter
}

// CreatePostgresAdapter creates a PostgresAdapter with the supplied database key
func CreatePostgresAdapter(dbKey string) *PostgresAdapter {
	return &PostgresAdapter{
		SQLAdapter: sqladapter.SQLAdapter{
			DriverName: "postgres",
			DbKey:      dbKey,
		},
	}
}
