package controllers_test

import (
	"authserver/controllers"
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/stretchr/testify/suite"
)

func CreateRequest(suite *suite.Suite, bearerToken string, body interface{}) *http.Request {
	var bodyReader io.Reader = nil

	if body != nil {
		bodyStr, err := json.Marshal(body)
		suite.Require().NoError(err)

		bodyReader = bytes.NewReader(bodyStr)
	}

	req, err := http.NewRequest("", "", bodyReader)
	suite.Require().NoError(err)

	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	return req
}

func ParseResponse(suite *suite.Suite, res *http.Response, body interface{}) (status int) {
	decoder := json.NewDecoder(res.Body)
	err := decoder.Decode(body)
	suite.Require().NoError(err)

	return res.StatusCode
}

func AssertSuccessResponse(suite *suite.Suite, res *http.Response) {
	var basicRes controllers.BasicResponse
	status := ParseResponse(suite, res, &basicRes)

	suite.Equal(http.StatusOK, status)
	suite.True(basicRes.Success)
}

func AssertErrorResponse(suite *suite.Suite, res *http.Response, expectedStatus int, expectedErrorSubStrings ...string) {
	var errRes controllers.ErrorResponse
	status := ParseResponse(suite, res, &errRes)

	suite.Equal(expectedStatus, status)
	suite.False(errRes.Success)

	for _, expectedError := range expectedErrorSubStrings {
		suite.Contains(errRes.Error, expectedError)
	}
}

func AssertInternalServerErrorResponse(suite *suite.Suite, res *http.Response) {
	AssertErrorResponse(suite, res, http.StatusInternalServerError, "an internal error occurred")
}

func AssertOAuthErrorResponse(suite *suite.Suite, res *http.Response, expectedStatus int, expectedError string, expectedDescriptionSubStrings ...string) {
	var errRes controllers.OAuthErrorResponse
	status := ParseResponse(suite, res, &errRes)

	suite.Equal(expectedStatus, status)
	suite.Contains(errRes.Error, expectedError)

	for _, expectedDescription := range expectedDescriptionSubStrings {
		suite.Contains(errRes.ErrorDescription, expectedDescription)
	}
}

func AssertAccessTokenResponse(suite *suite.Suite, res *http.Response, expectedTokenID string) {
	var tokenRes controllers.AccessTokenResponse
	status := ParseResponse(suite, res, &tokenRes)

	suite.Equal(http.StatusOK, status)
	suite.Equal(expectedTokenID, tokenRes.AccessToken)
	suite.Equal("bearer", tokenRes.TokenType)
}
