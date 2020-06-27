package postgresadapter

import (
	sqladapter "authserver/database/sql_adapter"
	//import the postgres driver
	_ "github.com/lib/pq"
)

type PostgresDB struct {
	sqladapter.SQLDB
}

// CreatePostgresDB creates a PostgresDB with the supplied database key
func CreatePostgresDB(dbKey string) *PostgresDB {
	return &PostgresDB{
		SQLDB: sqladapter.SQLDB{
			SQLAdapter: sqladapter.SQLAdapter{
				DriverName:          "postgres",
				DbKey:               dbKey,
				SQLScriptRepository: PostgresScriptRepository{},
			},
		},
	}
}
