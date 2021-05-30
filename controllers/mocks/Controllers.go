// Code generated by mockery v1.1.2. DO NOT EDIT.

package mocks

import (
	models "authserver/models"

	mock "github.com/stretchr/testify/mock"

	requesterror "authserver/common/request_error"

	uuid "github.com/google/uuid"
)

// Controllers is an autogenerated mock type for the Controllers type
type Controllers struct {
	mock.Mock
}

// CreateTokenFromPassword provides a mock function with given fields: username, password, clientID, scopeName
func (_m *Controllers) CreateTokenFromPassword(username string, password string, clientID uuid.UUID, scopeName string) (*models.AccessToken, requesterror.OAuthRequestError) {
	ret := _m.Called(username, password, clientID, scopeName)

	var r0 *models.AccessToken
	if rf, ok := ret.Get(0).(func(string, string, uuid.UUID, string) *models.AccessToken); ok {
		r0 = rf(username, password, clientID, scopeName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.AccessToken)
		}
	}

	var r1 requesterror.OAuthRequestError
	if rf, ok := ret.Get(1).(func(string, string, uuid.UUID, string) requesterror.OAuthRequestError); ok {
		r1 = rf(username, password, clientID, scopeName)
	} else {
		r1 = ret.Get(1).(requesterror.OAuthRequestError)
	}

	return r0, r1
}

// CreateUser provides a mock function with given fields: username, password
func (_m *Controllers) CreateUser(username string, password string) (*models.User, requesterror.RequestError) {
	ret := _m.Called(username, password)

	var r0 *models.User
	if rf, ok := ret.Get(0).(func(string, string) *models.User); ok {
		r0 = rf(username, password)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.User)
		}
	}

	var r1 requesterror.RequestError
	if rf, ok := ret.Get(1).(func(string, string) requesterror.RequestError); ok {
		r1 = rf(username, password)
	} else {
		r1 = ret.Get(1).(requesterror.RequestError)
	}

	return r0, r1
}

// DeleteToken provides a mock function with given fields: token
func (_m *Controllers) DeleteToken(token *models.AccessToken) requesterror.RequestError {
	ret := _m.Called(token)

	var r0 requesterror.RequestError
	if rf, ok := ret.Get(0).(func(*models.AccessToken) requesterror.RequestError); ok {
		r0 = rf(token)
	} else {
		r0 = ret.Get(0).(requesterror.RequestError)
	}

	return r0
}

// DeleteUser provides a mock function with given fields: id
func (_m *Controllers) DeleteUser(id uuid.UUID) requesterror.RequestError {
	ret := _m.Called(id)

	var r0 requesterror.RequestError
	if rf, ok := ret.Get(0).(func(uuid.UUID) requesterror.RequestError); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(requesterror.RequestError)
	}

	return r0
}

// UpdateUserPassword provides a mock function with given fields: user, oldPassword, newPassword
func (_m *Controllers) UpdateUserPassword(user *models.User, oldPassword string, newPassword string) requesterror.RequestError {
	ret := _m.Called(user, oldPassword, newPassword)

	var r0 requesterror.RequestError
	if rf, ok := ret.Get(0).(func(*models.User, string, string) requesterror.RequestError); ok {
		r0 = rf(user, oldPassword, newPassword)
	} else {
		r0 = ret.Get(0).(requesterror.RequestError)
	}

	return r0
}
