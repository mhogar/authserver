package sqladapter

import (
	"authserver/models"
	"errors"

	"github.com/google/uuid"
)

func (db *SQLAdapter) SaveAccessToken(token *models.AccessToken) error {
	return errors.New("not implemented yet")
}

func (db *SQLAdapter) GetAccessTokenByID(ID uuid.UUID) (*models.AccessToken, error) {
	return nil, errors.New("not implemented yet")
}

func (db *SQLAdapter) DeleteAccessToken(token *models.AccessToken) error {
	return errors.New("not implemented yet")
}
