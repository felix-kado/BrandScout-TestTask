package store

import (
	"context"
	"quote-api/internal/model"
)

type QuoteStore interface {
	Add(ctx context.Context, q model.Quote) model.Quote
	GetAll(ctx context.Context) []model.Quote
	GetByAuthor(ctx context.Context, author string) []model.Quote
	GetRandom(ctx context.Context) (model.Quote, error)
	Delete(ctx context.Context, id int) bool
}
