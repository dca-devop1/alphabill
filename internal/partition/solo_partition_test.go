package partition

import (
	gocrypto "crypto"
	"fmt"
	"reflect"
	"testing"
	"time"

	"gitdc.ee.guardtime.com/alphabill/alphabill/internal/block"
	"gitdc.ee.guardtime.com/alphabill/alphabill/internal/certificates"
	"gitdc.ee.guardtime.com/alphabill/alphabill/internal/crypto"
	"gitdc.ee.guardtime.com/alphabill/alphabill/internal/errors"
	"gitdc.ee.guardtime.com/alphabill/alphabill/internal/network"
	"gitdc.ee.guardtime.com/alphabill/alphabill/internal/network/protocol/blockproposal"
	"gitdc.ee.guardtime.com/alphabill/alphabill/internal/network/protocol/certification"
	"gitdc.ee.guardtime.com/alphabill/alphabill/internal/network/protocol/genesis"
	"gitdc.ee.guardtime.com/alphabill/alphabill/internal/partition/store"
	"gitdc.ee.guardtime.com/alphabill/alphabill/internal/rootchain"
	"gitdc.ee.guardtime.com/alphabill/alphabill/internal/rootchain/unicitytree"
	testnetwork "gitdc.ee.guardtime.com/alphabill/alphabill/internal/testutils/network"
	testsig "gitdc.ee.guardtime.com/alphabill/alphabill/internal/testutils/sig"
	"gitdc.ee.guardtime.com/alphabill/alphabill/internal/txsystem"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/stretchr/testify/require"
)

type AlwaysValidBlockProposalValidator struct{}
type AlwaysValidTransactionValidator struct{}

type SingleNodePartition struct {
	nodeConf   *configuration
	store      store.BlockStore
	partition  *Node
	rootState  *rootchain.State
	rootSigner crypto.Signer
	mockNet    *testnetwork.MockNet
}

func (t *AlwaysValidTransactionValidator) Validate(_ txsystem.GenericTransaction) error {
	return nil
}
func (t *AlwaysValidBlockProposalValidator) Validate(*blockproposal.BlockProposal, crypto.Verifier) error {
	return nil
}

func NewSingleNodePartition(t *testing.T, txSystem txsystem.TransactionSystem) *SingleNodePartition {
	p := createPeer(t)
	key, err := p.PublicKey()
	require.NoError(t, err)
	pubKeyBytes, err := key.Raw()
	require.NoError(t, err)

	// node genesis
	nodeSigner, _ := testsig.CreateSignerAndVerifier(t)

	systemId := []byte{1, 1, 1, 1}
	nodeGenesis, err := NewNodeGenesis(
		txSystem,
		WithPeerID("1"),
		WithSigningKey(nodeSigner),
		WithEncryptionPubKey(pubKeyBytes),
		WithSystemIdentifier(systemId),
		WithT2Timeout(2500),
	)
	require.NoError(t, err)

	// root genesis
	rootSigner, _ := testsig.CreateSignerAndVerifier(t)
	_, encPubKey := testsig.CreateSignerAndVerifier(t)

	rootGenesis, partitionGenesis, err := rootchain.NewGenesisFromPartitionNodes([]*genesis.PartitionNode{nodeGenesis}, rootSigner, encPubKey)
	if err != nil {
		t.Error(err)
	}
	require.NoError(t, err)

	// root chain
	rc, err := rootchain.NewStateFromGenesis(rootGenesis, rootSigner)
	require.NoError(t, err)

	net := testnetwork.NewMockNetwork()
	// partition
	n, err := New(
		p,
		nodeSigner,
		txSystem,
		partitionGenesis[0],
		net,
		WithT1Timeout(100*time.Minute),
		WithLeaderSelector(&TestLeaderSelector{
			leader:      "1",
			currentNode: "1",
		}),
		WithTxValidator(&AlwaysValidTransactionValidator{}),
	)
	require.NoError(t, err)
	n.blockProposalValidator = &AlwaysValidBlockProposalValidator{}

	partition := &SingleNodePartition{
		partition:  n,
		rootState:  rc,
		nodeConf:   n.configuration,
		store:      n.blockStore,
		rootSigner: rootSigner,
		mockNet:    net,
	}
	return partition
}

func (sn *SingleNodePartition) Close() {
	sn.partition.Close()
	close(sn.mockNet.MessageCh)
}

func (sn *SingleNodePartition) SubmitTx(tx *txsystem.Transaction) error {
	sn.mockNet.Receive(network.ReceivedMessage{
		From:     "from-test",
		Protocol: network.ProtocolInputForward,
		Message:  tx,
	})
	return nil
}

func (sn *SingleNodePartition) SubmitUnicityCertificate(uc *certificates.UnicityCertificate) error {
	return sn.partition.handleUnicityCertificate(uc)
}

func (sn *SingleNodePartition) HandleBlockProposal(prop *blockproposal.BlockProposal) error {
	return sn.partition.handleBlockProposal(prop)
}

func (sn *SingleNodePartition) SubmitBlockProposal(prop *blockproposal.BlockProposal) error {
	sn.mockNet.Receive(network.ReceivedMessage{
		From:     "from-test",
		Protocol: network.ProtocolBlockProposal,
		Message:  prop,
	})
	return nil
}

func (sn *SingleNodePartition) GetProposalTxs() []txsystem.GenericTransaction {
	return sn.partition.proposal
}

