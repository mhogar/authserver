package sqladapter_test

import (
	"authserver/config"
	sqladapter "authserver/database/sql_adapter"
	"authserver/dependencies"
	"authserver/models"

	"github.com/stretchr/testify/suite"
)

type CRUDTestSuite struct {
	suite.Suite
	TransactionFactory *sqladapter.SQLTransactionFactory
	Tx                 *sqladapter.SQLTransaction
}

func (suite *CRUDTestSuite) SetupSuite() {
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

func (suite *CRUDTestSuite) TearDownSuite() {
	suite.TransactionFactory.DB.CloseConnection()
}

func (suite *CRUDTestSuite) SetupTest() {
	//start a new transaction for every test
	tx, err := suite.TransactionFactory.CreateTransaction()
	suite.Require().NoError(err)

	suite.Tx = tx.(*sqladapter.SQLTransaction)
}

func (suite *CRUDTestSuite) TearDownTest() {
	//rollback the transaction after each test
	suite.Tx.RollbackTransaction()
}

func (suite *CRUDTestSuite) SaveUser(tx *sqladapter.SQLTransaction, user *models.User) {
	err := tx.SaveUser(user)
	suite.Require().NoError(err)
}

func (suite *CRUDTestSuite) SaveScope(tx *sqladapter.SQLTransaction, scope *models.Scope) {
	err := tx.SaveScope(scope)
	suite.Require().NoError(err)
}

func (suite *CRUDTestSuite) SaveClient(tx *sqladapter.SQLTransaction, client *models.Client) {
	err := tx.SaveClient(client)
	suite.Require().NoError(err)
}

func (suite *CRUDTestSuite) SaveAccessToken(tx *sqladapter.SQLTransaction, token *models.AccessToken) {
	err := tx.SaveAccessToken(token)
	suite.Require().NoError(err)
}

func (suite *CRUDTestSuite) SaveAccessTokenAndFields(tx *sqladapter.SQLTransaction, token *models.AccessToken) {
	suite.SaveUser(tx, token.User)
	suite.SaveClient(tx, token.Client)
	suite.SaveScope(tx, token.Scope)
	suite.SaveAccessToken(tx, token)
}
