package backend

import (
	"bytes"
	"context"
	"crypto"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"syscall"
	"testing"
	"time"

	abcrypto "github.com/alphabill-org/alphabill/crypto"
	"github.com/alphabill-org/alphabill/internal/testutils"
	"github.com/alphabill-org/alphabill/internal/testutils/http"
	"github.com/alphabill-org/alphabill/internal/testutils/net"
	"github.com/alphabill-org/alphabill/internal/testutils/observability"
	"github.com/alphabill-org/alphabill/internal/testutils/transaction"
	"github.com/alphabill-org/alphabill/predicates/templates"
	"github.com/alphabill-org/alphabill/txsystem/money"
	"github.com/alphabill-org/alphabill/types"
	"github.com/alphabill-org/alphabill/util"
	"github.com/alphabill-org/alphabill/wallet/account"

	"github.com/ainvaltin/httpsrv"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/fxamacker/cbor/v2"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"

	sdk "github.com/alphabill-org/alphabill/wallet"
)

const (
	pubkeyHex = "0x000000000000000000000000000000000000000000000000000000000000000000"
)

var (
	billID            = money.NewBillID(nil, []byte{1})
	feeCreditRecordID = money.NewFeeCreditRecordID(nil, []byte{1})
)

func TestListBillsRequest_Ok(t *testing.T) {
	expectedBill := &Bill{
		Id:             newBillID(1),
		Value:          1,
		OwnerPredicate: getOwnerPredicate(pubkeyHex),
	}
	walletBackend := newWalletBackend(t, withBills(expectedBill))
	port, _ := startServer(t, walletBackend)

	res := &ListBillsResponse{}
	httpRes, err := testhttp.DoGetJson(fmt.Sprintf("http://localhost:%d/api/v1/list-bills?pubkey=%s", port, pubkeyHex), res)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, httpRes.StatusCode)
	require.Len(t, res.Bills, 1)
	expectedRes := expectedBill.ToGenericBills().Bills
	require.Equal(t, expectedRes, res.Bills)
}

func TestListBillsRequest_NilPubKey(t *testing.T) {
	port, _ := startServer(t, newWalletBackend(t))

	res := &sdk.ErrorResponse{}
	httpRes, err := testhttp.DoGetJson(fmt.Sprintf("http://localhost:%d/api/v1/list-bills", port), res)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, httpRes.StatusCode)
	require.Contains(t, res.Message, "must be 68 characters long (including 0x prefix), got 0 characters")
}

func TestListBillsRequest_InvalidPubKey(t *testing.T) {
	port, _ := startServer(t, newWalletBackend(t))

	res := &sdk.ErrorResponse{}
	pk := "0x00"
	httpRes, err := testhttp.DoGetJson(fmt.Sprintf("http://localhost:%d/api/v1/list-bills?pubkey=%s", port, pk), res)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, httpRes.StatusCode)
	require.Contains(t, res.Message, "must be 68 characters long (including 0x prefix), got 4 characters starting 0x00")
}

func TestListBillsRequest_DCBillsIncluded(t *testing.T) {
	walletBackend := newWalletBackend(t, withBills(
		&Bill{
			Id:             newBillID(1),
			Value:          1,
			OwnerPredicate: getOwnerPredicate(pubkeyHex),
		},
		&Bill{
			Id:             newBillID(2),
			Value:          2,
			DCTargetUnitID: []byte{2},
			OwnerPredicate: getOwnerPredicate(pubkeyHex),
		},
	))
	port, _ := startServer(t, walletBackend)

	res := &ListBillsResponse{}
	httpRes, err := testhttp.DoGetJson(fmt.Sprintf("http://localhost:%d/api/v1/list-bills?pubkey=%s", port, pubkeyHex), res)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, httpRes.StatusCode)
	require.Len(t, res.Bills, 2)
	bill := res.Bills[0]
	require.EqualValues(t, 1, bill.Value)
	require.Nil(t, bill.DCTargetUnitID)
	bill = res.Bills[1]
	require.EqualValues(t, 2, bill.Value)
	require.NotNil(t, bill.DCTargetUnitID)
}

