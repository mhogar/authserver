// Code generated by mockery v1.1.2. DO NOT EDIT.

package mocks

import (
	models "authserver/models"
	http "net/http"

	mock "github.com/stretchr/testify/mock"

	requesterror "authserver/common/request_error"
)

// Authenticator is an autogenerated mock type for the Authenticator type
type Authenticator struct {
	mock.Mock
}

// Authenticate provides a mock function with given fields: req
func (_m *Authenticator) Authenticate(req *http.Request) (*models.AccessToken, requesterror.RequestError) {
	ret := _m.Called(req)

	var r0 *models.AccessToken
	if rf, ok := ret.Get(0).(func(*http.Request) *models.AccessToken); ok {
		r0 = rf(req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.AccessToken)
		}
	}

	var r1 requesterror.RequestError
	if rf, ok := ret.Get(1).(func(*http.Request) requesterror.RequestError); ok {
		r1 = rf(req)
	} else {
		r1 = ret.Get(1).(requesterror.RequestError)
	}

	return r0, r1
}
