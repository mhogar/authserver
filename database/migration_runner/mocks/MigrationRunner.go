// Code generated by mockery v1.1.2. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// MigrationRunner is an autogenerated mock type for the MigrationRunner type
type MigrationRunner struct {
	mock.Mock
}

// MigrateDown provides a mock function with given fields:
func (_m *MigrationRunner) MigrateDown() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MigrateUp provides a mock function with given fields:
func (_m *MigrationRunner) MigrateUp() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
