package forwarder

import (
	"testing"
	"time"

	"gitdc.ee.guardtime.com/alphabill/alphabill/internal/protocol"

	"gitdc.ee.guardtime.com/alphabill/alphabill/internal/network"
	testnetwork "gitdc.ee.guardtime.com/alphabill/alphabill/internal/testutils/network"
	testtransaction "gitdc.ee.guardtime.com/alphabill/alphabill/internal/testutils/transaction"
	"gitdc.ee.guardtime.com/alphabill/alphabill/internal/txsystem"
	golog "github.com/ipfs/go-log"
	"github.com/libp2p/go-libp2p-core/peerstore"
	"github.com/stretchr/testify/require"
)

var (
	transfer = testtransaction.RandomBillTransfer()
	split    = testtransaction.RandomBillSplit()
)

func init() {
	golog.SetAllLoggers(golog.LevelWarn) // change this to Debug if libp2p logs are needed
}

/*
func TestNew_PeerIsNil(t *testing.T) {
	_, err := New(nil, DefaultForwardingTimeout, func(tx *txsystem.Transaction) {})
	require.ErrorIs(t, err, ErrPeerIsNil)
}

func TestNew_TxHandlerIsNil(t *testing.T) {
	_, err := New(testnetwork.CreatePeer(t), DefaultForwardingTimeout, nil)
	require.ErrorIs(t, err, ErrReceivedMessageChIsNil)
}
*/
func TestTxHandler_ForwardTx(t *testing.T) {
	// init peer1
	peer1 := testnetwork.CreatePeer(t)
	defer peer1.Close()

	// init peer1
	peer2 := testnetwork.CreatePeer(t)
	defer peer2.Close()

	// init peerstores
	peer1.Network().Peerstore().AddAddrs(peer2.ID(), peer2.MultiAddresses(), peerstore.PermanentAddrTTL)
	peer2.Network().Peerstore().AddAddrs(peer1.ID(), peer1.MultiAddresses(), peerstore.PermanentAddrTTL)

	// init peer1 forwarder
	outCh := make(chan network.ReceivedMessage)
	defer close(outCh)
	_, err := protocol.NewReceiverProtocol[*txsystem.Transaction](peer1, ProtocolInputForward, outCh, func() *txsystem.Transaction { return &txsystem.Transaction{} })
	require.NoError(t, err)

	// init peer2 forwarder
	outCh2 := make(chan<- network.ReceivedMessage)
	defer close(outCh2)

	var peer2Tx *txsystem.Transaction
	_, err = protocol.NewReceiverProtocol[*txsystem.Transaction](peer2, ProtocolInputForward, outCh2, func() *txsystem.Transaction { return &txsystem.Transaction{} })
	require.NoError(t, err)

	peer2Sender, err := protocol.NewSendProtocol(peer2, ProtocolInputForward, time.Second)
	require.NoError(t, err)

	// peer2 forwards tx to peer1
	err = peer2Sender.Send(transfer, peer1.ID())
	require.NoError(t, err)

	peer1Tx := <-outCh
	require.Nil(t, peer2Tx)

	require.NotNil(t, peer1Tx)

	time.Sleep(time.Second)
	// peer1 forward tx to peer2
	//err = peer1Forwarder.Forward(split, peer2.ID())
	//require.NoError(t, err)
	//require.NotNil(t, peer2Tx)
}

/*
func TestTxHandler_UnknownPeer(t *testing.T) {
	// init peer1
	peer1 := testnetwork.CreatePeer(t)
	defer peer1.Close()

	// init peer1
	peer2 := testnetwork.CreatePeer(t)
	defer peer2.Close()

	// init peer2 forwarder
	peer2Forwarder, err := New(peer2, DefaultForwardingTimeout, func(tx *txsystem.Transaction) {})
	require.NoError(t, err)
	err = peer2Forwarder.Forward(split, peer1.ID())
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "failed to open stream"))
}

func TestTxHandler_PeerIsClosed(t *testing.T) {
	// init peer1
	peer1 := testnetwork.CreatePeer(t)
	defer peer1.Close()

	// init peer1
	peer2 := testnetwork.CreatePeer(t)
	defer peer2.Close()

	// init peerstores
	peer1.Network().Peerstore().AddAddrs(peer2.ID(), peer2.MultiAddresses(), peerstore.PermanentAddrTTL)
	peer2.Network().Peerstore().AddAddrs(peer1.ID(), peer1.MultiAddresses(), peerstore.PermanentAddrTTL)
	require.NoError(t, peer1.Close())

	// init peer2 forwarder
	peer2Forwarder, err := New(peer2, DefaultForwardingTimeout, func(tx *txsystem.Transaction) {})
	require.NoError(t, err)

	// peer2 forwards tx to peer1
	err = peer2Forwarder.Forward(transfer, peer1.ID())
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "connection refused"))
}

func TestTxHandler_Timeout(t *testing.T) {
	// init peer1
	peer1 := testnetwork.CreatePeer(t)
	defer peer1.Close()

	// init peer1
	peer2 := testnetwork.CreatePeer(t)
	defer peer2.Close()

	// init peerstores
	peer1.Network().Peerstore().AddAddrs(peer2.ID(), peer2.MultiAddresses(), peerstore.PermanentAddrTTL)
	peer2.Network().Peerstore().AddAddrs(peer1.ID(), peer1.MultiAddresses(), peerstore.PermanentAddrTTL)

	// init peer1 forwarder
	_, err := New(peer1, DefaultForwardingTimeout, func(tx *txsystem.Transaction) {
		time.Sleep(time.Second)
	})
	require.NoError(t, err)

	// init peer2 forwarder
	peer2Forwarder, err := New(peer2, time.Millisecond, func(tx *txsystem.Transaction) {
	})
	require.NoError(t, err)

	// peer2 forwards tx to peer1
	err = peer2Forwarder.Forward(transfer, peer1.ID())
	require.ErrorIs(t, err, ErrTimout)

}
*/
