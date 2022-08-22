package money

import (
	"bytes"
	"context"
	"crypto"
	"errors"
	"fmt"
	"sort"
	"strconv"

	"github.com/alphabill-org/alphabill/internal/block"
	abcrypto "github.com/alphabill-org/alphabill/internal/crypto"
	"github.com/alphabill-org/alphabill/internal/txsystem"
	"github.com/alphabill-org/alphabill/internal/txsystem/money"
	moneytx "github.com/alphabill-org/alphabill/internal/txsystem/money"
	"github.com/alphabill-org/alphabill/internal/txsystem/util"
	"github.com/alphabill-org/alphabill/pkg/wallet"
	"github.com/alphabill-org/alphabill/pkg/wallet/log"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/holiman/uint256"
	"github.com/robfig/cron/v3"
	"golang.org/x/sync/errgroup"
)

const (
	dcTimeoutBlockCount     = 10
	swapTimeoutBlockCount   = 60
	txTimeoutBlockCount     = 100
	dustBillDeletionTimeout = 300
)

var (
	ErrSwapInProgress       = errors.New("swap is in progress, synchronize your wallet to complete the process")
	ErrSwapNotEnoughBills   = errors.New("need to have more than 1 bill to perform swap")
	ErrInsufficientBalance  = errors.New("insufficient balance for transaction")
	ErrInvalidPubKey        = errors.New("invalid public key, public key must be in compressed secp256k1 format")
	ErrInvalidPassword      = errors.New("invalid password")
	ErrInvalidBlockSystemID = errors.New("invalid system identifier")
)

type (
	Wallet struct {
		*wallet.Wallet

		config           WalletConfig
		db               Db
		dustCollectorJob *cron.Cron
		dcWg             *dcWaitGroup
		accounts         *accounts
	}
)

// CreateNewWallet creates a new wallet. To synchronize wallet with a node call Sync.
// Shutdown needs to be called to release resources used by wallet.
// If mnemonic seed is empty then new mnemonic will ge generated, otherwise wallet is restored using given mnemonic.
func CreateNewWallet(mnemonic string, config WalletConfig) (*Wallet, error) {
	db, err := getDb(config, true)
	if err != nil {
		return nil, err
	}
	return createMoneyWallet(config, db, mnemonic)
}

func LoadExistingWallet(config WalletConfig) (*Wallet, error) {
	db, err := getDb(config, false)
	if err != nil {
		return nil, err
	}

	ok, err := db.Do().VerifyPassword()
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrInvalidPassword
	}

	accountKeys, err := db.Do().GetAccountKeys()
	if err != nil {
		return nil, err
	}
	accs := make([]account, len(accountKeys))
	for num, val := range accountKeys {
		accs[num] = account{
			accountNumber: uint64(num),
			accountKeys:   *val.PubKeyHash,
		}
	}
	mw := &Wallet{config: config, db: db, dustCollectorJob: cron.New(), dcWg: newDcWaitGroup(), accounts: &accounts{accounts: accs}}

	mw.Wallet = wallet.New().
		SetBlockProcessor(mw).
		SetABClientConf(config.AlphabillClientConfig).
		Build()

	return mw, nil
}

// IsEncrypted returns true if wallet exists and is encrypted and or false if wallet exists and is not encrypted,
// returns error if wallet does not exist.
func IsEncrypted(config WalletConfig) (bool, error) {
	db, err := getDb(config, false)
	if err != nil {
		return false, err
	}
	defer db.Close()
	return db.Do().IsEncrypted()
}

func (w *Wallet) ProcessBlock(b *block.Block) error {
	log.Info("processing block: " + strconv.FormatUint(b.BlockNumber, 10))
	if !bytes.Equal(alphabillMoneySystemId, b.GetSystemIdentifier()) {
		return ErrInvalidBlockSystemID
	}

	return w.db.WithTransaction(func(dbTx TxContext) error {
		lastBlockNumber, err := w.db.Do().GetBlockNumber()
		if err != nil {
			return err
		}
		err = validateBlockNumber(b.BlockNumber, lastBlockNumber)
		if err != nil {
			return err
		}
		for _, acc := range w.accounts.getAll() {
			for i, pbTx := range b.Transactions {
				err = w.collectBills(dbTx, pbTx, b, i, &acc)
				if err != nil {
					return err
				}
			}
		}
		return w.endBlock(dbTx, b)
	})
}

