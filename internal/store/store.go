package store

import "quote-api/internal/model"

type QuoteStore interface {
	Add(model.QuoteNote) model.QuoteNote
	GetAll() []model.QuoteNote
	GetByAuthor(author string) []model.QuoteNote
	GetRandom() (model.QuoteNote, error)
	Delete(id int) bool
}
