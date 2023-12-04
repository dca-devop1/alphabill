package abdrc

import (
	"bytes"
	"context"
	gocrypto "crypto"
	"crypto/rand"
	"fmt"
	"slices"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/alphabill-org/alphabill/internal/keyvaluedb/memorydb"
	"github.com/alphabill-org/alphabill/internal/rootchain/consensus/abdrc/storage"
	p2pcrypto "github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
	p2ptest "github.com/libp2p/go-libp2p/core/test"
	"github.com/stretchr/testify/require"

	"github.com/alphabill-org/alphabill/internal/crypto"
	"github.com/alphabill-org/alphabill/internal/network"
	"github.com/alphabill-org/alphabill/internal/network/protocol/abdrc"
	"github.com/alphabill-org/alphabill/internal/network/protocol/certification"
	"github.com/alphabill-org/alphabill/internal/network/protocol/genesis"
	"github.com/alphabill-org/alphabill/internal/rootchain/consensus"
	abdrctu "github.com/alphabill-org/alphabill/internal/rootchain/consensus/abdrc/testutils"
	abtypes "github.com/alphabill-org/alphabill/internal/rootchain/consensus/abdrc/types"
	rootgenesis "github.com/alphabill-org/alphabill/internal/rootchain/genesis"
	"github.com/alphabill-org/alphabill/internal/rootchain/partitions"
	"github.com/alphabill-org/alphabill/internal/rootchain/testutils"
	test "github.com/alphabill-org/alphabill/internal/testutils"
	testlogger "github.com/alphabill-org/alphabill/internal/testutils/logger"
	testnetwork "github.com/alphabill-org/alphabill/internal/testutils/network"
	"github.com/alphabill-org/alphabill/internal/testutils/observability"
	"github.com/alphabill-org/alphabill/internal/types"
	"github.com/alphabill-org/alphabill/pkg/logger"
)

var partitionID types.SystemID32 = 0x00FF0001
var partitionInputRecord = &types.InputRecord{
	PreviousHash: make([]byte, 32),
	Hash:         []byte{0, 0, 0, 1},
	BlockHash:    []byte{0, 0, 1, 2},
	SummaryValue: []byte{0, 0, 1, 3},
	RoundNumber:  1,
}

func readResult(ch <-chan *types.UnicityCertificate, timeout time.Duration) (*types.UnicityCertificate, error) {
	select {
	case result, ok := <-ch:
		if !ok {
			return nil, fmt.Errorf("failed to read from channel")
		}
		return result, nil
	case <-time.After(timeout):
		return nil, fmt.Errorf("timeout")
	}
}

func initConsensusManager(t *testing.T, net RootNet) (*ConsensusManager, *testutils.TestNode, []*testutils.TestNode, *genesis.RootGenesis) {
	partitionNodes, partitionRecord := testutils.CreatePartitionNodesAndPartitionRecord(t, partitionInputRecord, partitionID, 3)
	rootNode := testutils.NewTestNode(t)
	verifier := rootNode.Verifier
	rootPubKeyBytes, err := verifier.MarshalPublicKey()
	require.NoError(t, err)
	id := rootNode.PeerConf.ID
	rootGenesis, _, err := rootgenesis.NewRootGenesis(id.String(), rootNode.Signer, rootPubKeyBytes, []*genesis.PartitionRecord{partitionRecord})
	require.NoError(t, err)
	partitions, err := partitions.NewPartitionStoreFromGenesis(rootGenesis.Partitions)
	require.NoError(t, err)
	cm, err := NewDistributedAbConsensusManager(id, rootGenesis, partitions, net, rootNode.Signer, observability.NOPMetrics(), testlogger.New(t).With(logger.NodeID(id)))
	require.NoError(t, err)
	return cm, rootNode, partitionNodes, rootGenesis
}

func buildBlockCertificationRequest(t *testing.T, rg *genesis.RootGenesis, partitionNodes []*testutils.TestNode) []*certification.BlockCertificationRequest {
	t.Helper()
	newIR := &types.InputRecord{
		PreviousHash: rg.Partitions[0].Nodes[0].BlockCertificationRequest.InputRecord.Hash,
		Hash:         test.RandomBytes(32),
		BlockHash:    test.RandomBytes(32),
		SummaryValue: rg.Partitions[0].Nodes[0].BlockCertificationRequest.InputRecord.SummaryValue,
		RoundNumber:  2,
	}
	requests := make([]*certification.BlockCertificationRequest, len(partitionNodes))
	for i, n := range partitionNodes {
		requests[i] = testutils.CreateBlockCertificationRequest(t, newIR, partitionID, n)
	}
	return requests
}

func TestNewConsensusManager_Ok(t *testing.T) {
	mockNet := testnetwork.NewRootMockNetwork()
	cm, root, partitionNodes, rg := initConsensusManager(t, mockNet)
	ctx, ctxCancel := context.WithCancel(context.Background())
	defer ctxCancel()
	go func() { require.ErrorIs(t, cm.Run(ctx), context.Canceled) }()

	require.Len(t, partitionNodes, 3)
	require.NotNil(t, cm)
	require.NotNil(t, root)
	require.NotNil(t, rg)
}

func Test_ConsensusManager_onPartitionIRChangeReq(t *testing.T) {
	mockNet := testnetwork.NewRootMockNetwork()
	cm, _, partitionNodes, rg := initConsensusManager(t, mockNet)

	req := &consensus.IRChangeRequest{
		SystemIdentifier: partitionID,
		Reason:           consensus.Quorum,
		Requests:         buildBlockCertificationRequest(t, rg, partitionNodes),
	}

	// we need to init pacemaker into correct round, otherwise IR validation fails
	cm.pacemaker.Reset(cm.blockStore.GetHighQc().VoteInfo.RoundNumber, nil, nil)
	defer cm.pacemaker.Stop()

	require.NoError(t, cm.onPartitionIRChangeReq(context.Background(), req))
	// since there is only one root node, it is the next leader, the request will be buffered
	require.True(t, cm.irReqBuffer.IsChangeInBuffer(partitionID))
}

func Test_ConsensusManager_onIRChangeMsg_ErrInvalidSignature(t *testing.T) {
	mockNet := testnetwork.NewRootMockNetwork()
	cm, _, partitionNodes, rg := initConsensusManager(t, mockNet)

	req := &abdrc.IrChangeReqMsg{
		Author: cm.id.String(),
		IrChangeReq: &abtypes.IRChangeReq{
			SystemIdentifier: partitionID,
			CertReason:       abtypes.Quorum,
			Requests:         buildBlockCertificationRequest(t, rg, partitionNodes),
		},
		Signature: []byte{1, 2, 3, 4},
	}
	// verify that error is printed and author ID is also present
	require.ErrorContains(t, cm.onIRChangeMsg(context.Background(), req),
		fmt.Sprintf("invalid IR change request message from node %s: signature verification failed", cm.id.String()))
}