func (w *Wallet) endBlock(dbTx TxContext, b *block.Block) error {
	blockNumber := b.BlockNumber
	err := dbTx.SetBlockNumber(blockNumber)
	if err != nil {
		return err
	}
	for _, acc := range w.accounts.getAll() {
		err = w.deleteExpiredDcBills(dbTx, blockNumber, acc.accountNumber)
		if err != nil {
			return err
		}
		err = w.trySwap(dbTx, acc.accountNumber)
		if err != nil {
			return err
		}
		err = w.dcWg.DecrementSwaps(dbTx, blockNumber, acc.accountNumber)
		if err != nil {
			return err
		}
	}
	return nil
}

// Shutdown terminates connection to alphabill node, closes wallet db, cancels dust collector job and any background goroutines.
func (w *Wallet) Shutdown() {
	w.Wallet.Shutdown()

	if w.dustCollectorJob != nil {
		w.dustCollectorJob.Stop()
	}
	if w.dcWg != nil {
		w.dcWg.ResetWaitGroup()
	}
	if w.db != nil {
		w.db.Close()
	}
}

// DeleteDb deletes the wallet database.
func (w *Wallet) DeleteDb() {
	w.db.DeleteDb()
}

// CollectDust starts the dust collector process.
// Wallet needs to be synchronizing using Sync or SyncToMaxBlockNumber in order to receive transactions and finish the process.
// The function blocks until dust collector process is finished or timed out.
func (w *Wallet) CollectDust(ctx context.Context) error {
	errgrp, ctx := errgroup.WithContext(ctx)
	for _, acc := range w.accounts.getAll() {
		errgrp.Go(func() error {
			return w.collectDust(ctx, true, acc.accountNumber)
		})
	}
	return errgrp.Wait()
}

// StartDustCollectorJob starts the dust collector background process that runs every hour until wallet is shut down.
// Wallet needs to be synchronizing using Sync or SyncToMaxBlockNumber in order to receive transactions and finish the process.
// Returns error if the job failed to start.
func (w *Wallet) StartDustCollectorJob() error {
	_, err := w.startDustCollectorJob()
	return err
}

// GetBalance returns sum value of all bills currently owned by the wallet, for given account
// the value returned is the smallest denomination of alphabills.
func (w *Wallet) GetBalance(accountNumber uint64) (uint64, error) {
	return w.db.Do().GetBalance(accountNumber)
}

// GetBalances returns sum value of all bills currently owned by the wallet, for all accounts
// the value returned is the smallest denomination of alphabills.
func (w *Wallet) GetBalances() ([]uint64, error) {
	return w.db.Do().GetBalances()
}

// GetPublicKey returns public key of the wallet (compressed secp256k1 key 33 bytes)
func (w *Wallet) GetPublicKey(accountNumber uint64) ([]byte, error) {
	key, err := w.db.Do().GetAccountKey(accountNumber)
	if err != nil {
		return nil, err
	}
	return key.PubKey, nil
}

// GetPublicKeys returns public keys of the wallet, indexed by account numbers
func (w *Wallet) GetPublicKeys() ([][]byte, error) {
	accKeys, err := w.db.Do().GetAccountKeys()
	if err != nil {
		return nil, err
	}
	pubKeys := make([][]byte, len(accKeys))
	for accNum, accKey := range accKeys {
		pubKeys[accNum] = accKey.PubKey
	}
	return pubKeys, nil
}

// GetMnemonic returns mnemonic seed of the wallet
func (w *Wallet) GetMnemonic() (string, error) {
	return w.db.Do().GetMnemonic()
}

// AddAccount adds the next account in account key series to the wallet.
// New accounts are indexed only from the time of creation and not backwards in time.
// Returns added account number together with account's public key.
func (w *Wallet) AddAccount() (uint64, []byte, error) {
	masterKeyString, err := w.db.Do().GetMasterKey()
	if err != nil {
		return 0, nil, err
	}
	masterKey, err := hdkeychain.NewKeyFromString(masterKeyString)
	if err != nil {
		return 0, nil, err
	}

	accountNumber, err := w.db.Do().GetMaxAccountNumber()
	if err != nil {
		return 0, nil, err
	}
	accountNumber += 1

	derivationPath := wallet.NewDerivationPath(accountNumber)
	accountKey, err := wallet.NewAccountKey(masterKey, derivationPath)
	if err != nil {
		return 0, nil, err
	}
	err = w.db.WithTransaction(func(tx TxContext) error {
		err := tx.AddAccount(accountNumber, accountKey)
		if err != nil {
			return err
		}
		err = tx.SetMaxAccountNumber(accountNumber)
		if err != nil {
			return err
		}
		w.accounts.add(&account{accountNumber: accountNumber, accountKeys: *accountKey.PubKeyHash})
		return nil
	})
	if err != nil {
		return 0, nil, err
	}
	return accountNumber, accountKey.PubKey, nil
}

