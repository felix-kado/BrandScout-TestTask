package handler

import (
	"encoding/json"
	"net/http"
	"quote-api/internal/model"
	"quote-api/internal/render"
	"quote-api/internal/service"
	"strconv"

	"github.com/gorilla/mux"
)

type Handler struct {
	svc service.QuoteService
}

func New(svc service.QuoteService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) CreateQuote(w http.ResponseWriter, r *http.Request) {
	var req model.Quote
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid input"})
		return
	}

	q, err := h.svc.AddQuote(req.Author, req.Text)
	if err != nil {
		render.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	render.JSON(w, http.StatusCreated, q)
}

func (h *Handler) GetAllQuotes(w http.ResponseWriter, r *http.Request) {
	author := r.URL.Query().Get("author")
	quotes := h.svc.ListQuotes(author)
	render.JSON(w, http.StatusOK, quotes)
}

func (h *Handler) GetRandomQuote(w http.ResponseWriter, r *http.Request) {
	q, err := h.svc.RandomQuote()
	if err != nil {
		render.JSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	render.JSON(w, http.StatusOK, q)
}

func (h *Handler) DeleteQuote(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		render.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid ID"})
		return
	}
	if err = h.svc.DeleteQuote(id); err != nil {
		render.JSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
