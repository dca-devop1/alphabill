package partition

import (
	gocrypto "crypto"
	"testing"
	"time"

	"github.com/alphabill-org/alphabill-go-base/types/hex"
	"github.com/btcsuite/btcd/btcutil/base58"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/stretchr/testify/require"

	"github.com/alphabill-org/alphabill-go-base/crypto"
	"github.com/alphabill-org/alphabill-go-base/hash"
	"github.com/alphabill-org/alphabill-go-base/types"
	testsig "github.com/alphabill-org/alphabill/internal/testutils/sig"
	"github.com/alphabill-org/alphabill/network/protocol/genesis"
	"github.com/alphabill-org/alphabill/state"
)

var zeroHash hex.Bytes = make([]byte, 32)
var nodeID peer.ID = "test"

func TestNewGenesisPartitionNode_NotOk(t *testing.T) {
	signer, verifier := testsig.CreateSignerAndVerifier(t)
	pubKeyBytes, err := verifier.MarshalPublicKey()
	require.NoError(t, err)
	validPDR := types.PartitionDescriptionRecord{
		Version:             1,
		NetworkIdentifier:   5,
		PartitionIdentifier: 1,
		TypeIdLen:           8,
		UnitIdLen:           128,
		T2Timeout:           5 * time.Second,
	}

	type args struct {
		state *state.State
		pdr   types.PartitionDescriptionRecord
		opts  []GenesisOption
	}

	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name:    "state is nil",
			args:    args{state: nil, pdr: validPDR},
			wantErr: ErrStateIsNil,
		},
		{
			name: "client signer is nil",
			args: args{
				state: state.NewEmptyState(),
				pdr:   validPDR,
				opts:  []GenesisOption{WithPeerID("1"), WithEncryptionPubKey(pubKeyBytes)},
			},
			wantErr: ErrSignerIsNil,
		},
		{
			name: "encryption public key is nil",
			args: args{
				state: state.NewEmptyState(),
				pdr:   validPDR,
				opts: []GenesisOption{
					WithSigningKey(signer),
					WithEncryptionPubKey(nil),
					WithPeerID("1")},
			},
			wantErr: ErrEncryptionPubKeyIsNil,
		},
		{
			name: "peer ID is empty",
			args: args{
				state: state.NewEmptyState(),
				pdr:   validPDR,
				opts: []GenesisOption{
					WithSigningKey(signer),
					WithPeerID(""),
				},
			},
			wantErr: genesis.ErrNodeIdentifierIsEmpty,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewNodeGenesis(tt.args.state, tt.args.pdr, tt.args.opts...)
			require.Nil(t, got)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}

	// invalid partition identifier
	got, err := NewNodeGenesis(
		state.NewEmptyState(),
		types.PartitionDescriptionRecord{Version: 1, NetworkIdentifier: 5, PartitionIdentifier: 0},
		WithPeerID("1"),
		WithSigningKey(signer),
		WithEncryptionPubKey(pubKeyBytes),
		WithHashAlgorithm(gocrypto.SHA256),
	)
	require.Nil(t, got)
	require.EqualError(t, err, `calculating genesis block hash: block hash calculation failed: invalid block: partition identifier is unassigned`)
}

func TestNewGenesisPartitionNode_Ok(t *testing.T) {
	signer, verifier := testsig.CreateSignerAndVerifier(t)
	pubKey, err := verifier.MarshalPublicKey()
	require.NoError(t, err)
	pdr := types.PartitionDescriptionRecord{Version: 1, NetworkIdentifier: 5, PartitionIdentifier: 1, T2Timeout: 2500 * time.Millisecond}
	pn := createPartitionNode(t, signer, verifier, pdr, nodeID)
	require.NotNil(t, pn)
	require.Equal(t, base58.Encode([]byte(nodeID)), pn.NodeIdentifier)
	require.Equal(t, hex.Bytes(pubKey), pn.SigningPublicKey)
	blockCertificationRequestRequest := pn.BlockCertificationRequest
	require.Equal(t, pdr.PartitionIdentifier, blockCertificationRequestRequest.Partition)
	require.NoError(t, blockCertificationRequestRequest.IsValid(verifier))

	ir := blockCertificationRequestRequest.InputRecord
	expectedHash := hex.Bytes(make([]byte, 32))
	require.Equal(t, expectedHash, ir.Hash)
	require.Equal(t, calculateBlockHash(pdr.PartitionIdentifier, nil, true), ir.BlockHash)
	require.Equal(t, zeroHash, ir.PreviousHash)
}

func createPartitionNode(t *testing.T, nodeSigningKey crypto.Signer, nodeEncryptionPublicKey crypto.Verifier, pdr types.PartitionDescriptionRecord, nodeIdentifier peer.ID) *genesis.PartitionNode {
	t.Helper()
	encPubKeyBytes, err := nodeEncryptionPublicKey.MarshalPublicKey()
	require.NoError(t, err)
	pn, err := NewNodeGenesis(
		state.NewEmptyState(),
		pdr,
		WithPeerID(nodeIdentifier),
		WithSigningKey(nodeSigningKey),
		WithEncryptionPubKey(encPubKeyBytes),
	)
	require.NoError(t, err)
	return pn
}

func calculateBlockHash(partitionIdentifier types.PartitionID, previousHash []byte, isEmpty bool) hex.Bytes {
	// blockhash = hash(header_hash, raw_txs_hash, mt_root_hash)
	hasher := gocrypto.SHA256.New()
	if isEmpty {
		return zeroHash
	}
	hasher.Write(partitionIdentifier.Bytes())
	hasher.Write(previousHash)
	headerHash := hasher.Sum(nil)

	hasher.Reset()
	txsHash := hasher.Sum(nil)

	treeHash := make([]byte, gocrypto.SHA256.Size())

	return hash.Sum(gocrypto.SHA256, headerHash, txsHash, treeHash)
}
