package twb

import (
	"bytes"
	"context"
	"embed"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"golang.org/x/sync/semaphore"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/alphabill-org/alphabill/internal/hash"
	"github.com/alphabill-org/alphabill/internal/script"
	"github.com/alphabill-org/alphabill/internal/txsystem"
	"github.com/alphabill-org/alphabill/internal/txsystem/tokens"
)

type dataSource interface {
	GetBlockNumber() (uint64, error)
	GetTokenType(id TokenTypeID) (*TokenUnitType, error)
	QueryTokenType(kind Kind, creator PubKey, startKey TokenTypeID, count int) ([]*TokenUnitType, TokenTypeID, error)
	GetToken(id TokenID) (*TokenUnit, error)
	QueryTokens(kind Kind, owner Predicate, startKey TokenID, count int) ([]*TokenUnit, TokenID, error)
	SaveTokenTypeCreator(id TokenTypeID, kind Kind, creator PubKey) error
	GetTxProof(unitID UnitID, txHash TxHash) (*Proof, error)
}

type restAPI struct {
	db              dataSource
	sendTransaction func(context.Context, *txsystem.Transaction) (*txsystem.TransactionResponse, error)
	convertTx       func(tx *txsystem.Transaction) (txsystem.GenericTransaction, error)
	logErr          func(a ...any)
}

const maxResponseItems = 100

//go:embed swagger/*
var swaggerFiles embed.FS

func (api *restAPI) endpoints() http.Handler {
	apiRouter := mux.NewRouter().StrictSlash(true).PathPrefix("/api").Subrouter()

	// add cors middleware
	// content-type needs to be explicitly defined without this content-type header is not allowed and cors filter is not applied
	// Link header is needed for pagination support.
	// OPTIONS method needs to be explicitly defined for each handler func
	apiRouter.Use(handlers.CORS(
		handlers.AllowedHeaders([]string{"Content-Type"}),
		handlers.ExposedHeaders([]string{"Link"}),
	))

	// version v1 router
	apiV1 := apiRouter.PathPrefix("/v1").Subrouter()
	apiV1.HandleFunc("/tokens/{tokenId}", api.getToken).Methods("GET", "OPTIONS")
	apiV1.HandleFunc("/types/{typeId}/hierarchy", api.typeHierarchy).Methods("GET", "OPTIONS")
	apiV1.HandleFunc("/kinds/{kind}/owners/{owner}/tokens", api.listTokens).Methods("GET", "OPTIONS")
	apiV1.HandleFunc("/kinds/{kind}/types", api.listTypes).Methods("GET", "OPTIONS")
	apiV1.HandleFunc("/round-number", api.getRoundNumber).Methods("GET", "OPTIONS")
	apiV1.HandleFunc("/transactions/{pubkey}", api.postTransactions).Methods("POST", "OPTIONS")
	apiV1.HandleFunc("/units/{unitId}/transactions/{txHash}/proof", api.getTxProof).Methods("GET", "OPTIONS")

	apiV1.Handle("/swagger/{.*}", http.StripPrefix("/api/v1/", http.FileServer(http.FS(swaggerFiles)))).Methods("GET", "OPTIONS")
	apiV1.Handle("/swagger/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f, err := swaggerFiles.ReadFile("swagger/index.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "failed to read swagger/index.html file: %v", err)
			return
		}
		http.ServeContent(w, r, "index.html", time.Time{}, bytes.NewReader(f))
	})).Methods("GET", "OPTIONS")

	return apiRouter
}

func (api *restAPI) listTokens(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	owner, err := parsePubKey(vars["owner"], true)
	if err != nil {
		api.invalidParamResponse(w, "owner", err)
		return
	}

	kind, err := strToTokenKind(vars["kind"])
	if err != nil {
		api.invalidParamResponse(w, "kind", err)
		return
	}

	qp := r.URL.Query()
	startKey, err := parseHex[TokenID](qp.Get("offsetKey"), false)
	if err != nil {
		api.invalidParamResponse(w, "offsetKey", err)
		return
	}

	limit, err := parseMaxResponseItems(qp.Get("limit"), maxResponseItems)
	if err != nil {
		api.invalidParamResponse(w, "limit", err)
		return
	}

	data, next, err := api.db.QueryTokens(
		kind,
		script.PredicatePayToPublicKeyHashDefault(hash.Sum256(owner)),
		startKey,
		limit)
	if err != nil {
		api.writeErrorResponse(w, err)
		return
	}
	setLinkHeader(r.URL, w, encodeHex(next))
	api.writeResponse(w, data)
}

