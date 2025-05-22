package handler

import (
	"github.com/gorilla/mux"
	"net/http"
	"quote-api/internal/middleware"
	"quote-api/internal/service"
	"quote-api/internal/store"
)

func newRouter(h *Handler) *mux.Router {
	r := mux.NewRouter()
	r.Use(middleware.LoggingMW)
	r.HandleFunc("/quotes", h.CreateQuote).Methods("POST")
	r.HandleFunc("/quotes", h.GetAllQuotes).Methods("GET")
	r.HandleFunc("/quotes/random", h.GetRandomQuote).Methods("GET")
	r.HandleFunc("/quotes/{id}", h.DeleteQuote).Methods("DELETE")
	return r
}

func Router() http.Handler {
	st := store.NewInMemoryStore()
	svc := service.NewQuoteService(st)
	h := New(svc)
	return newRouter(h)
}