func TestIRChangeRequestFromRootValidator_RootTimeoutOnFirstRound(t *testing.T) {
	var lastProposalMsg *abdrc.ProposalMsg = nil
	var lastVoteMsg *abdrc.VoteMsg = nil
	var lastTimeoutMsg *abdrc.TimeoutMsg = nil

	mockNet := testnetwork.NewRootMockNetwork()
	cm, rootNode, _, _ := initConsensusManager(t, mockNet)
	ctx, ctxCancel := context.WithCancel(context.Background())
	defer ctxCancel()
	go func() { require.ErrorIs(t, cm.Run(ctx), context.Canceled) }()

	// Await proposal
	lastProposalMsg = testutils.MockAwaitMessage[*abdrc.ProposalMsg](t, mockNet, network.ProtocolRootProposal)
	require.Equal(t, uint64(2), lastProposalMsg.Block.Round)
	// Quick hack to trigger timeout
	// simulate local timeout by calling the method -> race/hack accessing from different go routines not safe
	cm.onLocalTimeout(ctx)
	// await timeout vote
	lastTimeoutMsg = testutils.MockAwaitMessage[*abdrc.TimeoutMsg](t, mockNet, network.ProtocolRootTimeout)
	require.Equal(t, uint64(2), lastTimeoutMsg.Timeout.Round)
	// simulate TC not achieved and make sure the same timeout message is sent again
	// Quick hack to trigger next timeout
	// simulate local timeout by calling the method -> race/hack accessing from different go routines not safe
	cm.onLocalTimeout(ctx)
	lastTimeoutMsg = testutils.MockAwaitMessage[*abdrc.TimeoutMsg](t, mockNet, network.ProtocolRootTimeout)
	require.Equal(t, uint64(2), lastTimeoutMsg.Timeout.Round)
	// route timeout message back
	// route the timeout message back to trigger timeout certificate and new round
	testutils.MockValidatorNetReceives(t, mockNet, rootNode.PeerConf.ID, network.ProtocolRootTimeout, lastTimeoutMsg)
	// This triggers TC and next round, wait for proposal
	lastProposalMsg = testutils.MockAwaitMessage[*abdrc.ProposalMsg](t, mockNet, network.ProtocolRootProposal)
	require.True(t, lastProposalMsg.Block.Payload.IsEmpty())
	require.Equal(t, uint64(3), lastProposalMsg.Block.Round)
	require.NotNil(t, lastProposalMsg.LastRoundTc)
	require.Equal(t, uint64(2), lastProposalMsg.LastRoundTc.Timeout.Round)
	// route the proposal back
	testutils.MockValidatorNetReceives(t, mockNet, rootNode.PeerConf.ID, network.ProtocolRootProposal, lastProposalMsg)
	// wait for the vote message
	lastVoteMsg = testutils.MockAwaitMessage[*abdrc.VoteMsg](t, mockNet, network.ProtocolRootVote)
	require.Equal(t, uint64(3), lastVoteMsg.VoteInfo.RoundNumber)
	// round 3 is skipped, as it timeouts
	require.Equal(t, uint64(1), lastVoteMsg.VoteInfo.ParentRoundNumber)
	require.Equal(t, uint64(0), lastVoteMsg.VoteInfo.Epoch)
	require.Nil(t, lastVoteMsg.LedgerCommitInfo.Hash)
}

func TestIRChangeRequestFromRootValidator_RootTimeout(t *testing.T) {
	var lastProposalMsg *abdrc.ProposalMsg = nil
	var lastVoteMsg *abdrc.VoteMsg = nil
	var lastTimeoutMsg *abdrc.TimeoutMsg = nil

	mockNet := testnetwork.NewRootMockNetwork()
	cm, rootNode, partitionNodes, rg := initConsensusManager(t, mockNet)
	ctx, ctxCancel := context.WithCancel(context.Background())
	defer ctxCancel()
	go func() { require.ErrorIs(t, cm.Run(ctx), context.Canceled) }()

	// simulate IR change request message
	irChReqMsg := &abdrc.IrChangeReqMsg{
		Author: rootNode.PeerConf.ID.String(),
		IrChangeReq: &abtypes.IRChangeReq{
			SystemIdentifier: partitionID,
			CertReason:       abtypes.Quorum,
			Requests:         buildBlockCertificationRequest(t, rg, partitionNodes[0:2]),
		},
	}
	require.NoError(t, irChReqMsg.Sign(rootNode.Signer))
	testutils.MockValidatorNetReceives(t, mockNet, rootNode.PeerConf.ID, network.ProtocolRootIrChangeReq, irChReqMsg)
	// As the node is the leader, next round will trigger a proposal
	lastProposalMsg = testutils.MockAwaitMessage[*abdrc.ProposalMsg](t, mockNet, network.ProtocolRootProposal)
	require.Equal(t, partitionID, lastProposalMsg.Block.Payload.Requests[0].SystemIdentifier)
	require.Equal(t, abtypes.Quorum, lastProposalMsg.Block.Payload.Requests[0].CertReason)
	require.Len(t, lastProposalMsg.Block.Payload.Requests[0].Requests, 2)
	// route the proposal back
	testutils.MockValidatorNetReceives(t, mockNet, rootNode.PeerConf.ID, network.ProtocolRootProposal, lastProposalMsg)
	// wait for the vote message
	lastVoteMsg = testutils.MockAwaitMessage[*abdrc.VoteMsg](t, mockNet, network.ProtocolRootVote)
	require.Equal(t, uint64(2), lastVoteMsg.VoteInfo.RoundNumber)
	require.Equal(t, uint64(1), lastVoteMsg.VoteInfo.ParentRoundNumber)
	require.Equal(t, uint64(0), lastVoteMsg.VoteInfo.Epoch)
	require.NotNil(t, lastVoteMsg.LedgerCommitInfo.Hash)

	// send vote back to validator
	testutils.MockValidatorNetReceives(t, mockNet, rootNode.PeerConf.ID, network.ProtocolRootVote, lastVoteMsg)
	// this will trigger next proposal since QC is achieved
	lastProposalMsg = testutils.MockAwaitMessage[*abdrc.ProposalMsg](t, mockNet, network.ProtocolRootProposal)
	// no additional requests have been received, meaning payload is empty
	require.True(t, lastProposalMsg.Block.Payload.IsEmpty())

	// route the proposal back to trigger new vote
	testutils.MockValidatorNetReceives(t, mockNet, rootNode.PeerConf.ID, network.ProtocolRootProposal, lastProposalMsg)
	// wait for the vote message
	lastVoteMsg = testutils.MockAwaitMessage[*abdrc.VoteMsg](t, mockNet, network.ProtocolRootVote)
	require.Equal(t, uint64(3), lastVoteMsg.VoteInfo.RoundNumber)
	require.NotNil(t, lastVoteMsg.LedgerCommitInfo.Hash)

	// Do not route the vote back, instead simulate round/view timeout
	// simulate local timeout by calling the method -> race/hack accessing from different go routines not safe
	cm.onLocalTimeout(ctx)
	// await timeout vote
	lastTimeoutMsg = testutils.MockAwaitMessage[*abdrc.TimeoutMsg](t, mockNet, network.ProtocolRootTimeout)
	require.Equal(t, uint64(3), lastTimeoutMsg.Timeout.Round)
	// route the timeout message back to trigger timeout certificate and new round
	testutils.MockValidatorNetReceives(t, mockNet, rootNode.PeerConf.ID, network.ProtocolRootTimeout, lastTimeoutMsg)
	// this will immediately trigger timeout certificate for the round
	// the following must be true now:
	// round is advanced
	require.Equal(t, uint64(4), cm.pacemaker.GetCurrentRound())
	// only changes from round 3 are removed, rest will still be active
	require.True(t, cm.blockStore.IsChangeInProgress(partitionID))
	// await the next proposal as well, the proposal must contain TC
	lastProposalMsg = testutils.MockAwaitMessage[*abdrc.ProposalMsg](t, mockNet, network.ProtocolRootProposal)
	require.NotNil(t, lastProposalMsg.LastRoundTc)
	require.Equal(t, uint64(3), lastProposalMsg.LastRoundTc.Timeout.Round)
	// query state
	getStateMsg := &abdrc.GetStateMsg{
		NodeId: partitionNodes[0].PeerConf.ID.String(),
	}
	// no change requests added, previous changes still not committed as timeout occurred
	require.True(t, lastProposalMsg.Block.Payload.IsEmpty())
	require.Equal(t, uint64(4), lastProposalMsg.Block.Round)
	// route the proposal back
	testutils.MockValidatorNetReceives(t, mockNet, rootNode.PeerConf.ID, network.ProtocolRootProposal, lastProposalMsg)
	// await vote
	lastVoteMsg = testutils.MockAwaitMessage[*abdrc.VoteMsg](t, mockNet, network.ProtocolRootVote)
	require.Equal(t, uint64(4), lastVoteMsg.VoteInfo.RoundNumber)
	require.Nil(t, lastVoteMsg.LedgerCommitInfo.Hash)
	// Check state before routing vote back to root
	testutils.MockValidatorNetReceives(t, mockNet, rootNode.PeerConf.ID, network.ProtocolRootStateReq, getStateMsg)
	stateMsg := testutils.MockAwaitMessage[*abdrc.StateMsg](t, mockNet, network.ProtocolRootStateResp)
	// commit head is still at round 1, as round 3 that would have committed 2 resulted in timeout
	require.Equal(t, uint64(1), stateMsg.CommittedHead.Block.Round)
	require.Equal(t, 2, len(stateMsg.BlockNode))
	// round 3 has been removed as it resulted in timeout quorum
	require.Equal(t, uint64(2), stateMsg.BlockNode[0].Block.Round)
	require.Equal(t, uint64(4), stateMsg.BlockNode[1].Block.Round)
	// send vote back to validator
	testutils.MockValidatorNetReceives(t, mockNet, rootNode.PeerConf.ID, network.ProtocolRootVote, lastVoteMsg)

	// await proposal again
	lastProposalMsg = testutils.MockAwaitMessage[*abdrc.ProposalMsg](t, mockNet, network.ProtocolRootProposal)
	require.Nil(t, lastProposalMsg.LastRoundTc)
	require.Equal(t, uint64(5), lastProposalMsg.Block.Round)
	// route the proposal back
	testutils.MockValidatorNetReceives(t, mockNet, rootNode.PeerConf.ID, network.ProtocolRootProposal, lastProposalMsg)

	lastVoteMsg = testutils.MockAwaitMessage[*abdrc.VoteMsg](t, mockNet, network.ProtocolRootVote)
	require.Equal(t, uint64(5), lastVoteMsg.VoteInfo.RoundNumber)
	require.NotNil(t, lastVoteMsg.LedgerCommitInfo.Hash)
	// Check state before routing vote back to root
	testutils.MockValidatorNetReceives(t, mockNet, rootNode.PeerConf.ID, network.ProtocolRootStateReq, getStateMsg)
	stateMsg = testutils.MockAwaitMessage[*abdrc.StateMsg](t, mockNet, network.ProtocolRootStateResp)
	// commit head is still at round 1, rounds 2, 4 and 5 are added, 5 will commit 4 when it reaches quorum, but
	// this will after vote is routed back, so current expected state is:
	require.Equal(t, uint64(1), stateMsg.CommittedHead.Block.Round)
	require.Equal(t, 3, len(stateMsg.BlockNode))
	// round 3 has been removed as it resulted in timeout quorum
	require.Equal(t, uint64(2), stateMsg.BlockNode[0].Block.Round)
	require.Equal(t, uint64(4), stateMsg.BlockNode[1].Block.Round)
	require.Equal(t, uint64(5), stateMsg.BlockNode[2].Block.Round)

	// send vote back to validator
	testutils.MockValidatorNetReceives(t, mockNet, rootNode.PeerConf.ID, network.ProtocolRootVote, lastVoteMsg)
	// after two successful rounds since timeout the IR change will be finally committed and UC is returned
	result, err := readResult(cm.CertificationResult(), time.Second)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, partitionID.ToSystemID(), result.UnicityTreeCertificate.SystemIdentifier)
	require.False(t, cm.blockStore.IsChangeInProgress(partitionID))
	// verify certificates have been updated when recovery query is sent
	getCertsMsg := &abdrc.GetStateMsg{
		NodeId: partitionNodes[0].PeerConf.ID.String(),
	}
	// simulate IR change request message
	testutils.MockValidatorNetReceives(t, mockNet, rootNode.PeerConf.ID, network.ProtocolRootStateReq, getCertsMsg)
	// As the node is the leader, next round will trigger a proposal
	certsMsg := testutils.MockAwaitMessage[*abdrc.StateMsg](t, mockNet, network.ProtocolRootStateResp)
	require.Equal(t, len(rg.Partitions), len(certsMsg.Certificates))
	idx := slices.IndexFunc(certsMsg.Certificates, func(c *types.UnicityCertificate) bool {
		return bytes.Equal(c.UnicityTreeCertificate.SystemIdentifier, partitionID.ToSystemID())
	})
	require.False(t, idx == -1)
	require.True(t, certsMsg.Certificates[idx].UnicitySeal.RootChainRoundNumber > uint64(1))
	// simulate IR change request message
	testutils.MockValidatorNetReceives(t, mockNet, rootNode.PeerConf.ID, network.ProtocolRootStateReq, getStateMsg)
	// As the node is the leader, next round will trigger a proposal
	stateMsg = testutils.MockAwaitMessage[*abdrc.StateMsg](t, mockNet, network.ProtocolRootStateResp)
	// at this stage the committed round is 4 and round 5 block is pending, if it reaches quorum it will commit 4
	require.Equal(t, uint64(4), stateMsg.CommittedHead.Block.Round)
	require.Equal(t, 1, len(stateMsg.BlockNode))
	require.Equal(t, uint64(5), stateMsg.BlockNode[0].Block.Round)
}

