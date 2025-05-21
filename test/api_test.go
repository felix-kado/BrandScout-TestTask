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

func TestAPI_CreateAndRetrieveQuote(t *testing.T) {
	r := setupRouter()

	quote := model.Quote{Author: "Test", Text: "Test quote"}
	body, _ := json.Marshal(quote)

	req := httptest.NewRequest(http.MethodPost, "/quotes", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusCreated, resp.Code)

	var created model.Quote
	_ = json.NewDecoder(resp.Body).Decode(&created)
	assert.Equal(t, quote.Author, created.Author)
	assert.NotZero(t, created.ID)
}

func TestAPI_FilterByAuthor(t *testing.T) {
	r := setupRouter()

	quotes := []model.Quote{
		{Author: "A", Text: "One"},
		{Author: "B", Text: "Two"},
		{Author: "A", Text: "Three"},
	}

	for _, q := range quotes {
		body, _ := json.Marshal(q)
		req := httptest.NewRequest(http.MethodPost, "/quotes", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(httptest.NewRecorder(), req)
	}

	req := httptest.NewRequest(http.MethodGet, "/quotes?author=A", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	var result []model.Quote
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.Len(t, result, 2)
}

func TestAPI_RandomQuote(t *testing.T) {
	r := setupRouter()
	quote := model.Quote{Author: "Rand", Text: "Surprise"}
	body, _ := json.Marshal(quote)
	postReq := httptest.NewRequest(http.MethodPost, "/quotes", bytes.NewReader(body))
	postReq.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(httptest.NewRecorder(), postReq)

	randomReq := httptest.NewRequest(http.MethodGet, "/quotes/random", nil)
	randomResp := httptest.NewRecorder()
	r.ServeHTTP(randomResp, randomReq)
	assert.Equal(t, http.StatusOK, randomResp.Code)

	var got model.Quote
	_ = json.NewDecoder(randomResp.Body).Decode(&got)
	assert.Equal(t, "Rand", got.Author)
}

func TestAPI_DeleteQuote(t *testing.T) {
	r := setupRouter()
	quote := model.Quote{Author: "Del", Text: "To be deleted"}
	body, _ := json.Marshal(quote)
	postReq := httptest.NewRequest(http.MethodPost, "/quotes", bytes.NewReader(body))
	postReq.Header.Set("Content-Type", "application/json")
	postResp := httptest.NewRecorder()
	r.ServeHTTP(postResp, postReq)

	var created model.Quote
	_ = json.NewDecoder(postResp.Body).Decode(&created)

	delReq := httptest.NewRequest(http.MethodDelete, "/quotes/"+strconv.Itoa(created.ID), nil)
	delResp := httptest.NewRecorder()
	r.ServeHTTP(delResp, delReq)
	assert.Equal(t, http.StatusNoContent, delResp.Code)

	delReq2 := httptest.NewRequest(http.MethodDelete, "/quotes/"+strconv.Itoa(created.ID), nil)
	delResp2 := httptest.NewRecorder()
	r.ServeHTTP(delResp2, delReq2)
	assert.Equal(t, http.StatusNotFound, delResp2.Code)
}
