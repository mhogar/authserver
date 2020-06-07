package sqladapter

import (
	"authserver/models"
	"errors"

	"github.com/google/uuid"
)

func (db *SQLAdapter) GetClientByID(ID uuid.UUID) (*models.Client, error) {
	return nil, errors.New("not implemented yet")
}
