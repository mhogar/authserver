package migrations

import (
	sqladapter "authserver/database/sql_adapter"
	commonhelpers "authserver/helpers/common"
)

type m20200628151601 struct {
	DB *sqladapter.SQLDB
}

func (m m20200628151601) GetTimestamp() string {
	return "20200628151601"
}

func (m m20200628151601) Up() error {
	//create the user table
	ctx, cancel := m.DB.CreateStandardTimeoutContext()
	_, err := m.DB.SQLExecuter.ExecContext(ctx, m.DB.SQLDriver.CreateUserTableScript())
	cancel()

	if err != nil {
		return commonhelpers.ChainError("error executing create user table script", err)
	}

	//create the client table
	ctx, cancel = m.DB.CreateStandardTimeoutContext()
	_, err = m.DB.SQLExecuter.ExecContext(ctx, m.DB.SQLDriver.CreateClientTableScript())
	cancel()

	if err != nil {
		return commonhelpers.ChainError("error executing create client table script", err)
	}

	return nil
}

func (m m20200628151601) Down() error {
	//drop the user table
	ctx, cancel := m.DB.CreateStandardTimeoutContext()
	_, err := m.DB.SQLExecuter.ExecContext(ctx, m.DB.SQLDriver.DropUserTableScript())
	cancel()

	if err != nil {
		return commonhelpers.ChainError("error executing drop user table script", err)
	}

	//drop the client table
	ctx, cancel = m.DB.CreateStandardTimeoutContext()
	_, err = m.DB.SQLExecuter.ExecContext(ctx, m.DB.SQLDriver.DropClientTableScript())
	cancel()

	if err != nil {
		return commonhelpers.ChainError("error executing drop client table script", err)
	}

	return nil
}
