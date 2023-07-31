package evm

import (
	"crypto"
	"crypto/sha256"
	"math/big"
	"testing"

	"github.com/alphabill-org/alphabill/internal/state"

	abcrypto "github.com/alphabill-org/alphabill/internal/crypto"
	"github.com/alphabill-org/alphabill/internal/script"
	testsig "github.com/alphabill-org/alphabill/internal/testutils/sig"
	testtransaction "github.com/alphabill-org/alphabill/internal/testutils/transaction"
	"github.com/alphabill-org/alphabill/internal/txsystem/evm/statedb"
	"github.com/alphabill-org/alphabill/internal/txsystem/fc"
	testfc "github.com/alphabill-org/alphabill/internal/txsystem/fc/testutils"
	"github.com/alphabill-org/alphabill/internal/txsystem/fc/transactions"
	"github.com/alphabill-org/alphabill/internal/types"
	"github.com/fxamacker/cbor/v2"
	"github.com/stretchr/testify/require"
)

func hashOfPrivateKey(t *testing.T, signer abcrypto.Signer) []byte {
	t.Helper()
	privKeyBytes, err := signer.MarshalPrivateKey()
	require.NoError(t, err)
	h := sha256.Sum256(privKeyBytes)
	return h[:]
}

func newTxPayload(t *testing.T, txType string, unitID []byte, timeout uint64, fcrID []byte, attr interface{}) *types.Payload {
	attrBytes, err := cbor.Marshal(attr)
	require.NoError(t, err)
	return &types.Payload{
		SystemID:   DefaultEvmTxSystemIdentifier,
		Type:       txType,
		UnitID:     unitID,
		Attributes: attrBytes,
		ClientMetadata: &types.ClientMetadata{
			Timeout:           timeout,
			MaxTransactionFee: 1,
			FeeCreditRecordID: fcrID,
		},
	}
}

func newAddFCTx(t *testing.T, unitID []byte, attr *transactions.AddFeeCreditAttributes, signer abcrypto.Signer, timeout uint64) *types.TransactionOrder {
	payload := newTxPayload(t, transactions.PayloadTypeAddFeeCredit, unitID, timeout, nil, attr)
	payloadBytes, err := payload.Bytes()
	require.NoError(t, err)
	sig, err := signer.SignBytes(payloadBytes)
	require.NoError(t, err)
	ver, err := signer.Verifier()
	require.NoError(t, err)
	pubKeyBytes, err := ver.MarshalPublicKey()
	require.NoError(t, err)
	return &types.TransactionOrder{
		Payload:    payload,
		OwnerProof: script.PredicateArgumentPayToPublicKeyHashDefault(sig, pubKeyBytes),
	}
}

func evmTestFeeCalculator() uint64 {
	return 2
}

