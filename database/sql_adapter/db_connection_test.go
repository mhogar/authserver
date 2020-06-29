package sqladapter_test

import (
	"authserver/config"
	sqladapter "authserver/database/sql_adapter"
	"authserver/dependencies"
	"authserver/helpers"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

type DbConnectionTestSuite struct {
	suite.Suite
	DB *sqladapter.SQLDB
}

func (suite *DbConnectionTestSuite) SetupTest() {
	viper.Reset()
	config.InitConfig()

	suite.DB = sqladapter.CreateSQLDB("integration", dependencies.ResolveSQLDriver())
}

func (suite *DbConnectionTestSuite) TestOpenConnection_WhereEnvironmentIsNotFound_ReturnsError() {
	//arrange
	env := "not a real environment"
	viper.Set("env", env)

	//act
	err := suite.DB.OpenConnection()

	//assert
	helpers.AssertError(&suite.Suite, err, "no database config", env)
}

func (suite *DbConnectionTestSuite) TestOpenConnection_WhereConnectionStringIsNotFound_ReturnsError() {
	//arrange
	dbKey := "not a real dbkey"
	suite.DB.DbKey = dbKey

	//act
	err := suite.DB.OpenConnection()

	//assert
	helpers.AssertError(&suite.Suite, err, "no connection string", dbKey)
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
