// Code generated by mockery v1.0.0

// This comment works around https://github.com/vektra/mockery/issues/155

package docker_helpers

import mock "github.com/stretchr/testify/mock"
import time "time"

// MockMachine is an autogenerated mock type for the Machine type
type MockMachine struct {
	mock.Mock
}

// CanConnect provides a mock function with given fields: name
func (_m *MockMachine) CanConnect(name string) bool {
	ret := _m.Called(name)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(name)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Create provides a mock function with given fields: driver, name, opts
func (_m *MockMachine) Create(driver string, name string, opts ...string) error {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, driver, name)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, ...string) error); ok {
		r0 = rf(driver, name, opts...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Credentials provides a mock function with given fields: name
func (_m *MockMachine) Credentials(name string) (DockerCredentials, error) {
	ret := _m.Called(name)

	var r0 DockerCredentials
	if rf, ok := ret.Get(0).(func(string) DockerCredentials); ok {
		r0 = rf(name)
	} else {
		r0 = ret.Get(0).(DockerCredentials)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Exist provides a mock function with given fields: name
func (_m *MockMachine) Exist(name string) bool {
	ret := _m.Called(name)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(name)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// List provides a mock function with given fields:
func (_m *MockMachine) List() ([]string, error) {
	ret := _m.Called()

	var r0 []string
	if rf, ok := ret.Get(0).(func() []string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Provision provides a mock function with given fields: name
func (_m *MockMachine) Provision(name string) error {
	ret := _m.Called(name)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(name)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Remove provides a mock function with given fields: name
func (_m *MockMachine) Remove(name string) error {
	ret := _m.Called(name)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(name)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Stop provides a mock function with given fields: name, timeout
func (_m *MockMachine) Stop(name string, timeout time.Duration) error {
	ret := _m.Called(name, timeout)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, time.Duration) error); ok {
		r0 = rf(name, timeout)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
