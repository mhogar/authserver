package sqladapter

import (
	"authserver/helpers"
	"authserver/models"
	"errors"

	"github.com/google/uuid"
)

func (adapter *SQLAdapter) SaveUser(user *models.User) error {
	verr := user.Validate()
	if verr.Status != models.ValidateUserValid {
		return helpers.ChainError("error validating user model", verr)
	}

	return errors.New("not implemented yet")
}

func (adapter *SQLAdapter) GetUserByID(id uuid.UUID) (*models.User, error) {
	return nil, errors.New("not implemented yet")
}

func (adapter *SQLAdapter) GetUserByUsername(email string) (*models.User, error) {
	return nil, errors.New("not implemented yet")
}

func (adapter *SQLAdapter) UpdateUser(user *models.User) error {
	return errors.New("not implemented yet")
}

func (adapter *SQLAdapter) DeleteUser(user *models.User) error {
	return errors.New("not implemented yet")
}
