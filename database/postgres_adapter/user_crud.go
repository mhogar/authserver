package postgresadapter

import (
	"authserver/helpers"
	"authserver/models"
	"errors"

	"github.com/google/uuid"
)

func (db *PostgresAdapter) CreateUser(user *models.User) error {
	verr := user.Validate()
	if verr.Status != models.ValidateUserValid {
		return helpers.ChainError("error validating user model", verr)
	}

	return errors.New("not implemented yet")
}

func (db *PostgresAdapter) GetUserByID(id uuid.UUID) (*models.User, error) {
	return nil, errors.New("not implemented yet")
}

func (db *PostgresAdapter) GetUserByUsername(email string) (*models.User, error) {
	return nil, errors.New("not implemented yet")
}

func (db *PostgresAdapter) GetUserBySessionID(sID uuid.UUID) (*models.User, error) {
	return nil, errors.New("not implemented yet")
}

func (db *PostgresAdapter) UpdateUser(user *models.User) error {
	return errors.New("not implemented yet")
}

func (db *PostgresAdapter) DeleteUser(user *models.User) error {
	return errors.New("not implemented yet")
}
