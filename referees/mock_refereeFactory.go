// Code generated by mockery v2.14.0. DO NOT EDIT.

package referees

import (
	logrus "github.com/sirupsen/logrus"
	mock "github.com/stretchr/testify/mock"
)

// mockRefereeFactory is an autogenerated mock type for the refereeFactory type
type mockRefereeFactory struct {
	mock.Mock
}

// Execute provides a mock function with given fields: executor, config, log
func (_m *mockRefereeFactory) Execute(executor interface{}, config *Config, log logrus.FieldLogger) Referee {
	ret := _m.Called(executor, config, log)

	var r0 Referee
	if rf, ok := ret.Get(0).(func(interface{}, *Config, logrus.FieldLogger) Referee); ok {
		r0 = rf(executor, config, log)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(Referee)
		}
	}

	return r0
}

type mockConstructorTestingTnewMockRefereeFactory interface {
	mock.TestingT
	Cleanup(func())
}

// newMockRefereeFactory creates a new instance of mockRefereeFactory. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockRefereeFactory(t mockConstructorTestingTnewMockRefereeFactory) *mockRefereeFactory {
	mock := &mockRefereeFactory{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}