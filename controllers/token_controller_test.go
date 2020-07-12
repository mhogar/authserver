package controllers_test

import (
	"authserver/controllers"
	"authserver/models"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"

	databasemocks "authserver/database/mocks"
	helpermocks "authserver/helpers/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type TokenControllerTestSuite struct {
	suite.Suite
	CRUDMock           databasemocks.CRUDOperations
	PasswordHasherMock helpermocks.PasswordHasher
	TokenController    controllers.TokenController
}

func (suite *TokenControllerTestSuite) SetupTest() {
	suite.CRUDMock = databasemocks.CRUDOperations{}
	suite.PasswordHasherMock = helpermocks.PasswordHasher{}
	suite.TokenController = controllers.TokenController{
		CRUD:           &suite.CRUDMock,
		PasswordHasher: &suite.PasswordHasherMock,
	}
}

func (suite *TokenControllerTestSuite) TestPostToken_WithInvalidJSONBody_ReturnsInvalidRequest() {
	//arrange
	w := httptest.NewRecorder()
	req := CreateRequest(&suite.Suite, "", "invalid")

	//act
	suite.TokenController.PostToken(w, req, nil)

	//assert
	AssertOAuthErrorResponse(&suite.Suite, w.Result(), http.StatusBadRequest, "invalid_request", "invalid json body")
}

func (suite *TokenControllerTestSuite) TestPostToken_WithMissingGrantType_ReturnsInvalidRequest() {
	//arrange
	w := httptest.NewRecorder()

	body := controllers.PostTokenBody{}
	req := CreateRequest(&suite.Suite, "", body)

	//act
	suite.TokenController.PostToken(w, req, nil)

	//assert
	AssertOAuthErrorResponse(&suite.Suite, w.Result(), http.StatusBadRequest, "invalid_request", "missing grant_type parameter")
}

func (suite *TokenControllerTestSuite) TestPostToken_WithUnsupportedGrantType_ReturnsUnsupportedGrantType() {
	//arrange
	w := httptest.NewRecorder()

	body := controllers.PostTokenBody{
		GrantType: "unsupported",
	}
	req := CreateRequest(&suite.Suite, "", body)

	//act
	suite.TokenController.PostToken(w, req, nil)

	//assert
	AssertOAuthErrorResponse(&suite.Suite, w.Result(), http.StatusBadRequest, "unsupported_grant_type", "")
}

func (suite *TokenControllerTestSuite) TestPostToken_PasswordGrant_WithMissingParameters_ReturnsInvalidRequest() {
	var grantBody controllers.PostTokenPasswordGrantBody
	var expectedErrorDescription string

	testCase := func() {
		//arrange
		w := httptest.NewRecorder()

		body := controllers.PostTokenBody{
			GrantType:                  "password",
			PostTokenPasswordGrantBody: grantBody,
		}
		req := CreateRequest(&suite.Suite, "", body)

		//act
		suite.TokenController.PostToken(w, req, nil)

		//assert
		AssertOAuthErrorResponse(&suite.Suite, w.Result(), http.StatusBadRequest, "invalid_request", expectedErrorDescription)
	}

	grantBody = controllers.PostTokenPasswordGrantBody{
		Password: "password",
		ClientID: "client id",
		Scope:    "scope",
	}
	expectedErrorDescription = "missing username parameter"
	suite.Run("MissingUsername", testCase)

	grantBody = controllers.PostTokenPasswordGrantBody{
		Username: "username",
		ClientID: "client id",
		Scope:    "scope",
	}
	expectedErrorDescription = "missing password parameter"
	suite.Run("MissingPassword", testCase)

	grantBody = controllers.PostTokenPasswordGrantBody{
		Username: "username",
		Password: "password",
		Scope:    "scope",
	}
	expectedErrorDescription = "missing client_id parameter"
	suite.Run("MissingClientID", testCase)

	grantBody = controllers.PostTokenPasswordGrantBody{
		Username: "username",
		Password: "password",
		ClientID: "client id",
	}
	expectedErrorDescription = "missing scope parameter"
	suite.Run("MissingScope", testCase)
}

