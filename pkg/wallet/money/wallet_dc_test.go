package money

import (
	"context"
	"crypto"
	"testing"

	"github.com/alphabill-org/alphabill/internal/script"

	"github.com/alphabill-org/alphabill/internal/hash"
	moneytesttx "github.com/alphabill-org/alphabill/internal/testutils/transaction/money"
	"github.com/alphabill-org/alphabill/internal/txsystem"
	billtx "github.com/alphabill-org/alphabill/internal/txsystem/money"
	"github.com/alphabill-org/alphabill/internal/util"
	"github.com/alphabill-org/alphabill/pkg/wallet/account"
	"github.com/alphabill-org/alphabill/pkg/wallet/log"
	"github.com/holiman/uint256"
	"github.com/stretchr/testify/require"
)

func TestDustCollectionWontRunForSingleBill(t *testing.T) {
	// create wallet with a single bill
	bills := []*Bill{addBill(1)}
	billsList := createBillListJsonResponse(bills)

	w, mockClient := CreateTestWallet(t, &backendMockReturnConf{customBillList: billsList})

	// when dc runs
	err := w.collectDust(context.Background(), false, 0)
	require.NoError(t, err)

	// then no txs are broadcast
	require.Len(t, mockClient.GetRecordedTransactions(), 0)
}

func TestDustCollectionMaxBillCount(t *testing.T) {
	// create wallet with max allowed bills for dc + 1
	bills := make([]*Bill, maxBillsForDustCollection+1)
	for i := 0; i < maxBillsForDustCollection+1; i++ {
		bills[i] = addBill(uint64(i))
	}
	billsList := createBillListJsonResponse(bills)
	proofList := createBlockProofJsonResponse(t, bills, nil, 0, dcTimeoutBlockCount, nil)

	w, mockClient := CreateTestWallet(t, &backendMockReturnConf{customBillList: billsList, proofList: proofList})

	// when dc runs
	err := w.collectDust(context.Background(), false, 0)
	require.NoError(t, err)

	// then dc tx count should be equal to max allowed bills for dc
	require.Len(t, mockClient.GetRecordedTransactions(), maxBillsForDustCollection)
}

func TestBasicDustCollection(t *testing.T) {
	// create wallet with 2 normal bills
	bills := []*Bill{addBill(1), addBill(2)}
	billsList := createBillListJsonResponse(bills)
	proofList := createBlockProofJsonResponse(t, bills, nil, 0, dcTimeoutBlockCount, nil)
	expectedDcNonce := calculateDcNonce(bills)

	w, mockClient := CreateTestWallet(t, &backendMockReturnConf{balance: 3, customBillList: billsList, proofList: proofList})
	k, _ := w.am.GetAccountKey(0)

	// when dc runs
	err := w.collectDust(context.Background(), false, 0)
	require.NoError(t, err)

	// then two dc txs are broadcast
	require.Len(t, mockClient.GetRecordedTransactions(), 2)
	for i, tx := range mockClient.GetRecordedTransactions() {
		dcTx := parseDcTx(t, tx)
		require.NotNil(t, dcTx)
		require.EqualValues(t, expectedDcNonce, dcTx.Nonce)
		require.EqualValues(t, bills[i].Value, dcTx.TargetValue)
		require.EqualValues(t, bills[i].TxHash, dcTx.Backlink)
		require.EqualValues(t, script.PredicatePayToPublicKeyHashDefault(k.PubKeyHash.Sha256), dcTx.TargetBearer)
	}

	// and expected swap is added to dc wait group
	require.Len(t, w.dcWg.swaps, 1)
	swap := w.dcWg.swaps[*util.BytesToUint256(expectedDcNonce)]
	require.EqualValues(t, expectedDcNonce, swap.dcNonce)
	require.EqualValues(t, 3, swap.dcSum)
	require.EqualValues(t, dcTimeoutBlockCount, swap.timeout)
}

