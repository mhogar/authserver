package sqladapter_test

import (
	"authserver/config"
	sqladapter "authserver/database/sql_adapter"
	"authserver/dependencies"
	commonhelpers "authserver/helpers/common"
	"authserver/models"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type ClientCRUDTestSuite struct {
	suite.Suite
	TransactionFactory *sqladapter.SQLTransactionFactory
	Tx                 *sqladapter.SQLTransaction
}

func (suite *ClientCRUDTestSuite) SetupSuite() {
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

func (suite *ClientCRUDTestSuite) TearDownSuite() {
	suite.TransactionFactory.DB.CloseConnection()
}

func (suite *ClientCRUDTestSuite) SetupTest() {
	//start a new transaction for every test
	tx, err := suite.TransactionFactory.CreateTransaction()
	suite.Require().NoError(err)

	suite.Tx = tx.(*sqladapter.SQLTransaction)
}

func (suite *ClientCRUDTestSuite) TearDownTest() {
	//rollback the transaction after each test
	err := suite.Tx.RollbackTransaction()
	suite.Require().NoError(err)
}

func (suite *ClientCRUDTestSuite) TestSaveClient_WithInvalidClient_ReturnsError() {
	//arrange
	client := &models.Client{
		ID: uuid.Nil,
	}

	//act
	err := suite.Tx.SaveClient(client)

	//assert
	commonhelpers.AssertError(&suite.Suite, err, "error", "model")
}

func (suite *ClientCRUDTestSuite) TestGetClientById_WhereClientNotFound_ReturnsNilClient() {
	//act
	client, err := suite.Tx.GetClientByID(uuid.New())

	//assert
	suite.NoError(err)
	suite.Nil(client)
}

func (suite *ClientCRUDTestSuite) TestGetClientById_GetsTheClientWithId() {
	//arrange
	client := models.CreateNewClient()
	err := suite.Tx.SaveClient(client)
	suite.Require().NoError(err)

	//act
	resultClient, err := suite.Tx.GetClientByID(client.ID)

	//assert
	suite.NoError(err)
	suite.EqualValues(client, resultClient)
}

func TestClientCRUDTestSuite(t *testing.T) {
	suite.Run(t, &ClientCRUDTestSuite{})
}
