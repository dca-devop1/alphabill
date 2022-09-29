package testtransaction

import (
	"math/rand"
	"testing"

	"github.com/alphabill-org/alphabill/internal/block"
	"github.com/alphabill-org/alphabill/internal/hash"
	"github.com/alphabill-org/alphabill/internal/script"
	"github.com/alphabill-org/alphabill/internal/txsystem"
	moneytx "github.com/alphabill-org/alphabill/internal/txsystem/money"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/anypb"
)

var moneySystemId = []byte{0, 0, 0, 0}

func RandomGenericBillTransfer(t *testing.T) txsystem.GenericTransaction {
	tx := randomTx()

	bt := &moneytx.TransferOrder{
		NewBearer: randomBytes(3),
		// #nosec G404
		TargetValue: rand.Uint64(),
		Backlink:    randomBytes(3),
	}
	// #nosec G104
	tx.TransactionAttributes.MarshalFrom(bt)
	genTx, err := moneytx.NewMoneyTx(moneySystemId, tx)
	require.NoError(t, err)
	return genTx
}

func RandomBillTransfer() *txsystem.Transaction {
	tx := randomTx()

	bt := &moneytx.TransferOrder{
		NewBearer: randomBytes(3),
		// #nosec G404
		TargetValue: rand.Uint64(),
		Backlink:    randomBytes(3),
	}
	// #nosec G104
	tx.TransactionAttributes.MarshalFrom(bt)
	return tx
}

func RandomBillSplit() *txsystem.Transaction {
	tx := randomTx()

	bt := &moneytx.SplitOrder{
		// #nosec G404
		Amount:       rand.Uint64(),
		TargetBearer: randomBytes(3),
		// #nosec G404
		RemainingValue: rand.Uint64(),
		Backlink:       randomBytes(3),
	}
	// #nosec G104
	tx.TransactionAttributes.MarshalFrom(bt)
	return tx
}

func randomTx() *txsystem.Transaction {
	return &txsystem.Transaction{
		SystemId:              moneySystemId,
		TransactionAttributes: new(anypb.Any),
		UnitId:                randomBytes(3),
		Timeout:               0,
		OwnerProof:            randomBytes(3),
	}
}

func randomBytes(len int) []byte {
	bytes := make([]byte, len)
	// #nosec G404
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}
	return bytes
}

func CreateBillTransferTx(pubKeyHash []byte) *anypb.Any {
	tx, _ := anypb.New(&moneytx.TransferOrder{
		TargetValue: 100,
		NewBearer:   script.PredicatePayToPublicKeyHashDefault(pubKeyHash),
		Backlink:    hash.Sum256([]byte{}),
	})
	return tx
}

func CreateDustTransferTx(pubKeyHash []byte) *anypb.Any {
	tx, _ := anypb.New(&moneytx.TransferDCOrder{
		TargetValue:  100,
		TargetBearer: script.PredicatePayToPublicKeyHashDefault(pubKeyHash),
		Backlink:     hash.Sum256([]byte{}),
	})
	return tx
}

func CreateBillSplitTx(pubKeyHash []byte, amount uint64, remainingValue uint64) *anypb.Any {
	tx, _ := anypb.New(&moneytx.SplitOrder{
		Amount:         amount,
		TargetBearer:   script.PredicatePayToPublicKeyHashDefault(pubKeyHash),
		RemainingValue: remainingValue,
		Backlink:       hash.Sum256([]byte{}),
	})
	return tx
}

func CreateRandomDcTx() *txsystem.Transaction {
	return &txsystem.Transaction{
		SystemId:              moneySystemId,
		UnitId:                hash.Sum256([]byte{0x00}),
		TransactionAttributes: CreateRandomDustTransferTx(),
		Timeout:               1000,
		OwnerProof:            script.PredicateArgumentEmpty(),
	}
}

func CreateRandomDustTransferTx() *anypb.Any {
	tx, _ := anypb.New(&moneytx.TransferDCOrder{
		TargetBearer: script.PredicateAlwaysTrue(),
		Backlink:     hash.Sum256([]byte{}),
		Nonce:        hash.Sum256([]byte{}),
		TargetValue:  100,
	})
	return tx
}

func CreateRandomSwapTransferTx(pubKeyHash []byte) *anypb.Any {
	tx, _ := anypb.New(&moneytx.SwapOrder{
		OwnerCondition:  script.PredicatePayToPublicKeyHashDefault(pubKeyHash),
		BillIdentifiers: [][]byte{},
		DcTransfers:     []*txsystem.Transaction{},
		Proofs:          []*block.BlockProof{},
		TargetValue:     100,
	})
	return tx
}
