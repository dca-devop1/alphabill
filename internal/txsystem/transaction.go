package txsystem

import (
	"bytes"

	"gitdc.ee.guardtime.com/alphabill/alphabill/internal/util"
)

// Bytes serializes the generic transaction order fields.
func (x *Transaction) Bytes() []byte {
	var b bytes.Buffer
	b.Write(x.SystemId)
	b.Write(x.UnitId)
	b.Write(util.Uint64ToBytes(x.Timeout))
	b.Write(x.OwnerProof)
	return b.Bytes()
}
