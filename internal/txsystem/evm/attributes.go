package evm

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

const PayloadTypeEVMCall = "evm"

type TxAttributes struct {
	_        struct{} `cbor:",toarray"`
	From     []byte
	To       []byte
	Data     []byte
	Value    *big.Int
	Gas      uint64
	GasPrice *big.Int
}

// GetFromAddr - returns From as Address, if nil empty address is returned
// From is mandatory field and must not be nil in a valid TxAttributes
func (t *TxAttributes) GetFromAddr() common.Address {
	if t == nil || t.From == nil {
		return common.Address{}
	}
	return common.BytesToAddress(t.From)
}

// GetToAddr - returns To as Address pointer.
// To field is optional and not present on contract creation calls.
// If To is nil then nil pointer is returned
func (t *TxAttributes) GetToAddr() *common.Address {
	if t == nil || t.To == nil {
		return nil
	}
	addr := common.BytesToAddress(t.To)
	return &addr
}