func TestListBillsRequest_DCBillsExcluded(t *testing.T) {
	walletBackend := newWalletBackend(t, withBills(
		&Bill{
			Id:             newBillID(1),
			Value:          1,
			OwnerPredicate: getOwnerPredicate(pubkeyHex),
		},
		&Bill{
			Id:             newBillID(2),
			Value:          2,
			DCTargetUnitID: []byte{2},
			OwnerPredicate: getOwnerPredicate(pubkeyHex),
		},
	))
	port, _ := startServer(t, walletBackend)

	res := &ListBillsResponse{}
	httpRes, err := testhttp.DoGetJson(fmt.Sprintf("http://localhost:%d/api/v1/list-bills?pubkey=%s&includeDcBills=false", port, pubkeyHex), res)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, httpRes.StatusCode)
	require.Len(t, res.Bills, 1)
	bill := res.Bills[0]
	require.EqualValues(t, 1, bill.Value)
	require.Nil(t, bill.DCTargetUnitID)
}

func Test_txHistory(t *testing.T) {
	walletService := newWalletBackend(t)
	port, api := startServer(t, walletService)

	// setup account
	dir := t.TempDir()
	am, err := account.NewManager(dir, "", true)
	require.NoError(t, err)
	defer am.Close()
	err = am.CreateKeys("")
	require.NoError(t, err)
	accKey, err := am.GetAccountKey(0)
	require.NoError(t, err)

	makePostTxRequest := func(owner sdk.PubKey, body []byte) *http.Response {
		req := httptest.NewRequest("POST", fmt.Sprintf("http://localhost:%d/api/v1/transactions/0x%x", port, owner), bytes.NewReader(body))
		req = mux.SetURLVars(req, map[string]string{"pubkey": sdk.EncodeHex(owner)})
		w := httptest.NewRecorder()
		api.postTransactions(w, req)
		return w.Result()
	}

	makeTxHistoryRequest := func(pubkey sdk.PubKey) *http.Response {
		req := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:%d/api/v1/tx-history/0x%x", port, pubkey), nil)
		req = mux.SetURLVars(req, map[string]string{"pubkey": sdk.EncodeHex(pubkey)})
		w := httptest.NewRecorder()
		api.txHistoryFunc(w, req)
		return w.Result()
	}

	pubkey2 := sdk.PubKey(test.RandomBytes(33))
	bearerPredicate := templates.NewP2pkh256BytesFromKeyHash(pubkey2.Hash())
	attrs := &money.TransferAttributes{NewBearer: bearerPredicate}
	txs := sdk.Transactions{Transactions: []*types.TransactionOrder{
		testtransaction.NewTransactionOrder(t, testtransaction.WithPayloadType(money.PayloadTypeTransfer), testtransaction.WithAttributes(attrs))},
	}
	txBytes, err := txs.Transactions[0].PayloadBytes()
	signer, _ := abcrypto.NewInMemorySecp256K1SignerFromKey(accKey.PrivKey)
	require.NoError(t, err)
	sigData, _ := signer.SignBytes(txBytes)
	txs.Transactions[0].OwnerProof = templates.NewP2pkh256SignatureBytes(sigData, accKey.PubKey)

	b, err := cbor.Marshal(txs)
	resp := makePostTxRequest(accKey.PubKey, b)
	require.Equal(t, http.StatusAccepted, resp.StatusCode)

	txHistResp := makeTxHistoryRequest(accKey.PubKey)
	require.Equal(t, http.StatusOK, txHistResp.StatusCode)

	buf, err := io.ReadAll(txHistResp.Body)
	require.NoError(t, err)
	var txHistory []*sdk.TxHistoryRecord
	require.NoError(t, cbor.Unmarshal(buf, &txHistory))
	require.Len(t, txHistory, 1)
	require.Equal(t, sdk.OUTGOING, txHistory[0].Kind)
	require.Equal(t, sdk.UNCONFIRMED, txHistory[0].State)
	require.EqualValues(t, pubkey2.Hash(), txHistory[0].CounterParty)
}

