package router_test

import (
	"authserver/common"
	requesterror "authserver/common/request_error"
	controllermocks "authserver/controllers/mocks"
	"authserver/models"
	"authserver/router"
	"authserver/router/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type TokenHandlerSuite struct {
	suite.Suite
	ControllersMock   controllermocks.Controllers
	AuthenticatorMock mocks.Authenticator
	Router            *httprouter.Router
}

func (suite *TokenHandlerSuite) SetupTest() {
	suite.ControllersMock = controllermocks.Controllers{}
	suite.AuthenticatorMock = mocks.Authenticator{}
	suite.Router = router.CreateRouter(&suite.ControllersMock, &suite.AuthenticatorMock)
}

func (suite *TokenHandlerSuite) TestPostToken_WithInvalidJSONBody_ReturnsInvalidRequest() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	req := common.CreateRequest(&suite.Suite, http.MethodPost, server.URL+"/token", "", "invalid")

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	common.AssertOAuthErrorResponse(&suite.Suite, res, http.StatusBadRequest, "invalid_request", "invalid json body")
}

func (suite *TokenHandlerSuite) TestPostToken_WithMissingGrantType_ReturnsInvalidRequest() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	body := router.PostTokenBody{}
	req := common.CreateRequest(&suite.Suite, http.MethodPost, server.URL+"/token", "", body)

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	common.AssertOAuthErrorResponse(&suite.Suite, res, http.StatusBadRequest, "invalid_request", "missing grant_type parameter")
}

func (suite *TokenHandlerSuite) TestPostToken_WithUnsupportedGrantType_ReturnsUnsupportedGrantType() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	body := router.PostTokenBody{
		GrantType: "unsupported",
	}
	req := common.CreateRequest(&suite.Suite, http.MethodPost, server.URL+"/token", "", body)

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	common.AssertOAuthErrorResponse(&suite.Suite, res, http.StatusBadRequest, "unsupported_grant_type", "")
}

func (suite *TokenHandlerSuite) TestPostToken_PasswordGrant_WithMissingParameters_ReturnsInvalidRequest() {
	var grantBody router.PostTokenPasswordGrantBody
	var expectedErrorDescription string

	testCase := func() {
		//arrange
		server := httptest.NewServer(suite.Router)
		defer server.Close()

		body := router.PostTokenBody{
			GrantType:                  "password",
			PostTokenPasswordGrantBody: grantBody,
		}
		req := common.CreateRequest(&suite.Suite, http.MethodPost, server.URL+"/token", "", body)

		//act
		res, err := http.DefaultClient.Do(req)
		suite.Require().NoError(err)

		//assert
		common.AssertOAuthErrorResponse(&suite.Suite, res, http.StatusBadRequest, "invalid_request", expectedErrorDescription)
	}

	grantBody = router.PostTokenPasswordGrantBody{
		Password: "password",
		ClientID: "client id",
		Scope:    "scope",
	}
	expectedErrorDescription = "missing username parameter"
	suite.Run("MissingUsername", testCase)

	grantBody = router.PostTokenPasswordGrantBody{
		Username: "username",
		ClientID: "client id",
		Scope:    "scope",
	}
	expectedErrorDescription = "missing password parameter"
	suite.Run("MissingPassword", testCase)

	grantBody = router.PostTokenPasswordGrantBody{
		Username: "username",
		Password: "password",
		Scope:    "scope",
	}
	expectedErrorDescription = "missing client_id parameter"
	suite.Run("MissingClientID", testCase)

	grantBody = router.PostTokenPasswordGrantBody{
		Username: "username",
		Password: "password",
		ClientID: "client id",
	}
	expectedErrorDescription = "missing scope parameter"
	suite.Run("MissingScope", testCase)
}

func (suite *TokenHandlerSuite) TestPostToken_PasswordGrant_WithClientErrorCreatingTokenFromPassword_ReturnsInvalidClient() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	body := router.PostTokenBody{
		GrantType: "password",
		PostTokenPasswordGrantBody: router.PostTokenPasswordGrantBody{
			Username: "username",
			Password: "password",
			ClientID: uuid.New().String(),
			Scope:    "scope",
		},
	}
	req := common.CreateRequest(&suite.Suite, http.MethodPost, server.URL+"/token", "", body)

	errorName := "error_name"
	message := "create token error"
	suite.ControllersMock.On("CreateTokenFromPassword", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil, requesterror.OAuthClientError(errorName, message))

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	common.AssertOAuthErrorResponse(&suite.Suite, res, http.StatusBadRequest, errorName, message)
}

func (suite *TokenHandlerSuite) TestPostToken_PasswordGrant_WithInternalErrorCreatingTokenFromPassword_ReturnsInternalServerError() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	body := router.PostTokenBody{
		GrantType: "password",
		PostTokenPasswordGrantBody: router.PostTokenPasswordGrantBody{
			Username: "username",
			Password: "password",
			ClientID: uuid.New().String(),
			Scope:    "scope",
		},
	}
	req := common.CreateRequest(&suite.Suite, http.MethodPost, server.URL+"/token", "", body)

	suite.ControllersMock.On("CreateTokenFromPassword", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil, requesterror.OAuthInternalError())

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	common.AssertInternalServerErrorResponse(&suite.Suite, res)
}

