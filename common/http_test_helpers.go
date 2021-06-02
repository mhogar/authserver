package common

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/stretchr/testify/suite"
)

// CreateRequest creates an http request object with the given parameters.
func CreateRequest(suite *suite.Suite, method string, url string, bearerToken string, body interface{}) *http.Request {
	var bodyReader io.Reader = nil

	if body != nil {
		bodyStr, err := json.Marshal(body)
		suite.Require().NoError(err)

		bodyReader = bytes.NewReader(bodyStr)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	suite.Require().NoError(err)

	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	return req
}

// ParseResponse parses the provided http response, return its status code and body
func ParseResponse(suite *suite.Suite, res *http.Response, body interface{}) (status int) {
	decoder := json.NewDecoder(res.Body)
	err := decoder.Decode(body)
	suite.Require().NoError(err)

	return res.StatusCode
}

// AssertSuccessResponse asserts the response is a success response
func AssertSuccessResponse(suite *suite.Suite, res *http.Response) {
	var basicRes BasicResponse
	status := ParseResponse(suite, res, &basicRes)

	suite.Equal(http.StatusOK, status)
	suite.True(basicRes.Success)
}

// AssertErrorResponse asserts the response is an error reponse with the expected status and error sub strings
func AssertErrorResponse(suite *suite.Suite, res *http.Response, expectedStatus int, expectedErrorSubStrings ...string) {
	var errRes ErrorResponse
	status := ParseResponse(suite, res, &errRes)

	suite.Equal(expectedStatus, status)
	suite.False(errRes.Success)

	AssertContainsSubstrings(suite, errRes.Error, expectedErrorSubStrings...)
}

// AssertInternalServerErrorResponse asserts the response is an internal server response
func AssertInternalServerErrorResponse(suite *suite.Suite, res *http.Response) {
	AssertErrorResponse(suite, res, http.StatusInternalServerError, "internal error")
}

// AssertOAuthErrorResponse asserts the response is an oauth error reponse with the expected status, error, and description sub strings
func AssertOAuthErrorResponse(suite *suite.Suite, res *http.Response, expectedStatus int, expectedError string, expectedDescriptionSubStrings ...string) {
	var errRes OAuthErrorResponse
	status := ParseResponse(suite, res, &errRes)

	suite.Equal(expectedStatus, status)
	suite.Contains(errRes.Error, expectedError)

	AssertContainsSubstrings(suite, errRes.ErrorDescription, expectedDescriptionSubStrings...)
}

// AssertAccessTokenResponse asserts the response is an access token response with the expect token
func AssertAccessTokenResponse(suite *suite.Suite, res *http.Response, expectedTokenID string) {
	var tokenRes AccessTokenResponse
	status := ParseResponse(suite, res, &tokenRes)

	suite.Equal(http.StatusOK, status)
	suite.Equal(expectedTokenID, tokenRes.AccessToken)
	suite.Equal("bearer", tokenRes.TokenType)
}

// AssertResponseOK asserts the response has an http OK status and returns the parsed result
func AssertResponseOK(suite *suite.Suite, res *http.Response, result interface{}) {
	status := ParseResponse(suite, res, result)
	suite.Equal(http.StatusOK, status)
}
