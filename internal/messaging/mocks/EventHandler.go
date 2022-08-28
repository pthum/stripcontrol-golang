// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	null "github.com/pthum/null"
	model "github.com/pthum/stripcontrol-golang/internal/model"
	mock "github.com/stretchr/testify/mock"
)

// EventHandler is an autogenerated mock type for the EventHandler type
type EventHandler struct {
	mock.Mock
}

type EventHandler_Expecter struct {
	mock *mock.Mock
}

func (_m *EventHandler) EXPECT() *EventHandler_Expecter {
	return &EventHandler_Expecter{mock: &_m.Mock}
}

// Close provides a mock function with given fields:
func (_m *EventHandler) Close() {
	_m.Called()
}

// EventHandler_Close_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Close'
type EventHandler_Close_Call struct {
	*mock.Call
}

// Close is a helper method to define mock.On call
func (_e *EventHandler_Expecter) Close() *EventHandler_Close_Call {
	return &EventHandler_Close_Call{Call: _e.mock.On("Close")}
}

func (_c *EventHandler_Close_Call) Run(run func()) *EventHandler_Close_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *EventHandler_Close_Call) Return() *EventHandler_Close_Call {
	_c.Call.Return()
	return _c
}

// PublishProfileDeleteEvent provides a mock function with given fields: id
func (_m *EventHandler) PublishProfileDeleteEvent(id null.Int) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(null.Int) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// EventHandler_PublishProfileDeleteEvent_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'PublishProfileDeleteEvent'
type EventHandler_PublishProfileDeleteEvent_Call struct {
	*mock.Call
}

// PublishProfileDeleteEvent is a helper method to define mock.On call
//   - id null.Int
func (_e *EventHandler_Expecter) PublishProfileDeleteEvent(id interface{}) *EventHandler_PublishProfileDeleteEvent_Call {
	return &EventHandler_PublishProfileDeleteEvent_Call{Call: _e.mock.On("PublishProfileDeleteEvent", id)}
}

func (_c *EventHandler_PublishProfileDeleteEvent_Call) Run(run func(id null.Int)) *EventHandler_PublishProfileDeleteEvent_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(null.Int))
	})
	return _c
}

func (_c *EventHandler_PublishProfileDeleteEvent_Call) Return(err error) *EventHandler_PublishProfileDeleteEvent_Call {
	_c.Call.Return(err)
	return _c
}

// PublishProfileSaveEvent provides a mock function with given fields: id, profile
func (_m *EventHandler) PublishProfileSaveEvent(id null.Int, profile model.ColorProfile) error {
	ret := _m.Called(id, profile)

	var r0 error
	if rf, ok := ret.Get(0).(func(null.Int, model.ColorProfile) error); ok {
		r0 = rf(id, profile)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// EventHandler_PublishProfileSaveEvent_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'PublishProfileSaveEvent'
type EventHandler_PublishProfileSaveEvent_Call struct {
	*mock.Call
}

// PublishProfileSaveEvent is a helper method to define mock.On call
//   - id null.Int
//   - profile model.ColorProfile
func (_e *EventHandler_Expecter) PublishProfileSaveEvent(id interface{}, profile interface{}) *EventHandler_PublishProfileSaveEvent_Call {
	return &EventHandler_PublishProfileSaveEvent_Call{Call: _e.mock.On("PublishProfileSaveEvent", id, profile)}
}

func (_c *EventHandler_PublishProfileSaveEvent_Call) Run(run func(id null.Int, profile model.ColorProfile)) *EventHandler_PublishProfileSaveEvent_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(null.Int), args[1].(model.ColorProfile))
	})
	return _c
}

func (_c *EventHandler_PublishProfileSaveEvent_Call) Return(err error) *EventHandler_PublishProfileSaveEvent_Call {
	_c.Call.Return(err)
	return _c
}

// PublishStripDeleteEvent provides a mock function with given fields: id
func (_m *EventHandler) PublishStripDeleteEvent(id null.Int) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(null.Int) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// EventHandler_PublishStripDeleteEvent_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'PublishStripDeleteEvent'
type EventHandler_PublishStripDeleteEvent_Call struct {
	*mock.Call
}

// PublishStripDeleteEvent is a helper method to define mock.On call
//   - id null.Int
func (_e *EventHandler_Expecter) PublishStripDeleteEvent(id interface{}) *EventHandler_PublishStripDeleteEvent_Call {
	return &EventHandler_PublishStripDeleteEvent_Call{Call: _e.mock.On("PublishStripDeleteEvent", id)}
}

func (_c *EventHandler_PublishStripDeleteEvent_Call) Run(run func(id null.Int)) *EventHandler_PublishStripDeleteEvent_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(null.Int))
	})
	return _c
}

func (_c *EventHandler_PublishStripDeleteEvent_Call) Return(err error) *EventHandler_PublishStripDeleteEvent_Call {
	_c.Call.Return(err)
	return _c
}

// PublishStripSaveEvent provides a mock function with given fields: id, strip
func (_m *EventHandler) PublishStripSaveEvent(id null.Int, strip model.LedStrip) error {
	ret := _m.Called(id, strip)

	var r0 error
	if rf, ok := ret.Get(0).(func(null.Int, model.LedStrip) error); ok {
		r0 = rf(id, strip)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// EventHandler_PublishStripSaveEvent_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'PublishStripSaveEvent'
type EventHandler_PublishStripSaveEvent_Call struct {
	*mock.Call
}

// PublishStripSaveEvent is a helper method to define mock.On call
//   - id null.Int
//   - strip model.LedStrip
func (_e *EventHandler_Expecter) PublishStripSaveEvent(id interface{}, strip interface{}) *EventHandler_PublishStripSaveEvent_Call {
	return &EventHandler_PublishStripSaveEvent_Call{Call: _e.mock.On("PublishStripSaveEvent", id, strip)}
}

func (_c *EventHandler_PublishStripSaveEvent_Call) Run(run func(id null.Int, strip model.LedStrip)) *EventHandler_PublishStripSaveEvent_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(null.Int), args[1].(model.LedStrip))
	})
	return _c
}

func (_c *EventHandler_PublishStripSaveEvent_Call) Return(err error) *EventHandler_PublishStripSaveEvent_Call {
	_c.Call.Return(err)
	return _c
}

type mockConstructorTestingTNewEventHandler interface {
	mock.TestingT
	Cleanup(func())
}

// NewEventHandler creates a new instance of EventHandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewEventHandler(t mockConstructorTestingTNewEventHandler) *EventHandler {
	mock := &EventHandler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