// Send creates, signs and broadcasts transactions, in total for the given amount,
// to the given public key, the public key must be in compressed secp256k1 format.
// Sends one transaction per bill, prioritzing larger bills.
func (w *Wallet) Send(receiverPubKey []byte, amount uint64, accountNumber uint64) error {
	if len(receiverPubKey) != abcrypto.CompressedSecp256K1PublicKeySize {
		return ErrInvalidPubKey
	}

	swapInProgress, err := w.isSwapInProgress(w.db.Do())
	if err != nil {
		return err
	}
	if swapInProgress {
		return ErrSwapInProgress
	}

	balance, err := w.GetBalance(accountNumber)
	if err != nil {
		return err
	}
	if amount > balance {
		return ErrInsufficientBalance
	}

	maxBlockNo, err := w.GetMaxBlockNumber()
	if err != nil {
		return err
	}
	timeout := maxBlockNo + txTimeoutBlockCount
	if err != nil {
		return err
	}

	k, err := w.db.Do().GetAccountKey(accountNumber)
	if err != nil {
		return err
	}

	bills, err := w.db.Do().GetBills(accountNumber)
	if err != nil {
		return err
	}

	txs, err := createTransactions(receiverPubKey, amount, bills, k, timeout)
	if err != nil {
		return err
	}
	for _, tx := range txs {
		res, err := w.SendTransaction(tx)
		if err != nil {
			return err
		}
		if !res.Ok {
			return errors.New("payment returned error code: " + res.Message)
		}
	}
	return nil
}

// Sync synchronises wallet from the last known block number with the given alphabill node.
// The function blocks forever or until alphabill connection is terminated.
// Returns immediately if already synchronizing.
func (w *Wallet) Sync(ctx context.Context) error {
	blockNumber, err := w.db.Do().GetBlockNumber()
	if err != nil {
		return err
	}
	return w.Wallet.Sync(ctx, blockNumber)
}

// Sync synchronises wallet from the last known block number with the given alphabill node.
// The function blocks until maximum block height, calculated at the start of the process, is reached.
// Returns immediately with ErrWalletAlreadySynchronizing if already synchronizing.
func (w *Wallet) SyncToMaxBlockNumber(ctx context.Context) error {
	blockNumber, err := w.db.Do().GetBlockNumber()
	if err != nil {
		return err
	}
	return w.Wallet.SyncToMaxBlockNumber(ctx, blockNumber)
}

