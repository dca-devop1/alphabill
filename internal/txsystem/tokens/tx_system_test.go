package tokens

import (
	gocrypto "crypto"
	"fmt"
	"hash"
	"testing"

	"github.com/alphabill-org/alphabill/internal/crypto"
	hasher "github.com/alphabill-org/alphabill/internal/hash"
	"github.com/alphabill-org/alphabill/internal/rma"
	"github.com/alphabill-org/alphabill/internal/script"
	testtransaction "github.com/alphabill-org/alphabill/internal/testutils/transaction"
	"github.com/alphabill-org/alphabill/internal/txsystem"
	"github.com/holiman/uint256"
	"github.com/stretchr/testify/require"
)

var (
	parent1Identifier = uint256.NewInt(1)
	parent2Identifier = uint256.NewInt(2)
	unitIdentifier    = uint256.NewInt(10)
)

func TestNewTokenTxSystem_DefaultOptions(t *testing.T) {
	txs, err := New()
	require.NoError(t, err)
	require.Equal(t, gocrypto.SHA256, txs.hashAlgorithm)
	require.Equal(t, DefaultTokenTxSystemIdentifier, txs.systemIdentifier)

	require.NotNil(t, txs.state)
	state, err := txs.State()
	require.NoError(t, err)
	require.Equal(t, make([]byte, gocrypto.SHA256.Size()), state.Root())
	require.Equal(t, zeroSummaryValue.Bytes(), state.Summary())
}

func TestNewTokenTxSystem_NilSystemIdentifier(t *testing.T) {
	txs, err := New(WithSystemIdentifier(nil))
	require.ErrorContains(t, err, ErrStrSystemIdentifierIsNil)
	require.Nil(t, txs)
}

func TestNewTokenTxSystem_UnsupportedHashAlgorithm(t *testing.T) {
	txs, err := New(WithHashAlgorithm(gocrypto.SHA1))
	require.ErrorContains(t, err, "invalid hash algorithm")
	require.Nil(t, txs)
}

func TestNewTokenTxSystem_OverrideDefaultOptions(t *testing.T) {
	systemIdentifier := []byte{0, 0, 0, 7}
	txs, err := New(WithSystemIdentifier(systemIdentifier), WithHashAlgorithm(gocrypto.SHA512))
	require.NoError(t, err)
	require.Equal(t, gocrypto.SHA512, txs.hashAlgorithm)
	require.Equal(t, systemIdentifier, txs.systemIdentifier)

	require.NotNil(t, txs.state)
	state, err := txs.State()
	require.NoError(t, err)
	require.Equal(t, make([]byte, gocrypto.SHA512.Size()), state.Root())
	require.Equal(t, zeroSummaryValue.Bytes(), state.Summary())
}

func TestExecuteCreateNFTType_WithoutParentID(t *testing.T) {
	txs := newTokenTxSystem(t)
	tx := testtransaction.NewGenericTransaction(
		t,
		txs.ConvertTx,
		testtransaction.WithUnitId(unitIdentifier.Bytes()),
		testtransaction.WithSystemID(DefaultTokenTxSystemIdentifier),
		testtransaction.WithAttributes(&CreateNonFungibleTokenTypeAttributes{
			Symbol:                   symbol,
			SubTypeCreationPredicate: subTypeCreationPredicate,
			TokenCreationPredicate:   tokenCreationPredicate,
			InvariantPredicate:       invariantPredicate,
			DataUpdatePredicate:      dataUpdatePredicate,
		}),
	)

	require.NoError(t, txs.Execute(tx))
	u, err := txs.state.GetUnit(unitIdentifier)
	require.NoError(t, err)
	require.Equal(t, tx.Hash(gocrypto.SHA256), u.StateHash)
	require.IsType(t, &nonFungibleTokenTypeData{}, u.Data)
	data := u.Data.(*nonFungibleTokenTypeData)
	require.Equal(t, zeroSummaryValue, data.Value())
	require.Equal(t, symbol, data.symbol)
	require.Equal(t, uint256.NewInt(0), data.parentTypeId)
	require.Equal(t, subTypeCreationPredicate, data.subTypeCreationPredicate)
	require.Equal(t, tokenCreationPredicate, data.tokenCreationPredicate)
	require.Equal(t, invariantPredicate, data.invariantPredicate)
	require.Equal(t, dataUpdatePredicate, data.dataUpdatePredicate)
}

