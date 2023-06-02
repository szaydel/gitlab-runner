// Code generated by mockery v2.28.2. DO NOT EDIT.

package common

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockJobTrace is an autogenerated mock type for the JobTrace type
type MockJobTrace struct {
	mock.Mock
}

// Abort provides a mock function with given fields:
func (_m *MockJobTrace) Abort() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Cancel provides a mock function with given fields:
func (_m *MockJobTrace) Cancel() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Fail provides a mock function with given fields: err, failureData
func (_m *MockJobTrace) Fail(err error, failureData JobFailureData) {
	_m.Called(err, failureData)
}

// IsStdout provides a mock function with given fields:
func (_m *MockJobTrace) IsStdout() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// SetAbortFunc provides a mock function with given fields: abortFunc
func (_m *MockJobTrace) SetAbortFunc(abortFunc context.CancelFunc) {
	_m.Called(abortFunc)
}

// SetCancelFunc provides a mock function with given fields: cancelFunc
func (_m *MockJobTrace) SetCancelFunc(cancelFunc context.CancelFunc) {
	_m.Called(cancelFunc)
}

// SetDebugModeEnabled provides a mock function with given fields: isEnabled
func (_m *MockJobTrace) SetDebugModeEnabled(isEnabled bool) {
	_m.Called(isEnabled)
}

// SetFailuresCollector provides a mock function with given fields: fc
func (_m *MockJobTrace) SetFailuresCollector(fc FailuresCollector) {
	_m.Called(fc)
}

// SetMasked provides a mock function with given fields: maskOptions
func (_m *MockJobTrace) SetMasked(maskOptions MaskOptions) {
	_m.Called(maskOptions)
}

// SetSupportedFailureReasonMapper provides a mock function with given fields: f
func (_m *MockJobTrace) SetSupportedFailureReasonMapper(f SupportedFailureReasonMapper) {
	_m.Called(f)
}

// Success provides a mock function with given fields:
func (_m *MockJobTrace) Success() {
	_m.Called()
}

// Write provides a mock function with given fields: p
func (_m *MockJobTrace) Write(p []byte) (int, error) {
	ret := _m.Called(p)

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func([]byte) (int, error)); ok {
		return rf(p)
	}
	if rf, ok := ret.Get(0).(func([]byte) int); ok {
		r0 = rf(p)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func([]byte) error); ok {
		r1 = rf(p)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewMockJobTrace interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockJobTrace creates a new instance of MockJobTrace. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockJobTrace(t mockConstructorTestingTNewMockJobTrace) *MockJobTrace {
	mock := &MockJobTrace{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
