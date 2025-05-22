package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"quote-api/internal/model"
	"quote-api/internal/service"
	"quote-api/internal/store"
)

func setupHandlerTest() (*Handler, *mux.Router) {
	svc := service.NewQuoteService(store.NewInMemoryStore())
	h := New(svc)
	r := mux.NewRouter()
	r.HandleFunc("/quotes", h.CreateQuote).Methods("POST")
	r.HandleFunc("/quotes", h.GetAllQuotes).Methods("GET")
	r.HandleFunc("/quotes/random", h.GetRandomQuote).Methods("GET")
	r.HandleFunc("/quotes/{id}", h.DeleteQuote).Methods("DELETE")
	return h, r
}

func TestCreateAndGet(t *testing.T) {
	_, r := setupHandlerTest()

	q := map[string]string{"author": "Tester", "quote": "Hello World"}
	body, _ := json.Marshal(q)
	req := httptest.NewRequest("POST", "/quotes", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var created model.Quote
	err := json.NewDecoder(w.Body).Decode(&created)
	assert.NoError(t, err)

	assert.Equal(t, "Tester", created.Author)

	getReq := httptest.NewRequest("GET", "/quotes", nil)
	getW := httptest.NewRecorder()
	r.ServeHTTP(getW, getReq)
	assert.Equal(t, http.StatusOK, getW.Code)

	var result []model.Quote
	err = json.NewDecoder(getW.Body).Decode(&result)
	assert.NoError(t, err)

	assert.GreaterOrEqual(t, len(result), 1)
}

func TestDeleteHandler(t *testing.T) {
	_, r := setupHandlerTest()

	q := map[string]string{"author": "Del", "quote": "To remove"}
	body, _ := json.Marshal(q)
	req := httptest.NewRequest("POST", "/quotes", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var created model.Quote
	err := json.NewDecoder(w.Body).Decode(&created)
	assert.NoError(t, err)

	delReq := httptest.NewRequest("DELETE", "/quotes/"+strconv.Itoa(created.ID), nil)
	delW := httptest.NewRecorder()
	r.ServeHTTP(delW, delReq)
	assert.Equal(t, http.StatusNoContent, delW.Code)

	delReq2 := httptest.NewRequest("DELETE", "/quotes/"+strconv.Itoa(created.ID), nil)
	delW2 := httptest.NewRecorder()
	r.ServeHTTP(delW2, delReq2)
	assert.Equal(t, http.StatusNotFound, delW2.Code)
}
