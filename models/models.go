package models

import "errors"

// Model is an interface for representing a model.
type Model interface {
	// Validate should check if a model if valid and return an appropriate ValidateError.
	Validate() ValidateError
}

// ValidateError is a struct for encapsulating the return value of Model's Validate method.
type ValidateError struct {
	// Status is an int that describes the type of error.
	Status int

	// error is the internal error object.
	error
}

// CreateValidateError creates a ValidateError with the provided status and an error with the provided message.
func CreateValidateError(status int, message string) ValidateError {
	return ValidateError{status, errors.New(message)}
}
