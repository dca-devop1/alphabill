package replication

import (
	"errors"
	"fmt"

	"github.com/alphabill-org/alphabill-go-base/types"
	"github.com/google/uuid"
)

const (
	Ok Status = iota
	InvalidRequestParameters
	UnknownPartitionIdentifier
	BlocksNotFound
	Unknown
)

var (
	ErrLedgerReplicationRespIsNil = errors.New("ledger replication response is nil")
	ErrLedgerResponseBlocksIsNil  = errors.New("ledger response blocks is nil")
	ErrLedgerReplicationReqIsNil  = errors.New("ledger replication requests is nil")
	ErrInvalidPartitionIdentifier = errors.New("invalid partition identifier")
	ErrNodeIdentifierIsMissing    = errors.New("node identifier is missing")
)

type (
	LedgerReplicationRequest struct {
		_                   struct{} `cbor:",toarray"`
		UUID                uuid.UUID
		PartitionIdentifier types.PartitionID
		NodeIdentifier      string
		BeginBlockNumber    uint64
		EndBlockNumber      uint64
	}

	LedgerReplicationResponse struct {
		_                struct{} `cbor:",toarray"`
		UUID             uuid.UUID
		Status           Status
		Message          string
		Blocks           []*types.Block
		FirstBlockNumber uint64
		LastBlockNumber  uint64
	}

	Status int
)

func (r *LedgerReplicationResponse) Pretty() string {
	count := len(r.Blocks)
	// error message or no blocks
	if r.Message != "" {
		return fmt.Sprintf("status: %s, message: %s, %v blocks, uuid: %s", r.Status.String(), r.Message, count, r.UUID.String())
	}
	blockInfo := ""
	if count > 0 {
		fb := r.Blocks[0]
		lb := r.Blocks[count-1]
		if fb != nil && lb != nil {
			fUC := &types.UnicityCertificate{}
			if err := fUC.UnmarshalCBOR(fb.UnicityCertificate); err != nil {
				blockInfo = " (error unmarshalling first block's UnicityCertificate)"
			} else {
				lUC := &types.UnicityCertificate{}
				if err := lUC.UnmarshalCBOR(lb.UnicityCertificate); err != nil {
					blockInfo = " (error unmarshalling last block's UnicityCertificate)"
				} else {
					blockInfo = fmt.Sprintf(" (round #%v (state %X) => #%v (state %X))", fUC.GetRoundNumber(), fUC.GetStateHash(), lUC.GetRoundNumber(), lUC.GetStateHash())
				}
			}
		}
	}
	return fmt.Sprintf("status: %s, %v blocks %s, uuid: %s", r.Status.String(), count, blockInfo, r.UUID.String())
}

func (r *LedgerReplicationResponse) IsValid() error {
	if r == nil {
		return ErrLedgerReplicationRespIsNil
	}
	if r.Status == Ok && r.Blocks == nil {
		return ErrLedgerResponseBlocksIsNil
	}
	return nil
}

func (r *LedgerReplicationRequest) IsValid() error {
	if r == nil {
		return ErrLedgerReplicationReqIsNil
	}
	if r.PartitionIdentifier == 0 {
		return ErrInvalidPartitionIdentifier
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
	case UnknownPartitionIdentifier:
		return "Unknown Partition Identifier"
	case Unknown:
		return "Unknown"
	}
	return "Unknown Status Code"
}
