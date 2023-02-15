package money

import (
	"context"
	"crypto"
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/alphabill-org/alphabill/internal/block"
	"github.com/alphabill-org/alphabill/internal/certificates"
	"github.com/alphabill-org/alphabill/internal/hash"
	"github.com/alphabill-org/alphabill/internal/script"
	testblock "github.com/alphabill-org/alphabill/internal/testutils/block"
	moneytesttx "github.com/alphabill-org/alphabill/internal/testutils/transaction/money"
	"github.com/alphabill-org/alphabill/internal/txsystem"
	"github.com/alphabill-org/alphabill/internal/util"
	"github.com/alphabill-org/alphabill/pkg/wallet/account"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/holiman/uint256"
	"github.com/stretchr/testify/require"
)

const (
	testMnemonic   = "dinosaur simple verify deliver bless ridge monkey design venue six problem lucky"
	testPubKey0Hex = "03c30573dc0c7fd43fcb801289a6a96cb78c27f4ba398b89da91ece23e9a99aca3"
	testPubKey1Hex = "02d36c574db299904b285aaeb57eb7b1fa145c43af90bec3c635c4174c224587b6"
	testPubKey2Hex = "02f6cbeacfd97ebc9b657081eb8b6c9ed3a588646d618ddbd03e198290af94c9d2"
)

func TestExistingWalletCanBeLoaded(t *testing.T) {
	walletDbPath, err := CopyWalletDBFile(t)
	require.NoError(t, err)

	am, err := account.NewManager(walletDbPath, "", true)
	w, err := LoadExistingWallet(WalletConfig{DbPath: walletDbPath}, am)
	require.NoError(t, err)
	t.Cleanup(func() {
		w.Shutdown()
	})
}

func TestWallet_GetPublicKey(t *testing.T) {
	w, _ := CreateTestWalletFromSeed(t)
	pubKey, err := w.am.GetPublicKey(0)
	require.NoError(t, err)
	require.EqualValues(t, "0x"+testPubKey0Hex, hexutil.Encode(pubKey))
}

func TestWallet_GetPublicKeys(t *testing.T) {
	w, _ := CreateTestWalletFromSeed(t)
	_, _, _ = w.AddAccount()

	pubKeys, err := w.am.GetPublicKeys()
	require.NoError(t, err)
	require.Len(t, pubKeys, 2)
	require.EqualValues(t, "0x"+testPubKey0Hex, hexutil.Encode(pubKeys[0]))
	require.EqualValues(t, "0x"+testPubKey1Hex, hexutil.Encode(pubKeys[1]))
}

func TestWallet_AddKey(t *testing.T) {
	w, _ := CreateTestWalletFromSeed(t)

	accIdx, accPubKey, err := w.AddAccount()
	require.NoError(t, err)
	require.EqualValues(t, 1, accIdx)
	require.EqualValues(t, "0x"+testPubKey1Hex, hexutil.Encode(accPubKey))
	accIdx, _ = w.am.GetMaxAccountIndex()
	require.EqualValues(t, 1, accIdx)

	accIdx, accPubKey, err = w.AddAccount()
	require.NoError(t, err)
	require.EqualValues(t, 2, accIdx)
	require.EqualValues(t, "0x"+testPubKey2Hex, hexutil.Encode(accPubKey))
	accIdx, _ = w.am.GetMaxAccountIndex()
	require.EqualValues(t, 2, accIdx)
}

func TestWallet_GetBalance(t *testing.T) {
	w, _ := CreateTestWalletFromSeed(t)
	balance, err := w.GetBalance(GetBalanceCmd{})
	require.NoError(t, err)
	require.EqualValues(t, 0, balance)
}