func TestExecuteCreateNFTType_WithParentID(t *testing.T) {
	txs := newTokenTxSystem(t)
	createParentTx := testtransaction.NewGenericTransaction(
		t,
		txs.ConvertTx,
		testtransaction.WithUnitId(parent1Identifier.Bytes()),
		testtransaction.WithSystemID(DefaultTokenTxSystemIdentifier),
		testtransaction.WithAttributes(&CreateNonFungibleTokenTypeAttributes{
			Symbol:                   symbol,
			SubTypeCreationPredicate: script.PredicateAlwaysTrue(),
		}),
	)

	require.NoError(t, txs.Execute(createParentTx))
	tx := testtransaction.NewGenericTransaction(
		t,
		txs.ConvertTx,
		testtransaction.WithUnitId(unitIdentifier.Bytes()),
		testtransaction.WithSystemID(DefaultTokenTxSystemIdentifier),
		testtransaction.WithOwnerProof(script.PredicateArgumentEmpty()),
		testtransaction.WithAttributes(
			&CreateNonFungibleTokenTypeAttributes{
				Symbol:                   symbol,
				ParentTypeId:             parent1Identifier.Bytes(),
				SubTypeCreationPredicate: script.PredicateAlwaysFalse(),
			},
		),
	)
	require.NoError(t, txs.Execute(tx))
}

func TestExecuteCreateNFTType_InheritanceChainWithP2PKHPredicates(t *testing.T) {
	// Inheritance Chain: parent1Identifier <- parent2Identifier <- unitIdentifier
	parent2Signer, parent2PubKey := createSigner(t)
	childSigner, childPublicKey := createSigner(t)

	// only parent2 can create sub-types from parent1
	parent1SubTypeCreationPredicate := script.PredicatePayToPublicKeyHashDefault(hasher.Sum256(parent2PubKey))

	// parent2 and child together can create a sub-type because SubTypeCreationPredicate are concatenated (ownerProof must contain both signatures)
	parent2SubTypeCreationPredicate := script.PredicatePayToPublicKeyHashDefault(hasher.Sum256(childPublicKey))
	parent2SubTypeCreationPredicate[0] = script.OpVerify // verify parent1SubTypeCreationPredicate signature verification result (replace script.StartByte byte with OpVerify)

	txs := newTokenTxSystem(t)

	// create parent1 type
	createParent1Tx := testtransaction.NewGenericTransaction(
		t,
		txs.ConvertTx,
		testtransaction.WithUnitId(parent1Identifier.Bytes()),
		testtransaction.WithSystemID(DefaultTokenTxSystemIdentifier),
		testtransaction.WithAttributes(&CreateNonFungibleTokenTypeAttributes{
			Symbol:                   symbol,
			SubTypeCreationPredicate: parent1SubTypeCreationPredicate,
		}),
	)
	require.NoError(t, txs.Execute(createParent1Tx))

	// create parent2 type
	createParent2Tx := testtransaction.NewTransaction(
		t,
		testtransaction.WithUnitId(parent2Identifier.Bytes()),
		testtransaction.WithSystemID(DefaultTokenTxSystemIdentifier),
		testtransaction.WithAttributes(
			&CreateNonFungibleTokenTypeAttributes{
				Symbol:                   symbol,
				ParentTypeId:             parent1Identifier.Bytes(),
				SubTypeCreationPredicate: parent2SubTypeCreationPredicate,
			},
		),
	)
	gtx, signature := signTx(t, txs, createParent2Tx, parent2Signer, parent2PubKey)
	require.NoError(t, txs.Execute(gtx))

	// create child sub-type
	createChildTx := testtransaction.NewTransaction(
		t,
		testtransaction.WithUnitId(unitIdentifier.Bytes()),
		testtransaction.WithSystemID(DefaultTokenTxSystemIdentifier),
		testtransaction.WithAttributes(
			&CreateNonFungibleTokenTypeAttributes{
				Symbol:                   symbol,
				ParentTypeId:             parent2Identifier.Bytes(),
				SubTypeCreationPredicate: script.PredicateAlwaysFalse(), // no sub-types
			},
		),
	)
	gtx, err := txs.ConvertTx(createChildTx)
	require.NoError(t, err)

	signature, err = childSigner.SignBytes(gtx.SigBytes())
	require.NoError(t, err)
	signature2, err := parent2Signer.SignBytes(gtx.SigBytes())
	require.NoError(t, err)

	// child owner proof must satisfy parent1 & parent2 SubTypeCreationPredicates
	createChildTx.OwnerProof = append(
		script.PredicateArgumentPayToPublicKeyHashDefault(signature, childPublicKey),        // parent2 predicate argument (with script.StartByte byte)
		script.PredicateArgumentPayToPublicKeyHashDefault(signature2, parent2PubKey)[1:]..., // parent1 predicate argument (without script.StartByte byte)
	)
	gtx, err = txs.ConvertTx(createChildTx)
	require.NoError(t, err)
	require.NoError(t, txs.Execute(gtx))
}

