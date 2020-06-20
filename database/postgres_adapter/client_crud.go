package postgresadapter

import (
	"authserver/models"
	"errors"

	"github.com/google/uuid"
)

func (db *PostgresAdapter) GetClientByID(ID uuid.UUID) (*models.Client, error) {
	return nil, errors.New("not implemented yet")
}
