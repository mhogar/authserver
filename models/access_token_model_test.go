package models_test

import (
	"testing"

	"authserver/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type AccessTokenSuite struct {
	suite.Suite
	Token *models.AccessToken
}

func (suite *AccessTokenSuite) SetupTest() {
	suite.Token = models.CreateNewAccessToken(uuid.New(), uuid.New(), uuid.New())
}

func (suite *AccessTokenSuite) TestCreateNewAccessToken_CreatesAccessTokenWithSuppliedFields() {
	//arrange
	userID := uuid.New()
	clientID := uuid.New()
	scopeID := uuid.New()

	//act
	token := models.CreateNewAccessToken(userID, clientID, scopeID)

	//assert
	suite.Require().NotNil(token)
	suite.NotEqual(token.ID, uuid.Nil)
	suite.NotEqual(token.UserID, uuid.Nil)
	suite.NotEqual(token.ClientID, uuid.Nil)
	suite.NotEqual(token.ScopeID, uuid.Nil)
}

func (suite *AccessTokenSuite) TestValidate_WithValidAccessToken_ReturnsValid() {
	//act
	err := suite.Token.Validate()

	//assert
	suite.Equal(models.ValidateAccessTokenValid, err.Status)
}

func (suite *AccessTokenSuite) TestValidate_WithNilID_ReturnsAccessTokenInvalidID() {
	//arrange
	suite.Token.ID = uuid.Nil

	//act
	err := suite.Token.Validate()

	//assert
	suite.Equal(models.ValidateAccessTokenInvalidID, err.Status)
}

func (suite *AccessTokenSuite) TestValidate_WithNilUserID_ReturnsAccessTokenInvalidUserID() {
	//arrange
	suite.Token.UserID = uuid.Nil

	//act
	err := suite.Token.Validate()

	//assert
	suite.Equal(models.ValidateAccessTokenInvalidUserID, err.Status)
}

func (suite *AccessTokenSuite) TestValidate_WithNilClientID_ReturnsAccessTokenInvalidClientID() {
	//arrange
	suite.Token.ClientID = uuid.Nil

	//act
	err := suite.Token.Validate()

	//assert
	suite.Equal(models.ValidateAccessTokenInvalidClientID, err.Status)
}

func (suite *AccessTokenSuite) TestValidate_WithNilScopeID_ReturnsAccessTokenInvalidScopeID() {
	//arrange
	suite.Token.ScopeID = uuid.Nil

	//act
	err := suite.Token.Validate()

	//assert
	suite.Equal(models.ValidateAccessTokenInvalidScopeID, err.Status)
}

func TestAccessTokenSuite(t *testing.T) {
	suite.Run(t, &AccessTokenSuite{})
}
