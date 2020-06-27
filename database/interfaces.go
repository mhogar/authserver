package database

import (
	"authserver/models"
)

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

type TransactionOperations interface {
	CommitTransaction() error
	RollbackTransaction() error
}

// Database is an interface that encapsulates the database connection interface and various model CRUD interfaces.
type Database interface {
	DBConnection
	CRUDOperations
}

type Transaction interface {
	TransactionOperations
	CRUDOperations
}

type TransactionFactory interface {
	CreateTransaction() (Transaction, error)
}
