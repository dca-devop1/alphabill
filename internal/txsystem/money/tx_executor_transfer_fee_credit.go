package money

import (
	"bytes"
	"crypto"
	"errors"
	"fmt"

	"github.com/alphabill-org/alphabill/internal/state"
	"github.com/alphabill-org/alphabill/internal/txsystem"
	"github.com/alphabill-org/alphabill/internal/txsystem/fc"
	"github.com/alphabill-org/alphabill/internal/txsystem/fc/transactions"
	"github.com/alphabill-org/alphabill/internal/types"
)

var (
	ErrTxNil                     = errors.New("tx is nil")
	ErrBillNil                   = errors.New("bill is nil")
	ErrTargetSystemIdentifierNil = errors.New("TargetSystemIdentifier is nil")
	ErrTargetRecordIDNil         = errors.New("TargetRecordID is nil")
	ErrTargetRecordIDEmpty       = errors.New("TargetRecordID is empty")
	ErrAdditionTimeInvalid       = errors.New("EarliestAdditonTime is greater than LatestAdditionTime")
	ErrRecordIDExists            = errors.New("fee tx cannot contain fee credit reference")
	ErrFeeProofExists            = errors.New("fee tx cannot contain fee authorization proof")
	ErrInvalidFCValue            = errors.New("the amount to transfer cannot exceed the value of the bill")
	ErrInvalidFeeValue           = errors.New("the transaction max fee cannot exceed the transferred amount")
	ErrInvalidBacklink           = errors.New("the transaction backlink is not equal to unit backlink")
)

func handleTransferFeeCreditTx(s *state.State, hashAlgorithm crypto.Hash, feeCreditTxRecorder *feeCreditTxRecorder, feeCalc fc.FeeCalculator) txsystem.GenericExecuteFunc[transactions.TransferFeeCreditAttributes] {
	return func(tx *types.TransactionOrder, attr *transactions.TransferFeeCreditAttributes, currentBlockNumber uint64) (*types.ServerMetadata, error) {
		log.Debug("Processing transferFC %v", tx)
		unitID := tx.UnitID()
		unit, _ := s.GetUnit(unitID, false)
		if unit == nil {
			return nil, fmt.Errorf("transferFC: unit not found %X", tx.UnitID())
		}
		billData, ok := unit.Data().(*BillData)
		if !ok {
			return nil, errors.New("transferFC: invalid unit type")
		}
		if err := validateTransferFC(tx, attr, billData); err != nil {
			return nil, fmt.Errorf("transferFC: validation failed: %w", err)
		}

		// remove value from source unit, zero value bills get removed later
		action := state.UpdateUnitData(unitID, func(data state.UnitData) (state.UnitData, error) {
			newBillData, ok := data.(*BillData)
			if !ok {
				return nil, fmt.Errorf("unit %v does not contain bill data", unitID)
			}
			newBillData.V -= attr.Amount
			newBillData.T = currentBlockNumber
			newBillData.Backlink = tx.Hash(hashAlgorithm)
			return newBillData, nil
		})
		if err := s.Apply(action); err != nil {
			return nil, fmt.Errorf("transferFC: failed to update state: %w", err)
		}

		fee := feeCalc()

		// record fee tx for end of the round consolidation
		feeCreditTxRecorder.recordTransferFC(&transferFeeCreditTx{
			tx:   tx,
			attr: attr,
			fee:  fee,
		})
		return &types.ServerMetadata{ActualFee: fee, TargetUnits: []types.UnitID{tx.UnitID()}}, nil
	}
}

func validateTransferFC(tx *types.TransactionOrder, attr *transactions.TransferFeeCreditAttributes, bd *BillData) error {
	if tx == nil {
		return ErrTxNil
	}
	if bd == nil {
		return ErrBillNil
	}
	if attr.TargetSystemIdentifier == nil {
		return ErrTargetSystemIdentifierNil
	}
	if attr.TargetRecordID == nil {
		return ErrTargetRecordIDNil
	}
	if len(attr.TargetRecordID) == 0 {
		return ErrTargetRecordIDEmpty
	}
	if attr.EarliestAdditionTime > attr.LatestAdditionTime {
		return ErrAdditionTimeInvalid
	}
	if attr.Amount > bd.V {
		return ErrInvalidFCValue
	}
	if tx.Payload.ClientMetadata.MaxTransactionFee > attr.Amount {
		return ErrInvalidFeeValue
	}
	if !bytes.Equal(attr.Backlink, bd.Backlink) {
		return ErrInvalidBacklink
	}
	if tx.GetClientFeeCreditRecordID() != nil {
		return ErrRecordIDExists
	}
	if tx.FeeProof != nil {
		return ErrFeeProofExists
	}
	return nil
}