func TestListBillsRequest_Paging(t *testing.T) {
	// given set of bills
	var bills []*Bill
	for i := uint64(1); i <= 200; i++ {
		bills = append(bills, &Bill{
			Id:             newBillID(byte(i)),
			Value:          i,
			OwnerPredicate: getOwnerPredicate(pubkeyHex),
		})
	}
	walletService := newWalletBackend(t, withBills(bills...))
	port, _ := startServer(t, walletService)

	// verify by default first 100 elements are returned
	res := &ListBillsResponse{}
	httpRes, err := testhttp.DoGetJson(fmt.Sprintf("http://localhost:%d/api/v1/list-bills?pubkey=%s", port, pubkeyHex), res)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, httpRes.StatusCode)
	require.Len(t, res.Bills, 100)
	require.EqualValues(t, 1, res.Bills[0].Value)
	require.EqualValues(t, 100, res.Bills[99].Value)
	verifyLinkHeader(t, httpRes, bills[100].Id)

	// verify offsetKey=100 returns next 100 elements
	res = &ListBillsResponse{}
	httpRes, err = testhttp.DoGetJson(fmt.Sprintf("http://localhost:%d/api/v1/list-bills?pubkey=%s&offsetKey=%s", port, pubkeyHex, hexutil.Encode(bills[100].Id)), res)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, httpRes.StatusCode)
	require.Len(t, res.Bills, 100)
	require.EqualValues(t, 101, res.Bills[0].Value)
	require.EqualValues(t, 200, res.Bills[99].Value)
	verifyNoLinkHeader(t, httpRes)

	// verify limit limits result size
	res = &ListBillsResponse{}
	httpRes, err = testhttp.DoGetJson(fmt.Sprintf("http://localhost:%d/api/v1/list-bills?pubkey=%s&offsetKey=%s&limit=50", port, pubkeyHex, hexutil.Encode(bills[100].Id)), res)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, httpRes.StatusCode)
	require.Len(t, res.Bills, 50)
	require.EqualValues(t, 101, res.Bills[0].Value)
	require.EqualValues(t, 150, res.Bills[49].Value)
	verifyLinkHeader(t, httpRes, bills[150].Id)

	// verify out of bounds offset returns nothing
	res = &ListBillsResponse{}
	httpRes, err = testhttp.DoGetJson(fmt.Sprintf("http://localhost:%d/api/v1/list-bills?pubkey=%s&offsetKey=%s", port, pubkeyHex, hexutil.Encode(util.Uint64ToBytes32(201))), res)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, httpRes.StatusCode)
	require.Len(t, res.Bills, 0)
	verifyNoLinkHeader(t, httpRes)

	// verify limit gets capped to 100
	res = &ListBillsResponse{}
	httpRes, err = testhttp.DoGetJson(fmt.Sprintf("http://localhost:%d/api/v1/list-bills?pubkey=%s&limit=200", port, pubkeyHex), res)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, httpRes.StatusCode)
	require.Len(t, res.Bills, 100)
	require.EqualValues(t, 1, res.Bills[0].Value)
	require.EqualValues(t, 100, res.Bills[99].Value)
	verifyLinkHeader(t, httpRes, bills[100].Id)

	// verify out of bounds offset+limit returns all data starting from offset
	res = &ListBillsResponse{}
	httpRes, err = testhttp.DoGetJson(fmt.Sprintf("http://localhost:%d/api/v1/list-bills?pubkey=%s&offsetKey=%s&limit=100", port, pubkeyHex, hexutil.Encode(bills[190].Id)), res)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, httpRes.StatusCode)
	require.Len(t, res.Bills, 10)
	require.EqualValues(t, 191, res.Bills[0].Value)
	require.EqualValues(t, 200, res.Bills[9].Value)
	verifyNoLinkHeader(t, httpRes)
}

