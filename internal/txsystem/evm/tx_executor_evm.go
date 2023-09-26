package evm

import (
	"fmt"
	"math/big"
	"os"

	"github.com/alphabill-org/alphabill/internal/keyvaluedb"
	"github.com/alphabill-org/alphabill/internal/txsystem"
	"github.com/alphabill-org/alphabill/internal/txsystem/evm/statedb"
	"github.com/alphabill-org/alphabill/internal/types"
	"github.com/alphabill-org/alphabill/internal/util"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/vm"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/eth/tracers/logger"
	"github.com/fxamacker/cbor/v2"
)

type (
	StateTransition struct {
		gp         *core.GasPool
		msg        *TxAttributes
		gas        uint64
		gasPrice   *big.Int
		initialGas uint64
		value      *big.Int
		data       []byte
		state      vm.StateDB
		evm        *vm.EVM
	}

	ProcessingDetails struct {
		_            struct{} `cbor:",toarray"`
		ErrorDetails string
		ReturnData   []byte
		ContractAddr common.Address
		Logs         []*statedb.LogEntry
	}
)

func errorToStr(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}
func (d *ProcessingDetails) Bytes() ([]byte, error) {
	return cbor.Marshal(d)
}

func handleEVMTx(systemIdentifier []byte, opts *Options, blockGas *core.GasPool, blockDB keyvaluedb.KeyValueDB) txsystem.GenericExecuteFunc[TxAttributes] {
	return func(tx *types.TransactionOrder, attr *TxAttributes, currentBlockNumber uint64) (sm *types.ServerMetadata, err error) {
		from := common.BytesToAddress(attr.From)
		stateDB := statedb.NewStateDB(opts.state)
		if !stateDB.Exist(from) {
			return nil, fmt.Errorf(" address %v does not exist", from)
		}
		defer func() {
			if err == nil {
				err = stateDB.Finalize()
			}
		}()
		return Execute(currentBlockNumber, stateDB, blockDB, attr, systemIdentifier, blockGas, opts.gasUnitPrice, false)
	}
}

func calcGasPrice(gas uint64, gasPrice *big.Int) *big.Int {
	cost := new(big.Int).SetUint64(gas)
	return cost.Mul(cost, gasPrice)
}

func Execute(currentBlockNumber uint64, stateDB *statedb.StateDB, blockDB keyvaluedb.KeyValueDB, attr *TxAttributes, systemIdentifier []byte, gp *core.GasPool, gasUnitPrice *big.Int, fake bool) (*types.ServerMetadata, error) {
	if err := validate(attr); err != nil {
		return nil, err
	}
	blockCtx := newBlockContext(currentBlockNumber, blockDB)
	evm := vm.NewEVM(blockCtx, newTxContext(attr, gasUnitPrice), stateDB, newChainConfig(new(big.Int).SetBytes(systemIdentifier)), newVMConfig())
	msg := attr.AsMessage(gasUnitPrice, fake)
	// Apply the transaction to the current state (included in the env)
	execResult, err := core.ApplyMessage(evm, msg, gp)
	if err != nil {
		return nil, err
	}
	success := types.TxStatusSuccessful
	var errorDetail error
	if execResult.Unwrap() != nil || stateDB.DBError() != nil {
		success = types.TxStatusFailed
		if execResult.Unwrap() != nil {
			errorDetail = fmt.Errorf("evm runtime error: %w", execResult.Unwrap())
		}
		if stateDB.DBError() != nil {
			errorDetail = fmt.Errorf("%w state db error: %w", errorDetail, stateDB.DBError())
		}
	}
	// The contract address can be derived from the transaction itself
	var contractAddress common.Address
	if attr.ToAddr() == nil {
		// Deriving the signer is expensive, only do if it's actually needed
		contractAddress = ethcrypto.CreateAddress(attr.FromAddr(), attr.Nonce)
	}
	evmProcessingDetails := &ProcessingDetails{
		ReturnData:   execResult.ReturnData,
		ContractAddr: contractAddress,
		ErrorDetails: errorToStr(errorDetail),
	}
	if errorDetail == nil {
		evmProcessingDetails.Logs = stateDB.GetLogs()
	}
	detailBytes, err := evmProcessingDetails.Bytes()
	if err != nil {
		return nil, fmt.Errorf("evm result encode error %w", err)
	}
	txPrice := calcGasPrice(execResult.UsedGas, gasUnitPrice)
	// calculate gas based fee in alpha and refund the remainder
	fee, remainderWei := weiToAlphaWithReminder(txPrice)
	// HACK: AB-1207 - quick hack for first evm release, refund remainder back to the account
	// Todo: Create a proper solution and implement ApplyMessage in this project
	stateDB.AddBalance(msg.From, remainderWei)

	log.Trace("total gas: %v gas units, price in alpha %v", execResult.UsedGas, fee)
	return &types.ServerMetadata{ActualFee: fee.Uint64(), TargetUnits: stateDB.GetUpdatedUnits(), SuccessIndicator: success, ProcessingDetails: detailBytes}, nil
}

func newBlockContext(currentBlockNumber uint64, blockDB keyvaluedb.KeyValueDB) vm.BlockContext {
	return vm.BlockContext{
		CanTransfer: core.CanTransfer,
		Transfer:    core.Transfer,
		GetHash: func(u uint64) common.Hash {
			// NB! SIGSEGV if blockDB is nil, this must not happen
			it := blockDB.Find(util.Uint64ToBytes(u))
			if !it.Valid() {
				return common.Hash{}
			}
			b := &types.Block{}
			if err := it.Value(b); err != nil {
				return common.Hash{}
			}
			return common.BytesToHash(b.UnicityCertificate.InputRecord.BlockHash)
		},
		Coinbase:      common.Address{},
		GasLimit:      DefaultBlockGasLimit,
		BlockNumber:   new(big.Int).SetUint64(currentBlockNumber),
		Time:          1,
		Difficulty:    big.NewInt(0),
		BaseFee:       big.NewInt(0),
		Random:        nil,
		ExcessBlobGas: nil,
	}
}

func newTxContext(attr *TxAttributes, gasPrice *big.Int) vm.TxContext {
	return vm.TxContext{
		Origin:   common.BytesToAddress(attr.From),
		GasPrice: gasPrice,
	}
}

func newVMConfig() vm.Config {
	return vm.Config{
		// TODO use AB logger
		Tracer:                  logger.NewJSONLogger(nil, os.Stdout),
		NoBaseFee:               true,
		EnablePreimageRecording: false, // Enables recording of SHA3/keccak preimages
	}
}

// validate - validate EVM call attributes
func validate(attr *TxAttributes) error {
	if attr.From == nil {
		return fmt.Errorf("invalid evm tx, from addr is nil")
	}
	if attr.Value == nil {
		return fmt.Errorf("invalid evm tx, value is nil")
	}
	if attr.Value.Sign() < 0 {
		return fmt.Errorf("invalid evm tx, value is negative")
	}
	return nil
}
