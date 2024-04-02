package partition

import (
	"bytes"
	"context"
	"crypto"
	"errors"
	"io"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	test "github.com/alphabill-org/alphabill/internal/testutils"
	testlogger "github.com/alphabill-org/alphabill/internal/testutils/logger"
	"github.com/alphabill-org/alphabill/keyvaluedb/boltdb"
	"github.com/alphabill-org/alphabill/keyvaluedb/memorydb"
	"github.com/alphabill-org/alphabill/state"
	"github.com/alphabill-org/alphabill/types"
)

func TestNewProofIndexer_history_2(t *testing.T) {
	proofDB, err := memorydb.New()
	require.NoError(t, err)
	logger := testlogger.New(t)
	indexer := NewProofIndexer(crypto.SHA256, proofDB, 2, logger)
	require.Equal(t, proofDB, indexer.GetDB())

	// start indexing loop
	ctx := context.Background()
	unitID := make([]byte, 32)
	blockRound1 := simulateInput(1, unitID)
	require.NoError(t, indexer.create(ctx, blockRound1.Block, blockRound1.State))
	blockRound2 := simulateInput(2, unitID)
	require.NoError(t, indexer.create(ctx, blockRound2.Block, blockRound2.State))

	// add the same block again
	require.ErrorContains(t, indexer.create(ctx, blockRound2.Block, blockRound2.State), "block 2 already indexed")
	blockRound3 := simulateInput(3, unitID)
	require.NoError(t, indexer.create(ctx, blockRound3.Block, blockRound3.State))

	// run clean-up
	require.NoError(t, indexer.historyCleanup(ctx, 3))
	require.EqualValues(t, 3, indexer.latestIndexedBlockNumber())

	// verify round number 1 is correctly cleaned up
	for _, txr := range blockRound1.Block.Transactions {
		// verify tx index is not deleted
		txoHash := txr.TransactionOrder.Hash(crypto.SHA256)
		var index *TxIndex
		f, err := proofDB.Read(txoHash, &index)
		require.NoError(t, err)
		require.True(t, f)

		// verify unit proofs are deleted
		var unitProof *types.UnitDataAndProof
		unitProofKey := bytes.Join([][]byte{unitID, txoHash}, nil)
		f, err = proofDB.Read(unitProofKey, &unitProof)
		require.NoError(t, err)
		require.False(t, f)
	}
}

func TestNewProofIndexer_NothingIsWrittenIfBlockIsEmpty(t *testing.T) {
	proofDB, err := memorydb.New()
	require.NoError(t, err)
	logger := testlogger.New(t)
	indexer := NewProofIndexer(crypto.SHA256, proofDB, 2, logger)
	require.Equal(t, proofDB, indexer.GetDB())
	// start indexing loop
	ctx := context.Background()
	blockRound1 := simulateEmptyInput(1)
	require.NoError(t, indexer.create(ctx, blockRound1.Block, blockRound1.State))
	blockRound2 := simulateEmptyInput(2)
	require.NoError(t, indexer.create(ctx, blockRound2.Block, blockRound2.State))
	// add the same block again
	require.ErrorContains(t, indexer.create(ctx, blockRound2.Block, blockRound2.State), "block 2 already indexed")
	blockRound3 := simulateEmptyInput(3)
	require.NoError(t, indexer.create(ctx, blockRound3.Block, blockRound3.State))
	// run clean-up
	require.NoError(t, indexer.historyCleanup(ctx, 3))
	require.EqualValues(t, 3, indexer.latestIndexedBlockNumber())
	// index db contains only latest round number
	dbIt := proofDB.First()
	require.True(t, dbIt.Valid())
	require.Equal(t, keyLatestRoundNumber, dbIt.Key())
	dbIt.Next()
	require.False(t, dbIt.Valid())
	require.NoError(t, dbIt.Close())
}

