package abdrc

import (
	"crypto"
	"fmt"
	"slices"

	abcrypto "github.com/alphabill-org/alphabill/crypto"
	drctypes "github.com/alphabill-org/alphabill/rootchain/consensus/abdrc/types"
	"github.com/alphabill-org/alphabill/types"
)

type GetStateMsg struct {
	_ struct{} `cbor:",toarray"`
	// ID of the node which requested the state, ie response should
	// be sent to that node
	NodeId string `json:"nodeId,omitempty"`
}

type InputData struct {
	_     struct{}           `cbor:",toarray"`
	SysID types.SystemID     `json:"sysID,omitempty"`
	Ir    *types.InputRecord `json:"ir,omitempty"`
	Sdrh  []byte             `json:"sdrh,omitempty"`
}

type CommittedBlock struct {
	_     struct{}            `cbor:",toarray"`
	Block *drctypes.BlockData `json:"block,omitempty"`
	Ir    []*InputData        `json:"ir,omitempty"`
}

type StateMsg struct {
	_             struct{}                    `cbor:",toarray"`
	Certificates  []*types.UnicityCertificate `json:"certificates,omitempty"`
	CommittedHead *CommittedBlock             `json:"committedHead,omitempty"`
	BlockData     []*drctypes.BlockData       `json:"blockData,omitempty"`
}

/*
CanRecoverToRound returns non-nil error when the state message is not suitable for recovery into round "round".
*/
func (sm *StateMsg) CanRecoverToRound(round uint64) error {
	if sm.CommittedHead == nil {
		return fmt.Errorf("committed block is nil")
	}
	if round < sm.CommittedHead.Block.GetRound() {
		return fmt.Errorf("can't recover to round %d with committed block for round %d", round, sm.CommittedHead.Block.GetRound())
	}
	// commit head matches recover round
	if sm.CommittedHead.Block.GetRound() == round {
		return nil
	}
	if !slices.ContainsFunc(sm.BlockData, func(b *drctypes.BlockData) bool { return b.GetRound() == round }) {
		return fmt.Errorf("state has no data block for round %d", round)
	}

	return nil
}

func (sm *StateMsg) Verify(hashAlgorithm crypto.Hash, quorum uint32, verifiers map[string]abcrypto.Verifier) error {
	if sm.CommittedHead == nil {
		return fmt.Errorf("commit head is nil")
	}
	if err := sm.CommittedHead.IsValid(); err != nil {
		return fmt.Errorf("invalid commit head block: %w", err)
	}
	if err := sm.CommittedHead.Block.Qc.Verify(quorum, verifiers); err != nil {
		return fmt.Errorf("commit head qc verification error: %w", err)
	}
	commitQcFound := false
	// verify node blocks
	for _, n := range sm.BlockData {
		if err := n.IsValid(); err != nil {
			return fmt.Errorf("invalid block node: %w", err)
		}
		if n.Qc != nil {
			if err := n.Qc.Verify(quorum, verifiers); err != nil {
				return fmt.Errorf("block node qc verification error: %w", err)
			}
			// check if this is the head block QC
			if n.Qc.LedgerCommitInfo != nil {
				if sm.CommittedHead.Block.GetRound() == n.Qc.LedgerCommitInfo.RootChainRoundNumber {
					commitQcFound = true
				}
			}
		}
	}
	if !commitQcFound {
		return fmt.Errorf("commit QC for head block not found")
	}
	for _, c := range sm.Certificates {
		if err := c.IsValid(verifiers, hashAlgorithm, c.UnicityTreeCertificate.SystemIdentifier, c.UnicityTreeCertificate.SystemDescriptionHash); err != nil {
			return fmt.Errorf("certificate for %X is invalid: %w", c.UnicityTreeCertificate.SystemIdentifier, err)
		}
	}
	return nil
}

func (r *CommittedBlock) GetRound() uint64 {
	if r != nil {
		return r.Block.GetRound()
	}
	return 0
}

func (r *CommittedBlock) IsValid() error {
	if len(r.Ir) == 0 {
		return fmt.Errorf("missing input record state")
	}
	for _, ir := range r.Ir {
		if err := ir.IsValid(); err != nil {
			return fmt.Errorf("invalid input record: %w", err)
		}
	}
	if r.Block == nil {
		return fmt.Errorf("block data is nil")
	}
	if err := r.Block.IsValid(); err != nil {
		return fmt.Errorf("block data error: %w", err)
	}
	return nil
}

func (i *InputData) IsValid() error {
	if i.Ir == nil {
		return fmt.Errorf("input record is nil")
	}
	if err := i.Ir.IsValid(); err != nil {
		return fmt.Errorf("input record error: %w", err)
	}
	if len(i.Sdrh) == 0 {
		return fmt.Errorf("system descrition hash not set")
	}
	return nil
}
