package postgresadapter

import (
	"authserver/helpers"
	"authserver/models"
	"errors"
)

func (db *PostgresAdapter) CreateSession(session *models.Session) error {
	verr := session.Validate()
	if verr.Status != models.ValidateSessionValid {
		return helpers.ChainError("error validating session model", verr)
	}

	return errors.New("not implemented yet")
}
