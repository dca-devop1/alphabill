package money

import (
	"crypto"
	"testing"

	test "gitdc.ee.guardtime.com/alphabill/alphabill/internal/testutils"

	"gitdc.ee.guardtime.com/alphabill/alphabill/internal/util"

	txutil "gitdc.ee.guardtime.com/alphabill/alphabill/internal/txsystem/util"

	"google.golang.org/protobuf/types/known/anypb"

	"gitdc.ee.guardtime.com/alphabill/alphabill/internal/rma"
	"gitdc.ee.guardtime.com/alphabill/alphabill/internal/script"
	"gitdc.ee.guardtime.com/alphabill/alphabill/internal/txsystem"
	"github.com/holiman/uint256"
	"github.com/stretchr/testify/require"
)

const (
	addItemId                = 0
	addItemOwner             = 1
	addItemData              = 2
	updateDataUpdateFunction = 1
)

var initialBill = &InitialBill{ID: uint256.NewInt(77), Value: 110, Owner: script.PredicateAlwaysTrue()}

const initialDustCollectorMoneyAmount uint64 = 100

func TestNewMoneyScheme(t *testing.T) {
	mockRevertibleState, err := rma.New(&rma.Config{
		HashAlgorithm:     crypto.SHA256,
		RecordingDisabled: false,
	})
	require.NoError(t, err)

	initialBill := &InitialBill{ID: uint256.NewInt(2), Value: 100, Owner: nil}
	dcMoneyAmount := uint64(222)

	txSystem, err := NewMoneyTxSystem(crypto.SHA256, initialBill, dcMoneyAmount, SchemeOpts.RevertibleState(mockRevertibleState))
	require.NoError(t, err)
	u, err := txSystem.revertibleState.GetUnit(initialBill.ID)
	require.NoError(t, err)
	require.NotNil(t, u)
	require.Equal(t, rma.Uint64SummaryValue(initialBill.Value), u.Data.Value())
	require.Equal(t, initialBill.Owner, u.Bearer)

	d, err := txSystem.revertibleState.GetUnit(dustCollectorMoneySupplyID)
	require.NoError(t, err)
	require.NotNil(t, d)

	require.Equal(t, rma.Uint64SummaryValue(dcMoneyAmount), d.Data.Value())
	require.Equal(t, rma.Predicate(dustCollectorPredicate), d.Bearer)
}

func TestNewMoneyScheme_InitialBillIsNil(t *testing.T) {
	_, err := NewMoneyTxSystem(crypto.SHA256, nil, 10)
	require.ErrorIs(t, err, ErrInitialBillIsNil)
}

func TestNewMoneyScheme_InvalidInitialBillID(t *testing.T) {
	ib := &InitialBill{ID: uint256.NewInt(0), Value: 100, Owner: nil}
	_, err := NewMoneyTxSystem(crypto.SHA256, ib, 10)
	require.ErrorIs(t, err, ErrInvalidInitialBillID)
}

func TestExecute_TransferOk(t *testing.T) {
	rmaTree, txSystem := createRMATreeAndTxSystem(t)
	unit, data := getBill(t, rmaTree, initialBill.ID)

	transferOk, err := NewMoneyTx(createBillTransfer(initialBill.ID, initialBill.Value, nil))
	require.NoError(t, err)
	roundNumber := uint64(1)
	txSystem.BeginBlock(roundNumber)
	err = txSystem.Execute(transferOk)
	txSystem.GetRootHash()
	require.NoError(t, err)
	unit2, data2 := getBill(t, rmaTree, initialBill.ID)
	require.Equal(t, data.Value(), data2.Value())
	require.NotEqual(t, transferOk.OwnerProof(), unit2.Bearer)
	require.NotEqual(t, unit.StateHash, unit2.StateHash)
	require.EqualValues(t, transferOk.Hash(crypto.SHA256), data2.Backlink)
	require.Equal(t, roundNumber, data2.T)
}

