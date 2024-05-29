package tokenenc

import (
	"errors"

	"github.com/alphabill-org/alphabill-go-base/txsystem/tokens"
	"github.com/alphabill-org/alphabill-go-base/types"
	"github.com/alphabill-org/alphabill/predicates/wasm/wvm/encoder"
)

func RegisterUnitDataEncoders(reg func(ud any, enc encoder.UnitDataEncoder) error) error {
	return errors.Join(
		reg(&tokens.NonFungibleTokenData{}, udeNonFungibleTokenData),
		reg(&tokens.NonFungibleTokenTypeData{}, udeNonFungibleTokenTypeData),
	)
}

func udeNonFungibleTokenData(data types.UnitData, ver uint32) ([]byte, error) {
	value := data.(*tokens.NonFungibleTokenData)
	buf := encoder.TVEnc{}
	buf.EncodeTagged(1, value.TypeID)
	buf.EncodeTagged(2, value.Name)
	buf.EncodeTagged(3, value.URI)
	buf.EncodeTagged(4, value.Data)
	buf.EncodeTagged(5, value.T)
	buf.EncodeTagged(6, value.Counter)
	buf.EncodeTagged(7, value.Locked)
	return buf.Bytes()
}

func udeNonFungibleTokenTypeData(data types.UnitData, ver uint32) ([]byte, error) {
	value := data.(*tokens.NonFungibleTokenTypeData)
	buf := encoder.TVEnc{}
	buf.EncodeTagged(1, value.ParentTypeID)
	buf.EncodeTagged(2, value.Symbol)
	buf.EncodeTagged(3, value.Name)
	return buf.Bytes()
}
