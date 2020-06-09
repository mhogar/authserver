package controllers_test

import (
	"authserver/controllers"
	"authserver/models"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"

	helpermocks "authserver/helpers/mocks"
	modelmocks "authserver/models/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type TokenControllerTestSuite struct {
	suite.Suite
	UserCRUDMock        modelmocks.UserCRUD
	ClientCRUDMock      modelmocks.ClientCRUD
	ScopeCRUDMock       modelmocks.ScopeCRUD
	AccessTokenCRUDMock modelmocks.AccessTokenCRUD
	PasswordHasherMock  helpermocks.PasswordHasher
	TokenController     controllers.TokenController
}

func (suite *TokenControllerTestSuite) SetupTest() {
	suite.UserCRUDMock = modelmocks.UserCRUD{}
	suite.ClientCRUDMock = modelmocks.ClientCRUD{}
	suite.ScopeCRUDMock = modelmocks.ScopeCRUD{}
	suite.AccessTokenCRUDMock = modelmocks.AccessTokenCRUD{}
	suite.PasswordHasherMock = helpermocks.PasswordHasher{}
	suite.TokenController = controllers.TokenController{
		UserCRUD:        &suite.UserCRUDMock,
		ClientCRUD:      &suite.ClientCRUDMock,
		ScopeCRUD:       &suite.ScopeCRUDMock,
		AccessTokenCRUD: &suite.AccessTokenCRUDMock,
		PasswordHasher:  &suite.PasswordHasherMock,
	}
}

func (suite *TokenControllerTestSuite) TestPostToken_WithInvalidJSONBody_ReturnsInvalidRequest() {
	//arrange
	w := httptest.NewRecorder()

	req, err := http.NewRequest("", "", strings.NewReader("0"))
	suite.Require().NoError(err)

	//act
	suite.TokenController.PostToken(w, req, nil)

	//assert
	assertOAuthErrorResponse(&suite.Suite, w.Result(), http.StatusBadRequest, "invalid_request", "invalid json body")
}

func (suite *TokenControllerTestSuite) TestPostToken_WithMissingGrantType_ReturnsInvalidRequest() {
	//arrange
	w := httptest.NewRecorder()

	body := controllers.PostTokenBody{}
	req := createRequestWithJSONBody(&suite.Suite, body)

	//act
	suite.TokenController.PostToken(w, req, nil)

	//assert
	assertOAuthErrorResponse(&suite.Suite, w.Result(), http.StatusBadRequest, "invalid_request", "missing grant_type parameter")
}

