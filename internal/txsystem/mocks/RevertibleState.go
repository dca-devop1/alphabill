// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	state "gitdc.ee.guardtime.com/alphabill/alphabill/internal/txsystem/state"
	mock "github.com/stretchr/testify/mock"

	uint256 "github.com/holiman/uint256"
)

// RevertibleState is an autogenerated mock type for the RevertibleState type
type RevertibleState struct {
	mock.Mock
}

// AddItem provides a mock function with given fields: id, owner, data, stateHash
func (_m *RevertibleState) AddItem(id *uint256.Int, owner state.Predicate, data state.UnitData, stateHash []byte) error {
	ret := _m.Called(id, owner, data, stateHash)

	var r0 error
	if rf, ok := ret.Get(0).(func(*uint256.Int, state.Predicate, state.UnitData, []byte) error); ok {
		r0 = rf(id, owner, data, stateHash)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Commit provides a mock function with given fields:
func (_m *RevertibleState) Commit() {
	_m.Called()
}

// DeleteItem provides a mock function with given fields: id
func (_m *RevertibleState) DeleteItem(id *uint256.Int) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(*uint256.Int) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetRootHash provides a mock function with given fields:
func (_m *RevertibleState) GetRootHash() []byte {
	ret := _m.Called()

	var r0 []byte
	if rf, ok := ret.Get(0).(func() []byte); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	return r0
}

// Revert provides a mock function with given fields:
func (_m *RevertibleState) Revert() {
	_m.Called()
}

// SetOwner provides a mock function with given fields: id, owner, stateHash
func (_m *RevertibleState) SetOwner(id *uint256.Int, owner state.Predicate, stateHash []byte) error {
	ret := _m.Called(id, owner, stateHash)

	var r0 error
	if rf, ok := ret.Get(0).(func(*uint256.Int, state.Predicate, []byte) error); ok {
		r0 = rf(id, owner, stateHash)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// TotalValue provides a mock function with given fields:
func (_m *RevertibleState) TotalValue() state.SummaryValue {
	ret := _m.Called()

	var r0 state.SummaryValue
	if rf, ok := ret.Get(0).(func() state.SummaryValue); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(state.SummaryValue)
		}
	}

	return r0
}

// UpdateData provides a mock function with given fields: id, f, stateHash
func (_m *RevertibleState) UpdateData(id *uint256.Int, f state.UpdateFunction, stateHash []byte) error {
	ret := _m.Called(id, f, stateHash)

	var r0 error
	if rf, ok := ret.Get(0).(func(*uint256.Int, state.UpdateFunction, []byte) error); ok {
		r0 = rf(id, f, stateHash)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
