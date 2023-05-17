package types

import (
	"crypto"
	"errors"
	"fmt"

	"github.com/alphabill-org/alphabill/internal/mt"
)

var (
	errBlockIsNil             = errors.New("block is nil")
	errBlockHeaderIsNil       = errors.New("block header is nil")
	errPrevBlockHashIsNil     = errors.New("previous block hash is nil")
	errBlockProposerIDMissing = errors.New("block proposer node identifier is missing")
	errTransactionsIsNil      = errors.New("transactions is nil")
	errSystemIDIsNil          = errors.New("system identifier is nil")
)

type (
	Block struct {
		_                  struct{} `cbor:",toarray"`
		Header             *Header
		Transactions       []*TransactionRecord
		UnicityCertificate *UnicityCertificate
	}

	Header struct {
		_                 struct{} `cbor:",toarray"`
		SystemID          SystemID
		ShardID           []byte
		ProposerID        string
		PreviousBlockHash []byte
	}
)

// Hash returns the hash of the block. Hash of a block is computed as hash of block header fields and tree hash
// of transactions.
func (b *Block) Hash(algorithm crypto.Hash) ([]byte, error) {
	if len(b.Transactions) == 0 {
		return make([]byte, algorithm.Size()), nil
	}
	// calculate merkle tree root hash from transactions
	tree, err := mt.New(algorithm, b.Transactions)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate merkle tree root hash: %w", err)
	}
	merkleRoot := tree.GetRootHash()

	// header hash
	headerHash := b.HeaderHash(algorithm)

	// header || merkle_root hash
	hasher := algorithm.New()
	hasher.Write(headerHash)
	hasher.Write(merkleRoot)
	return hasher.Sum(nil), nil
}

func (b *Block) HeaderHash(algorithm crypto.Hash) []byte {
	hasher := algorithm.New()
	hasher.Write(b.Header.SystemID)
	hasher.Write(b.Header.ShardID)
	hasher.Write(b.Header.PreviousBlockHash)
	hasher.Write([]byte(b.Header.ProposerID))
	return hasher.Sum(nil)
}

func (b *Block) GetRoundNumber() uint64 {
	if b != nil {
		return b.UnicityCertificate.GetRoundNumber()
	}
	return 0
}

func (b *Block) IsValid(v func(uc *UnicityCertificate) error) error {
	if b == nil {
		return errBlockIsNil
	}
	if b.Header == nil {
		return errBlockHeaderIsNil
	}
	if len(b.Header.SystemID) != 4 {
		return errSystemIDIsNil
	}
	// skip shard identifier for now, it is not used
	if b.Header.PreviousBlockHash == nil {
		return errPrevBlockHashIsNil
	}
	if len(b.Header.ProposerID) == 0 {
		return errBlockProposerIDMissing
	}
	if b.Transactions == nil {
		return errTransactionsIsNil
	}
	if b.UnicityCertificate == nil {
		return ErrUCIsNil
	}
	if err := v(b.UnicityCertificate); err != nil {
		return fmt.Errorf("unicity certificate validation failed, %w", err)
	}
	return nil
}

func (b *Block) GetProposerID() string {
	if b == nil || b.Header == nil {
		return ""
	}
	return b.Header.ProposerID
}
