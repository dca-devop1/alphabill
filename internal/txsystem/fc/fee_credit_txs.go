package fc

import (
	"bytes"
	"crypto"
	"errors"
	"fmt"
	"hash"

	"github.com/alphabill-org/alphabill/internal/txsystem"
	"github.com/alphabill-org/alphabill/internal/util"
	"github.com/holiman/uint256"
)

type (
	Wrapper struct {
		Transaction *txsystem.Transaction
		hashFunc    crypto.Hash
		hashValue   []byte
	}

	TransferFeeCreditWrapper struct {
		Wrapper
		TransferFC *TransferFeeCreditOrder
	}

	AddFeeCreditWrapper struct {
		Wrapper
		AddFC *AddFeeCreditOrder

		// TransferFC the "fee credit transfer" that also exist inside AddFeeCreditOrder as *txsystem.Transaction
		// needed to correctly serialize bytes
		TransferFC *TransferFeeCreditWrapper
	}

	CloseFeeCreditWrapper struct {
		Wrapper
		CloseFC *CloseFeeCreditOrder
	}

	ReclaimFeeCreditWrapper struct {
		Wrapper
		ReclaimFC *ReclaimFeeCreditOrder

		// CloseFCTransfer the "close fee credit" transfer that also exist inside ReclaimFeeCreditOrder as *txsystem.Transaction
		// needed to correctly serialize bytes
		CloseFCTransfer *CloseFeeCreditWrapper
	}
)

func NewFeeCreditTx(tx *txsystem.Transaction) (txsystem.GenericTransaction, error) {
	if tx == nil {
		return nil, errors.New("cannot create fee credit tx wrapper: tx is nil")
	}
	switch tx.TransactionAttributes.TypeUrl {
	case txsystem.TypeURLTransferFeeCreditOrder:
		pb := &TransferFeeCreditOrder{}
		err := tx.TransactionAttributes.UnmarshalTo(pb)
		if err != nil {
			return nil, err
		}
		return &TransferFeeCreditWrapper{
			Wrapper:    Wrapper{Transaction: tx},
			TransferFC: pb,
		}, nil
	case txsystem.TypeURLAddFeeCreditOrder:
		pb := &AddFeeCreditOrder{}
		err := tx.TransactionAttributes.UnmarshalTo(pb)
		if err != nil {
			return nil, err
		}
		fcGen, err := NewFeeCreditTx(pb.FeeCreditTransfer)
		if err != nil {
			return nil, fmt.Errorf("add fee credit wrapping failed: %w", err)
		}
		fcWrapper, ok := fcGen.(*TransferFeeCreditWrapper)
		if !ok {
			return nil, fmt.Errorf("transfer FC wrapper is invalid type: %T", fcWrapper)
		}
		return &AddFeeCreditWrapper{
			Wrapper:    Wrapper{Transaction: tx},
			AddFC:      pb,
			TransferFC: fcWrapper,
		}, nil
	case txsystem.TypeURLCloseFeeCreditOrder:
		pb := &CloseFeeCreditOrder{}
		err := tx.TransactionAttributes.UnmarshalTo(pb)
		if err != nil {
			return nil, err
		}
		return &CloseFeeCreditWrapper{
			Wrapper: Wrapper{Transaction: tx},
			CloseFC: pb,
		}, nil
	case txsystem.TypeURLReclaimFeeCreditOrder:
		pb := &ReclaimFeeCreditOrder{}
		err := tx.TransactionAttributes.UnmarshalTo(pb)
		if err != nil {
			return nil, err
		}
		fcGen, err := NewFeeCreditTx(pb.CloseFeeCreditTransfer)
		if err != nil {
			return nil, fmt.Errorf("reclaim fee credit wrapping failed: %w", err)
		}
		fcWrapper, ok := fcGen.(*CloseFeeCreditWrapper)
		if !ok {
			return nil, fmt.Errorf("close fee credit wrapper is invalid type: %T", fcWrapper)
		}
		return &ReclaimFeeCreditWrapper{
			Wrapper:         Wrapper{Transaction: tx},
			ReclaimFC:       pb,
			CloseFCTransfer: fcWrapper,
		}, nil
	default:
		return nil, nil
	}
}

