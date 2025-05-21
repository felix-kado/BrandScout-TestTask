package store

import "testing"

func TestInMemoryStore_CommonTests(t *testing.T) {
	runStoreTests(t, func() QuoteStore {
		return NewInMemoryStore()
	})
}
