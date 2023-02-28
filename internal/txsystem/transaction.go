package txsystem

import (
	"bytes"

	"github.com/alphabill-org/alphabill/internal/util"
	"github.com/holiman/uint256"
)

// Bytes serializes the generic transaction order fields.
func (x *Transaction) Bytes() []byte {
	var b bytes.Buffer
	b.Write(x.SystemId)
	b.Write(x.UnitId)
	b.Write(x.OwnerProof)
	b.Write(x.FeeProof)
	if x.ClientMetadata != nil {
		b.Write(x.ClientMetadata.Bytes())
	}
	if x.ServerMetadata != nil {
		b.Write(x.ServerMetadata.Bytes())
	}
	return b.Bytes()
}

func (x *Transaction) GetClientFeeCreditRecordID() *uint256.Int {
	clientFeeCreditID := x.ClientMetadata.FeeCreditRecordId
	return uint256.NewInt(0).SetBytes(clientFeeCreditID)
}

func (x *ClientMetadata) Bytes() []byte {
	var b bytes.Buffer
	b.Write(util.Uint64ToBytes(x.Timeout))
	b.Write(util.Uint64ToBytes(x.MaxFee))
	b.Write(x.FeeCreditRecordId)
	return b.Bytes()
}

func (x *ServerMetadata) Bytes() []byte {
	var b bytes.Buffer
	b.Write(util.Uint64ToBytes(x.Fee))
	return b.Bytes()
}

// TxBytes returns Bytes without ServerMetadata field.
func (x *Transaction) TxBytes() []byte {
	var b bytes.Buffer
	b.Write(x.SystemId)
	b.Write(x.UnitId)
	b.Write(x.OwnerProof)
	b.Write(x.FeeProof)
	if x.ClientMetadata != nil {
		b.Write(x.ClientMetadata.Bytes())
	}
	return b.Bytes()
}

// Timeout returns timeout from ClientMetadata, defaults to 0 if client metadata is nil.
func (x *Transaction) Timeout() uint64 {
	cm := x.ClientMetadata
	if cm != nil {
		return cm.Timeout
	}
	return 0
}
