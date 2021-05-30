package migrations

import (
	"authserver/common"
	sqladapter "authserver/database/sql_adapter"
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
		return common.ChainError("error executing create user table script", err)
	}

	//create the client table
	ctx, cancel = m.DB.CreateStandardTimeoutContext()
	_, err = m.DB.SQLExecuter.ExecContext(ctx, m.DB.SQLDriver.CreateClientTableScript())
	cancel()

	if err != nil {
		return common.ChainError("error executing create client table script", err)
	}

	//create the scope table
	ctx, cancel = m.DB.CreateStandardTimeoutContext()
	_, err = m.DB.SQLExecuter.ExecContext(ctx, m.DB.SQLDriver.CreateScopeTableScript())
	cancel()

	if err != nil {
		return common.ChainError("error executing create scope table script", err)
	}

	//create the access_token table
	ctx, cancel = m.DB.CreateStandardTimeoutContext()
	_, err = m.DB.SQLExecuter.ExecContext(ctx, m.DB.SQLDriver.CreateAccessTokenTableScript())
	cancel()

	if err != nil {
		return common.ChainError("error executing create access token table script", err)
	}

	return nil
}

func (m m20200628151601) Down() error {
	//drop the access token table
	ctx, cancel := m.DB.CreateStandardTimeoutContext()
	_, err := m.DB.SQLExecuter.ExecContext(ctx, m.DB.SQLDriver.DropAccessTokenTableScript())
	cancel()

	if err != nil {
		return common.ChainError("error executing drop access token table script", err)
	}

	//drop the scope table
	ctx, cancel = m.DB.CreateStandardTimeoutContext()
	_, err = m.DB.SQLExecuter.ExecContext(ctx, m.DB.SQLDriver.DropScopeTableScript())
	cancel()

	if err != nil {
		return common.ChainError("error executing drop scope table script", err)
	}

	//drop the client table
	ctx, cancel = m.DB.CreateStandardTimeoutContext()
	_, err = m.DB.SQLExecuter.ExecContext(ctx, m.DB.SQLDriver.DropClientTableScript())
	cancel()

	if err != nil {
		return common.ChainError("error executing drop client table script", err)
	}

	//drop the user table
	ctx, cancel = m.DB.CreateStandardTimeoutContext()
	_, err = m.DB.SQLExecuter.ExecContext(ctx, m.DB.SQLDriver.DropUserTableScript())
	cancel()

	if err != nil {
		return common.ChainError("error executing drop user table script", err)
	}

	return nil
}
