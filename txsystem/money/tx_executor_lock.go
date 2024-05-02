package money

import (
	"errors"
	"fmt"

	"github.com/alphabill-org/alphabill-go-sdk/txsystem/money"
	"github.com/alphabill-org/alphabill-go-sdk/types"

	"github.com/alphabill-org/alphabill/state"
	"github.com/alphabill-org/alphabill/txsystem"
)

var ErrInvalidLockStatus = errors.New("invalid lock status: expected non-zero value, got zero value")

func (m *Module) executeLockTx(tx *types.TransactionOrder, attr *money.LockAttributes, exeCtx *txsystem.TxExecutionContext) (*types.ServerMetadata, error) {
	// lock the unit
	unitID := tx.UnitID()
	action := state.UpdateUnitData(unitID, func(data types.UnitData) (types.UnitData, error) {
		newBillData, ok := data.(*money.BillData)
		if !ok {
			return nil, fmt.Errorf("unit %v does not contain bill data", unitID)
		}
		newBillData.Locked = attr.LockStatus
		newBillData.T = exeCtx.CurrentBlockNr
		newBillData.Counter += 1
		return newBillData, nil
	})
	if err := m.state.Apply(action); err != nil {
		return nil, fmt.Errorf("lock tx: failed to update state: %w", err)
	}
	return &types.ServerMetadata{ActualFee: m.feeCalculator(), TargetUnits: []types.UnitID{unitID}}, nil
}

func (m *Module) validateLockTx(tx *types.TransactionOrder, attr *money.LockAttributes, exeCtx *txsystem.TxExecutionContext) error {
	unitID := tx.UnitID()
	unit, err := m.state.GetUnit(unitID, false)
	if err != nil {
		return fmt.Errorf("lock tx: get unit error: %w", err)
	}
	if err = m.execPredicate(unit.Bearer(), tx.OwnerProof, tx); err != nil {
		return err
	}
	billData, ok := unit.Data().(*money.BillData)
	if !ok {
		return errors.New("lock tx: invalid unit type")
	}
	if attr == nil {
		return ErrTxAttrNil
	}
	// billData cannot be nil - it is an interface that must implement some methods
	if billData == nil {
		return ErrBillNil
	}
	if billData.IsLocked() {
		return errors.New("bill is already locked")
	}
	if attr.LockStatus == 0 {
		return ErrInvalidLockStatus
	}
	if billData.Counter != attr.Counter {
		return ErrInvalidCounter
	}
	return nil
}
