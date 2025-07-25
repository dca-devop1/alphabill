package storage

import (
	"crypto"
	"encoding/hex"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/alphabill-org/alphabill-go-base/cbor"
	"github.com/alphabill-org/alphabill-go-base/types"
	"github.com/alphabill-org/alphabill/internal/testutils/logger"
	"github.com/alphabill-org/alphabill/network/protocol/certification"
	drctypes "github.com/alphabill-org/alphabill/rootchain/consensus/types"
	"github.com/alphabill-org/alphabill/rootchain/testutils"
)

type mockIRVerifier struct {
	verify func(round uint64, irChReq *drctypes.IRChangeReq) (*types.InputRecord, error)
}

func (x mockIRVerifier) VerifyIRChangeReq(round uint64, irChReq *drctypes.IRChangeReq) (*types.InputRecord, error) {
	return x.verify(round, irChReq)
}

func TestNewGenesisBlock(t *testing.T) {
	b, err := NewGenesisBlock(5, crypto.SHA256)
	require.NoError(t, err)
	require.Equal(t, b.HashAlgo, crypto.SHA256)
	require.Empty(t, b.ShardState.States)
	require.Empty(t, b.ShardState.Changed)
	require.Nil(t, b.RootHash)
	require.NotNil(t, b.BlockData)
	require.Equal(t, uint64(1), b.BlockData.Round)
	require.Equal(t, "genesis", b.BlockData.Author)
	require.Nil(t, b.BlockData.Qc)
	require.NotNil(t, b.Qc)
	require.NoError(t, b.Qc.IsValid())
	require.NotNil(t, b.CommitQc)
	require.NoError(t, b.CommitQc.IsValid())
}

