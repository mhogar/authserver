package postgresadapter_test

import (
	"authserver/config"
	"authserver/database"
	postgresadapter "authserver/database/postgres_adapter"
	sqladapter "authserver/database/sql_adapter"
	"authserver/helpers"
	"testing"

	"github.com/stretchr/testify/suite"
)

type MigrationCRUDTestSuite struct {
	suite.Suite
	TransactionFactory *sqladapter.SQLTransactionFactory
	Transaction        database.Transaction
}

func (suite *MigrationCRUDTestSuite) SetupSuite() {
	config.InitConfig()

	//create the database and open its connection
	db := postgresadapter.CreatePostgresDB("integration")

	err := db.OpenConnection()
	suite.Require().NoError(err)

	err = db.Ping()
	suite.Require().NoError(err)

	suite.TransactionFactory = &sqladapter.SQLTransactionFactory{
		DB: &db.SQLDB,
	}
}

func (suite *MigrationCRUDTestSuite) TearDownSuite() {
	suite.TransactionFactory.DB.CloseConnection()
}

func (suite *MigrationCRUDTestSuite) SetupTest() {
	//start a new transaction for every test
	tx, err := suite.TransactionFactory.CreateTransaction()
	suite.Require().NoError(err)

	suite.Transaction = tx
}

func (suite *MigrationCRUDTestSuite) TearDownTest() {
	//rollback the transaction after each test
	err := suite.Transaction.RollbackTransaction()
	suite.Require().NoError(err)
}

func (suite *MigrationCRUDTestSuite) TestCreateMigration_WithInvalidTimestamp_ReturnsError() {
	//act
	err := suite.Transaction.CreateMigration("invalid")

	//assert
	helpers.AssertError(&suite.Suite, err, "error", "model")
}

func (suite *MigrationCRUDTestSuite) TestCreateMigration_AddsRowToMigrationTable() {
	//arrange
	timestamp := "00000000000001"

	//act
	err := suite.Transaction.CreateMigration(timestamp)

	//assert
	suite.Require().NoError(err)

	migration, err := suite.Transaction.GetMigrationByTimestamp(timestamp)
	suite.Require().NoError(err)
	suite.Require().NotNil(migration)

	suite.Equal(timestamp, migration.Timestamp)
}

func TestMigrationCRUDTestSuite(t *testing.T) {
	suite.Run(t, &MigrationCRUDTestSuite{})
}