func TestIRChangeRequestFromRootValidator(t *testing.T) {
	var lastProposalMsg *abdrc.ProposalMsg = nil
	var lastVoteMsg *abdrc.VoteMsg = nil

	mockNet := testnetwork.NewRootMockNetwork()
	cm, rootNode, partitionNodes, rg := initConsensusManager(t, mockNet)
	ctx, ctxCancel := context.WithCancel(context.Background())
	defer ctxCancel()
	go func() { require.ErrorIs(t, cm.Run(ctx), context.Canceled) }()

	// simulate IR change request message
	irChReqMsg := &abdrc.IrChangeReqMsg{
		Author: rootNode.PeerConf.ID.String(),
		IrChangeReq: &abtypes.IRChangeReq{
			SystemIdentifier: partitionID,
			CertReason:       abtypes.Quorum,
			Requests:         buildBlockCertificationRequest(t, rg, partitionNodes[0:2]),
		},
	}
	require.NoError(t, irChReqMsg.Sign(rootNode.Signer))
	testutils.MockValidatorNetReceives(t, mockNet, rootNode.PeerConf.ID, network.ProtocolRootIrChangeReq, irChReqMsg)
	// As the node is the leader, next round will trigger a proposal
	lastProposalMsg = testutils.MockAwaitMessage[*abdrc.ProposalMsg](t, mockNet, network.ProtocolRootProposal)
	require.Equal(t, partitionID, lastProposalMsg.Block.Payload.Requests[0].SystemIdentifier)
	require.Equal(t, abtypes.Quorum, lastProposalMsg.Block.Payload.Requests[0].CertReason)
	require.Len(t, lastProposalMsg.Block.Payload.Requests[0].Requests, 2)
	// route the proposal back
	testutils.MockValidatorNetReceives(t, mockNet, rootNode.PeerConf.ID, network.ProtocolRootProposal, lastProposalMsg)
	// wait for the vote message
	lastVoteMsg = testutils.MockAwaitMessage[*abdrc.VoteMsg](t, mockNet, network.ProtocolRootVote)
	require.Equal(t, uint64(2), lastVoteMsg.VoteInfo.RoundNumber)
	require.Equal(t, uint64(1), lastVoteMsg.VoteInfo.ParentRoundNumber)
	require.Equal(t, uint64(0), lastVoteMsg.VoteInfo.Epoch)
	require.NotNil(t, lastVoteMsg.LedgerCommitInfo.Hash)

	// send vote back to validator
	testutils.MockValidatorNetReceives(t, mockNet, rootNode.PeerConf.ID, network.ProtocolRootVote, lastVoteMsg)
	// this will trigger next proposal since QC is achieved
	lastProposalMsg = testutils.MockAwaitMessage[*abdrc.ProposalMsg](t, mockNet, network.ProtocolRootProposal)
	// no additional requests have been received, meaning payload is empty
	require.True(t, lastProposalMsg.Block.Payload.IsEmpty())

	// route the proposal back to trigger new vote
	testutils.MockValidatorNetReceives(t, mockNet, rootNode.PeerConf.ID, network.ProtocolRootProposal, lastProposalMsg)
	// wait for the vote message
	lastVoteMsg = testutils.MockAwaitMessage[*abdrc.VoteMsg](t, mockNet, network.ProtocolRootVote)
	require.Equal(t, uint64(3), lastVoteMsg.VoteInfo.RoundNumber)
	require.NotNil(t, lastVoteMsg.LedgerCommitInfo.Hash)
	// send vote back to validator
	testutils.MockValidatorNetReceives(t, mockNet, rootNode.PeerConf.ID, network.ProtocolRootVote, lastVoteMsg)
	// after two successful rounds the IR change will be committed and UC is returned
	result, err := readResult(cm.CertificationResult(), time.Second)
	trustBase := map[string]crypto.Verifier{rootNode.PeerConf.ID.String(): rootNode.Verifier}
	sdrh := rg.Partitions[0].GetSystemDescriptionRecord().Hash(gocrypto.SHA256)
	require.NoError(t, result.IsValid(trustBase, gocrypto.SHA256, partitionID.ToSystemID(), sdrh))

	// roor will continue and next proposal is also triggered by the same QC
	lastProposalMsg = testutils.MockAwaitMessage[*abdrc.ProposalMsg](t, mockNet, network.ProtocolRootProposal)
	require.True(t, lastProposalMsg.Block.Payload.IsEmpty())
	require.NoError(t, err)
}

