package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"quote-api/internal/service"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"quote-api/internal/handler"
	"quote-api/internal/middleware"
	"quote-api/internal/store"
)

func main() {
	st := store.NewInMemoryStore()
	svc := service.NewQuoteService(st)
	h := handler.New(svc)

	r := mux.NewRouter()
	r.Use(middleware.LoggingMW)
	r.HandleFunc("/quotes", h.CreateQuote).Methods("POST")
	r.HandleFunc("/quotes", h.GetAllQuotes).Methods("GET")
	r.HandleFunc("/quotes/random", h.GetRandomQuote).Methods("GET")
	r.HandleFunc("/quotes/{id}", h.DeleteQuote).Methods("DELETE")

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		log.Println("🚀 Server started on :8080")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown: %v", err)
	}
	log.Println("Server gracefully stopped")

}
