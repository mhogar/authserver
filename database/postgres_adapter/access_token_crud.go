package postgresadapter

import (
	"authserver/models"
	"errors"
)

func (db *PostgresAdapter) SaveAccessToken(token *models.AccessToken) error {
	return errors.New("not implemented yet")
}