func (w *Wallet) collectBills(dbTx TxContext, txPb *txsystem.Transaction, b *block.Block, txIdx int, acc *account) error {
	gtx, err := moneytx.NewMoneyTx(alphabillMoneySystemId, txPb)
	if err != nil {
		return err
	}
	stx := gtx.(txsystem.GenericTransaction)
	switch tx := stx.(type) {
	case money.Transfer:
		isOwner, err := verifyOwner(acc, tx.NewBearer())
		if err != nil {
			return err
		}
		if isOwner {
			log.Info("received transfer order")
			err := w.saveWithProof(dbTx, b, txIdx, &bill{
				Id:     tx.UnitID(),
				Value:  tx.TargetValue(),
				TxHash: tx.Hash(crypto.SHA256),
			}, acc.accountNumber)
			if err != nil {
				return err
			}
		} else {
			err := dbTx.RemoveBill(acc.accountNumber, tx.UnitID())
			if err != nil {
				return err
			}
		}
	case money.TransferDC:
		isOwner, err := verifyOwner(acc, tx.TargetBearer())
		if err != nil {
			return err
		}
		if isOwner {
			log.Info("received TransferDC order")
			err := w.saveWithProof(dbTx, b, txIdx, &bill{
				Id:                  tx.UnitID(),
				Value:               tx.TargetValue(),
				TxHash:              tx.Hash(crypto.SHA256),
				IsDcBill:            true,
				DcTx:                txPb,
				DcTimeout:           tx.Timeout(),
				DcNonce:             tx.Nonce(),
				DcExpirationTimeout: b.BlockNumber + dustBillDeletionTimeout,
			}, acc.accountNumber)
			if err != nil {
				return err
			}
		} else {
			err := dbTx.RemoveBill(acc.accountNumber, tx.UnitID())
			if err != nil {
				return err
			}
		}
	case money.Split:
		// split tx contains two bills: existing bill and new bill
		// if any of these bills belong to wallet then we have to
		// 1) update the existing bill and
		// 2) add the new bill
		containsBill, err := dbTx.ContainsBill(acc.accountNumber, tx.UnitID())
		if err != nil {
			return err
		}
		if containsBill {
			log.Info("received split order (existing bill)")
			err := w.saveWithProof(dbTx, b, txIdx, &bill{
				Id:     tx.UnitID(),
				Value:  tx.RemainingValue(),
				TxHash: tx.Hash(crypto.SHA256),
			}, acc.accountNumber)
			if err != nil {
				return err
			}
		}
		isOwner, err := verifyOwner(acc, tx.TargetBearer())
		if err != nil {
			return err
		}
		if isOwner {
			log.Info("received split order (new bill)")
			err := w.saveWithProof(dbTx, b, txIdx, &bill{
				Id:     util.SameShardId(tx.UnitID(), tx.HashForIdCalculation(crypto.SHA256)),
				Value:  tx.Amount(),
				TxHash: tx.Hash(crypto.SHA256),
			}, acc.accountNumber)
			if err != nil {
				return err
			}
		}
	case money.Swap:
		isOwner, err := verifyOwner(acc, tx.OwnerCondition())
		if err != nil {
			return err
		}
		if isOwner {
			log.Info("received swap order")
			err := w.saveWithProof(dbTx, b, txIdx, &bill{
				Id:     tx.UnitID(),
				Value:  tx.TargetValue(),
				TxHash: tx.Hash(crypto.SHA256),
			}, acc.accountNumber)
			if err != nil {
				return err
			}
			// clear dc metadata
			err = dbTx.SetDcMetadata(txPb.UnitId, nil)
			if err != nil {
				return err
			}
			for _, dustTransfer := range tx.DCTransfers() {
				err := dbTx.RemoveBill(acc.accountNumber, dustTransfer.UnitID())
				if err != nil {
					return err
				}
			}
		} else {
			err := dbTx.RemoveBill(acc.accountNumber, tx.UnitID())
			if err != nil {
				return err
			}
		}
	default:
		log.Warning(fmt.Sprintf("received unknown transaction type, skipping processing: %s", tx))
		return nil
	}
	return nil
}

func (w *Wallet) saveWithProof(dbTx TxContext, b *block.Block, txIdx int, bill *bill, accountNumber uint64) error {
	blockProof, err := ExtractBlockProof(b, txIdx, crypto.SHA256)
	if err != nil {
		return err
	}
	bill.BlockProof = blockProof
	return dbTx.SetBill(accountNumber, bill)
}

