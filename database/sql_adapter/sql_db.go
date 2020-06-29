package sqladapter

import (
	"database/sql"
)

// SQLDB is a SQL implementation of the Database interface.
type SQLDB struct {
	SQLAdapter

	// DB is the sql database instance.
	DB *sql.DB
}

// CreateSQLDB creates a SQLDB with the supplied database key
func CreateSQLDB(dbKey string, SQLDriver SQLDriver) *SQLDB {
	return &SQLDB{
		SQLAdapter: SQLAdapter{
			DbKey:     dbKey,
			SQLDriver: SQLDriver,
		},
	}
}
