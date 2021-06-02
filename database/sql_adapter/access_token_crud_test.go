package sqladapter_test

import (
	"authserver/common"
	"authserver/config"
	sqladapter "authserver/database/sql_adapter"
	"authserver/dependencies"
	"authserver/models"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type AccessTokenCRUDTestSuite struct {
	suite.Suite
	TransactionFactory *sqladapter.SQLTransactionFactory
	Tx                 *sqladapter.SQLTransaction
}

func (suite *AccessTokenCRUDTestSuite) SetupSuite() {
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

func (suite *AccessTokenCRUDTestSuite) TearDownSuite() {
	suite.TransactionFactory.DB.CloseConnection()
}

func (suite *AccessTokenCRUDTestSuite) SetupTest() {
	//start a new transaction for every test
	tx, err := suite.TransactionFactory.CreateTransaction()
	suite.Require().NoError(err)

	suite.Tx = tx.(*sqladapter.SQLTransaction)
}

func (suite *AccessTokenCRUDTestSuite) TearDownTest() {
	//rollback the transaction after each test
	err := suite.Tx.RollbackTransaction()
	suite.Require().NoError(err)
}

func (suite *AccessTokenCRUDTestSuite) TestSaveAccessToken_WithInvalidAccessToken_ReturnsError() {
	//act
	err := suite.Tx.SaveAccessToken(models.CreateNewAccessToken(nil, nil, nil))

	//assert
	common.AssertError(&suite.Suite, err, "error", "access token model")
}

func (suite *AccessTokenCRUDTestSuite) TestGetAccessTokenById_WhereAccessTokenNotFound_ReturnsNilAccessToken() {
	//act
	token, err := suite.Tx.GetAccessTokenByID(uuid.New())

	//assert
	suite.NoError(err)
	suite.Nil(token)
}

func (suite *AccessTokenCRUDTestSuite) TestGetAccessTokenById_GetsTheAccessTokenWithId() {
	//arrange
	token := models.CreateNewAccessToken(
		models.CreateNewUser("username", []byte("password")),
		models.CreateNewClient(),
		models.CreateNewScope("name"),
	)
	SaveAccessToken(&suite.Suite, suite.Tx, token)

	//act
	resultAccessToken, err := suite.Tx.GetAccessTokenByID(token.ID)

	//assert
	suite.NoError(err)
	suite.EqualValues(token, resultAccessToken)
}

func (suite *AccessTokenCRUDTestSuite) TestDeleteAccessToken_WithNoAccessTokenToDelete_ReturnsNilError() {
	//act
	err := suite.Tx.DeleteAccessToken(models.CreateNewAccessToken(nil, nil, nil))

	//assert
	suite.NoError(err)
}

func (suite *AccessTokenCRUDTestSuite) TestDeleteAccessToken_DeletesAccessTokenWithId() {
	//arrange
	token := models.CreateNewAccessToken(
		models.CreateNewUser("username", []byte("password")),
		models.CreateNewClient(),
		models.CreateNewScope("name"),
	)
	SaveAccessToken(&suite.Suite, suite.Tx, token)

	//act
	err := suite.Tx.DeleteAccessToken(token)

	//assert
	suite.Require().NoError(err)

	resultAccessToken, err := suite.Tx.GetAccessTokenByID(token.ID)
	suite.NoError(err)
	suite.Nil(resultAccessToken)
}

//TODO: test cascade delete with users, scopes, and clients

func TestAccessTokenCRUDTestSuite(t *testing.T) {
	suite.Run(t, &AccessTokenCRUDTestSuite{})
}

func SaveAccessToken(suite *suite.Suite, tx *sqladapter.SQLTransaction, token *models.AccessToken) {
	SaveUser(suite, tx, token.User)
	SaveClient(suite, tx, token.Client)
	SaveScope(suite, tx, token.Scope)

	err := tx.SaveAccessToken(token)
	suite.Require().NoError(err)
}
