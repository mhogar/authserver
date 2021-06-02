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

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type UserHandlerSuite struct {
	suite.Suite
	ControllersMock   controllermocks.Controllers
	AuthenticatorMock mocks.Authenticator
	Router            *httprouter.Router
}

func (suite *UserHandlerSuite) SetupTest() {
	suite.ControllersMock = controllermocks.Controllers{}
	suite.AuthenticatorMock = mocks.Authenticator{}
	suite.Router = router.CreateRouter(&suite.ControllersMock, &suite.AuthenticatorMock)
}

func (suite *UserHandlerSuite) TestPostUser_WithInvalidJSONBody_ReturnsBadRequest() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	req := common.CreateRequest(&suite.Suite, http.MethodPost, server.URL+"/user", "", "invalid")

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	common.AssertErrorResponse(&suite.Suite, res, http.StatusBadRequest, "invalid json body")
}

func (suite *UserHandlerSuite) TestPostUser_WithClientErrorCreatingUser_ReturnsBadRequest() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	body := router.PostUserBody{
		Username: "username",
		Password: "password",
	}
	req := common.CreateRequest(&suite.Suite, http.MethodPost, server.URL+"/user", "", body)

	message := "create user error"
	suite.ControllersMock.On("CreateUser", mock.Anything, mock.Anything).Return(nil, requesterror.ClientError(message))

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	common.AssertErrorResponse(&suite.Suite, res, http.StatusBadRequest, message)
}

func (suite *UserHandlerSuite) TestPostUser_WithInternalErrorCreatingUser_ReturnsInternalServerError() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	body := router.PostUserBody{
		Username: "username",
		Password: "password",
	}
	req := common.CreateRequest(&suite.Suite, http.MethodPost, server.URL+"/user", "", body)

	suite.ControllersMock.On("CreateUser", mock.Anything, mock.Anything).Return(nil, requesterror.InternalError())

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	common.AssertInternalServerErrorResponse(&suite.Suite, res)
}

func (suite *UserHandlerSuite) TestPostUser_WithValidRequest_ReturnsSuccess() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	body := router.PostUserBody{
		Username: "username",
		Password: "password",
	}
	req := common.CreateRequest(&suite.Suite, http.MethodPost, server.URL+"/user", "", body)

	suite.ControllersMock.On("CreateUser", mock.Anything, mock.Anything).Return(nil, requesterror.NoError())

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	suite.ControllersMock.AssertCalled(suite.T(), "CreateUser", body.Username, body.Password)
	common.AssertSuccessResponse(&suite.Suite, res)
}

func (suite *UserHandlerSuite) TestPostUser_WithPanicTriggered_ReturnsInternalServerError() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	body := router.PostUserBody{
		Username: "username",
		Password: "password",
	}
	req := common.CreateRequest(&suite.Suite, http.MethodPost, server.URL+"/user", "", body)

	suite.ControllersMock.On("CreateUser", mock.Anything, mock.Anything).Run(func(_ mock.Arguments) {
		panic("test panic handler")
	})

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	common.AssertInternalServerErrorResponse(&suite.Suite, res)
}

func (suite *UserHandlerSuite) TestDeleteUser_WithClientErrorAuthenticatingUser_ReturnsUnauthorized() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	req := common.CreateRequest(&suite.Suite, http.MethodDelete, server.URL+"/user", "", nil)

	message := "authenticate error"
	suite.AuthenticatorMock.On("Authenticate", mock.Anything).Return(nil, requesterror.ClientError(message))

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	common.AssertErrorResponse(&suite.Suite, res, http.StatusUnauthorized, message)
}

func (suite *UserHandlerSuite) TestDeleteUser_WithInternalErrorAuthenticatingUser_ReturnsInternalServerError() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	req := common.CreateRequest(&suite.Suite, http.MethodDelete, server.URL+"/user", "", nil)

	suite.AuthenticatorMock.On("Authenticate", mock.Anything).Return(nil, requesterror.InternalError())

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	common.AssertInternalServerErrorResponse(&suite.Suite, res)
}

func (suite *UserHandlerSuite) TestDeleteUser_WithClientErrorDeletingUser_ReturnsBadRequest() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	token := &models.AccessToken{User: &models.User{}}
	req := common.CreateRequest(&suite.Suite, http.MethodDelete, server.URL+"/user", "", nil)

	message := "delete user error"
	suite.AuthenticatorMock.On("Authenticate", mock.Anything).Return(token, requesterror.NoError())
	suite.ControllersMock.On("DeleteUser", mock.Anything).Return(requesterror.ClientError(message))

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	common.AssertErrorResponse(&suite.Suite, res, http.StatusBadRequest, message)
}

func (suite *UserHandlerSuite) TestDeleteUser_WithInternalDeletingUser_ReturnsInternalServerError() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	token := &models.AccessToken{User: &models.User{}}
	req := common.CreateRequest(&suite.Suite, http.MethodDelete, server.URL+"/user", "", nil)

	suite.AuthenticatorMock.On("Authenticate", mock.Anything).Return(token, requesterror.NoError())
	suite.ControllersMock.On("DeleteUser", mock.Anything).Return(requesterror.InternalError())

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	common.AssertInternalServerErrorResponse(&suite.Suite, res)
}

func (suite *UserHandlerSuite) TestDeleteUser_WithValidRequest_ReturnsSuccess() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	token := &models.AccessToken{User: &models.User{}}
	req := common.CreateRequest(&suite.Suite, http.MethodDelete, server.URL+"/user", "", nil)

	suite.AuthenticatorMock.On("Authenticate", mock.Anything).Return(token, requesterror.NoError())
	suite.ControllersMock.On("DeleteUser", mock.Anything).Return(requesterror.NoError())

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	suite.AuthenticatorMock.AssertCalled(suite.T(), "Authenticate", mock.Anything)
	suite.ControllersMock.AssertCalled(suite.T(), "DeleteUser", token.User)
	common.AssertSuccessResponse(&suite.Suite, res)
}

