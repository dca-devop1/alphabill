package evm

import (
	"fmt"

	"github.com/alphabill-org/alphabill/internal/rma"
	"github.com/alphabill-org/alphabill/internal/script"
	"github.com/alphabill-org/alphabill/internal/txsystem"
	"github.com/alphabill-org/alphabill/internal/txsystem/evm/statedb"
	"github.com/alphabill-org/alphabill/internal/txsystem/fc/transactions"
	"github.com/alphabill-org/alphabill/internal/types"
)

func isFeeCreditTx(tx *types.TransactionOrder) bool {
	typeUrl := tx.PayloadType()
	return typeUrl == transactions.PayloadTypeAddFeeCredit ||
		typeUrl == transactions.PayloadTypeCloseFeeCredit
}

func checkFeeAccountBalance(state *rma.Tree) txsystem.GenericTransactionValidator {
	return func(ctx *txsystem.TxValidationContext) error {
		if isFeeCreditTx(ctx.Tx) {
			addr, err := getAddressFromPredicateArg(ctx.Tx.OwnerProof)
			if err != nil {
				return fmt.Errorf("failed to extract address from public key bytes, %w", err)
			}
			stateDB := statedb.NewStateDB(state)
			feeData := stateDB.GetFeeData(addr)
			if feeData == nil && ctx.Tx.PayloadType() == transactions.PayloadTypeCloseFeeCredit {
				return fmt.Errorf("no fee credit info found for unit %X", ctx.Tx.UnitID())
			}
			if feeData == nil && ctx.Tx.PayloadType() == transactions.PayloadTypeAddFeeCredit {
				// account creation
				return nil
			}
			// owner proof verifies correctly
			payloadBytes, err := ctx.Tx.PayloadBytes()
			if err != nil {
				return fmt.Errorf("failed to marshal payload bytes: %w", err)
			}

			if err = script.RunScript(ctx.Tx.OwnerProof, feeData.Bearer, payloadBytes); err != nil {
				return fmt.Errorf("invalid owner proof: %w", err)
			}
		}
		// todo: To be defined: TX is either deploy or call SC
		return nil
	}
}