func Test_addFeeCreditTx(t *testing.T) {
	type args struct {
		order       *types.TransactionOrder
		blockNumber uint64
	}
	signer, ver := testsig.CreateSignerAndVerifier(t)
	tb := map[string]abcrypto.Verifier{"test": ver}
	pubKeyBytes, err := ver.MarshalPublicKey()
	require.NoError(t, err)
	pubHash := sha256.Sum256(pubKeyBytes)
	privKeyHash := hashOfPrivateKey(t, signer)
	addExecFn := addFeeCreditTx(
		state.NewEmptyState(),
		crypto.SHA256,
		evmTestFeeCalculator,
		fc.NewDefaultFeeCreditTxValidator([]byte{0, 0, 0, 0}, DefaultEvmTxSystemIdentifier, crypto.SHA256, tb))

	tests := []struct {
		name       string
		args       args
		wantErrStr string
	}{
		{
			name:       "err - invalid owner proof",
			args:       args{order: testfc.NewAddFC(t, signer, nil), blockNumber: 5},
			wantErrStr: "failed to extract public key from fee credit owner proof",
		},
		{
			name:       "err - attr:FeeCreditOwnerCondition is nil",
			args:       args{order: newAddFCTx(t, privKeyHash, nil, signer, 7), blockNumber: 5},
			wantErrStr: "addFC tx validation failed: fee credit owner condition is nil",
		},
		{
			name: "ok",
			args: args{order: newAddFCTx(t,
				hashOfPrivateKey(t, signer),
				testfc.NewAddFCAttr(t, signer, testfc.WithTransferFCTx(
					&types.TransactionRecord{
						TransactionOrder: testfc.NewTransferFC(t, testfc.NewTransferFCAttr(testfc.WithAmount(100), testfc.WithTargetRecordID(privKeyHash), testfc.WithTargetSystemID(DefaultEvmTxSystemIdentifier)),
							testtransaction.WithSystemID([]byte{0, 0, 0, 0}), testtransaction.WithOwnerProof(script.PredicatePayToPublicKeyHashDefault(pubHash[:]))),
						ServerMetadata: nil,
					})),
				signer,
				7,
			),
				blockNumber: 5},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attr := new(transactions.AddFeeCreditAttributes)
			require.NoError(t, tt.args.order.UnmarshalAttributes(attr))
			metaData, err := addExecFn(tt.args.order, attr, tt.args.blockNumber)
			if tt.wantErrStr != "" {
				require.ErrorContains(t, err, tt.wantErrStr)
				require.Nil(t, metaData)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_getTransferPayloadAttributes(t *testing.T) {
	type args struct {
		transfer *types.TransactionRecord
	}
	signer, _ := testsig.CreateSignerAndVerifier(t)
	addFcAttr := testfc.NewAddFCAttr(t, signer, testfc.WithTransferFCTx(
		&types.TransactionRecord{
			TransactionOrder: testfc.NewTransferFC(t, nil, testtransaction.WithSystemID([]byte("not money partition"))),
			ServerMetadata:   nil,
		}))
	closeFCAmount := uint64(20)
	closeFCFee := uint64(2)
	closeFCAttr := testfc.NewCloseFCAttr(testfc.WithCloseFCAmount(closeFCAmount))
	closureTx := testfc.WithReclaimFCClosureTx(
		&types.TransactionRecord{
			TransactionOrder: testfc.NewCloseFC(t, closeFCAttr),
			ServerMetadata:   &types.ServerMetadata{ActualFee: closeFCFee},
		},
	)
	newReclaimFCAttr := testfc.NewReclaimFCAttr(t, signer, closureTx)
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "nil",
			args:    args{transfer: nil},
			wantErr: true,
		},
		{
			name:    "incorrect type",
			args:    args{transfer: newReclaimFCAttr.CloseFeeCreditTransfer},
			wantErr: true,
		},
		{
			name:    "ok",
			args:    args{transfer: addFcAttr.FeeCreditTransfer},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := getTransferPayloadAttributes(tt.args.transfer)
			if (err != nil) != tt.wantErr {
				t.Errorf("getTransferPayloadAttributes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_addFeeCreditTxAndUpdate(t *testing.T) {
	stateTree := state.NewEmptyState()
	signer, ver := testsig.CreateSignerAndVerifier(t)
	tb := map[string]abcrypto.Verifier{"test": ver}
	pubKeyBytes, err := ver.MarshalPublicKey()
	require.NoError(t, err)
	pubHash := sha256.Sum256(pubKeyBytes)
	privKeyHash := hashOfPrivateKey(t, signer)
	addExecFn := addFeeCreditTx(
		stateTree,
		crypto.SHA256,
		evmTestFeeCalculator,
		fc.NewDefaultFeeCreditTxValidator([]byte{0, 0, 0, 0}, DefaultEvmTxSystemIdentifier, crypto.SHA256, tb))
	addFeeOrder := newAddFCTx(t,
		privKeyHash,
		testfc.NewAddFCAttr(t, signer, testfc.WithTransferFCTx(
			&types.TransactionRecord{
				TransactionOrder: testfc.NewTransferFC(t, testfc.NewTransferFCAttr(testfc.WithAmount(100), testfc.WithTargetRecordID(privKeyHash), testfc.WithTargetSystemID(DefaultEvmTxSystemIdentifier)),
					testtransaction.WithSystemID([]byte{0, 0, 0, 0}), testtransaction.WithOwnerProof(script.PredicatePayToPublicKeyHashDefault(pubHash[:]))),
				ServerMetadata: nil,
			})),
		signer, 7)
	backlink := addFeeOrder.Hash(crypto.SHA256)
	attr := new(transactions.AddFeeCreditAttributes)
	require.NoError(t, addFeeOrder.UnmarshalAttributes(attr))
	metaData, err := addExecFn(addFeeOrder, attr, 5)
	require.NoError(t, err)
	require.NotNil(t, metaData)
	require.EqualValues(t, evmTestFeeCalculator(), metaData.ActualFee)
	// validate stateDB
	stateDB := statedb.NewStateDB(stateTree)
	addr, err := generateAddress(pubKeyBytes)
	require.NoError(t, err)
	balance := stateDB.GetBalance(addr)
	// balance is equal to 100-"fee = 2" to wei
	require.EqualValues(t, balance, new(big.Int).Sub(alphaToWei(100), alphaToWei(evmTestFeeCalculator())))
	// add more funds
	addFeeOrder = newAddFCTx(t,
		privKeyHash,
		testfc.NewAddFCAttr(t, signer, testfc.WithTransferFCTx(
			&types.TransactionRecord{
				TransactionOrder: testfc.NewTransferFC(t, testfc.NewTransferFCAttr(testfc.WithAmount(10), testfc.WithTargetRecordID(privKeyHash), testfc.WithTargetSystemID(DefaultEvmTxSystemIdentifier), testfc.WithNonce(backlink)),
					testtransaction.WithSystemID([]byte{0, 0, 0, 0}), testtransaction.WithOwnerProof(script.PredicatePayToPublicKeyHashDefault(pubHash[:]))),
				ServerMetadata: nil,
			})),
		signer, 7)
	require.NoError(t, addFeeOrder.UnmarshalAttributes(attr))
	metaData, err = addExecFn(addFeeOrder, attr, 5)
	require.NoError(t, err)
	balance = stateDB.GetBalance(addr)
	// balance is equal to 100-"fee = 2" to wei
	require.EqualValues(t, balance, new(big.Int).Sub(alphaToWei(110), alphaToWei(2*evmTestFeeCalculator())))
}
