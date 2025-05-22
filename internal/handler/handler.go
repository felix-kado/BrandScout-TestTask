package handler

import (
	"encoding/json"
	"net/http"
	"quote-api/internal/model"
	"strconv"

	"github.com/gorilla/mux"
	"quote-api/internal/service"
)

type Handler struct {
	svc service.QuoteService
}

func New(svc service.QuoteService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) CreateQuote(w http.ResponseWriter, r *http.Request) {
	var q model.Quote
	if err := json.NewDecoder(r.Body).Decode(&q); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}
	q, err := h.svc.AddQuote(q.Author, q.Text)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(q); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func (h *Handler) GetAllQuotes(w http.ResponseWriter, r *http.Request) {
	author := r.URL.Query().Get("author")
	quotes := h.svc.ListQuotes(author)
	if err := json.NewEncoder(w).Encode(quotes); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func (h *Handler) GetRandomQuote(w http.ResponseWriter, r *http.Request) {
	q, err := h.svc.RandomQuote()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err = json.NewEncoder(w).Encode(q); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func (h *Handler) DeleteQuote(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid ID", http.StatusBadRequest)
		return
	}
	if err = h.svc.DeleteQuote(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
