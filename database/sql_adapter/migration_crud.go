package sqladapter

import (
	"authserver/helpers"
	"authserver/models"
)

// Setup creates the migration table if it does not already exist.
func (adapter *SQLAdapter) Setup() error {
	_, err := adapter.ExecStatement(adapter.SQLScriptRepository.GetSQLScript("CreateMigrationTable"))
	if err != nil {
		return err
	}

	return nil
}

// CreateMigration validates the given timestamp and inserts it into the migration table.
// Returns any errors.
func (adapter *SQLAdapter) CreateMigration(timestamp string) error {
	//create and validate migration model
	migration := models.CreateNewMigration(timestamp)
	verr := migration.Validate()
	if verr.Status != models.ValidateMigrationValid {
		return helpers.ChainError("error validating migration model", verr)
	}

	_, err := adapter.ExecStatement(adapter.SQLScriptRepository.GetSQLScript("SaveMigration"), migration.Timestamp)
	if err != nil {
		return err
	}

	return nil
}

// GetMigrationByTimestamp gets the row in the migration table with the matching timestamp, and creates a new migration model using its data.
// Returns the model and any errors.
func (adapter *SQLAdapter) GetMigrationByTimestamp(timestamp string) (*models.Migration, error) {
	rows, err := adapter.ExecQuery(adapter.SQLScriptRepository.GetSQLScript("GetMigrationByTimestamp"), timestamp)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	//check if there was a result
	if !rows.Next() {
		err := rows.Err()
		if err != nil {
			return nil, helpers.ChainError("error preparing next row", err)
		}

		//return no resultsz
		return nil, nil
	}

	//get the result
	migration := &models.Migration{}
	err = rows.Scan(&migration.Timestamp)
	if err != nil {
		return nil, helpers.ChainError("error reading row", err)
	}

	return migration, nil
}

// GetLatestTimestamp returns the latest timestamp of all rows in the migration table.
// If the table is empty, hasLatest will be false, else it will be true.
// Returns any errors.
func (adapter *SQLAdapter) GetLatestTimestamp() (timestamp string, hasLatest bool, err error) {
	rows, err := adapter.ExecQuery(adapter.SQLScriptRepository.GetSQLScript("GetLatestTimestamp"))
	if err != nil {
		return "", false, err
	}
	defer rows.Close()

	//check if there was a result
	if !rows.Next() {
		err := rows.Err()
		if err != nil {
			return "", false, helpers.ChainError("error preparing next row", err)
		}

		//return no results
		return "", false, nil
	}

	//get the result
	err = rows.Scan(&timestamp)
	if err != nil {
		return "", false, helpers.ChainError("error reading row", err)
	}

	return timestamp, true, nil
}

// DeleteMigrationByTimestamp deletes up to one row from the migartion table with a matching timestamp.
// Returns any errors.
func (adapter *SQLAdapter) DeleteMigrationByTimestamp(timestamp string) error {
	_, err := adapter.ExecStatement(adapter.SQLScriptRepository.GetSQLScript("DeleteMigrationByTimestamp"), timestamp)
	if err != nil {
		return err
	}

	return nil
}