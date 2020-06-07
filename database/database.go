package database

import (
	"authserver/models"

	"github.com/google/uuid"
)

// Database is an interface that encapsulates the other database intefaces.
type Database interface {
	DBConnection
	UserCRUD
	ClientCRUD
	ScopeCRUD
	AccessTokenCRUD
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

// UserCRUD is an interface for performing CRUD operations on a user.
type UserCRUD interface {
	// CreateUser creates a new user and returns any errors.
	CreateUser(user *models.User) error

	// GetUserByID fetches the user associated with the id.
	// If no users are found, returns nil user. Also returns any errors.
	GetUserByID(ID uuid.UUID) (*models.User, error)

	// GetUserBySessionID fetches the user associated with the session id.
	// If no users are found, returns nil user. Also returns any errors.
	GetUserBySessionID(sID uuid.UUID) (*models.User, error)

	// GetUserByUsername fetches the user with the matching username.
	// If no users are found, returns nil user. Also returns any errors.
	GetUserByUsername(username string) (*models.User, error)

	// CreateUser updates the user and returns any errors.
	UpdateUser(user *models.User) error

	// DeleteUser deletes the user associated with the provided user model.
	// Returns an error if the user could not be deleted, as well as any other errors.
	DeleteUser(user *models.User) error
}

// ClientCRUD is an interface for performing CRUD operations on a client.
type ClientCRUD interface {
	// GetClientByID fetches the client associated with the id.
	// If no clients are found, returns nil client. Also returns any errors.
	GetClientByID(ID uuid.UUID) (*models.Client, error)
}

// ScopeCRUD is an interface for performing CRUD operations on a scope.
type ScopeCRUD interface {
	// GetScopeByName fetches the scope with the matching name.
	// If no scopes are found, returns nil scope. Also returns any errors.
	GetScopeByName(name string) (*models.Scope, error)
}

// AccessTokenCRUD is an interface for performing CRUD operations on an access token.
type AccessTokenCRUD interface {
	// CreateAccessToken creates a new access token and returns any errors.
	CreateAccessToken(token *models.AccessToken) error
}