func TestDustCollectionWithSwap(t *testing.T) {
	// create wallet with 2 normal bills
	tempNonce := uint256.NewInt(1)
	am, err := account.NewManager(t.TempDir(), "", true)
	require.NoError(t, err)
	_ = am.CreateKeys("")
	k, _ := am.GetAccountKey(0)
	bills := []*Bill{addBill(1), addBill(2)}
	expectedDcNonce := calculateDcNonce(bills)
	billsList := createBillListJsonResponse(bills)
	// proofs are polled twice, one for the regular bills and one for dc bills
	proofList := createBlockProofJsonResponse(t, bills, nil, 0, dcTimeoutBlockCount, k)
	proofList = append(proofList, createBlockProofJsonResponse(t, []*Bill{addDcBill(t, k, tempNonce, expectedDcNonce, 1, dcTimeoutBlockCount), addDcBill(t, k, tempNonce, expectedDcNonce, 2, dcTimeoutBlockCount)}, expectedDcNonce, 0, dcTimeoutBlockCount, k)...)

	w, mockClient := CreateTestWalletWithManager(t, &backendMockReturnConf{balance: 3, customBillList: billsList, proofList: proofList}, am)

	// when dc runs
	err = w.collectDust(context.Background(), false, 0)
	require.NoError(t, err)

	// then two dc txs + one swap tx are broadcast
	require.Len(t, mockClient.GetRecordedTransactions(), 3)
	for _, tx := range mockClient.GetRecordedTransactions()[0:2] {
		require.NotNil(t, parseDcTx(t, tx))
	}
	txSwap := parseSwapTx(t, mockClient.GetRecordedTransactions()[2])
	require.EqualValues(t, 3, txSwap.TargetValue)
	require.EqualValues(t, [][]byte{util.Uint256ToBytes(tempNonce), util.Uint256ToBytes(tempNonce)}, txSwap.BillIdentifiers)
	require.EqualValues(t, script.PredicatePayToPublicKeyHashDefault(k.PubKeyHash.Sha256), txSwap.OwnerCondition)
	require.Len(t, txSwap.DcTransfers, 2)
	require.Len(t, txSwap.Proofs, 2)

	// and expected swap is updated with swap timeout
	require.Len(t, w.dcWg.swaps, 1)
	swap := w.dcWg.swaps[*util.BytesToUint256(expectedDcNonce)]
	require.EqualValues(t, expectedDcNonce, swap.dcNonce)
	require.EqualValues(t, 3, swap.dcSum)
	require.EqualValues(t, swapTimeoutBlockCount, swap.timeout)
}

func TestSwapWithExistingDCBillsBeforeDCTimeout(t *testing.T) {
	// create wallet with 2 dc bills
	roundNr := uint64(5)
	tempNonce := uint256.NewInt(1)
	am, err := account.NewManager(t.TempDir(), "", true)
	require.NoError(t, err)
	_ = am.CreateKeys("")
	k, _ := am.GetAccountKey(0)
	bills := []*Bill{addDcBill(t, k, tempNonce, util.Uint256ToBytes(tempNonce), 1, dcTimeoutBlockCount), addDcBill(t, k, tempNonce, util.Uint256ToBytes(tempNonce), 2, dcTimeoutBlockCount)}
	billsList := createBillListJsonResponse(bills)
	proofList := createBlockProofJsonResponse(t, bills, util.Uint256ToBytes(tempNonce), 0, dcTimeoutBlockCount, k)
	w, mockClient := CreateTestWalletWithManager(t, &backendMockReturnConf{balance: 3, customBillList: billsList, proofList: proofList}, am)
	// set specific round number
	mockClient.SetMaxRoundNumber(roundNr)

	// when dc runs
	err = w.collectDust(context.Background(), false, 0)
	require.NoError(t, err)

	// then a swap tx is broadcast
	require.Len(t, mockClient.GetRecordedTransactions(), 1)
	txSwap := parseSwapTx(t, mockClient.GetRecordedTransactions()[0])
	require.EqualValues(t, 3, txSwap.TargetValue)
	require.EqualValues(t, [][]byte{util.Uint256ToBytes(tempNonce), util.Uint256ToBytes(tempNonce)}, txSwap.BillIdentifiers)
	require.EqualValues(t, script.PredicatePayToPublicKeyHashDefault(k.PubKeyHash.Sha256), txSwap.OwnerCondition)
	require.Len(t, txSwap.DcTransfers, 2)
	require.Len(t, txSwap.Proofs, 2)

	// and expected swap is updated with swap timeout + round number
	require.Len(t, w.dcWg.swaps, 1)
	swap := w.dcWg.swaps[*tempNonce]
	require.EqualValues(t, util.Uint256ToBytes(tempNonce), swap.dcNonce)
	require.EqualValues(t, 3, swap.dcSum)
	require.EqualValues(t, swapTimeoutBlockCount+roundNr, swap.timeout)
}

