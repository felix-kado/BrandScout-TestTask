package handler

import (
	"net/http"
	"quote-api/internal/service"
	"quote-api/internal/store"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"quote-api/internal/testutil"
)

func TestCreateAndGet(t *testing.T) {
	st := store.NewInMemoryStore()
	svc := service.NewQuoteService(st)
	h := NewRouter(svc)

	t.Run("CreateQuote", func(t *testing.T) {
		payload := map[string]string{"author": "Tester", "quote": "Hello World"}
		w := testutil.DoRequest(h, "POST", "/quotes", payload)
		testutil.AssertStatus(t, w, http.StatusCreated)

		var created struct {
			ID     int    `json:"id"`
			Author string `json:"author"`
			Quote  string `json:"quote"`
		}
		testutil.ParseJSON(t, w, &created)
		assert.Equal(t, "Tester", created.Author)
	})

	t.Run("GetAllQuotes", func(t *testing.T) {
		w := testutil.DoRequest(h, "GET", "/quotes", nil)
		testutil.AssertStatus(t, w, http.StatusOK)

		var list []struct {
			ID     int    `json:"id"`
			Author string `json:"author"`
			Quote  string `json:"quote"`
		}
		testutil.ParseJSON(t, w, &list)
		assert.GreaterOrEqual(t, len(list), 1)
	})
}

func TestDeleteHandler(t *testing.T) {
	st := store.NewInMemoryStore()
	svc := service.NewQuoteService(st)
	h := NewRouter(svc)
	// Создаём цитату
	w := testutil.DoRequest(h, "POST", "/quotes", map[string]string{"author": "Del", "quote": "To remove"})
	testutil.AssertStatus(t, w, http.StatusCreated)

	var created struct {
		ID int `json:"id"`
	}
	testutil.ParseJSON(t, w, &created)

	// Удаляем цитату
	delPath := "/quotes/" + strconv.Itoa(created.ID)
	wDel := testutil.DoRequest(h, "DELETE", delPath, nil)
	testutil.AssertStatus(t, wDel, http.StatusNoContent)

	// Повторное удаление -> 404
	wDel2 := testutil.DoRequest(h, "DELETE", delPath, nil)
	testutil.AssertStatus(t, wDel2, http.StatusNotFound)
}
