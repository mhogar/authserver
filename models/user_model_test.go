package models_test

import (
	"testing"

	"authserver/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type UserTestSuite struct {
	suite.Suite
	User *models.User
}

func (suite *UserTestSuite) SetupTest() {
	suite.User = models.CreateNewUser("username", []byte("password"))
}

func (suite *UserTestSuite) TestCreateNewUser_CreatesUserWithSuppliedFields() {
	//arrange
	username := "this is a test username"
	hash := []byte("this is a password")

	//act
	user := models.CreateNewUser(username, hash)

	//assert
	suite.Require().NotNil(user)
	suite.NotEqual(user.ID, uuid.Nil)
	suite.EqualValues(username, user.Username)
	suite.EqualValues(hash, user.PasswordHash)
}

func (suite *UserTestSuite) TestValidate_WithValidUser_ReturnsModelValid() {
	//act
	err := suite.User.Validate()

	//assert
	suite.EqualValues(models.ValidateUserValid, err.Status)
}

func (suite *UserTestSuite) TestValidate_WithNilUserID_ReturnsUserInvalidID() {
	//arrange
	suite.User.ID = uuid.Nil

	//act
	err := suite.User.Validate()

	//assert
	suite.EqualValues(models.ValidateUserInvalidID, err.Status)
}

func (suite *UserTestSuite) TestValidate_WithEmptyUsername_ReturnsUserInvalidUsername() {
	//arrange
	suite.User.Username = ""

	//act
	err := suite.User.Validate()

	//assert
	suite.EqualValues(models.ValidateUserInvalidUsername, err.Status)
}

func (suite *UserTestSuite) TestValidate_WithEmptyPasswordHash_ReturnsUserInvalidPasswordHash() {
	//arrange
	suite.User.PasswordHash = make([]byte, 0)

	//act
	err := suite.User.Validate()

	//assert
	suite.EqualValues(models.ValidateUserInvalidPasswordHash, err.Status)
}

func TestUserTestSuite(t *testing.T) {
	suite.Run(t, &UserTestSuite{})
}
