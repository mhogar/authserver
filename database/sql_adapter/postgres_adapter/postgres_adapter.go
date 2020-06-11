package postgresadapter

import (
	sqladapter "authserver/database/sql_adapter"

	//import the postgres driver
	_ "github.com/lib/pq"
)

// PostgresAdapter is a sql adapter that uses the postgres driver
type PostgresAdapter struct {
	sqladapter.SQLAdapter
}

// CreatePostgresAdapter creates a PostgresAdapter with the provided db key.
func CreatePostgresAdapter(dbKey string) *PostgresAdapter {
	return &PostgresAdapter{
		SQLAdapter: sqladapter.SQLAdapter{
			DriverName: "postgres",
			DbKey:      dbKey,
		},
	}
}