func TestWallet_GetBalances(t *testing.T) {
	w, _ := CreateTestWalletFromSeed(t)
	_ = w.db.Do().SetBill(0, &Bill{Id: uint256.NewInt(0), Value: 1})
	_ = w.db.Do().SetBill(0, &Bill{Id: uint256.NewInt(1), Value: 1})

	_, _, _ = w.AddAccount()
	_ = w.db.Do().SetBill(1, &Bill{Id: uint256.NewInt(2), Value: 2})
	_ = w.db.Do().SetBill(1, &Bill{Id: uint256.NewInt(3), Value: 2})

	balances, err := w.GetBalances(GetBalanceCmd{})
	require.NoError(t, err)
	require.EqualValues(t, 2, balances[0])
	require.EqualValues(t, 4, balances[1])
}

func TestBlockProcessing(t *testing.T) {
	w, _ := CreateTestWallet(t)

	k, err := w.am.GetAccountKey(0)
	require.NoError(t, err)

	blocks := []*block.Block{
		{
			SystemIdentifier:  w.SystemID(),
			PreviousBlockHash: hash.Sum256([]byte{}),
			Transactions: []*txsystem.Transaction{
				// random dust transfer can be processed
				{
					SystemId:              w.SystemID(),
					UnitId:                hash.Sum256([]byte{0x00}),
					TransactionAttributes: moneytesttx.CreateRandomDustTransferTx(),
					Timeout:               1000,
					OwnerProof:            script.PredicateArgumentEmpty(),
				},
				// receive transfer of 100 bills
				{
					SystemId:              w.SystemID(),
					UnitId:                hash.Sum256([]byte{0x01}),
					TransactionAttributes: moneytesttx.CreateBillTransferTx(k.PubKeyHash.Sha256),
					Timeout:               1000,
					OwnerProof:            script.PredicateArgumentPayToPublicKeyHashDefault([]byte{}, k.PubKey),
				},
				// receive split of 100 bills
				{
					SystemId:              w.SystemID(),
					UnitId:                hash.Sum256([]byte{0x02}),
					TransactionAttributes: moneytesttx.CreateBillSplitTx(k.PubKeyHash.Sha256, 100, 100),
					Timeout:               1000,
					OwnerProof:            script.PredicateArgumentPayToPublicKeyHashDefault([]byte{}, k.PubKey),
				},
				// receive swap of 100 bills
				{
					SystemId:              w.SystemID(),
					UnitId:                hash.Sum256([]byte{0x03}),
					TransactionAttributes: moneytesttx.CreateRandomSwapTransferTx(k.PubKeyHash.Sha256),
					Timeout:               1000,
					OwnerProof:            script.PredicateArgumentPayToPublicKeyHashDefault([]byte{}, k.PubKey),
				},
			},
			UnicityCertificate: &certificates.UnicityCertificate{InputRecord: &certificates.InputRecord{RoundNumber: 1}},
		},
	}

	// verify block number 0 before processing
	blockNumber, err := w.db.Do().GetBlockNumber()
	require.EqualValues(t, 0, blockNumber)
	require.NoError(t, err)

	// verify balance 0 before processing
	balance, err := w.db.Do().GetBalance(GetBalanceCmd{})
	require.EqualValues(t, 0, balance)
	require.NoError(t, err)

	// process blocks
	for _, b := range blocks {
		err = w.ProcessBlock(b)
		require.NoError(t, err)
	}

	// verify block number after block processing
	blockNumber, err = w.db.Do().GetBlockNumber()
	require.EqualValues(t, 1, blockNumber)
	require.NoError(t, err)

	// verify balance after block processing
	balance, err = w.db.Do().GetBalance(GetBalanceCmd{})
	require.EqualValues(t, 300, balance)
	require.NoError(t, err)
}

func TestBlockProcessing_InvalidSystemID(t *testing.T) {
	w, _ := CreateTestWallet(t)

	b := &block.Block{
		SystemIdentifier:   []byte{0, 0, 0, 1},
		PreviousBlockHash:  hash.Sum256([]byte{}),
		Transactions:       []*txsystem.Transaction{},
		UnicityCertificate: &certificates.UnicityCertificate{InputRecord: &certificates.InputRecord{RoundNumber: 1}},
	}

	err := w.ProcessBlock(b)
	require.ErrorContains(t, err, "invalid system identifier")
}