func TestSwapWithExistingExpiredDCBills(t *testing.T) {
	// create wallet with 2 timed out dc bills
	tempNonce := uint256.NewInt(1)
	am, err := account.NewManager(t.TempDir(), "", true)
	require.NoError(t, err)
	_ = am.CreateKeys("")
	k, _ := am.GetAccountKey(0)
	bills := []*Bill{addDcBill(t, k, tempNonce, util.Uint256ToBytes(tempNonce), 1, 0), addDcBill(t, k, tempNonce, util.Uint256ToBytes(tempNonce), 2, 0)}
	billsList := createBillListJsonResponse(bills)
	proofList := createBlockProofJsonResponse(t, bills, util.Uint256ToBytes(tempNonce), 0, 0, k)
	w, mockClient := CreateTestWalletWithManager(t, &backendMockReturnConf{balance: 3, customBillList: billsList, proofList: proofList}, am)

	// when dc runs
	err = w.collectDust(context.Background(), false, 0)
	require.NoError(t, err)

	// then a swap tx is broadcast
	require.Len(t, mockClient.GetRecordedTransactions(), 1)
	txSwap := parseSwapTx(t, mockClient.GetRecordedTransactions()[0])
	require.EqualValues(t, 3, txSwap.TargetValue)
	require.EqualValues(t, [][]byte{util.Uint256ToBytes(tempNonce), util.Uint256ToBytes(tempNonce)}, txSwap.BillIdentifiers)
	require.EqualValues(t, script.PredicatePayToPublicKeyHashDefault(k.PubKeyHash.Sha256), txSwap.OwnerCondition)
	require.Len(t, txSwap.DcTransfers, 2)
	require.Len(t, txSwap.Proofs, 2)

	// and expected swap is updated with swap timeout
	require.Len(t, w.dcWg.swaps, 1)
	swap := w.dcWg.swaps[*tempNonce]
	require.EqualValues(t, util.Uint256ToBytes(tempNonce), swap.dcNonce)
	require.EqualValues(t, 3, swap.dcSum)
	require.EqualValues(t, swapTimeoutBlockCount, swap.timeout)
}

func TestDcNonceHashIsCalculatedInCorrectBillOrder(t *testing.T) {
	bills := []*Bill{
		{Id: uint256.NewInt(2)},
		{Id: uint256.NewInt(1)},
		{Id: uint256.NewInt(0)},
	}
	hasher := crypto.SHA256.New()
	for i := len(bills) - 1; i >= 0; i-- {
		hasher.Write(bills[i].GetID())
	}
	expectedNonce := hasher.Sum(nil)

	nonce := calculateDcNonce(bills)
	require.EqualValues(t, expectedNonce, nonce)
}

func TestSwapTxValuesAreCalculatedInCorrectBillOrder(t *testing.T) {
	w, _ := CreateTestWallet(t, nil)
	k, _ := w.am.GetAccountKey(0)

	dcBills := []*Bill{
		{Id: uint256.NewInt(2), BlockProof: &BlockProof{Tx: moneytesttx.CreateRandomDcTx()}},
		{Id: uint256.NewInt(1), BlockProof: &BlockProof{Tx: moneytesttx.CreateRandomDcTx()}},
		{Id: uint256.NewInt(0), BlockProof: &BlockProof{Tx: moneytesttx.CreateRandomDcTx()}},
	}
	dcNonce := calculateDcNonce(dcBills)
	var dcBillIds [][]byte
	for _, dcBill := range dcBills {
		dcBillIds = append(dcBillIds, dcBill.GetID())
	}

	tx, err := createSwapTx(k, w.SystemID(), dcBills, dcNonce, dcBillIds, 10)
	require.NoError(t, err)
	swapTx := parseSwapTx(t, tx)

	// verify bill ids in swap tx are in correct order (equal hash values)
	hasher := crypto.SHA256.New()
	for _, billId := range swapTx.BillIdentifiers {
		hasher.Write(billId)
	}
	actualDcNonce := hasher.Sum(nil)
	require.EqualValues(t, dcNonce, actualDcNonce)
}

