package sqladapter_test

import (
	"authserver/config"
	postgresadapter "authserver/database/sql_adapter/postgres_adapter"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

type DbConnectionTestSuite struct {
	suite.Suite
	Adapter *postgresadapter.PostgresAdapter
}

func (suite *DbConnectionTestSuite) SetupTest() {
	viper.Reset()
	config.InitConfig()

	suite.Adapter = postgresadapter.CreatePostgresAdapter("integration")
}

func (suite *DbConnectionTestSuite) TestOpenConnection_WhereEnvironmentIsNotFound_ReturnsError() {
	//arrange
	env := "not a real environment"
	viper.Set("env", env)

	//act
	err := suite.Adapter.OpenConnection()

	//assert
	suite.Require().Error(err)
	suite.Contains(err.Error(), env)
}

func (suite *DbConnectionTestSuite) TestCloseConnection_WithValidConnection_ReturnsNoError() {
	//arrange
	err := suite.Adapter.OpenConnection()
	suite.Require().NoError(err)

	//act
	err = suite.Adapter.CloseConnection()

	//assert
	suite.NoError(err)
	suite.Nil(suite.Adapter.DB)
}

func (suite *DbConnectionTestSuite) TestPing_WithValidConnection_ReturnsNoError() {
	//arrange
	err := suite.Adapter.OpenConnection()
	suite.Require().NoError(err)

	//act
	err = suite.Adapter.Ping()

	//assert
	suite.NoError(err)
}

func TestDbConnectionTestSuite(t *testing.T) {
	suite.Run(t, &DbConnectionTestSuite{})
}