func (suite *TokenControllerTestSuite) TestPostToken_PasswordGrant_WithClientIDinInvalidFormat_ReturnsInvalidClient() {
	//arrange
	w := httptest.NewRecorder()

	body := controllers.PostTokenBody{
		GrantType: "password",
		PostTokenPasswordGrantBody: controllers.PostTokenPasswordGrantBody{
			Username: "username",
			Password: "password",
			ClientID: "invalid",
			Scope:    "scope",
		},
	}
	req := CreateRequest(&suite.Suite, "", body)

	//act
	suite.TokenController.PostToken(w, req, nil)

	//assert
	AssertOAuthErrorResponse(&suite.Suite, w.Result(), http.StatusUnauthorized, "invalid_client", "client_id", "invalid format")
}

func (suite *TokenControllerTestSuite) TestPostToken_PasswordGrant_WithErrorGettingClientByID_ReturnsInternalServerError() {
	//arrange
	w := httptest.NewRecorder()

	body := controllers.PostTokenBody{
		GrantType: "password",
		PostTokenPasswordGrantBody: controllers.PostTokenPasswordGrantBody{
			Username: "username",
			Password: "password",
			ClientID: uuid.New().String(),
			Scope:    "scope",
		},
	}
	req := CreateRequest(&suite.Suite, "", body)

	suite.CRUDMock.On("GetClientByID", mock.Anything).Return(nil, errors.New(""))

	//act
	suite.TokenController.PostToken(w, req, nil)

	//assert
	AssertInternalServerErrorResponse(&suite.Suite, w.Result())
}

func (suite *TokenControllerTestSuite) TestPostToken_PasswordGrant_WhereClientWithIDisNotFound_ReturnsInvalidClient() {
	//arrange
	w := httptest.NewRecorder()

	body := controllers.PostTokenBody{
		GrantType: "password",
		PostTokenPasswordGrantBody: controllers.PostTokenPasswordGrantBody{
			Username: "username",
			Password: "password",
			ClientID: uuid.New().String(),
			Scope:    "scope",
		},
	}
	req := CreateRequest(&suite.Suite, "", body)

	suite.CRUDMock.On("GetClientByID", mock.Anything).Return(nil, nil)

	//act
	suite.TokenController.PostToken(w, req, nil)

	//assert
	AssertOAuthErrorResponse(&suite.Suite, w.Result(), http.StatusUnauthorized, "invalid_client", "")
}

func (suite *TokenControllerTestSuite) TestPostToken_PasswordGrant_WithErrorGettingScopeByName_ReturnsInternalServerError() {
	//arrange
	w := httptest.NewRecorder()

	body := controllers.PostTokenBody{
		GrantType: "password",
		PostTokenPasswordGrantBody: controllers.PostTokenPasswordGrantBody{
			Username: "username",
			Password: "password",
			ClientID: uuid.New().String(),
			Scope:    "scope",
		},
	}
	req := CreateRequest(&suite.Suite, "", body)

	suite.CRUDMock.On("GetClientByID", mock.Anything).Return(&models.Client{}, nil)
	suite.CRUDMock.On("GetScopeByName", mock.Anything).Return(nil, errors.New(""))

	//act
	suite.TokenController.PostToken(w, req, nil)

	//assert
	AssertInternalServerErrorResponse(&suite.Suite, w.Result())
}

func (suite *TokenControllerTestSuite) TestPostToken_PasswordGrant_WhereNoScopeWithNameisNotFound_ReturnsInvalidScope() {
	//arrange
	w := httptest.NewRecorder()

	body := controllers.PostTokenBody{
		GrantType: "password",
		PostTokenPasswordGrantBody: controllers.PostTokenPasswordGrantBody{
			Username: "username",
			Password: "password",
			ClientID: uuid.New().String(),
			Scope:    "scope",
		},
	}
	req := CreateRequest(&suite.Suite, "", body)

	suite.CRUDMock.On("GetClientByID", mock.Anything).Return(&models.Client{}, nil)
	suite.CRUDMock.On("GetScopeByName", mock.Anything).Return(nil, nil)

	//act
	suite.TokenController.PostToken(w, req, nil)

	//assert
	AssertOAuthErrorResponse(&suite.Suite, w.Result(), http.StatusBadRequest, "invalid_scope", "")
}