func TestPartitionTimeoutFromRootValidator(t *testing.T) {
	var lastProposalMsg *abdrc.ProposalMsg = nil
	var lastVoteMsg *abdrc.VoteMsg = nil

	mockNet := testnetwork.NewRootMockNetwork()
	cm, rootNode, _, rg := initConsensusManager(t, mockNet)
	ctx, ctxCancel := context.WithCancel(context.Background())
	defer ctxCancel()
	go func() { require.ErrorIs(t, cm.Run(ctx), context.Canceled) }()

	roundNo := uint64(1) // 1 is genesis
	// run a loop of 11 rounds to produce a root chain timeout
	for i := 0; i < int(rg.Partitions[0].SystemDescriptionRecord.T2Timeout/(rg.Root.Consensus.BlockRateMs/2)); i++ {
		// proposal rounds 2..
		roundNo++
		lastProposalMsg = testutils.MockAwaitMessage[*abdrc.ProposalMsg](t, mockNet, network.ProtocolRootProposal)
		require.Empty(t, lastProposalMsg.Block.Payload.Requests)
		// route the proposal back
		testutils.MockValidatorNetReceives(t, mockNet, rootNode.PeerConf.ID, network.ProtocolRootProposal, lastProposalMsg)
		// wait for the vote message
		lastVoteMsg = testutils.MockAwaitMessage[*abdrc.VoteMsg](t, mockNet, network.ProtocolRootVote)
		require.Equal(t, roundNo, lastVoteMsg.VoteInfo.RoundNumber)
		testutils.MockValidatorNetReceives(t, mockNet, rootNode.PeerConf.ID, network.ProtocolRootVote, lastVoteMsg)
	}
	// proposal round 7 with timeout
	roundNo++
	lastProposalMsg = testutils.MockAwaitMessage[*abdrc.ProposalMsg](t, mockNet, network.ProtocolRootProposal)
	require.NotEmpty(t, lastProposalMsg.Block.Payload.Requests)
	require.Equal(t, abtypes.T2Timeout, lastProposalMsg.Block.Payload.Requests[0].CertReason)
	testutils.MockValidatorNetReceives(t, mockNet, rootNode.PeerConf.ID, network.ProtocolRootProposal, lastProposalMsg)
	// wait for the vote message
	lastVoteMsg = testutils.MockAwaitMessage[*abdrc.VoteMsg](t, mockNet, network.ProtocolRootVote)
	require.Equal(t, roundNo, lastVoteMsg.VoteInfo.RoundNumber)
	testutils.MockValidatorNetReceives(t, mockNet, rootNode.PeerConf.ID, network.ProtocolRootVote, lastVoteMsg)
	// new proposal round 8
	roundNo++
	lastProposalMsg = testutils.MockAwaitMessage[*abdrc.ProposalMsg](t, mockNet, network.ProtocolRootProposal)
	require.Empty(t, lastProposalMsg.Block.Payload.Requests)
	testutils.MockValidatorNetReceives(t, mockNet, rootNode.PeerConf.ID, network.ProtocolRootProposal, lastProposalMsg)
	// voting round 8
	lastVoteMsg = testutils.MockAwaitMessage[*abdrc.VoteMsg](t, mockNet, network.ProtocolRootVote)
	require.Equal(t, roundNo, lastVoteMsg.VoteInfo.RoundNumber)
	testutils.MockValidatorNetReceives(t, mockNet, rootNode.PeerConf.ID, network.ProtocolRootVote, lastVoteMsg)
	// triggers timeout certificates for in round 8 to be committed
	result, err := readResult(cm.CertificationResult(), time.Second)
	require.NoError(t, err)
	require.NotNil(t, result)
	// round 7 got certified in round 8
	require.Equal(t, roundNo-1, result.UnicitySeal.RootChainRoundNumber)
	// proposal in round 9 should be empty again
	roundNo++
	lastProposalMsg = testutils.MockAwaitMessage[*abdrc.ProposalMsg](t, mockNet, network.ProtocolRootProposal)
	require.True(t, lastProposalMsg.Block.Payload.IsEmpty())
	testutils.MockValidatorNetReceives(t, mockNet, rootNode.PeerConf.ID, network.ProtocolRootProposal, lastProposalMsg)
	// vote round 9
	lastVoteMsg = testutils.MockAwaitMessage[*abdrc.VoteMsg](t, mockNet, network.ProtocolRootVote)
	require.Equal(t, roundNo, lastVoteMsg.VoteInfo.RoundNumber)
	testutils.MockValidatorNetReceives(t, mockNet, rootNode.PeerConf.ID, network.ProtocolRootVote, lastVoteMsg)
	// proposal round 10
	lastProposalMsg = testutils.MockAwaitMessage[*abdrc.ProposalMsg](t, mockNet, network.ProtocolRootProposal)
	require.True(t, lastProposalMsg.Block.Payload.IsEmpty())
	testutils.MockValidatorNetReceives(t, mockNet, rootNode.PeerConf.ID, network.ProtocolRootProposal, lastProposalMsg)
}

func TestGetState(t *testing.T) {
	mockNet := testnetwork.NewRootMockNetwork()
	cm, rootNode, partitionNodes, _ := initConsensusManager(t, mockNet)
	ctx, ctxCancel := context.WithCancel(context.Background())
	defer ctxCancel()
	go func() { require.ErrorIs(t, cm.Run(ctx), context.Canceled) }()

	getStateMsg := &abdrc.GetStateMsg{
		NodeId: partitionNodes[0].PeerConf.ID.String(),
	}
	// simulate IR change request message
	testutils.MockValidatorNetReceives(t, mockNet, rootNode.PeerConf.ID, network.ProtocolRootStateReq, getStateMsg)
	// As the node is the leader, next round will trigger a proposal
	stateMsg := testutils.MockAwaitMessage[*abdrc.StateMsg](t, mockNet, network.ProtocolRootStateResp)
	// at this stage there is only genesis block
	require.Equal(t, uint64(1), stateMsg.CommittedHead.Block.Round)
	require.Equal(t, 0, len(stateMsg.BlockNode))
	require.Len(t, stateMsg.Certificates, 1)
}

