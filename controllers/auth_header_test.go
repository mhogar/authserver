package controllers_test

import (
	modelmocks "authserver/models/mocks"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func RunAuthHeaderTests(suite *suite.Suite, accessTokenCRUDMock *modelmocks.AccessTokenCRUD, actFunc httprouter.Handle) {
	suite.Run("NoBearerToken_ReturnsUnauthorized", func() {
		var req *http.Request

		testCase := func() {
			//arrange
			w := httptest.NewRecorder()

			//act
			actFunc(w, req, nil)

			//assert
			assertErrorResponse(suite, w.Result(), http.StatusUnauthorized, "no bearer token")
		}

		req = createEmptyRequest(suite)
		suite.Run("NoAuthorizationHeader", testCase)

		req.Header.Set("Authorization", "invalid")
		suite.Run("AuthorizationHeaderDoesNotContainBearerToken", testCase)
	})

	suite.Run("BearerTokenInInvalidFormat_ReturnsUnauthorized", func() {
		//arrange
		w := httptest.NewRecorder()
		req := createRequestWithAuthorizationHeader(suite, "invalid")

		//act
		actFunc(w, req, nil)

		//assert
		assertErrorResponse(suite, w.Result(), http.StatusUnauthorized, "bearer token", "invalid format")
	})

	suite.Run("ErrorFetchingAccessTokenByID_ReturnsInternalServerError", func() {
		//arrange
		w := httptest.NewRecorder()
		req := createRequestWithAuthorizationHeader(suite, uuid.New().String())

		accessTokenCRUDMock.On("GetAccessTokenByID", mock.Anything).Return(nil, errors.New(""))

		//act
		actFunc(w, req, nil)

		//assert
		assertInternalServerErrorResponse(suite, w.Result())
	})
}