func TestNewProofIndexer_RunLoop(t *testing.T) {
	t.Run("run loop - no history clean-up", func(t *testing.T) {
		proofDB, err := memorydb.New()
		require.NoError(t, err)
		logger := testlogger.New(t)
		indexer := NewProofIndexer(crypto.SHA256, proofDB, 0, logger)
		// start indexing loop
		ctx := context.Background()
		nctx, cancel := context.WithCancel(ctx)
		done := make(chan error)
		t.Cleanup(func() {
			cancel()
			select {
			case err := <-done:
				require.ErrorIs(t, err, context.Canceled)
			case <-time.After(200 * time.Millisecond):
				t.Error("indexer loop did not exit in time")
			}
		})
		// start loop
		go func(dn chan error) {
			done <- indexer.loop(nctx)
		}(done)
		unitID := make([]byte, 32)
		blockRound1 := simulateInput(1, unitID)
		stateMock := mockStateStoreOK{}
		indexer.Handle(nctx, blockRound1.Block, stateMock)
		blockRound2 := simulateInput(2, unitID)
		indexer.Handle(nctx, blockRound2.Block, stateMock)
		blockRound3 := simulateInput(3, unitID)
		indexer.Handle(nctx, blockRound3.Block, stateMock)
		require.Eventually(t, func() bool {
			return indexer.latestIndexedBlockNumber() == 3
		}, test.WaitDuration, test.WaitTick)
		// verify history for round 1 is cleaned up
		for _, transaction := range blockRound1.Block.Transactions {
			oderHash := transaction.TransactionOrder.Hash(crypto.SHA256)
			index := &struct {
				RoundNumber  uint64
				TxOrderIndex int
			}{}
			f, err := proofDB.Read(oderHash, index)
			require.NoError(t, err)
			require.True(t, f)
		}
	})
	t.Run("run loop - keep last 2 rounds", func(t *testing.T) {
		proofDB, err := memorydb.New()
		require.NoError(t, err)
		logger := testlogger.New(t)
		indexer := NewProofIndexer(crypto.SHA256, proofDB, 2, logger)

		// start indexing loop
		ctx := context.Background()
		nctx, cancel := context.WithCancel(ctx)
		done := make(chan error)
		t.Cleanup(func() {
			cancel()
			select {
			case err := <-done:
				require.ErrorIs(t, err, context.Canceled)
			case <-time.After(200 * time.Millisecond):
				t.Error("indexer loop did not exit in time")
			}
		})
		go func(dn chan error) {
			done <- indexer.loop(nctx)
		}(done)
		unitID := make([]byte, 32)
		blockRound1 := simulateInput(1, unitID)
		stateMock := mockStateStoreOK{}
		indexer.Handle(nctx, blockRound1.Block, stateMock)
		blockRound2 := simulateInput(2, unitID)
		indexer.Handle(nctx, blockRound2.Block, stateMock)
		blockRound3 := simulateInput(3, unitID)
		indexer.Handle(nctx, blockRound3.Block, stateMock)
		require.Eventually(t, func() bool {
			return indexer.latestIndexedBlockNumber() == 3
		}, test.WaitDuration, test.WaitTick)

		// verify history for round 1 is correctly cleaned up
		for _, transaction := range blockRound1.Block.Transactions {
			// verify tx index is not deleted
			txoHash := transaction.TransactionOrder.Hash(crypto.SHA256)
			var index *TxIndex
			f, err := proofDB.Read(txoHash, &index)
			require.NoError(t, err)
			require.True(t, f)

			// verify unit proofs are deleted
			var unitProof *types.UnitDataAndProof
			unitProofKey := bytes.Join([][]byte{unitID, txoHash}, nil)
			f, err = proofDB.Read(unitProofKey, &unitProof)
			require.NoError(t, err)
			require.False(t, f)
		}
	})
}

func TestProofIndexer_BoltDBTx(t *testing.T) {
	proofDB, err := boltdb.New(filepath.Join(t.TempDir(), "tempdb.db"))
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = proofDB.Close()
	})
	logger := testlogger.New(t)
	indexer := NewProofIndexer(crypto.SHA256, proofDB, 2, logger)

	// simulate error when indexing a block
	ctx := context.Background()
	bas := simulateInput(1, []byte{1})
	bas.State = mockStateStoreOK{err: errors.New("some error")}

	err = indexer.create(ctx, bas.Block, bas.State)
	require.ErrorContains(t, err, "some error")

	// verify index db does not contain the stored round number (tx is rolled back)
	dbIt := proofDB.First()
	t.Cleanup(func() {
		_ = dbIt.Close()
	})
	require.False(t, dbIt.Valid())
}

type mockStateStoreOK struct {
	err error
}

func (m mockStateStoreOK) GetUnit(id types.UnitID, committed bool) (*state.Unit, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &state.Unit{}, nil
}

func (m mockStateStoreOK) CreateUnitStateProof(id types.UnitID, logIndex int) (*types.UnitStateProof, error) {
	return &types.UnitStateProof{}, nil
}

func (m mockStateStoreOK) CreateIndex(state.KeyExtractor[string]) (state.Index[string], error) {
	return nil, nil
}

func (m mockStateStoreOK) Serialize(writer io.Writer, committed bool) error {
	return nil
}

func simulateInput(round uint64, unitID []byte) *BlockAndState {
	block := &types.Block{
		Header: &types.Header{SystemID: 1},
		Transactions: []*types.TransactionRecord{
			{
				TransactionOrder: &types.TransactionOrder{},
				ServerMetadata:   &types.ServerMetadata{TargetUnits: []types.UnitID{unitID}},
			},
		},
		UnicityCertificate: &types.UnicityCertificate{
			InputRecord: &types.InputRecord{RoundNumber: round},
		},
	}
	return &BlockAndState{
		Block: block,
		State: mockStateStoreOK{},
	}
}

func simulateEmptyInput(round uint64) *BlockAndState {
	block := &types.Block{
		Header:       &types.Header{SystemID: 1},
		Transactions: []*types.TransactionRecord{},
		UnicityCertificate: &types.UnicityCertificate{
			InputRecord: &types.InputRecord{RoundNumber: round},
		},
	}
	return &BlockAndState{
		Block: block,
		State: mockStateStoreOK{},
	}
}