func Test_ConsensusManager_onVoteMsg(t *testing.T) {
	t.Parallel()

	// partition data used/shared by tests
	_, partitionRecord := testutils.CreatePartitionNodesAndPartitionRecord(t, partitionInputRecord, partitionID, 2)

	makeVoteMsg := func(t *testing.T, cms []*ConsensusManager, round uint64) *abdrc.VoteMsg {
		t.Helper()
		qcRoundInfo := abdrctu.NewDummyRootRoundInfo(round - 2)
		commitInfo := abdrctu.NewDummyCommitInfo(gocrypto.SHA256, qcRoundInfo)
		highQc := &abtypes.QuorumCert{
			VoteInfo:         qcRoundInfo,
			LedgerCommitInfo: commitInfo,
			Signatures:       map[string][]byte{},
		}
		cib := commitInfo.Bytes()
		for _, cm := range cms {
			sig, err := cm.safety.signer.SignBytes(cib)
			require.NoError(t, err)
			highQc.Signatures[cm.id.String()] = sig
		}

		voteRoundInfo := abdrctu.NewDummyRootRoundInfo(round)
		voteMsg := &abdrc.VoteMsg{
			VoteInfo: voteRoundInfo,
			LedgerCommitInfo: &types.UnicitySeal{
				PreviousHash: voteRoundInfo.Hash(gocrypto.SHA256),
			},
			HighQc: highQc,
			Author: cms[0].id.String(),
		}
		require.NoError(t, voteMsg.Sign(cms[0].safety.signer))
		return voteMsg
	}

	t.Run("stale vote", func(t *testing.T) {
		cms, _, _ := createConsensusManagers(t, 1, []*genesis.PartitionRecord{partitionRecord})
		cms[0].pacemaker.Reset(8, nil, nil)
		defer cms[0].pacemaker.Stop()

		vote := makeVoteMsg(t, cms, 7)
		err := cms[0].onVoteMsg(context.Background(), vote)
		require.EqualError(t, err, `stale vote for round 7 from `+cms[0].id.String())
		require.Empty(t, cms[0].voteBuffer)
	})

	t.Run("invalid vote: verify fails", func(t *testing.T) {
		// here we just test that only verified votes are processed, all the possible
		// vote verification failures should be tested by vote.Verify unit tests...
		const votedRound = 10
		cms, _, _ := createConsensusManagers(t, 1, []*genesis.PartitionRecord{partitionRecord})
		cms[0].pacemaker.Reset(votedRound-1, nil, nil)
		defer cms[0].pacemaker.Stop()

		vote := makeVoteMsg(t, cms, votedRound)
		vote.Author = "foobar"
		err := cms[0].onVoteMsg(context.Background(), vote)
		require.EqualError(t, err, `invalid vote: author 'foobar' is not in the trustbase`)
		require.Empty(t, cms[0].voteBuffer)
	})

	t.Run("vote for next round should be buffered", func(t *testing.T) {
		const votedRound = 10
		// need at least two CMs so that we do not trigger recovery because of having
		// received enough votes for the quorum
		cms, _, _ := createConsensusManagers(t, 2, []*genesis.PartitionRecord{partitionRecord})
		cms[0].pacemaker.Reset(votedRound-1, nil, nil)
		defer cms[0].pacemaker.Stop()

		vote := makeVoteMsg(t, cms, votedRound+1)
		err := cms[0].onVoteMsg(context.Background(), vote)
		require.NoError(t, err)
		require.Equal(t, vote, cms[0].voteBuffer[vote.Author])
	})

	t.Run("repeat vote for next round should be ignored (not buffered twice)", func(t *testing.T) {
		const votedRound = 10
		// need at least two CMs so that we do not trigger recovery because of having
		// received enough votes for the quorum
		cms, _, _ := createConsensusManagers(t, 2, []*genesis.PartitionRecord{partitionRecord})
		cms[0].pacemaker.Reset(votedRound-1, nil, nil)
		defer cms[0].pacemaker.Stop()

		vote := makeVoteMsg(t, cms, votedRound+1)
		err := cms[0].onVoteMsg(context.Background(), vote)
		require.NoError(t, err)
		require.Equal(t, vote, cms[0].voteBuffer[vote.Author])
		// send the vote again - should not trigger recovery ie vote is not counted again
		require.NoError(t, cms[0].onVoteMsg(context.Background(), vote))
		require.Equal(t, vote, cms[0].voteBuffer[vote.Author], "expected original vote still to be in the buffer")
		require.Len(t, cms[0].voteBuffer, 1, "expected only one vote to be buffered")
	})

	t.Run("quorum of votes for next round should trigger recovery", func(t *testing.T) {
		const votedRound = 10
		cms, _, _ := createConsensusManagers(t, 1, []*genesis.PartitionRecord{partitionRecord})
		cms[0].pacemaker.Reset(votedRound-1, nil, nil)
		defer cms[0].pacemaker.Stop()

		// as we have single CM vote means quorum and recovery should be triggered as CM hasn't
		// seen proposal yet
		vote := makeVoteMsg(t, cms, votedRound+1)
		err := cms[0].onVoteMsg(context.Background(), vote)
		require.EqualError(t, err, `have received 1 votes but no proposal, entering recovery`)
		require.Equal(t, vote, cms[0].voteBuffer[vote.Author], "expected vote to be buffered")
	})

	t.Run("not the leader of the (next) round", func(t *testing.T) {
		const votedRound = 10
		cms, _, _ := createConsensusManagers(t, 2, []*genesis.PartitionRecord{partitionRecord})
		cms[0].leaderSelector = constLeader{leader: cms[1].id, nodes: cms[1].leaderSelector.GetNodes()} // make sure this CM won't be the leader
		cms[0].pacemaker.Reset(votedRound-1, nil, nil)
		defer cms[0].pacemaker.Stop()

		vote := makeVoteMsg(t, cms, votedRound)
		err := cms[0].onVoteMsg(context.Background(), vote)
		require.EqualError(t, err, fmt.Sprintf("validator is not the leader for round %d", votedRound+1))
		require.Empty(t, cms[0].voteBuffer)
	})
}

func Test_ConsensusManager_handleRootNetMsg(t *testing.T) {
	t.Parallel()

	t.Run("untyped nil msg", func(t *testing.T) {
		cm := &ConsensusManager{}
		require.NoError(t, cm.initMetrics(observability.NOPMetrics()))
		err := cm.handleRootNetMsg(context.Background(), nil)
		require.EqualError(t, err, `unknown message type <nil>`)
	})

	t.Run("type not known for the handler", func(t *testing.T) {
		cm := &ConsensusManager{}
		require.NoError(t, cm.initMetrics(observability.NOPMetrics()))
		err := cm.handleRootNetMsg(context.Background(), "foobar")
		require.EqualError(t, err, `unknown message type string`)
	})
}

