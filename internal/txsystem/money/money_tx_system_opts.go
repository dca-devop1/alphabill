package money

import (
	gocrypto "crypto"

	"github.com/alphabill-org/alphabill/internal/crypto"
	"github.com/alphabill-org/alphabill/internal/network/protocol/genesis"
	"github.com/alphabill-org/alphabill/internal/state"
	"github.com/alphabill-org/alphabill/internal/txsystem/fc"
)

var DefaultSystemIdentifier = []byte{0, 0, 0, 0}

type (
	Options struct {
		systemIdentifier         []byte
		state                    *state.State
		hashAlgorithm            gocrypto.Hash
		trustBase                map[string]crypto.Verifier
		initialBill              *InitialBill
		dcMoneyAmount            uint64
		dustBillDeletionTimeout  uint64
		systemDescriptionRecords []*genesis.SystemDescriptionRecord
		feeCalculator            fc.FeeCalculator
	}

	Option func(*Options)
)

func DefaultOptions() *Options {
	return &Options{
		systemIdentifier:        DefaultSystemIdentifier,
		hashAlgorithm:           gocrypto.SHA256,
		state:                   state.NewEmptyState(),
		trustBase:               make(map[string]crypto.Verifier),
		dcMoneyAmount:           0,
		dustBillDeletionTimeout: defaultDustBillDeletionTimeout,
		feeCalculator:           fc.FixedFee(1),
	}
}

func WithSystemIdentifier(systemIdentifier []byte) Option {
	return func(g *Options) {
		g.systemIdentifier = systemIdentifier
	}
}

func WithState(s *state.State) Option {
	return func(g *Options) {
		g.state = s
	}
}

func WithTrustBase(trust map[string]crypto.Verifier) Option {
	return func(options *Options) {
		options.trustBase = trust
	}
}

func WithHashAlgorithm(hashAlgorithm gocrypto.Hash) Option {
	return func(g *Options) {
		g.hashAlgorithm = hashAlgorithm
	}
}

func WithInitialBill(bill *InitialBill) Option {
	return func(g *Options) {
		g.initialBill = bill
	}
}

func WithDCMoneyAmount(a uint64) Option {
	return func(g *Options) {
		g.dcMoneyAmount = a
	}
}

func WithDCBillDeletionTimeout(t uint64) Option {
	return func(g *Options) {
		g.dustBillDeletionTimeout = t
	}
}

func WithSystemDescriptionRecords(records []*genesis.SystemDescriptionRecord) Option {
	return func(g *Options) {
		g.systemDescriptionRecords = records
	}
}

func WithFeeCalculator(calc fc.FeeCalculator) Option {
	return func(g *Options) {
		g.feeCalculator = calc
	}
}
