package money

import (
	"errors"
	"fmt"

	"github.com/alphabill-org/alphabill-go-sdk/txsystem/money"
	"github.com/alphabill-org/alphabill-go-sdk/types"

	"github.com/alphabill-org/alphabill/state"
	"github.com/alphabill-org/alphabill/txsystem"
)

var (
	ErrBillUnlocked = errors.New("bill is already unlocked")
)

func (m *Module) executeUnlockTx(tx *types.TransactionOrder, _ *money.UnlockAttributes, exeCtx *txsystem.TxExecutionContext) (*types.ServerMetadata, error) {
	// unlock the unit
	unitID := tx.UnitID()
	action := state.UpdateUnitData(unitID, func(data types.UnitData) (types.UnitData, error) {
		newBillData, ok := data.(*money.BillData)
		if !ok {
			return nil, fmt.Errorf("unlock tx: unit %v does not contain bill data", unitID)
		}
		newBillData.Locked = 0
		newBillData.T = exeCtx.CurrentBlockNr
		newBillData.Counter += 1
		return newBillData, nil
	})
	if err := m.state.Apply(action); err != nil {
		return nil, fmt.Errorf("unlock tx: failed to update state: %w", err)
	}
	return &types.ServerMetadata{ActualFee: m.feeCalculator(), TargetUnits: []types.UnitID{unitID}}, nil
}

func (m *Module) validateUnlockTx(tx *types.TransactionOrder, attr *money.UnlockAttributes, exeCtx *txsystem.TxExecutionContext) error {
	unitID := tx.UnitID()
	unit, err := m.state.GetUnit(unitID, false)
	if err != nil {
		return fmt.Errorf("unlock tx: get unit error: %w", err)
	}
	if err = m.execPredicate(unit.Bearer(), tx.OwnerProof, tx); err != nil {
		return err
	}
	billData, ok := unit.Data().(*money.BillData)
	if !ok {
		return errors.New("unlock tx: invalid unit type")
	}
	if attr == nil {
		return ErrTxAttrNil
	}
	if billData == nil {
		return ErrBillNil
	}
	if !billData.IsLocked() {
		return ErrBillUnlocked
	}
	if billData.Counter != attr.Counter {
		return ErrInvalidCounter
	}
	return nil
}
