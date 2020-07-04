package sqladapter

import (
	commonhelpers "authserver/helpers/common"
	"authserver/models"
	"database/sql"

	"github.com/google/uuid"
)

// SaveClient validates the client model is valid and inserts a new row into the client table.
// Returns any errors.
func (adapter *SQLAdapter) SaveClient(client *models.Client) error {
	verr := client.Validate()
	if verr.Status != models.ValidateClientValid {
		return commonhelpers.ChainError("error validating client model", verr)
	}

	ctx, cancel := adapter.CreateStandardTimeoutContext()
	_, err := adapter.SQLExecuter.ExecContext(ctx, adapter.SQLDriver.SaveClientScript(), client.ID)
	cancel()

	if err != nil {
		return commonhelpers.ChainError("error executing save client statement", err)
	}

	return nil
}

// GetClientByID gets the row in the client table with the matching id, and creates a new client model using its data.
// Returns the model and any errors.
func (adapter *SQLAdapter) GetClientByID(ID uuid.UUID) (*models.Client, error) {
	ctx, cancel := adapter.CreateStandardTimeoutContext()
	rows, err := adapter.SQLExecuter.QueryContext(ctx, adapter.SQLDriver.GetClientByIdScript(), ID)
	defer cancel()

	if err != nil {
		return nil, commonhelpers.ChainError("error executing get client by id query", err)
	}
	defer rows.Close()

	return readClientData(rows)
}

func readClientData(rows *sql.Rows) (*models.Client, error) {
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
	client := &models.Client{}
	err := rows.Scan(&client.ID)
	if err != nil {
		return nil, commonhelpers.ChainError("error reading row", err)
	}

	return client, nil
}
