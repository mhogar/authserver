package controllers_test

import (
	"authserver/controllers"
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/stretchr/testify/suite"
)

func createRequestWithJSONBody(suite *suite.Suite, body interface{}) *http.Request {
	bodyStr, err := json.Marshal(body)
	suite.Require().NoError(err)

	req, err := http.NewRequest("", "", bytes.NewReader(bodyStr))
	suite.Require().NoError(err)

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

func assertErrorResponse(suite *suite.Suite, res *http.Response, expectedStatus int, expectedError string) {
	var errRes controllers.ErrorResponse
	status := parseResponse(suite, res, &errRes)

	suite.Equal(expectedStatus, status)
	suite.False(errRes.Success)
	suite.Contains(errRes.Error, expectedError)
}

func assertInternalServerErrorResponse(suite *suite.Suite, res *http.Response) {
	assertErrorResponse(suite, res, http.StatusInternalServerError, "an internal error occurred")
}
