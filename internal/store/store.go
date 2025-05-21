package store

import "quote-api/internal/model"

type QuoteStore interface {
	Add(model.Quote) model.Quote
	GetAll() []model.Quote
	GetByAuthor(author string) []model.Quote
	GetRandom() (model.Quote, error)
	Delete(id int) bool
}
