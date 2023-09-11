package api

import (
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alphabill-org/alphabill/internal/rpc"
	abstate "github.com/alphabill-org/alphabill/internal/state"
	test "github.com/alphabill-org/alphabill/internal/testutils"
	"github.com/alphabill-org/alphabill/internal/txsystem/evm/statedb"
	"github.com/ethereum/go-ethereum/common"
	"github.com/fxamacker/cbor/v2"
	"github.com/stretchr/testify/require"
)

func TestAPI_Balance_OK(t *testing.T) {
	tree := abstate.NewEmptyState()
	stateDB := statedb.NewStateDB(tree)
	address := common.BytesToAddress(test.RandomBytes(20))

	stateDB.CreateAccount(address)
	balance := big.NewInt(101)
	stateDB.AddBalance(address, balance)
	_, _, err := tree.CalculateRoot()
	require.NoError(t, err)
	require.NoError(t, tree.Commit())

	a := &API{
		state:            tree,
		systemIdentifier: []byte{0, 0, 0, 1},
		gasUnitPrice:     big.NewInt(10),
		blockGasLimit:    10000,
	}

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/evm/balance/%X", address.Bytes()), nil)
	recorder := httptest.NewRecorder()

	rpc.NewRESTServer("", 2000, a).Handler.ServeHTTP(recorder, req)
	require.Equal(t, http.StatusOK, recorder.Code)
	resp := &struct {
		_       struct{} `cbor:",toarray"`
		Balance string
	}{}
	require.NoError(t, cbor.NewDecoder(recorder.Body).Decode(resp))
	require.Equal(t, balance.String(), resp.Balance)
}

func TestAPI_Balance_NotFound(t *testing.T) {
	a := &API{
		state:            abstate.NewEmptyState(),
		systemIdentifier: []byte{0, 0, 0, 1},
		gasUnitPrice:     big.NewInt(10),
		blockGasLimit:    10000,
	}

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/evm/balance/%X", test.RandomBytes(20)), nil)
	recorder := httptest.NewRecorder()

	rpc.NewRESTServer("", 2000, a).Handler.ServeHTTP(recorder, req)
	require.Equal(t, http.StatusNotFound, recorder.Code)
}
