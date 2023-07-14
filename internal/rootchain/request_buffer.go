package rootchain

import (
	"crypto/sha256"
	"errors"
	"sync"

	"github.com/alphabill-org/alphabill/internal/network/protocol"
	"github.com/alphabill-org/alphabill/internal/network/protocol/certification"
	"github.com/alphabill-org/alphabill/internal/rootchain/partitions"
)

type (
	QuorumResult uint8

	CertRequestBuffer struct {
		mu    sync.RWMutex
		store map[protocol.SystemIdentifier]*requestBuffer
	}

	// requestBuffer keeps track of received certification nodeRequest and counts state hashes.
	sha256Hash    [sha256.Size]byte
	requestBuffer struct {
		nodeRequest map[string]sha256Hash                                     // node request register - one request per node id per round (key is node identifier)
		requests    map[sha256Hash][]*certification.BlockCertificationRequest // unique requests, IR record hash to certification request
	}
)

const (
	QuorumInProgress QuorumResult = iota
	Quorum
	QuorumNotPossible
)

// NewCertificationRequestBuffer create new certification nodeRequest buffer
func NewCertificationRequestBuffer() *CertRequestBuffer {
	return &CertRequestBuffer{
		store: make(map[protocol.SystemIdentifier]*requestBuffer),
	}
}

// Add request to certification store. Per node id first valid request is stored. Rest are either duplicate or
// equivocating and in both cases error is returned. Clear or Reset in order to receive new nodeRequest
func (c *CertRequestBuffer) Add(request *certification.BlockCertificationRequest, tb partitions.PartitionTrustBase) (QuorumResult, []*certification.BlockCertificationRequest, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	rs := c.get(protocol.SystemIdentifier(request.SystemIdentifier))
	return rs.add(request, tb)
}

// IsConsensusReceived has partition with id reached consensus. Required nrOfNodes as input to calculate consensus
func (c *CertRequestBuffer) IsConsensusReceived(id protocol.SystemIdentifier, tb partitions.PartitionTrustBase) QuorumResult {
	c.mu.RLock()
	defer c.mu.RUnlock()
	rs := c.get(id)
	_, res := rs.isConsensusReceived(tb)
	return res
}

// Reset removed all incoming nodeRequest from all stores
func (c *CertRequestBuffer) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, rs := range c.store {
		rs.reset()
	}
}

// Clear clears nodeRequest in one partition
func (c *CertRequestBuffer) Clear(id protocol.SystemIdentifier) {
	c.mu.Lock()
	defer c.mu.Unlock()
	rs := c.get(id)
	rs.reset()
}

// get returns an existing store for system identifier or registers and returns a new one if none existed
func (c *CertRequestBuffer) get(id protocol.SystemIdentifier) *requestBuffer {
	rs, f := c.store[id]
	if !f {
		rs = newRequestStore()
		c.store[id] = rs
	}
	return rs
}

// newRequestStore creates a new empty requestBuffer.
func newRequestStore() *requestBuffer {
	s := &requestBuffer{
		nodeRequest: make(map[string]sha256Hash),                                     // request per node is allowed per partition round
		requests:    make(map[sha256Hash][]*certification.BlockCertificationRequest), // count matching requests
	}
	return s
}

// add stores a new input record received from the node.
func (rs *requestBuffer) add(req *certification.BlockCertificationRequest, tb partitions.PartitionTrustBase) (QuorumResult, []*certification.BlockCertificationRequest, error) {
	_, f := rs.nodeRequest[req.NodeIdentifier]
	if f {
		proof, res := rs.isConsensusReceived(tb)
		return res, proof, errors.New("request in this round already stored, rejected")
	}
	// first request, calculate hash of IR and store
	hash := sha256.Sum256(req.InputRecord.Bytes())
	// no request in buffered yet, add
	rs.nodeRequest[req.NodeIdentifier] = hash
	rs.requests[hash] = append(rs.requests[hash], req)
	proof, res := rs.isConsensusReceived(tb)
	return res, proof, nil
}

func (rs *requestBuffer) reset() {
	rs.nodeRequest = make(map[string]sha256Hash)
	rs.requests = make(map[sha256Hash][]*certification.BlockCertificationRequest)
}

func (rs *requestBuffer) isConsensusReceived(tb partitions.PartitionTrustBase) ([]*certification.BlockCertificationRequest, QuorumResult) {
	var (
		reqCount                                            = 0
		bcReqs   []*certification.BlockCertificationRequest = nil
	)
	// find IR hash that has with most matching requests
	for _, reqs := range rs.requests {
		if reqCount < len(reqs) {
			reqCount = len(reqs)
			bcReqs = reqs
		}
	}
	quorum := int(tb.GetQuorum())
	if reqCount >= quorum {
		// consensus received
		return bcReqs, Quorum
	}
	// enough nodes have voted and even if the rest of the votes are for the most popular, quorum is still not possible
	if int(tb.GetTotalNodes())-len(rs.nodeRequest)+reqCount < quorum {
		// consensus not possible
		allReq := make([]*certification.BlockCertificationRequest, 0, len(rs.nodeRequest))
		for _, req := range rs.requests {
			allReq = append(allReq, req...)
		}
		return allReq, QuorumNotPossible
	}
	// consensus possible in the future
	logger.Trace("Consensus possible in the future, hash count: %v, needed count: %v", reqCount, quorum)
	return nil, QuorumInProgress
}
