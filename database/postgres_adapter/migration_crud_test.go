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

func (suite *MigrationCRUDTestSuite) TestGetMigrationByTimestamp_WhereTimestampNotFound_ReturnsNilMigration() {
	//act
	migration, err := suite.Transaction.GetMigrationByTimestamp("DNE")

	//assert
	suite.NoError(err)
	suite.Nil(migration)
}

func (suite *MigrationCRUDTestSuite) TestGetMigrationByTimestamp_FindsMigration() {
	//arrange
	timestamp := "00000000000001"
	err := suite.Transaction.CreateMigration(timestamp)
	suite.Require().NoError(err)

	//act
	migration, err := suite.Transaction.GetMigrationByTimestamp(timestamp)

	//assert
	suite.NoError(err)
	suite.Require().NotNil(migration)
	suite.Equal(timestamp, migration.Timestamp)
}

func (suite *MigrationCRUDTestSuite) TestGetLatestTimestamp_WithNoLatestTimestamp_ReturnsHasLatestFalse() {
	//act
	_, hasLatest, err := suite.Transaction.GetLatestTimestamp()

	//assert
	suite.False(hasLatest)
	suite.NoError(err)
}

func (suite *MigrationCRUDTestSuite) TestGetLatestTimestamp_ReturnsLatestTimestamp() {
	//arrange
	timestamps := []string{
		"00000000000001",
		"00000000000005",
		"00000000000002",
		"00000000000003",
	}

	for _, timestamp := range timestamps {
		err := suite.Transaction.CreateMigration(timestamp)
		suite.Require().NoError(err)
	}

	//act
	timestamp, hasLatest, err := suite.Transaction.GetLatestTimestamp()

	//assert
	suite.Equal(timestamps[1], timestamp)
	suite.True(hasLatest)
	suite.NoError(err)
}

func (suite *MigrationCRUDTestSuite) TestDeleteMigrationByTimestamp_WithNoMigrationToDelete_ReturnsNilError() {
	//act
	err := suite.Transaction.DeleteMigrationByTimestamp("DNE")

	//assert
	suite.NoError(err)
}

func (suite *MigrationCRUDTestSuite) TestDeleteMigrationByTimestamp_DeletesMigration() {
	//arrange
	timestamp := "00000000000001"
	err := suite.Transaction.CreateMigration(timestamp)
	suite.Require().NoError(err)

	//act
	err = suite.Transaction.DeleteMigrationByTimestamp(timestamp)

	//assert
	suite.Require().NoError(err)

	migration, err := suite.Transaction.GetMigrationByTimestamp(timestamp)
	suite.NoError(err)
	suite.Nil(migration)
}

func TestMigrationCRUDTestSuite(t *testing.T) {
	suite.Run(t, &MigrationCRUDTestSuite{})
}
