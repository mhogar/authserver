package postgresadapter

import (
	"authserver/models"
	"errors"
)

func (db *PostgresAdapter) CreateAccessToken(token *models.AccessToken) error {
	return errors.New("not implemented yet")
}
