package vd

import (
	gocrypto "crypto"

	"github.com/alphabill-org/alphabill/internal/crypto"
	"github.com/alphabill-org/alphabill/internal/state"
	"github.com/alphabill-org/alphabill/internal/txsystem/fc"
	"github.com/alphabill-org/alphabill/internal/txsystem/money"
)

var DefaultSystemIdentifier = []byte{0, 0, 0, 1}

type (
	Options struct {
		systemIdentifier      []byte
		moneySystemIdentifier []byte
		hashAlgorithm         gocrypto.Hash
		trustBase             map[string]crypto.Verifier
		state                 *state.State
		feeCalculator         fc.FeeCalculator
	}

	Option func(*Options)
)

func defaultOptions() (*Options, error) {
	return &Options{
		systemIdentifier:      DefaultSystemIdentifier,
		moneySystemIdentifier: money.DefaultSystemIdentifier,
		hashAlgorithm:         gocrypto.SHA256,
		trustBase:             map[string]crypto.Verifier{},
		state:                 state.NewEmptyState(),
		feeCalculator:         fc.FixedFee(1),
	}, nil
}

func WithState(s *state.State) Option {
	return func(c *Options) {
		c.state = s
	}
}

func WithSystemIdentifier(systemIdentifier []byte) Option {
	return func(c *Options) {
		c.systemIdentifier = systemIdentifier
	}
}

func WithMoneySystemIdentifier(moneySystemIdentifier []byte) Option {
	return func(c *Options) {
		c.moneySystemIdentifier = moneySystemIdentifier
	}
}

func WithFeeCalculator(calc fc.FeeCalculator) Option {
	return func(g *Options) {
		g.feeCalculator = calc
	}
}

func WithHashAlgorithm(algorithm gocrypto.Hash) Option {
	return func(c *Options) {
		c.hashAlgorithm = algorithm
	}
}

func WithTrustBase(trustBase map[string]crypto.Verifier) Option {
	return func(c *Options) {
		c.trustBase = trustBase
	}
}
