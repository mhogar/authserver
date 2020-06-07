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

// CreateNewUser creates a user model with new id and the provided fields.
func CreateNewUser(username string, passwordHash []byte) *User {
	return &User{
		ID:           uuid.New(),
		Username:     username,
		PasswordHash: passwordHash,
	}
}

// CreateValidateUserValid creates a ValidateError with status ValidateUserValid and nil error.
func CreateValidateUserValid() ValidateError {
	return ValidateError{ValidateUserValid, nil}
}

// Validate validates the the user model has valid fields.
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

	return CreateValidateUserValid()
}
