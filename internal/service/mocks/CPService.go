// Code generated by mockery v2.36.0. DO NOT EDIT.

package servicemocks

import (
	model "github.com/pthum/stripcontrol-golang/internal/model"
	mock "github.com/stretchr/testify/mock"
)

// CPService is an autogenerated mock type for the CPService type
type CPService struct {
	mock.Mock
}

type CPService_Expecter struct {
	mock *mock.Mock
}

func (_m *CPService) EXPECT() *CPService_Expecter {
	return &CPService_Expecter{mock: &_m.Mock}
}

// CreateColorProfile provides a mock function with given fields: mdl
func (_m *CPService) CreateColorProfile(mdl *model.ColorProfile) error {
	ret := _m.Called(mdl)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.ColorProfile) error); ok {
		r0 = rf(mdl)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CPService_CreateColorProfile_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateColorProfile'
type CPService_CreateColorProfile_Call struct {
	*mock.Call
}

// CreateColorProfile is a helper method to define mock.On call
//   - mdl *model.ColorProfile
func (_e *CPService_Expecter) CreateColorProfile(mdl interface{}) *CPService_CreateColorProfile_Call {
	return &CPService_CreateColorProfile_Call{Call: _e.mock.On("CreateColorProfile", mdl)}
}

func (_c *CPService_CreateColorProfile_Call) Run(run func(mdl *model.ColorProfile)) *CPService_CreateColorProfile_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*model.ColorProfile))
	})
	return _c
}

func (_c *CPService_CreateColorProfile_Call) Return(_a0 error) *CPService_CreateColorProfile_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *CPService_CreateColorProfile_Call) RunAndReturn(run func(*model.ColorProfile) error) *CPService_CreateColorProfile_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteColorProfile provides a mock function with given fields: id
func (_m *CPService) DeleteColorProfile(id string) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CPService_DeleteColorProfile_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteColorProfile'
type CPService_DeleteColorProfile_Call struct {
	*mock.Call
}

// DeleteColorProfile is a helper method to define mock.On call
//   - id string
func (_e *CPService_Expecter) DeleteColorProfile(id interface{}) *CPService_DeleteColorProfile_Call {
	return &CPService_DeleteColorProfile_Call{Call: _e.mock.On("DeleteColorProfile", id)}
}

func (_c *CPService_DeleteColorProfile_Call) Run(run func(id string)) *CPService_DeleteColorProfile_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *CPService_DeleteColorProfile_Call) Return(_a0 error) *CPService_DeleteColorProfile_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *CPService_DeleteColorProfile_Call) RunAndReturn(run func(string) error) *CPService_DeleteColorProfile_Call {
	_c.Call.Return(run)
	return _c
}

// GetAll provides a mock function with given fields:
func (_m *CPService) GetAll() ([]model.ColorProfile, error) {
	ret := _m.Called()

	var r0 []model.ColorProfile
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]model.ColorProfile, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []model.ColorProfile); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.ColorProfile)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CPService_GetAll_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAll'
type CPService_GetAll_Call struct {
	*mock.Call
}

// GetAll is a helper method to define mock.On call
func (_e *CPService_Expecter) GetAll() *CPService_GetAll_Call {
	return &CPService_GetAll_Call{Call: _e.mock.On("GetAll")}
}

func (_c *CPService_GetAll_Call) Run(run func()) *CPService_GetAll_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *CPService_GetAll_Call) Return(_a0 []model.ColorProfile, _a1 error) *CPService_GetAll_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *CPService_GetAll_Call) RunAndReturn(run func() ([]model.ColorProfile, error)) *CPService_GetAll_Call {
	_c.Call.Return(run)
	return _c
}

// GetColorProfile provides a mock function with given fields: id
func (_m *CPService) GetColorProfile(id string) (*model.ColorProfile, error) {
	ret := _m.Called(id)

	var r0 *model.ColorProfile
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*model.ColorProfile, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(string) *model.ColorProfile); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.ColorProfile)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CPService_GetColorProfile_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetColorProfile'
type CPService_GetColorProfile_Call struct {
	*mock.Call
}

// GetColorProfile is a helper method to define mock.On call
//   - id string
func (_e *CPService_Expecter) GetColorProfile(id interface{}) *CPService_GetColorProfile_Call {
	return &CPService_GetColorProfile_Call{Call: _e.mock.On("GetColorProfile", id)}
}

func (_c *CPService_GetColorProfile_Call) Run(run func(id string)) *CPService_GetColorProfile_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *CPService_GetColorProfile_Call) Return(_a0 *model.ColorProfile, _a1 error) *CPService_GetColorProfile_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *CPService_GetColorProfile_Call) RunAndReturn(run func(string) (*model.ColorProfile, error)) *CPService_GetColorProfile_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateColorProfile provides a mock function with given fields: id, updMdl
func (_m *CPService) UpdateColorProfile(id string, updMdl model.ColorProfile) error {
	ret := _m.Called(id, updMdl)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, model.ColorProfile) error); ok {
		r0 = rf(id, updMdl)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CPService_UpdateColorProfile_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateColorProfile'
type CPService_UpdateColorProfile_Call struct {
	*mock.Call
}

// UpdateColorProfile is a helper method to define mock.On call
//   - id string
//   - updMdl model.ColorProfile
func (_e *CPService_Expecter) UpdateColorProfile(id interface{}, updMdl interface{}) *CPService_UpdateColorProfile_Call {
	return &CPService_UpdateColorProfile_Call{Call: _e.mock.On("UpdateColorProfile", id, updMdl)}
}

func (_c *CPService_UpdateColorProfile_Call) Run(run func(id string, updMdl model.ColorProfile)) *CPService_UpdateColorProfile_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(model.ColorProfile))
	})
	return _c
}

func (_c *CPService_UpdateColorProfile_Call) Return(_a0 error) *CPService_UpdateColorProfile_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *CPService_UpdateColorProfile_Call) RunAndReturn(run func(string, model.ColorProfile) error) *CPService_UpdateColorProfile_Call {
	_c.Call.Return(run)
	return _c
}

// NewCPService creates a new instance of CPService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCPService(t interface {
	mock.TestingT
	Cleanup(func())
}) *CPService {
	mock := &CPService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
