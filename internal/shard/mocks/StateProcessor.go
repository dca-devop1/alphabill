// Code generated by mockery v2.10.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"

	txsystem "gitdc.ee.guardtime.com/alphabill/alphabill/internal/txsystem"
)

// StateProcessor is an autogenerated mock type for the StateProcessor type
type StateProcessor struct {
	mock.Mock
}

// Process provides a mock function with given fields: tx
func (_m *StateProcessor) Process(tx txsystem.GenericTransaction) error {
	ret := _m.Called(tx)

	var r0 error
	if rf, ok := ret.Get(0).(func(txsystem.GenericTransaction) error); ok {
		r0 = rf(tx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
