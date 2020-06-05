package common

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/stretchr/testify/suite"
)

// CreateRequestWithJSONBody is a testing helper method to create a new http request with the provided json body.
func CreateRequestWithJSONBody(suite *suite.Suite, body interface{}) *http.Request {
	bodyStr, err := json.Marshal(body)
	suite.Require().NoError(err)

	req, err := http.NewRequest("", "", bytes.NewReader(bodyStr))
	suite.Require().NoError(err)

	return req
}

// ParseResponse is a testing helper method to parse the http response into body.
func ParseResponse(suite *suite.Suite, res *http.Response, body interface{}) (status int) {
	decoder := json.NewDecoder(res.Body)
	err := decoder.Decode(body)
	suite.Require().NoError(err)

	return res.StatusCode
}

// AssertSuccessResponse is a testing helper method to assert that the response is a success reponse.
func AssertSuccessResponse(suite *suite.Suite, res *http.Response) {
	var basicRes BasicResponse
	status := ParseResponse(suite, res, &basicRes)

	suite.Equal(http.StatusOK, status)
	suite.True(basicRes.Success)
}

// AssertErrorResponse is a testing helper method to assert that the response is an error reponse with the provided status and error message.
func AssertErrorResponse(suite *suite.Suite, res *http.Response, expectedStatus int, expectedError string) {
	var errRes ErrorResponse
	status := ParseResponse(suite, res, &errRes)

	suite.Equal(expectedStatus, status)
	suite.False(errRes.Success)
	suite.Contains(errRes.Error, expectedError)
}

// AssertInternalServerErrorResponse is a testing helper method to assert that the response is an internal server error.
func AssertInternalServerErrorResponse(suite *suite.Suite, res *http.Response) {
	AssertErrorResponse(suite, res, http.StatusInternalServerError, "an internal error occurred")
}
