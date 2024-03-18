package main

import (
	"bytes"
	"context"
	"crypto"
	"errors"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/alphabill-org/alphabill/hash"
	"github.com/alphabill-org/alphabill/predicates/templates"
	"github.com/alphabill-org/alphabill/rpc/alphabill"
	"github.com/alphabill-org/alphabill/txsystem/fc/transactions"
	"github.com/alphabill-org/alphabill/txsystem/money"
	"github.com/alphabill-org/alphabill/types"
	"github.com/alphabill-org/alphabill/util"
)

/*
Example usage
go run scripts/money/spend_initial_bill.go --pubkey 0x03c30573dc0c7fd43fcb801289a6a96cb78c27f4ba398b89da91ece23e9a99aca3 --alphabill-uri localhost:26766 --bill-id 1 --bill-value 1000000000000000000 --timeout 10
*/
func main() {
	// parse command line parameters
	pubKeyHex := flag.String("pubkey", "", "public key of the new bill owner")
	billIdUint := flag.Uint64("bill-id", 0, "bill id of the spendable bill")
	billValue := flag.Uint64("bill-value", 0, "bill value of the spendable bill")
	timeout := flag.Uint64("timeout", 0, "transaction timeout (block number)")
	uri := flag.String("alphabill-uri", "", "alphabill node uri where to send the transaction")
	flag.Parse()

	// verify command line parameters
	if *pubKeyHex == "" {
		log.Fatal("pubkey is required")
	}
	if *billIdUint == 0 {
		log.Fatal("bill-id is required")
	}
	if *billValue == 0 {
		log.Fatal("bill-value is required")
	}
	if *timeout == 0 {
		log.Fatal("timeout is required")
	}
	if *uri == "" {
		log.Fatal("alphabill-uri is required")
	}

	// process command line parameters
	pubKey, err := hexutil.Decode(*pubKeyHex)
	if err != nil {
		log.Fatal(err)
	}

	billID := money.NewBillID(nil, util.Uint64ToBytes(*billIdUint))

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, *uri, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	txClient := alphabill.NewAlphabillServiceClient(conn)
	err = execInitialBill(ctx, txClient, *timeout, billID, *billValue, pubKey)
	if err != nil {
		log.Fatal(err)
	}
}

func execInitialBill(ctx context.Context, client alphabill.AlphabillServiceClient, timeout uint64, billID types.UnitID, billValue uint64, pubKey []byte) error {
	res, err := client.GetRoundNumber(ctx, &emptypb.Empty{})
	if err != nil {
		return fmt.Errorf("getting round number: %w", err)
	}
	absoluteTimeout := res.RoundNumber + timeout

	txFee := uint64(1)
	feeAmount := uint64(2)
	// Make the initial fcrID different from the default
	// sha256(pubKey), so that wallet can later create it's own
	// fcrID for the same account with a different owner condition
	fcrID := money.NewFeeCreditRecordID(billID, hash.Sum256(hash.Sum256(pubKey)))

	// create transferFC
	transferFC, err := createTransferFC(feeAmount+txFee, billID, fcrID, res.RoundNumber, absoluteTimeout)
	if err != nil {
		return fmt.Errorf("creating transfer FC transaction: %w", err)
	}
	transferFCBytes, err := types.Cbor.Marshal(transferFC)
	if err != nil {
		return fmt.Errorf("marshalling transfer FC transaction: %w", err)
	}
	protoTransferFC := &alphabill.Transaction{Order: transferFCBytes}

	// send transferFC
	_, err = client.ProcessTransaction(ctx, protoTransferFC)
	if err != nil {
		return fmt.Errorf("processing transfer FC transaction: %w", err)
	}
	log.Println("sent transferFC transaction")

	// wait for transferFC proof
	transferFCProof, err := waitForConfirmation(ctx, client, transferFC, res.RoundNumber, absoluteTimeout)
	if err != nil {
		return fmt.Errorf("failed to confirm transferFC transaction %v", err)
	} else {
		log.Println("confirmed transferFC transaction")
	}

	// create addFC
	addFC, err := createAddFC(fcrID, templates.AlwaysTrueBytes(), transferFCProof.TxRecord, transferFCProof.TxProof, absoluteTimeout, feeAmount)
	if err != nil {
		return fmt.Errorf("creating add FC transaction: %w", err)
	}
	addFCBytes, err := types.Cbor.Marshal(addFC)
	if err != nil {
		return fmt.Errorf("marshalling add FC transaction: %w", err)
	}
	protoAddFC := &alphabill.Transaction{Order: addFCBytes}

	// send addFC
	_, err = client.ProcessTransaction(ctx, protoAddFC)
	if err != nil {
		return fmt.Errorf("processing add FC transaction: %w", err)
	}
	log.Println("sent addFC transaction")

	// wait for addFC confirmation
	_, err = waitForConfirmation(ctx, client, addFC, res.RoundNumber, absoluteTimeout)
	if err != nil {
		return fmt.Errorf("failed to confirm addFC transaction %v", err)
	} else {
		log.Println("confirmed addFC transaction")
	}

	// create transfer tx
	tx, err := createTransferTx(pubKey, billID, billValue-feeAmount-txFee, fcrID, absoluteTimeout, transferFC.Hash(crypto.SHA256))
	if err != nil {
		return fmt.Errorf("creating transfer transaction: %w", err)
	}
	txBytes, err := types.Cbor.Marshal(tx)
	if err != nil {
		return fmt.Errorf("marshalling transfer transaction: %w", err)
	}

	// send transfer tx
	protoTransferTx := &alphabill.Transaction{Order: txBytes}
	if _, err := client.ProcessTransaction(ctx, protoTransferTx); err != nil {
		return fmt.Errorf("processing transfer transaction: %w", err)
	}
	log.Println("successfully sent initial bill transfer transaction")

	return nil
}

