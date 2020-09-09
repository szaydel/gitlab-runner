// Code generated by mockery v1.1.0. DO NOT EDIT.

package service

import mock "github.com/stretchr/testify/mock"

// MockEngine is an autogenerated mock type for the Engine type
type MockEngine struct {
	mock.Mock
}

// EngineName provides a mock function with given fields:
func (_m *MockEngine) EngineName() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// EnginePath provides a mock function with given fields:
func (_m *MockEngine) EnginePath() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}