package service

import (
	"context"
	"quote-api/internal/model"
	"quote-api/internal/store"
)

type QuoteService interface {
	AddQuote(ctx context.Context, author, text string) (model.Quote, error)
	ListQuotes(ctx context.Context, authorFilter string) []model.Quote
	RandomQuote(ctx context.Context) (model.Quote, error)
	DeleteQuote(ctx context.Context, id int) error
}

type service struct {
	store store.QuoteStore
}

func NewQuoteService(st store.QuoteStore) QuoteService {
	return &service{store: st}
}

func (s *service) AddQuote(ctx context.Context, author, text string) (model.Quote, error) {
	if author == "" {
		return model.Quote{}, model.ErrEmptyQuoteAuthor
	}
	if text == "" {
		return model.Quote{}, model.ErrEmptyQuoteText
	}

	q := s.store.Add(ctx, model.Quote{Author: author, Text: text})
	return q, nil
}

func (s *service) ListQuotes(ctx context.Context, authorFilter string) []model.Quote {
	if authorFilter != "" {
		return s.store.GetByAuthor(ctx, authorFilter)
	}
	return s.store.GetAll(ctx)
}

func (s *service) RandomQuote(ctx context.Context) (model.Quote, error) {
	return s.store.GetRandom(ctx)
}

func (s *service) DeleteQuote(ctx context.Context, id int) error {
	if ok := s.store.Delete(ctx, id); !ok {
		return model.ErrQuoteNotFound
	}
	return nil
}