func TestExecute_SplitOk(t *testing.T) {
	rmaTree, txSystem := createRMATreeAndTxSystem(t)
	totalValue := rmaTree.TotalValue()
	initBill, initBillData := getBill(t, rmaTree, initialBill.ID)
	var remaining uint64 = 10
	amount := initialBill.Value - remaining
	splitOk, err := NewMoneyTx(createSplit(initialBill.ID, amount, remaining, script.PredicateAlwaysTrue(), initBillData.Backlink))
	require.NoError(t, err)
	roundNumber := uint64(1)
	txSystem.BeginBlock(roundNumber)
	err = txSystem.Execute(splitOk)
	txSystem.GetRootHash()
	require.NoError(t, err)
	initBillAfterUpdate, initBillDataAfterUpdate := getBill(t, rmaTree, initialBill.ID)

	// bill value was reduced
	require.NotEqual(t, initBillData.V, initBillDataAfterUpdate.V)
	require.Equal(t, remaining, initBillDataAfterUpdate.V)
	// total value was not changed
	require.Equal(t, totalValue, rmaTree.TotalValue())
	// bearer of the initial bill was not changed
	require.Equal(t, initBill.Bearer, initBillAfterUpdate.Bearer)
	require.Equal(t, roundNumber, initBillDataAfterUpdate.T)

	splitWrapper := splitOk.(*billSplitWrapper)
	expectedNewUnitId := txutil.SameShardId(splitOk.UnitID(), unitIdFromTransaction(splitWrapper))
	newBill, newBillData := getBill(t, rmaTree, expectedNewUnitId)
	require.NotNil(t, newBill)
	require.NotNil(t, newBillData)
	require.Equal(t, amount, newBillData.V)
	require.EqualValues(t, splitOk.Hash(crypto.SHA256), newBillData.Backlink)
	require.Equal(t, rma.Predicate(splitWrapper.billSplit.TargetBearer), newBill.Bearer)
	require.Equal(t, roundNumber, newBillData.T)
}

func TestExecute_TransferDCOk(t *testing.T) {
	rmaTree, txSystem := createRMATreeAndTxSystem(t)
	_, initBillData := getBill(t, rmaTree, initialBill.ID)
	var remaining uint64 = 10
	amount := initialBill.Value - remaining
	splitOk, err := NewMoneyTx(createSplit(initialBill.ID, amount, remaining, script.PredicateAlwaysTrue(), initBillData.Backlink))
	require.NoError(t, err)
	roundNumber := uint64(10)
	txSystem.BeginBlock(roundNumber)
	err = txSystem.Execute(splitOk)

	splitWrapper := splitOk.(*billSplitWrapper)
	billID := txutil.SameShardId(splitOk.UnitID(), unitIdFromTransaction(splitWrapper))
	_, splitBillData := getBill(t, rmaTree, billID)
	transferDCOk, err := NewMoneyTx(createDCTransfer(billID, splitBillData.V, splitBillData.Backlink, test.RandomBytes(32)))
	require.NoError(t, err)

	err = txSystem.Execute(transferDCOk)
	require.NoError(t, err)

	transferDCBill, transferDCBillData := getBill(t, rmaTree, billID)
	require.NotEqual(t, dustCollectorPredicate, transferDCBill.Bearer)
	require.Equal(t, splitBillData.Value(), transferDCBillData.Value())
	require.Equal(t, roundNumber, transferDCBillData.T)
	require.EqualValues(t, transferDCOk.Hash(crypto.SHA256), transferDCBillData.Backlink)
}

func TestExecute_SwapOk(t *testing.T) {
	rmaTree, txSystem := createRMATreeAndTxSystem(t)
	_, initBillData := getBill(t, rmaTree, initialBill.ID)
	var remaining uint64 = 99
	amount := initialBill.Value - remaining
	splitOk, err := NewMoneyTx(createSplit(initialBill.ID, amount, remaining, script.PredicateAlwaysTrue(), initBillData.Backlink))
	require.NoError(t, err)
	roundNumber := uint64(10)
	txSystem.BeginBlock(roundNumber)
	err = txSystem.Execute(splitOk)

	splitWrapper := splitOk.(*billSplitWrapper)
	splitBillID := txutil.SameShardId(splitOk.UnitID(), unitIdFromTransaction(splitWrapper))

	dcTransfers, swapTx := createDCTransferAndSwapTxs(t, []*uint256.Int{splitBillID}, rmaTree)

	for _, dcTransfer := range dcTransfers {
		tx, err := NewMoneyTx(dcTransfer)
		require.NoError(t, err)
		err = txSystem.Execute(tx)
		require.NoError(t, err)
	}
	rmaTree.GetRootHash()
	swap, err := NewMoneyTx(swapTx)
	err = txSystem.Execute(swap)
	require.NoError(t, err)
	_, newBillData := getBill(t, rmaTree, swap.UnitID())
	require.Equal(t, amount, newBillData.V)
	_, dustBill := getBill(t, rmaTree, dustCollectorMoneySupplyID)
	require.Equal(t, amount, initialDustCollectorMoneyAmount-dustBill.V)
}

