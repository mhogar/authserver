package controllers_test

import (
	"authserver/controllers"
	databasemocks "authserver/database/mocks"
	"authserver/helpers"
	helpermocks "authserver/helpers/mocks"
	"authserver/models"
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type UserControllerTestSuite struct {
	suite.Suite
	CRUDMock                      databasemocks.CRUDOperations
	PasswordHasherMock            helpermocks.PasswordHasher
	PasswordCriteriaValidatorMock helpermocks.PasswordCriteriaValidator
	UserController                controllers.UserController
}

func (suite *UserControllerTestSuite) SetupTest() {
	suite.CRUDMock = databasemocks.CRUDOperations{}
	suite.PasswordHasherMock = helpermocks.PasswordHasher{}
	suite.PasswordCriteriaValidatorMock = helpermocks.PasswordCriteriaValidator{}
	suite.UserController = controllers.UserController{
		CRUD:                      &suite.CRUDMock,
		PasswordHasher:            &suite.PasswordHasherMock,
		PasswordCriteriaValidator: &suite.PasswordCriteriaValidatorMock,
	}
}

func (suite *UserControllerTestSuite) TestPostUser_AuthorizationHeaderTests() {
	setupTest := func() {
		suite.CRUDMock = databasemocks.CRUDOperations{}
		suite.UserController.CRUD = &suite.CRUDMock
	}

	RunAuthHeaderTests(&suite.Suite, &suite.CRUDMock, setupTest, suite.UserController.PostUser)
}

func (suite *UserControllerTestSuite) TestPostUser_WithInvalidJSONBody_ReturnsBadRequest() {
	//arrange
	w := httptest.NewRecorder()
	req := CreateRequest(&suite.Suite, uuid.New().String(), "invalid")

	suite.CRUDMock.On("GetAccessTokenByID", mock.Anything).Return(&models.AccessToken{}, nil)

	//act
	suite.UserController.PostUser(w, req, nil)

	//assert
	AssertErrorResponse(&suite.Suite, w.Result(), http.StatusBadRequest, "invalid json body")
}

func (suite *UserControllerTestSuite) TestPostUser_WithInvalidBodyFields_ReturnsBadRequest() {
	var body controllers.PostUserBody

	testCase := func() {
		//arrange
		w := httptest.NewRecorder()
		req := CreateRequest(&suite.Suite, uuid.New().String(), body)

		suite.CRUDMock.On("GetAccessTokenByID", mock.Anything).Return(&models.AccessToken{}, nil)

		//act
		suite.UserController.PostUser(w, req, nil)

		//assert
		AssertErrorResponse(&suite.Suite, w.Result(), http.StatusBadRequest, "username", "password", "cannot be empty")
	}

	body = controllers.PostUserBody{
		Username: "",
		Password: "password",
	}
	suite.Run("EmptyUsername", testCase)

	body = controllers.PostUserBody{
		Username: "username",
		Password: "",
	}
	suite.Run("EmptyPassword", testCase)
}

func (suite *UserControllerTestSuite) TestPostUser_WithErrorGettingUserByUsername_ReturnsInternalServerError() {
	//arrange
	w := httptest.NewRecorder()

	body := controllers.PostUserBody{
		Username: "username",
		Password: "password",
	}
	req := CreateRequest(&suite.Suite, uuid.New().String(), body)

	suite.CRUDMock.On("GetAccessTokenByID", mock.Anything).Return(&models.AccessToken{}, nil)
	suite.CRUDMock.On("GetUserByUsername", body.Username).Return(nil, errors.New(""))

	//act
	suite.UserController.PostUser(w, req, nil)

	//assert
	AssertInternalServerErrorResponse(&suite.Suite, w.Result())
}

func (suite *UserControllerTestSuite) TestPostUser_WithNonUniqueUsername_ReturnsBadRequest() {
	//arrange
	w := httptest.NewRecorder()

	body := controllers.PostUserBody{
		Username: "username",
		Password: "password",
	}
	req := CreateRequest(&suite.Suite, uuid.New().String(), body)

	suite.CRUDMock.On("GetAccessTokenByID", mock.Anything).Return(&models.AccessToken{}, nil)
	suite.CRUDMock.On("GetUserByUsername", body.Username).Return(&models.User{}, nil)

	//act
	suite.UserController.PostUser(w, req, nil)

	//assert
	AssertErrorResponse(&suite.Suite, w.Result(), http.StatusBadRequest, "username already exists")
}

func (suite *UserControllerTestSuite) TestPostUser_WherePasswordDoesNotMeetCriteria_ReturnsBadRequest() {
	//arrange
	w := httptest.NewRecorder()

	body := controllers.PostUserBody{
		Username: "username",
		Password: "password",
	}
	req := CreateRequest(&suite.Suite, uuid.New().String(), body)

	suite.CRUDMock.On("GetAccessTokenByID", mock.Anything).Return(&models.AccessToken{}, nil)
	suite.CRUDMock.On("GetUserByUsername", body.Username).Return(nil, nil)
	suite.PasswordCriteriaValidatorMock.On("ValidatePasswordCriteria", mock.Anything).Return(helpers.CreateValidatePasswordCriteriaError(helpers.ValidatePasswordCriteriaTooShort, ""))

	//act
	suite.UserController.PostUser(w, req, nil)

	//assert
	AssertErrorResponse(&suite.Suite, w.Result(), http.StatusBadRequest, "password", "not", "minimum criteria")
}

func (suite *UserControllerTestSuite) TestPostUser_WithErrorHashingNewPassword_ReturnsInternalServerError() {
	//arrange
	w := httptest.NewRecorder()

	body := controllers.PostUserBody{
		Username: "username",
		Password: "password",
	}
	req := CreateRequest(&suite.Suite, uuid.New().String(), body)

	suite.CRUDMock.On("GetAccessTokenByID", mock.Anything).Return(&models.AccessToken{}, nil)
	suite.CRUDMock.On("GetUserByUsername", body.Username).Return(nil, nil)
	suite.PasswordCriteriaValidatorMock.On("ValidatePasswordCriteria", mock.Anything).Return(helpers.CreateValidatePasswordCriteriaValid())
	suite.PasswordHasherMock.On("HashPassword", mock.Anything).Return(nil, errors.New(""))

	//act
	suite.UserController.PostUser(w, req, nil)

	//assert
	AssertInternalServerErrorResponse(&suite.Suite, w.Result())
}

func (suite *UserControllerTestSuite) TestPostUser_WithErrorCreatingUser_ReturnsInternalServerError() {
	//arrange
	w := httptest.NewRecorder()

	body := controllers.PostUserBody{
		Username: "username",
		Password: "password",
	}
	req := CreateRequest(&suite.Suite, uuid.New().String(), body)

	suite.CRUDMock.On("GetAccessTokenByID", mock.Anything).Return(&models.AccessToken{}, nil)
	suite.CRUDMock.On("GetUserByUsername", body.Username).Return(nil, nil)
	suite.PasswordHasherMock.On("HashPassword", mock.Anything).Return(nil, nil)
	suite.PasswordCriteriaValidatorMock.On("ValidatePasswordCriteria", mock.Anything).Return(helpers.CreateValidatePasswordCriteriaValid())
	suite.CRUDMock.On("SaveUser", mock.Anything).Return(errors.New(""))

	//act
	suite.UserController.PostUser(w, req, nil)

	//assert
	AssertInternalServerErrorResponse(&suite.Suite, w.Result())
}

func (suite *UserControllerTestSuite) TestPostUser_WithValidRequest_ReturnsOK() {
	//arrange
	w := httptest.NewRecorder()

	body := controllers.PostUserBody{
		Username: "username",
		Password: "password",
	}
	tokenID := uuid.New()

	req := CreateRequest(&suite.Suite, tokenID.String(), body)

	hash := []byte("password hash")

	suite.CRUDMock.On("GetAccessTokenByID", mock.Anything).Return(&models.AccessToken{}, nil)
	suite.CRUDMock.On("GetUserByUsername", body.Username).Return(nil, nil)
	suite.PasswordCriteriaValidatorMock.On("ValidatePasswordCriteria", mock.Anything).Return(helpers.CreateValidatePasswordCriteriaValid())
	suite.PasswordHasherMock.On("HashPassword", mock.Anything).Return(hash, nil)
	suite.CRUDMock.On("SaveUser", mock.Anything).Return(nil)

	//act
	suite.UserController.PostUser(w, req, nil)

	//assert
	suite.CRUDMock.AssertCalled(suite.T(), "GetAccessTokenByID", tokenID)
	suite.CRUDMock.AssertCalled(suite.T(), "GetUserByUsername", body.Username)
	suite.PasswordCriteriaValidatorMock.AssertCalled(suite.T(), "ValidatePasswordCriteria", body.Password)
	suite.PasswordHasherMock.AssertCalled(suite.T(), "HashPassword", body.Password)
	suite.CRUDMock.AssertCalled(suite.T(), "SaveUser", mock.MatchedBy(func(u *models.User) bool {
		return u.Username == body.Username && bytes.Equal(u.PasswordHash, hash)
	}))

	AssertSuccessResponse(&suite.Suite, w.Result())
}

func (suite *UserControllerTestSuite) TestDeleteUser_AuthorizationHeaderTests() {
	setupTest := func() {
		suite.CRUDMock = databasemocks.CRUDOperations{}
		suite.UserController.CRUD = &suite.CRUDMock
	}

	RunAuthHeaderTests(&suite.Suite, &suite.CRUDMock, setupTest, suite.UserController.DeleteUser)
}

func (suite *UserControllerTestSuite) TestDeleteUser_WithoutIdInParams_ReturnsBadRequest() {
	//arrange
	w := httptest.NewRecorder()
	req := CreateRequest(&suite.Suite, uuid.New().String(), nil)

	suite.CRUDMock.On("GetAccessTokenByID", mock.Anything).Return(&models.AccessToken{}, nil)

	//act
	suite.UserController.DeleteUser(w, req, make(httprouter.Params, 0))

	//assert
	AssertErrorResponse(&suite.Suite, w.Result(), http.StatusBadRequest, "id must be present")
}

func (suite *UserControllerTestSuite) TestDeleteUser_WithIdInInvalidFormat_ReturnsBadRequest() {
	//arrange
	w := httptest.NewRecorder()
	req := CreateRequest(&suite.Suite, uuid.New().String(), nil)

	id := 0
	params := httprouter.Params{
		httprouter.Param{Key: "id", Value: string(id)},
	}

	suite.CRUDMock.On("GetAccessTokenByID", mock.Anything).Return(&models.AccessToken{}, nil)

	//act
	suite.UserController.DeleteUser(w, req, params)

	//assert
	AssertErrorResponse(&suite.Suite, w.Result(), http.StatusBadRequest, "id", "invalid format")
}

func (suite *UserControllerTestSuite) TestDeleteUser_WithErrorGettingUserById_ReturnsInternalServerError() {
	//arrange
	w := httptest.NewRecorder()
	req := CreateRequest(&suite.Suite, uuid.New().String(), nil)

	id := uuid.New()
	params := httprouter.Params{
		httprouter.Param{Key: "id", Value: id.String()},
	}

	suite.CRUDMock.On("GetAccessTokenByID", mock.Anything).Return(&models.AccessToken{}, nil)
	suite.CRUDMock.On("GetUserByID", mock.Anything).Return(nil, errors.New(""))

	//act
	suite.UserController.DeleteUser(w, req, params)

	//assert
	AssertInternalServerErrorResponse(&suite.Suite, w.Result())
}

func (suite *UserControllerTestSuite) TestDeleteUser_WhereUserIsNotFound_ReturnsBadRequest() {
	//arrange
	w := httptest.NewRecorder()
	req := CreateRequest(&suite.Suite, uuid.New().String(), nil)

	id := uuid.New()
	params := httprouter.Params{
		httprouter.Param{Key: "id", Value: id.String()},
	}

	suite.CRUDMock.On("GetAccessTokenByID", mock.Anything).Return(&models.AccessToken{}, nil)
	suite.CRUDMock.On("GetUserByID", mock.Anything).Return(nil, nil)

	//act
	suite.UserController.DeleteUser(w, req, params)

	//assert
	AssertErrorResponse(&suite.Suite, w.Result(), http.StatusBadRequest, "user not found")
}

func (suite *UserControllerTestSuite) TestDeleteUser_WithErrorDeletingUser_ReturnsInternalServerError() {
	//arrange
	w := httptest.NewRecorder()
	req := CreateRequest(&suite.Suite, uuid.New().String(), nil)

	user := models.CreateNewUser("username", []byte("password hash"))
	params := httprouter.Params{
		httprouter.Param{Key: "id", Value: user.ID.String()},
	}

	suite.CRUDMock.On("GetAccessTokenByID", mock.Anything).Return(&models.AccessToken{}, nil)
	suite.CRUDMock.On("GetUserByID", mock.Anything).Return(user, nil)
	suite.CRUDMock.On("DeleteUser", mock.Anything).Return(errors.New(""))

	//act
	suite.UserController.DeleteUser(w, req, params)

	//assert
	AssertInternalServerErrorResponse(&suite.Suite, w.Result())
}

func (suite *UserControllerTestSuite) TestDeleteUser_WithValidRequest_ReturnsOK() {
	//arrange
	w := httptest.NewRecorder()

	tokenID := uuid.New()
	req := CreateRequest(&suite.Suite, tokenID.String(), nil)

	user := models.CreateNewUser("username", []byte("password hash"))
	params := httprouter.Params{
		httprouter.Param{Key: "id", Value: user.ID.String()},
	}

	suite.CRUDMock.On("GetAccessTokenByID", mock.Anything).Return(&models.AccessToken{}, nil)
	suite.CRUDMock.On("GetUserByID", mock.Anything).Return(user, nil)
	suite.CRUDMock.On("DeleteUser", mock.Anything).Return(nil)

	//act
	suite.UserController.DeleteUser(w, req, params)

	//assert
	suite.CRUDMock.AssertCalled(suite.T(), "GetAccessTokenByID", tokenID)
	suite.CRUDMock.AssertCalled(suite.T(), "GetUserByID", user.ID)
	suite.CRUDMock.AssertCalled(suite.T(), "DeleteUser", user)

	AssertSuccessResponse(&suite.Suite, w.Result())
}

func (suite *UserControllerTestSuite) TestPatchUserPassword_AuthorizationHeaderTests() {
	setupTest := func() {
		suite.CRUDMock = databasemocks.CRUDOperations{}
		suite.UserController.CRUD = &suite.CRUDMock
	}

	RunAuthHeaderTests(&suite.Suite, &suite.CRUDMock, setupTest, suite.UserController.PatchUserPassword)
}

func (suite *UserControllerTestSuite) TestPatchUserPassword_WithErrorGettingUserById_ReturnsInternalServerError() {
	//arrange
	w := httptest.NewRecorder()
	req := CreateRequest(&suite.Suite, uuid.New().String(), nil)

	suite.CRUDMock.On("GetAccessTokenByID", mock.Anything).Return(&models.AccessToken{}, nil)
	suite.CRUDMock.On("GetUserByID", mock.Anything).Return(nil, errors.New(""))

	//act
	suite.UserController.PatchUserPassword(w, req, nil)

	//assert
	AssertInternalServerErrorResponse(&suite.Suite, w.Result())
}

func (suite *UserControllerTestSuite) TestPatchUserPassword_WhereNoUserIsFound_ReturnsUnauthorized() {
	//arrange
	w := httptest.NewRecorder()
	req := CreateRequest(&suite.Suite, uuid.New().String(), nil)

	suite.CRUDMock.On("GetAccessTokenByID", mock.Anything).Return(&models.AccessToken{}, nil)
	suite.CRUDMock.On("GetUserByID", mock.Anything).Return(nil, nil)

	//act
	suite.UserController.PatchUserPassword(w, req, nil)

	//assert
	AssertErrorResponse(&suite.Suite, w.Result(), http.StatusUnauthorized, "no user", "access token")
}

func (suite *UserControllerTestSuite) TestPatchUserPassword_WithInvalidJSONBody_ReturnsBadRequest() {
	//arrange
	w := httptest.NewRecorder()
	req := CreateRequest(&suite.Suite, uuid.New().String(), "invalid")

	suite.CRUDMock.On("GetAccessTokenByID", mock.Anything).Return(&models.AccessToken{}, nil)
	suite.CRUDMock.On("GetUserByID", mock.Anything).Return(&models.User{}, nil)

	//act
	suite.UserController.PatchUserPassword(w, req, nil)

	//assert
	AssertErrorResponse(&suite.Suite, w.Result(), http.StatusBadRequest, "invalid json body")
}

func (suite *UserControllerTestSuite) TestPatchUserPassword_WithInvalidBodyFields_ReturnsBadRequest() {
	var body controllers.PatchUserPasswordBody

	testCase := func() {
		//arrange
		w := httptest.NewRecorder()
		req := CreateRequest(&suite.Suite, uuid.New().String(), body)

		suite.CRUDMock.On("GetAccessTokenByID", mock.Anything).Return(&models.AccessToken{}, nil)
		suite.CRUDMock.On("GetUserByID", mock.Anything).Return(&models.User{}, nil)

		//act
		suite.UserController.PatchUserPassword(w, req, nil)

		//assert
		AssertErrorResponse(&suite.Suite, w.Result(), http.StatusBadRequest, "old password", "new password", "cannot be empty")
	}

	body = controllers.PatchUserPasswordBody{
		OldPassword: "",
		NewPassword: "new password",
	}
	suite.Run("EmptyUsername", testCase)

	body = controllers.PatchUserPasswordBody{
		OldPassword: "old password",
		NewPassword: "",
	}
	suite.Run("EmptyPassword", testCase)
}

func (suite *UserControllerTestSuite) TestPatchUserPassword_WhereOldPasswordIsInvalid_ReturnsBadRequest() {
	//arrange
	w := httptest.NewRecorder()

	body := controllers.PatchUserPasswordBody{
		OldPassword: "old password",
		NewPassword: "new password",
	}
	req := CreateRequest(&suite.Suite, uuid.New().String(), body)

	suite.CRUDMock.On("GetAccessTokenByID", mock.Anything).Return(&models.AccessToken{}, nil)
	suite.CRUDMock.On("GetUserByID", mock.Anything).Return(&models.User{}, nil)
	suite.PasswordHasherMock.On("ComparePasswords", mock.Anything, mock.Anything).Return(errors.New(""))

	//act
	suite.UserController.PatchUserPassword(w, req, nil)

	//assert
	AssertErrorResponse(&suite.Suite, w.Result(), http.StatusBadRequest, "old password", "invalid")
}

func (suite *UserControllerTestSuite) TestPatchUserPassword_WhereNewPasswordDoesNotMeetCriteria_ReturnsBadRequest() {
	//arrange
	w := httptest.NewRecorder()

	body := controllers.PatchUserPasswordBody{
		OldPassword: "old password",
		NewPassword: "new password",
	}
	req := CreateRequest(&suite.Suite, uuid.New().String(), body)

	suite.CRUDMock.On("GetAccessTokenByID", mock.Anything).Return(&models.AccessToken{}, nil)
	suite.CRUDMock.On("GetUserByID", mock.Anything).Return(&models.User{}, nil)
	suite.PasswordHasherMock.On("ComparePasswords", mock.Anything, mock.Anything).Return(nil)
	suite.PasswordCriteriaValidatorMock.On("ValidatePasswordCriteria", mock.Anything).Return(helpers.CreateValidatePasswordCriteriaError(helpers.ValidatePasswordCriteriaTooShort, ""))

	//act
	suite.UserController.PatchUserPassword(w, req, nil)

	//assert
	AssertErrorResponse(&suite.Suite, w.Result(), http.StatusBadRequest, "password", "not", "minimum criteria")
}

func (suite *UserControllerTestSuite) TestPatchUserPassword_WithErrorHashingNewPassword_ReturnsInternalServerError() {
	//arrange
	w := httptest.NewRecorder()

	body := controllers.PatchUserPasswordBody{
		OldPassword: "old password",
		NewPassword: "new password",
	}
	req := CreateRequest(&suite.Suite, uuid.New().String(), body)

	suite.CRUDMock.On("GetAccessTokenByID", mock.Anything).Return(&models.AccessToken{}, nil)
	suite.CRUDMock.On("GetUserByID", mock.Anything).Return(&models.User{}, nil)
	suite.PasswordHasherMock.On("ComparePasswords", mock.Anything, mock.Anything).Return(nil)
	suite.PasswordCriteriaValidatorMock.On("ValidatePasswordCriteria", mock.Anything).Return(helpers.CreateValidatePasswordCriteriaValid())
	suite.PasswordHasherMock.On("HashPassword", mock.Anything).Return(nil, errors.New(""))

	//act
	suite.UserController.PatchUserPassword(w, req, nil)

	//assert
	AssertInternalServerErrorResponse(&suite.Suite, w.Result())
}

func (suite *UserControllerTestSuite) TestPatchUserPassword_WithErrorUpdatingUser_ReturnsInternalServerError() {
	//arrange
	w := httptest.NewRecorder()

	body := controllers.PatchUserPasswordBody{
		OldPassword: "old password",
		NewPassword: "new password",
	}
	req := CreateRequest(&suite.Suite, uuid.New().String(), body)

	suite.CRUDMock.On("GetAccessTokenByID", mock.Anything).Return(&models.AccessToken{}, nil)
	suite.CRUDMock.On("GetUserByID", mock.Anything).Return(&models.User{}, nil)
	suite.PasswordHasherMock.On("ComparePasswords", mock.Anything, mock.Anything).Return(nil)
	suite.PasswordCriteriaValidatorMock.On("ValidatePasswordCriteria", mock.Anything).Return(helpers.CreateValidatePasswordCriteriaValid())
	suite.PasswordHasherMock.On("HashPassword", mock.Anything).Return(nil, nil)
	suite.CRUDMock.On("UpdateUser", mock.Anything).Return(errors.New(""))

	//act
	suite.UserController.PatchUserPassword(w, req, nil)

	//assert
	AssertInternalServerErrorResponse(&suite.Suite, w.Result())
}

func (suite *UserControllerTestSuite) TestPatchUserPassword_WithValidRequest_ReturnsOK() {
	//arrange
	w := httptest.NewRecorder()

	body := controllers.PatchUserPasswordBody{
		OldPassword: "old password",
		NewPassword: "new password",
	}
	tokenID := uuid.New()

	req := CreateRequest(&suite.Suite, tokenID.String(), body)

	token := &models.AccessToken{UserID: uuid.New()}
	oldPasswordHash := []byte("hashed old password")
	newPasswordHash := []byte("hashed new password")
	user := models.CreateNewUser("username", oldPasswordHash)

	suite.CRUDMock.On("GetAccessTokenByID", mock.Anything).Return(token, nil)
	suite.CRUDMock.On("GetUserByID", mock.Anything).Return(user, nil)
	suite.PasswordHasherMock.On("ComparePasswords", mock.Anything, mock.Anything).Return(nil)
	suite.PasswordCriteriaValidatorMock.On("ValidatePasswordCriteria", mock.Anything).Return(helpers.CreateValidatePasswordCriteriaValid())
	suite.PasswordHasherMock.On("HashPassword", mock.Anything).Return(newPasswordHash, nil)
	suite.CRUDMock.On("UpdateUser", mock.Anything).Return(nil)

	//act
	suite.UserController.PatchUserPassword(w, req, nil)

	//assert
	suite.CRUDMock.AssertCalled(suite.T(), "GetAccessTokenByID", tokenID)
	suite.CRUDMock.AssertCalled(suite.T(), "GetUserByID", token.UserID)
	suite.PasswordHasherMock.AssertCalled(suite.T(), "ComparePasswords", oldPasswordHash, body.OldPassword)
	suite.PasswordCriteriaValidatorMock.AssertCalled(suite.T(), "ValidatePasswordCriteria", body.NewPassword)
	suite.PasswordHasherMock.AssertCalled(suite.T(), "HashPassword", body.NewPassword)
	suite.CRUDMock.AssertCalled(suite.T(), "UpdateUser", mock.MatchedBy(func(u *models.User) bool {
		return bytes.Equal(u.PasswordHash, newPasswordHash)
	}))

	AssertSuccessResponse(&suite.Suite, w.Result())
}

func TestUserControllerTestSuite(t *testing.T) {
	suite.Run(t, &UserControllerTestSuite{})
}
