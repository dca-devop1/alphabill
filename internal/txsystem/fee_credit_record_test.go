package txsystem

import (
	"crypto"

	"testing"

	"github.com/alphabill-org/alphabill/internal/rma"
	test "github.com/alphabill-org/alphabill/internal/testutils"
	"github.com/alphabill-org/alphabill/internal/util"
	"github.com/stretchr/testify/require"
)

func TestFCR_HashIsCalculatedCorrectly(t *testing.T) {
	fcr := &FeeCreditRecord{
		balance: 1,
		hash:    test.RandomBytes(32),
		timeout: 2,
	}
	// calculate actual hash
	hasher := crypto.SHA256.New()
	fcr.AddToHasher(hasher)
	actualHash := hasher.Sum(nil)

	// calculate expected hash
	hasher.Reset()
	hasher.Write(util.Uint64ToBytes(uint64(fcr.balance)))
	hasher.Write(fcr.hash)
	hasher.Write(util.Uint64ToBytes(fcr.timeout))
	expectedHash := hasher.Sum(nil)

	require.Equal(t, expectedHash, actualHash)
}

func TestFCR_SummaryValueIsZero(t *testing.T) {
	fcr := &FeeCreditRecord{
		balance: 1,
		hash:    test.RandomBytes(32),
		timeout: 2,
	}
	require.Equal(t, rma.Uint64SummaryValue(0), fcr.Value())
}
