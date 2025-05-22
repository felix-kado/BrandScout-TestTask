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

func NewHandler(svc service.QuoteService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) CreateQuote(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req model.Quote
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse(w, err)
		return
	}

	q, err := h.svc.AddQuote(ctx, req.Author, req.Text)
	if err != nil {
		errorResponse(w, err)
		return
	}

	render.JSON(w, http.StatusCreated, q)
}

func (h *Handler) GetAllQuotes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	author := r.URL.Query().Get("author")
	quotes := h.svc.ListQuotes(ctx, author)
	render.JSON(w, http.StatusOK, quotes)
}

func (h *Handler) GetRandomQuote(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	q, err := h.svc.RandomQuote(ctx)
	if err != nil {
		errorResponse(w, err)
		return
	}
	render.JSON(w, http.StatusOK, q)
}

func (h *Handler) DeleteQuote(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		errorResponse(w, model.ErrInvalidQuoteID)
		return
	}
	if err = h.svc.DeleteQuote(ctx, id); err != nil {
		errorResponse(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