func TestBlockProcessing_VerifyBlockProofs(t *testing.T) {
	w, _ := CreateTestWallet(t)
	k, _ := w.am.GetAccountKey(0)

	testBlock := &block.Block{
		SystemIdentifier:  w.SystemID(),
		PreviousBlockHash: hash.Sum256([]byte{}),
		Transactions: []*txsystem.Transaction{
			// receive transfer of 100 bills
			{
				SystemId:              w.SystemID(),
				UnitId:                hash.Sum256([]byte{0x00}),
				TransactionAttributes: moneytesttx.CreateBillTransferTx(k.PubKeyHash.Sha256),
				Timeout:               1000,
				OwnerProof:            script.PredicateArgumentPayToPublicKeyHashDefault([]byte{}, k.PubKey),
			},
			// receive dc transfer of 100 bills
			{
				SystemId:              w.SystemID(),
				UnitId:                hash.Sum256([]byte{0x01}),
				TransactionAttributes: moneytesttx.CreateDustTransferTx(k.PubKeyHash.Sha256),
				Timeout:               1000,
				OwnerProof:            script.PredicateArgumentPayToPublicKeyHashDefault([]byte{}, k.PubKey),
			},
			// receive split of 100 bills
			{
				SystemId:              w.SystemID(),
				UnitId:                hash.Sum256([]byte{0x02}),
				TransactionAttributes: moneytesttx.CreateBillSplitTx(k.PubKeyHash.Sha256, 100, 100),
				Timeout:               1000,
				OwnerProof:            script.PredicateArgumentPayToPublicKeyHashDefault([]byte{}, k.PubKey),
			},
			// receive swap of 100 bills
			{
				SystemId:              w.SystemID(),
				UnitId:                hash.Sum256([]byte{0x03}),
				TransactionAttributes: moneytesttx.CreateRandomSwapTransferTx(k.PubKeyHash.Sha256),
				Timeout:               1000,
				OwnerProof:            script.PredicateArgumentPayToPublicKeyHashDefault([]byte{}, k.PubKey),
			},
		},
		UnicityCertificate: &certificates.UnicityCertificate{InputRecord: &certificates.InputRecord{RoundNumber: 1}},
	}
	txc := NewTxConverter(w.SystemID())

	certifiedBlock, verifiers := testblock.CertifyBlock(t, testBlock, txc)
	err := w.ProcessBlock(certifiedBlock)
	require.NoError(t, err)

	bills, _ := w.db.Do().GetBills(0)
	require.Len(t, bills, 4)
	for _, b := range bills {
		err = b.BlockProof.Verify(b.GetID(), verifiers, crypto.SHA256, txc)
		require.NoError(t, err)
		require.Equal(t, block.ProofType_PRIM, b.BlockProof.Proof.ProofType)
	}
}

func TestSyncOnClosedWalletShouldNotHang(t *testing.T) {
	w, _ := CreateTestWallet(t)
	addBill(t, w, 100)

	// when wallet is closed
	w.Shutdown()

	// and Sync is called
	err := w.Sync(context.Background())
	require.ErrorContains(t, err, "database not open")
}

func TestWalletDbIsNotCreatedOnWalletCreationError(t *testing.T) {
	// create wallet with invalid seed
	dir := t.TempDir()
	c := WalletConfig{DbPath: dir}
	invalidSeed := "this pond palace oblige remind glory lens popular iron decide coral"
	am, err := account.NewManager(dir, "", true)
	require.NoError(t, err)
	_, err = CreateNewWallet(am, invalidSeed, c)
	require.ErrorContains(t, err, "invalid mnemonic")

	// verify database is not created
	require.False(t, util.FileExists(path.Join(os.TempDir(), WalletFileName)))
}

