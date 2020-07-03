package models

import (
	"github.com/google/uuid"
)

// Client ValidateError statuses.
const (
	ValidateClientValid     = iota
	ValidateClientInvalidID = iota
)

// Client represents the client model.
type Client struct {
	ID uuid.UUID
}

// ClientCRUD is an interface for performing CRUD operations on a client.
type ClientCRUD interface {
	// SaveClient saves the client and returns any errors.
	SaveClient(client *Client) error

	// GetClientByID fetches the client associated with the id.
	// If no clients are found, returns nil client. Also returns any errors.
	GetClientByID(ID uuid.UUID) (*Client, error)
}

// CreateNewClient creates a client model with a new id and the provided fields.
func CreateNewClient() *Client {
	return &Client{
		ID: uuid.New(),
	}
}

// Validate validates the client model has valid fields.
// Returns a ValidateError indicating its result.
func (c *Client) Validate() ValidateError {
	if c.ID == uuid.Nil {
		return CreateValidateError(ValidateClientInvalidID, "id cannot be nil")
	}

	return ValidateError{ValidateClientValid, nil}
}
