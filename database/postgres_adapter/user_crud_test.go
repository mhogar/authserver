package postgresadapter_test

import (
	"authserver/config"
	"authserver/database"
	postgresadapter "authserver/database/postgres_adapter"
	sqladapter "authserver/database/sql_adapter"
	"authserver/helpers"
	"authserver/models"
	"testing"

	"github.com/stretchr/testify/suite"
)

type UserCRUDTestSuite struct {
	suite.Suite
	TransactionFactory *sqladapter.SQLTransactionFactory
	Transaction        database.Transaction
}

func (suite *UserCRUDTestSuite) SetupSuite() {
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

func (suite *UserCRUDTestSuite) TearDownSuite() {
	suite.TransactionFactory.DB.CloseConnection()
}

func (suite *UserCRUDTestSuite) SetupTest() {
	//start a new transaction for every test
	tx, err := suite.TransactionFactory.CreateTransaction()
	suite.Require().NoError(err)

	suite.Transaction = tx
}

func (suite *UserCRUDTestSuite) TearDownTest() {
	//rollback the transaction after each test
	err := suite.Transaction.RollbackTransaction()
	suite.Require().NoError(err)
}

func (suite *UserCRUDTestSuite) TestSaveUser_WithInvalidUser_ReturnsError() {
	//act
	err := suite.Transaction.SaveUser(models.CreateNewUser("", nil))

	//assert
	helpers.AssertError(&suite.Suite, err, "error", "model")
}

func (suite *UserCRUDTestSuite) TestGetUserById_GetsTheUser() {
	//arrange
	user := models.CreateNewUser("user.name", []byte("password"))
	err := suite.Transaction.SaveUser(user)
	suite.Require().NoError(err)

	//act
	resultUser, err := suite.Transaction.GetUserByID(user.ID)

	//assert
	suite.NoError(err)
	suite.EqualValues(user, resultUser)
}

func TestUserCRUDTestSuite(t *testing.T) {
	suite.Run(t, &UserCRUDTestSuite{})
}
