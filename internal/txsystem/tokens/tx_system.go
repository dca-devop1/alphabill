package tokens

import (
	"crypto"
	"reflect"

	"github.com/alphabill-org/alphabill/internal/errors"
	"github.com/alphabill-org/alphabill/internal/rma"
	"github.com/alphabill-org/alphabill/internal/txsystem"
)

const (
	zeroSummaryValue = rma.Uint64SummaryValue(0)
	uriMaxSize       = 4 * 1024
	dataMaxSize      = 64 * 1024
	maxSymbolLength  = 64
	maxDecimalPlaces = 8

	ErrStrSystemIdentifierIsNil = "system identifier is nil"
	ErrStrUnitIDIsZero          = "unit ID cannot be zero"
	ErrStrInvalidSymbolName     = "symbol name exceeds the allowed maximum length of 64 bytes"
)

type (
	tokensTxSystem struct {
		systemIdentifier   []byte
		state              *rma.Tree
		hashAlgorithm      crypto.Hash
		currentBlockNumber uint64
		executors          map[reflect.Type]txExecutor
	}

	txExecutor interface {
		Execute(tx txsystem.GenericTransaction, currentBlockNr uint64) error
	}
)

func New(opts ...Option) (*tokensTxSystem, error) {
	options := defaultOptions()
	for _, opt := range opts {
		opt(options)
	}
	if options.systemIdentifier == nil {
		return nil, errors.New(ErrStrSystemIdentifierIsNil)
	}
	state, err := rma.New(&rma.Config{
		HashAlgorithm: options.hashAlgorithm,
	})
	if err != nil {
		return nil, err
	}

	txs := &tokensTxSystem{
		systemIdentifier: options.systemIdentifier,
		hashAlgorithm:    options.hashAlgorithm,
		state:            state,
		executors:        initExecutors(state, options),
	}
	logger.Info("TokensTransactionSystem initialized: systemIdentifier=%X, hashAlgorithm=%v", options.systemIdentifier, options.hashAlgorithm)
	return txs, nil
}

func (t *tokensTxSystem) State() (txsystem.State, error) {
	if t.state.ContainsUncommittedChanges() {
		return nil, txsystem.ErrStateContainsUncommittedChanges
	}
	return t.getState(), nil
}

func (t *tokensTxSystem) ConvertTx(tx *txsystem.Transaction) (txsystem.GenericTransaction, error) {
	return NewGenericTx(tx)
}

func (t *tokensTxSystem) Execute(tx txsystem.GenericTransaction) error {
	err := txsystem.ValidateGenericTransaction(&txsystem.TxValidationContext{Tx: tx, Bd: nil, SystemIdentifier: t.systemIdentifier, BlockNumber: t.currentBlockNumber})
	if err != nil {
		return err
	}
	txType := reflect.TypeOf(tx)
	executor := t.executors[txType]
	if executor == nil {
		return errors.Errorf("unknown tx type %T", tx)
	}
	return executor.Execute(tx, t.currentBlockNumber)
}

func (t *tokensTxSystem) BeginBlock(blockNr uint64) {
	t.currentBlockNumber = blockNr
}

func (t *tokensTxSystem) EndBlock() (txsystem.State, error) {
	return t.getState(), nil
}

func (t *tokensTxSystem) Revert() {
	t.state.Revert()
}

func (t *tokensTxSystem) Commit() {
	t.state.Commit()
}

func (t *tokensTxSystem) getState() txsystem.State {
	if t.state.GetRootHash() == nil {
		return txsystem.NewStateSummary(make([]byte, t.hashAlgorithm.Size()), zeroSummaryValue.Bytes())
	}
	return txsystem.NewStateSummary(t.state.GetRootHash(), zeroSummaryValue.Bytes())
}

func initExecutors(state *rma.Tree, options *Options) map[reflect.Type]txExecutor {
	executors := make(map[reflect.Type]txExecutor)
	// non-fungible token tx executors
	commonNFTTxExecutor := &baseTxExecutor[*nonFungibleTokenTypeData]{
		state:         state,
		hashAlgorithm: options.hashAlgorithm,
	}
	executors[reflect.TypeOf(&createNonFungibleTokenTypeWrapper{})] = &createNonFungibleTokenTypeTxExecutor{commonNFTTxExecutor}
	executors[reflect.TypeOf(&mintNonFungibleTokenWrapper{})] = &mintNonFungibleTokenTxExecutor{commonNFTTxExecutor}
	executors[reflect.TypeOf(&transferNonFungibleTokenWrapper{})] = &transferNonFungibleTokenTxExecutor{commonNFTTxExecutor}
	executors[reflect.TypeOf(&updateNonFungibleTokenWrapper{})] = &updateNonFungibleTokenTxExecutor{commonNFTTxExecutor}

	// fungible token tx executors
	commonFungibleTokenTxExecutor := &baseTxExecutor[*fungibleTokenTypeData]{
		state:         state,
		hashAlgorithm: options.hashAlgorithm,
	}
	executors[reflect.TypeOf(&createFungibleTokenTypeWrapper{})] = &createFungibleTokenTypeTxExecutor{commonFungibleTokenTxExecutor}
	executors[reflect.TypeOf(&mintFungibleTokenWrapper{})] = &mintFungibleTokenTxExecutor{commonFungibleTokenTxExecutor}
	executors[reflect.TypeOf(&transferFungibleTokenWrapper{})] = &transferFungibleTokenTxExecutor{commonFungibleTokenTxExecutor}
	executors[reflect.TypeOf(&splitFungibleTokenWrapper{})] = &splitFungibleTokenTxExecutor{commonFungibleTokenTxExecutor}
	executors[reflect.TypeOf(&burnFungibleTokenWrapper{})] = &burnFungibleTokenTxExecutor{commonFungibleTokenTxExecutor}

	return executors
}