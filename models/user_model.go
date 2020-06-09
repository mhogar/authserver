package models

import (
	"github.com/google/uuid"
)

// User ValidateError statuses.
const (
	ValidateUserValid               = iota
	ValidateUserInvalidID           = iota
	ValidateUserInvalidUsername     = iota
	ValidateUserInvalidPasswordHash = iota
)

// User represents the user model.
type User struct {
	ID           uuid.UUID
	Username     string
	PasswordHash []byte
}

// UserCRUD is an interface for performing CRUD operations on a user.
type UserCRUD interface {
	// SaveUser saves the user and returns any errors.
	SaveUser(user *User) error

	// GetUserByID fetches the user associated with the id.
	// If no users are found, returns nil user. Also returns any errors.
	GetUserByID(ID uuid.UUID) (*User, error)

	// GetUserBySessionID fetches the user associated with the session id.
	// If no users are found, returns nil user. Also returns any errors.
	GetUserBySessionID(sID uuid.UUID) (*User, error)

	// GetUserByUsername fetches the user with the matching username.
	// If no users are found, returns nil user. Also returns any errors.
	GetUserByUsername(username string) (*User, error)

	// UpdateUser updates the user and returns any errors.
	UpdateUser(user *User) error

	// DeleteUser deletes the user and returns any errors.
	DeleteUser(user *User) error
}

// CreateNewUser creates a user model with a new id and the provided fields.
func CreateNewUser(username string, passwordHash []byte) *User {
	return &User{
		ID:           uuid.New(),
		Username:     username,
		PasswordHash: passwordHash,
	}
}

// Validate validates the user model has valid fields.
// Returns a ValidateError indicating its result.
func (u *User) Validate() ValidateError {
	if u.ID == uuid.Nil {
		return CreateValidateError(ValidateUserInvalidID, "id cannot be nil")
	}

	if u.Username == "" {
		return CreateValidateError(ValidateUserInvalidUsername, "username cannot be empty")
	}

	if len(u.PasswordHash) == 0 {
		return CreateValidateError(ValidateUserInvalidPasswordHash, "password hash cannot be nil")
	}

	return ValidateError{ValidateUserValid, nil}
}
