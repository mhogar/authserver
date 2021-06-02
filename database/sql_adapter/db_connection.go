package sqladapter

import (
	"authserver/common"
	"authserver/config"
	"context"
	"database/sql"
	"errors"

	"github.com/spf13/viper"
)

// OpenConnection opens the connection to SQL database server using the fields from the database config.
// Initializes the adapter's context and cancel function, as well as its db instance.
// Returns any errors.
func (DB *SQLDB) OpenConnection() error {
	//load the database config
	dbConfig := viper.Get("database").(config.DatabaseConfig)

	//get conection string
	connectionStr, ok := dbConfig.ConnectionStrings[DB.DbKey]
	if !ok {
		return errors.New("no connection string found for database key " + DB.DbKey)
	}

	DB.context, DB.cancelFunc = context.WithCancel(context.Background())
	DB.timeout = dbConfig.Timeout

	//connect to the db
	db, err := sql.Open(DB.SQLDriver.GetDriverName(), connectionStr)
	if err != nil {
		return common.ChainError("error opening database connection", err)
	}

	DB.DB = db
	DB.SQLExecuter = db

	return nil
}

// CloseConnection closes the connection to the SQL database server and resets its db instance.
// The adapter also calls its cancel function to cancel any child requests that may still be running.
// Niether the adapter's db instance or context should be used after calling this function.
// Returns any errors.
func (DB *SQLDB) CloseConnection() error {
	err := DB.DB.Close()
	if err != nil {
		return common.ChainError("error closing database connection", err)
	}

	//cancel any remaining requests that may still be running
	DB.cancelFunc()

	//clean up resources
	DB.DB = nil

	return nil
}

// Ping pings the SQL database server to verify it can still be reached.
// Returns an error if it cannot, or if any other errors are encountered.
func (DB *SQLDB) Ping() error {
	ctx, cancel := DB.CreateStandardTimeoutContext()
	err := DB.DB.PingContext(ctx)
	cancel()

	if err != nil {
		return common.ChainError("error pinging database", err)
	}

	return nil
}
