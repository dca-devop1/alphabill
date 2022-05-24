package testpartition

import (
	"testing"

	test "gitdc.ee.guardtime.com/alphabill/alphabill/internal/testutils"
	testtransaction "gitdc.ee.guardtime.com/alphabill/alphabill/internal/testutils/transaction"
	"gitdc.ee.guardtime.com/alphabill/alphabill/internal/txsystem"
	"github.com/stretchr/testify/require"
)

var systemIdentifier = []byte{1, 2, 4, 1}

func TestNewNetwork_Ok(t *testing.T) {
	network, err := NewNetwork(3, func() txsystem.TransactionSystem {
		return &CounterTxSystem{}
	}, systemIdentifier)
	require.NoError(t, err)
	defer func() {
		err = network.Close()
		require.NoError(t, err)
	}()
	require.NotNil(t, network.RootChain)
	require.Equal(t, 3, len(network.Nodes))

	tx := randomTx(systemIdentifier)
	require.NoError(t, network.SubmitTx(tx))
	require.Eventually(t, BlockchainContainsTx(tx, network), test.WaitDuration, test.WaitTick)

	tx = randomTx(systemIdentifier)
	err = network.BroadcastTx(tx)
	require.Eventually(t, BlockchainContainsTx(tx, network), test.WaitDuration, test.WaitTick)
}

func randomTx(systemIdentifier []byte) *txsystem.Transaction {
	tx := testtransaction.RandomBillTransfer()
	tx.SystemId = systemIdentifier
	tx.Timeout = 100
	return tx
}