func TestExecutedBlock_Extend(t *testing.T) {
	_, shardNodeInfos := testutils.CreateTestNodes(t, 3)
	pdrEpoch1 := types.PartitionDescriptionRecord{
		PartitionID: 1,
		ShardID:     types.ShardID{},
		Epoch:       0,
		EpochStart:  1,
		Validators:  shardNodeInfos,
	}

	psID := types.PartitionShardID{PartitionID: pdrEpoch1.PartitionID, ShardID: pdrEpoch1.ShardID.Key()}
	orchestration := mockOrchestration{
		shardConfig: func(partitionID types.PartitionID, shardID types.ShardID, rootRound uint64) (*types.PartitionDescriptionRecord, error) {
			return &pdrEpoch1, nil
		},
		shardConfigs: func(rootRound uint64) (map[types.PartitionShardID]*types.PartitionDescriptionRecord, error) {
			return map[types.PartitionShardID]*types.PartitionDescriptionRecord{psID: &pdrEpoch1}, nil
		},
	}

	certReq := &certification.BlockCertificationRequest{
		PartitionID: pdrEpoch1.PartitionID,
		ShardID:     pdrEpoch1.ShardID,
		NodeID:      shardNodeInfos[0].NodeID,
		InputRecord: &types.InputRecord{
			Version:         1,
			PreviousHash:    []byte{1, 1, 1, 1},
			Hash:            []byte{2, 2, 2, 2},
			BlockHash:       []byte{3, 3, 3, 3},
			SummaryValue:    []byte{4, 4, 4, 4},
			RoundNumber:     4,
			SumOfEarnedFees: 3,
		},
	}
	newBlock := drctypes.BlockData{
		Author:    "test",
		Round:     drctypes.GenesisRootRound + 1,
		Epoch:     0,
		Timestamp: 12,
		Payload: &drctypes.Payload{
			Requests: []*drctypes.IRChangeReq{{
				Partition:  certReq.PartitionID,
				CertReason: drctypes.Quorum,
				Requests:   []*certification.BlockCertificationRequest{certReq},
			}},
		},
		Qc: nil, // not important in this context
	}

	reqVer := mockIRVerifier{verify: func(round uint64, irChReq *drctypes.IRChangeReq) (*types.InputRecord, error) {
		return irChReq.Requests[0].InputRecord, nil
	}}

	// current root block for tests to extend from. it's ok to extend from
	// the same block multiple times (mustn't affect the parent block)
	parent := genesisBlockWithShard(t, &pdrEpoch1)
	require.Len(t, parent.ShardState.States, 1)
	require.Contains(t, parent.ShardState.States, psID)
	require.Len(t, parent.ShardState.Changed, 1)
	require.Contains(t, parent.ShardState.Changed, psID)

	t.Run("orchestration error", func(t *testing.T) {
		expErr := errors.New("no configs")
		orc := mockOrchestration{
			shardConfigs: func(rootRound uint64) (map[types.PartitionShardID]*types.PartitionDescriptionRecord, error) {
				return nil, expErr
			},
		}
		executedBlock, err := parent.Extend(&newBlock, reqVer, orc, crypto.SHA256, logger.New(t))
		require.ErrorIs(t, err, expErr)
		require.Nil(t, executedBlock)
	})

	t.Run("invalid request", func(t *testing.T) {
		// verifying the request fails
		expErr := errors.New("invalid request")
		reqVer := mockIRVerifier{
			verify: func(round uint64, irChReq *drctypes.IRChangeReq) (*types.InputRecord, error) { return nil, expErr },
		}
		executedBlock, err := parent.Extend(&newBlock, reqVer, orchestration, crypto.SHA256, logger.New(t))
		require.ErrorIs(t, err, expErr)
		require.Nil(t, executedBlock)
	})

	t.Run("empty block", func(t *testing.T) {
		emptyBlock := newBlock
		emptyBlock.Payload = &drctypes.Payload{}
		executedBlock, err := parent.Extend(&emptyBlock, reqVer, orchestration, crypto.SHA512, logger.New(t))
		require.NoError(t, err)
		require.Equal(t, &emptyBlock, executedBlock.BlockData)
		require.Empty(t, executedBlock.ShardState.Changed, "expected no changes")
		require.Contains(t, executedBlock.ShardState.States, psID)
		require.Len(t, executedBlock.ShardState.States, 1)
		require.Equal(t, parent.ShardState.States, executedBlock.ShardState.States)
		require.Equal(t, crypto.SHA512, executedBlock.HashAlgo)
		require.Nil(t, executedBlock.Qc)
		require.Nil(t, executedBlock.CommitQc)
	})

	t.Run("non-empty block", func(t *testing.T) {
		executedBlock, err := parent.Extend(&newBlock, &reqVer, orchestration, crypto.SHA256, logger.New(t))
		require.NoError(t, err)
		require.Equal(t, &newBlock, executedBlock.BlockData)
		require.Len(t, executedBlock.ShardState.Changed, 1)
		require.Contains(t, executedBlock.ShardState.Changed, psID)
		require.Contains(t, executedBlock.ShardState.States, psID)
		require.Len(t, executedBlock.ShardState.States, 1)
		require.Equal(t, certReq.InputRecord, executedBlock.ShardState.States[psID].IR)
		require.Equal(t, crypto.SHA256, executedBlock.HashAlgo)
		// can't compare against hardcoded hash as fee hash and leader id change on each run (we generate partitionRecord)
		//require.EqualValues(t, "99AD3740E3CFC07EC1C1C04ED60D930BC3E2DC01AD5B3E8631C119C50EAF4520", fmt.Sprintf("%X", executedBlock.RootHash))
		// block has not got QC nor commit QC yet
		require.Nil(t, executedBlock.Qc)
		require.Nil(t, executedBlock.CommitQc)
	})

	t.Run("next epoch of a shard", func(t *testing.T) {
		// parent block is created with pdrEpoch1, now we return pdrEpoch2
		pdrEpoch2 := pdrEpoch1
		pdrEpoch2.Epoch++
		orchestration := mockOrchestration{
			shardConfigs: func(rootRound uint64) (map[types.PartitionShardID]*types.PartitionDescriptionRecord, error) {
				return map[types.PartitionShardID]*types.PartitionDescriptionRecord{psID: &pdrEpoch2}, nil
			},
		}
		// if shard has no ChangeRequest in the block then Epoch doesn't change!
		emptyBlock := newBlock
		emptyBlock.Payload = &drctypes.Payload{}
		executedBlock, err := parent.Extend(&emptyBlock, &reqVer, orchestration, crypto.SHA256, logger.New(t))
		require.NoError(t, err)
		require.Equal(t, &emptyBlock, executedBlock.BlockData)
		require.Empty(t, executedBlock.ShardState.Changed)
		require.Len(t, executedBlock.ShardState.States, 1)
		if assert.Contains(t, executedBlock.ShardState.States, psID) {
			si := executedBlock.ShardState.States[psID]
			require.Equal(t, pdrEpoch1.Epoch, si.TR.Epoch, "epoch should stay the same")
		}

		// next block with shard sending ChangeRequest - the TR should now indicate next epoch
		executedBlock, err = executedBlock.Extend(&newBlock, &reqVer, orchestration, crypto.SHA256, logger.New(t))
		require.NoError(t, err)
		if assert.Contains(t, executedBlock.ShardState.States, psID) {
			si := executedBlock.ShardState.States[psID]
			require.Equal(t, pdrEpoch2.Epoch, si.TR.Epoch, "signal new epoch in the TR")
		}
		require.Equal(t, &newBlock, executedBlock.BlockData)
		require.Len(t, executedBlock.ShardState.Changed, 1)
		require.Contains(t, executedBlock.ShardState.Changed, psID)
		require.Len(t, executedBlock.ShardState.States, 1)
		require.Contains(t, executedBlock.ShardState.States, psID)
		require.Equal(t, certReq.InputRecord, executedBlock.ShardState.States[psID].IR)
	})

	t.Run("new shard introduced", func(t *testing.T) {
		pdrPart2 := pdrEpoch1
		pdrPart2.PartitionID++
		psID2 := types.PartitionShardID{PartitionID: pdrPart2.PartitionID, ShardID: pdrPart2.ShardID.Key()}
		orchestration := mockOrchestration{
			shardConfigs: func(rootRound uint64) (map[types.PartitionShardID]*types.PartitionDescriptionRecord, error) {
				return map[types.PartitionShardID]*types.PartitionDescriptionRecord{psID: &pdrEpoch1, psID2: &pdrPart2}, nil
			},
		}

		block := newBlock
		block.Payload = &drctypes.Payload{
			Requests: []*drctypes.IRChangeReq{{
				Partition:  certReq.PartitionID,
				CertReason: drctypes.Quorum,
				Requests:   []*certification.BlockCertificationRequest{certReq},
			}},
		}
		executedBlock, err := parent.Extend(&block, &reqVer, orchestration, crypto.SHA256, logger.New(t))
		require.NoError(t, err)
		require.Equal(t, &block, executedBlock.BlockData)
		require.Len(t, executedBlock.ShardState.Changed, 2)
		require.Contains(t, executedBlock.ShardState.Changed, psID)
		require.Contains(t, executedBlock.ShardState.Changed, psID2)
		require.Len(t, executedBlock.ShardState.States, 2)
		require.Contains(t, executedBlock.ShardState.States, psID)
		require.Contains(t, executedBlock.ShardState.States, psID2)
	})

	t.Run("shard removed", func(t *testing.T) {
		// split shard - two new added and original removed
		pdrPartA := pdrEpoch1
		pdrPartB := pdrEpoch1
		pdrPartA.ShardID, pdrPartB.ShardID = pdrEpoch1.ShardID.Split()
		psIDA := types.PartitionShardID{PartitionID: pdrPartA.PartitionID, ShardID: pdrPartA.ShardID.Key()}
		psIDB := types.PartitionShardID{PartitionID: pdrPartB.PartitionID, ShardID: pdrPartB.ShardID.Key()}
		orchestration := mockOrchestration{
			shardConfigs: func(rootRound uint64) (map[types.PartitionShardID]*types.PartitionDescriptionRecord, error) {
				return map[types.PartitionShardID]*types.PartitionDescriptionRecord{psIDA: &pdrPartA, psIDB: &pdrPartB}, nil
			},
		}
		emptyBlock := newBlock
		emptyBlock.Payload = &drctypes.Payload{}
		executedBlock, err := parent.Extend(&emptyBlock, &reqVer, orchestration, crypto.SHA256, logger.New(t))
		require.NoError(t, err)
		require.Len(t, executedBlock.ShardState.Changed, 2)
		require.Contains(t, executedBlock.ShardState.Changed, psIDA)
		require.Contains(t, executedBlock.ShardState.Changed, psIDB)
		require.Len(t, executedBlock.ShardState.States, 2)
		require.Contains(t, executedBlock.ShardState.States, psIDA)
		require.Contains(t, executedBlock.ShardState.States, psIDB)
	})
}

