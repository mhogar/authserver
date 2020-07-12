package controllers_test

import (
	databasemocks "authserver/database/mocks"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func RunAuthHeaderTests(suite *suite.Suite, CRUDMock *databasemocks.CRUDOperations, setupTest func(), actFunc httprouter.Handle) {
	setupTest()
	suite.Run("NoBearerToken_ReturnsUnauthorized", func() {
		var req *http.Request

		testCase := func() {
			//arrange
			w := httptest.NewRecorder()

			//act
			actFunc(w, req, nil)

			//assert
			AssertErrorResponse(suite, w.Result(), http.StatusUnauthorized, "no bearer token")
		}

		req = CreateRequest(suite, "", nil)
		suite.Run("NoAuthorizationHeader", testCase)

		req.Header.Set("Authorization", "invalid")
		suite.Run("AuthorizationHeaderDoesNotContainBearerToken", testCase)
	})

	setupTest()
	suite.Run("BearerTokenInInvalidFormat_ReturnsUnauthorized", func() {
		//arrange
		w := httptest.NewRecorder()
		req := CreateRequest(suite, "invalid", nil)

		//act
		actFunc(w, req, nil)

		//assert
		AssertErrorResponse(suite, w.Result(), http.StatusUnauthorized, "bearer token", "invalid format")
	})

	setupTest()
	suite.Run("ErrorFetchingAccessTokenByID_ReturnsInternalServerError", func() {
		//arrange
		w := httptest.NewRecorder()
		req := CreateRequest(suite, uuid.New().String(), nil)

		CRUDMock.On("GetAccessTokenByID", mock.Anything).Return(nil, errors.New(""))

		//act
		actFunc(w, req, nil)

		//assert
		AssertInternalServerErrorResponse(suite, w.Result())
	})

	setupTest()
	suite.Run("AccessTokenWithIDisNotFound_ReturnsUnauthorized", func() {
		//arrange
		w := httptest.NewRecorder()
		req := CreateRequest(suite, uuid.New().String(), nil)

		CRUDMock.On("GetAccessTokenByID", mock.Anything).Return(nil, nil)

		//act
		actFunc(w, req, nil)

		//assert
		AssertErrorResponse(suite, w.Result(), http.StatusUnauthorized, "bearer token", "invalid")
	})
}
