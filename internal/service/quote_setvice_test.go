package service

import (
	"context"
	"quote-api/internal/model"
	"quote-api/internal/store"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuoteService(t *testing.T) {
	svc := NewQuoteService(store.NewInMemoryStore())
	ctx := context.Background()

	t.Run("AddQuote validates input", func(t *testing.T) {
		_, err := svc.AddQuote(ctx, "", "")
		assert.Error(t, err)
	})

	t.Run("AddQuote returns added quote", func(t *testing.T) {
		q, err := svc.AddQuote(ctx, "Author", "Text")
		assert.NoError(t, err)
		assert.Equal(t, "Author", q.Author)
		assert.Equal(t, "Text", q.Text)
	})

	t.Run("ListQuotes filters by author", func(t *testing.T) {
		_, err := svc.AddQuote(ctx, "Alice", "A1")
		assert.NoError(t, err)
		_, err = svc.AddQuote(ctx, "Bob", "B1")
		assert.NoError(t, err)
		_, err = svc.AddQuote(ctx, "Alice", "A2")
		assert.NoError(t, err)

		aliceQuotes := svc.ListQuotes(ctx, "Alice")
		assert.Len(t, aliceQuotes, 2)
	})

	t.Run("RandomQuote returns error when empty", func(t *testing.T) {
		emptySvc := NewQuoteService(store.NewInMemoryStore())
		_, err := emptySvc.RandomQuote(ctx)
		assert.Error(t, err)
	})

	t.Run("DeleteQuote works and errors when not found", func(t *testing.T) {
		q, _ := svc.AddQuote(ctx, "ToDel", "Remove me")
		assert.NoError(t, svc.DeleteQuote(ctx, q.ID))
		exErr := svc.DeleteQuote(ctx, q.ID)
		assert.Equal(t, model.ErrQuoteNotFound, exErr)
	})
}