func TestSwapContainsUnconfirmedDustBillIds(t *testing.T) {
	// create wallet with three bills
	_ = log.InitStdoutLogger(log.INFO)
	b1 := addBill(1)
	b2 := addBill(2)
	b3 := addBill(3)
	nonce := calculateDcNonce([]*Bill{b1, b2, b3})
	am, err := account.NewManager(t.TempDir(), "", true)
	require.NoError(t, err)
	_ = am.CreateKeys("")
	k, _ := am.GetAccountKey(0)

	billsList := createBillListJsonResponse([]*Bill{b1, b2, b3})
	// proofs are polled twice, one for the regular bills and one for dc bills
	proofList := createBlockProofJsonResponse(t, []*Bill{b1, b2, b3}, nil, 0, dcTimeoutBlockCount, k)
	proofList = append(proofList, createBlockProofJsonResponse(t, []*Bill{addDcBill(t, k, b1.Id, nonce, 1, dcTimeoutBlockCount), addDcBill(t, k, b2.Id, nonce, 2, dcTimeoutBlockCount), addDcBill(t, k, b3.Id, nonce, 3, dcTimeoutBlockCount)}, nonce, 0, dcTimeoutBlockCount, k)...)
	w, mockClient := CreateTestWalletWithManager(t, &backendMockReturnConf{balance: 3, customBillList: billsList, proofList: proofList}, am)

	// when dc runs
	err = w.collectDust(context.Background(), false, 0)
	require.NoError(t, err)

	verifyBlockHeight(t, w, 0)

	// and three dc txs are broadcast
	dcTxs := mockClient.GetRecordedTransactions()
	require.Len(t, dcTxs, 4)
	for _, tx := range dcTxs[0:3] {
		require.NotNil(t, parseDcTx(t, tx))
	}

	// and swap should contain all bill ids
	tx := mockClient.GetRecordedTransactions()[3]
	swapOrder := parseSwapTx(t, tx)
	require.EqualValues(t, nonce, tx.UnitId)
	require.Len(t, swapOrder.BillIdentifiers, 3)
	require.Equal(t, b1.Id, uint256.NewInt(0).SetBytes(swapOrder.BillIdentifiers[0]))
	require.Equal(t, b2.Id, uint256.NewInt(0).SetBytes(swapOrder.BillIdentifiers[1]))
	require.Equal(t, b3.Id, uint256.NewInt(0).SetBytes(swapOrder.BillIdentifiers[2]))
	require.Len(t, swapOrder.DcTransfers, 3)
	require.Equal(t, dcTxs[0], swapOrder.DcTransfers[0])
	require.Equal(t, dcTxs[1], swapOrder.DcTransfers[1])
	require.Equal(t, dcTxs[2], swapOrder.DcTransfers[2])
}

func addBill(value uint64) *Bill {
	b1 := Bill{
		Id:     uint256.NewInt(value),
		Value:  value,
		TxHash: hash.Sum256([]byte{byte(value)}),
	}
	return &b1
}

func addDcBill(t *testing.T, k *account.AccountKey, id *uint256.Int, nonce []byte, value uint64, timeout uint64) *Bill {
	b := Bill{
		Id:     id,
		Value:  value,
		TxHash: hash.Sum256([]byte{byte(value)}),
	}

	tx, err := createDustTx(k, []byte{0, 0, 0, 0}, &b, nonce, timeout)
	require.NoError(t, err)
	b.BlockProof = &BlockProof{Tx: tx}

	b.IsDcBill = true
	b.DcNonce = nonce
	b.DcTimeout = timeout
	b.DcExpirationTimeout = dustBillDeletionTimeout

	require.NoError(t, err)
	return &b
}

func verifyBlockHeight(t *testing.T, w *Wallet, blockHeight uint64) {
	actualBlockHeight, _, err := w.AlphabillClient.GetMaxBlockNumber(context.Background())
	require.NoError(t, err)
	require.Equal(t, blockHeight, actualBlockHeight)
}

func parseBillTransferTx(t *testing.T, tx *txsystem.Transaction) *billtx.TransferOrder {
	btTx := &billtx.TransferOrder{}
	err := tx.TransactionAttributes.UnmarshalTo(btTx)
	require.NoError(t, err)
	return btTx
}

func parseDcTx(t *testing.T, tx *txsystem.Transaction) *billtx.TransferDCOrder {
	dcTx := &billtx.TransferDCOrder{}
	err := tx.TransactionAttributes.UnmarshalTo(dcTx)
	require.NoError(t, err)
	return dcTx
}

func parseSwapTx(t *testing.T, tx *txsystem.Transaction) *billtx.SwapOrder {
	txSwap := &billtx.SwapOrder{}
	err := tx.TransactionAttributes.UnmarshalTo(txSwap)
	require.NoError(t, err)
	return txSwap
}