func TestListBillsRequest_PagingWithDCBills(t *testing.T) {
	// create 30 bills where first 10 and last 10 are dc-bills
	var bills []*Bill
	for i := byte(1); i <= 30; i++ {
		var b *Bill
		if i <= 10 || i > 20 {
			b = &Bill{
				Id:                   newBillID(i),
				Value:                uint64(i),
				OwnerPredicate:       getOwnerPredicate(pubkeyHex),
				DCTargetUnitID:       test.RandomBytes(32),
				DCTargetUnitBacklink: test.RandomBytes(32),
			}
		} else {
			b = &Bill{
				Id:             newBillID(i),
				Value:          uint64(i),
				OwnerPredicate: getOwnerPredicate(pubkeyHex),
			}
		}
		bills = append(bills, b)
	}
	walletService := newWalletBackend(t, withBills(bills...))
	port, _ := startServer(t, walletService)

	// verify first 10 non-dc bills are returned
	res := &ListBillsResponse{}
	httpRes, err := testhttp.DoGetJson(fmt.Sprintf("http://localhost:%d/api/v1/list-bills?pubkey=%s&includeDcBills=false&limit=10", port, pubkeyHex), res)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, httpRes.StatusCode)
	require.Len(t, res.Bills, 10)
	require.EqualValues(t, 11, res.Bills[0].Value)
	require.EqualValues(t, 20, res.Bills[9].Value)
	verifyLinkHeader(t, httpRes, bills[20].Id)
}

func TestBalanceRequest_Ok(t *testing.T) {
	port, _ := startServer(t, newWalletBackend(t, withBills(
		&Bill{
			Id:             newBillID(1),
			Value:          1,
			OwnerPredicate: getOwnerPredicate(pubkeyHex),
		})))

	res := &BalanceResponse{}
	httpRes, err := testhttp.DoGetJson(fmt.Sprintf("http://localhost:%d/api/v1/balance?pubkey=%s", port, pubkeyHex), res)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, httpRes.StatusCode)
	require.EqualValues(t, 1, res.Balance)
}

func TestBalanceRequest_NilPubKey(t *testing.T) {
	port, _ := startServer(t, newWalletBackend(t))

	res := &sdk.ErrorResponse{}
	httpRes, err := testhttp.DoGetJson(fmt.Sprintf("http://localhost:%d/api/v1/balance", port), res)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, httpRes.StatusCode)
	require.Contains(t, res.Message, "must be 68 characters long (including 0x prefix), got 0 characters")
}

func TestBalanceRequest_InvalidPubKey(t *testing.T) {
	port, _ := startServer(t, newWalletBackend(t))

	res := &sdk.ErrorResponse{}
	pk := "0x00"
	httpRes, err := testhttp.DoGetJson(fmt.Sprintf("http://localhost:%d/api/v1/balance?pubkey=%s", port, pk), res)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, httpRes.StatusCode)
	require.Contains(t, res.Message, "must be 68 characters long (including 0x prefix), got 4 characters starting 0x00")
}

func TestBalanceRequest_DCBillNotIncluded(t *testing.T) {
	walletBackend := newWalletBackend(t, withBills(
		&Bill{
			Id:             newBillID(1),
			Value:          1,
			OwnerPredicate: getOwnerPredicate(pubkeyHex),
		},
		&Bill{
			Id:             newBillID(2),
			Value:          2,
			DCTargetUnitID: []byte{2},
			OwnerPredicate: getOwnerPredicate(pubkeyHex),
		}),
	)
	port, _ := startServer(t, walletBackend)

	res := &BalanceResponse{}
	httpRes, err := testhttp.DoGetJson(fmt.Sprintf("http://localhost:%d/api/v1/balance?pubkey=%s", port, pubkeyHex), res)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, httpRes.StatusCode)
	require.EqualValues(t, 1, res.Balance)
}

func TestProofRequest_Ok(t *testing.T) {
	tr := testtransaction.NewTransactionRecord(t)
	txHash := tr.TransactionOrder.Hash(crypto.SHA256)
	b := &Bill{
		Id:             money.NewBillID(nil, []byte{1}),
		Value:          1,
		TxHash:         txHash,
		OwnerPredicate: getOwnerPredicate(pubkeyHex),
	}
	p := &sdk.Proof{
		TxRecord: tr,
		TxProof: &types.TxProof{
			BlockHeaderHash:    []byte{0},
			Chain:              []*types.GenericChainItem{{Hash: []byte{0}}},
			UnicityCertificate: &types.UnicityCertificate{InputRecord: &types.InputRecord{RoundNumber: 1}},
		},
	}
	walletBackend := newWalletBackend(t, withBillProofs(&billProof{b, p}))
	port, _ := startServer(t, walletBackend)

	response := &sdk.Proof{}
	httpRes, err := testhttp.DoGetCbor(fmt.Sprintf("http://localhost:%d/api/v1/units/0x%s/transactions/0x%x/proof", port, billID, b.TxHash), response)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, httpRes.StatusCode)
	require.Equal(t, b.TxHash, response.TxRecord.TransactionOrder.Hash(crypto.SHA256))
	//
	require.Equal(t, p.TxProof.UnicityCertificate.GetRoundNumber(), response.TxProof.UnicityCertificate.GetRoundNumber())
	require.EqualValues(t, p.TxRecord.TransactionOrder.UnitID(), response.TxRecord.TransactionOrder.UnitID())
	require.EqualValues(t, p.TxProof.BlockHeaderHash, response.TxProof.BlockHeaderHash)
}

