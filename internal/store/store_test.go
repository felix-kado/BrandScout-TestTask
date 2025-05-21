package store

import (
	"github.com/stretchr/testify/assert"
	"quote-api/internal/model"
	"testing"
)

func runStoreTests(t *testing.T, newStore func() QuoteStore) {
	t.Run("Add and GetAll", func(t *testing.T) {
		s := newStore()
		q := model.QuoteNote{Author: "X", Quote: "Y"}
		added := s.Add(q)
		assert.NotZero(t, added.ID)
		assert.Equal(t, "X", added.Author)
		all := s.GetAll()
		assert.Len(t, all, 1)
	})

	t.Run("Filter by Author", func(t *testing.T) {
		s := newStore()
		s.Add(model.QuoteNote{Author: "A", Quote: "1"})
		s.Add(model.QuoteNote{Author: "B", Quote: "2"})
		s.Add(model.QuoteNote{Author: "A", Quote: "3"})
		a := s.GetByAuthor("A")
		assert.Len(t, a, 2)
	})

	t.Run("Delete", func(t *testing.T) {
		s := newStore()
		q := s.Add(model.QuoteNote{Author: "Del", Quote: "Me"})
		assert.True(t, s.Delete(q.ID))
		assert.False(t, s.Delete(q.ID))
	})

	t.Run("GetRandom", func(t *testing.T) {
		s := newStore()
		_, err := s.GetRandom()
		assert.Error(t, err)
		s.Add(model.QuoteNote{Author: "Rand", Quote: "Q"})
		q, err := s.GetRandom()
		assert.NoError(t, err)
		assert.Equal(t, "Rand", q.Author)
	})
}
