package postgresadapter

import (
	"authserver/helpers"
	"authserver/models"
)

// Setup creates the migration table if it does not already exist.
func (adapter *PostgresAdapter) Setup() error {
	stmt := `
		CREATE TABLE IF NOT EXISTS public.migration (
			"timestamp" varchar(14) NOT NULL,
			CONSTRAINT migration_pk PRIMARY KEY ("timestamp")
		);
	`

	_, err := adapter.ExecStatement(stmt)
	if err != nil {
		return err
	}

	return nil
}

// CreateMigration validates the given timestamp and inserts it into the migration table.
// Returns any errors.
func (adapter *PostgresAdapter) CreateMigration(timestamp string) error {
	//create and validate migration model
	migration := models.CreateNewMigration(timestamp)
	verr := migration.Validate()
	if verr.Status != models.ValidateMigrationValid {
		return helpers.ChainError("error validating migration model", verr)
	}

	stmt := `
		INSERT INTO migration ("timestamp")
			VALUES ($1)
	`

	_, err := adapter.ExecStatement(stmt, migration.Timestamp)
	if err != nil {
		return err
	}

	return nil
}

// GetMigrationByTimestamp gets the row in the migration table with the matching timestamp, and creates a new migration model using its data.
// Returns the model and any errors.
func (adapter *PostgresAdapter) GetMigrationByTimestamp(timestamp string) (*models.Migration, error) {
	query := `
		SELECT m."timestamp" 
			FROM migration m 
			WHERE m."timestamp" = $1
	`

	rows, err := adapter.ExecQuery(query, timestamp)
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

		//return no results
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
func (adapter *PostgresAdapter) GetLatestTimestamp() (timestamp string, hasLatest bool, err error) {
	query := `
		SELECT m."timestamp" FROM migration m
			ORDER BY m."timestamp" DESC
			LIMIT 1
	`

	rows, err := adapter.ExecQuery(query)
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
func (adapter *PostgresAdapter) DeleteMigrationByTimestamp(timestamp string) error {
	stmt := `
		DELETE FROM migration
			WHERE "timestamp" = $1
	`

	_, err := adapter.ExecStatement(stmt, timestamp)
	if err != nil {
		return err
	}

	return nil
}
