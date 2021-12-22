package wallet

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInMemoryWalletCanBeCreated(t *testing.T) {
	w, err := NewInMemoryWallet()
	require.NoError(t, err)
	require.EqualValues(t, 0, w.GetBalance())
}
