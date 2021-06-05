package database

import (
	"authserver/models"
)

// CRUDOperations is an interface that encapsulates various model CRUD interfaces.
type CRUDOperations interface {
	models.MigrationCRUD
	models.UserCRUD
	models.ClientCRUD
	models.ScopeCRUD
	models.AccessTokenCRUD
}

// DBConnection is an interface for controlling the connection to the database.
type DBConnection interface {
	// OpenConnection opens the connection to the database. Returns any errors.
	OpenConnection() error

	// CloseConnection closes the connection to the database and cleanup associated resources. Returns any errors.
	CloseConnection() error

	// Ping pings the database to verify it can be reached.
	// Returns an error if the database can't be reached or if any other errors occur.
	Ping() error
}

// TransactionOperations is an interface for performing transaction operations.
type TransactionOperations interface {
	// CommitTransaction commits the transaction.
	CommitTransaction() error

	// RollbackTransaction rollbacks the transaction.
	RollbackTransaction()
}

// Database is an interface that encapsulates the database connection and CRUD operations interfaces.
type Database interface {
	DBConnection
	CRUDOperations
}

// Transaction is an interface that encapsulates the transaction operations and CRUD operations interfaces.
type Transaction interface {
	TransactionOperations
	CRUDOperations
}

// TransactionFactory is an interface for creating new transactions.
type TransactionFactory interface {
	// CreateTransaction creates a new transaction and returns and errors.
	CreateTransaction() (Transaction, error)
}
