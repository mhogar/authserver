package postgresadapter_test

import (
	"authserver/config"
	postgresadapter "authserver/database/postgres_adapter"
	"authserver/helpers"
	"database/sql"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

type MigrationCRUDTestSuite struct {
	suite.Suite
	Tx      *sql.Tx
	Adapter *postgresadapter.PostgresAdapter
}

func (suite *MigrationCRUDTestSuite) SetupSuite() {
	config.InitConfig()

	suite.Adapter = postgresadapter.CreatePostgresAdapter(viper.GetString("test_db"))

	err := suite.Adapter.OpenConnection()
	suite.Require().NoError(err)

	err = suite.Adapter.Ping()
	suite.Require().NoError(err)
}

func (suite *MigrationCRUDTestSuite) TearDownSuite() {
	suite.Adapter.CloseConnection()
}

func (suite *MigrationCRUDTestSuite) SetupTest() {
	//start a new transaction for every test
	err := suite.Adapter.BeginTransaction()
	suite.Require().NoError(err)
}

func (suite *MigrationCRUDTestSuite) TearDownTest() {
	//rollback the transaction after each test
	err := suite.Adapter.RollbackTransaction()
	suite.Require().NoError(err)
}

func (suite *MigrationCRUDTestSuite) TestCreateMigration_WithInvalidTimestamp_ReturnsError() {
	//act
	err := suite.Adapter.CreateMigration("invalid")

	//assert
	helpers.AssertError(&suite.Suite, err, "error", "model")
}

func (suite *MigrationCRUDTestSuite) TestCreateMigration_AddsRowToMigrationTable() {
	//arrange
	timestamp := "00000000000001"

	//act
	err := suite.Adapter.CreateMigration(timestamp)

	//assert
	suite.Require().NoError(err)

	migration, err := suite.Adapter.GetMigrationByTimestamp(timestamp)
	suite.Require().NoError(err)
	suite.Require().NotNil(migration)

	suite.Equal(timestamp, migration.Timestamp)
}

func TestMigrationCRUDTestSuite(t *testing.T) {
	suite.Run(t, &MigrationCRUDTestSuite{})
}
