package models

import (
	"github.com/google/uuid"
)

// Client represents the client model.
type Client struct {
	ID uuid.UUID
}

// CreateNewClient creates a client model with a new id and the provided fields.
func CreateNewClient() *Client {
	return &Client{
		ID: uuid.New(),
	}
}
