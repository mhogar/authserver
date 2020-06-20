package postgresadapter_test

import (
	"authserver/config"
	postgresadapter "authserver/database/postgres_adapter"
	"authserver/helpers"
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

	suite.Adapter = postgresadapter.CreatePostgresAdapter(viper.GetString("test_db"))
}

func (suite *DbConnectionTestSuite) TestOpenConnection_WhereEnvironmentIsNotFound_ReturnsError() {
	//arrange
	env := "not a real environment"
	viper.Set("env", env)

	//act
	err := suite.Adapter.OpenConnection()

	//assert
	helpers.AssertError(&suite.Suite, err, "no database config", env)
}

func (suite *DbConnectionTestSuite) TestOpenConnection_WhereConnectionStringIsNotFound_ReturnsError() {
	//arrange
	dbKey := "not a real dbkey"
	suite.Adapter.DbKey = dbKey

	//act
	err := suite.Adapter.OpenConnection()

	//assert
	helpers.AssertError(&suite.Suite, err, "no connection string", dbKey)
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
