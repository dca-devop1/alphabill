package backend

import (
	"context"
	"errors"
	"time"

	"github.com/alphabill-org/alphabill/internal/block"
	abcrypto "github.com/alphabill-org/alphabill/internal/crypto"
	"github.com/alphabill-org/alphabill/internal/txsystem"
	"github.com/alphabill-org/alphabill/pkg/wallet"
	wlog "github.com/alphabill-org/alphabill/pkg/wallet/log"
	"github.com/alphabill-org/alphabill/pkg/wallet/money/tx_verifier"
)

var alphabillMoneySystemId = []byte{0, 0, 0, 0}

var (
	errKeyNotIndexed  = errors.New("pubkey is not indexed")
	errBillsIsNil     = errors.New("bills input is empty")
	errEmptyBillsList = errors.New("bills list is empty")
)

type (
	WalletBackend struct {
		store         BillStore
		genericWallet *wallet.Wallet
		verifiers     map[string]abcrypto.Verifier
		cancelSyncCh  chan bool
	}

	Bills struct {
		Bills []*Bill `json:"bills"`
	}

	Bill struct {
		Id       []byte `json:"id"`
		Value    uint64 `json:"value"`
		TxHash   []byte `json:"txHash"`
		IsDCBill bool   `json:"isDcBill"`
		// OrderNumber insertion order of given bill in pubkey => list of bills bucket, needed for determistic paging
		OrderNumber uint64   `json:"orderNumber"`
		TxProof     *TxProof `json:"txProof"`
	}

	TxProof struct {
		BlockNumber uint64                `json:"blockNumber"`
		Tx          *txsystem.Transaction `json:"tx"`
		Proof       *block.BlockProof     `json:"proof"`
	}

	Pubkey struct {
		Pubkey     []byte            `json:"pubkey"`
		PubkeyHash *wallet.KeyHashes `json:"pubkeyHash"`
	}

	BillStore interface {
		GetBlockNumber() (uint64, error)
		SetBlockNumber(blockNumber uint64) error
		GetBills(pubKey []byte) ([]*Bill, error)
		RemoveBill(pubKey []byte, id []byte) error
		ContainsBill(pubkey []byte, unitID []byte) (bool, error)
		GetBill(billId []byte) (*Bill, error)
		SetBills(pubkey []byte, bills ...*Bill) error
		GetKeys() ([]*Pubkey, error)
		GetKey(pubkey []byte) (*Pubkey, error)
		AddKey(key *Pubkey) error
	}
)

// New creates a new wallet backend service which can be started by calling the Start or StartProcess method.
// Shutdown method should be called to close resources used by the service.
func New(wallet *wallet.Wallet, store BillStore, verifiers map[string]abcrypto.Verifier) *WalletBackend {
	return &WalletBackend{store: store, genericWallet: wallet, verifiers: verifiers, cancelSyncCh: make(chan bool, 1)}
}

// NewPubkey creates a new hashed Pubkey
func NewPubkey(pubkey []byte) *Pubkey {
	return &Pubkey{
		Pubkey:     pubkey,
		PubkeyHash: wallet.NewKeyHash(pubkey),
	}
}

// Start starts downloading blocks and indexing bills by their owner's public key.
// Blocks forever or until alphabill connection is terminated.
func (w *WalletBackend) Start(ctx context.Context) error {
	blockNumber, err := w.store.GetBlockNumber()
	if err != nil {
		return err
	}
	return w.genericWallet.Sync(ctx, blockNumber)
}

// StartProcess calls Start in a retry loop, can be canceled by cancelling context or calling Shutdown method.
func (w *WalletBackend) StartProcess(ctx context.Context) {
	wlog.Info("starting wallet-backend synchronization")
	defer wlog.Info("wallet-backend synchronization ended")
	retryCount := 0
	for {
		select {
		case <-ctx.Done(): // canceled from context
			return
		case <-w.cancelSyncCh: // canceled from shutdown method
			return
		default:
			if retryCount > 0 {
				wlog.Info("sleeping 10s before retrying alphabill connection")
				time.Sleep(10 * time.Second)
			}
			err := w.Start(ctx)
			if err != nil {
				wlog.Error("error synchronizing wallet-backend: ", err)
			}
			retryCount++
		}
	}
}

// GetBills returns all bills for given public key.
func (w *WalletBackend) GetBills(pubkey []byte) ([]*Bill, error) {
	return w.store.GetBills(pubkey)
}

// GetBill returns most recently seen bill with given unit id.
func (w *WalletBackend) GetBill(unitId []byte) (*Bill, error) {
	return w.store.GetBill(unitId)
}

// SetBill adds new bill to the index.
// Bill most have a valid block proof.
// Overwrites existing bill, if one exists.
// Returns error if given pubkey is not indexed.
func (w *WalletBackend) SetBills(pubkey []byte, bills *block.Bills) error {
	if bills == nil {
		return errBillsIsNil
	}
	if len(bills.Bills) == 0 {
		return errEmptyBillsList
	}
	err := bills.Verify(txConverter, w.verifiers)
	if err != nil {
		return err
	}
	key, err := w.store.GetKey(pubkey)
	if err != nil {
		return err
	}
	if key == nil {
		return errKeyNotIndexed
	}
	pubkeyHash := wallet.NewKeyHash(pubkey)
	domainBills := newBillsFromProto(bills)
	for _, bill := range domainBills {
		tx, err := txConverter.ConvertTx(bill.TxProof.Tx)
		if err != nil {
			return err
		}
		err = txverifier.VerifyTxP2PKHOwner(tx, pubkeyHash)
		if err != nil {
			return err
		}
	}
	return w.store.SetBills(pubkey, domainBills...)
}

// AddKey adds new public key to list of tracked keys.
// Returns ErrKeyAlreadyExists error if key already exists.
func (w *WalletBackend) AddKey(pubkey []byte) error {
	return w.store.AddKey(NewPubkey(pubkey))
}

// Shutdown terminates wallet backend service.
func (w *WalletBackend) Shutdown() {
	// send signal to cancel channel if channel is not full
	select {
	case w.cancelSyncCh <- true:
	default:
	}
	w.genericWallet.Shutdown()
}

func (b *Bill) toProto() *block.Bill {
	return &block.Bill{
		Id:       b.Id,
		Value:    b.Value,
		TxHash:   b.TxHash,
		IsDcBill: b.IsDCBill,
		TxProof:  b.TxProof.toProto(),
	}
}

func (b *TxProof) toProto() *block.TxProof {
	return &block.TxProof{
		BlockNumber: b.BlockNumber,
		Tx:          b.Tx,
		Proof:       b.Proof,
	}
}

func (b *Bill) toProtoBills() *block.Bills {
	return &block.Bills{
		Bills: []*block.Bill{
			b.toProto(),
		},
	}
}

func newBillsFromProto(src *block.Bills) []*Bill {
	dst := make([]*Bill, len(src.Bills))
	for i, b := range src.Bills {
		dst[i] = newBill(b)
	}
	return dst
}

func newBill(b *block.Bill) *Bill {
	return &Bill{
		Id:       b.Id,
		Value:    b.Value,
		TxHash:   b.TxHash,
		IsDCBill: b.IsDcBill,
		TxProof:  newTxProof(b.TxProof),
	}
}

func newTxProof(b *block.TxProof) *TxProof {
	return &TxProof{
		BlockNumber: b.BlockNumber,
		Tx:          b.Tx,
		Proof:       b.Proof,
	}
}
