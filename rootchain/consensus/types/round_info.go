package types

import (
	"bytes"
	gocrypto "crypto"
	"errors"
	"fmt"

	"github.com/alphabill-org/alphabill-go-base/types/hex"
	"github.com/alphabill-org/alphabill-go-base/util"
)

var (
	errRoundNumberUnassigned   = errors.New("round number is not assigned")
	errParentRoundUnassigned   = errors.New("parent round number is not assigned")
	errRootHashUnassigned      = errors.New("root hash is not assigned")
	errRoundCreationTimeNotSet = errors.New("round creation time is not set")
)

type RoundInfo struct {
	_                 struct{}  `cbor:",toarray"`
	RoundNumber       uint64    `json:"rootChainRoundNumber"`
	Epoch             uint64    `json:"rootEpoch"`
	Timestamp         uint64    `json:"timestamp"`
	ParentRoundNumber uint64    `json:"rootChainParentRoundNumber"`
	CurrentRootHash   hex.Bytes `json:"currentRootHash"`
}

func (x *RoundInfo) GetRound() uint64 {
	if x != nil {
		return x.RoundNumber
	}
	return 0
}

func (x *RoundInfo) GetParentRound() uint64 {
	if x != nil {
		return x.ParentRoundNumber
	}
	return 0
}

func (x *RoundInfo) Hash(hash gocrypto.Hash) []byte {
	hasher := hash.New()
	hasher.Write(x.Bytes())
	return hasher.Sum(nil)
}

func (x *RoundInfo) Bytes() []byte {
	var b bytes.Buffer
	b.Write(util.Uint64ToBytes(x.RoundNumber))
	b.Write(util.Uint64ToBytes(x.Epoch))
	b.Write(util.Uint64ToBytes(x.Timestamp))
	b.Write(util.Uint64ToBytes(x.ParentRoundNumber))
	b.Write(x.CurrentRootHash)
	return b.Bytes()
}

func (x *RoundInfo) IsValid() error {
	if x.RoundNumber == 0 {
		return errRoundNumberUnassigned
	}
	if x.ParentRoundNumber == 0 && x.RoundNumber > 1 {
		return errParentRoundUnassigned
	}
	if x.RoundNumber <= x.ParentRoundNumber {
		return fmt.Errorf("invalid round number %d - must be greater than parent round %d", x.RoundNumber, x.ParentRoundNumber)
	}
	if len(x.CurrentRootHash) < 1 {
		return errRootHashUnassigned
	}
	if x.Timestamp == 0 {
		return errRoundCreationTimeNotSet
	}
	return nil
}
