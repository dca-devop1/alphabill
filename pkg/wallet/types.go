package wallet

import (
	"github.com/alphabill-org/alphabill/internal/hash"
	"github.com/alphabill-org/alphabill/internal/types"
)

type BlockProcessor interface {
	// ProcessBlock signals given block to be processed
	// any error returned here signals block processor to terminate,
	ProcessBlock(b *types.Block) error
}

type TxHash []byte

type Transactions struct {
	_            struct{} `cbor:",toarray"`
	Transactions []*types.TransactionOrder
}

type UnitID []byte

type Predicate []byte

type PubKey []byte

type PubKeyHash []byte

// TxProof type alias for block.TxProof, can be removed once block package is moved out of internal
type TxProof = types.TxProof

// Proof wrapper struct around TxRecord and TxProof
type Proof struct {
	_        struct{}                 `cbor:",toarray"`
	TxRecord *types.TransactionRecord `json:"txRecord"`
	TxProof  *types.TxProof           `json:"txProof"`
}

func (pk PubKey) Hash() PubKeyHash {
	return hash.Sum256(pk)
}

type (
	TxHistoryRecord struct {
		_            struct{} `cbor:",toarray"`
		UnitID       UnitID
		TxHash       TxHash
		CounterParty PubKeyHash
		Confirmed    bool
		Kind         TxHistoryRecordKind
	}

	TxHistoryRecordKind uint8
)

const (
	OUTGOING TxHistoryRecordKind = iota
	INCOMING
)
