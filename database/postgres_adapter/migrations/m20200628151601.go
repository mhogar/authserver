package migrations

import (
	postgresadapter "authserver/database/postgres_adapter"
	"authserver/helpers"
)

type m20200628151601 struct {
	DB *postgresadapter.PostgresDB
}

func (m m20200628151601) GetTimestamp() string {
	return "20200628151601"
}

func (m m20200628151601) Up() error {
	//create the user table
	ctx, cancel := m.DB.CreateStandardTimeoutContext()
	_, err := m.DB.SQLExecuter.ExecContext(ctx, m.DB.SQLScriptRepository.CreateUserTableScript())
	defer cancel()

	if err != nil {
		return helpers.ChainError("error executing create user table script", err)
	}

	return nil
}

func (m m20200628151601) Down() error {
	//drop the user table
	ctx, cancel := m.DB.CreateStandardTimeoutContext()
	_, err := m.DB.SQLExecuter.ExecContext(ctx, `DROP TABLE public."user";`)
	defer cancel()

	if err != nil {
		return helpers.ChainError("error executing drop user table script", err)
	}

	return nil
}
