package postgresadapter

import (
	"authserver/models"
	"errors"
)

func (db *PostgresAdapter) GetScopeByName(name string) (*models.Scope, error) {
	return nil, errors.New("not implemented yet")
}
