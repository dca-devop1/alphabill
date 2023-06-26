package state

import (
	"crypto"
	"testing"

	hasherUtil "github.com/alphabill-org/alphabill/internal/hash"
	"github.com/alphabill-org/alphabill/internal/script"
	test "github.com/alphabill-org/alphabill/internal/testutils"
	"github.com/alphabill-org/alphabill/internal/types"
	"github.com/alphabill-org/alphabill/internal/util"
	"github.com/stretchr/testify/require"
)

func TestAdd(t *testing.T) {
	type args struct {
		id     types.UnitID
		bearer []byte
		data   UnitData
	}
	type testCase struct {
		name            string
		args            args
		initialState    *State
		executionErrStr string
		expectedUnit    *Unit
	}
	tests := []testCase{
		{
			name: "unit id is nil",
			args: args{
				id: nil,
			},
			initialState:    newEmptyState(t),
			executionErrStr: "id is nil",
		},
		{
			name: "unit ID exists",
			args: args{
				id:     []byte{1, 1, 1, 1},
				bearer: script.PredicateAlwaysTrue(),
				data: &TestData{
					Value: 100,
				},
			},
			initialState:    newStateWithUnits(t),
			executionErrStr: "unable to add unit: key 01010101 exists",
		},
		{
			name: "ok",
			args: args{
				id:     []byte{1},
				bearer: script.PredicateAlwaysTrue(),
				data:   &TestData{Value: 123},
			},
			initialState: newEmptyState(t),
			expectedUnit: &Unit{
				unitLedgerHeadHash:  nil,
				logs:                nil,
				logRoot:             nil,
				bearer:              script.PredicateAlwaysTrue(),
				data:                &TestData{Value: 123},
				subTreeSummaryValue: 123,
				subTreeSummaryHash: hasherUtil.Sum(crypto.SHA256,
					[]byte{1},
					nil, // h_s is nil (we do not have a log entry)
					util.Uint64ToBytes(123),
					util.Uint64ToBytes(0),
					make([]byte, 32),
					util.Uint64ToBytes(0),
					make([]byte, 32),
				),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			add := AddUnit(tt.args.id, tt.args.bearer, tt.args.data)
			err := add(tt.initialState.latestSavepoint(), crypto.SHA256)
			if tt.executionErrStr != "" {
				require.ErrorContains(t, err, tt.executionErrStr)
			}
			if tt.expectedUnit != nil {
				assertUnit(t, tt.initialState, tt.args.id, tt.expectedUnit, false)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	type args struct {
		id types.UnitID
		f  UpdateFunction
	}
	type testCase struct {
		name            string
		args            args
		initialState    *State
		executionErrStr string
		expectedUnit    *Unit
	}
	tests := []testCase{
		{
			name: "not found",
			args: args{
				id: []byte{1},
				f: func(data UnitData) UnitData {
					return data
				},
			},
			initialState:    newEmptyState(t),
			executionErrStr: "failed to get unit: item 01 does not exist: not found",
		},
		{
			name: "update function is nil",
			args: args{
				id: test.RandomBytes(32),
			},
			initialState:    newEmptyState(t),
			executionErrStr: "update function is nil",
		},
		{
			name: "ok",
			args: args{
				id: []byte{1, 1, 1, 1},
				f: func(data UnitData) UnitData {
					return &TestData{Value: 200}
				},
			},
			initialState: newStateWithUnits(t),
			expectedUnit: &Unit{
				unitLedgerHeadHash:  nil,
				logs:                nil,
				logRoot:             nil,
				bearer:              script.PredicateAlwaysTrue(),
				data:                &TestData{Value: 200},
				subTreeSummaryValue: 10,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := UpdateUnitData(tt.args.id, tt.args.f)
			err := f(tt.initialState.latestSavepoint(), crypto.SHA256)
			if tt.executionErrStr != "" {
				require.ErrorContains(t, err, tt.executionErrStr)
			}
			if tt.expectedUnit != nil {
				assertUnit(t, tt.initialState, tt.args.id, tt.expectedUnit, false)
			}
		})
	}
}

func TestDelete(t *testing.T) {

	type testCase struct {
		name            string
		unitID          types.UnitID
		initialState    *State
		executionErrStr string
	}
	tests := []testCase{
		{
			name:            "unit ID is nil",
			unitID:          nil,
			initialState:    newEmptyState(t),
			executionErrStr: "id is nil",
		},
		{
			name:            "unit ID not found",
			unitID:          []byte{1},
			initialState:    newEmptyState(t),
			executionErrStr: "unable to delete unit",
		},
		{
			name:         "ok",
			unitID:       []byte{1, 1, 1, 1},
			initialState: newStateWithUnits(t),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := DeleteUnit(tt.unitID)(tt.initialState.latestSavepoint(), crypto.SHA256)
			if tt.executionErrStr != "" {
				require.ErrorContains(t, err, tt.executionErrStr)
				return
			}
			require.NoError(t, err)
			u, err := tt.initialState.latestSavepoint().Get(tt.unitID)
			require.ErrorContains(t, err, "not found")
			require.Nil(t, u)
		})
	}
}

func TestSetOwner(t *testing.T) {
	type args struct {
		id       types.UnitID
		newOwner []byte
	}
	type testCase struct {
		name            string
		args            args
		initialState    *State
		executionErrStr string
		expectedUnit    *Unit
	}
	tests := []testCase{
		{
			name:            "unit ID is nil",
			args:            args{},
			initialState:    newEmptyState(t),
			executionErrStr: "id is nil",
		},
		{
			name: "unit ID not found",
			args: args{
				id: []byte{1},
			},
			initialState:    newEmptyState(t),
			executionErrStr: "not found",
		},
		{
			name: "ok",
			args: args{
				id:       []byte{1, 1, 1, 1},
				newOwner: []byte{1, 2, 3, 4, 5},
			},
			initialState: newStateWithUnits(t),
			expectedUnit: &Unit{
				unitLedgerHeadHash:  nil,
				logs:                nil,
				logRoot:             nil,
				bearer:              []byte{1, 2, 3, 4, 5},
				data:                &TestData{Value: 10},
				subTreeSummaryValue: 10,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SetOwner(tt.args.id, tt.args.newOwner)(tt.initialState.latestSavepoint(), crypto.SHA256)
			if tt.executionErrStr != "" {
				require.ErrorContains(t, err, tt.executionErrStr)
				return
			}
			require.NoError(t, err)
			assertUnit(t, tt.initialState, tt.args.id, tt.expectedUnit, false)
		})
	}
}

func newEmptyState(t *testing.T) *State {
	emptyState, err := New()
	require.NoError(t, err)
	return emptyState
}

func newStateWithUnits(t *testing.T) *State {
	s := newEmptyState(t)
	require.NoError(t,
		s.Apply(
			AddUnit(
				[]byte{1, 1, 1, 1},
				script.PredicateAlwaysTrue(),
				&TestData{Value: 10},
			),
		),
	)
	return s
}

func assertUnit(t *testing.T, state *State, unitID types.UnitID, expectedUnit *Unit, committed bool) {
	t.Helper()
	unit, err := state.latestSavepoint().Get(unitID)
	require.NoError(t, err)
	require.NotNil(t, unit)
	assertUnitEqual(t, expectedUnit, unit, committed)

	committedUnit, err := state.committedTree.Get(unitID)
	if !committed {
		require.ErrorContains(t, err, "not found")
	} else {
		require.NoError(t, err)
		require.NotNil(t, committedUnit)
		assertUnitEqual(t, expectedUnit, unit, committed)
	}
}

func assertUnitEqual(t *testing.T, expectedUnit *Unit, unit *Unit, committed bool) {
	require.Equal(t, expectedUnit.data, unit.data)
	if !committed {
		require.Nil(t, unit.unitLedgerHeadHash)
	} else {
		require.NotNil(t, unit.unitLedgerHeadHash)
	}
	require.Equal(t, expectedUnit.subTreeSummaryValue, unit.subTreeSummaryValue)
	require.Equal(t, expectedUnit.bearer, unit.bearer)
}
