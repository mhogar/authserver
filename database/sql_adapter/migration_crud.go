package sqladapter

// GetLatestTimestamp returns the latest timestamp of all rows in the migration table.
// If the table is empty or hasn't been created yet, hasLatest will be false, else it will be true.
// Returns any errors.
func (adapter *SQLAdapter) GetLatestTimestamp() (timestamp string, hasLatest bool, err error) {
	return "", false, nil
}

// CreateMigration validates the given timestamp and inserts it into the migration table.
// Returns any errors.
func (adapter *SQLAdapter) CreateMigration(timestamp string) error {
	return nil
}

// DeleteMigrationByTimestamp deletes up to one row from the migartion table with a matching timestamp.
// Returns any errors.
func (adapter *SQLAdapter) DeleteMigrationByTimestamp(timestamp string) error {
	return nil
}
