// Code generated by mockery v1.1.2. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// DBConnection is an autogenerated mock type for the DBConnection type
type DBConnection struct {
	mock.Mock
}

// CloseConnection provides a mock function with given fields:
func (_m *DBConnection) CloseConnection() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// OpenConnection provides a mock function with given fields:
func (_m *DBConnection) OpenConnection() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Ping provides a mock function with given fields:
func (_m *DBConnection) Ping() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