func createTransferFC(feeAmount uint64, unitID []byte, targetUnitID []byte, t1, t2 uint64) (*types.TransactionOrder, error) {
	attr, err := types.Cbor.Marshal(
		&transactions.TransferFeeCreditAttributes{
			Amount:                 feeAmount,
			TargetSystemIdentifier: 1,
			TargetRecordID:         targetUnitID,
			EarliestAdditionTime:   t1,
			LatestAdditionTime:     t2,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transferFC attributes: %w", err)
	}
	tx := &types.TransactionOrder{
		Payload: &types.Payload{
			SystemID:       1,
			Type:           transactions.PayloadTypeTransferFeeCredit,
			UnitID:         unitID,
			Attributes:     attr,
			ClientMetadata: &types.ClientMetadata{Timeout: t2, MaxTransactionFee: 1},
		},
		OwnerProof: nil,
	}
	return tx, nil
}

func createAddFC(unitID []byte, ownerCondition []byte, transferFC *types.TransactionRecord, transferFCProof *types.TxProof, timeout uint64, maxFee uint64) (*types.TransactionOrder, error) {
	attr, err := types.Cbor.Marshal(
		&transactions.AddFeeCreditAttributes{
			FeeCreditTransfer:       transferFC,
			FeeCreditTransferProof:  transferFCProof,
			FeeCreditOwnerCondition: ownerCondition,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transferFC attributes: %w", err)
	}
	return &types.TransactionOrder{
		Payload: &types.Payload{
			SystemID:       1,
			Type:           transactions.PayloadTypeAddFeeCredit,
			UnitID:         unitID,
			Attributes:     attr,
			ClientMetadata: &types.ClientMetadata{Timeout: timeout, MaxTransactionFee: maxFee},
		},
		OwnerProof: nil,
	}, nil
}

func createTransferTx(pubKey []byte, unitID []byte, billValue uint64, fcrID []byte, timeout uint64, backlink []byte) (*types.TransactionOrder, error) {
	attr, err := types.Cbor.Marshal(
		&money.TransferAttributes{
			NewBearer:   templates.NewP2pkh256BytesFromKeyHash(hash.Sum256(pubKey)),
			TargetValue: billValue,
			Backlink:    backlink,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transferFC attributes: %w", err)
	}
	return &types.TransactionOrder{
		Payload: &types.Payload{
			SystemID:   1,
			Type:       money.PayloadTypeTransfer,
			UnitID:     unitID,
			Attributes: attr,
			ClientMetadata: &types.ClientMetadata{
				Timeout:           timeout,
				MaxTransactionFee: 1,
				FeeCreditRecordID: fcrID,
			},
		},
		OwnerProof: nil,
	}, nil
}

func waitForConfirmation(ctx context.Context, abClient alphabill.AlphabillServiceClient, pendingTx *types.TransactionOrder, latestRoundNumber, timeout uint64) (*Proof, error) {
	for latestRoundNumber <= timeout {
		res, err := abClient.GetBlock(ctx, &alphabill.GetBlockRequest{BlockNo: latestRoundNumber})
		if err != nil {
			return nil, err
		}
		blockBytes := res.Block
		if blockBytes == nil || (len(blockBytes) == 1 && blockBytes[0] == 0xf6) { // 0xf6 cbor Null
			// block might be empty, check latest round number
			res, err := abClient.GetRoundNumber(ctx, &emptypb.Empty{})
			if err != nil {
				return nil, err
			}
			if res.RoundNumber > latestRoundNumber {
				latestRoundNumber++
			} else {
				// wait for some time before retrying to fetch new block
				select {
				case <-time.After(time.Second):
					continue
				case <-ctx.Done():
					return nil, nil
				}
			}
		} else {
			block := &types.Block{}
			if err := types.Cbor.Unmarshal(blockBytes, block); err != nil {
				return nil, fmt.Errorf("failed to unmarshal block: %w", err)
			}
			for i, tx := range block.Transactions {
				if bytes.Equal(tx.TransactionOrder.UnitID(), pendingTx.UnitID()) {
					return NewTxProof(i, block, crypto.SHA256)
				}
			}
			latestRoundNumber++
		}
	}
	return nil, errors.New("error tx failed to confirm")
}

// Proof wrapper struct around TxRecord and TxProof
type Proof struct {
	_        struct{}                 `cbor:",toarray"`
	TxRecord *types.TransactionRecord `json:"txRecord"`
	TxProof  *types.TxProof           `json:"txProof"`
}

func NewTxProof(txIdx int, b *types.Block, hashAlgorithm crypto.Hash) (*Proof, error) {
	txProof, txRecord, err := types.NewTxProof(b, txIdx, hashAlgorithm)
	if err != nil {
		return nil, err
	}
	return &Proof{
		TxRecord: txRecord,
		TxProof:  txProof,
	}, nil
}