func TestBillData_Value(t *testing.T) {
	bd := &BillData{
		V:        10,
		T:        0,
		Backlink: nil,
	}

	actualSumValue := bd.Value()
	require.Equal(t, rma.Uint64SummaryValue(10), actualSumValue)
}

func TestBillData_AddToHasher(t *testing.T) {
	bd := &BillData{
		V:        10,
		T:        50,
		Backlink: []byte("backlink"),
	}

	hasher := crypto.SHA256.New()
	hasher.Write(util.Uint64ToBytes(bd.V))
	hasher.Write(util.Uint64ToBytes(bd.T))
	hasher.Write(bd.Backlink)
	expectedHash := hasher.Sum(nil)
	hasher.Reset()
	bd.AddToHasher(hasher)
	actualHash := hasher.Sum(nil)
	require.Equal(t, expectedHash, actualHash)
}

func TestBillSummary_Concatenate(t *testing.T) {
	self := rma.Uint64SummaryValue(10)
	left := rma.Uint64SummaryValue(2)
	right := rma.Uint64SummaryValue(3)

	actualSum := self.Concatenate(left, right)
	require.Equal(t, rma.Uint64SummaryValue(15), actualSum)

	actualSum = self.Concatenate(nil, nil)
	require.Equal(t, rma.Uint64SummaryValue(10), actualSum)

	actualSum = self.Concatenate(left, nil)
	require.Equal(t, rma.Uint64SummaryValue(12), actualSum)

	actualSum = self.Concatenate(nil, right)
	require.Equal(t, rma.Uint64SummaryValue(13), actualSum)
}

func TestBillSummary_AddToHasher(t *testing.T) {
	bs := rma.Uint64SummaryValue(10)

	hasher := crypto.SHA256.New()
	hasher.Write(util.Uint64ToBytes(10))
	expectedHash := hasher.Sum(nil)
	hasher.Reset()

	bs.AddToHasher(hasher)
	actualHash := hasher.Sum(nil)
	require.Equal(t, expectedHash, actualHash)
}

func TestEndBlock_DustBillsAreRemoved(t *testing.T) {
	rmaTree, txSystem := createRMATreeAndTxSystem(t)
	_, initBillData := getBill(t, rmaTree, initialBill.ID)
	remaining := initBillData.V
	var splitBillIDs = make([]*uint256.Int, 10)
	backlink := initBillData.Backlink
	for i := 0; i < 10; i++ {
		remaining--
		splitOk, err := NewMoneyTx(createSplit(initialBill.ID, 1, remaining, script.PredicateAlwaysTrue(), backlink))

		require.NoError(t, err)
		roundNumber := uint64(10)
		txSystem.BeginBlock(roundNumber)
		err = txSystem.Execute(splitOk)
		require.NoError(t, err)
		splitWrapper := splitOk.(*billSplitWrapper)
		splitBillIDs[i] = txutil.SameShardId(splitOk.UnitID(), unitIdFromTransaction(splitWrapper))

		_, data := getBill(t, rmaTree, initialBill.ID)
		backlink = data.Backlink
	}

	dcTransfers, swapTx := createDCTransferAndSwapTxs(t, splitBillIDs, rmaTree)

	for _, dcTransfer := range dcTransfers {
		tx, err := NewMoneyTx(dcTransfer)
		require.NoError(t, err)
		err = txSystem.Execute(tx)
		require.NoError(t, err)
	}
	rmaTree.GetRootHash()
	swap, err := NewMoneyTx(swapTx)
	err = txSystem.Execute(swap)
	require.NoError(t, err)
	_, newBillData := getBill(t, rmaTree, swap.UnitID())
	require.Equal(t, uint64(10), newBillData.V)
	_, dustBill := getBill(t, rmaTree, dustCollectorMoneySupplyID)
	require.Equal(t, uint64(10), initialDustCollectorMoneyAmount-dustBill.V)
	_, err = txSystem.EndBlock()
	require.NoError(t, err)
	txSystem.Commit()

	txSystem.BeginBlock(dustBillDeletionTimeout + 10)
	_, err = txSystem.EndBlock()
	require.NoError(t, err)
	txSystem.Commit()

	_, dustBill = getBill(t, rmaTree, dustCollectorMoneySupplyID)
	require.Equal(t, initialDustCollectorMoneyAmount, dustBill.V)
}