func (suite *TokenControllerTestSuite) TestPostToken_PasswordGrant_WithErrorGettingUserByUsername_ReturnsInternalServerError() {
	//arrange
	w := httptest.NewRecorder()

	body := controllers.PostTokenBody{
		GrantType: "password",
		PostTokenPasswordGrantBody: controllers.PostTokenPasswordGrantBody{
			Username: "username",
			Password: "password",
			ClientID: uuid.New().String(),
			Scope:    "scope",
		},
	}
	req := CreateRequest(&suite.Suite, "", body)

	suite.CRUDMock.On("GetClientByID", mock.Anything).Return(&models.Client{}, nil)
	suite.CRUDMock.On("GetScopeByName", mock.Anything).Return(&models.Scope{}, nil)
	suite.CRUDMock.On("GetUserByUsername", mock.Anything).Return(nil, errors.New(""))

	//act
	suite.TokenController.PostToken(w, req, nil)

	//assert
	AssertInternalServerErrorResponse(&suite.Suite, w.Result())
}

func (suite *TokenControllerTestSuite) TestPostToken_PasswordGrant_WhereUserWithUsernameIsNotFound_ReturnsBadRequest() {
	//arrange
	w := httptest.NewRecorder()

	body := controllers.PostTokenBody{
		GrantType: "password",
		PostTokenPasswordGrantBody: controllers.PostTokenPasswordGrantBody{
			Username: "username",
			Password: "password",
			ClientID: uuid.New().String(),
			Scope:    "scope",
		},
	}
	req := CreateRequest(&suite.Suite, "", body)

	suite.CRUDMock.On("GetClientByID", mock.Anything).Return(&models.Client{}, nil)
	suite.CRUDMock.On("GetScopeByName", mock.Anything).Return(&models.Scope{}, nil)
	suite.CRUDMock.On("GetUserByUsername", mock.Anything).Return(nil, nil)

	//act
	suite.TokenController.PostToken(w, req, nil)

	//assert
	AssertErrorResponse(&suite.Suite, w.Result(), http.StatusBadRequest, "invalid", "username", "password")
}

func (suite *TokenControllerTestSuite) TestPostToken_PasswordGrant_WherePasswordDoesNotMatch_ReturnsBadRequest() {
	//arrange
	w := httptest.NewRecorder()

	body := controllers.PostTokenBody{
		GrantType: "password",
		PostTokenPasswordGrantBody: controllers.PostTokenPasswordGrantBody{
			Username: "username",
			Password: "password",
			ClientID: uuid.New().String(),
			Scope:    "scope",
		},
	}
	req := CreateRequest(&suite.Suite, "", body)

	suite.CRUDMock.On("GetClientByID", mock.Anything).Return(&models.Client{}, nil)
	suite.CRUDMock.On("GetScopeByName", mock.Anything).Return(&models.Scope{}, nil)
	suite.CRUDMock.On("GetUserByUsername", mock.Anything).Return(&models.User{}, nil)
	suite.PasswordHasherMock.On("ComparePasswords", mock.Anything, mock.Anything).Return(errors.New(""))

	//act
	suite.TokenController.PostToken(w, req, nil)

	//assert
	AssertErrorResponse(&suite.Suite, w.Result(), http.StatusBadRequest, "invalid", "username", "password")
}

func (suite *TokenControllerTestSuite) TestPostToken_PasswordGrant_WithErrorSavingAccessToken_ReturnsInternalServerError() {
	//arrange
	w := httptest.NewRecorder()

	body := controllers.PostTokenBody{
		GrantType: "password",
		PostTokenPasswordGrantBody: controllers.PostTokenPasswordGrantBody{
			Username: "username",
			Password: "password",
			ClientID: uuid.New().String(),
			Scope:    "scope",
		},
	}
	req := CreateRequest(&suite.Suite, "", body)

	suite.CRUDMock.On("GetClientByID", mock.Anything).Return(&models.Client{}, nil)
	suite.CRUDMock.On("GetScopeByName", mock.Anything).Return(&models.Scope{}, nil)
	suite.CRUDMock.On("GetUserByUsername", mock.Anything).Return(&models.User{}, nil)
	suite.PasswordHasherMock.On("ComparePasswords", mock.Anything, mock.Anything).Return(nil)
	suite.CRUDMock.On("SaveAccessToken", mock.Anything).Return(errors.New(""))

	//act
	suite.TokenController.PostToken(w, req, nil)

	//assert
	AssertInternalServerErrorResponse(&suite.Suite, w.Result())
}

