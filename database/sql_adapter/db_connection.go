package sqladapter

import (
	"authserver/config"
	"authserver/helpers"
	"context"
	"database/sql"
	"errors"

	"github.com/spf13/viper"
)

// OpenConnection opens the connection to SQL database server using the fields from the database config.
// Initializes the adapter's context and cancel function, as well as its db instance.
// Returns any errors.
func (adapter *SQLAdapter) OpenConnection() error {
	//load the database config
	env := viper.GetString("env")
	mapResult, ok := viper.GetStringMap("database")[env]
	if !ok {
		return errors.New("no database config found for environment " + env)
	}
	dbConfig := mapResult.(config.DatabaseConfig)

	//get conection string
	connectionStr, ok := dbConfig.ConnectionStrings[adapter.DbKey]
	if !ok {
		return errors.New("no connection string found for database key " + adapter.DbKey)
	}

	adapter.context, adapter.cancelFunc = context.WithCancel(context.Background())
	adapter.timeout = dbConfig.Timeout

	//connect to the db
	db, err := sql.Open(adapter.DriverName, connectionStr)
	if err != nil {
		return helpers.ChainError("error opening database connection", err)
	}
	adapter.DB = db

	return nil
}

// CloseConnection closes the connection to the SQL database server and resets its db instance.
// The adapter also calls its cancel function to cancel any child requests that may still be running.
// Niether the adapter's db instance or context should be used after calling this function.
// Returns any errors.
func (adapter *SQLAdapter) CloseConnection() error {
	err := adapter.DB.Close()
	if err != nil {
		return helpers.ChainError("error closing database connection", err)
	}

	//cancel any remaining requests that may still be running
	adapter.cancelFunc()

	//clean up resources
	adapter.DB = nil

	return nil
}

// Ping pings the SQL database server to verify it can still be reached.
// Returns an error if it cannot, or if any other errors are encountered.
func (adapter *SQLAdapter) Ping() error {
	ctx, cancel := adapter.CreateStandardTimeoutContext()
	err := adapter.DB.PingContext(ctx)
	cancel()

	if err != nil {
		return helpers.ChainError("error pinging database", err)
	}

	return nil
}
