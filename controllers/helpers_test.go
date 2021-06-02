package controllers_test

import (
	"authserver/common"
	requesterror "authserver/common/request_error"

	"github.com/stretchr/testify/suite"
)

func AssertNoError(suite *suite.Suite, err requesterror.RequestError) {
	suite.Require().NotNil(err)
	suite.Equal(requesterror.ErrorTypeNone, err.Type)
}

func AssertClientError(suite *suite.Suite, err requesterror.RequestError, expectedSubStrs ...string) {
	suite.Require().NotNil(err)
	suite.Equal(requesterror.ErrorTypeClient, err.Type)
	common.AssertContainsSubstrings(suite, err.Error(), expectedSubStrs...)
}

func AssertInternalError(suite *suite.Suite, err requesterror.RequestError) {
	suite.Require().NotNil(err)
	suite.Equal(requesterror.ErrorTypeInternal, err.Type)
	common.AssertContainsSubstrings(suite, err.Error(), "internal error")
}

func AssertOAuthNoError(suite *suite.Suite, err requesterror.OAuthRequestError) {
	AssertNoError(suite, err.RequestError)
}

func AssertOAuthClientError(suite *suite.Suite, err requesterror.OAuthRequestError, expectedErrorName string, expectedMessageSubStrs ...string) {
	AssertClientError(suite, err.RequestError, expectedMessageSubStrs...)
	suite.Equal(expectedErrorName, err.ErrorName)
}

func AssertOAuthInternalError(suite *suite.Suite, err requesterror.OAuthRequestError) {
	AssertInternalError(suite, err.RequestError)
}
