package sqladapter

import (
	"authserver/models"
	"errors"
)

func (db *SQLAdapter) GetScopeByName(name string) (*models.Scope, error) {
	return nil, errors.New("not implemented yet")
}