func (w *Wallet) deleteExpiredDcBills(dbTx TxContext, blockNumber uint64, accountNumber uint64) error {
	bills, err := dbTx.GetBills(accountNumber)
	if err != nil {
		return err
	}
	for _, b := range bills {
		if b.isExpired(blockNumber) {
			log.Info(fmt.Sprintf("deleting expired dc bill: value=%d id=%s", b.Value, b.Id.String()))
			err = dbTx.RemoveBill(accountNumber, b.Id)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (w *Wallet) trySwap(tx TxContext, accountNumber uint64) error {
	blockHeight, err := tx.GetBlockNumber()
	if err != nil {
		return err
	}
	maxBlockNo, err := w.GetMaxBlockNumber()
	if err != nil {
		return err
	}
	bills, err := tx.GetBills(accountNumber)
	if err != nil {
		return err
	}
	dcBillGroups := groupDcBills(bills)
	for nonce, billGroup := range dcBillGroups {
		nonce32 := nonce.Bytes32()
		dcMeta, err := tx.GetDcMetadata(nonce32[:])
		if err != nil {
			return err
		}
		if dcMeta != nil && dcMeta.isSwapRequired(blockHeight, billGroup.valueSum) {
			timeout := maxBlockNo + swapTimeoutBlockCount
			err := w.swapDcBills(tx, billGroup.dcBills, billGroup.dcNonce, timeout, accountNumber)
			if err != nil {
				return err
			}
			w.dcWg.UpdateTimeout(billGroup.dcNonce, timeout)
		}
	}

	// delete expired metadata
	nonceMetadataMap, err := tx.GetDcMetadataMap()
	if err != nil {
		return err
	}
	for nonce, m := range nonceMetadataMap {
		if m.timeoutReached(blockHeight) {
			nonce32 := nonce.Bytes32()
			err := tx.SetDcMetadata(nonce32[:], nil)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// collectDust sends dust transfer for every bill for given account in wallet and records metadata.
// Once the dust transfers get confirmed on the ledger then swap transfer is broadcast and metadata cleared.
// If blocking is true then the function blocks until swap has been completed or timed out,
// if blocking is false then the function returns after sending the dc transfers.
func (w *Wallet) collectDust(ctx context.Context, blocking bool, accountNumber uint64) error {
	err := w.db.WithTransaction(func(dbTx TxContext) error {
		blockHeight, err := dbTx.GetBlockNumber()
		if err != nil {
			return err
		}
		maxBlockNo, err := w.GetMaxBlockNumber()
		if err != nil {
			return err
		}
		bills, err := dbTx.GetBills(accountNumber)
		if err != nil {
			return err
		}
		if len(bills) < 2 {
			return ErrSwapNotEnoughBills
		}
		var expectedSwaps []expectedSwap
		dcBillGroups := groupDcBills(bills)
		if len(dcBillGroups) > 0 {
			for _, v := range dcBillGroups {
				if blockHeight >= v.dcTimeout {
					swapTimeout := maxBlockNo + swapTimeoutBlockCount
					err = w.swapDcBills(dbTx, v.dcBills, v.dcNonce, swapTimeout, accountNumber)
					if err != nil {
						return err
					}
					expectedSwaps = append(expectedSwaps, expectedSwap{dcNonce: v.dcNonce, timeout: swapTimeout})
				} else {
					// expecting to receive swap during dcTimeout
					expectedSwaps = append(expectedSwaps, expectedSwap{dcNonce: v.dcNonce, timeout: v.dcTimeout})
				}
			}
		} else {
			swapInProgress, err := w.isSwapInProgress(dbTx)
			if err != nil {
				return err
			}
			if swapInProgress {
				return ErrSwapInProgress
			}

			k, err := dbTx.GetAccountKey(accountNumber)
			if err != nil {
				return err
			}

			dcNonce := calculateDcNonce(bills)
			dcTimeout := maxBlockNo + dcTimeoutBlockCount
			var dcValueSum uint64
			for _, b := range bills {
				dcValueSum += b.Value
				tx, err := createDustTx(k, b, dcNonce, dcTimeout)
				if err != nil {
					return err
				}

				log.Info("sending dust transfer tx for bill ", b.Id)
				res, err := w.SendTransaction(tx)
				if err != nil {
					return err
				}
				if !res.Ok {
					return errors.New("dust transfer returned error code: " + res.Message)
				}
			}
			expectedSwaps = append(expectedSwaps, expectedSwap{dcNonce: dcNonce, timeout: dcTimeout})
			err = dbTx.SetDcMetadata(dcNonce, &dcMetadata{
				DcValueSum: dcValueSum,
				DcTimeout:  dcTimeout,
			})
			if err != nil {
				return err
			}
		}
		if blocking {
			w.dcWg.AddExpectedSwaps(expectedSwaps)
		}
		return nil
	})
	if err != nil {
		return err
	}
	if blocking {
		log.Info("waiting for blocking collect dust (wallet needs to be synchronizing to finish this process)")

		// wrap wg.Wait() as channel
		done := make(chan struct{})
		go func() {
			w.dcWg.wg.Wait()
			close(done)
		}()

		select {
		case <-ctx.Done():
			// context canceled externally
		case <-done:
			// dust collection finished (swap received or timed out)
		}
		log.Info("finished waiting for blocking collect dust")
	}
	return nil
}

func (w *Wallet) swapDcBills(tx TxContext, dcBills []*bill, dcNonce []byte, timeout uint64, accountNumber uint64) error {
	k, err := tx.GetAccountKey(accountNumber)
	if err != nil {
		return err
	}
	swap, err := createSwapTx(k, dcBills, dcNonce, timeout)
	if err != nil {
		return err
	}
	log.Info(fmt.Sprintf("sending swap tx: nonce=%s timeout=%d", hexutil.Encode(dcNonce), timeout))
	res, err := w.SendTransaction(swap)
	if err != nil {
		return err
	}
	if !res.Ok {
		return errors.New("swap tx returned error code: " + res.Message)
	}
	return tx.SetDcMetadata(dcNonce, &dcMetadata{SwapTimeout: timeout})
}

// isSwapInProgress returns true if there's a running dc process managed by the wallet
func (w *Wallet) isSwapInProgress(dbTx TxContext) (bool, error) {
	blockHeight, err := dbTx.GetBlockNumber()
	if err != nil {
		return false, err
	}
	dcMetadataMap, err := dbTx.GetDcMetadataMap()
	if err != nil {
		return false, err
	}
	for _, m := range dcMetadataMap {
		if m.DcValueSum > 0 { // value sum is set only for dc process that was started by wallet
			return blockHeight < m.DcTimeout || blockHeight < m.SwapTimeout, nil
		}
	}
	return false, nil
}

func (w *Wallet) startDustCollectorJob() (cron.EntryID, error) {
	return w.dustCollectorJob.AddFunc("@hourly", func() {
		for _, acc := range w.accounts.getAll() {
			err := w.collectDust(context.Background(), false, acc.accountNumber)
			if err != nil {
				log.Error("error in dust collector job: ", err)
			}
		}
	})
}

func createMoneyWallet(config WalletConfig, db Db, mnemonic string) (mw *Wallet, err error) {
	mw = &Wallet{config: config, db: db, dustCollectorJob: cron.New(), dcWg: newDcWaitGroup(), accounts: newAccountsCache()}
	defer func() {
		if err != nil {
			// delete database if any error occurs after creating it
			mw.DeleteDb()
		}
	}()

	keys, err := wallet.NewKeys(mnemonic)
	if err != nil {
		return
	}

	mw.Wallet = wallet.New().
		SetBlockProcessor(mw).
		SetABClientConf(config.AlphabillClientConfig).
		Build()

	err = saveKeys(db, keys, config.WalletPass)
	if err != nil {
		return
	}

	mw.accounts.add(&account{
		accountNumber: 0,
		accountKeys:   *keys.AccountKey.PubKeyHash,
	})
	return
}

func calculateDcNonce(bills []*bill) []byte {
	var billIds [][]byte
	for _, b := range bills {
		billIds = append(billIds, b.getId())
	}

	// sort billIds in ascending order
	sort.Slice(billIds, func(i, j int) bool {
		return bytes.Compare(billIds[i], billIds[j]) < 0
	})

	hasher := crypto.Hash.New(crypto.SHA256)
	for _, billId := range billIds {
		hasher.Write(billId)
	}
	return hasher.Sum(nil)
}

// groupDcBills groups bills together by dc nonce
func groupDcBills(bills []*bill) map[uint256.Int]*dcBillGroup {
	m := map[uint256.Int]*dcBillGroup{}
	for _, b := range bills {
		if b.IsDcBill {
			k := *uint256.NewInt(0).SetBytes(b.DcNonce)
			billContainer, exists := m[k]
			if !exists {
				billContainer = &dcBillGroup{}
				m[k] = billContainer
			}
			billContainer.valueSum += b.Value
			billContainer.dcBills = append(billContainer.dcBills, b)
			billContainer.dcNonce = b.DcNonce
			billContainer.dcTimeout = b.DcTimeout
		}
	}
	return m
}

func validateBlockNumber(blockNumber uint64, lastBlockNumber uint64) error {
	// verify that we are processing blocks sequentially
	// TODO verify last prev block hash?
	if blockNumber-lastBlockNumber != 1 {
		return errors.New(fmt.Sprintf("Invalid block height. Received blockNumber %d current wallet blockNumber %d", blockNumber, lastBlockNumber))
	}
	return nil
}

func getDb(config WalletConfig, create bool) (Db, error) {
	if config.Db != nil {
		return config.Db, nil
	}
	if create {
		return createNewDb(config)
	}
	return OpenDb(config)
}

func saveKeys(db Db, keys *wallet.Keys, walletPass string) error {
	return db.WithTransaction(func(tx TxContext) error {
		err := tx.SetEncrypted(walletPass != "")
		if err != nil {
			return err
		}
		err = tx.SetMnemonic(keys.Mnemonic)
		if err != nil {
			return err
		}
		err = tx.SetMasterKey(keys.MasterKey.String())
		if err != nil {
			return err
		}
		err = tx.AddAccount(0, keys.AccountKey)
		if err != nil {
			return err
		}
		return tx.SetMaxAccountNumber(0)
	})
}
