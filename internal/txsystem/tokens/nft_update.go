package tokens

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/alphabill-org/alphabill/internal/rma"
	"github.com/alphabill-org/alphabill/internal/txsystem"
	"github.com/alphabill-org/alphabill/internal/txsystem/fc"
	"github.com/holiman/uint256"
)

func handleUpdateNonFungibleTokenTx(options *Options) txsystem.GenericExecuteFunc[*updateNonFungibleTokenWrapper] {
	return func(tx *updateNonFungibleTokenWrapper, currentBlockNr uint64) error {
		logger.Debug("Processing Update Non-Fungible Token tx: %v", tx)
		if err := validateUpdateNonFungibleToken(tx, options.state); err != nil {
			return fmt.Errorf("invalid update none-fungible token tx: %w", err)
		}
		h := tx.Hash(options.hashAlgorithm)
		fee := options.feeCalculator()
		tx.SetServerMetadata(&txsystem.ServerMetadata{Fee: fee})
		// update state
		fcrID := tx.transaction.GetClientFeeCreditRecordID()
		return options.state.AtomicUpdate(
			fc.DecrCredit(fcrID, fee, h),
			rma.UpdateData(tx.UnitID(), func(data rma.UnitData) (newData rma.UnitData) {
				d, ok := data.(*nonFungibleTokenData)
				if !ok {
					return data
				}
				d.data = tx.attributes.Data
				d.t = currentBlockNr
				d.backlink = tx.Hash(options.hashAlgorithm)
				return data
			}, h))
	}
}

func validateUpdateNonFungibleToken(tx *updateNonFungibleTokenWrapper, state *rma.Tree) error {
	if len(tx.attributes.Data) > dataMaxSize {
		return fmt.Errorf("data exceeds the maximum allowed size of %v KB", dataMaxSize)
	}
	unitID := tx.UnitID()
	u, err := state.GetUnit(unitID)
	if err != nil {
		return err
	}
	data, ok := u.Data.(*nonFungibleTokenData)
	if !ok {
		return fmt.Errorf("unit %v is not a non-fungible token type", unitID)
	}
	if !bytes.Equal(data.backlink, tx.attributes.Backlink) {
		return errors.New("invalid backlink")
	}
	predicates, err := getChainedPredicates[*nonFungibleTokenTypeData](
		state,
		data.nftType,
		func(d *nonFungibleTokenTypeData) []byte {
			return d.dataUpdatePredicate
		},
		func(d *nonFungibleTokenTypeData) *uint256.Int {
			return d.parentTypeId
		},
	)
	if err != nil {
		return err
	}
	predicates = append([]Predicate{data.dataUpdatePredicate}, predicates...)
	return verifyPredicates(predicates, tx.DataUpdateSignatures(), tx.SigBytes())
}
