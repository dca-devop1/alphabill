package forwarder

import (
	"context"
	"fmt"
	"time"

	"gitdc.ee.guardtime.com/alphabill/alphabill/internal/rpc/transaction"

	"gitdc.ee.guardtime.com/alphabill/alphabill/internal/errors"
	"gitdc.ee.guardtime.com/alphabill/alphabill/internal/errors/errstr"
	log "gitdc.ee.guardtime.com/alphabill/alphabill/internal/logger"
	"gitdc.ee.guardtime.com/alphabill/alphabill/internal/network"
	"gitdc.ee.guardtime.com/alphabill/alphabill/internal/protocol"
	"gitdc.ee.guardtime.com/alphabill/alphabill/internal/txbuffer"
	libp2pNetwork "github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
)

var logger = log.CreateForPackage()

const (
	// ProtocolIdTxForwarder is the protocol.ID of the AlphaBill transaction forwarding protocol.
	ProtocolIdTxForwarder            = "/ab/tx/1.0.0"
	DefaultForwardingTimeout         = 2 * time.Second
	UnknownLeader            peer.ID = ""
)

type (
	// LeaderSelector interface is used to get the next leader.
	LeaderSelector interface {
		// NextLeader returns the identifier of the next leader.
		NextLeader() (peer.ID, error)
	}

	// TxForwarder sends transactions, as they arrive, to the expected next leader.
	TxForwarder struct {
		leaderSelector LeaderSelector
		txBuffer       *txbuffer.TxBuffer
		self           *network.Peer
	}
)

// New constructs a new *TxForwarder and activates it by attaching its stream handler to the given network.Peer.
func New(self *network.Peer, leaderSelector LeaderSelector, txBuffer *txbuffer.TxBuffer) (*TxForwarder, error) {
	if self == nil {
		return nil, errors.New(errstr.NilArgument)
	}
	if leaderSelector == nil {
		return nil, errors.New(errstr.NilArgument)
	}
	if txBuffer == nil {
		return nil, errors.New(errstr.NilArgument)
	}
	tf := &TxForwarder{
		txBuffer:       txBuffer,
		leaderSelector: leaderSelector,
		self:           self,
	}
	self.RegisterProtocolHandler(ProtocolIdTxForwarder, tf.handleStream)
	return tf, nil
}

// Handle handles the incoming transaction. If current node isn't the leader then the transaction is forwarded to the
// expected next leader. If current node is the leader then the transaction is added the txbuffer.TxBuffer.
func (tf *TxForwarder) Handle(ctx context.Context, req *transaction.Transaction) error {
	nextLeader, err := tf.leaderSelector.NextLeader()
	if err != nil {
		return err
	}
	if nextLeader == UnknownLeader || nextLeader == tf.self.ID() {
		// leader is unknown or the current node is the leader
		return tf.handleTx(req)
	}
	// forward transaction to the leader
	return tf.forwardTx(ctx, req, nextLeader)
}

// Close shuts down the TxForwarder.
func (tf *TxForwarder) Close() error {
	tf.self.RemoveProtocolHandler(ProtocolIdTxForwarder)
	return nil
}

// forwardTx forwards the transaction to the receiver.
func (tf *TxForwarder) forwardTx(ctx context.Context, req *transaction.Transaction, receiver peer.ID) error {
	ctx, cancel := context.WithTimeout(ctx, DefaultForwardingTimeout)
	defer cancel()

	s, err := tf.self.CreateStream(ctx, receiver, ProtocolIdTxForwarder)
	if err != nil {
		return err
	}
	defer s.Close()
	w := protocol.NewProtoBufWriter(s)
	if err := w.Write(req); err != nil {
		_ = s.Reset()
		return fmt.Errorf("failed to forward transaction, %w", err)
	}
	logger.Debug("forwarded tx to peer %v", receiver)
	return nil
}

// handleStream receives incoming transactions from other peers in the network.
func (tf *TxForwarder) handleStream(s libp2pNetwork.Stream) {
	r := protocol.NewProtoBufReader(s)
	defer r.Close()

	req := &transaction.Transaction{}
	err := r.Read(req)
	if err != nil {
		logger.Warning("Failed to read the transaction: %v", err)
		return
	}
	// TODO some transaction must be included to the txbuffer
	logger.Debug("Got a new transaction %v", req)
	nextLeader, err := tf.leaderSelector.NextLeader()
	if err != nil {
		logger.Warning("Ignoring tx: %v", err)
		return
	}
	if nextLeader != tf.self.ID() {
		logger.Warning("Ignoring tx. Current node isn't the next leader.")
		return
	}
	err = tf.handleTx(req)
	if err != nil {
		logger.Warning("Transaction was not added to the TxBuffer: %v", err)
	}
}

func (tf *TxForwarder) handleTx(req *transaction.Transaction) error {
	genericTx, err := transaction.New(req)
	if err != nil {
		return err
	}
	return tf.txBuffer.Add(genericTx)
}
