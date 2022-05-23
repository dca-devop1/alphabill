package money

import (
	"crypto"
	"testing"

	test "gitdc.ee.guardtime.com/alphabill/alphabill/internal/testutils"
	"gitdc.ee.guardtime.com/alphabill/alphabill/internal/txsystem"
	"github.com/holiman/uint256"
	"github.com/stretchr/testify/require"
)

func TestTransfer(t *testing.T) {
	tests := []struct {
		name string
		bd   *BillData
		tx   *transferWrapper
		res  error
	}{
		{
			name: "Ok",
			bd:   newBillData(100, []byte{6}),
			tx:   newTransfer(t, 100, []byte{6}),
			res:  nil,
		},
		{
			name: "InvalidBalance",
			bd:   newBillData(100, []byte{6}),
			tx:   newTransfer(t, 101, []byte{6}),
			res:  ErrInvalidBillValue,
		},
		{
			name: "InvalidBacklink",
			bd:   newBillData(100, []byte{6}),
			tx:   newTransfer(t, 100, []byte{5}),
			res:  txsystem.ErrInvalidBacklink,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateTransfer(tt.bd, tt.tx)
			if tt.res == nil {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, tt.res)
			}
		})
	}
}

func TestTransferDC(t *testing.T) {
	tests := []struct {
		name string
		bd   *BillData
		tx   *transferDCWrapper
		res  error
	}{
		{
			name: "Ok",
			bd:   newBillData(100, []byte{6}),
			tx:   newTransferDC(t, 100, []byte{6}, []byte{1}, test.RandomBytes(32)),
			res:  nil,
		},
		{
			name: "InvalidBalance",
			bd:   newBillData(100, []byte{6}),
			tx:   newTransferDC(t, 101, []byte{6}, []byte{1}, test.RandomBytes(32)),
			res:  ErrInvalidBillValue,
		},
		{
			name: "InvalidBacklink",
			bd:   newBillData(100, []byte{6}),
			tx:   newTransferDC(t, 100, []byte{5}, []byte{1}, test.RandomBytes(32)),
			res:  txsystem.ErrInvalidBacklink,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateTransferDC(tt.bd, tt.tx)
			if tt.res == nil {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, tt.res)
			}
		})
	}
}

func TestSplit(t *testing.T) {
	tests := []struct {
		name string
		bd   *BillData
		tx   *billSplitWrapper
		res  error
	}{
		{
			name: "Ok",
			bd:   newBillData(100, []byte{6}),
			tx:   newSplit(t, 50, 50, []byte{6}),
			res:  nil,
		},
		{
			name: "AmountExceedsBillValue",
			bd:   newBillData(100, []byte{6}),
			tx:   newSplit(t, 101, 100, []byte{6}),
			res:  ErrInvalidBillValue,
		},
		{
			name: "AmountEqualsBillValue",
			bd:   newBillData(100, []byte{6}),
			tx:   newSplit(t, 100, 0, []byte{6}),
			res:  ErrInvalidBillValue,
		},
		{
			name: "InvalidRemainingValue",
			bd:   newBillData(100, []byte{6}),
			tx:   newSplit(t, 50, 51, []byte{6}),
			res:  ErrInvalidBillValue,
		},
		{
			name: "InvalidBacklink",
			bd:   newBillData(100, []byte{6}),
			tx:   newSplit(t, 50, 50, []byte{5}),
			res:  txsystem.ErrInvalidBacklink,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateSplit(tt.bd, tt.tx)
			if tt.res == nil {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, tt.res)
			}
		})
	}
}

func TestSwap(t *testing.T) {
	tests := []struct {
		name string
		tx   *swapWrapper
		res  error
	}{
		{
			name: "Ok",
			tx:   newValidSwap(t),
			res:  nil,
		},
		{
			name: "InvalidTargetValue",
			tx:   newInvalidTargetValueSwap(t),
			res:  ErrSwapInvalidTargetValue,
		},
		{
			name: "InvalidBillIdentifiers",
			tx:   newInvalidBillIdentifierSwap(t),
			res:  ErrSwapInvalidBillIdentifiers,
		},
		{
			name: "InvalidBillId",
			tx:   newInvalidBillIdSwap(t),
			res:  ErrSwapInvalidBillId,
		},
		{
			name: "InvalidNonce",
			tx:   newInvalidNonceSwap(t),
			res:  ErrSwapInvalidNonce,
		},
		{
			name: "InvalidTargetBearer",
			tx:   newInvalidTargetBearerSwap(t),
			res:  ErrSwapInvalidTargetBearer,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateSwap(tt.tx, crypto.SHA256)
			if tt.res == nil {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, tt.res)
			}
		})
	}
}

func newTransfer(t *testing.T, v uint64, backlink []byte) *transferWrapper {
	tx, err := NewMoneyTx(newPBTransactionOrder([]byte{1}, []byte{3}, 2, &TransferOrder{
		NewBearer:   []byte{4},
		TargetValue: v,
		Backlink:    backlink,
	}))
	require.NoError(t, err)
	require.IsType(t, tx, &transferWrapper{})
	return tx.(*transferWrapper)
}

func newTransferDC(t *testing.T, v uint64, backlink []byte, unitID []byte, nonce []byte) *transferDCWrapper {
	order := newPBTransactionOrder(unitID, []byte{3}, 2, &TransferDCOrder{
		Nonce:        nonce,
		TargetBearer: []byte{4},
		TargetValue:  v,
		Backlink:     backlink,
	})
	order.SystemId = []byte{0}
	tx, err := NewMoneyTx(order)
	require.NoError(t, err)
	require.IsType(t, tx, &transferDCWrapper{})
	return tx.(*transferDCWrapper)
}

