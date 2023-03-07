package sc

import (
	"errors"
	"fmt"

	"github.com/alphabill-org/alphabill/internal/txsystem"
)

func New(systemIdentifier []byte, opts ...Option) (*txsystem.GenericTxSystem, error) {
	options := defaultOptions()
	for _, opt := range opts {
		opt(options)
	}
	if systemIdentifier == nil {
		return nil, errors.New("system identifier is nil")
	}
	sc, err := NewSmartContractModule(systemIdentifier, options)
	if err != nil {
		return nil, fmt.Errorf("failed to load smart contract module: %w", err)
	}
	return txsystem.NewGenericTxSystem(
		[]txsystem.Module{sc},
		txsystem.WithSystemIdentifier(systemIdentifier),
		txsystem.WithHashAlgorithm(options.hashAlgorithm),
		txsystem.WithState(options.state),
	)
}
