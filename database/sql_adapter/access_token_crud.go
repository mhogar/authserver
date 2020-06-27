package sqladapter

import (
	"authserver/models"
	"errors"

	"github.com/google/uuid"
)

func (adapter *SQLAdapter) SaveAccessToken(token *models.AccessToken) error {
	return errors.New("not implemented yet")
}

func (adapter *SQLAdapter) GetAccessTokenByID(ID uuid.UUID) (*models.AccessToken, error) {
	return nil, errors.New("not implemented yet")
}

func (adapter *SQLAdapter) DeleteAccessToken(token *models.AccessToken) error {
	return errors.New("not implemented yet")
}