func TestProofRequest_InvalidBillIdLength(t *testing.T) {
	port, _ := startServer(t, newWalletBackend(t))

	// verify bill id larger than 33 bytes returns error
	res := &sdk.ErrorResponse{}
	billID := test.RandomBytes(34)
	httpRes, err := testhttp.DoGetJson(fmt.Sprintf("http://localhost:%d/api/v1/units/0x%x/transactions/0x00/proof", port, billID), res)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, httpRes.StatusCode)
	require.Equal(t, errInvalidBillIDLength.Error(), res.Message)

	// verify bill id smaller than 33 bytes returns error
	res = &sdk.ErrorResponse{}
	httpRes, err = testhttp.DoGetJson(fmt.Sprintf("http://localhost:%d/api/v1/units/0x01/transactions/0x00/proof", port), res)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, httpRes.StatusCode)
	require.Equal(t, errInvalidBillIDLength.Error(), res.Message)

	// verify bill id with correct length but missing prefix returns error
	res = &sdk.ErrorResponse{}
	httpRes, err = testhttp.DoGetJson(fmt.Sprintf("http://localhost:%d/api/v1/units/%x/transactions/0x00/proof", port, billID), res)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, httpRes.StatusCode)
	require.Contains(t, res.Message, "hex string without 0x prefix")
}

func TestProofRequest_ProofDoesNotExist(t *testing.T) {
	port, _ := startServer(t, newWalletBackend(t))

	res := &sdk.ErrorResponse{}
	httpRes, err := testhttp.DoGetJson(fmt.Sprintf("http://localhost:%d/api/v1/units/0x%s/transactions/0x00/proof", port, billID), res)
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, httpRes.StatusCode)
	require.Contains(t, res.Message, fmt.Sprintf("no proof found for tx 0x00 (unit 0x%s)", billID))
}

func TestRoundNumberRequest_Ok(t *testing.T) {
	nodeRoundNumber := uint64(150)
	backendRoundNumber := uint64(10)
	alphabillClient := &mockABClient{
		getRoundNumber: func(ctx context.Context) (uint64, error) { return nodeRoundNumber, nil },
	}
	service := newWalletBackend(t, withABClient(alphabillClient))
	port, _ := startServer(t, service)
	err := service.store.Do().SetBlockNumber(backendRoundNumber)
	require.NoError(t, err)

	res := &sdk.RoundNumber{}
	httpRes, err := testhttp.DoGetJson(fmt.Sprintf("http://localhost:%d/api/v1/round-number", port), res)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, httpRes.StatusCode)
	require.EqualValues(t, nodeRoundNumber, res.RoundNumber)
	require.EqualValues(t, backendRoundNumber, res.LastIndexedRoundNumber)
}

func TestInvalidUrl_NotFound(t *testing.T) {
	port, _ := startServer(t, newWalletBackend(t))

	// verify request to to non-existent /api2 endpoint returns 404
	httpRes, err := http.Get(fmt.Sprintf("http://localhost:%d/api2/v1/list-bills", port))
	require.NoError(t, err)
	require.Equal(t, 404, httpRes.StatusCode)

	// verify request to to non-existent version endpoint returns 404
	httpRes, err = http.Get(fmt.Sprintf("http://localhost:%d/api/v5/list-bills", port))
	require.NoError(t, err)
	require.Equal(t, 404, httpRes.StatusCode)
}

