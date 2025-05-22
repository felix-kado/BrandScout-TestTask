package service

import (
	"errors"
	"quote-api/internal/store"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuoteService(t *testing.T) {
	s := NewQuoteService(store.NewInMemoryStore())

	t.Run("AddQuote validates input", func(t *testing.T) {
		_, err := s.AddQuote("", "")
		assert.Error(t, err)
	})

	t.Run("AddQuote returns added quote", func(t *testing.T) {
		q, err := s.AddQuote("Author", "Text")
		assert.NoError(t, err)
		assert.Equal(t, "Author", q.Author)
		assert.Equal(t, "Text", q.Text)
	})

	t.Run("ListQuotes filters by author", func(t *testing.T) {
		_, err := s.AddQuote("Alice", "A1")
		assert.NoError(t, err)
		_, err = s.AddQuote("Bob", "B1")
		assert.NoError(t, err)
		_, err = s.AddQuote("Alice", "A2")
		assert.NoError(t, err)

		aliceQuotes := s.ListQuotes("Alice")
		assert.Len(t, aliceQuotes, 2)
	})

	t.Run("RandomQuote returns error when empty", func(t *testing.T) {
		emptyService := NewQuoteService(store.NewInMemoryStore())
		_, err := emptyService.RandomQuote()
		assert.Error(t, err)
	})

	t.Run("DeleteQuote works", func(t *testing.T) {
		q, _ := s.AddQuote("ToDel", "Remove me")
		assert.NoError(t, s.DeleteQuote(q.ID))
		assert.Equal(t, errors.New("quote not found"), s.DeleteQuote(q.ID))
	})
}