func Test_ConsensusManager_messages(t *testing.T) {
	t.Parallel()

	waitExit := func(t *testing.T, ctxCancel context.CancelFunc, doneCh chan struct{}) {
		t.Helper()
		ctxCancel()
		// and wait for cm to exit
		select {
		case <-time.After(1300 * time.Millisecond):
			t.Fatal("consensus manager did not exit in time")
		case <-doneCh:
		}
	}

	// partition data used/shared by tests
	partitionNodes, partitionRecord := testutils.CreatePartitionNodesAndPartitionRecord(t, partitionInputRecord, partitionID, 2)

	t.Run("IR change request from partition included in proposal", func(t *testing.T) {
		cms, rootNet, rootG := createConsensusManagers(t, 1, []*genesis.PartitionRecord{partitionRecord})

		// proposal will be broadcasted so eavesdrop the network and make copy of it
		propCh := make(chan *abdrc.ProposalMsg, 1)
		rootNet.SetFirewall(func(from, to peer.ID, msg any) bool {
			if msg, ok := msg.(*abdrc.ProposalMsg); ok {
				propCh <- msg
			}
			return false
		})

		ctx, cancel := context.WithCancel(context.Background())
		done := make(chan struct{})
		go func() { defer close(done); require.ErrorIs(t, cms[0].Run(ctx), context.Canceled) }()
		defer waitExit(t, cancel, done)

		// simulate root validator node sending IRCR to consensus manager
		irCReq := consensus.IRChangeRequest{
			SystemIdentifier: partitionID,
			Reason:           consensus.Quorum,
			Requests:         buildBlockCertificationRequest(t, rootG, partitionNodes),
		}

		select {
		case <-time.After(cms[0].pacemaker.minRoundLen):
			t.Fatal("CM doesn't consume IR change request")
		case cms[0].RequestCertification() <- irCReq:
		}

		// IRCR must be included into proposal
		select {
		case <-time.After(cms[0].pacemaker.maxRoundLen):
			t.Fatal("haven't got the proposal before timeout")
		case prop := <-propCh:
			require.NotNil(t, prop)
			require.NotNil(t, prop.Block)
			require.NotNil(t, prop.Block.Payload)
			require.Len(t, prop.Block.Payload.Requests, 1)
			require.EqualValues(t, irCReq.SystemIdentifier, prop.Block.Payload.Requests[0].SystemIdentifier)
			require.ElementsMatch(t, irCReq.Requests, prop.Block.Payload.Requests[0].Requests)
		}
	})

	t.Run("IR change request from partition forwarded to leader", func(t *testing.T) {
		// we create two CMs but only non-leader node has to be running as we test
		// that it will forward message to leader by monitoring network
		cms, rootNet, rootG := createConsensusManagers(t, 2, []*genesis.PartitionRecord{partitionRecord})
		cmLeader := cms[0]
		nonLeaderNode := cms[1]
		nonLeaderNode.leaderSelector = constLeader{leader: cmLeader.id, nodes: cmLeader.leaderSelector.GetNodes()} // use "const leader" to take leader selection out of test
		// eavesdrop the network and copy IR change message sent by non-leader to leader
		irCh := make(chan *abdrc.IrChangeReqMsg, 1)
		rootNet.SetFirewall(func(from, to peer.ID, msg any) bool {
			if msg, ok := msg.(*abdrc.IrChangeReqMsg); ok && from == nonLeaderNode.id && to == cmLeader.id {
				irCh <- msg
			}
			return false
		})

		ctx, cancel := context.WithCancel(context.Background())
		done := make(chan struct{})
		go func() { defer close(done); require.ErrorIs(t, nonLeaderNode.Run(ctx), context.Canceled) }()
		defer waitExit(t, cancel, done)

		// simulate partition request sent to non-leader node
		irCReq := consensus.IRChangeRequest{
			SystemIdentifier: partitionID,
			Reason:           consensus.Quorum,
			Requests:         buildBlockCertificationRequest(t, rootG, partitionNodes),
		}
		select {
		case <-time.After(nonLeaderNode.pacemaker.minRoundLen):
			t.Fatal("CM doesn't consume IR change request")
		case nonLeaderNode.RequestCertification() <- irCReq:
		}

		select {
		case <-time.After(cmLeader.pacemaker.maxRoundLen):
			t.Fatal("haven't got the IR Change message before timeout")
		case irMsg := <-irCh:
			require.NotNil(t, irMsg)
			require.Equal(t, irMsg.Author, nonLeaderNode.id.String())
			require.EqualValues(t, irCReq.SystemIdentifier, irMsg.IrChangeReq.SystemIdentifier)
			require.ElementsMatch(t, irCReq.Requests, irMsg.IrChangeReq.Requests)
		}
	})

	t.Run("IR change request forwarded by peer included in proposal", func(t *testing.T) {
		// we create two CMs but only leader is running, the other is just needed for
		// valid peer ID in the genesis so IRCR can be signed and validated
		cms, rootNet, rootG := createConsensusManagers(t, 2, []*genesis.PartitionRecord{partitionRecord})
		cmOther := cms[1]
		cmLeader := cms[0]
		cmLeader.leaderSelector = constLeader{leader: cmLeader.id, nodes: cmLeader.leaderSelector.GetNodes()} // use "const leader" to take leader selection out of test
		ctx, cancel := context.WithCancel(context.Background())
		done := make(chan struct{})
		go func() { defer close(done); require.ErrorIs(t, cmLeader.Run(ctx), context.Canceled) }()
		defer waitExit(t, cancel, done)

		// leader is expected to broadcast proposal message, snoop for it
		wire := make(chan *abdrc.ProposalMsg, 1)
		rootNet.SetFirewall(ForwardMsgs(cmLeader.id, cmOther.id, wire))

		// simulate IR change request message, "other root node" forwarding IRCR to leader
		irChReqMsg := &abdrc.IrChangeReqMsg{
			Author: cmOther.id.String(),
			IrChangeReq: &abtypes.IRChangeReq{
				SystemIdentifier: partitionID,
				CertReason:       abtypes.Quorum,
				Requests:         buildBlockCertificationRequest(t, rootG, partitionNodes[0:2]),
			},
		}
		require.NoError(t, cmOther.safety.Sign(irChReqMsg))
		require.NoError(t, cmOther.net.Send(ctx, irChReqMsg, cmLeader.id))

		// IRCR must be included into broadcast proposal, either this or next round
		sawIRCR := false
		for cnt := 0; cnt < 2 && !sawIRCR; cnt++ {
			select {
			case <-time.After(cmLeader.pacemaker.maxRoundLen):
				t.Fatal("haven't got the proposal before timeout")
			case prop := <-wire:
				require.NotNil(t, prop)
				require.NotNil(t, prop.Block)
				require.NotNil(t, prop.Block.Payload)
				if len(prop.Block.Payload.Requests) == 1 {
					require.EqualValues(t, irChReqMsg.IrChangeReq.SystemIdentifier, prop.Block.Payload.Requests[0].SystemIdentifier)
					require.ElementsMatch(t, irChReqMsg.IrChangeReq.Requests, prop.Block.Payload.Requests[0].Requests)
					sawIRCR = true
				}
			}
		}
		require.True(t, sawIRCR, "expected to see the IRCR in one of the next two proposals")
	})

	t.Run("IR change request arrives late and is forwarded to the next leader", func(t *testing.T) {
		// mimic situation where nonLeaderNode was the leader and IRCR was sent to it. However, by
		// the time msg arrives leader has changed to cmLeader so we expect nonLeaderNode to
		// forward the message.
		cms, rootNet, rootG := createConsensusManagers(t, 2, []*genesis.PartitionRecord{partitionRecord})
		cmLeader := cms[0]
		nonLeaderNode := cms[1]
		nonLeaderNode.leaderSelector = constLeader{leader: cmLeader.id, nodes: cmLeader.leaderSelector.GetNodes()}

		irCh := make(chan *abdrc.IrChangeReqMsg, 1)
		rootNet.SetFirewall(func(from, to peer.ID, msg any) bool {
			if msg, ok := msg.(*abdrc.IrChangeReqMsg); ok && from == nonLeaderNode.id && to == cmLeader.id {
				irCh <- msg
			}
			return false
		})

		ctx, cancel := context.WithCancel(context.Background())
		done := make(chan struct{})
		go func() { defer close(done); require.ErrorIs(t, nonLeaderNode.Run(ctx), context.Canceled) }()
		defer waitExit(t, cancel, done)

		// send IRCR to non-leader, simulating message arriving late, leader has changed
		irChReqMsg := &abdrc.IrChangeReqMsg{
			Author: cmLeader.id.String(),
			IrChangeReq: &abtypes.IRChangeReq{
				SystemIdentifier: partitionID,
				CertReason:       abtypes.Quorum,
				Requests:         buildBlockCertificationRequest(t, rootG, partitionNodes[0:2]),
			},
		}
		require.NoError(t, cmLeader.safety.Sign(irChReqMsg))
		rootNet.Send(irChReqMsg, nonLeaderNode.id)

		// non-leader is not the next leader and must forward the request to the leader node
		select {
		case <-time.After(cmLeader.pacemaker.maxRoundLen):
			t.Fatal("haven't seen forwarded proposal before timeout")
		case irMsg := <-irCh:
			require.NotNil(t, irMsg)
			require.Equal(t, irMsg.Author, cmLeader.id.String())
			require.EqualValues(t, irChReqMsg.IrChangeReq.SystemIdentifier, irMsg.IrChangeReq.SystemIdentifier)
			require.ElementsMatch(t, irChReqMsg.IrChangeReq.Requests, irMsg.IrChangeReq.Requests)
		}
	})

	t.Run("state request triggers response", func(t *testing.T) {
		cms, _, _ := createConsensusManagers(t, 2, []*genesis.PartitionRecord{partitionRecord})
		cmA, cmB := cms[0], cms[1]

		// only launch cmA, we manage cmB "manually"
		ctx, cancel := context.WithCancel(context.Background())
		done := make(chan struct{})
		go func() { defer close(done); require.ErrorIs(t, cmA.Run(ctx), context.Canceled) }()
		defer waitExit(t, cancel, done)

		// cmB sends state request to cmA
		msg := &abdrc.GetStateMsg{NodeId: cmB.id.String()}
		require.NoError(t, cmB.net.Send(ctx, msg, cmA.id))

		// cmB should receive state response
		select {
		case <-time.After(1000 * time.Millisecond):
			t.Fatal("timeout while waiting for recovery response")
		case msg := <-cmB.net.ReceivedChannel():
			state := msg.(*abdrc.StateMsg)
			require.NotNil(t, state)
			require.EqualValues(t, 1, state.CommittedHead.Block.Round)
			require.Empty(t, state.BlockNode)
			require.Len(t, state.Certificates, 1)
		}
	})
}

