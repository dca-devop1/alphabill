package storage

import (
	"testing"

	"github.com/alphabill-org/alphabill/internal/keyvaluedb/memorydb"
	"github.com/alphabill-org/alphabill/internal/network/protocol/abdrc"
	abtypes "github.com/alphabill-org/alphabill/internal/rootchain/consensus/abdrc/types"
	"github.com/stretchr/testify/require"
)

func TestCertKey(t *testing.T) {
	// appends to prefix
	require.Equal(t, []byte("cert_test"), certKey([]byte("test")))
	require.Equal(t, []byte("cert_00000001"), certKey([]byte("00000001")))
	// nil case
	require.Equal(t, []byte(certPrefix), certKey(nil))
}

func TestBlockKey(t *testing.T) {
	const round uint64 = 1
	require.Equal(t, []byte("block_\000\000\000\000\000\000\000\001"), blockKey(round))
}

func TestWriteReadLastVote(t *testing.T) {
	t.Run("error - store proposal", func(t *testing.T) {
		db := memorydb.New()
		proposal := abdrc.ProposalMsg{}
		require.ErrorContains(t, WriteVote(db, proposal), "unknown vote type")
	})
	t.Run("read blank store", func(t *testing.T) {
		db := memorydb.New()
		msg, err := ReadVote(db)
		require.NoError(t, err)
		require.Nil(t, msg)
	})
	t.Run("ok - store vote", func(t *testing.T) {
		db := memorydb.New()
		vote := &abdrc.VoteMsg{Author: "test"}
		require.NoError(t, WriteVote(db, vote))
		// read back
		msg, err := ReadVote(db)
		require.NoError(t, err)
		require.IsType(t, &abdrc.VoteMsg{}, msg)
		require.Equal(t, "test", msg.(*abdrc.VoteMsg).Author)
	})
	t.Run("ok - store timeout vote", func(t *testing.T) {
		db := memorydb.New()
		vote := &abdrc.TimeoutMsg{Timeout: &abtypes.Timeout{Round: 1}, Author: "test"}
		require.NoError(t, WriteVote(db, vote))
		// read back
		msg, err := ReadVote(db)
		require.NoError(t, err)
		require.IsType(t, &abdrc.TimeoutMsg{}, msg)
		require.Equal(t, "test", msg.(*abdrc.TimeoutMsg).Author)
		require.EqualValues(t, 1, msg.(*abdrc.TimeoutMsg).Timeout.Round)
	})
}
