// Code generated by mockery v2.28.2. DO NOT EDIT.

package pull

import mock "github.com/stretchr/testify/mock"

// mockPullLogger is an autogenerated mock type for the pullLogger type
type mockPullLogger struct {
	mock.Mock
}

// Debugln provides a mock function with given fields: args
func (_m *mockPullLogger) Debugln(args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// Infoln provides a mock function with given fields: args
func (_m *mockPullLogger) Infoln(args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// Println provides a mock function with given fields: args
func (_m *mockPullLogger) Println(args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// Warningln provides a mock function with given fields: args
func (_m *mockPullLogger) Warningln(args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

type mockConstructorTestingTnewMockPullLogger interface {
	mock.TestingT
	Cleanup(func())
}

// newMockPullLogger creates a new instance of mockPullLogger. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockPullLogger(t mockConstructorTestingTnewMockPullLogger) *mockPullLogger {
	mock := &mockPullLogger{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