func Test_ConsensusManager_sendCertificates(t *testing.T) {
	t.Parallel()

	_, partitionRecord := testutils.CreatePartitionNodesAndPartitionRecord(t, partitionInputRecord, partitionID, 2)

	// generate UCs for given systems (with random data in QC)
	makeUCs := func(sysID ...types.SystemID32) map[types.SystemID32]*types.UnicityCertificate {
		rUC := make(map[types.SystemID32]*types.UnicityCertificate)
		for _, id := range sysID {
			uc := &types.UnicityCertificate{
				UnicityTreeCertificate: &types.UnicityTreeCertificate{
					SystemIdentifier:      id.ToSystemID(),
					SystemDescriptionHash: test.RandomBytes(32),
				},
			}
			rUC[id] = uc
		}
		return rUC
	}

	// consumeUCs reads UCs from "cm"-s output and stores them into map it returns.
	// it reads until "timeout" has passed.
	consumeUCs := func(cm *ConsensusManager, timeout time.Duration) map[types.SystemID32]*types.UnicityCertificate {
		to := time.After(timeout)
		rUC := make(map[types.SystemID32]*types.UnicityCertificate)
		for {
			select {
			case uc := <-cm.CertificationResult():
				id, err := uc.UnicityTreeCertificate.SystemIdentifier.Id32()
				require.NoError(t, err)
				rUC[id] = uc
			case <-to:
				return rUC
			}
		}
	}

	outputMustBeEmpty := func(t *testing.T, cm *ConsensusManager) {
		t.Helper()
		select {
		case uc := <-cm.CertificationResult():
			t.Errorf("unexpected data from cert chan: %#v", uc)
		default:
		}
	}

	t.Run("consume before next input is sent", func(t *testing.T) {
		cms, _, _ := createConsensusManagers(t, 1, []*genesis.PartitionRecord{partitionRecord})
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		go func() { require.ErrorIs(t, cms[0].sendCertificates(ctx), context.Canceled) }()

		// send some certificates into sink...
		ucs := makeUCs(types.SystemID32(1), types.SystemID32(2))
		select {
		case cms[0].ucSink <- ucs:
		default:
			t.Fatal("expected that input would be accepted immediately, sink should be empty")
		}
		//...and consume them
		rUC := consumeUCs(cms[0], 100*time.Millisecond)
		require.Equal(t, ucs, rUC)
		outputMustBeEmpty(t, cms[0])

		// and repeat the exercise with different systemIDs
		ucs = makeUCs(types.SystemID32(3), types.SystemID32(4))
		select {
		case cms[0].ucSink <- ucs:
		default:
			t.Fatal("expected that input would be accepted immediately, sink should be empty")
		}

		rUC = consumeUCs(cms[0], 100*time.Millisecond)
		require.Equal(t, ucs, rUC)
		outputMustBeEmpty(t, cms[0])
	})

	t.Run("overwriting unconsumed QC", func(t *testing.T) {
		cms, _, _ := createConsensusManagers(t, 1, []*genesis.PartitionRecord{partitionRecord})
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		go func() { require.ErrorIs(t, cms[0].sendCertificates(ctx), context.Canceled) }()

		// exp - expected result in the end of test, we add/overwrite certs as we send them
		exp := map[types.SystemID32]*types.UnicityCertificate{}

		ucs := makeUCs(types.SystemID32(1), types.SystemID32(2))
		for k, v := range ucs {
			exp[k] = v
		}

		select {
		case cms[0].ucSink <- ucs:
		default:
			t.Fatal("expected that input would be accepted immediately, sink should be empty")
		}

		// as we haven't consumed anything sending new set of certs into the sink should
		// overwrite {0,0,0,2} and add {0,0,0,3}
		ucs = makeUCs(types.SystemID32(3), types.SystemID32(2))
		for k, v := range ucs {
			exp[k] = v
		}

		select {
		case cms[0].ucSink <- ucs:
		case <-time.After(100 * time.Millisecond):
			t.Error("next input hasn't been consumed fast enough")
		}

		rUC := consumeUCs(cms[0], 100*time.Millisecond)
		require.Len(t, rUC, 3, "number of different systemIDs")
		require.Equal(t, exp, rUC)
		outputMustBeEmpty(t, cms[0])
	})

	t.Run("adding without overwriting", func(t *testing.T) {
		cms, _, _ := createConsensusManagers(t, 1, []*genesis.PartitionRecord{partitionRecord})
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		go func() { require.ErrorIs(t, cms[0].sendCertificates(ctx), context.Canceled) }()
		// exp - expected result in the end of test, we add/overwrite certs as we send them
		exp := map[types.SystemID32]*types.UnicityCertificate{}

		ucs := makeUCs(types.SystemID32(1), types.SystemID32(2))
		for k, v := range ucs {
			exp[k] = v
		}

		select {
		case cms[0].ucSink <- ucs:
		default:
			t.Fatal("expected that input would be accepted immediately, sink should be empty")
		}

		// send another set of certs, unique sysIDs, ie no overwrites, just add (nothing
		// has been consumed yet)
		ucs = makeUCs(types.SystemID32(3), types.SystemID32(4))
		for k, v := range ucs {
			exp[k] = v
		}

		select {
		case cms[0].ucSink <- ucs:
		case <-time.After(100 * time.Millisecond):
			t.Error("next input hasn't been consumed fast enough")
		}

		rUC := consumeUCs(cms[0], 100*time.Millisecond)
		require.Len(t, rUC, 4, "number of different systemIDs")
		require.Equal(t, exp, rUC)
		outputMustBeEmpty(t, cms[0])
	})

	t.Run("concurrency", func(t *testing.T) {
		// concurrent read and writes to trip race detector
		cms, _, _ := createConsensusManagers(t, 1, []*genesis.PartitionRecord{partitionRecord})
		done := make(chan struct{})
		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			defer close(done)
			require.ErrorIs(t, cms[0].sendCertificates(ctx), context.Canceled)
		}()

		go func() {
			for {
				select {
				case cms[0].ucSink <- makeUCs(types.SystemID32(1), types.SystemID32(2), types.SystemID32(3)):
				case <-ctx.Done():
					return
				}
			}
		}()

		go func() {
			for {
				select {
				case <-cms[0].CertificationResult():
				case <-ctx.Done():
					return
				}
			}
		}()

		time.Sleep(time.Second)
		cancel()
		<-done
	})
}