func (suite *TokenControllerTestSuite) TestPostToken_PasswordGrant_WithValidRequest_ReturnsOK() {
	//arrange
	w := httptest.NewRecorder()

	clientID := uuid.New()
	body := controllers.PostTokenBody{
		GrantType: "password",
		PostTokenPasswordGrantBody: controllers.PostTokenPasswordGrantBody{
			Username: "username",
			Password: "password",
			ClientID: clientID.String(),
			Scope:    "scope",
		},
	}
	req := CreateRequest(&suite.Suite, "", body)

	var token *models.AccessToken
	client := &models.Client{ID: clientID}
	scope := &models.Scope{ID: uuid.New()}
	user := &models.User{ID: uuid.New()}

	suite.CRUDMock.On("GetClientByID", mock.Anything).Return(client, nil)
	suite.CRUDMock.On("GetScopeByName", mock.Anything).Return(scope, nil)
	suite.CRUDMock.On("GetUserByUsername", mock.Anything).Return(user, nil)
	suite.PasswordHasherMock.On("ComparePasswords", mock.Anything, mock.Anything).Return(nil)
	suite.CRUDMock.On("SaveAccessToken", mock.Anything).Run(func(args mock.Arguments) {
		token = args.Get(0).(*models.AccessToken)
	}).Return(nil)

	//act
	suite.TokenController.PostToken(w, req, nil)

	//assert
	suite.CRUDMock.AssertCalled(suite.T(), "GetClientByID", clientID)
	suite.CRUDMock.AssertCalled(suite.T(), "GetScopeByName", body.Scope)
	suite.CRUDMock.AssertCalled(suite.T(), "GetUserByUsername", body.Username)
	suite.PasswordHasherMock.AssertCalled(suite.T(), "ComparePasswords", mock.Anything, body.Password)

	suite.Equal(client, token.Client)
	suite.Equal(scope, token.Scope)
	suite.Equal(user, token.User)

	AssertAccessTokenResponse(&suite.Suite, w.Result(), token.ID.String())
}

func (suite *TokenControllerTestSuite) TestDeleteToken_AuthorizationHeaderTests() {
	setupTest := func() {
		suite.CRUDMock = databasemocks.CRUDOperations{}
		suite.TokenController.CRUD = &suite.CRUDMock
	}

	RunAuthHeaderTests(&suite.Suite, &suite.CRUDMock, setupTest, suite.TokenController.DeleteToken)
}

func (suite *TokenControllerTestSuite) TestDeleteToken_WithErrorDeletingAccessToken_ReturnsInternalServerError() {
	//arrange
	w := httptest.NewRecorder()
	req := CreateRequest(&suite.Suite, uuid.New().String(), nil)

	suite.CRUDMock.On("GetAccessTokenByID", mock.Anything).Return(&models.AccessToken{}, nil)
	suite.CRUDMock.On("DeleteAccessToken", mock.Anything).Return(errors.New(""))

	//act
	suite.TokenController.DeleteToken(w, req, nil)

	//assert
	AssertInternalServerErrorResponse(&suite.Suite, w.Result())
}

func (suite *TokenControllerTestSuite) TestDeleteToken_WithValidRequest_ReturnsOK() {
	//arrange
	w := httptest.NewRecorder()

	tokenID := uuid.New()
	accessToken := &models.AccessToken{}

	req := CreateRequest(&suite.Suite, tokenID.String(), nil)

	suite.CRUDMock.On("GetAccessTokenByID", mock.Anything).Return(accessToken, nil)
	suite.CRUDMock.On("DeleteAccessToken", mock.Anything).Return(nil)

	//act
	suite.TokenController.DeleteToken(w, req, nil)

	//assert
	suite.CRUDMock.AssertCalled(suite.T(), "GetAccessTokenByID", tokenID)
	suite.CRUDMock.AssertCalled(suite.T(), "DeleteAccessToken", accessToken)

	AssertSuccessResponse(&suite.Suite, w.Result())
}

func TestTokenControllerTestSuite(t *testing.T) {
	suite.Run(t, &TokenControllerTestSuite{})
}
