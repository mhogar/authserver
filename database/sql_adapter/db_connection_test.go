package sqladapter_test

import (
	"authserver/common"
	"authserver/config"
	sqladapter "authserver/database/sql_adapter"
	"authserver/dependencies"
	"testing"

	"github.com/stretchr/testify/suite"
)

type DbConnectionTestSuite struct {
	suite.Suite
	DB *sqladapter.SQLDB
}

func (suite *DbConnectionTestSuite) SetupTest() {
	err := config.InitConfig("../..")
	suite.Require().NoError(err)

	suite.DB = sqladapter.CreateSQLDB("integration", dependencies.ResolveSQLDriver())
}

func (suite *DbConnectionTestSuite) TestOpenConnection_WhereConnectionStringIsNotFound_ReturnsError() {
	//arrange
	dbKey := "not a real dbkey"
	suite.DB.DbKey = dbKey

	//act
	err := suite.DB.OpenConnection()

	//assert
	common.AssertError(&suite.Suite, err, "no connection string", dbKey)
}

func (suite *DbConnectionTestSuite) TestCloseConnection_WithValidConnection_ReturnsNoError() {
	//arrange
	err := suite.DB.OpenConnection()
	suite.Require().NoError(err)

	//act
	err = suite.DB.CloseConnection()

	//assert
	suite.NoError(err)
	suite.Nil(suite.DB.DB)
}

func (suite *DbConnectionTestSuite) TestPing_WithValidConnection_ReturnsNoError() {
	//arrange
	err := suite.DB.OpenConnection()
	suite.Require().NoError(err)

	//act
	err = suite.DB.Ping()

	//assert
	suite.NoError(err)
}

func TestDbConnectionTestSuite(t *testing.T) {
	suite.Run(t, &DbConnectionTestSuite{})
}