func TestExecutedBlock_GenerateCertificates(t *testing.T) {
	const partitionID1 types.PartitionID = 1
	const partitionID2 types.PartitionID = 2
	rh, err := hex.DecodeString("51592107828763663BE3378AD1F4BAE7D9C1A921DEEC1A6B28247770A8B4F526")
	require.NoError(t, err)

	validBlock := func() *ExecutedBlock {
		return &ExecutedBlock{
			BlockData: &drctypes.BlockData{
				Author:  "test",
				Round:   2,
				Payload: &drctypes.Payload{},
				Qc:      nil,
			},
			ShardState: ShardStates{
				States: map[types.PartitionShardID]*ShardInfo{
					{PartitionID: partitionID1, ShardID: types.ShardID{}.Key()}: {
						PartitionID: partitionID1,
						IR: &types.InputRecord{
							Version:         1,
							PreviousHash:    []byte{1, 1, 1, 1},
							Hash:            []byte{2, 2, 2, 2},
							BlockHash:       []byte{3, 3, 3, 3},
							SummaryValue:    []byte{4, 4, 4, 4},
							ETHash:          []byte{5, 5, 5, 5},
							RoundNumber:     3,
							SumOfEarnedFees: 4,
							Timestamp:       20241113,
						},
						ShardConfHash: []byte{1, 2, 3, 4},
					},
					{PartitionID: partitionID2, ShardID: types.ShardID{}.Key()}: {
						PartitionID: partitionID2,
						IR: &types.InputRecord{
							Version:         1,
							PreviousHash:    []byte{1, 1, 1, 1},
							Hash:            []byte{4, 4, 4, 4},
							BlockHash:       []byte{3, 3, 3, 3},
							SummaryValue:    []byte{4, 4, 4, 4},
							ETHash:          []byte{5, 5, 5, 5},
							RoundNumber:     3,
							SumOfEarnedFees: 6,
							Timestamp:       20241113,
						},
						ShardConfHash: []byte{4, 5, 6, 7},
					},
				},
				Changed: ShardSet{
					{PartitionID: partitionID1, ShardID: types.ShardID{}.Key()}: {},
					{PartitionID: partitionID2, ShardID: types.ShardID{}.Key()}: {},
				},
			},
			HashAlgo: crypto.SHA256,
			RootHash: rh,
		}
	}

	validCommitQc := func() *drctypes.QuorumCert {
		return &drctypes.QuorumCert{
			VoteInfo: &drctypes.RoundInfo{
				RoundNumber:       3,
				ParentRoundNumber: 2,
				CurrentRootHash:   make([]byte, crypto.SHA256.Size()),
			},
			LedgerCommitInfo: &types.UnicitySeal{
				Version:      1,
				PreviousHash: []byte{0, 0, 0, 0},
				Hash:         rh,
			},
		}
	}

	t.Run("invalid state", func(t *testing.T) {
		// block is in invalid state - the schema and shard info do not match
		commitQc := validCommitQc()
		block := validBlock()
		// scheme lists partition for which there is no shard info
		block.ShardState.schemes = map[types.PartitionID]types.ShardingScheme{66: {}}
		certs, err := block.GenerateCertificates(commitQc)
		require.EqualError(t, err, `failed to generate root hash: creating unicity tree: missing shard info for 00000042_80`)
		require.Empty(t, certs)
	})

	t.Run("root hash of the block differs", func(t *testing.T) {
		commitQc := validCommitQc()
		block := validBlock()
		block.RootHash = []byte{1}
		certs, err := block.GenerateCertificates(commitQc)
		require.EqualError(t, err, "root hash does not match previously calculated root hash")
		require.Nil(t, certs)
	})

	t.Run("root hash of the commitQc does not match", func(t *testing.T) {
		commitQc := validCommitQc()
		commitQc.LedgerCommitInfo.Hash = []byte{2}
		block := validBlock()
		certs, err := block.GenerateCertificates(commitQc)
		require.EqualError(t, err, "root hash does not match hash in commit QC")
		require.Nil(t, certs)
	})

	t.Run("success, no changes", func(t *testing.T) {
		commitQc := validCommitQc()
		block := validBlock()
		block.ShardState.Changed = nil
		certs, err := block.GenerateCertificates(commitQc)
		require.NoError(t, err)
		require.Empty(t, certs)
	})

	t.Run("success with changes", func(t *testing.T) {
		commitQc := validCommitQc()
		block := validBlock()
		certs, err := block.GenerateCertificates(commitQc)
		require.NoError(t, err)
		require.Len(t, certs, 2)
		si, ok := block.ShardState.States[types.PartitionShardID{PartitionID: partitionID1, ShardID: types.ShardID{}.Key()}]
		require.True(t, ok)
		require.NotNil(t, si.LastCR)
	})
}

