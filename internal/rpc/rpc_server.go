package rpc

import (
	"context"
	"fmt"

	"github.com/alphabill-org/alphabill/internal/block"
	"github.com/alphabill-org/alphabill/internal/errors"
	"github.com/alphabill-org/alphabill/internal/errors/errstr"
	"github.com/alphabill-org/alphabill/internal/rpc/alphabill"
	"github.com/alphabill-org/alphabill/internal/txsystem"
	"github.com/alphabill-org/alphabill/internal/util"
)

type (
	rpcServer struct {
		alphabill.UnimplementedAlphabillServiceServer
		node partitionNode
	}

	partitionNode interface {
		SubmitTx(tx *txsystem.Transaction) error
		GetBlock(blockNr uint64) (*block.Block, error)
		GetLatestBlock() *block.Block
	}
)

func NewRpcServer(node partitionNode) (*rpcServer, error) {
	if node == nil {
		return nil, errors.Wrap(errors.ErrInvalidArgument, errstr.NilArgument)
	}
	return &rpcServer{
		node: node,
	}, nil
}

func (r *rpcServer) ProcessTransaction(_ context.Context, tx *txsystem.Transaction) (*txsystem.TransactionResponse, error) {
	err := r.node.SubmitTx(tx)
	if err != nil {
		return &txsystem.TransactionResponse{
			Ok:      false,
			Message: err.Error(),
		}, err
	}
	return &txsystem.TransactionResponse{
		Ok:      true,
		Message: "",
	}, nil
}

func (r *rpcServer) GetBlock(_ context.Context, req *alphabill.GetBlockRequest) (*alphabill.GetBlockResponse, error) {
	b, err := r.node.GetBlock(req.BlockNo)
	if err != nil {
		return &alphabill.GetBlockResponse{ErrorMessage: err.Error()}, err
	}
	return &alphabill.GetBlockResponse{Block: b}, nil
}

func (r *rpcServer) GetMaxBlockNo(_ context.Context, req *alphabill.GetMaxBlockNoRequest) (*alphabill.GetMaxBlockNoResponse, error) {
	maxBlockNumber := r.node.GetLatestBlock().GetBlockNumber()
	return &alphabill.GetMaxBlockNoResponse{BlockNo: maxBlockNumber}, nil
}

func (r *rpcServer) GetBlocks(_ context.Context, req *alphabill.GetBlocksRequest) (*alphabill.GetBlocksResponse, error) {
	latestBlock := r.node.GetLatestBlock()
	err := verifyRequest(req, latestBlock.BlockNumber)
	if err != nil {
		return &alphabill.GetBlocksResponse{ErrorMessage: err.Error()}, err
	}
	maxBlockCount := util.Min(req.BlockCount, 100)
	maxAvailableBlockNumber := util.Min(req.BlockNumber+maxBlockCount-1, latestBlock.BlockNumber)
	batchSize := maxAvailableBlockNumber - req.BlockNumber
	res := make([]*block.Block, 0, batchSize)
	for blockNumber := req.BlockNumber; blockNumber <= maxAvailableBlockNumber; blockNumber++ {
		b, err := r.node.GetBlock(blockNumber)
		if err != nil {
			return nil, err
		}
		res = append(res, b)
	}
	return &alphabill.GetBlocksResponse{Blocks: res, MaxBlockNumber: latestBlock.BlockNumber}, nil
}

func verifyRequest(req *alphabill.GetBlocksRequest, latestBlockNumber uint64) error {
	if req.BlockNumber < 1 {
		return errors.New("block number cannot be less than one")
	}
	if req.BlockCount < 1 {
		return errors.New("block count cannot be less than one")
	}
	if req.BlockNumber > latestBlockNumber {
		return errors.New(fmt.Sprintf("block number cannot be larger than latest block number, got %d have %d", req.BlockNumber, latestBlockNumber))
	}
	return nil
}
