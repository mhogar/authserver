package sqladapter

import (
	"authserver/helpers"
	"authserver/models"
	"errors"

	"github.com/google/uuid"
)

// SaveUser validates the user model is valid and inserts a new row into the user table.
// Returns any errors.
func (adapter *SQLAdapter) SaveUser(user *models.User) error {
	verr := user.Validate()
	if verr.Status != models.ValidateUserValid {
		return helpers.ChainError("error validating user model", verr)
	}

	ctx, cancel := adapter.CreateStandardTimeoutContext()
	_, err := adapter.SQLExecuter.ExecContext(ctx, adapter.SQLScriptRepository.SaveUserScript(),
		user.ID, user.Username, user.PasswordHash)
	cancel()

	if err != nil {
		return helpers.ChainError("error executing save user statement", err)
	}

	return nil
}

// GetUserByID gets the row in the user table with the matching id, and creates a new user model using its data.
// Returns the model and any errors.
func (adapter *SQLAdapter) GetUserByID(id uuid.UUID) (*models.User, error) {
	ctx, cancel := adapter.CreateStandardTimeoutContext()
	rows, err := adapter.SQLExecuter.QueryContext(ctx, adapter.SQLScriptRepository.GetUserByIDScript(), id)
	defer cancel()

	if err != nil {
		return nil, helpers.ChainError("error executing get user by id query", err)
	}
	defer rows.Close()

	//check if there was a result
	if !rows.Next() {
		err := rows.Err()
		if err != nil {
			return nil, helpers.ChainError("error preparing next row", err)
		}

		//return no results
		return nil, nil
	}

	//get the result
	user := &models.User{}
	err = rows.Scan(&user.ID, &user.Username, &user.PasswordHash)
	if err != nil {
		return nil, helpers.ChainError("error reading row", err)
	}

	return user, nil
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