func TestGetFeeCreditBillRequest_Ok(t *testing.T) {
	b := &Bill{
		Id:             feeCreditRecordID,
		Value:          1,
		TxHash:         []byte{0},
		OwnerPredicate: getOwnerPredicate(pubkeyHex),
	}
	walletBackend := newWalletBackend(t, withFeeCreditBills(b))
	port, _ := startServer(t, walletBackend)

	response := &sdk.Bill{}
	httpRes, err := testhttp.DoGetJson(fmt.Sprintf("http://localhost:%d/api/v1/fee-credit-bills/0x%s", port, feeCreditRecordID), response)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, httpRes.StatusCode)
	require.Equal(t, b.Id, response.Id)
	require.Equal(t, b.Value, response.Value)
	require.Equal(t, b.TxHash, response.TxHash)
	require.Equal(t, b.DCTargetUnitID, response.DCTargetUnitID)
}

func TestGetFeeCreditBillRequest_InvalidBillIdLength(t *testing.T) {
	port, _ := startServer(t, newWalletBackend(t))

	// verify bill id larger than 33 bytes returns error
	res := &sdk.ErrorResponse{}
	billID := "0x00000000000000000000000000000000000000000000000000000000000000000101"
	httpRes, err := testhttp.DoGetJson(fmt.Sprintf("http://localhost:%d/api/v1/fee-credit-bills/%s", port, billID), res)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, httpRes.StatusCode)
	require.Equal(t, "bill_id hex string must be 68 characters long (with 0x prefix)", res.Message)

	// verify bill id smaller than 33 bytes returns error
	res = &sdk.ErrorResponse{}
	httpRes, err = testhttp.DoGetJson(fmt.Sprintf("http://localhost:%d/api/v1/fee-credit-bills/0x01", port), res)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, httpRes.StatusCode)
	require.Equal(t, "bill_id hex string must be 68 characters long (with 0x prefix)", res.Message)

	// verify bill id with correct length but missing prefix returns error
	res = &sdk.ErrorResponse{}
	httpRes, err = testhttp.DoGetJson(fmt.Sprintf("http://localhost:%d/api/v1/fee-credit-bills/%s", port, feeCreditRecordID), res)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, httpRes.StatusCode)
	require.Contains(t, res.Message, "hex string without 0x prefix")
}

func TestGetFeeCreditBillRequest_BillDoesNotExist(t *testing.T) {
	port, _ := startServer(t, newWalletBackend(t))

	res := &sdk.ErrorResponse{}
	httpRes, err := testhttp.DoGetJson(fmt.Sprintf("http://localhost:%d/api/v1/fee-credit-bills/0x%s", port, billID), res)
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, httpRes.StatusCode)
	require.Equal(t, "fee credit bill does not exist", res.Message)
}

func TestPostTransactionsRequest_InvalidPubkey(t *testing.T) {
	walletBackend := newWalletBackend(t)
	port, _ := startServer(t, walletBackend)

	res := &sdk.ErrorResponse{}
	httpRes, err := testhttp.DoPost(fmt.Sprintf("http://localhost:%d/api/v1/transactions/%s", port, "invalid"), nil, res)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, httpRes.StatusCode)
	require.Contains(t, res.Message, "failed to parse sender pubkey")
}

func TestPostTransactionsRequest_EmptyBody(t *testing.T) {
	walletBackend := newWalletBackend(t)
	port, _ := startServer(t, walletBackend)

	res := &sdk.ErrorResponse{}
	httpRes, err := testhttp.DoPostCBOR(fmt.Sprintf("http://localhost:%d/api/v1/transactions/%s", port, pubkeyHex), nil, res)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, httpRes.StatusCode)
	require.Equal(t, "request body contained no transactions to process", res.Message)
}