func (suite *TokenControllerTestSuite) TestPostToken_WithUnsupportedGrantType_ReturnsUnsupportedGrantType() {
	//arrange
	w := httptest.NewRecorder()

	body := controllers.PostTokenBody{
		GrantType: "unsupported",
	}
	req := createRequestWithJSONBody(&suite.Suite, body)

	//act
	suite.TokenController.PostToken(w, req, nil)

	//assert
	assertOAuthErrorResponse(&suite.Suite, w.Result(), http.StatusBadRequest, "unsupported_grant_type", "")
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
		req := createRequestWithJSONBody(&suite.Suite, body)

		//act
		suite.TokenController.PostToken(w, req, nil)

		//assert
		assertOAuthErrorResponse(&suite.Suite, w.Result(), http.StatusBadRequest, "invalid_request", expectedErrorDescription)
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
	req := createRequestWithJSONBody(&suite.Suite, body)

	//act
	suite.TokenController.PostToken(w, req, nil)

	//assert
	assertOAuthErrorResponse(&suite.Suite, w.Result(), http.StatusUnauthorized, "invalid_client", "client_id was in invalid format")
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
	req := createRequestWithJSONBody(&suite.Suite, body)

	suite.ClientCRUDMock.On("GetClientByID", mock.Anything).Return(nil, errors.New(""))

	//act
	suite.TokenController.PostToken(w, req, nil)

	//assert
	assertInternalServerErrorResponse(&suite.Suite, w.Result())
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
	req := createRequestWithJSONBody(&suite.Suite, body)

	suite.ClientCRUDMock.On("GetClientByID", mock.Anything).Return(nil, nil)

	//act
	suite.TokenController.PostToken(w, req, nil)

	//assert
	assertOAuthErrorResponse(&suite.Suite, w.Result(), http.StatusUnauthorized, "invalid_client", "")
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
	req := createRequestWithJSONBody(&suite.Suite, body)

	suite.ClientCRUDMock.On("GetClientByID", mock.Anything).Return(&models.Client{}, nil)
	suite.ScopeCRUDMock.On("GetScopeByName", mock.Anything).Return(nil, errors.New(""))

	//act
	suite.TokenController.PostToken(w, req, nil)

	//assert
	assertInternalServerErrorResponse(&suite.Suite, w.Result())
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
	req := createRequestWithJSONBody(&suite.Suite, body)

	suite.ClientCRUDMock.On("GetClientByID", mock.Anything).Return(&models.Client{}, nil)
	suite.ScopeCRUDMock.On("GetScopeByName", mock.Anything).Return(nil, nil)

	//act
	suite.TokenController.PostToken(w, req, nil)

	//assert
	assertOAuthErrorResponse(&suite.Suite, w.Result(), http.StatusBadRequest, "invalid_scope", "")
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
	req := createRequestWithJSONBody(&suite.Suite, body)

	suite.ClientCRUDMock.On("GetClientByID", mock.Anything).Return(&models.Client{}, nil)
	suite.ScopeCRUDMock.On("GetScopeByName", mock.Anything).Return(&models.Scope{}, nil)
	suite.UserCRUDMock.On("GetUserByUsername", mock.Anything).Return(nil, errors.New(""))

	//act
	suite.TokenController.PostToken(w, req, nil)

	//assert
	assertInternalServerErrorResponse(&suite.Suite, w.Result())
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
	req := createRequestWithJSONBody(&suite.Suite, body)

	suite.ClientCRUDMock.On("GetClientByID", mock.Anything).Return(&models.Client{}, nil)
	suite.ScopeCRUDMock.On("GetScopeByName", mock.Anything).Return(&models.Scope{}, nil)
	suite.UserCRUDMock.On("GetUserByUsername", mock.Anything).Return(nil, nil)

	//act
	suite.TokenController.PostToken(w, req, nil)

	//assert
	assertErrorResponse(&suite.Suite, w.Result(), http.StatusBadRequest, "invalid username and/or password")
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
	req := createRequestWithJSONBody(&suite.Suite, body)

	suite.ClientCRUDMock.On("GetClientByID", mock.Anything).Return(&models.Client{}, nil)
	suite.ScopeCRUDMock.On("GetScopeByName", mock.Anything).Return(&models.Scope{}, nil)
	suite.UserCRUDMock.On("GetUserByUsername", mock.Anything).Return(&models.User{}, nil)
	suite.PasswordHasherMock.On("ComparePasswords", mock.Anything, mock.Anything).Return(errors.New(""))

	//act
	suite.TokenController.PostToken(w, req, nil)

	//assert
	assertErrorResponse(&suite.Suite, w.Result(), http.StatusBadRequest, "invalid username and/or password")
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
	req := createRequestWithJSONBody(&suite.Suite, body)

	suite.ClientCRUDMock.On("GetClientByID", mock.Anything).Return(&models.Client{}, nil)
	suite.ScopeCRUDMock.On("GetScopeByName", mock.Anything).Return(&models.Scope{}, nil)
	suite.UserCRUDMock.On("GetUserByUsername", mock.Anything).Return(&models.User{}, nil)
	suite.PasswordHasherMock.On("ComparePasswords", mock.Anything, mock.Anything).Return(nil)
	suite.AccessTokenCRUDMock.On("SaveAccessToken", mock.Anything).Return(errors.New(""))

	//act
	suite.TokenController.PostToken(w, req, nil)

	//assert
	assertInternalServerErrorResponse(&suite.Suite, w.Result())
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
	req := createRequestWithJSONBody(&suite.Suite, body)

	var token *models.AccessToken
	client := &models.Client{ID: clientID}
	scope := &models.Scope{ID: uuid.New()}
	user := &models.User{ID: uuid.New()}

	suite.ClientCRUDMock.On("GetClientByID", mock.Anything).Return(client, nil)
	suite.ScopeCRUDMock.On("GetScopeByName", mock.Anything).Return(scope, nil)
	suite.UserCRUDMock.On("GetUserByUsername", mock.Anything).Return(user, nil)
	suite.PasswordHasherMock.On("ComparePasswords", mock.Anything, mock.Anything).Return(nil)
	suite.AccessTokenCRUDMock.On("SaveAccessToken", mock.Anything).Run(func(args mock.Arguments) {
		token = args.Get(0).(*models.AccessToken)
	}).Return(nil)

	//act
	suite.TokenController.PostToken(w, req, nil)

	//assert
	suite.ClientCRUDMock.AssertCalled(suite.T(), "GetClientByID", clientID)
	suite.ScopeCRUDMock.AssertCalled(suite.T(), "GetScopeByName", body.Scope)
	suite.UserCRUDMock.AssertCalled(suite.T(), "GetUserByUsername", body.Username)
	suite.PasswordHasherMock.AssertCalled(suite.T(), "ComparePasswords", mock.Anything, body.Password)

	suite.Equal(client.ID, token.ClientID)
	suite.Equal(scope.ID, token.ScopeID)
	suite.Equal(user.ID, token.UserID)

	assertAccessTokenResponse(&suite.Suite, w.Result(), token.ID.String())
}

