package money

import (
	"context"
	"testing"

	"github.com/holiman/uint256"
	"github.com/stretchr/testify/require"
)

func TestDcJobWithExistingDcBills(t *testing.T) {
	// wallet contains 2 dc bills with the same nonce that have timed out
	w, _ := CreateTestWallet(t, nil)
	k, _ := w.am.GetAccountKey(0)
	bills := []*Bill{addDcBill(t, w, k, uint256.NewInt(1), 1, dcTimeoutBlockCount), addDcBill(t, w, k, uint256.NewInt(1), 2, dcTimeoutBlockCount)}
	nonce := calculateDcNonce(bills)
	billsList := createBillListJsonResponse(bills)
	proofList := createBlockProofJsonResponse(t, bills, nonce, 0, dcTimeoutBlockCount)
	w, mockClient := CreateTestWallet(t, &backendMockReturnConf{balance: 3, customBillList: billsList, proofList: proofList})
	mockClient.SetMaxBlockNumber(100)

	// when dust collector runs
	err := w.collectDust(context.Background(), false, 0)
	require.NoError(t, err)

	// then swap tx is broadcast
	require.Len(t, mockClient.GetRecordedTransactions(), 1)
	tx := mockClient.GetRecordedTransactions()[0]
	txSwap := parseSwapTx(t, tx)

	// and verify each dc tx id = nonce = swap.id
	require.Len(t, txSwap.DcTransfers, 2)
	for i := 0; i < len(txSwap.DcTransfers); i++ {
		dcTx := parseDcTx(t, txSwap.DcTransfers[i])
		require.EqualValues(t, nonce, dcTx.Nonce)
		require.EqualValues(t, nonce, tx.UnitId)
	}
}

func TestDcJobWithExistingDcAndNonDcBills(t *testing.T) {
	// wallet contains timed out dc bill and normal bill
	w, _ := CreateTestWallet(t, nil)
	k, _ := w.am.GetAccountKey(0)
	bill := addBill(1)
	dc := addDcBill(t, w, k, uint256.NewInt(1), 2, dcTimeoutBlockCount)
	nonce := calculateDcNonce([]*Bill{bill, dc})
	billsList := createBillListJsonResponse([]*Bill{bill, dc})
	proofList := createBlockProofJsonResponse(t, []*Bill{bill, dc}, nonce, 0, dcTimeoutBlockCount)

	w, mockClient := CreateTestWallet(t, &backendMockReturnConf{balance: 3, customBillList: billsList, proofList: proofList})
	mockClient.SetMaxBlockNumber(100)

	// when dust collector runs
	err := w.collectDust(context.Background(), false, 0)
	require.NoError(t, err)

	// then swap tx is sent for the timed out dc bill
	require.Len(t, mockClient.GetRecordedTransactions(), 1)
	tx := mockClient.GetRecordedTransactions()[0]
	txSwap := parseSwapTx(t, tx)

	// and verify nonce = swap.id = dc tx id
	require.Len(t, txSwap.DcTransfers, 1)
	for i := 0; i < len(txSwap.DcTransfers); i++ {
		dcTx := parseDcTx(t, txSwap.DcTransfers[i])
		require.EqualValues(t, nonce, dcTx.Nonce)
		require.EqualValues(t, nonce, tx.UnitId)
	}
}

func TestDcJobWithExistingNonDcBills(t *testing.T) {
	// wallet contains 2 non dc bills
	bills := []*Bill{addBill(1), addBill(2)}
	billsList := createBillListJsonResponse(bills)
	proofList := createBlockProofJsonResponse(t, bills, nil, 0, dcTimeoutBlockCount)

	w, mockClient := CreateTestWallet(t, &backendMockReturnConf{balance: 3, customBillList: billsList, proofList: proofList})
	mockClient.SetMaxBlockNumber(100)

	// when dust collector runs
	err := w.collectDust(context.Background(), false, 0)
	require.NoError(t, err)

	// then dust txs are broadcast
	require.Len(t, mockClient.GetRecordedTransactions(), 2)

	// and nonces are equal
	dcTx0 := parseDcTx(t, mockClient.GetRecordedTransactions()[0])
	dcTx1 := parseDcTx(t, mockClient.GetRecordedTransactions()[1])
	require.EqualValues(t, dcTx0.Nonce, dcTx1.Nonce)
}

func TestDcJobSendsSwapsIfDcBillTimeoutHasBeenReached(t *testing.T) {
	// wallet contains 2 dc bills that both have timed out
	w, _ := CreateTestWallet(t, nil)
	k, _ := w.am.GetAccountKey(0)
	bills := []*Bill{addDcBill(t, w, k, uint256.NewInt(1), 1, dcTimeoutBlockCount), addDcBill(t, w, k, uint256.NewInt(1), 2, dcTimeoutBlockCount)}
	nonce := calculateDcNonce(bills)
	billsList := createBillListJsonResponse(bills)
	proofList := createBlockProofJsonResponse(t, bills, nonce, 0, dcTimeoutBlockCount)
	w, mockClient := CreateTestWallet(t, &backendMockReturnConf{balance: 3, customBillList: billsList, proofList: proofList})

	// when dust collector runs
	err := w.collectDust(context.Background(), false, 0)
	require.NoError(t, err)

	// then 2 swap txs must be broadcast
	require.Len(t, mockClient.GetRecordedTransactions(), 1)
	for _, tx := range mockClient.GetRecordedTransactions() {
		require.NotNil(t, parseSwapTx(t, tx))
	}
}
