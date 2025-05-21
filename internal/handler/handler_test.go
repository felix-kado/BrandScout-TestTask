package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"quote-api/internal/model"
	"quote-api/internal/store"
)

func setupTestHandler() *Handler {
	return New(store.NewInMemoryStore())
}

func TestCreateQuoteHandler(t *testing.T) {
	h := setupTestHandler()
	r := mux.NewRouter()
	r.HandleFunc("/quotes", h.CreateQuote).Methods("POST")

	quote := model.Quote{Author: "H", Text: "Hello"}
	body, _ := json.Marshal(quote)

	req := httptest.NewRequest(http.MethodPost, "/quotes", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusCreated, resp.Code)
}

func TestDeleteQuoteHandler_NotFound(t *testing.T) {
	h := setupTestHandler()
	r := mux.NewRouter()
	r.HandleFunc("/quotes/{id}", h.DeleteQuote).Methods("DELETE")

	req := httptest.NewRequest(http.MethodDelete, "/quotes/999", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusNotFound, resp.Code)
}
