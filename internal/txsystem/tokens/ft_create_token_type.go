package tokens

import (
	"errors"
	"fmt"

	"github.com/alphabill-org/alphabill/internal/rma"
	"github.com/alphabill-org/alphabill/internal/script"
	"github.com/alphabill-org/alphabill/internal/txsystem"
	"github.com/alphabill-org/alphabill/internal/txsystem/fc"
	"github.com/holiman/uint256"
)

func handleCreateFungibleTokenTypeTx(options *Options) txsystem.GenericExecuteFunc[*createFungibleTokenTypeWrapper] {
	return func(tx *createFungibleTokenTypeWrapper, currentBlockNr uint64) error {
		logger.Debug("Processing Create Fungible Token Type tx: %v", tx)
		if err := validateCreateFungibleTokenType(tx, options.state); err != nil {
			return fmt.Errorf("invalid create fungible token tx: %w", err)
		}
		fee := options.feeCalculator()
		tx.SetServerMetadata(&txsystem.ServerMetadata{Fee: fee})
		// update state
		fcrID := tx.transaction.GetClientFeeCreditRecordID()
		h := tx.Hash(options.hashAlgorithm)
		return options.state.AtomicUpdate(
			fc.DecrCredit(fcrID, fee, h),
			rma.AddItem(tx.UnitID(), script.PredicateAlwaysTrue(), newFungibleTokenTypeData(tx), h),
		)
	}
}

func validateCreateFungibleTokenType(tx *createFungibleTokenTypeWrapper, state *rma.Tree) error {
	unitID := tx.UnitID()
	if unitID.IsZero() {
		return errors.New(ErrStrUnitIDIsZero)
	}
	if len(tx.attributes.Symbol) > maxSymbolLength {
		return errors.New(ErrStrInvalidSymbolName)
	}
	decimalPlaces := tx.attributes.DecimalPlaces
	if decimalPlaces > maxDecimalPlaces {
		return fmt.Errorf("invalid decimal places. maximum allowed value %v, got %v", maxDecimalPlaces, decimalPlaces)
	}

	u, err := state.GetUnit(unitID)
	if u != nil {
		return fmt.Errorf("unit %v exists", unitID)
	}
	if !errors.Is(err, rma.ErrUnitNotFound) {
		return err
	}

	parentUnitID := tx.ParentTypeIdInt()
	if !parentUnitID.IsZero() {
		_, parentData, err := getUnit[*fungibleTokenTypeData](state, parentUnitID)
		if err != nil {
			return err
		}
		if decimalPlaces != parentData.decimalPlaces {
			return fmt.Errorf("invalid decimal places. allowed %v, got %v", parentData.decimalPlaces, decimalPlaces)
		}
	}
	predicates, err := getChainedPredicates[*fungibleTokenTypeData](
		state,
		tx.ParentTypeIdInt(),
		func(d *fungibleTokenTypeData) []byte {
			return d.subTypeCreationPredicate
		},
		func(d *fungibleTokenTypeData) *uint256.Int {
			return d.parentTypeId
		},
	)
	if err != nil {
		return err
	}
	return verifyPredicates(predicates, tx.SubTypeCreationPredicateSignatures(), tx.SigBytes())
}