func (api *restAPI) getToken(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenId, err := parseHex[TokenID](vars["tokenId"], true)
	if err != nil {
		api.invalidParamResponse(w, "tokenId", err)
		return
	}

	token, err := api.db.GetToken(tokenId)
	if err != nil {
		api.writeErrorResponse(w, err)
		return
	}
	api.writeResponse(w, token)
}

func (api *restAPI) listTypes(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	kind, err := strToTokenKind(vars["kind"])
	if err != nil {
		api.invalidParamResponse(w, "kind", err)
		return
	}

	qp := r.URL.Query()
	creator, err := parsePubKey(qp.Get("creator"), false)
	if err != nil {
		api.invalidParamResponse(w, "creator", err)
		return
	}

	startKey, err := parseHex[TokenTypeID](qp.Get("offsetKey"), false)
	if err != nil {
		api.invalidParamResponse(w, "offsetKey", err)
		return
	}

	limit, err := parseMaxResponseItems(qp.Get("limit"), maxResponseItems)
	if err != nil {
		api.invalidParamResponse(w, "limit", err)
		return
	}

	data, next, err := api.db.QueryTokenType(kind, creator, startKey, limit)
	if err != nil {
		api.writeErrorResponse(w, err)
		return
	}
	setLinkHeader(r.URL, w, encodeHex(next))
	api.writeResponse(w, data)
}

func (api *restAPI) typeHierarchy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	typeId, err := parseHex[TokenTypeID](vars["typeId"], true)
	if err != nil {
		api.invalidParamResponse(w, "typeId", err)
		return
	}

	var rsp []*TokenUnitType
	for len(typeId) > 0 && !bytes.Equal(typeId, NoParent) {
		tokTyp, err := api.db.GetTokenType(typeId)
		if err != nil {
			api.writeErrorResponse(w, fmt.Errorf("failed to load type with id %x: %w", typeId, err))
			return
		}
		rsp = append(rsp, tokTyp)
		typeId = tokTyp.ParentTypeID
	}
	api.writeResponse(w, rsp)
}

func (api *restAPI) getRoundNumber(w http.ResponseWriter, r *http.Request) {
	rn, err := api.db.GetBlockNumber()
	if err != nil {
		api.writeErrorResponse(w, err)
		return
	}
	api.writeResponse(w, RoundNumberResponse{RoundNumber: rn})
}

