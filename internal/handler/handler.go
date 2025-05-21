package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"quote-api/internal/model"
	"quote-api/internal/store"
	"strconv"
)

type Handler struct {
	store store.QuoteStore
}

func New(s store.QuoteStore) *Handler {
	return &Handler{store: s}
}

func (h *Handler) CreateQuote(w http.ResponseWriter, r *http.Request) {
	var q model.Quote
	if err := json.NewDecoder(r.Body).Decode(&q); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}
	added := h.store.Add(q)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(added)
}

func (h *Handler) GetAllQuotes(w http.ResponseWriter, r *http.Request) {
	author := r.URL.Query().Get("author")
	if author != "" {
		json.NewEncoder(w).Encode(h.store.GetByAuthor(author))
		return
	}
	json.NewEncoder(w).Encode(h.store.GetAll())
}

func (h *Handler) GetRandomQuote(w http.ResponseWriter, r *http.Request) {
	q, err := h.store.GetRandom()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(q)
}

func (h *Handler) DeleteQuote(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid ID", http.StatusBadRequest)
		return
	}
	if !h.store.Delete(id) {
		http.Error(w, "quote not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