func TestExecutedBlock_GetRound(t *testing.T) {
	var b *ExecutedBlock
	require.Equal(t, uint64(0), b.GetRound())
	b = &ExecutedBlock{BlockData: nil}
	require.Equal(t, uint64(0), b.GetRound())
	b = &ExecutedBlock{BlockData: &drctypes.BlockData{Round: 2}}
	require.Equal(t, uint64(2), b.GetRound())
}

func TestExecutedBlock_GetParentRound(t *testing.T) {
	var b *ExecutedBlock
	require.Equal(t, uint64(0), b.GetParentRound())
	b = &ExecutedBlock{BlockData: &drctypes.BlockData{}}
	require.Equal(t, uint64(0), b.GetParentRound())
	b = &ExecutedBlock{BlockData: &drctypes.BlockData{Qc: &drctypes.QuorumCert{}}}
	require.Equal(t, uint64(0), b.GetParentRound())
	b = &ExecutedBlock{BlockData: &drctypes.BlockData{Qc: &drctypes.QuorumCert{VoteInfo: &drctypes.RoundInfo{}}}}
	require.Equal(t, uint64(0), b.GetParentRound())
	b = &ExecutedBlock{BlockData: &drctypes.BlockData{Qc: &drctypes.QuorumCert{VoteInfo: &drctypes.RoundInfo{RoundNumber: 2}}}}
	require.Equal(t, uint64(2), b.GetParentRound())
}

