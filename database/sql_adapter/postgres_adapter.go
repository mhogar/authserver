package sqladapter

// SQLAdapter is a PostgreSQL implementation of the Database interface.
type SQLAdapter struct {
	// DbKey is the key that will be used to resolve the database's name.
	DbKey string
}
