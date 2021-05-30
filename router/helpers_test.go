package router_test

import (
	commonhelpers "authserver/helpers/common"
	"authserver/router"
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/stretchr/testify/suite"
)

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

func ParseResponse(suite *suite.Suite, res *http.Response, body interface{}) (status int) {
	decoder := json.NewDecoder(res.Body)
	err := decoder.Decode(body)
	suite.Require().NoError(err)

	return res.StatusCode
}

func AssertSuccessResponse(suite *suite.Suite, res *http.Response) {
	var basicRes router.BasicResponse
	status := ParseResponse(suite, res, &basicRes)

	suite.Equal(http.StatusOK, status)
	suite.True(basicRes.Success)
}

func AssertErrorResponse(suite *suite.Suite, res *http.Response, expectedStatus int, expectedErrorSubStrings ...string) {
	var errRes router.ErrorResponse
	status := ParseResponse(suite, res, &errRes)

	suite.Equal(expectedStatus, status)
	suite.False(errRes.Success)

	commonhelpers.AssertContainsSubstrings(suite, errRes.Error, expectedErrorSubStrings...)
}

func AssertInternalServerErrorResponse(suite *suite.Suite, res *http.Response) {
	AssertErrorResponse(suite, res, http.StatusInternalServerError, "internal error")
}

func AssertOAuthErrorResponse(suite *suite.Suite, res *http.Response, expectedStatus int, expectedError string, expectedDescriptionSubStrings ...string) {
	var errRes router.OAuthErrorResponse
	status := ParseResponse(suite, res, &errRes)

	suite.Equal(expectedStatus, status)
	suite.Contains(errRes.Error, expectedError)

	commonhelpers.AssertContainsSubstrings(suite, errRes.ErrorDescription, expectedDescriptionSubStrings...)
}

func AssertAccessTokenResponse(suite *suite.Suite, res *http.Response, expectedTokenID string) {
	var tokenRes router.AccessTokenResponse
	status := ParseResponse(suite, res, &tokenRes)

	suite.Equal(http.StatusOK, status)
	suite.Equal(expectedTokenID, tokenRes.AccessToken)
	suite.Equal("bearer", tokenRes.TokenType)
}
