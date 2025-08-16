package http

import (
	"encoding/json"
	stdhttp "net/http"

	"github.com/go-chi/chi/v5"

	"github.com/rickyreddygari/walletsdk/internal/service"
)

type RouteBuilder struct {
	wallets  service.WalletService
	balances service.BalanceService
}

func NewRouteBuilder(wallets service.WalletService, balances service.BalanceService) *RouteBuilder {
	return &RouteBuilder{wallets: wallets, balances: balances}
}

func (b *RouteBuilder) Register(r *chi.Mux) {
	r.Route("/v1", func(r chi.Router) {
		r.Post("/wallets", b.createWallet)
		r.Get("/wallets/{id}", b.getWallet)
		r.Get("/wallets", b.listWallets)
		r.Post("/wallets/{id}/sign-message", b.signMessage)
		r.Post("/wallets/{id}/sign-transaction", b.signTransaction)
		r.Get("/wallets/{id}/balance", b.getBalance)
	})
}

func (b *RouteBuilder) createWallet(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	var payload struct {
		Network string `json:"network"`
	}

	if err := decodeJSON(r, &payload); err != nil {
		writeError(w, stdhttp.StatusBadRequest, "invalid request body")
		return
	}

	wallet, err := b.wallets.CreateWallet(r.Context(), payload.Network)
	if err != nil {
		writeError(w, stdhttp.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, stdhttp.StatusCreated, wallet)
}

func (b *RouteBuilder) getWallet(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		writeError(w, stdhttp.StatusBadRequest, "wallet id is required")
		return
	}

	wallet, err := b.wallets.GetWallet(r.Context(), id)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	writeJSON(w, stdhttp.StatusOK, wallet)
}

func (b *RouteBuilder) listWallets(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	network := r.URL.Query().Get("network")
	wallets, err := b.wallets.ListWallets(r.Context(), network)
	if err != nil {
		writeError(w, stdhttp.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, stdhttp.StatusOK, wallets)
}

func (b *RouteBuilder) signMessage(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		writeError(w, stdhttp.StatusBadRequest, "wallet id is required")
		return
	}
	var payload struct {
		Message string `json:"message"`
	}

	if err := decodeJSON(r, &payload); err != nil {
		writeError(w, stdhttp.StatusBadRequest, "invalid request body")
		return
	}

	result, err := b.wallets.SignMessage(r.Context(), id, []byte(payload.Message))
	if err != nil {
		handleServiceError(w, err)
		return
	}

	writeJSON(w, stdhttp.StatusOK, result)
}

func (b *RouteBuilder) signTransaction(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		writeError(w, stdhttp.StatusBadRequest, "wallet id is required")
		return
	}
	var payload service.Transaction

	if err := decodeJSON(r, &payload); err != nil {
		writeError(w, stdhttp.StatusBadRequest, "invalid request body")
		return
	}

	signed, err := b.wallets.SignTransaction(r.Context(), id, &payload)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	writeJSON(w, stdhttp.StatusOK, map[string]string{"signedTransaction": signed})
}

func (b *RouteBuilder) getBalance(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		writeError(w, stdhttp.StatusBadRequest, "wallet id is required")
		return
	}
	balance, err := b.balances.GetBalance(r.Context(), id)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	writeJSON(w, stdhttp.StatusOK, balance)
}

func handleServiceError(w stdhttp.ResponseWriter, err error) {
	switch err {
	case service.ErrNotFound:
		writeError(w, stdhttp.StatusNotFound, "resource not found")
	case service.ErrValidation:
		writeError(w, stdhttp.StatusBadRequest, err.Error())
	case service.ErrNotImplemented:
		writeError(w, stdhttp.StatusNotImplemented, err.Error())
	default:
		writeError(w, stdhttp.StatusInternalServerError, err.Error())
	}
}

func writeJSON(w stdhttp.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func writeError(w stdhttp.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func decodeJSON(r *stdhttp.Request, v interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(v)
}
