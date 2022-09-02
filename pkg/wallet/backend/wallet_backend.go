package backend

import (
	"context"
	"fmt"
	"time"

	"github.com/alphabill-org/alphabill/internal/errors"
	"github.com/alphabill-org/alphabill/internal/hash"
	"github.com/alphabill-org/alphabill/internal/logger"
	"github.com/alphabill-org/alphabill/internal/proof"
	"github.com/alphabill-org/alphabill/internal/util"
	"github.com/alphabill-org/alphabill/pkg/client"
	"github.com/alphabill-org/alphabill/pkg/wallet"
	"github.com/holiman/uint256"
)

var (
	alphabillMoneySystemId = []byte{0, 0, 0, 0}
	log                    = logger.CreateForPackage()
)

type (
	WalletBackend struct {
		trackedPubKeys []*pubkey
		txProcessor    *txProcessor
		abclient       client.ABClient
		store          BillStore
	}

	bill struct {
		Id    *uint256.Int
		Value uint64
	}

	blockProof struct {
		BillId      []byte            `json:"billId"`
		BlockNumber uint64            `json:"blockNumber"`
		BlockProof  *proof.BlockProof `json:"blockProof"`
	}

	pubkey struct {
		pubkey     []byte
		pubkeyHash *wallet.KeyHashes
	}

	BillStore interface {
		GetBlockNumber() (uint64, error)
		SetBlockNumber(blockNumber uint64) error
		GetBills(pubKey []byte) ([]*bill, error)
		AddBill(pubKey []byte, bill *bill) error
		RemoveBill(pubKey []byte, id *uint256.Int) error
		ContainsBill(pubKey []byte, id *uint256.Int) (bool, error)
		GetBlockProof(billId []byte) (*blockProof, error)
		SetBlockProof(proof *blockProof) error
	}
)

func New(pubkeys [][]byte, abclient client.ABClient, store BillStore) *WalletBackend {
	var trackedPubKeys []*pubkey
	for _, pk := range pubkeys {
		trackedPubKeys = append(trackedPubKeys, &pubkey{
			pubkey: pk,
			pubkeyHash: &wallet.KeyHashes{
				Sha256: hash.Sum256(pk),
				Sha512: hash.Sum512(pk),
			},
		})
	}
	return &WalletBackend{trackedPubKeys: trackedPubKeys, abclient: abclient, txProcessor: newTxProcessor(store), store: store}
}

// Start downloading blocks and indexing bills by their owner's public key
func (w *WalletBackend) Start(ctx context.Context) error {
	blockNumber, _ := w.store.GetBlockNumber()
	for {
		select {
		case <-ctx.Done():
			log.Info("wallet backend canceled by user")
			return nil
		default:
			maxBlockNumber, _ := w.abclient.GetMaxBlockNumber()
			if blockNumber == maxBlockNumber {
				log.Debug("sleeping at max block height...")
				time.Sleep(time.Second)
				continue
			}

			batchSize := util.Min(100, maxBlockNumber-blockNumber)
			res, _ := w.abclient.GetBlocks(blockNumber+1, batchSize)
			for _, b := range res.Blocks {
				log.Info("processing block %d", b.BlockNumber)
				if blockNumber+1 != b.BlockNumber {
					return errors.New(fmt.Sprintf("invalid block number, got %d current %d", b.BlockNumber, blockNumber))
				}

				for _, tx := range b.Transactions {
					for i, pubKey := range w.trackedPubKeys {
						_ = w.txProcessor.processTx(tx, b, i, pubKey)
					}
				}
				w.store.SetBlockNumber(b.BlockNumber)
				blockNumber = b.BlockNumber
			}
		}
	}
}

// GetBills returns all bills for given public key
func (w *WalletBackend) GetBills(pubkey []byte) ([]*bill, error) {
	return w.store.GetBills(pubkey)
}

// GetBlockProof returns most recent proof for given unit id
func (w *WalletBackend) GetBlockProof(unitId []byte) (*blockProof, error) {
	return w.store.GetBlockProof(unitId)
}
