package service

import (
	"errors"
	"quote-api/internal/model"
	"quote-api/internal/store"
)

type QuoteService interface {
	AddQuote(author, text string) (model.Quote, error)
	ListQuotes(authorFilter string) []model.Quote
	RandomQuote() (model.Quote, error)
	DeleteQuote(id int) error
}

type service struct {
	store store.QuoteStore
}

func NewQuoteService(st store.QuoteStore) QuoteService {
	return &service{store: st}
}

func (s *service) AddQuote(author, text string) (model.Quote, error) {
	if author == "" || text == "" {
		return model.Quote{}, errors.New("author and quote text must be provided")
	}
	return s.store.Add(model.Quote{Author: author, Text: text}), nil
}

func (s *service) ListQuotes(authorFilter string) []model.Quote {
	if authorFilter != "" {
		return s.store.GetByAuthor(authorFilter)
	}
	return s.store.GetAll()
}

func (s *service) RandomQuote() (model.Quote, error) {
	return s.store.GetRandom()
}

func (s *service) DeleteQuote(id int) error {
	if ok := s.store.Delete(id); !ok {
		return errors.New("quote not found")
	}
	return nil
}
