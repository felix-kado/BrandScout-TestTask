package handler

import (
	"github.com/gorilla/mux"
	"quote-api/internal/middleware"
	"quote-api/internal/service"
)

func NewRouter(quoteService service.QuoteService) *mux.Router {
	h := NewHandler(quoteService)
	r := mux.NewRouter()
	r.Use(middleware.LoggingMW)
	r.HandleFunc("/quotes", h.CreateQuote).Methods("POST")
	r.HandleFunc("/quotes", h.GetAllQuotes).Methods("GET")
	r.HandleFunc("/quotes/random", h.GetRandomQuote).Methods("GET")
	r.HandleFunc("/quotes/{id}", h.DeleteQuote).Methods("DELETE")
	return r
}