func TestExecuteCreateNFTType_UnitTypeIsZero(t *testing.T) {
	txs := newTokenTxSystem(t)
	tx := testtransaction.NewGenericTransaction(
		t,
		txs.ConvertTx,
		testtransaction.WithUnitId(uint256.NewInt(0).Bytes()),
		testtransaction.WithSystemID(DefaultTokenTxSystemIdentifier),
		testtransaction.WithAttributes(&CreateNonFungibleTokenTypeAttributes{}),
	)
	require.ErrorContains(t, txs.Execute(tx), ErrStrUnitIDIsZero)
}

func TestExecuteCreateNFTType_UnitIDExists(t *testing.T) {
	txs := newTokenTxSystem(t)
	tx := testtransaction.NewGenericTransaction(
		t,
		txs.ConvertTx,
		testtransaction.WithUnitId(unitIdentifier.Bytes()),
		testtransaction.WithSystemID(DefaultTokenTxSystemIdentifier),
		testtransaction.WithAttributes(&CreateNonFungibleTokenTypeAttributes{
			Symbol:                   symbol,
			SubTypeCreationPredicate: subTypeCreationPredicate,
		}),
	)
	require.NoError(t, txs.Execute(tx))
	require.ErrorContains(t, txs.Execute(tx), fmt.Sprintf("unit %v exists", unitIdentifier))
}

func TestExecuteCreateNFTType_ParentDoesNotExist(t *testing.T) {
	txs := newTokenTxSystem(t)
	tx := testtransaction.NewGenericTransaction(
		t,
		txs.ConvertTx,
		testtransaction.WithUnitId(unitIdentifier.Bytes()),
		testtransaction.WithSystemID(DefaultTokenTxSystemIdentifier),
		testtransaction.WithAttributes(&CreateNonFungibleTokenTypeAttributes{
			Symbol:                   symbol,
			ParentTypeId:             parent1Identifier.Bytes(),
			SubTypeCreationPredicate: subTypeCreationPredicate,
		}),
	)
	require.ErrorContains(t, txs.Execute(tx), fmt.Sprintf("item %v does not exist", parent1Identifier))
}

func TestExecuteCreateNFTType_InvalidParentType(t *testing.T) {
	txs := newTokenTxSystem(t)
	txs.state.AddItem(parent1Identifier, script.PredicateAlwaysTrue(), &mockUnitData{}, []byte{})
	tx := testtransaction.NewGenericTransaction(
		t,
		txs.ConvertTx,
		testtransaction.WithUnitId(unitIdentifier.Bytes()),
		testtransaction.WithSystemID(DefaultTokenTxSystemIdentifier),
		testtransaction.WithAttributes(&CreateNonFungibleTokenTypeAttributes{
			Symbol:                   symbol,
			ParentTypeId:             parent1Identifier.Bytes(),
			SubTypeCreationPredicate: subTypeCreationPredicate,
		}),
	)
	require.ErrorContains(t, txs.Execute(tx), fmt.Sprintf("unit %v is not a non-fungible token type", parent1Identifier))
}

