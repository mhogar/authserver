package sqladapter

import (
	commonhelpers "authserver/helpers/common"
	"authserver/models"
	"database/sql"

	"github.com/google/uuid"
)

// SaveUser validates the user model is valid and inserts a new row into the user table.
// Returns any errors.
func (adapter *SQLAdapter) SaveUser(user *models.User) error {
	verr := user.Validate()
	if verr.Status != models.ValidateUserValid {
		return commonhelpers.ChainError("error validating user model", verr)
	}

	ctx, cancel := adapter.CreateStandardTimeoutContext()
	_, err := adapter.SQLExecuter.ExecContext(ctx, adapter.SQLDriver.SaveUserScript(),
		user.ID, user.Username, user.PasswordHash)
	cancel()

	if err != nil {
		return commonhelpers.ChainError("error executing save user statement", err)
	}

	return nil
}

// GetUserByID gets the row in the user table with the matching id, and creates a new user model using its data.
// Returns the model and any errors.
func (adapter *SQLAdapter) GetUserByID(ID uuid.UUID) (*models.User, error) {
	ctx, cancel := adapter.CreateStandardTimeoutContext()
	rows, err := adapter.SQLExecuter.QueryContext(ctx, adapter.SQLDriver.GetUserByIdScript(), ID)
	defer cancel()

	if err != nil {
		return nil, commonhelpers.ChainError("error executing get user by id query", err)
	}
	defer rows.Close()

	return readUserData(rows)
}

// GetUserByUsername gets the row in the user table with the matching username, and creates a new user model using its data.
// Returns the model and any errors.
func (adapter *SQLAdapter) GetUserByUsername(username string) (*models.User, error) {
	ctx, cancel := adapter.CreateStandardTimeoutContext()
	rows, err := adapter.SQLExecuter.QueryContext(ctx, adapter.SQLDriver.GetUserByUsernameScript(), username)
	defer cancel()

	if err != nil {
		return nil, commonhelpers.ChainError("error executing get user by username query", err)
	}
	defer rows.Close()

	return readUserData(rows)
}

// UpdateUser validates the user model is valid and updates the row in the user table with the matching id.
// Returns any errors.
func (adapter *SQLAdapter) UpdateUser(user *models.User) error {
	verr := user.Validate()
	if verr.Status != models.ValidateUserValid {
		return commonhelpers.ChainError("error validating user model", verr)
	}

	ctx, cancel := adapter.CreateStandardTimeoutContext()
	_, err := adapter.SQLExecuter.ExecContext(ctx, adapter.SQLDriver.UpdateUserScript(),
		user.ID, user.Username, user.PasswordHash)
	cancel()

	if err != nil {
		return commonhelpers.ChainError("error executing update user statement", err)
	}

	return nil
}

// DeleteUser deletes the row in the user table with the matching id.
// Returns any errors.
func (adapter *SQLAdapter) DeleteUser(user *models.User) error {
	ctx, cancel := adapter.CreateStandardTimeoutContext()
	_, err := adapter.SQLExecuter.ExecContext(ctx, adapter.SQLDriver.DeleteUserScript(), user.ID)
	cancel()

	if err != nil {
		return commonhelpers.ChainError("error executing delete user statement", err)
	}

	return nil
}

func readUserData(rows *sql.Rows) (*models.User, error) {
	//check if there was a result
	if !rows.Next() {
		err := rows.Err()
		if err != nil {
			return nil, commonhelpers.ChainError("error preparing next row", err)
		}

		//return no results
		return nil, nil
	}

	//get the result
	user := &models.User{}
	err := rows.Scan(&user.ID, &user.Username, &user.PasswordHash)
	if err != nil {
		return nil, commonhelpers.ChainError("error reading row", err)
	}

	return user, nil
}
