package controllers_test

import (
	"authserver/controllers"
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/stretchr/testify/suite"
)

func createEmptyRequest(suite *suite.Suite) *http.Request {
	req, err := http.NewRequest("", "", nil)
	suite.Require().NoError(err)

	return req
}

func createRequestWithJSONBody(suite *suite.Suite, body interface{}) *http.Request {
	bodyStr, err := json.Marshal(body)
	suite.Require().NoError(err)

	req, err := http.NewRequest("", "", bytes.NewReader(bodyStr))
	suite.Require().NoError(err)

	return req
}

func createRequestWithAuthorizationHeader(suite *suite.Suite, token string) *http.Request {
	req, err := http.NewRequest("", "", nil)
	suite.Require().NoError(err)

	req.Header.Set("Authorization", "Bearer "+token)

	return req
}

func parseResponse(suite *suite.Suite, res *http.Response, body interface{}) (status int) {
	decoder := json.NewDecoder(res.Body)
	err := decoder.Decode(body)
	suite.Require().NoError(err)

	return res.StatusCode
}

func assertSuccessResponse(suite *suite.Suite, res *http.Response) {
	var basicRes controllers.BasicResponse
	status := parseResponse(suite, res, &basicRes)

	suite.Equal(http.StatusOK, status)
	suite.True(basicRes.Success)
}

func assertErrorResponse(suite *suite.Suite, res *http.Response, expectedStatus int, expectedErrorSubStrings ...string) {
	var errRes controllers.ErrorResponse
	status := parseResponse(suite, res, &errRes)

	suite.Equal(expectedStatus, status)
	suite.False(errRes.Success)

	for _, expectedError := range expectedErrorSubStrings {
		suite.Contains(errRes.Error, expectedError)
	}
}

func assertInternalServerErrorResponse(suite *suite.Suite, res *http.Response) {
	assertErrorResponse(suite, res, http.StatusInternalServerError, "an internal error occurred")
}

func assertOAuthErrorResponse(suite *suite.Suite, res *http.Response, expectedStatus int, expectedError string, expectedDescription string) {
	var errRes controllers.OAuthErrorResponse
	status := parseResponse(suite, res, &errRes)

	suite.Equal(expectedStatus, status)
	suite.Contains(errRes.Error, expectedError)
	suite.Contains(errRes.ErrorDescription, expectedDescription)
}

func assertAccessTokenResponse(suite *suite.Suite, res *http.Response, expectedTokenID string) {
	var tokenRes controllers.AccessTokenResponse
	status := parseResponse(suite, res, &tokenRes)

	suite.Equal(http.StatusOK, status)
	suite.Equal(expectedTokenID, tokenRes.AccessToken)
	suite.Equal("bearer", tokenRes.TokenType)
}
