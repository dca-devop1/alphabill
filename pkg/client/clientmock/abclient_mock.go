package clientmock

import (
	"context"

	"github.com/alphabill-org/alphabill/internal/block"
	"github.com/alphabill-org/alphabill/internal/rpc/alphabill"
	"github.com/alphabill-org/alphabill/internal/txsystem"
)

// MockAlphabillClient for testing. NOT thread safe.
type (
	MockAlphabillClient struct {
		recordedTxs              []*txsystem.Transaction
		txResponse               error
		maxBlockNumber           uint64
		maxRoundNumber           uint64
		shutdown                 bool
		blocks                   map[uint64]*block.Block
		txListener               func(tx *txsystem.Transaction)
		incrementOnFetch         bool // if true, maxBlockNumber will be incremented on each GetBlocks call
		lastRequestedBlockNumber uint64
	}
	Option func(c *MockAlphabillClient)
)

func NewMockAlphabillClient(options ...Option) *MockAlphabillClient {
	mockClient := &MockAlphabillClient{blocks: map[uint64]*block.Block{}}
	for _, o := range options {
		o(mockClient)
	}
	return mockClient
}

func WithMaxBlockNumber(blockNumber uint64) Option {
	return func(c *MockAlphabillClient) {
		c.SetMaxBlockNumber(blockNumber)
	}
}

func WithMaxRoundNumber(roundNumber uint64) Option {
	return func(c *MockAlphabillClient) {
		c.maxRoundNumber = roundNumber
	}
}

func WithBlocks(blocks map[uint64]*block.Block) Option {
	return func(c *MockAlphabillClient) {
		c.blocks = blocks
	}
}

func (c *MockAlphabillClient) SendTransaction(ctx context.Context, tx *txsystem.Transaction) error {
	c.recordedTxs = append(c.recordedTxs, tx)
	if c.txListener != nil {
		c.txListener(tx)
	}
	return c.txResponse
}

func (c *MockAlphabillClient) GetBlock(ctx context.Context, blockNumber uint64) (*block.Block, error) {
	if c.incrementOnFetch {
		defer c.SetMaxBlockNumber(blockNumber + 1)
	}
	if c.blocks != nil {
		b := c.blocks[blockNumber]
		return b, nil
	}
	return nil, nil
}

func (c *MockAlphabillClient) GetBlocks(ctx context.Context, blockNumber, blockCount uint64) (*alphabill.GetBlocksResponse, error) {
	c.lastRequestedBlockNumber = blockNumber
	if c.incrementOnFetch {
		defer c.SetMaxBlockNumber(blockNumber + 1)
	}
	batchMaxBlockNumber := blockNumber
	if blockNumber <= c.maxBlockNumber {
		var blocks []*block.Block
		b, f := c.blocks[blockNumber]
		if f {
			blocks = []*block.Block{b}
			batchMaxBlockNumber = b.UnicityCertificate.InputRecord.RoundNumber
		} else {
			blocks = []*block.Block{}
		}
		return &alphabill.GetBlocksResponse{
			MaxBlockNumber:      c.maxBlockNumber,
			MaxRoundNumber:      c.maxRoundNumber,
			Blocks:              blocks,
			BatchMaxBlockNumber: batchMaxBlockNumber,
		}, nil
	}
	return &alphabill.GetBlocksResponse{
		MaxBlockNumber:      c.maxBlockNumber,
		MaxRoundNumber:      c.maxRoundNumber,
		Blocks:              []*block.Block{},
		BatchMaxBlockNumber: batchMaxBlockNumber,
	}, nil
}

func (c *MockAlphabillClient) GetRoundNumber(ctx context.Context) (uint64, error) {
	return c.maxRoundNumber, nil
}

func (c *MockAlphabillClient) Shutdown() error {
	c.shutdown = true
	return nil
}

func (c *MockAlphabillClient) IsShutdown() bool {
	return c.shutdown
}

func (c *MockAlphabillClient) SetTxResponse(txResponse error) {
	c.txResponse = txResponse
}

func (c *MockAlphabillClient) SetMaxBlockNumber(blockNumber uint64) {
	c.maxBlockNumber = blockNumber
	if c.maxBlockNumber > c.maxRoundNumber {
		c.SetMaxRoundNumber(c.maxBlockNumber)
	}
}

func (c *MockAlphabillClient) SetMaxRoundNumber(roundNumber uint64) {
	if c.maxBlockNumber > roundNumber {
		panic("round number cannot be behind the block number")
	}
	c.maxRoundNumber = roundNumber
}

func (c *MockAlphabillClient) SetBlock(b *block.Block) {
	c.blocks[b.UnicityCertificate.InputRecord.RoundNumber] = b
}

func (c *MockAlphabillClient) GetRecordedTransactions() []*txsystem.Transaction {
	return c.recordedTxs
}

func (c *MockAlphabillClient) ClearRecordedTransactions() {
	c.recordedTxs = make([]*txsystem.Transaction, 0)
}

func (c *MockAlphabillClient) SetTxListener(txListener func(tx *txsystem.Transaction)) {
	c.txListener = txListener
}

func (c *MockAlphabillClient) SetIncrementOnFetch(incrementOnFetch bool) {
	c.incrementOnFetch = incrementOnFetch
}

func (c *MockAlphabillClient) GetLastRequestedBlockNumber() uint64 {
	return c.lastRequestedBlockNumber
}