func (api *restAPI) postTransactions(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	buf, err := io.ReadAll(r.Body)
	if err != nil {
		api.writeErrorResponse(w, fmt.Errorf("failed to read request body: %w", err))
		return
	}

	vars := mux.Vars(r)
	owner, err := parsePubKey(vars["pubkey"], true)
	if err != nil {
		api.invalidParamResponse(w, "pubkey", err)
		return
	}

	txs := &txsystem.Transactions{}
	if err = protojson.Unmarshal(buf, txs); err != nil {
		api.errorResponse(w, http.StatusBadRequest, fmt.Errorf("failed to decode request body: %w", err))
		return
	}
	if len(txs.GetTransactions()) == 0 {
		api.errorResponse(w, http.StatusBadRequest, fmt.Errorf("request body contained no transactions to process"))
		return
	}

	if errs := api.saveTxs(r.Context(), txs.GetTransactions(), owner); len(errs) > 0 {
		w.WriteHeader(http.StatusInternalServerError)
		api.writeResponse(w, errs)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (api *restAPI) saveTxs(ctx context.Context, txs []*txsystem.Transaction, owner []byte) map[string]string {
	errs := make(map[string]string)
	var m sync.Mutex

	const maxWorkers = 5
	sem := semaphore.NewWeighted(maxWorkers)
	for _, tx := range txs {
		if err := sem.Acquire(ctx, 1); err != nil {
			break
		}
		go func(tx *txsystem.Transaction) {
			defer sem.Release(1)
			if err := api.saveTx(ctx, tx, owner); err != nil {
				m.Lock()
				errs[hex.EncodeToString(tx.GetUnitId())] = err.Error()
				m.Unlock()
			}
		}(tx)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := sem.Acquire(ctx, maxWorkers); err != nil {
		m.Lock()
		errs["waiting-for-workers"] = err.Error()
		m.Unlock()
	}
	return errs
}

func (api *restAPI) saveTx(ctx context.Context, tx *txsystem.Transaction, owner []byte) error {
	// if "creator type tx" then save the type->owner relation
	gtx, err := api.convertTx(tx)
	if err != nil {
		return fmt.Errorf("failed to convert transaction: %w", err)
	}
	kind := Any
	switch gtx.(type) {
	case tokens.CreateFungibleTokenType:
		kind = Fungible
	case tokens.CreateNonFungibleTokenType:
		kind = NonFungible
	}
	if kind != Any {
		if err := api.db.SaveTokenTypeCreator(tx.UnitId, kind, owner); err != nil {
			return fmt.Errorf("failed to save creator relation: %w", err)
		}
	}

	rsp, err := api.sendTransaction(ctx, tx)
	if err != nil {
		return fmt.Errorf("failed to forward tx: %w", err)
	}
	if !rsp.GetOk() {
		return fmt.Errorf("transaction was not accepted: %s", rsp.GetMessage())
	}
	return nil
}

func (api *restAPI) getTxProof(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	unitID, err := parseHex[UnitID](vars["unitId"], true)
	if err != nil {
		api.invalidParamResponse(w, "unitId", err)
		return
	}
	txHash, err := parseHex[TxHash](vars["txHash"], true)
	if err != nil {
		api.invalidParamResponse(w, "txHash", err)
		return
	}

	proof, err := api.db.GetTxProof(unitID, txHash)
	if err != nil {
		api.writeErrorResponse(w, fmt.Errorf("failed to load proof of tx 0x%X (unit 0x%X): %w", txHash, unitID, err))
		return
	}
	if proof == nil {
		api.errorResponse(w, http.StatusNotFound, fmt.Errorf("no proof found for tx 0x%X (unit 0x%X)", txHash, unitID))
		return
	}

	api.writeResponse(w, proof)
}

func (api *restAPI) writeResponse(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		api.logError(fmt.Errorf("failed to encode response data as json: %w", err))
	}
}

func (api *restAPI) writeErrorResponse(w http.ResponseWriter, err error) {
	if errors.Is(err, errRecordNotFound) {
		api.errorResponse(w, http.StatusNotFound, err)
		return
	}

	api.errorResponse(w, http.StatusInternalServerError, err)
	api.logError(err)
}

func (api *restAPI) invalidParamResponse(w http.ResponseWriter, name string, err error) {
	api.errorResponse(w, http.StatusBadRequest, fmt.Errorf("invalid parameter %q: %w", name, err))
}

func (api *restAPI) errorResponse(w http.ResponseWriter, code int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(ErrorResponse{Message: err.Error()}); err != nil {
		api.logError(fmt.Errorf("failed to encode error response as json: %w", err))
	}
}

func (api *restAPI) logError(err error) {
	if api.logErr != nil {
		api.logErr(err)
	}
}

func setLinkHeader(u *url.URL, w http.ResponseWriter, next string) {
	if next == "" {
		w.Header().Del("Link")
		return
	}
	qp := u.Query()
	qp.Set("offsetKey", next)
	u.RawQuery = qp.Encode()
	w.Header().Set("Link", fmt.Sprintf(`<%s>; rel="next"`, u))
}

type (
	RoundNumberResponse struct {
		RoundNumber uint64 `json:"roundNumber,string"`
	}

	ErrorResponse struct {
		Message string `json:"message"`
	}
)