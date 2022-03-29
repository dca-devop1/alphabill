package txsystem

import (
	"crypto"
	"github.com/holiman/uint256"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTransfer(t *testing.T) {
	tests := []struct {
		name string
		bd   *BillData
		tx   *transfer
		res  error
	}{
		{
			name: "Ok",
			bd:   newBillData(100, []byte{6}),
			tx:   newTransfer(100, []byte{6}),
			res:  nil,
		},
		{
			name: "InvalidBalance",
			bd:   newBillData(100, []byte{6}),
			tx:   newTransfer(101, []byte{6}),
			res:  ErrInvalidBillValue,
		},
		{
			name: "InvalidBacklink",
			bd:   newBillData(100, []byte{6}),
			tx:   newTransfer(100, []byte{5}),
			res:  ErrInvalidBacklink,
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
		tx   *transferDC
		res  error
	}{
		{
			name: "Ok",
			bd:   newBillData(100, []byte{6}),
			tx:   newTransferDC(100, []byte{6}),
			res:  nil,
		},
		{
			name: "InvalidBalance",
			bd:   newBillData(100, []byte{6}),
			tx:   newTransferDC(101, []byte{6}),
			res:  ErrInvalidBillValue,
		},
		{
			name: "InvalidBacklink",
			bd:   newBillData(100, []byte{6}),
			tx:   newTransferDC(100, []byte{5}),
			res:  ErrInvalidBacklink,
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
		tx   *split
		res  error
	}{
		{
			name: "Ok",
			bd:   newBillData(100, []byte{6}),
			tx:   newSplit(50, 50, []byte{6}),
			res:  nil,
		},
		{
			name: "AmountExceedsBillValue",
			bd:   newBillData(100, []byte{6}),
			tx:   newSplit(101, 100, []byte{6}),
			res:  ErrInvalidBillValue,
		},
		{
			name: "AmountEqualsBillValue",
			bd:   newBillData(100, []byte{6}),
			tx:   newSplit(100, 0, []byte{6}),
			res:  ErrInvalidBillValue,
		},
		{
			name: "InvalidRemainingValue",
			bd:   newBillData(100, []byte{6}),
			tx:   newSplit(50, 51, []byte{6}),
			res:  ErrInvalidBillValue,
		},
		{
			name: "InvalidBacklink",
			bd:   newBillData(100, []byte{6}),
			tx:   newSplit(50, 50, []byte{5}),
			res:  ErrInvalidBacklink,
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
		tx   *swap
		res  error
	}{
		{
			name: "Ok",
			tx:   newValidSwap(),
			res:  nil,
		},
		{
			name: "InvalidTargetValue",
			tx:   newInvalidTargetValueSwap(),
			res:  ErrSwapInvalidTargetValue,
		},
		{
			name: "InvalidBillIdentifiers",
			tx:   newInvalidBillIdentifierSwap(),
			res:  ErrSwapInvalidBillIdentifiers,
		},
		{
			name: "InvalidBillId",
			tx:   newInvalidBillIdSwap(),
			res:  ErrSwapInvalidBillId,
		},
		{
			name: "InvalidNonce",
			tx:   newInvalidNonceSwap(),
			res:  ErrSwapInvalidNonce,
		},
		{
			name: "InvalidTargetBearer",
			tx:   newInvalidTargetBearerSwap(),
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

func TestGenericTxValidation(t *testing.T) {
	tests := []struct {
		name    string
		gtx     GenericTransaction
		blockNo uint64
		res     error
	}{
		{
			name:    "Ok",
			gtx:     newTransferWithTimeout(11),
			blockNo: 10,
			res:     nil,
		},
		{
			name:    "TransactionExpired",
			gtx:     newTransferWithTimeout(10),
			blockNo: 10,
			res:     ErrTransactionExpired,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateGenericTransaction(tt.gtx, tt.blockNo)
			if tt.res == nil {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, tt.res)
			}
		})
	}
}

func newTransfer(v uint64, backlink []byte) *transfer {
	return &transfer{
		genericTx: genericTx{
			unitId:     uint256.NewInt(1),
			timeout:    2,
			ownerProof: []byte{3},
		},
		newBearer:   []byte{4},
		targetValue: v,
		backlink:    backlink,
	}
}

func newTransferWithTimeout(timeout uint64) *transfer {
	return &transfer{
		genericTx: genericTx{
			unitId:     uint256.NewInt(1),
			timeout:    timeout,
			ownerProof: []byte{3},
		},
		newBearer:   []byte{4},
		targetValue: 5,
		backlink:    []byte{6},
	}
}

func newTransferDC(v uint64, backlink []byte) *transferDC {
	return &transferDC{
		genericTx: genericTx{
			unitId:     uint256.NewInt(1),
			timeout:    2,
			ownerProof: []byte{3},
		},
		targetBearer: []byte{4},
		targetValue:  v,
		backlink:     backlink,
		nonce:        []byte{7},
	}
}

func newSplit(amount uint64, remainingValue uint64, backlink []byte) *split {
	return &split{
		genericTx: genericTx{
			unitId:     uint256.NewInt(1),
			timeout:    2,
			ownerProof: []byte{3},
		},
		amount:         amount,
		targetBearer:   []byte{5},
		remainingValue: remainingValue,
		backlink:       backlink,
	}
}

func newInvalidTargetValueSwap() *swap {
	s := newValidSwap()
	s.targetValue = s.targetValue - 1
	return s
}

func newInvalidBillIdentifierSwap() *swap {
	s := newValidSwap()

	// add extra swap dc transfer that is not in swap bill identifier list
	dc := newRandomTransferDC()
	dc.nonce = s.dcTransfers[0].Nonce() // correct nonce
	s.dcTransfers = append(s.dcTransfers, dc)

	// update swap target sum with new dc transfer
	s.targetValue += dc.targetValue

	return s
}

func newInvalidBillIdSwap() *swap {
	s := newValidSwap()
	s.unitId = uint256.NewInt(1)
	return s
}

func newInvalidNonceSwap() *swap {
	s := newValidSwap()
	s.dcTransfers[0].Nonce()[0] = 0
	return s
}

func newInvalidTargetBearerSwap() *swap {
	s := newValidSwap()
	dc := s.dcTransfers[0].(*transferDC)
	dc.targetBearer = uint256.NewInt(2).Bytes()
	return s
}

func newValidSwap() *swap {
	dcTransfer := newRandomTransferDC()

	// swap tx bill id = hash of dc transfers
	hasher := crypto.SHA256.New()
	hasher.Write(dcTransfer.unitId.Bytes())
	billId := hasher.Sum(nil)

	// dc transfer nonce must be equal to swap tx id
	dcTransfer.nonce = billId

	return &swap{
		genericTx: genericTx{
			systemID:   []byte{0},
			unitId:     uint256.NewInt(0).SetBytes(billId),
			timeout:    2,
			ownerProof: []byte{3},
		},
		ownerCondition:  dcTransfer.targetBearer,
		billIdentifiers: []*uint256.Int{dcTransfer.unitId},
		dcTransfers:     []TransferDC{dcTransfer},
		proofs:          [][]byte{{9}, {10}},
		targetValue:     dcTransfer.targetValue,
	}
}

func newBillData(v uint64, backlink []byte) *BillData {
	return &BillData{V: v, Backlink: backlink}
}
