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
	"quote-api/internal/service"
	"quote-api/internal/store"
)

func setupRouter() *mux.Router {
	svc := service.NewQuoteService(store.NewInMemoryStore())
	h := handler.New(svc)
	r := mux.NewRouter()
	r.HandleFunc("/quotes", h.CreateQuote).Methods("POST")
	r.HandleFunc("/quotes", h.GetAllQuotes).Methods("GET")
	r.HandleFunc("/quotes/random", h.GetRandomQuote).Methods("GET")
	r.HandleFunc("/quotes/{id}", h.DeleteQuote).Methods("DELETE")
	return r
}

func TestAPI_FullFlow(t *testing.T) {
	r := setupRouter()

	quotes := []map[string]string{
		{"author": "Alice", "quote": "First"},
		{"author": "Bob", "quote": "Second"},
		{"author": "Alice", "quote": "Third"},
	}

	var created []map[string]interface{}
	for _, q := range quotes {
		body, _ := json.Marshal(q)
		req := httptest.NewRequest("POST", "/quotes", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusCreated, resp.Code)
		var result map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&result)
		assert.NoError(t, err)
		created = append(created, result)
	}

	// Список всех цитат
	req := httptest.NewRequest("GET", "/quotes", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	var all []map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&all)
	assert.NoError(t, err)

	assert.Len(t, all, 3)

	// Фильтр по автору
	req = httptest.NewRequest("GET", "/quotes?author=Alice", nil)
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	var filtered []map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&filtered)
	assert.NoError(t, err)

	assert.Len(t, filtered, 2)

	// Случайная цитата
	req = httptest.NewRequest("GET", "/quotes/random", nil)
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	var random map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&random)
	assert.NoError(t, err)
	assert.NotEmpty(t, random["quote"])

	// Удалить цитату
	id := int(created[0]["id"].(float64))
	delReq := httptest.NewRequest("DELETE", "/quotes/"+strconv.Itoa(id), nil)
	delResp := httptest.NewRecorder()
	r.ServeHTTP(delResp, delReq)
	assert.Equal(t, http.StatusNoContent, delResp.Code)

	// Удалить её же опять -> 404
	delReq2 := httptest.NewRequest("DELETE", "/quotes/"+strconv.Itoa(id), nil)
	delResp2 := httptest.NewRecorder()
	r.ServeHTTP(delResp2, delReq2)
	assert.Equal(t, http.StatusNotFound, delResp2.Code)
}

func TestAPI_InvalidInput(t *testing.T) {
	r := setupRouter()

	// --- Invalid POST (missing body)
	req := httptest.NewRequest("POST", "/quotes", nil)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusBadRequest, resp.Code)

	// --- Invalid DELETE ID
	req = httptest.NewRequest("DELETE", "/quotes/abc", nil)
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusBadRequest, resp.Code)

	// --- Рандомная цитата из пустой памяти
	rEmpty := setupRouter()
	req = httptest.NewRequest("GET", "/quotes/random", nil)
	resp = httptest.NewRecorder()
	rEmpty.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusNotFound, resp.Code)
}
