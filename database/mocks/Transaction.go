// Code generated by mockery v1.1.2. DO NOT EDIT.

package mocks

import (
	models "authserver/models"

	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// Transaction is an autogenerated mock type for the Transaction type
type Transaction struct {
	mock.Mock
}

// CommitTransaction provides a mock function with given fields:
func (_m *Transaction) CommitTransaction() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateMigration provides a mock function with given fields: timestamp
func (_m *Transaction) CreateMigration(timestamp string) error {
	ret := _m.Called(timestamp)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(timestamp)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteAccessToken provides a mock function with given fields: token
func (_m *Transaction) DeleteAccessToken(token *models.AccessToken) error {
	ret := _m.Called(token)

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.AccessToken) error); ok {
		r0 = rf(token)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteAllOtherUserTokens provides a mock function with given fields: token
func (_m *Transaction) DeleteAllOtherUserTokens(token *models.AccessToken) error {
	ret := _m.Called(token)

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.AccessToken) error); ok {
		r0 = rf(token)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteMigrationByTimestamp provides a mock function with given fields: timestamp
func (_m *Transaction) DeleteMigrationByTimestamp(timestamp string) error {
	ret := _m.Called(timestamp)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(timestamp)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteUser provides a mock function with given fields: user
func (_m *Transaction) DeleteUser(user *models.User) error {
	ret := _m.Called(user)

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.User) error); ok {
		r0 = rf(user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAccessTokenByID provides a mock function with given fields: ID
func (_m *Transaction) GetAccessTokenByID(ID uuid.UUID) (*models.AccessToken, error) {
	ret := _m.Called(ID)

	var r0 *models.AccessToken
	if rf, ok := ret.Get(0).(func(uuid.UUID) *models.AccessToken); ok {
		r0 = rf(ID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.AccessToken)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = rf(ID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetClientByID provides a mock function with given fields: ID
func (_m *Transaction) GetClientByID(ID uuid.UUID) (*models.Client, error) {
	ret := _m.Called(ID)

	var r0 *models.Client
	if rf, ok := ret.Get(0).(func(uuid.UUID) *models.Client); ok {
		r0 = rf(ID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Client)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = rf(ID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetLatestTimestamp provides a mock function with given fields:
func (_m *Transaction) GetLatestTimestamp() (string, bool, error) {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 bool
	if rf, ok := ret.Get(1).(func() bool); ok {
		r1 = rf()
	} else {
		r1 = ret.Get(1).(bool)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func() error); ok {
		r2 = rf()
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetMigrationByTimestamp provides a mock function with given fields: timestamp
func (_m *Transaction) GetMigrationByTimestamp(timestamp string) (*models.Migration, error) {
	ret := _m.Called(timestamp)

	var r0 *models.Migration
	if rf, ok := ret.Get(0).(func(string) *models.Migration); ok {
		r0 = rf(timestamp)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Migration)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(timestamp)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetScopeByName provides a mock function with given fields: name
func (_m *Transaction) GetScopeByName(name string) (*models.Scope, error) {
	ret := _m.Called(name)

	var r0 *models.Scope
	if rf, ok := ret.Get(0).(func(string) *models.Scope); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Scope)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserByID provides a mock function with given fields: ID
func (_m *Transaction) GetUserByID(ID uuid.UUID) (*models.User, error) {
	ret := _m.Called(ID)

	var r0 *models.User
	if rf, ok := ret.Get(0).(func(uuid.UUID) *models.User); ok {
		r0 = rf(ID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = rf(ID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserByUsername provides a mock function with given fields: username
func (_m *Transaction) GetUserByUsername(username string) (*models.User, error) {
	ret := _m.Called(username)

	var r0 *models.User
	if rf, ok := ret.Get(0).(func(string) *models.User); ok {
		r0 = rf(username)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(username)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RollbackTransaction provides a mock function with given fields:
func (_m *Transaction) RollbackTransaction() {
	_m.Called()
}

// SaveAccessToken provides a mock function with given fields: token
func (_m *Transaction) SaveAccessToken(token *models.AccessToken) error {
	ret := _m.Called(token)

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.AccessToken) error); ok {
		r0 = rf(token)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SaveClient provides a mock function with given fields: client
func (_m *Transaction) SaveClient(client *models.Client) error {
	ret := _m.Called(client)

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.Client) error); ok {
		r0 = rf(client)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SaveScope provides a mock function with given fields: scope
func (_m *Transaction) SaveScope(scope *models.Scope) error {
	ret := _m.Called(scope)

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.Scope) error); ok {
		r0 = rf(scope)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SaveUser provides a mock function with given fields: user
func (_m *Transaction) SaveUser(user *models.User) error {
	ret := _m.Called(user)

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.User) error); ok {
		r0 = rf(user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Setup provides a mock function with given fields:
func (_m *Transaction) Setup() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateUser provides a mock function with given fields: user
func (_m *Transaction) UpdateUser(user *models.User) error {
	ret := _m.Called(user)

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.User) error); ok {
		r0 = rf(user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
