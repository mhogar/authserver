package sqladapter

import (
	"authserver/models"
	"errors"
)

func (db *SQLAdapter) SaveAccessToken(token *models.AccessToken) error {
	return errors.New("not implemented yet")
}