func (suite *TokenControllerTestSuite) TestDeleteToken_WithNoBearerToken_ReturnsUnauthorized() {
	var req *http.Request

	testCase := func() {
		//arrange
		w := httptest.NewRecorder()

		//act
		suite.TokenController.DeleteToken(w, req, nil)

		//assert
		assertErrorResponse(&suite.Suite, w.Result(), http.StatusUnauthorized, "no bearer token")
	}

	req = createEmptyRequest(&suite.Suite)
	suite.Run("NoAuthorizationHeader", testCase)

	req.Header.Set("Authorization", "invalid")
	suite.Run("AuthorizationHeaderDoesNotContainBearerToken", testCase)
}

func (suite *TokenControllerTestSuite) TestDeleteToken_WithBearerTokenInInvalidFormat_ReturnsUnauthorized() {
	//arrange
	w := httptest.NewRecorder()
	req := createRequestWithAuthorizationHeader(&suite.Suite, "invalid")

	//act
	suite.TokenController.DeleteToken(w, req, nil)

	//assert
	assertErrorResponse(&suite.Suite, w.Result(), http.StatusUnauthorized, "bearer token", "invalid format")
}

func (suite *TokenControllerTestSuite) TestDeleteToken_WithErrorFetchingAccessTokenByID_ReturnsInternalServerError() {
	//arrange
	w := httptest.NewRecorder()
	req := createRequestWithAuthorizationHeader(&suite.Suite, uuid.New().String())

	suite.AccessTokenCRUDMock.On("GetAccessTokenByID", mock.Anything).Return(nil, errors.New(""))

	//act
	suite.TokenController.DeleteToken(w, req, nil)

	//assert
	assertInternalServerErrorResponse(&suite.Suite, w.Result())
}

func (suite *TokenControllerTestSuite) TestDeleteToken_WhereAccessTokenWithIDisNotFound_ReturnsUnauthorized() {
	//arrange
	w := httptest.NewRecorder()
	req := createRequestWithAuthorizationHeader(&suite.Suite, uuid.New().String())

	suite.AccessTokenCRUDMock.On("GetAccessTokenByID", mock.Anything).Return(nil, nil)

	//act
	suite.TokenController.DeleteToken(w, req, nil)

	//assert
	assertErrorResponse(&suite.Suite, w.Result(), http.StatusUnauthorized, "bearer token", "invalid")
}

func (suite *TokenControllerTestSuite) TestDeleteToken_WithErrorDeletingAccessToken_ReturnsInternalServerError() {
	//arrange
	w := httptest.NewRecorder()
	req := createRequestWithAuthorizationHeader(&suite.Suite, uuid.New().String())

	suite.AccessTokenCRUDMock.On("GetAccessTokenByID", mock.Anything).Return(&models.AccessToken{}, nil)
	suite.AccessTokenCRUDMock.On("DeleteAccessToken", mock.Anything).Return(errors.New(""))

	//act
	suite.TokenController.DeleteToken(w, req, nil)

	//assert
	assertInternalServerErrorResponse(&suite.Suite, w.Result())
}

func (suite *TokenControllerTestSuite) TestDeleteToken_WithValidRequest_ReturnsOK() {
	//arrange
	w := httptest.NewRecorder()
	req := createRequestWithAuthorizationHeader(&suite.Suite, uuid.New().String())

	suite.AccessTokenCRUDMock.On("GetAccessTokenByID", mock.Anything).Return(&models.AccessToken{}, nil)
	suite.AccessTokenCRUDMock.On("DeleteAccessToken", mock.Anything).Return(nil)

	//act
	suite.TokenController.DeleteToken(w, req, nil)

	//assert
	assertSuccessResponse(&suite.Suite, w.Result())
}

func TestTokenControllerTestSuite(t *testing.T) {
	suite.Run(t, &TokenControllerTestSuite{})
}