func Test_selectRandomNodeIdsFromSignatureMap(t *testing.T) {
	t.Parallel()

	// generate some valid peer IDs for tests to use
	peerIDs := make(peer.IDSlice, 3)
	for i := range peerIDs {
		_, publicKey, err := p2pcrypto.GenerateSecp256k1Key(rand.Reader)
		require.NoError(t, err)
		pubKeyBytes, err := publicKey.Raw()
		require.NoError(t, err)
		peerIDs[i], err = network.NodeIDFromPublicKeyBytes(pubKeyBytes)
		require.NoError(t, err)
	}
	idA, idB, idC := peerIDs[0], peerIDs[1], peerIDs[2]

	t.Run("empty inputs", func(t *testing.T) {
		nodes := selectRandomNodeIdsFromSignatureMap(nil, 2)
		require.Empty(t, nodes)

		nodes = selectRandomNodeIdsFromSignatureMap(map[string][]byte{}, 2)
		require.Empty(t, nodes)
	})

	t.Run("no duplicates added", func(t *testing.T) {
		nodes := selectRandomNodeIdsFromSignatureMap(map[string][]byte{idA.String(): nil}, 2)
		require.ElementsMatch(t, []peer.ID{idA}, nodes)

		nodes = selectRandomNodeIdsFromSignatureMap(map[string][]byte{idA.String(): nil, idB.String(): nil}, 2)
		require.ElementsMatch(t, []peer.ID{idA, idB}, nodes)

		nodes = selectRandomNodeIdsFromSignatureMap(map[string][]byte{idA.String(): nil, idB.String(): nil}, 3)
		require.ElementsMatch(t, []peer.ID{idA, idB}, nodes)
	})

	t.Run("invalid IDs are ignored", func(t *testing.T) {
		nodes := selectRandomNodeIdsFromSignatureMap(map[string][]byte{"foo bar": nil}, 1)
		require.Empty(t, nodes)

		nodes = selectRandomNodeIdsFromSignatureMap(map[string][]byte{"foo bar": nil, idB.String(): nil}, 2)
		require.ElementsMatch(t, []peer.ID{idB}, nodes)
	})

	t.Run("max count items is returned", func(t *testing.T) {
		inp := map[string][]byte{idA.String(): nil, idB.String(): nil, idC.String(): nil}

		nodes := selectRandomNodeIdsFromSignatureMap(inp, 1)
		require.Len(t, nodes, 1)

		nodes = selectRandomNodeIdsFromSignatureMap(inp, 2)
		require.Len(t, nodes, 2)
		require.NotEqual(t, nodes[0], nodes[1])
	})
}

func Test_rootNetworkRunning(t *testing.T) {
	t.Parallel()
	// this test is mostly useful for debugging - modify conditions in the test,
	// launch the test and observe logs...

	// for quorum we need ⅔+1 validators (root nodes) to be healthy
	const rootNodeCnt = 4
	// destination round - until which round (minimum) the test should run. test stops as soon as
	// one node is in that round so last round might be "incomplete". System starts with round 2.
	const destRound = 10

	// consumeUC acts as a validator node consuming the UC-s generated by CM (until ctx is cancelled).
	// strictly speaking not needed as current implementation should continue working even when there
	// is no-one consuming UCs.
	consumeUC := func(ctx context.Context, cm *ConsensusManager) {
		for {
			select {
			case <-ctx.Done():
				return
			case uc := <-cm.certResultCh:
				t.Logf("%s UC for round %d sent to validator", cm.id.ShortString(), uc.GetRoundNumber())
			}
		}
	}

	partitionRecs := []*genesis.PartitionRecord{
		createPartitionRecord(t, partitionID.ToSystemID(), partitionInputRecord, 1),
	}

	cms, rootNet, _ := createConsensusManagers(t, rootNodeCnt, partitionRecs)

	var totalMsgCnt atomic.Uint32
	rootNet.SetFirewall(func(from, to peer.ID, msg any) bool {
		msgCnt := totalMsgCnt.Add(1)
		block := msgCnt%200 == 0 // drop every n-th message from the network
		t.Logf("%t # %s -> %s : %T", block, from.ShortString(), to.ShortString(), msg)
		return block
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan struct{})
	start := time.Now()
	go func() {
		wg := sync.WaitGroup{}
		wg.Add(len(cms))
		for _, v := range cms {
			go func(cm *ConsensusManager) {
				defer wg.Done()
				require.ErrorIs(t, cm.Run(ctx), context.Canceled)
			}(v)
			go func(cm *ConsensusManager) { consumeUC(ctx, cm) }(v)
		}
		wg.Wait()
		close(done)
	}()
	cm := cms[0]
	// assume rounds are successful and each round takes between minRoundLen and roundTimeout on average
	maxTestDuration := destRound * (cm.pacemaker.minRoundLen + (cm.pacemaker.maxRoundLen-cm.pacemaker.minRoundLen)/2)
	require.Eventually(t, func() bool { return cm.pacemaker.GetCurrentRound() >= destRound }, maxTestDuration, 100*time.Millisecond, "waiting for round %d to be achieved", destRound)
	stop := time.Now()
	cancel()
	// when calculating expected message counts keep in mind that last round might not be complete
	// ie the test ended before all nodes had a chance to post their message. so use -1 rounds!
	// and we start from round 2 so thats another -1 completed rounds.
	completeRounds := cm.pacemaker.GetCurrentRound() - 2
	avgRoundLen := time.Duration(int64(stop.Sub(start)) / int64(completeRounds))

	// output some statistics
	t.Logf("total msg count: %d during %s", totalMsgCnt.Load(), stop.Sub(start))

	// check some expectations
	// we expect to see proposal + vote per node per round. when some round timeouts msg count is higher!
	require.GreaterOrEqual(t, totalMsgCnt.Load(), uint32(completeRounds*rootNodeCnt*2), "total number of messages in the network")
	// average round duration should be between minRoundLen and maxRoundLen (aka timeout)
	// potentially flaky as there is delay between starting CMs and starting the clock!
	require.GreaterOrEqual(t, avgRoundLen, cm.pacemaker.minRoundLen, "minimum round duration for %d rounds", completeRounds)
	require.GreaterOrEqual(t, cm.pacemaker.maxRoundLen, avgRoundLen, "maximum round duration for %d rounds", completeRounds)
	// wait for cm routine to exit, otherwise logger may be destructed before last usage
	select {
	case <-time.After(1000 * time.Millisecond):
		t.Fatal("consensus managers did not exit in time")
	case <-done:
	}
}

func TestConsensusManger_ResoreVote(t *testing.T) {
	net := testnetwork.NewRootMockNetwork()
	_, partitionRecord := testutils.CreatePartitionNodesAndPartitionRecord(t, partitionInputRecord, partitionID, 3)
	rootNode := testutils.NewTestNode(t)
	verifier := rootNode.Verifier
	rootPubKeyBytes, err := verifier.MarshalPublicKey()
	require.NoError(t, err)
	id := rootNode.PeerConf.ID
	rootGenesis, _, err := rootgenesis.NewRootGenesis(id.String(),
		rootNode.Signer,
		rootPubKeyBytes,
		[]*genesis.PartitionRecord{partitionRecord},
		rootgenesis.WithBlockRate(200),
		rootgenesis.WithConsensusTimeout(2200),
	)
	require.NoError(t, err)
	partitions, err := partitions.NewPartitionStoreFromGenesis(rootGenesis.Partitions)
	require.NoError(t, err)
	// load timeout vote to DB
	db := memorydb.New()
	// init DB from genesis
	_, err = storage.New(gocrypto.SHA256, rootGenesis.Partitions, db)
	require.NoError(t, err)
	timeoutVote := &abdrc.TimeoutMsg{Timeout: &abtypes.Timeout{Round: 2}, Author: "test"}
	require.NoError(t, storage.WriteVote(db, timeoutVote))
	cm, err := NewDistributedAbConsensusManager(id, rootGenesis, partitions, net, rootNode.Signer,
		observability.NOPMetrics(),
		testlogger.New(t).With(logger.NodeID(id)),
		consensus.WithStorage(db),
	)
	// replace leader selector
	allNodes := cm.leaderSelector.GetNodes()
	leaderId, err := p2ptest.RandPeerID()
	require.NoError(t, err)
	allNodes = append(allNodes, leaderId)
	cm.leaderSelector = constLeader{leader: leaderId, nodes: allNodes}
	require.NoError(t, err)
	require.NotNil(t, cm)
	ctx, ctxCancel := context.WithCancel(context.Background())
	defer ctxCancel()
	go func() { require.ErrorIs(t, cm.Run(ctx), context.Canceled) }()
	lastTimeoutMsg := testutils.MockAwaitMessage[*abdrc.TimeoutMsg](t, net, network.ProtocolRootTimeout)
	require.NotNil(t, lastTimeoutMsg)
	// make sure the stored timeout vote is broadcast
	require.EqualValues(t, 2, lastTimeoutMsg.Timeout.Round)
	require.EqualValues(t, "test", lastTimeoutMsg.Author)
}
