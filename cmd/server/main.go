package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"quote-api/internal/handler"
	"quote-api/internal/middleware"
	"quote-api/internal/store"
)

func main() {
	store := store.NewInMemoryStore()
	h := handler.New(store)

	r := mux.NewRouter()
	r.Use(middleware.LoggingMW)

	r.HandleFunc("/quotes", h.CreateQuote).Methods("POST")
	r.HandleFunc("/quotes", h.GetAllQuotes).Methods("GET")
	r.HandleFunc("/quotes/random", h.GetRandomQuote).Methods("GET")
	r.HandleFunc("/quotes/{id}", h.DeleteQuote).Methods("DELETE")

	log.Println("🚀 Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
