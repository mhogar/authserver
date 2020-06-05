package api_test

import (
	"authserver/controllers/api"
	"authserver/controllers/common"
	databasemocks "authserver/database/mocks"
	"authserver/helpers"
	helpermocks "authserver/helpers/mocks"
	"authserver/models"
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type UserControllerTestSuite struct {
	suite.Suite
	SID                           uuid.UUID
	SessionCookie                 *http.Cookie
	UserCRUDMock                  databasemocks.UserCRUD
	PasswordHasherMock            helpermocks.PasswordHasher
	PasswordCriteriaValidatorMock helpermocks.PasswordCriteriaValidator
	UserController                api.UserController
}

func (suite *UserControllerTestSuite) SetupTest() {
	suite.SID = uuid.New()
	suite.SessionCookie = &http.Cookie{
		Name:  "session",
		Value: suite.SID.String(),
	}

	suite.UserCRUDMock = databasemocks.UserCRUD{}
	suite.PasswordHasherMock = helpermocks.PasswordHasher{}
	suite.PasswordCriteriaValidatorMock = helpermocks.PasswordCriteriaValidator{}
	suite.UserController = api.UserController{
		UserCRUD:                  &suite.UserCRUDMock,
		PasswordHasher:            &suite.PasswordHasherMock,
		PasswordCriteriaValidator: &suite.PasswordCriteriaValidatorMock,
	}
}

func (suite *UserControllerTestSuite) TestPostAPIUser_WithInvalidJSONBody_ReturnsBadRequest() {
	//arrange
	w := httptest.NewRecorder()

	req, err := http.NewRequest("", "", strings.NewReader("0"))
	suite.Require().NoError(err)

	//act
	suite.UserController.PostAPIUser(w, req, nil)

	//assert
	common.AssertErrorResponse(&suite.Suite, w.Result(), http.StatusBadRequest, "invalid json body")
}

func (suite *UserControllerTestSuite) TestPostAPIUser_WithInvalidBodyFields_ReturnsBadRequest() {
	var body api.PostAPIUserBody

	testCase := func() {
		//arrange
		w := httptest.NewRecorder()
		req := common.CreateRequestWithJSONBody(&suite.Suite, body)

		//act
		suite.UserController.PostAPIUser(w, req, nil)

		//assert
		common.AssertErrorResponse(&suite.Suite, w.Result(), http.StatusBadRequest, "username and password cannot be empty")
	}

	body = api.PostAPIUserBody{
		Username: "",
		Password: "password",
	}
	suite.Run("EmptyUsername", testCase)

	body = api.PostAPIUserBody{
		Username: "username",
		Password: "",
	}
	suite.Run("EmptyPassword", testCase)
}

func (suite *UserControllerTestSuite) TestPostAPIUser_WithErrorGettingUserByUsername_ReturnsInternalServerError() {
	//arrange
	w := httptest.NewRecorder()

	body := api.PostAPIUserBody{
		Username: "username",
		Password: "password",
	}
	req := common.CreateRequestWithJSONBody(&suite.Suite, body)

	suite.UserCRUDMock.On("GetUserByUsername", body.Username).Return(nil, errors.New(""))

	//act
	suite.UserController.PostAPIUser(w, req, nil)

	//assert
	common.AssertInternalServerErrorResponse(&suite.Suite, w.Result())
}

func (suite *UserControllerTestSuite) TestPostAPIUser_WithNonUniqueUsername_ReturnsBadRequest() {
	//arrange
	w := httptest.NewRecorder()

	body := api.PostAPIUserBody{
		Username: "username",
		Password: "password",
	}
	req := common.CreateRequestWithJSONBody(&suite.Suite, body)

	suite.UserCRUDMock.On("GetUserByUsername", body.Username).Return(&models.User{}, nil)

	//act
	suite.UserController.PostAPIUser(w, req, nil)

	//assert
	common.AssertErrorResponse(&suite.Suite, w.Result(), http.StatusBadRequest, "username already exists")
}

func (suite *UserControllerTestSuite) TestPostAPIUser_WherePasswordDoesNotMeetCriteria_ReturnsBadRequest() {
	//arrange
	w := httptest.NewRecorder()

	body := api.PostAPIUserBody{
		Username: "username",
		Password: "password",
	}
	req := common.CreateRequestWithJSONBody(&suite.Suite, body)

	suite.UserCRUDMock.On("GetUserByUsername", body.Username).Return(nil, nil)
	suite.PasswordCriteriaValidatorMock.On("ValidatePasswordCriteria", mock.Anything).Return(helpers.CreateValidatePasswordCriteriaValid())
	suite.PasswordHasherMock.On("HashPassword", mock.Anything).Return(nil, errors.New(""))

	//act
	suite.UserController.PostAPIUser(w, req, nil)

	//assert
	common.AssertInternalServerErrorResponse(&suite.Suite, w.Result())
}

func (suite *UserControllerTestSuite) TestPostAPIUser_WithErrorHashingNewPassword_ReturnsInternalServerError() {
	//arrange
	w := httptest.NewRecorder()

	body := api.PostAPIUserBody{
		Username: "username",
		Password: "password",
	}
	req := common.CreateRequestWithJSONBody(&suite.Suite, body)

	suite.UserCRUDMock.On("GetUserByUsername", body.Username).Return(nil, nil)
	suite.PasswordCriteriaValidatorMock.On("ValidatePasswordCriteria", mock.Anything).Return(helpers.CreateValidatePasswordCriteriaError(helpers.ValidatePasswordCriteriaTooShort, ""))

	//act
	suite.UserController.PostAPIUser(w, req, nil)

	//assert
	common.AssertErrorResponse(&suite.Suite, w.Result(), http.StatusBadRequest, "password does not meet minimum criteria")
}

func (suite *UserControllerTestSuite) TestPostAPIUser_WithErrorCreatingUser_ReturnsInternalServerError() {
	//arrange
	w := httptest.NewRecorder()

	body := api.PostAPIUserBody{
		Username: "username",
		Password: "password",
	}
	req := common.CreateRequestWithJSONBody(&suite.Suite, body)

	suite.UserCRUDMock.On("GetUserByUsername", body.Username).Return(nil, nil)
	suite.PasswordHasherMock.On("HashPassword", mock.Anything).Return(nil, nil)
	suite.PasswordCriteriaValidatorMock.On("ValidatePasswordCriteria", mock.Anything).Return(helpers.CreateValidatePasswordCriteriaValid())
	suite.UserCRUDMock.On("CreateUser", mock.Anything).Return(errors.New(""))

	//act
	suite.UserController.PostAPIUser(w, req, nil)

	//assert
	common.AssertInternalServerErrorResponse(&suite.Suite, w.Result())
}

func (suite *UserControllerTestSuite) TestPostAPIUser_WithValidRequest_ReturnsOK() {
	//arrange
	w := httptest.NewRecorder()

	body := api.PostAPIUserBody{
		Username: "username",
		Password: "password",
	}
	req := common.CreateRequestWithJSONBody(&suite.Suite, body)

	hash := []byte("password hash")

	suite.UserCRUDMock.On("GetUserByUsername", body.Username).Return(nil, nil)
	suite.PasswordCriteriaValidatorMock.On("ValidatePasswordCriteria", mock.Anything).Return(helpers.CreateValidatePasswordCriteriaValid())
	suite.PasswordHasherMock.On("HashPassword", mock.Anything).Return(hash, nil)
	suite.UserCRUDMock.On("CreateUser", mock.Anything).Return(nil)

	//act
	suite.UserController.PostAPIUser(w, req, nil)

	//assert
	suite.UserCRUDMock.AssertCalled(suite.T(), "GetUserByUsername", body.Username)
	suite.PasswordCriteriaValidatorMock.AssertCalled(suite.T(), "ValidatePasswordCriteria", body.Password)
	suite.PasswordHasherMock.AssertCalled(suite.T(), "HashPassword", body.Password)
	suite.UserCRUDMock.AssertCalled(suite.T(), "CreateUser", mock.MatchedBy(func(u *models.User) bool {
		return u.Username == body.Username && bytes.Equal(u.PasswordHash, hash)
	}))

	common.AssertSuccessResponse(&suite.Suite, w.Result())
}

func (suite *UserControllerTestSuite) TestDeleteAPIUser_WithoutIdInParams_ReturnsBadRequest() {
	//arrange
	w := httptest.NewRecorder()

	//act
	suite.UserController.DeleteAPIUser(w, nil, make(httprouter.Params, 0))

	//assert
	common.AssertErrorResponse(&suite.Suite, w.Result(), http.StatusBadRequest, "id must be present")
}

func (suite *UserControllerTestSuite) TestDeleteAPIUser_WithIdInInvalidFormat_ReturnsBadRequest() {
	//arrange
	w := httptest.NewRecorder()

	id := 0
	params := httprouter.Params{
		httprouter.Param{Key: "id", Value: string(id)},
	}

	//act
	suite.UserController.DeleteAPIUser(w, nil, params)

	//assert
	common.AssertErrorResponse(&suite.Suite, w.Result(), http.StatusBadRequest, "id is in invalid format")
}

func (suite *UserControllerTestSuite) TestDeleteAPIUser_WithErrorGettingUserById_ReturnsInternalServerError() {
	//arrange
	w := httptest.NewRecorder()

	id := uuid.New()
	params := httprouter.Params{
		httprouter.Param{Key: "id", Value: id.String()},
	}

	suite.UserCRUDMock.On("GetUserByID", mock.Anything).Return(nil, errors.New(""))

	//act
	suite.UserController.DeleteAPIUser(w, nil, params)

	//assert
	common.AssertInternalServerErrorResponse(&suite.Suite, w.Result())
}

func (suite *UserControllerTestSuite) TestDeleteAPIUser_WhereUserIsNotFound_ReturnsBadRequest() {
	//arrange
	w := httptest.NewRecorder()

	id := uuid.New()
	params := httprouter.Params{
		httprouter.Param{Key: "id", Value: id.String()},
	}

	suite.UserCRUDMock.On("GetUserByID", mock.Anything).Return(nil, nil)

	//act
	suite.UserController.DeleteAPIUser(w, nil, params)

	//assert
	common.AssertErrorResponse(&suite.Suite, w.Result(), http.StatusBadRequest, "user not found")
}

func (suite *UserControllerTestSuite) TestDeleteAPIUser_WithErrorDeletingUser_ReturnsInternalServerError() {
	//arrange
	w := httptest.NewRecorder()

	user := models.CreateNewUser("username", []byte("password hash"))
	params := httprouter.Params{
		httprouter.Param{Key: "id", Value: user.ID.String()},
	}

	suite.UserCRUDMock.On("GetUserByID", mock.Anything).Return(user, nil)
	suite.UserCRUDMock.On("DeleteUser", mock.Anything).Return(errors.New(""))

	//act
	suite.UserController.DeleteAPIUser(w, nil, params)

	//assert
	common.AssertInternalServerErrorResponse(&suite.Suite, w.Result())
}

func (suite *UserControllerTestSuite) TestDeleteAPIUser_WithValidRequest_ReturnsOK() {
	//arrange
	w := httptest.NewRecorder()

	user := models.CreateNewUser("username", []byte("password hash"))
	params := httprouter.Params{
		httprouter.Param{Key: "id", Value: user.ID.String()},
	}

	suite.UserCRUDMock.On("GetUserByID", mock.Anything).Return(user, nil)
	suite.UserCRUDMock.On("DeleteUser", mock.Anything).Return(nil)

	//act
	suite.UserController.DeleteAPIUser(w, nil, params)

	//assert
	suite.UserCRUDMock.AssertCalled(suite.T(), "GetUserByID", user.ID)
	suite.UserCRUDMock.AssertCalled(suite.T(), "DeleteUser", user)

	common.AssertSuccessResponse(&suite.Suite, w.Result())
}

func (suite *UserControllerTestSuite) TestPatchAPIUserPassword_WithNoSessionId_ReturnsUnauthorized() {
	//arrange
	w := httptest.NewRecorder()

	req, err := http.NewRequest("", "", nil)
	suite.Require().NoError(err)

	//act
	suite.UserController.PatchAPIUserPassword(w, req, nil)

	//assert
	common.AssertErrorResponse(&suite.Suite, w.Result(), http.StatusUnauthorized, "token not provided")
}

func (suite *UserControllerTestSuite) TestPatchAPIUserPassword_WithInvalidSessionId_ReturnsUnauthorized() {
	//arrange
	w := httptest.NewRecorder()

	req, err := http.NewRequest("", "", nil)
	suite.Require().NoError(err)

	suite.SessionCookie.Value = "invalid session id"
	req.AddCookie(suite.SessionCookie)

	//act
	suite.UserController.PatchAPIUserPassword(w, req, nil)

	//assert
	common.AssertErrorResponse(&suite.Suite, w.Result(), http.StatusUnauthorized, "invalid format")
}

func (suite *UserControllerTestSuite) TestPatchAPIUserPassword_WithErrorGettingUserBySessionId_ReturnsInternalServerError() {
	//arrange
	w := httptest.NewRecorder()

	req, err := http.NewRequest("", "", nil)
	suite.Require().NoError(err)
	req.AddCookie(suite.SessionCookie)

	suite.UserCRUDMock.On("GetUserBySessionID", mock.Anything).Return(nil, errors.New(""))

	//act
	suite.UserController.PatchAPIUserPassword(w, req, nil)

	//assert
	common.AssertInternalServerErrorResponse(&suite.Suite, w.Result())
}

func (suite *UserControllerTestSuite) TestPatchAPIUserPassword_WhereNoUserIsFound_ReturnsUnauthorized() {
	//arrange
	w := httptest.NewRecorder()

	req, err := http.NewRequest("", "", nil)
	suite.Require().NoError(err)
	req.AddCookie(suite.SessionCookie)

	suite.UserCRUDMock.On("GetUserBySessionID", mock.Anything).Return(nil, nil)

	//act
	suite.UserController.PatchAPIUserPassword(w, req, nil)

	//assert
	common.AssertErrorResponse(&suite.Suite, w.Result(), http.StatusUnauthorized, "no user")
}

func (suite *UserControllerTestSuite) TestPatchAPIUserPassword_WithInvalidJSONBody_ReturnsBadRequest() {
	//arrange
	w := httptest.NewRecorder()

	req, err := http.NewRequest("", "", strings.NewReader("0"))
	suite.Require().NoError(err)
	req.AddCookie(suite.SessionCookie)

	suite.UserCRUDMock.On("GetUserBySessionID", mock.Anything).Return(&models.User{}, nil)

	//act
	suite.UserController.PatchAPIUserPassword(w, req, nil)

	//assert
	common.AssertErrorResponse(&suite.Suite, w.Result(), http.StatusBadRequest, "invalid json body")
}

func (suite *UserControllerTestSuite) TestPatchAPIUserPassword_WithInvalidBodyFields_ReturnsBadRequest() {
	var body api.PatchAPIUserPasswordBody

	testCase := func() {
		//arrange
		w := httptest.NewRecorder()

		req := common.CreateRequestWithJSONBody(&suite.Suite, body)
		req.AddCookie(suite.SessionCookie)

		suite.UserCRUDMock.On("GetUserBySessionID", mock.Anything).Return(&models.User{}, nil)

		//act
		suite.UserController.PatchAPIUserPassword(w, req, nil)

		//assert
		common.AssertErrorResponse(&suite.Suite, w.Result(), http.StatusBadRequest, "old password and new password cannot be empty")
	}

	body = api.PatchAPIUserPasswordBody{
		OldPassword: "",
		NewPassword: "new password",
	}
	suite.Run("EmptyUsername", testCase)

	body = api.PatchAPIUserPasswordBody{
		OldPassword: "old password",
		NewPassword: "",
	}
	suite.Run("EmptyPassword", testCase)
}

func (suite *UserControllerTestSuite) TestPatchAPIUserPassword_WhereOldPasswordIsInvalid_ReturnsBadRequest() {
	//arrange
	w := httptest.NewRecorder()

	body := api.PatchAPIUserPasswordBody{
		OldPassword: "old password",
		NewPassword: "new password",
	}

	req := common.CreateRequestWithJSONBody(&suite.Suite, body)
	req.AddCookie(suite.SessionCookie)

	suite.UserCRUDMock.On("GetUserBySessionID", mock.Anything).Return(&models.User{}, nil)
	suite.PasswordHasherMock.On("ComparePasswords", mock.Anything, mock.Anything).Return(errors.New(""))

	//act
	suite.UserController.PatchAPIUserPassword(w, req, nil)

	//assert
	common.AssertErrorResponse(&suite.Suite, w.Result(), http.StatusBadRequest, "old password is invalid")
}

func (suite *UserControllerTestSuite) TestPatchAPIUserPassword_WhereNewPasswordDoesNotMeetCriteria_ReturnsBadRequest() {
	//arrange
	w := httptest.NewRecorder()

	body := api.PatchAPIUserPasswordBody{
		OldPassword: "old password",
		NewPassword: "new password",
	}

	req := common.CreateRequestWithJSONBody(&suite.Suite, body)
	req.AddCookie(suite.SessionCookie)

	suite.UserCRUDMock.On("GetUserBySessionID", mock.Anything).Return(&models.User{}, nil)
	suite.PasswordHasherMock.On("ComparePasswords", mock.Anything, mock.Anything).Return(nil)
	suite.PasswordCriteriaValidatorMock.On("ValidatePasswordCriteria", mock.Anything).Return(helpers.CreateValidatePasswordCriteriaError(helpers.ValidatePasswordCriteriaTooShort, ""))

	//act
	suite.UserController.PatchAPIUserPassword(w, req, nil)

	//assert
	common.AssertErrorResponse(&suite.Suite, w.Result(), http.StatusBadRequest, "password does not meet minimum criteria")
}

func (suite *UserControllerTestSuite) TestPatchAPIUserPassword_WithErrorHashingNewPassword_ReturnsInternalServerError() {
	//arrange
	w := httptest.NewRecorder()

	body := api.PatchAPIUserPasswordBody{
		OldPassword: "old password",
		NewPassword: "new password",
	}

	req := common.CreateRequestWithJSONBody(&suite.Suite, body)
	req.AddCookie(suite.SessionCookie)

	suite.UserCRUDMock.On("GetUserBySessionID", mock.Anything).Return(&models.User{}, nil)
	suite.PasswordHasherMock.On("ComparePasswords", mock.Anything, mock.Anything).Return(nil)
	suite.PasswordCriteriaValidatorMock.On("ValidatePasswordCriteria", mock.Anything).Return(helpers.CreateValidatePasswordCriteriaValid())
	suite.PasswordHasherMock.On("HashPassword", mock.Anything).Return(nil, errors.New(""))

	//act
	suite.UserController.PatchAPIUserPassword(w, req, nil)

	//assert
	common.AssertInternalServerErrorResponse(&suite.Suite, w.Result())
}

func (suite *UserControllerTestSuite) TestPatchAPIUserPassword_WithErrorUpdatingUser_ReturnsInternalServerError() {
	//arrange
	w := httptest.NewRecorder()

	body := api.PatchAPIUserPasswordBody{
		OldPassword: "old password",
		NewPassword: "new password",
	}

	req := common.CreateRequestWithJSONBody(&suite.Suite, body)
	req.AddCookie(suite.SessionCookie)

	suite.UserCRUDMock.On("GetUserBySessionID", mock.Anything).Return(&models.User{}, nil)
	suite.PasswordHasherMock.On("ComparePasswords", mock.Anything, mock.Anything).Return(nil)
	suite.PasswordCriteriaValidatorMock.On("ValidatePasswordCriteria", mock.Anything).Return(helpers.CreateValidatePasswordCriteriaValid())
	suite.PasswordHasherMock.On("HashPassword", mock.Anything).Return(nil, nil)
	suite.UserCRUDMock.On("UpdateUser", mock.Anything).Return(errors.New(""))

	//act
	suite.UserController.PatchAPIUserPassword(w, req, nil)

	//assert
	common.AssertInternalServerErrorResponse(&suite.Suite, w.Result())
}

func (suite *UserControllerTestSuite) TestPatchAPIUserPassword_WithValidRequest_ReturnsOK() {
	//arrange
	w := httptest.NewRecorder()

	body := api.PatchAPIUserPasswordBody{
		OldPassword: "old password",
		NewPassword: "new password",
	}

	req := common.CreateRequestWithJSONBody(&suite.Suite, body)
	req.AddCookie(suite.SessionCookie)

	oldPasswordHash := []byte("hashed old password")
	newPasswordHash := []byte("hashed new password")
	user := models.CreateNewUser("username", oldPasswordHash)

	suite.UserCRUDMock.On("GetUserBySessionID", mock.Anything).Return(user, nil)
	suite.PasswordHasherMock.On("ComparePasswords", mock.Anything, mock.Anything).Return(nil)
	suite.PasswordCriteriaValidatorMock.On("ValidatePasswordCriteria", mock.Anything).Return(helpers.CreateValidatePasswordCriteriaValid())
	suite.PasswordHasherMock.On("HashPassword", mock.Anything).Return(newPasswordHash, nil)
	suite.UserCRUDMock.On("UpdateUser", mock.Anything).Return(nil)

	//act
	suite.UserController.PatchAPIUserPassword(w, req, nil)

	//assert
	suite.UserCRUDMock.AssertCalled(suite.T(), "GetUserBySessionID", suite.SID)
	suite.PasswordHasherMock.AssertCalled(suite.T(), "ComparePasswords", oldPasswordHash, body.OldPassword)
	suite.PasswordCriteriaValidatorMock.AssertCalled(suite.T(), "ValidatePasswordCriteria", body.NewPassword)
	suite.PasswordHasherMock.AssertCalled(suite.T(), "HashPassword", body.NewPassword)
	suite.UserCRUDMock.AssertCalled(suite.T(), "UpdateUser", mock.MatchedBy(func(u *models.User) bool {
		return bytes.Equal(u.PasswordHash, newPasswordHash)
	}))

	common.AssertSuccessResponse(&suite.Suite, w.Result())
}

func TestUserControllerTestSuite(t *testing.T) {
	suite.Run(t, &UserControllerTestSuite{})
}