// GenericTransaction methods
func (w *Wrapper) SystemID() []byte {
	return w.Transaction.SystemId
}
func (w *Wrapper) UnitID() *uint256.Int {
	return uint256.NewInt(0).SetBytes(w.Transaction.UnitId)
}
func (w *Wrapper) Timeout() uint64 {
	return w.Transaction.Timeout()
}
func (w *Wrapper) OwnerProof() []byte {
	return w.Transaction.OwnerProof
}
func (w *Wrapper) ToProtoBuf() *txsystem.Transaction {
	return w.Transaction
}
func (w *Wrapper) IsPrimary() bool {
	return true
}
func (w *Wrapper) TargetUnits(_ crypto.Hash) []*uint256.Int {
	return []*uint256.Int{w.UnitID()}
}
func (w *Wrapper) SetServerMetadata(sm *txsystem.ServerMetadata) {
	w.ToProtoBuf().ServerMetadata = sm
	w.resetHasher()
}
func (w *Wrapper) resetHasher() {
	w.hashValue = nil
}
func (w *Wrapper) transactionSigBytes(b *bytes.Buffer) {
	b.Write(w.Transaction.SystemId)
	b.Write(w.Transaction.UnitId)
	if w.Transaction.ClientMetadata != nil {
		b.Write(w.Transaction.ClientMetadata.Bytes())
	}
}
func (w *Wrapper) addTransactionFieldsToHasher(hasher hash.Hash) {
	hasher.Write(w.Transaction.SystemId)
	hasher.Write(w.Transaction.UnitId)
	hasher.Write(w.Transaction.OwnerProof)
	hasher.Write(w.Transaction.FeeProof)
	if w.Transaction.ClientMetadata != nil {
		hasher.Write(w.Transaction.ClientMetadata.Bytes())
	}
	if w.Transaction.ServerMetadata != nil {
		hasher.Write(w.Transaction.ServerMetadata.Bytes())
	}
}
func (w *Wrapper) hashComputed(hashFunc crypto.Hash) bool {
	return w.hashFunc == hashFunc && w.hashValue != nil
}

// GenericTransaction methods (transaction specific)
func (w *TransferFeeCreditWrapper) Hash(hashFunc crypto.Hash) []byte {
	if w.hashComputed(hashFunc) {
		return w.hashValue
	}
	hasher := hashFunc.New()
	w.AddToHasher(hasher)

	w.hashValue = hasher.Sum(nil)
	w.hashFunc = hashFunc
	return w.hashValue
}
func (w *TransferFeeCreditWrapper) AddToHasher(hasher hash.Hash) {
	w.Wrapper.addTransactionFieldsToHasher(hasher)
	w.TransferFC.addFieldsToHasher(hasher)
}
func (w *TransferFeeCreditWrapper) SigBytes() []byte {
	var b bytes.Buffer
	w.transactionSigBytes(&b)
	w.TransferFC.sigBytes(&b)
	return b.Bytes()
}

func (w *AddFeeCreditWrapper) Hash(hashFunc crypto.Hash) []byte {
	if w.hashComputed(hashFunc) {
		return w.hashValue
	}
	hasher := hashFunc.New()
	w.AddToHasher(hasher)

	w.hashValue = hasher.Sum(nil)
	w.hashFunc = hashFunc
	return w.hashValue
}
func (w *AddFeeCreditWrapper) AddToHasher(hasher hash.Hash) {
	w.Wrapper.addTransactionFieldsToHasher(hasher)
	w.addFieldsToHasher(hasher)
}
func (w *AddFeeCreditWrapper) SigBytes() []byte {
	var b bytes.Buffer
	w.transactionSigBytes(&b)
	w.sigBytes(&b)
	return b.Bytes()
}
func (w *AddFeeCreditWrapper) addFieldsToHasher(hasher hash.Hash) {
	hasher.Write(w.AddFC.FeeCreditOwnerCondition)
	w.TransferFC.AddToHasher(hasher)
	w.AddFC.FeeCreditTransferProof.AddToHasher(hasher)
}
func (w *AddFeeCreditWrapper) sigBytes(b *bytes.Buffer) {
	b.Write(w.AddFC.FeeCreditOwnerCondition)
	b.Write(w.TransferFC.SigBytes())
	b.Write(w.TransferFC.OwnerProof())
	b.Write(w.AddFC.FeeCreditTransferProof.Bytes())
}