func newSplit(t *testing.T, amount uint64, remainingValue uint64, backlink []byte) *billSplitWrapper {
	order := newPBTransactionOrder([]byte{1}, []byte{3}, 2, &SplitOrder{
		Amount:         amount,
		TargetBearer:   []byte{5},
		RemainingValue: remainingValue,
		Backlink:       backlink,
	})
	order.SystemId = []byte{0}
	tx, err := NewMoneyTx(order)
	require.NoError(t, err)
	require.IsType(t, tx, &billSplitWrapper{})
	return tx.(*billSplitWrapper)
}

func newInvalidTargetValueSwap(t *testing.T) *swapWrapper {
	id := uint256.NewInt(1)
	transferDCID, swapId := calculateSwapID(id)
	dcTransfer := newTransferDC(t, 100, []byte{6}, transferDCID, swapId)
	order := newPBTransactionOrder(swapId, []byte{3}, 2, &SwapOrder{
		OwnerCondition:  dcTransfer.TargetBearer(),
		BillIdentifiers: [][]byte{transferDCID},
		DcTransfers:     []*txsystem.Transaction{dcTransfer.transaction},
		Proofs:          [][]byte{{9}, {10}},
		TargetValue:     dcTransfer.TargetValue() - 1,
	})
	order.SystemId = []byte{0}
	tx, err := NewMoneyTx(order)
	require.NoError(t, err)
	require.IsType(t, tx, &swapWrapper{})
	return tx.(*swapWrapper)
}

func newInvalidBillIdentifierSwap(t *testing.T) *swapWrapper {
	id := uint256.NewInt(1)
	transferId, swapId := calculateSwapID(id)
	dcTransfer := newTransferDC(t, 100, []byte{6}, test.RandomBytes(3), swapId)
	order := newPBTransactionOrder(swapId, []byte{3}, 2, newSwapOrder(dcTransfer, transferId))
	tx, err := NewMoneyTx(order)
	require.NoError(t, err)
	require.IsType(t, tx, &swapWrapper{})
	return tx.(*swapWrapper)
}

func newInvalidBillIdSwap(t *testing.T) *swapWrapper {
	id := uint256.NewInt(1)
	transferId, swapId := calculateSwapID(id)
	dcTransfer := newTransferDC(t, 100, []byte{6}, transferId, swapId)
	order := newPBTransactionOrder([]byte{0}, []byte{3}, 2, newSwapOrder(dcTransfer, transferId))
	order.SystemId = []byte{0}
	tx, err := NewMoneyTx(order)
	require.NoError(t, err)
	require.IsType(t, tx, &swapWrapper{})
	return tx.(*swapWrapper)
}

func newInvalidNonceSwap(t *testing.T) *swapWrapper {
	id := uint256.NewInt(1)
	transferId, swapId := calculateSwapID(id)
	dcTransfer := newTransferDC(t, 100, []byte{6}, transferId, []byte{0})
	order := newPBTransactionOrder(swapId, []byte{3}, 2, newSwapOrder(dcTransfer, transferId))
	tx, err := NewMoneyTx(order)
	require.NoError(t, err)
	require.IsType(t, tx, &swapWrapper{})
	return tx.(*swapWrapper)
}

func newInvalidTargetBearerSwap(t *testing.T) *swapWrapper {
	id := uint256.NewInt(1)
	transferId, swapId := calculateSwapID(id)
	dcTransfer := newTransferDC(t, 100, []byte{6}, transferId, swapId)
	order := newPBTransactionOrder(swapId, []byte{3}, 2, &SwapOrder{
		OwnerCondition:  test.RandomBytes(32),
		BillIdentifiers: [][]byte{transferId},
		DcTransfers:     []*txsystem.Transaction{dcTransfer.transaction},
		Proofs:          [][]byte{{9}, {10}},
		TargetValue:     dcTransfer.TargetValue(),
	})
	tx, err := NewMoneyTx(order)
	require.NoError(t, err)
	require.IsType(t, tx, &swapWrapper{})
	return tx.(*swapWrapper)
}

func newValidSwap(t *testing.T) *swapWrapper {
	id := uint256.NewInt(1)
	transferId, swapId := calculateSwapID(id)
	dcTransfer := newTransferDC(t, 100, []byte{6}, transferId, swapId)
	order := newPBTransactionOrder(swapId, []byte{3}, 2, newSwapOrder(dcTransfer, transferId))
	tx, err := NewMoneyTx(order)
	require.NoError(t, err)
	require.IsType(t, tx, &swapWrapper{})
	return tx.(*swapWrapper)
}

func newSwapOrder(dcTransfer *transferDCWrapper, transferDCID []byte) *SwapOrder {
	return &SwapOrder{
		OwnerCondition:  dcTransfer.TargetBearer(),
		BillIdentifiers: [][]byte{transferDCID},
		DcTransfers:     []*txsystem.Transaction{dcTransfer.transaction},
		Proofs:          [][]byte{{9}, {10}},
		TargetValue:     dcTransfer.TargetValue(),
	}
}

func calculateSwapID(id *uint256.Int) ([]byte, []byte) {
	hasher := crypto.SHA256.New()
	bytes32 := id.Bytes32()
	hasher.Write(bytes32[:])
	swapId := hasher.Sum(nil)
	return bytes32[:], swapId
}

func newBillData(v uint64, backlink []byte) *BillData {
	return &BillData{V: v, Backlink: backlink}
}
