package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"quote-api/internal/handler"
	"quote-api/internal/model"
	"quote-api/internal/store"
)

func setupRouter() *mux.Router {
	s := store.NewInMemoryStore()
	h := handler.New(s)

	r := mux.NewRouter()
	r.HandleFunc("/quotes", h.CreateQuote).Methods("POST")
	r.HandleFunc("/quotes", h.GetAllQuotes).Methods("GET")

	r.HandleFunc("/quotes/random", h.GetRandomQuote).Methods("GET")
	r.HandleFunc("/quotes/{id}", h.DeleteQuote).Methods("DELETE")

	return r
}

func TestAPI_FullCycle(t *testing.T) {
	r := setupRouter()

	quote := model.QuoteNote{Author: "Test", Quote: "Test quote"}
	body, _ := json.Marshal(quote)

	req := httptest.NewRequest(http.MethodPost, "/quotes", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusCreated, resp.Code)

	var created model.QuoteNote
	_ = json.NewDecoder(resp.Body).Decode(&created)
	assert.Equal(t, quote.Author, created.Author)

	reqGet := httptest.NewRequest(http.MethodGet, "/quotes", nil)
	respGet := httptest.NewRecorder()
	r.ServeHTTP(respGet, reqGet)
	assert.Equal(t, http.StatusOK, respGet.Code)

	reqDel := httptest.NewRequest(http.MethodDelete, "/quotes/"+strconv.Itoa(created.ID), nil)
	respDel := httptest.NewRecorder()
	r.ServeHTTP(respDel, reqDel)
	assert.Equal(t, http.StatusNoContent, respDel.Code)
}
