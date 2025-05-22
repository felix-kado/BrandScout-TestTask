package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"quote-api/internal/model"
)

func runStoreTests(t *testing.T, newStore func() QuoteStore) {
	t.Run("Add and GetAll", func(t *testing.T) {
		s := newStore()
		ctx := t.Context()

		q := model.Quote{Author: "X", Text: "Y"}
		added := s.Add(ctx, q)
		assert.NotZero(t, added.ID)
		assert.Equal(t, "X", added.Author)

		all := s.GetAll(ctx)
		assert.Len(t, all, 1)
	})

	t.Run("Filter by Author", func(t *testing.T) {
		s := newStore()
		ctx := t.Context()

		s.Add(ctx, model.Quote{Author: "A", Text: "1"})
		s.Add(ctx, model.Quote{Author: "B", Text: "2"})
		s.Add(ctx, model.Quote{Author: "A", Text: "3"})

		a := s.GetByAuthor(ctx, "A")
		assert.Len(t, a, 2)
	})

	t.Run("Delete", func(t *testing.T) {
		s := newStore()
		ctx := t.Context()

		q := s.Add(ctx, model.Quote{Author: "Del", Text: "Me"})
		assert.True(t, s.Delete(ctx, q.ID))
		assert.False(t, s.Delete(ctx, q.ID))
	})

	t.Run("GetRandom", func(t *testing.T) {
		s := newStore()
		ctx := t.Context()

		_, err := s.GetRandom(ctx)
		assert.Error(t, err)

		s.Add(ctx, model.Quote{Author: "Rand", Text: "Q"})
		q, err := s.GetRandom(ctx)
		assert.NoError(t, err)
		assert.Equal(t, "Rand", q.Author)
	})
}