func (sn *SingleNodePartition) CreateUnicityCertificate(ir *certificates.InputRecord, roundNumber uint64, previousRoundRootHash []byte) (*certificates.UnicityCertificate, error) {
	id := sn.nodeConf.GetSystemIdentifier()
	sdrHash := sn.nodeConf.genesis.SystemDescriptionRecord.Hash(gocrypto.SHA256)
	data := []*unicitytree.Data{{
		SystemIdentifier:            id,
		InputRecord:                 ir,
		SystemDescriptionRecordHash: sdrHash,
	},
	}
	ut, err := unicitytree.New(gocrypto.SHA256.New(), data)
	if err != nil {
		return nil, err
	}
	rootHash := ut.GetRootHash()
	unicitySeal, err := sn.createUnicitySeal(roundNumber, previousRoundRootHash, rootHash)
	if err != nil {
		return nil, err
	}
	cert, err := ut.GetCertificate(id)
	if err != nil {
		// this should never happen. if it does then exit with panic because we cannot generate
		// unicity tree certificates.
		panic(err)
	}

	return &certificates.UnicityCertificate{
		InputRecord: ir,
		UnicityTreeCertificate: &certificates.UnicityTreeCertificate{
			SystemIdentifier:      cert.SystemIdentifier,
			SiblingHashes:         cert.SiblingHashes,
			SystemDescriptionHash: sdrHash,
		},
		UnicitySeal: unicitySeal,
	}, nil
}

func (sn *SingleNodePartition) createUnicitySeal(roundNumber uint64, previousRoundRootHash, rootHash []byte) (*certificates.UnicitySeal, error) {
	u := &certificates.UnicitySeal{
		RootChainRoundNumber: roundNumber,
		PreviousHash:         previousRoundRootHash,
		Hash:                 rootHash,
	}
	return u, u.Sign(sn.rootSigner)
}

func (sn *SingleNodePartition) GetLatestBlock() *block.Block {
	return sn.store.LatestBlock()
}

func (sn *SingleNodePartition) CreateBlock() error {
	sn.partition.handleT1TimeoutEvent()

	certificationRequests := sn.mockNet.SentMessages[network.ProtocolBlockCertification]
	if len(certificationRequests) != 1 {
		return errors.New("block certification request not found")
	}
	req := certificationRequests[0].Message.(*certification.BlockCertificationRequest)
	sn.mockNet.SentMessages[network.ProtocolBlockCertification] = []testnetwork.PeerMessage{}
	_, err := sn.rootState.HandleBlockCertificationRequest(req)
	if err != nil {
		return err
	}
	systemIds, err := sn.rootState.CreateUnicityCertificates()
	if err != nil {
		return err
	}
	if len(systemIds) != 1 {
		return errors.New("uc not created")
	}
	uc := sn.rootState.GetLatestUnicityCertificate(systemIds[0])
	sn.mockNet.Receive(network.ReceivedMessage{
		From:     "from-test",
		Protocol: network.ProtocolUnicityCertificates,
		Message:  uc,
	})
	return nil
}

type TestLeaderSelector struct {
	leader      peer.ID
	currentNode peer.ID
}

func (l *TestLeaderSelector) SelfID() peer.ID {
	return l.currentNode
}

// IsCurrentNodeLeader returns true it current node is the leader and must propose the next block.
func (l *TestLeaderSelector) IsCurrentNodeLeader() bool {
	return l.leader == l.SelfID()
}

func (l *TestLeaderSelector) UpdateLeader(seal *certificates.UnicitySeal) {
	if seal == nil {
		l.leader = ""
		return
	}
	l.leader = l.currentNode
	return
}

func (l *TestLeaderSelector) GetLeaderID() peer.ID {
	return l.leader
}

func (l *TestLeaderSelector) LeaderFromUnicitySeal(seal *certificates.UnicitySeal) peer.ID {
	if seal == nil {
		return ""
	}
	return l.currentNode
}

func createPeer(t *testing.T) *network.Peer {
	conf := &network.PeerConfiguration{}
	peer, err := network.NewPeer(conf)
	require.NoError(t, err)

	pubKey, err := peer.PublicKey()
	require.NoError(t, err)

	pubKeyBytes, err := pubKey.Raw()
	require.NoError(t, err)

	conf.PersistentPeers = []*network.PeerInfo{{
		Address:   fmt.Sprintf("%v", peer.MultiAddresses()[0]),
		PublicKey: pubKeyBytes,
	}}
	return peer
}

func ProposalSize(tp *SingleNodePartition, i int) func() bool {
	return func() bool {
		return len(tp.GetProposalTxs()) == i
	}
}

func NextBlockReceived(tp *SingleNodePartition, prevBlock *block.Block) func() bool {
	return func() bool {
		b := tp.GetLatestBlock()
		return b.UnicityCertificate.UnicitySeal.RootChainRoundNumber == prevBlock.UnicityCertificate.UnicitySeal.GetRootChainRoundNumber()+1
	}
}

func ProposalContains(tp *SingleNodePartition, t *txsystem.Transaction) func() bool {
	return func() bool {
		for _, tx := range tp.GetProposalTxs() {
			if reflect.DeepEqual(tx.ToProtoBuf(), t) {
				return true
			}
		}
		return false
	}
}

func ContainsTransaction(block *block.Block, tx *txsystem.Transaction) bool {
	for _, t := range block.Transactions {
		if t == tx {
			return true
		}
	}
	return false
}

func CertificationRequestReceived(tp *SingleNodePartition) func() bool {
	return func() bool {
		messages := tp.mockNet.SentMessages[network.ProtocolBlockCertification]
		return len(messages) > 0
	}
}