func TestValidateSwap_InsufficientDcMoneySupply(t *testing.T) {
	rmaTree, txSystem := createRMATreeAndTxSystem(t)
	roundNumber := uint64(10)
	txSystem.BeginBlock(roundNumber)
	dcTransfers, swapTx := createDCTransferAndSwapTxs(t, []*uint256.Int{initialBill.ID}, rmaTree)

	for _, dcTransfer := range dcTransfers {
		tx, err := NewMoneyTx(dcTransfer)
		require.NoError(t, err)
		err = txSystem.Execute(tx)
		require.NoError(t, err)
	}
	tx, err := NewMoneyTx(swapTx)
	require.NoError(t, err)
	err = txSystem.Execute(tx)
	require.ErrorIs(t, err, ErrSwapInsufficientDCMoneySupply)
}

func TestValidateSwap_SwapBillAlreadyExists(t *testing.T) {
	rmaTree, txSystem := createRMATreeAndTxSystem(t)
	_, initBillData := getBill(t, rmaTree, initialBill.ID)
	roundNumber := uint64(10)
	txSystem.BeginBlock(roundNumber)

	var remaining uint64 = 99
	amount := initialBill.Value - remaining
	splitOk, err := NewMoneyTx(createSplit(initialBill.ID, amount, remaining, script.PredicateAlwaysTrue(), initBillData.Backlink))
	require.NoError(t, err)
	txSystem.BeginBlock(roundNumber)
	err = txSystem.Execute(splitOk)

	splitWrapper := splitOk.(*billSplitWrapper)
	splitBillID := txutil.SameShardId(splitOk.UnitID(), unitIdFromTransaction(splitWrapper))

	dcTransfers, swapTx := createDCTransferAndSwapTxs(t, []*uint256.Int{splitBillID}, rmaTree)

	err = rmaTree.AddItem(uint256.NewInt(0).SetBytes(swapTx.UnitId), script.PredicateAlwaysTrue(), &BillData{}, []byte{})
	require.NoError(t, err)
	for _, dcTransfer := range dcTransfers {
		tx, err := NewMoneyTx(dcTransfer)
		require.NoError(t, err)
		err = txSystem.Execute(tx)
		require.NoError(t, err)
	}
	tx, err := NewMoneyTx(swapTx)
	require.NoError(t, err)
	err = txSystem.Execute(tx)
	require.ErrorIs(t, err, ErrSwapBillAlreadyExists)
}

func TestRegisterData_Revert(t *testing.T) {
	rmaTree, txSystem := createRMATreeAndTxSystem(t)
	_, initBillData := getBill(t, rmaTree, initialBill.ID)

	vdState, err := txSystem.State()
	require.NoError(t, err)

	var remaining uint64 = 10
	amount := initialBill.Value - remaining
	splitOk, err := NewMoneyTx(createSplit(initialBill.ID, amount, remaining, script.PredicateAlwaysTrue(), initBillData.Backlink))
	require.NoError(t, err)
	roundNumber := uint64(10)
	txSystem.BeginBlock(roundNumber)
	err = txSystem.Execute(splitOk)
	require.NoError(t, err)
	_, err = txSystem.State()
	require.ErrorIs(t, err, txsystem.ErrStateContainsUncommittedChanges)

	txSystem.Revert()
	state, err := txSystem.State()
	require.NoError(t, err)
	require.Equal(t, vdState, state)
}

func unitIdFromTransaction(tx *billSplitWrapper) []byte {
	hasher := crypto.SHA256.New()
	idBytes := tx.UnitID().Bytes32()
	hasher.Write(idBytes[:])
	tx.addAttributesToHasher(hasher)
	hasher.Write(util.Uint64ToBytes(tx.Timeout()))
	return hasher.Sum(nil)
}

