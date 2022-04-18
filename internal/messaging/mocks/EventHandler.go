// Code generated by mockery v2.10.6. DO NOT EDIT.

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

// Close provides a mock function with given fields:
func (_m *EventHandler) Close() {
	_m.Called()
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
