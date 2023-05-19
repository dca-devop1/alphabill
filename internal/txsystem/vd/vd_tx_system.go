package vd

import (
	"bytes"
	"crypto"
	"hash"
	"fmt"

	"github.com/alphabill-org/alphabill/internal/errors"
	hasherUtil "github.com/alphabill-org/alphabill/internal/hash"
	"github.com/alphabill-org/alphabill/internal/rma"
	"github.com/alphabill-org/alphabill/internal/script"
	"github.com/alphabill-org/alphabill/internal/txsystem"
	"github.com/alphabill-org/alphabill/internal/txsystem/fc/transactions"
	"github.com/alphabill-org/alphabill/internal/util"
	"github.com/alphabill-org/alphabill/pkg/logger"
	"github.com/holiman/uint256"
)

var DefaultSystemIdentifier = []byte{0, 0, 0, 1}

const zeroSummaryValue = rma.Uint64SummaryValue(0)

var (
	ErrOwnerProofPresent = errors.New("'register data' transaction cannot have an owner proof")
	log                  = logger.CreateForPackage()
	zeroRootHash         = make([]byte, 32)
)

type (
	vdTransaction struct {
		transaction *txsystem.Transaction
		hashFunc    crypto.Hash
		hashValue   []byte
	}

	txSystem struct {
		systemIdentifier   []byte
		stateTree          *rma.Tree
		hashAlgorithm      crypto.Hash
		currentBlockNumber uint64
	}

	unit struct {
		dataHash    []byte
		blockNumber uint64
	}
)

func NewTxSystem(systemId []byte) (*txSystem, error) {
	conf := &rma.Config{HashAlgorithm: crypto.SHA256}
	s, err := rma.New(conf)
	if err != nil {
		return nil, err
	}

	vdTxSystem := &txSystem{
		systemIdentifier: systemId,
		stateTree:        s,
		hashAlgorithm:    conf.HashAlgorithm,
	}

	return vdTxSystem, nil
}

// NewVDTx creates a new wrapper, returns an error if unknown transaction type is given as argument.
func NewVDTx(systemID []byte, tx *txsystem.Transaction) (txsystem.GenericTransaction, error) {
	if !bytes.Equal(systemID, tx.GetSystemId()) {
		return nil, fmt.Errorf("transaction has invalid system identifier %X, expected %X", tx.GetSystemId(), systemID)
	}
	if tx.TransactionAttributes != nil {
		feeTx, err := transactions.NewFeeCreditTx(tx)
		if err != nil {
			return nil, err
		}
		if feeTx != nil {
			return feeTx, nil
		}
		return nil, errors.New("invalid vd transaction: transactionAttributes present")
	}

	return &vdTransaction{
		transaction: tx,
	}, nil
}

func (d *txSystem) State() (txsystem.State, error) {
	if d.stateTree.ContainsUncommittedChanges() {
		return nil, txsystem.ErrStateContainsUncommittedChanges
	}
	return d.getState(), nil
}

func (d *txSystem) BeginBlock(blockNumber uint64) {
	d.currentBlockNumber = blockNumber
}

func (d *txSystem) EndBlock() (txsystem.State, error) {
	return d.getState(), nil
}

func (d *txSystem) Revert() {
	d.stateTree.Revert()
}

func (d *txSystem) Commit() {
	d.stateTree.Commit()
}

func (d *txSystem) Execute(tx txsystem.GenericTransaction) error {
	log.Debug("Processing register data tx: '%v', UnitID=%x", tx, tx.UnitID())
	if len(tx.OwnerProof()) > 0 {
		return ErrOwnerProofPresent
	}
	h := tx.Hash(d.hashAlgorithm)
	err := d.stateTree.AtomicUpdate(
		rma.AddItem(tx.UnitID(),
			script.PredicateAlwaysFalse(),
			&unit{
				dataHash:    hasherUtil.Sum256(util.Uint256ToBytes(tx.UnitID())),
				blockNumber: d.currentBlockNumber,
			},
			h,
		))
	if err != nil {
		return errors.Wrapf(err, "could not add item: %v", err)
	}
	return nil
}

func (d *txSystem) ConvertTx(tx *txsystem.Transaction) (txsystem.GenericTransaction, error) {
	return NewVDTx(d.systemIdentifier, tx)
}

func (d *txSystem) getState() txsystem.State {
	if d.stateTree.GetRootHash() == nil {
		return txsystem.NewStateSummary(zeroRootHash, zeroSummaryValue.Bytes())
	}
	return txsystem.NewStateSummary(d.stateTree.GetRootHash(), zeroSummaryValue.Bytes())
}

func (u *unit) AddToHasher(hasher hash.Hash) {
	hasher.Write(u.dataHash)
	hasher.Write(util.Uint64ToBytes(u.blockNumber))
}

func (u *unit) Value() rma.SummaryValue {
	return zeroSummaryValue
}

func (w *vdTransaction) Hash(hashFunc crypto.Hash) []byte {
	if w.hashComputed(hashFunc) {
		return w.hashValue
	}
	hasher := hashFunc.New()
	w.AddToHasher(hasher)

	w.hashValue = hasher.Sum(nil)
	w.hashFunc = hashFunc
	return w.hashValue
}

func (w *vdTransaction) AddToHasher(hasher hash.Hash) {
	hasher.Write(w.transaction.Bytes())
}

func (w *vdTransaction) SigBytes() []byte {
	return nil
}

func (w *vdTransaction) UnitID() *uint256.Int {
	return uint256.NewInt(0).SetBytes(w.transaction.UnitId)
}

func (w *vdTransaction) Timeout() uint64 {
	return w.transaction.Timeout()
}

func (w *vdTransaction) SystemID() []byte {
	return w.transaction.SystemId
}

func (w *vdTransaction) OwnerProof() []byte {
	return w.transaction.OwnerProof
}

func (w *vdTransaction) ToProtoBuf() *txsystem.Transaction {
	return w.transaction
}

func (w *vdTransaction) IsPrimary() bool {
	return true
}

func (w *vdTransaction) TargetUnits(_ crypto.Hash) []*uint256.Int {
	return []*uint256.Int{w.UnitID()}
}

func (w *vdTransaction) SetServerMetadata(sm *txsystem.ServerMetadata) {
	w.ToProtoBuf().ServerMetadata = sm
	w.resetHasher()
}

func (w *vdTransaction) resetHasher() {
	w.hashValue = nil
}

func (w *vdTransaction) sigBytes(b *bytes.Buffer) {
	b.Write(w.transaction.SystemId)
	b.Write(w.transaction.UnitId)
	if w.transaction.ClientMetadata != nil {
		b.Write(w.transaction.ClientMetadata.Bytes())
	}
}

func (w *vdTransaction) hashComputed(hashFunc crypto.Hash) bool {
	return w.hashFunc == hashFunc && w.hashValue != nil
}
