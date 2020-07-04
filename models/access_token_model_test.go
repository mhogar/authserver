package models_test

import (
	"testing"

	"authserver/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type AccessTokenTestSuite struct {
	suite.Suite
	Token *models.AccessToken
}

func (suite *AccessTokenTestSuite) SetupTest() {
	suite.Token = models.CreateNewAccessToken(
		models.CreateNewUser("username", []byte("password")),
		models.CreateNewClient(),
		models.CreateNewScope("name"),
	)
}

func (suite *AccessTokenTestSuite) TestCreateNewAccessToken_CreatesAccessTokenWithSuppliedFields() {
	//arrange
	user := models.CreateNewUser("", nil)
	client := models.CreateNewClient()
	scope := models.CreateNewScope("")

	//act
	token := models.CreateNewAccessToken(user, client, scope)

	//assert
	suite.Require().NotNil(token)
	suite.NotEqual(token.ID, uuid.Nil)
	suite.Equal(token.User, user)
	suite.Equal(token.Client, client)
	suite.Equal(token.Scope, scope)
}

func (suite *AccessTokenTestSuite) TestValidate_WithValidAccessToken_ReturnsValid() {
	//act
	err := suite.Token.Validate()

	//assert
	suite.Equal(models.ValidateAccessTokenValid, err.Status)
}

func (suite *AccessTokenTestSuite) TestValidate_WithNilID_ReturnsAccessTokenInvalidID() {
	//arrange
	suite.Token.ID = uuid.Nil

	//act
	err := suite.Token.Validate()

	//assert
	suite.Equal(models.ValidateAccessTokenNilID, err.Status)
}

func (suite *AccessTokenTestSuite) TestValidate_WithNilUser_ReturnsAccessTokenNilUser() {
	//arrange
	suite.Token.User = nil

	//act
	err := suite.Token.Validate()

	//assert
	suite.Equal(models.ValidateAccessTokenNilUser, err.Status)
}

func (suite *AccessTokenTestSuite) TestValidate_WithInvalidUser_ReturnsAccessTokenInvalidUser() {
	//arrange
	suite.Token.User = models.CreateNewUser("", nil)

	//act
	err := suite.Token.Validate()

	//assert
	suite.Equal(models.ValidateAccessTokenInvalidUser, err.Status)
}

func (suite *AccessTokenTestSuite) TestValidate_WithNilClient_ReturnsAccessTokenNilClient() {
	//arrange
	suite.Token.Client = nil

	//act
	err := suite.Token.Validate()

	//assert
	suite.Equal(models.ValidateAccessTokenNilClient, err.Status)
}

func (suite *AccessTokenTestSuite) TestValidate_WithInvalidClient_ReturnsAccessTokenInvalidClient() {
	//arrange
	suite.Token.Client = &models.Client{ID: uuid.Nil}

	//act
	err := suite.Token.Validate()

	//assert
	suite.Equal(models.ValidateAccessTokenInvalidClient, err.Status)
}

func (suite *AccessTokenTestSuite) TestValidate_WithNilScope_ReturnsAccessTokenNilScope() {
	//arrange
	suite.Token.Scope = nil

	//act
	err := suite.Token.Validate()

	//assert
	suite.Equal(models.ValidateAccessTokenNilScope, err.Status)
}

func (suite *AccessTokenTestSuite) TestValidate_WithInvalidScope_ReturnsAccessTokenInvalidScope() {
	//arrange
	suite.Token.Scope = models.CreateNewScope("")

	//act
	err := suite.Token.Validate()

	//assert
	suite.Equal(models.ValidateAccessTokenInvalidScope, err.Status)
}

func TestAccessTokenTestSuite(t *testing.T) {
	suite.Run(t, &AccessTokenTestSuite{})
}
