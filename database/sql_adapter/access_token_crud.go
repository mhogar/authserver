package sqladapter

import (
	commonhelpers "authserver/helpers/common"
	"authserver/models"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

// SaveAccessToken validates the access token model is valid and inserts a new row into the access_token table.
// Returns any errors.
func (adapter *SQLAdapter) SaveAccessToken(token *models.AccessToken) error {
	verr := token.Validate()
	if verr != models.ValidateAccessTokenValid {
		return errors.New(fmt.Sprint("error validating access token model:", verr))
	}

	ctx, cancel := adapter.CreateStandardTimeoutContext()
	_, err := adapter.SQLExecuter.ExecContext(ctx, adapter.SQLDriver.SaveAccessTokenScript(),
		token.ID, token.User.ID, token.Client.ID, token.Scope.ID)
	cancel()

	if err != nil {
		return commonhelpers.ChainError("error executing save access token statement", err)
	}

	return nil
}

// GetAccessTokenByID gets the row in the access_token table with the matching id, and creates a new access token model with associated models using its data.
// Returns the model and any errors.
func (adapter *SQLAdapter) GetAccessTokenByID(ID uuid.UUID) (*models.AccessToken, error) {
	ctx, cancel := adapter.CreateStandardTimeoutContext()
	rows, err := adapter.SQLExecuter.QueryContext(ctx, adapter.SQLDriver.GetAccessTokenByIdScript(), ID)
	defer cancel()

	if err != nil {
		return nil, commonhelpers.ChainError("error executing get access token by id query", err)
	}
	defer rows.Close()

	return readAccessTokenData(rows)
}

// DeleteAccessToken deletes the row in the access_token table with the matching id.
// Returns any errors.
func (adapter *SQLAdapter) DeleteAccessToken(token *models.AccessToken) error {
	ctx, cancel := adapter.CreateStandardTimeoutContext()
	_, err := adapter.SQLExecuter.ExecContext(ctx, adapter.SQLDriver.DeleteAccessTokenScript(), token.ID)
	cancel()

	if err != nil {
		return commonhelpers.ChainError("error executing delete access token statement", err)
	}

	return nil
}

func readAccessTokenData(rows *sql.Rows) (*models.AccessToken, error) {
	//check if there was a result
	if !rows.Next() {
		err := rows.Err()
		if err != nil {
			return nil, commonhelpers.ChainError("error preparing next row", err)
		}

		//return no results
		return nil, nil
	}

	token := &models.AccessToken{
		User:   &models.User{},
		Client: &models.Client{},
		Scope:  &models.Scope{},
	}

	//get the result
	err := rows.Scan(
		&token.ID,
		&token.User.ID, &token.User.Username, &token.User.PasswordHash,
		&token.Client.ID,
		&token.Scope.ID, &token.Scope.Name,
	)
	if err != nil {
		return nil, commonhelpers.ChainError("error reading row", err)
	}

	return token, nil
}
