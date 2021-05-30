package sqladapter_test

import (
	"authserver/common"
	"authserver/config"
	sqladapter "authserver/database/sql_adapter"
	"authserver/dependencies"
	"authserver/models"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ScopeCRUDTestSuite struct {
	suite.Suite
	TransactionFactory *sqladapter.SQLTransactionFactory
	Tx                 *sqladapter.SQLTransaction
}

func (suite *ScopeCRUDTestSuite) SetupSuite() {
	err := config.InitConfig("../..")
	suite.Require().NoError(err)

	//create the database and open its connection
	db := sqladapter.CreateSQLDB("integration", dependencies.ResolveSQLDriver())

	err = db.OpenConnection()
	suite.Require().NoError(err)

	err = db.Ping()
	suite.Require().NoError(err)

	suite.TransactionFactory = &sqladapter.SQLTransactionFactory{
		DB: db,
	}
}

func (suite *ScopeCRUDTestSuite) TearDownSuite() {
	suite.TransactionFactory.DB.CloseConnection()
}

func (suite *ScopeCRUDTestSuite) SetupTest() {
	//start a new transaction for every test
	tx, err := suite.TransactionFactory.CreateTransaction()
	suite.Require().NoError(err)

	suite.Tx = tx.(*sqladapter.SQLTransaction)
}

func (suite *ScopeCRUDTestSuite) TearDownTest() {
	//rollback the transaction after each test
	err := suite.Tx.RollbackTransaction()
	suite.Require().NoError(err)
}

func (suite *ScopeCRUDTestSuite) TestSaveScope_WithInvalidScope_ReturnsError() {
	//arrange
	scope := models.CreateNewScope("")

	//act
	err := suite.Tx.SaveScope(scope)

	//assert
	common.AssertError(&suite.Suite, err, "error", "scope model")
}

func (suite *ScopeCRUDTestSuite) TestGetScopeByName_WhereScopeNotFound_ReturnsNilScope() {
	//act
	scope, err := suite.Tx.GetScopeByName("not a real name")

	//assert
	suite.NoError(err)
	suite.Nil(scope)
}

func (suite *ScopeCRUDTestSuite) TestGetScopeByName_GetsTheScopeWithName() {
	//arrange
	scope := models.CreateNewScope("name")
	SaveScope(&suite.Suite, suite.Tx, scope)

	//act
	resultScope, err := suite.Tx.GetScopeByName(scope.Name)

	//assert
	suite.NoError(err)
	suite.EqualValues(scope, resultScope)
}

func TestScopeCRUDTestSuite(t *testing.T) {
	suite.Run(t, &ScopeCRUDTestSuite{})
}

func SaveScope(suite *suite.Suite, tx *sqladapter.SQLTransaction, scope *models.Scope) {
	err := tx.SaveScope(scope)
	suite.Require().NoError(err)
}