func TestPostTransactionsRequest_Ok(t *testing.T) {
	walletBackend := newWalletBackend(t)
	port, _ := startServer(t, walletBackend)

	// setup account
	dir := t.TempDir()
	am, err := account.NewManager(dir, "", true)
	require.NoError(t, err)
	defer am.Close()
	err = am.CreateKeys("")
	require.NoError(t, err)
	accKey, err := am.GetAccountKey(0)
	require.NoError(t, err)

	txs := &sdk.Transactions{Transactions: []*types.TransactionOrder{
		testtransaction.NewTransactionOrder(t),
		testtransaction.NewTransactionOrder(t),
		testtransaction.NewTransactionOrder(t),
	}}
	signer, _ := abcrypto.NewInMemorySecp256K1SignerFromKey(accKey.PrivKey)
	txBytes1, _ := txs.Transactions[0].PayloadBytes()
	sigData, _ := signer.SignBytes(txBytes1)
	txs.Transactions[0].OwnerProof = templates.NewP2pkh256SignatureBytes(sigData, accKey.PubKey)
	txBytes2, _ := txs.Transactions[1].PayloadBytes()
	sigData, _ = signer.SignBytes(txBytes2)
	txs.Transactions[1].OwnerProof = templates.NewP2pkh256SignatureBytes(sigData, accKey.PubKey)
	txBytes3, _ := txs.Transactions[2].PayloadBytes()
	sigData, _ = signer.SignBytes(txBytes3)
	txs.Transactions[2].OwnerProof = templates.NewP2pkh256SignatureBytes(sigData, accKey.PubKey)

	httpRes, err := testhttp.DoPostCBOR(fmt.Sprintf("http://localhost:%d/api/v1/transactions/%s", port, hexutil.Encode(accKey.PubKey)), txs, nil)
	require.NoError(t, err)
	require.Equal(t, http.StatusAccepted, httpRes.StatusCode)
}

func TestPostTransactionsRequest_InvalidOwner(t *testing.T) {
	walletBackend := newWalletBackend(t)
	port, _ := startServer(t, walletBackend)

	// setup account
	dir := t.TempDir()
	am, err := account.NewManager(dir, "", true)
	require.NoError(t, err)
	defer am.Close()
	err = am.CreateKeys("")
	require.NoError(t, err)
	accKey, err := am.GetAccountKey(0)
	require.NoError(t, err)

	txs := &sdk.Transactions{Transactions: []*types.TransactionOrder{
		testtransaction.NewTransactionOrder(t),
		testtransaction.NewTransactionOrder(t),
		testtransaction.NewTransactionOrder(t),
	}}
	signer, _ := abcrypto.NewInMemorySecp256K1SignerFromKey(accKey.PrivKey)
	txBytes1, _ := txs.Transactions[0].PayloadBytes()
	sigData, _ := signer.SignBytes(txBytes1)
	txs.Transactions[0].OwnerProof = templates.NewP2pkh256SignatureBytes(sigData, accKey.PubKey)
	txBytes2, _ := txs.Transactions[1].PayloadBytes()
	sigData, _ = signer.SignBytes(txBytes2)
	txs.Transactions[1].OwnerProof = templates.NewP2pkh256SignatureBytes(sigData, accKey.PubKey)
	txBytes3, _ := txs.Transactions[2].PayloadBytes()
	sigData, _ = signer.SignBytes(txBytes3)
	txs.Transactions[2].OwnerProof = templates.NewP2pkh256SignatureBytes(sigData, test.RandomBytes(33))

	httpRes, err := testhttp.DoPostCBOR(fmt.Sprintf("http://localhost:%d/api/v1/transactions/%s", port, hexutil.Encode(accKey.PubKey)), txs, nil)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, httpRes.StatusCode)
}

func TestInfoRequest_Ok(t *testing.T) {
	service := newWalletBackend(t)
	port, _ := startServer(t, service)

	var res *sdk.InfoResponse
	httpRes, err := testhttp.DoGetJson(fmt.Sprintf("http://localhost:%d/api/v1/info", port), &res)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, httpRes.StatusCode)
	require.Equal(t, "00000000", res.SystemID)
	require.Equal(t, "money backend", res.Name)
}

func verifyLinkHeader(t *testing.T, httpRes *http.Response, nextKey []byte) {
	var linkHdrMatcher = regexp.MustCompile("<(.*)>")
	match := linkHdrMatcher.FindStringSubmatch(httpRes.Header.Get(sdk.HeaderLink))
	if len(match) != 2 {
		t.Errorf("Link header didn't result in expected match\nHeader: %s\nmatches: %v\n", httpRes.Header.Get(sdk.HeaderLink), match)
	} else {
		u, err := url.Parse(match[1])
		if err != nil {
			t.Fatal("failed to parse Link header:", err)
		}
		if s := u.Query().Get(sdk.QueryParamOffsetKey); s != hexutil.Encode(nextKey) {
			t.Errorf("expected %x got %s", nextKey, s)
		}
	}
}