func (suite *TokenHandlerSuite) TestPostToken_PasswordGrant_WithValidRequest_ReturnsAccessToken() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	clientID := uuid.New()
	token := models.CreateNewAccessToken(nil, nil, nil)

	body := router.PostTokenBody{
		GrantType: "password",
		PostTokenPasswordGrantBody: router.PostTokenPasswordGrantBody{
			Username: "username",
			Password: "password",
			ClientID: clientID.String(),
			Scope:    "scope",
		},
	}
	req := common.CreateRequest(&suite.Suite, http.MethodPost, server.URL+"/token", "", body)

	suite.ControllersMock.On("CreateTokenFromPassword", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(token, requesterror.OAuthNoError())

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	suite.ControllersMock.AssertCalled(suite.T(), "CreateTokenFromPassword", body.Username, body.Password, clientID, body.Scope)
	common.AssertAccessTokenResponse(&suite.Suite, res, token.ID.String())
}

func (suite *TokenHandlerSuite) TestPostToken_PasswordGrant_WithPanicTriggered_ReturnsInternalServerError() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	body := router.PostTokenBody{
		GrantType: "password",
		PostTokenPasswordGrantBody: router.PostTokenPasswordGrantBody{
			Username: "username",
			Password: "password",
			ClientID: uuid.New().String(),
			Scope:    "scope",
		},
	}
	req := common.CreateRequest(&suite.Suite, http.MethodPost, server.URL+"/token", "", body)

	suite.ControllersMock.On("CreateTokenFromPassword", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Run(func(_ mock.Arguments) {
		panic("test panic handler")
	})

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	common.AssertInternalServerErrorResponse(&suite.Suite, res)
}

func (suite *TokenHandlerSuite) TestDeleteToken_WithClientErrorAuthenticatingUser_ReturnsUnauthorized() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	req := common.CreateRequest(&suite.Suite, http.MethodDelete, server.URL+"/token", "", nil)

	message := "authenticate error"
	suite.AuthenticatorMock.On("Authenticate", mock.Anything).Return(nil, requesterror.ClientError(message))

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	common.AssertErrorResponse(&suite.Suite, res, http.StatusUnauthorized, message)
}

func (suite *TokenHandlerSuite) TestDeleteToken_WithInternalErrorAuthenticatingUser_ReturnsInternalServerError() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	req := common.CreateRequest(&suite.Suite, http.MethodDelete, server.URL+"/token", "", nil)

	suite.AuthenticatorMock.On("Authenticate", mock.Anything).Return(nil, requesterror.InternalError())

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	common.AssertInternalServerErrorResponse(&suite.Suite, res)
}

func (suite *TokenHandlerSuite) TestDeleteToken_WithClientErrorDeletingToken_ReturnsBadRequest() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	token := &models.AccessToken{}
	req := common.CreateRequest(&suite.Suite, http.MethodDelete, server.URL+"/token", "", nil)

	message := "delete token error"
	suite.AuthenticatorMock.On("Authenticate", mock.Anything).Return(token, requesterror.NoError())
	suite.ControllersMock.On("DeleteToken", mock.Anything).Return(requesterror.ClientError(message))

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	common.AssertErrorResponse(&suite.Suite, res, http.StatusBadRequest, message)
}

func (suite *TokenHandlerSuite) TestDeleteToken_WithInternalErrorDeletingToken_ReturnsInternalServerError() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	token := &models.AccessToken{}
	req := common.CreateRequest(&suite.Suite, http.MethodDelete, server.URL+"/token", "", nil)

	suite.AuthenticatorMock.On("Authenticate", mock.Anything).Return(token, requesterror.NoError())
	suite.ControllersMock.On("DeleteToken", mock.Anything).Return(requesterror.InternalError())

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	common.AssertInternalServerErrorResponse(&suite.Suite, res)
}

func (suite *TokenHandlerSuite) TestDeleteToken_WithValidRequest_ReturnsSuccess() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	token := &models.AccessToken{}
	req := common.CreateRequest(&suite.Suite, http.MethodDelete, server.URL+"/token", "", nil)

	suite.AuthenticatorMock.On("Authenticate", mock.Anything).Return(token, requesterror.NoError())
	suite.ControllersMock.On("DeleteToken", mock.Anything).Return(requesterror.NoError())

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	suite.AuthenticatorMock.AssertCalled(suite.T(), "Authenticate", mock.Anything)
	suite.ControllersMock.AssertCalled(suite.T(), "DeleteToken", token)
	common.AssertSuccessResponse(&suite.Suite, res)
}

func (suite *TokenHandlerSuite) TestDeleteToken_WithPanicTriggered_ReturnsInternalServerError() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	token := &models.AccessToken{}
	req := common.CreateRequest(&suite.Suite, http.MethodDelete, server.URL+"/token", "", nil)

	suite.AuthenticatorMock.On("Authenticate", mock.Anything).Return(token, requesterror.NoError())
	suite.ControllersMock.On("DeleteToken", mock.Anything).Run(func(_ mock.Arguments) {
		panic("test panic handler")
	})

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	common.AssertInternalServerErrorResponse(&suite.Suite, res)
}

func TestTokenHandlerTestSuite(t *testing.T) {
	suite.Run(t, &TokenHandlerSuite{})
}
