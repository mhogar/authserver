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

type UserCRUDTestSuite struct {
	suite.Suite
	TransactionFactory *sqladapter.SQLTransactionFactory
	Tx                 *sqladapter.SQLTransaction
}

func (suite *UserCRUDTestSuite) SetupSuite() {
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

func (suite *UserCRUDTestSuite) TearDownSuite() {
	suite.TransactionFactory.DB.CloseConnection()
}

func (suite *UserCRUDTestSuite) SetupTest() {
	//start a new transaction for every test
	tx, err := suite.TransactionFactory.CreateTransaction()
	suite.Require().NoError(err)

	suite.Tx = tx.(*sqladapter.SQLTransaction)
}

func (suite *UserCRUDTestSuite) TearDownTest() {
	//rollback the transaction after each test
	err := suite.Tx.RollbackTransaction()
	suite.Require().NoError(err)
}

func (suite *UserCRUDTestSuite) TestSaveUser_WithInvalidUser_ReturnsError() {
	//act
	err := suite.Tx.SaveUser(models.CreateNewUser("", nil))

	//assert
	common.AssertError(&suite.Suite, err, "error", "user model")
}

func (suite *UserCRUDTestSuite) TestGetUserById_WhereUserNotFound_ReturnsNilUser() {
	//act
	user, err := suite.Tx.GetUserByID(uuid.New())

	//assert
	suite.NoError(err)
	suite.Nil(user)
}

func (suite *UserCRUDTestSuite) TestGetUserById_GetsTheUserWithId() {
	//arrange
	user := models.CreateNewUser("username", []byte("password"))
	SaveUser(&suite.Suite, suite.Tx, user)

	//act
	resultUser, err := suite.Tx.GetUserByID(user.ID)

	//assert
	suite.NoError(err)
	suite.EqualValues(user, resultUser)
}

func (suite *UserCRUDTestSuite) TestGetUserByUsername_WhereUserNotFound_ReturnsNilUser() {
	//act
	user, err := suite.Tx.GetUserByUsername("DNE")

	//assert
	suite.NoError(err)
	suite.Nil(user)
}

func (suite *UserCRUDTestSuite) TestGetUserByUsernameGetsTheUserWithUsername() {
	//arrange
	user := models.CreateNewUser("username", []byte("password"))
	SaveUser(&suite.Suite, suite.Tx, user)

	//act
	resultUser, err := suite.Tx.GetUserByUsername(user.Username)

	//assert
	suite.NoError(err)
	suite.EqualValues(user, resultUser)
}

func (suite *UserCRUDTestSuite) TestUpdateUser_WithInvalidUser_ReturnsError() {
	//act
	err := suite.Tx.UpdateUser(models.CreateNewUser("", nil))

	//assert
	common.AssertError(&suite.Suite, err, "error", "user model")
}

func (suite *UserCRUDTestSuite) TestUpdateUser_WithNoUserToUpdate_ReturnsNilError() {
	//act
	err := suite.Tx.UpdateUser(models.CreateNewUser("username", []byte("password")))

	//assert
	suite.NoError(err)
}

func (suite *UserCRUDTestSuite) TestUpdateUser_UpdatesUserWithId() {
	//arrange
	user := models.CreateNewUser("username", []byte("password"))
	SaveUser(&suite.Suite, suite.Tx, user)

	//act
	user.Username = "username2"
	err := suite.Tx.UpdateUser(user)

	//assert
	suite.Require().NoError(err)

	resultUser, err := suite.Tx.GetUserByID(user.ID)
	suite.NoError(err)
	suite.EqualValues(user, resultUser)
}

func (suite *UserCRUDTestSuite) TestDeleteUser_WithNoUserToDelete_ReturnsNilError() {
	//act
	err := suite.Tx.DeleteUser(models.CreateNewUser("", nil))

	//assert
	suite.NoError(err)
}

func (suite *UserCRUDTestSuite) TestDeleteUser_DeletesUserWithId() {
	//arrange
	user := models.CreateNewUser("username", []byte("password"))
	SaveUser(&suite.Suite, suite.Tx, user)

	//act
	err := suite.Tx.DeleteUser(user)

	//assert
	suite.Require().NoError(err)

	resultUser, err := suite.Tx.GetUserByID(user.ID)
	suite.NoError(err)
	suite.Nil(resultUser)
}

func TestUserCRUDTestSuite(t *testing.T) {
	suite.Run(t, &UserCRUDTestSuite{})
}

func SaveUser(suite *suite.Suite, tx *sqladapter.SQLTransaction, user *models.User) {
	err := tx.SaveUser(user)
	suite.Require().NoError(err)
}
