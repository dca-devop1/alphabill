package tokens

import (
	"bytes"
	"crypto"
	"errors"
	"fmt"

	"github.com/fxamacker/cbor/v2"

	"github.com/alphabill-org/alphabill/hash"
	"github.com/alphabill-org/alphabill/predicates/templates"
	"github.com/alphabill-org/alphabill/state"
	"github.com/alphabill-org/alphabill/txsystem"
	"github.com/alphabill-org/alphabill/types"
)

var (
	DustCollectorPredicate = templates.NewP2pkh256BytesFromKeyHash(hash.Sum256([]byte("dust collector")))
)

func handleBurnFungibleTokenTx(options *Options) txsystem.GenericExecuteFunc[BurnFungibleTokenAttributes] {
	return func(tx *types.TransactionOrder, attr *BurnFungibleTokenAttributes, currentBlockNo uint64) (*types.ServerMetadata, error) {
		if err := validateBurnFungibleToken(tx, attr, options.state, options.hashAlgorithm); err != nil {
			return nil, fmt.Errorf("invalid burn fungible token transaction: %w", err)
		}
		fee := options.feeCalculator()
		unitID := tx.UnitID()
		txHash := tx.Hash(options.hashAlgorithm)

		// 1. SetOwner(ι, DC)
		setOwnerFn := state.SetOwner(unitID, DustCollectorPredicate)

		// 2. UpdateData(ι, f), where f(D) = (0, S.n, H(P))
		updateUnitFn := state.UpdateUnitData(unitID,
			func(data state.UnitData) (state.UnitData, error) {
				ftData, ok := data.(*FungibleTokenData)
				if !ok {
					return nil, fmt.Errorf("unit %v does not contain fungible token data", unitID)
				}
				ftData.Value = 0
				ftData.T = currentBlockNo
				ftData.Backlink = txHash
				return ftData, nil
			},
		)

		if err := options.state.Apply(setOwnerFn, updateUnitFn); err != nil {
			return nil, fmt.Errorf("burnFToken: failed to update state: %w", err)
		}
		return &types.ServerMetadata{ActualFee: fee, TargetUnits: []types.UnitID{unitID}, SuccessIndicator: types.TxStatusSuccessful}, nil
	}
}

func validateBurnFungibleToken(tx *types.TransactionOrder, attr *BurnFungibleTokenAttributes, s *state.State, hashAlgorithm crypto.Hash) error {
	bearer, d, err := getFungibleTokenData(tx.UnitID(), s)
	if err != nil {
		return err
	}
	if d.Locked != 0 {
		return errors.New("token is locked")
	}
	if !bytes.Equal(d.TokenType, attr.TypeID) {
		return fmt.Errorf("type of token to burn does not matches the actual type of the token: expected %s, got %s", d.TokenType, attr.TypeID)
	}
	if attr.Value != d.Value {
		return fmt.Errorf("invalid token value: expected %v, got %v", d.Value, attr.Value)
	}
	if !bytes.Equal(d.Backlink, attr.Backlink) {
		return fmt.Errorf("invalid backlink: expected %X, got %X", d.Backlink, attr.Backlink)
	}
	predicates, err := getChainedPredicates[*FungibleTokenTypeData](
		hashAlgorithm,
		s,
		d.TokenType,
		func(d *FungibleTokenTypeData) []byte {
			return d.InvariantPredicate
		},
		func(d *FungibleTokenTypeData) types.UnitID {
			return d.ParentTypeId
		},
	)
	if err != nil {
		return err
	}
	return verifyOwnership(bearer, predicates, &burnFungibleTokenOwnershipProver{tx: tx, attr: attr})
}

type burnFungibleTokenOwnershipProver struct {
	tx   *types.TransactionOrder
	attr *BurnFungibleTokenAttributes
}

func (t *burnFungibleTokenOwnershipProver) OwnerProof() []byte {
	return t.tx.OwnerProof
}

func (t *burnFungibleTokenOwnershipProver) InvariantPredicateSignatures() [][]byte {
	return t.attr.InvariantPredicateSignatures
}

func (t *burnFungibleTokenOwnershipProver) SigBytes() ([]byte, error) {
	return t.tx.Payload.BytesWithAttributeSigBytes(t.attr)
}

func (b *BurnFungibleTokenAttributes) SigBytes() ([]byte, error) {
	// TODO: AB-1016 exclude InvariantPredicateSignatures from the payload hash because otherwise we have "chicken and egg" problem.
	signatureAttr := &BurnFungibleTokenAttributes{
		TypeID:                       b.TypeID,
		Value:                        b.Value,
		TargetTokenID:                b.TargetTokenID,
		TargetTokenBacklink:          b.TargetTokenBacklink,
		Backlink:                     b.Backlink,
		InvariantPredicateSignatures: nil,
	}
	return cbor.Marshal(signatureAttr)
}

func (b *BurnFungibleTokenAttributes) GetTypeID() types.UnitID {
	return b.TypeID
}

func (b *BurnFungibleTokenAttributes) SetTypeID(typeID types.UnitID) {
	b.TypeID = typeID
}

func (b *BurnFungibleTokenAttributes) GetValue() uint64 {
	return b.Value
}

func (b *BurnFungibleTokenAttributes) SetValue(value uint64) {
	b.Value = value
}

func (b *BurnFungibleTokenAttributes) GetTargetTokenID() []byte {
	return b.TargetTokenID
}

func (b *BurnFungibleTokenAttributes) SetTargetTokenID(targetTokenID []byte) {
	b.TargetTokenID = targetTokenID
}

func (b *BurnFungibleTokenAttributes) GetTargetTokenBacklink() []byte {
	return b.TargetTokenBacklink
}

func (b *BurnFungibleTokenAttributes) SetTargetTokenBacklink(targetTokenBacklink []byte) {
	b.TargetTokenBacklink = targetTokenBacklink
}

func (b *BurnFungibleTokenAttributes) GetBacklink() []byte {
	return b.Backlink
}

func (b *BurnFungibleTokenAttributes) SetBacklink(backlink []byte) {
	b.Backlink = backlink
}

func (b *BurnFungibleTokenAttributes) GetInvariantPredicateSignatures() [][]byte {
	return b.InvariantPredicateSignatures
}

func (b *BurnFungibleTokenAttributes) SetInvariantPredicateSignatures(signatures [][]byte) {
	b.InvariantPredicateSignatures = signatures
}
