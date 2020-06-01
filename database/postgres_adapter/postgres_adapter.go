package postgresadapter

// PostgresAdapter is a PostgreSQL implementation of the Database interface.
type PostgresAdapter struct {
	// DbKey is the key that will be used to resolve the database's name.
	DbKey string
}
