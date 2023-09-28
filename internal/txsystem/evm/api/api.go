package api

import (
	"math/big"
	"net/http"

	"github.com/alphabill-org/alphabill/internal/types"

	"github.com/alphabill-org/alphabill/internal/state"
	"github.com/gorilla/mux"
)

type API struct {
	state            *state.State
	systemIdentifier types.SystemID
	gasUnitPrice     *big.Int
	blockGasLimit    uint64
}

func NewAPI(s *state.State, systemIdentifier types.SystemID, gasUnitPrice *big.Int, blockGasLimit uint64) *API {
	return &API{
		state:            s,
		systemIdentifier: systemIdentifier,
		gasUnitPrice:     gasUnitPrice,
		blockGasLimit:    blockGasLimit,
	}
}

func (a *API) Register(r *mux.Router) {
	evmRouter := r.PathPrefix("/evm").Subrouter()
	evmRouter.HandleFunc("/call", a.CallEVM).Methods(http.MethodPost, http.MethodOptions)
	evmRouter.HandleFunc("/estimateGas", a.EstimateGas).Methods(http.MethodPost, http.MethodOptions)
	evmRouter.HandleFunc("/balance/{address}", a.Balance).Methods(http.MethodGet, http.MethodOptions)
	evmRouter.HandleFunc("/transactionCount/{address}", a.TransactionCount).Methods(http.MethodGet, http.MethodOptions)
}
