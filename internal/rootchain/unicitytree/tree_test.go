package unicitytree

import (
	"crypto"
	"crypto/sha256"
	"testing"

	"github.com/alphabill-org/alphabill/internal/smt"
	"google.golang.org/protobuf/proto"

	"github.com/alphabill-org/alphabill/internal/certificates"

	"github.com/stretchr/testify/require"
)

var inputRecord = &certificates.InputRecord{
	PreviousHash: []byte{0x00},
	Hash:         []byte{0x01},
	BlockHash:    []byte{0x02},
	SummaryValue: []byte{0x03},
}

func TestNewUnicityTree(t *testing.T) {
	unicityTree, err := New(sha256.New(), []*Data{
		{
			SystemIdentifier:            []byte{0, 0, 0, 1},
			InputRecord:                 inputRecord,
			SystemDescriptionRecordHash: []byte{1, 2, 3, 4},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, unicityTree)
}

func TestGetCertificate_Ok(t *testing.T) {
	key := []byte{0x00, 0x00, 0x00, 0x01}
	data := []*Data{
		{
			SystemIdentifier:            key,
			InputRecord:                 inputRecord,
			SystemDescriptionRecordHash: []byte{1, 2, 3, 4},
		},
	}
	unicityTree, err := New(sha256.New(), data)
	require.NoError(t, err)
	cert, err := unicityTree.GetCertificate(key)
	require.NoError(t, err)
	require.NotNil(t, cert)
	require.Equal(t, key, cert.SystemIdentifier)
	require.Equal(t, systemIdentifierLength*8, len(cert.SiblingHashes))

	hasher := crypto.SHA256.New()
	data[0].AddToHasher(hasher)
	dataHash := hasher.Sum(nil)
	hasher.Reset()

	root, err := smt.CalculatePathRoot(cert.SiblingHashes, dataHash, key, crypto.SHA256)
	require.Equal(t, unicityTree.GetRootHash(), root)
	ir, err := unicityTree.GetIR(key)
	require.NoError(t, err)
	require.True(t, proto.Equal(ir, data[0].InputRecord))
}

func TestGetCertificate_InvalidKey(t *testing.T) {
	unicityTree, err := New(sha256.New(), []*Data{
		{
			SystemIdentifier:            []byte{1, 2, 3, 1},
			InputRecord:                 inputRecord,
			SystemDescriptionRecordHash: []byte{1, 2, 3, 4},
		},
	})
	require.NoError(t, err)
	cert, err := unicityTree.GetCertificate([]byte{0x00, 0x00})

	require.Nil(t, cert)
	require.ErrorIs(t, err, ErrInvalidSystemIdentifierLength)
}

func TestGetCertificate_KeyNotFound(t *testing.T) {
	unicityTree, err := New(sha256.New(), []*Data{
		{
			SystemIdentifier:            []byte{1, 2, 3, 1},
			InputRecord:                 inputRecord,
			SystemDescriptionRecordHash: []byte{1, 2, 3, 4},
		},
	})
	require.NoError(t, err)
	cert, err := unicityTree.GetCertificate([]byte{0, 0, 0, 0})
	require.Nil(t, cert)
	require.ErrorContains(t, err, "certificate for system id 00000000 not found")
	ir, err := unicityTree.GetIR([]byte{0, 0, 0, 0})
	require.ErrorContains(t, err, "ir for system id 00000000 not found")
	require.Nil(t, ir)
}