func (w *CloseFeeCreditWrapper) Hash(hashFunc crypto.Hash) []byte {
	if w.hashComputed(hashFunc) {
		return w.hashValue
	}
	hasher := hashFunc.New()
	w.AddToHasher(hasher)

	w.hashValue = hasher.Sum(nil)
	w.hashFunc = hashFunc
	return w.hashValue
}
func (w *CloseFeeCreditWrapper) AddToHasher(hasher hash.Hash) {
	w.Wrapper.addTransactionFieldsToHasher(hasher)
	w.CloseFC.addFieldsToHasher(hasher)
}
func (w *CloseFeeCreditWrapper) SigBytes() []byte {
	var b bytes.Buffer
	w.transactionSigBytes(&b)
	w.CloseFC.sigBytes(&b)
	return b.Bytes()
}

func (w *ReclaimFeeCreditWrapper) Hash(hashFunc crypto.Hash) []byte {
	if w.hashComputed(hashFunc) {
		return w.hashValue
	}
	hasher := hashFunc.New()
	w.AddToHasher(hasher)

	w.hashValue = hasher.Sum(nil)
	w.hashFunc = hashFunc
	return w.hashValue
}
func (w *ReclaimFeeCreditWrapper) AddToHasher(hasher hash.Hash) {
	w.Wrapper.addTransactionFieldsToHasher(hasher)
	w.addFieldsToHasher(hasher)
}
func (w *ReclaimFeeCreditWrapper) SigBytes() []byte {
	var b bytes.Buffer
	w.transactionSigBytes(&b)
	w.sigBytes(&b)
	return b.Bytes()
}
func (w *ReclaimFeeCreditWrapper) addFieldsToHasher(hasher hash.Hash) {
	w.CloseFCTransfer.AddToHasher(hasher)
	w.ReclaimFC.CloseFeeCreditProof.AddToHasher(hasher)
	hasher.Write(w.ReclaimFC.Backlink)
}
func (w *ReclaimFeeCreditWrapper) sigBytes(b *bytes.Buffer) {
	b.Write(w.CloseFCTransfer.SigBytes())
	b.Write(w.CloseFCTransfer.OwnerProof())
	b.Write(w.ReclaimFC.CloseFeeCreditProof.Bytes())
	b.Write(w.ReclaimFC.Backlink)
}

// Protobuf transaction struct methods
func (x *TransferFeeCreditOrder) addFieldsToHasher(hasher hash.Hash) {
	hasher.Write(util.Uint64ToBytes(x.Amount))
	hasher.Write(x.TargetSystemIdentifier)
	hasher.Write(x.TargetRecordId)
	hasher.Write(util.Uint64ToBytes(x.EarliestAdditionTime))
	hasher.Write(util.Uint64ToBytes(x.LatestAdditionTime))
	hasher.Write(x.Nonce)
	hasher.Write(x.Backlink)
}
func (x *TransferFeeCreditOrder) sigBytes(b *bytes.Buffer) {
	b.Write(util.Uint64ToBytes(x.Amount))
	b.Write(x.TargetSystemIdentifier)
	b.Write(x.TargetRecordId)
	b.Write(util.Uint64ToBytes(x.EarliestAdditionTime))
	b.Write(util.Uint64ToBytes(x.LatestAdditionTime))
	b.Write(x.Nonce)
	b.Write(x.Backlink)
}

func (x *CloseFeeCreditOrder) addFieldsToHasher(hasher hash.Hash) {
	hasher.Write(util.Uint64ToBytes(x.Amount))
	hasher.Write(x.TargetUnitId)
	hasher.Write(x.Nonce)
}
func (x *CloseFeeCreditOrder) sigBytes(b *bytes.Buffer) {
	b.Write(util.Uint64ToBytes(x.Amount))
	b.Write(x.TargetUnitId)
	b.Write(x.Nonce)
}