func TestExecuteCreateNFTType_InvalidSystemIdentifier(t *testing.T) {
	txs := newTokenTxSystem(t)
	tx := testtransaction.NewGenericTransaction(
		t,
		txs.ConvertTx,
		testtransaction.WithUnitId(unitIdentifier.Bytes()),
		testtransaction.WithSystemID([]byte{0, 0, 0, 0}),
		testtransaction.WithAttributes(&CreateNonFungibleTokenTypeAttributes{}),
	)
	require.ErrorContains(t, txs.Execute(tx), "invalid system identifier")
}

func TestExecuteCreateNFTType_InvalidTxType(t *testing.T) {
	txs, err := New(WithSystemIdentifier([]byte{0, 0, 0, 0}))
	require.NoError(t, err)
	tx := testtransaction.RandomGenericBillTransfer(t)
	require.ErrorContains(t, txs.Execute(tx), "unknown tx type")
}

func TestRevertTransaction_Ok(t *testing.T) {
	txs := newTokenTxSystem(t)
	tx := testtransaction.NewGenericTransaction(
		t,
		txs.ConvertTx,
		testtransaction.WithUnitId(unitIdentifier.Bytes()),
		testtransaction.WithSystemID(DefaultTokenTxSystemIdentifier),
		testtransaction.WithAttributes(&CreateNonFungibleTokenTypeAttributes{}),
	)
	require.NoError(t, txs.Execute(tx))
	txs.Revert()
	_, err := txs.state.GetUnit(unitIdentifier)
	require.ErrorContains(t, err, fmt.Sprintf("item %v does not exist", unitIdentifier))
}

func TestExecuteCreateNFTType_InvalidSymbolName(t *testing.T) {
	s := "♥♥♥♥♥♥♥♥ We ♥ Alphabill ♥♥♥♥♥♥♥♥"
	txs := newTokenTxSystem(t)
	tx := testtransaction.NewGenericTransaction(
		t,
		txs.ConvertTx,
		testtransaction.WithUnitId(unitIdentifier.Bytes()),
		testtransaction.WithSystemID(DefaultTokenTxSystemIdentifier),
		testtransaction.WithAttributes(&CreateNonFungibleTokenTypeAttributes{
			Symbol: s,
		}),
	)
	require.ErrorContains(t, txs.Execute(tx), ErrStringInvalidSymbolName)
}

type mockUnitData struct{}

func (m mockUnitData) AddToHasher(hash.Hash) {}

func (m mockUnitData) Value() rma.SummaryValue { return zeroSummaryValue }

func createSigner(t *testing.T) (crypto.Signer, []byte) {
	t.Helper()
	signer, err := crypto.NewInMemorySecp256K1Signer()
	require.NoError(t, err)

	verifier, err := signer.Verifier()
	require.NoError(t, err)

	pubKey, err := verifier.MarshalPublicKey()
	require.NoError(t, err)
	return signer, pubKey
}

func newTokenTxSystem(t *testing.T) *tokensTxSystem {
	txs, err := New()
	require.NoError(t, err)
	return txs
}

func signTx(t *testing.T, txs *tokensTxSystem, tx *txsystem.Transaction, signer crypto.Signer, pubKey []byte) (txsystem.GenericTransaction, []byte) {
	gtx, err := txs.ConvertTx(tx)
	require.NoError(t, err)

	signature, err := signer.SignBytes(gtx.SigBytes())
	require.NoError(t, err)

	tx.OwnerProof = script.PredicateArgumentPayToPublicKeyHashDefault(signature, pubKey)
	gtx, err = txs.ConvertTx(tx)
	require.NoError(t, err)
	return gtx, signature
}