func Test_ExecutedBlock_serialization(t *testing.T) {
	t.Run("Changed set", func(t *testing.T) {
		// empty set
		// we init the Changed manually to non-nil value as require.EqualValues
		// considers nil and empty map as different. In code the ExecutedBlock
		// values are constructed via constructors which init the Changed field.
		b1 := ExecutedBlock{ShardState: ShardStates{Changed: ShardSet{}}}
		buf, err := cbor.Marshal(b1)
		require.NoError(t, err)

		var b2 ExecutedBlock
		require.NoError(t, cbor.Unmarshal(buf, &b2))
		require.EqualValues(t, b1.ShardState.Changed, b2.ShardState.Changed)

		// set with one item (empty shard ID)
		b1.ShardState.Changed = map[types.PartitionShardID]struct{}{{PartitionID: 1, ShardID: types.ShardID{}.Key()}: {}}
		buf, err = cbor.Marshal(b1)
		require.NoError(t, err)

		require.NoError(t, cbor.Unmarshal(buf, &b2))
		require.EqualValues(t, b1.ShardState.Changed, b2.ShardState.Changed)

		// set with two shards
		s0, s1 := types.ShardID{}.Split()
		b1.ShardState.Changed = map[types.PartitionShardID]struct{}{
			{PartitionID: 2, ShardID: s0.Key()}: {},
			{PartitionID: 2, ShardID: s1.Key()}: {},
		}
		buf, err = cbor.Marshal(b1)
		require.NoError(t, err)

		require.NoError(t, cbor.Unmarshal(buf, &b2))
		require.EqualValues(t, b1.ShardState.Changed, b2.ShardState.Changed)
	})

	t.Run("ShardInfo", func(t *testing.T) {
		// empty map
		b1 := ExecutedBlock{ShardState: ShardStates{States: map[types.PartitionShardID]*ShardInfo{}}}
		buf, err := cbor.Marshal(b1)
		require.NoError(t, err)

		var b2 ExecutedBlock
		require.NoError(t, cbor.Unmarshal(buf, &b2))
		require.EqualValues(t, b1.ShardState.States, b2.ShardState.States)

		// non-empty map
		si := ShardInfo{
			PartitionID:   9,
			ShardID:       types.ShardID{},
			RootHash:      []byte{3, 3, 3},
			PrevEpochStat: []byte{0x43, 4, 4, 4}, // array(3)
			PrevEpochFees: []byte{0x43, 5, 5, 5},
			Fees:          map[string]uint64{"A": 10},
			LastCR: &certification.CertificationResponse{
				Partition: 9,
				Shard:     types.ShardID{},
				Technical: certification.TechnicalRecord{
					Round:  2,
					Epoch:  3,
					Leader: "ldr",
				},
				UC: types.UnicityCertificate{
					Version: 1,
				},
			},
		}
		psKey := types.PartitionShardID{PartitionID: si.LastCR.Partition, ShardID: si.LastCR.Shard.Key()}
		b1.ShardState.States[psKey] = &si
		buf, err = cbor.Marshal(b1)
		require.NoError(t, err)

		require.NoError(t, cbor.Unmarshal(buf, &b2))
		require.Equal(t, b1.ShardState.States, b2.ShardState.States)
		require.Equal(t, &si, b2.ShardState.States[psKey])
	})
}
