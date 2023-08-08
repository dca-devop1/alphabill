package replication

import (
	"errors"
	"fmt"

	"github.com/alphabill-org/alphabill/internal/types"
)

const (
	Ok Status = iota
	InvalidRequestParameters
	UnknownSystemIdentifier
	BlocksNotFound
	Unknown
)

var (
	ErrLedgerReplicationRespIsNil = errors.New("ledger replication response is nil")
	ErrLedgerResponseBlocksIsNil  = errors.New("ledger response blocks is nil")
	ErrLedgerReplicationReqIsNil  = errors.New("ledger replication requests is nil")
	ErrInvalidSystemIdentifier    = errors.New("invalid system identifier")
	ErrNodeIdentifierIsMissing    = errors.New("node identifier is missing")
)

type (
	LedgerReplicationRequest struct {
		_                struct{} `cbor:",toarray"`
		SystemIdentifier []byte
		NodeIdentifier   string
		BeginBlockNumber uint64
		EndBlockNumber   uint64
	}

	LedgerReplicationResponse struct {
		_              struct{} `cbor:",toarray"`
		Status         Status
		Message        string
		Blocks         []*types.Block
		MaxRoundNumber uint64
	}

	Status int
)

func (r *LedgerReplicationResponse) Pretty() string {
	count := len(r.Blocks)
	// error message or no blocks
	if r.Message != "" {
		return fmt.Sprintf("status: %s, message: %s, %v blocks", r.Status.String(), r.Message, count)
	}
	return fmt.Sprintf("status: %s, %v blocks, max round number: %v", r.Status.String(), count, r.MaxRoundNumber)
}

func (r *LedgerReplicationResponse) IsValid() error {
	if r == nil {
		return ErrLedgerReplicationRespIsNil
	}
	if r.Blocks == nil {
		return ErrLedgerResponseBlocksIsNil
	}
	return nil
}

func (r *LedgerReplicationRequest) IsValid() error {
	if r == nil {
		return ErrLedgerReplicationReqIsNil
	}
	if len(r.SystemIdentifier) != 4 {
		return ErrInvalidSystemIdentifier
	}
	if r.NodeIdentifier == "" {
		return ErrNodeIdentifierIsMissing
	}
	if r.EndBlockNumber != 0 && r.EndBlockNumber < r.BeginBlockNumber {
		return fmt.Errorf("invalid block request range from %v to %v", r.BeginBlockNumber, r.EndBlockNumber)
	}
	return nil
}

func (s Status) String() string {
	switch s {
	case Ok:
		return "OK"
	case BlocksNotFound:
		return "Blocks Not Found"
	case InvalidRequestParameters:
		return "Invalid Request Parameters"
	case UnknownSystemIdentifier:
		return "Unknown System Identifier"
	case Unknown:
		return "Unknown"
	}
	return "Unknown Status Code"
}
