package models

import (
	"github.com/google/uuid"
)

// Client represents the client model.
type Client struct {
	ID uuid.UUID
}

// ClientCRUD is an interface for performing CRUD operations on a client.
type ClientCRUD interface {
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
