// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	models "github.com/bitrise-io/bitrise-step-analytics/models"
	mock "github.com/stretchr/testify/mock"
)

// Tracker is an autogenerated mock type for the Tracker type
type Tracker struct {
	mock.Mock
}

// Send provides a mock function with given fields: analytics
func (_m *Tracker) Send(analytics models.TrackEvent) error {
	ret := _m.Called(analytics)

	var r0 error
	if rf, ok := ret.Get(0).(func(models.TrackEvent) error); ok {
		r0 = rf(analytics)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}