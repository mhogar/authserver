package sqladapter

import (
	commonhelpers "authserver/helpers/common"
	"authserver/models"
	"database/sql"
	"errors"
	"fmt"
)

// SaveScope validates the scope model is valid and inserts a new row into the scope table.
// Returns any errors.
func (adapter *SQLAdapter) SaveScope(scope *models.Scope) error {
	verr := scope.Validate()
	if verr != models.ValidateScopeValid {
		return errors.New(fmt.Sprint("error validating scope model:", verr))
	}

	ctx, cancel := adapter.CreateStandardTimeoutContext()
	_, err := adapter.SQLExecuter.ExecContext(ctx, adapter.SQLDriver.SaveScopeScript(), scope.ID, scope.Name)
	cancel()

	if err != nil {
		return commonhelpers.ChainError("error executing save scope statement", err)
	}

	return nil
}

// GetScopeByName gets the row in the scope table with the matching name, and creates a new scope model using its data.
// Returns the scope and any errors.
func (adapter *SQLAdapter) GetScopeByName(name string) (*models.Scope, error) {
	ctx, cancel := adapter.CreateStandardTimeoutContext()
	rows, err := adapter.SQLExecuter.QueryContext(ctx, adapter.SQLDriver.GetScopeByNameScript(), name)
	defer cancel()

	if err != nil {
		return nil, commonhelpers.ChainError("error executing get scope by name query", err)
	}
	defer rows.Close()

	return readScopeData(rows)
}

func readScopeData(rows *sql.Rows) (*models.Scope, error) {
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
	scope := &models.Scope{}
	err := rows.Scan(&scope.ID, &scope.Name)
	if err != nil {
		return nil, commonhelpers.ChainError("error reading row", err)
	}

	return scope, nil
}
