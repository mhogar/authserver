package database

import (
	"authserver/models"
)

// Database is an interface that encapsulates the database connection interface and various model CRUD interfaces.
type Database interface {
	DBConnection
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