func verifyNoLinkHeader(t *testing.T, httpRes *http.Response) {
	if link := httpRes.Header.Get(sdk.HeaderLink); link != "" {
		t.Errorf("unexpectedly the Link header is not empty, got %q", link)
	}
}

type (
	option func(service *WalletBackend) error
)

func newWalletBackend(t *testing.T, options ...option) *WalletBackend {
	storage := createTestBillStore(t)
	mabc := &mockABClient{
		getRoundNumber:  func(ctx context.Context) (uint64, error) { return 0, nil },
		sendTransaction: func(ctx context.Context, tx *types.TransactionOrder) error { return nil },
	}
	service := &WalletBackend{store: storage, abc: mabc}
	for _, o := range options {
		err := o(service)
		require.NoError(t, err)
	}
	return service
}

func withBills(bills ...*Bill) option {
	return func(s *WalletBackend) error {
		return s.store.WithTransaction(func(tx BillStoreTx) error {
			for _, bill := range bills {
				err := tx.SetBill(bill, nil)
				if err != nil {
					return err
				}
			}
			return nil
		})
	}
}

type billProof struct {
	bill  *Bill
	proof *sdk.Proof
}

func withBillProofs(bills ...*billProof) option {
	return func(s *WalletBackend) error {
		return s.store.WithTransaction(func(tx BillStoreTx) error {
			for _, bill := range bills {
				err := tx.SetBill(bill.bill, bill.proof)
				if err != nil {
					return err
				}
			}
			return nil
		})
	}
}

func withABClient(client ABClient) option {
	return func(s *WalletBackend) error {
		s.abc = client
		return nil
	}
}

func withFeeCreditBills(bills ...*Bill) option {
	return func(s *WalletBackend) error {
		return s.store.WithTransaction(func(tx BillStoreTx) error {
			for _, bill := range bills {
				err := tx.SetFeeCreditBill(bill, nil)
				if err != nil {
					return err
				}
			}
			return nil
		})
	}
}

func startServer(t *testing.T, service WalletBackendService) (port int, api *moneyRestAPI) {
	port, err := net.GetFreePort()
	require.NoError(t, err)
	observe := observability.Default(t)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		api = &moneyRestAPI{Service: service, ListBillsPageLimit: 100, rw: &sdk.ResponseWriter{}, log: observe.Logger(), tracer: observe.Tracer("moneyRestAPI"), SystemID: moneySystemID}
		err := httpsrv.Run(ctx,
			http.Server{
				Addr:              fmt.Sprintf("localhost:%d", port),
				Handler:           api.Router(),
				ReadTimeout:       3 * time.Second,
				ReadHeaderTimeout: time.Second,
				WriteTimeout:      5 * time.Second,
				IdleTimeout:       30 * time.Second,
			},
			httpsrv.ShutdownTimeout(5*time.Second),
		)
		require.ErrorIs(t, err, context.Canceled)
	}()
	// stop the server
	t.Cleanup(func() { cancel() })

	// wait until server is up
	tout := time.After(1500 * time.Millisecond)
	for {
		if _, err := http.Get(fmt.Sprintf("http://localhost:%d", port)); err != nil {
			if !errors.Is(err, syscall.ECONNREFUSED) {
				t.Fatalf("unexpected error from http server: %v", err)
			}
		} else {
			return port, api
		}

		select {
		case <-time.After(50 * time.Millisecond):
		case <-tout:
			t.Fatalf("http server didn't become available within timeout")
		}
	}
}

type mockABClient struct {
	sendTransaction func(ctx context.Context, tx *types.TransactionOrder) error
	getRoundNumber  func(ctx context.Context) (uint64, error)
}

func (mc *mockABClient) SendTransaction(ctx context.Context, tx *types.TransactionOrder) error {
	return mc.sendTransaction(ctx, tx)
}
func (mc *mockABClient) GetRoundNumber(ctx context.Context) (uint64, error) {
	return mc.getRoundNumber(ctx)
}