func getBill(t *testing.T, rmaTree *rma.Tree, billID *uint256.Int) (*rma.Unit, *BillData) {
	t.Helper()
	ib, err := rmaTree.GetUnit(billID)
	require.NoError(t, err)
	require.IsType(t, ib.Data, &BillData{})
	return ib, ib.Data.(*BillData)
}

func createBillTransfer(fromID *uint256.Int, value uint64, backlink []byte) *txsystem.Transaction {
	tx := createTx(fromID)
	bt := &TransferOrder{
		NewBearer: script.PredicateAlwaysTrue(),
		// #nosec G404
		TargetValue: value,
		Backlink:    backlink,
	}
	// #nosec G104
	tx.TransactionAttributes.MarshalFrom(bt)
	return tx
}

func createDCTransferAndSwapTxs(
	t *testing.T,
	ids []*uint256.Int, // bills to swap
	rmaTree *rma.Tree) ([]*txsystem.Transaction, *txsystem.Transaction) {
	t.Helper()

	// calculate new bill ID
	hasher := crypto.SHA256.New()
	idsByteArray := make([][]byte, len(ids))
	for i, id := range ids {
		bytes32 := id.Bytes32()
		hasher.Write(bytes32[:])
		idsByteArray[i] = bytes32[:]
	}
	newBillID := hasher.Sum(nil)

	// create dc transfers
	dcTransfers := make([]*txsystem.Transaction, len(ids))
	var targetValue uint64 = 0
	for i, id := range ids {
		_, billData := getBill(t, rmaTree, id)
		// NB! dc transfer nonce must be equal to swap tx unit id
		targetValue += billData.V
		dcTransfers[i] = createDCTransfer(id, billData.V, billData.Backlink, newBillID)
	}

	tx := &txsystem.Transaction{
		SystemId:              []byte{0, 0, 0, 0},
		UnitId:                newBillID,
		Timeout:               20,
		TransactionAttributes: &anypb.Any{},
		OwnerProof:            script.PredicateArgumentEmpty(),
	}

	bt := &SwapOrder{
		OwnerCondition:  script.PredicateArgumentEmpty(),
		BillIdentifiers: idsByteArray,
		DcTransfers:     dcTransfers,
		// TODO ledger proofs AB-211
		Proofs:      [][]byte{{9}, {10}},
		TargetValue: targetValue,
	}
	// #nosec G104
	tx.TransactionAttributes.MarshalFrom(bt)
	return dcTransfers, tx
}

func createDCTransfer(fromID *uint256.Int, targetValue uint64, backlink []byte, nonce []byte) *txsystem.Transaction {
	tx := createTx(fromID)
	bt := &TransferDCOrder{
		Nonce:        nonce,
		TargetBearer: script.PredicateArgumentEmpty(),
		TargetValue:  targetValue,
		Backlink:     backlink,
	}
	// #nosec G104
	tx.TransactionAttributes.MarshalFrom(bt)
	return tx
}

func createSplit(fromID *uint256.Int, amount uint64, remainingValue uint64, targetBearer, backlink []byte) *txsystem.Transaction {
	tx := createTx(fromID)
	bt := &SplitOrder{
		Amount:         amount,
		TargetBearer:   targetBearer,
		RemainingValue: remainingValue,
		Backlink:       backlink,
	}
	// #nosec G104
	tx.TransactionAttributes.MarshalFrom(bt)
	return tx
}

func createTx(fromID *uint256.Int) *txsystem.Transaction {
	tx := &txsystem.Transaction{
		SystemId:              []byte{0, 0, 0, 0},
		UnitId:                fromID.Bytes(),
		Timeout:               20,
		TransactionAttributes: &anypb.Any{},
		OwnerProof:            script.PredicateArgumentEmpty(),
	}
	return tx
}

func createRMATreeAndTxSystem(t *testing.T) (*rma.Tree, *moneyTxSystem) {
	rmaTree, err := rma.New(&rma.Config{
		HashAlgorithm:     crypto.SHA256,
		RecordingDisabled: false,
	})
	require.NoError(t, err)
	mss, err := NewMoneyTxSystem(
		crypto.SHA256,
		initialBill,
		initialDustCollectorMoneyAmount,
		SchemeOpts.RevertibleState(rmaTree),
	)
	require.NoError(t, err)
	return rmaTree, mss
}