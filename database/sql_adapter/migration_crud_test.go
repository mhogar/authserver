package sqladapter_test

import (
	"authserver/config"
	sqladapter "authserver/database/sql_adapter"
	"authserver/dependencies"
	"authserver/helpers"
	"testing"

	"github.com/stretchr/testify/suite"
)

type MigrationCRUDTestSuite struct {
	suite.Suite
	TransactionFactory *sqladapter.SQLTransactionFactory
	Tx                 *sqladapter.SQLTransaction
}

func (suite *MigrationCRUDTestSuite) SetupSuite() {
	config.InitConfig()

	//create the database and open its connection
	db := sqladapter.CreateSQLDB("integration", dependencies.ResolveSQLDriver())

	err := db.OpenConnection()
	suite.Require().NoError(err)

	err = db.Ping()
	suite.Require().NoError(err)

	suite.TransactionFactory = &sqladapter.SQLTransactionFactory{
		DB: db,
	}
}

func (suite *MigrationCRUDTestSuite) TearDownSuite() {
	suite.TransactionFactory.DB.CloseConnection()
}

func (suite *MigrationCRUDTestSuite) SetupTest() {
	//start a new transaction for every test
	tx, err := suite.TransactionFactory.CreateTransaction()
	suite.Require().NoError(err)

	suite.Tx = tx.(*sqladapter.SQLTransaction)
}

func (suite *MigrationCRUDTestSuite) TearDownTest() {
	//rollback the transaction after each test
	err := suite.Tx.RollbackTransaction()
	suite.Require().NoError(err)
}

func (suite *MigrationCRUDTestSuite) TestCreateMigration_WithInvalidTimestamp_ReturnsError() {
	//act
	err := suite.Tx.CreateMigration("invalid")

	//assert
	helpers.AssertError(&suite.Suite, err, "error", "model")
}

func (suite *MigrationCRUDTestSuite) TestGetMigrationByTimestamp_WhereTimestampNotFound_ReturnsNilMigration() {
	//act
	migration, err := suite.Tx.GetMigrationByTimestamp("DNE")

	//assert
	suite.NoError(err)
	suite.Nil(migration)
}

func (suite *MigrationCRUDTestSuite) TestGetMigrationByTimestamp_FindsMigration() {
	//arrange
	timestamp := "00000000000001"
	err := suite.Tx.CreateMigration(timestamp)
	suite.Require().NoError(err)

	//act
	migration, err := suite.Tx.GetMigrationByTimestamp(timestamp)

	//assert
	suite.NoError(err)
	suite.Require().NotNil(migration)
	suite.Equal(timestamp, migration.Timestamp)
}

func (suite *MigrationCRUDTestSuite) TestGetLatestTimestamp_WithNoLatestTimestamp_ReturnsHasLatestFalse() {
	//arrange
	ctx, cancel := suite.TransactionFactory.DB.CreateStandardTimeoutContext()
	_, err := suite.Tx.SQLExecuter.ExecContext(ctx, `DELETE FROM "migration"`)
	cancel()
	suite.Require().NoError(err)

	//act
	_, hasLatest, err := suite.Tx.GetLatestTimestamp()

	//assert
	suite.False(hasLatest)
	suite.NoError(err)
}

func (suite *MigrationCRUDTestSuite) TestGetLatestTimestamp_ReturnsLatestTimestamp() {
	//arrange
	timestamps := []string{
		"99990000000001",
		"99990000000005",
		"99990000000002",
		"99990000000003",
	}

	for _, timestamp := range timestamps {
		err := suite.Tx.CreateMigration(timestamp)
		suite.Require().NoError(err)
	}

	//act
	timestamp, hasLatest, err := suite.Tx.GetLatestTimestamp()

	//assert
	suite.Equal(timestamps[1], timestamp)
	suite.True(hasLatest)
	suite.NoError(err)
}

func (suite *MigrationCRUDTestSuite) TestDeleteMigrationByTimestamp_WithNoMigrationToDelete_ReturnsNilError() {
	//act
	err := suite.Tx.DeleteMigrationByTimestamp("DNE")

	//assert
	suite.NoError(err)
}

func (suite *MigrationCRUDTestSuite) TestDeleteMigrationByTimestamp_DeletesMigration() {
	//arrange
	timestamp := "00000000000001"
	err := suite.Tx.CreateMigration(timestamp)
	suite.Require().NoError(err)

	//act
	err = suite.Tx.DeleteMigrationByTimestamp(timestamp)

	//assert
	suite.Require().NoError(err)

	migration, err := suite.Tx.GetMigrationByTimestamp(timestamp)
	suite.NoError(err)
	suite.Nil(migration)
}

func TestMigrationCRUDTestSuite(t *testing.T) {
	suite.Run(t, &MigrationCRUDTestSuite{})
}
