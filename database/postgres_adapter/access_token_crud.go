package postgresadapter

import (
	"authserver/models"
	"errors"

	"github.com/google/uuid"
)

func (db *PostgresAdapter) SaveAccessToken(token *models.AccessToken) error {
	return errors.New("not implemented yet")
}

func (db *PostgresAdapter) GetAccessTokenByID(ID uuid.UUID) (*models.AccessToken, error) {
	return nil, errors.New("not implemented yet")
}

func (db *PostgresAdapter) DeleteAccessToken(token *models.AccessToken) error {
	return errors.New("not implemented yet")
}