func (suite *UserHandlerSuite) TestDeleteUser_WithPanicTriggered_ReturnsInternalServerError() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	token := &models.AccessToken{User: &models.User{}}
	req := common.CreateRequest(&suite.Suite, http.MethodDelete, server.URL+"/user", "", nil)

	suite.AuthenticatorMock.On("Authenticate", mock.Anything).Return(token, requesterror.NoError())
	suite.ControllersMock.On("DeleteUser", mock.Anything).Run(func(_ mock.Arguments) {
		panic("test panic handler")
	})

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	common.AssertInternalServerErrorResponse(&suite.Suite, res)
}

func (suite *UserHandlerSuite) TestUpdateUserPassword_WithClientErrorAuthenticatingUser_ReturnsUnauthorized() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	req := common.CreateRequest(&suite.Suite, http.MethodPatch, server.URL+"/user/password", "", nil)

	message := "authenticate error"
	suite.AuthenticatorMock.On("Authenticate", mock.Anything).Return(nil, requesterror.ClientError(message))

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	common.AssertErrorResponse(&suite.Suite, res, http.StatusUnauthorized, message)
}

func (suite *UserHandlerSuite) TestUpdateUserPassword_WithInternalErrorAuthenticatingUser_ReturnsInternalServerError() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	req := common.CreateRequest(&suite.Suite, http.MethodPatch, server.URL+"/user/password", "", nil)

	suite.AuthenticatorMock.On("Authenticate", mock.Anything).Return(nil, requesterror.InternalError())

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	common.AssertInternalServerErrorResponse(&suite.Suite, res)
}

func (suite *UserHandlerSuite) TestUpdateUserPassword_WithInvalidJSONBody_ReturnsBadRequest() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	req := common.CreateRequest(&suite.Suite, http.MethodPatch, server.URL+"/user/password", "", "invalid")

	token := &models.AccessToken{User: &models.User{}}
	suite.AuthenticatorMock.On("Authenticate", mock.Anything).Return(token, requesterror.NoError())

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	common.AssertErrorResponse(&suite.Suite, res, http.StatusBadRequest, "invalid json body")
}

func (suite *UserHandlerSuite) TestUpdateUserPassword_WithClientErrorUpdatingUserPassword_ReturnsBadRequest() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	body := router.PatchUserPasswordBody{
		OldPassword: "old password",
		NewPassword: "new password",
	}

	req := common.CreateRequest(&suite.Suite, http.MethodPatch, server.URL+"/user/password", "", body)

	token := &models.AccessToken{User: &models.User{}}
	suite.AuthenticatorMock.On("Authenticate", mock.Anything).Return(token, requesterror.NoError())

	message := "update user password error"
	suite.ControllersMock.On("UpdateUserPassword", mock.Anything, mock.Anything, mock.Anything).Return(requesterror.ClientError(message))

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	common.AssertErrorResponse(&suite.Suite, res, http.StatusBadRequest, message)
}

func (suite *UserHandlerSuite) TestUpdateUserPassword_WithInternalErrorUpdatingUserPassword_ReturnsInternalServerError() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	body := router.PatchUserPasswordBody{
		OldPassword: "old password",
		NewPassword: "new password",
	}

	req := common.CreateRequest(&suite.Suite, http.MethodPatch, server.URL+"/user/password", "", body)

	token := &models.AccessToken{User: &models.User{}}
	suite.AuthenticatorMock.On("Authenticate", mock.Anything).Return(token, requesterror.NoError())
	suite.ControllersMock.On("UpdateUserPassword", mock.Anything, mock.Anything, mock.Anything).Return(requesterror.InternalError())

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	common.AssertInternalServerErrorResponse(&suite.Suite, res)
}

func (suite *UserHandlerSuite) TestUpdateUserPassword_WithValidRequest_ReturnsSuccess() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	body := router.PatchUserPasswordBody{
		OldPassword: "old password",
		NewPassword: "new password",
	}

	req := common.CreateRequest(&suite.Suite, http.MethodPatch, server.URL+"/user/password", "", body)

	token := &models.AccessToken{User: &models.User{}}
	suite.AuthenticatorMock.On("Authenticate", mock.Anything).Return(token, requesterror.NoError())
	suite.ControllersMock.On("UpdateUserPassword", mock.Anything, mock.Anything, mock.Anything).Return(requesterror.NoError())

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	suite.AuthenticatorMock.AssertCalled(suite.T(), "Authenticate", mock.Anything)
	suite.ControllersMock.AssertCalled(suite.T(), "UpdateUserPassword", token.User, body.OldPassword, body.NewPassword)
	common.AssertSuccessResponse(&suite.Suite, res)
}

func (suite *UserHandlerSuite) TestUpdateUserPassword_WithPanicTriggered_ReturnsInternalServerError() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	body := router.PatchUserPasswordBody{
		OldPassword: "old password",
		NewPassword: "new password",
	}

	req := common.CreateRequest(&suite.Suite, http.MethodPatch, server.URL+"/user/password", "", body)

	token := &models.AccessToken{User: &models.User{}}
	suite.AuthenticatorMock.On("Authenticate", mock.Anything).Return(token, requesterror.NoError())
	suite.ControllersMock.On("UpdateUserPassword", mock.Anything, mock.Anything, mock.Anything).Run(func(_ mock.Arguments) {
		panic("test panic handler")
	})

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	common.AssertInternalServerErrorResponse(&suite.Suite, res)
}

func TestUserHandlerTestSuite(t *testing.T) {
	suite.Run(t, &UserHandlerSuite{})
}
