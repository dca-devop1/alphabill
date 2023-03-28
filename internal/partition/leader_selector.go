package partition

import (
	"crypto"
	"errors"
	"sync"

	"github.com/alphabill-org/alphabill/internal/certificates"
	"github.com/alphabill-org/alphabill/internal/network"
	"github.com/alphabill-org/alphabill/internal/util"
	"github.com/holiman/uint256"
	"github.com/libp2p/go-libp2p/core/peer"
)

const (
	UnknownLeader = ""

	ErrStrSystemIdentifierIsNil = "system identifier is nil"
)

type (
	LeaderSelector interface {
		UpdateLeader(seal *certificates.UnicityCertificate)
		LeaderFromUnicitySeal(seal *certificates.UnicityCertificate) peer.ID
		IsCurrentNodeLeader() bool
		GetLeaderID() peer.ID
		SelfID() peer.ID
	}

	// DefaultLeaderSelector is used to select a leader from the validator pool.
	DefaultLeaderSelector struct {
		mutex            sync.Mutex
		leader           peer.ID // current leader ID
		systemIdentifier []byte
		self             *network.Peer
	}
)

func NewDefaultLeaderSelector(self *network.Peer, systemIdentifier []byte) (*DefaultLeaderSelector, error) {
	if self == nil {
		return nil, ErrPeerIsNil
	}
	if systemIdentifier == nil {
		return nil, errors.New(ErrStrSystemIdentifierIsNil)
	}
	return &DefaultLeaderSelector{self: self, leader: UnknownLeader, systemIdentifier: systemIdentifier}, nil
}

func (l *DefaultLeaderSelector) SelfID() peer.ID {
	return l.self.ID()
}

// IsCurrentNodeLeader returns true it current node is the leader and must propose the next block.
func (l *DefaultLeaderSelector) IsCurrentNodeLeader() bool {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	logger.Info("current leader: '%v', current node '%v'", l.leader, l.self.ID())
	return l.leader == l.self.ID()
}

func (l *DefaultLeaderSelector) GetLeaderID() peer.ID {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	return l.leader
}

// UpdateLeader updates the next block proposer. If input is nil then leader is set to UnknownLeader.
func (l *DefaultLeaderSelector) UpdateLeader(uc *certificates.UnicityCertificate) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.leader = l.LeaderFromUnicitySeal(uc)
	logger.Debug("Leader set to '%v'", l.leader)
}

func (l *DefaultLeaderSelector) LeaderFromUnicitySeal(uc *certificates.UnicityCertificate) peer.ID {
	// We don't need the lock because we don't change the state of the struct.
	if uc == nil {
		return UnknownLeader
	}
	peerCount := uint64(len(l.self.Configuration().PersistentPeers))
	hasher := crypto.SHA256.New()
	hasher.Write(util.Uint64ToBytes(uc.UnicitySeal.RootChainRoundNumber))
	hasher.Write(uc.UnicitySeal.Hash)
	hasher.Write(l.systemIdentifier)
	hasher.Write(uc.InputRecord.Bytes())
	hash := hasher.Sum(nil)
	x := uint256.NewInt(0).SetBytes(hash)
	i := uint256.NewInt(0)
	i.Mod(x, uint256.NewInt(peerCount))
	index := i.Uint64()
	if index > peerCount {
		logger.Warning("Invalid leader index.")
		return UnknownLeader
	}
	leader, err := l.self.Configuration().PersistentPeers[index].GetID()
	if err != nil {
		logger.Warning("Failed to get leader ID: %v", err)
		return UnknownLeader
	}
	return leader
}