func TestWalletGetBills_Ok(t *testing.T) {
	w, _ := CreateTestWallet(t)
	addBill(t, w, 100)
	addBill(t, w, 200)
	bills, err := w.GetBills(0)
	require.NoError(t, err)
	require.Len(t, bills, 2)
	require.Equal(t, "0000000000000000000000000000000000000000000000000000000000000064", fmt.Sprintf("%X", bills[0].GetID()))
	require.Equal(t, "00000000000000000000000000000000000000000000000000000000000000C8", fmt.Sprintf("%X", bills[1].GetID()))
}

func TestWalletGetAllBills_Ok(t *testing.T) {
	w, _ := CreateTestWallet(t)
	_, _, _ = w.AddAccount()
	_ = w.db.Do().SetBill(0, &Bill{
		Id:     uint256.NewInt(100),
		Value:  100,
		TxHash: hash.Sum256([]byte{byte(100)}),
	})
	_ = w.db.Do().SetBill(1, &Bill{
		Id:     uint256.NewInt(200),
		Value:  200,
		TxHash: hash.Sum256([]byte{byte(200)}),
	})

	accBills, err := w.GetAllBills()
	require.NoError(t, err)
	require.Len(t, accBills, 2)

	acc0Bills := accBills[0]
	require.Len(t, acc0Bills, 1)
	require.EqualValues(t, acc0Bills[0].Value, 100)

	acc1Bills := accBills[1]
	require.Len(t, acc1Bills, 1)
	require.EqualValues(t, acc1Bills[0].Value, 200)
}

func TestWalletGetBill(t *testing.T) {
	// setup wallet with a bill
	w, _ := CreateTestWallet(t)
	b1 := addBill(t, w, 100)

	// verify getBill returns existing bill
	b, err := w.GetBill(0, b1.GetID())
	require.NoError(t, err)
	require.NotNil(t, b)

	// verify non-existent bill returns BillNotFound error
	b, err = w.GetBill(0, []byte{0})
	require.ErrorIs(t, err, errBillNotFound)
	require.Nil(t, b)
}

func TestWalletAddBill(t *testing.T) {
	// setup wallet
	w, _ := CreateTestWalletFromSeed(t)
	pubkey, _ := w.am.GetPublicKey(0)

	// verify nil bill
	err := w.AddBill(0, nil)
	require.ErrorContains(t, err, "bill is nil")

	// verify bill id is nil
	err = w.AddBill(0, &Bill{Id: nil})
	require.ErrorContains(t, err, "bill id is nil")

	// verify bill tx is nil
	err = w.AddBill(0, &Bill{Id: uint256.NewInt(0)})
	require.ErrorContains(t, err, "bill tx hash is nil")

	// verify bill block proof is nil
	err = w.AddBill(0, &Bill{Id: uint256.NewInt(0), TxHash: []byte{}})
	require.ErrorContains(t, err, "bill block proof is nil")

	err = w.AddBill(0, &Bill{Id: uint256.NewInt(0), TxHash: []byte{}, BlockProof: &BlockProof{}})
	require.ErrorContains(t, err, "bill block proof tx is nil")

	// verify invalid bearer predicate
	invalidPubkey := []byte{0}
	err = w.AddBill(0, &Bill{
		Id:         uint256.NewInt(0),
		TxHash:     []byte{},
		BlockProof: &BlockProof{Tx: createTransferTxForPubKey(w.SystemID(), invalidPubkey)},
	})
	require.ErrorContains(t, err, "invalid bearer predicate")

	// verify valid bill no error
	err = w.AddBill(0, &Bill{
		Id:         uint256.NewInt(0),
		TxHash:     []byte{},
		BlockProof: &BlockProof{Tx: createTransferTxForPubKey(w.SystemID(), pubkey)},
	})
	require.NoError(t, err)
}

func createTransferTxForPubKey(systemId, pubkey []byte) *txsystem.Transaction {
	return &txsystem.Transaction{
		SystemId:              systemId,
		UnitId:                hash.Sum256([]byte{0x01}),
		TransactionAttributes: moneytesttx.CreateBillTransferTx(hash.Sum256(pubkey)),
		Timeout:               1000,
		OwnerProof:            script.PredicateArgumentPayToPublicKeyHashDefault([]byte{}, pubkey),
	}
}
