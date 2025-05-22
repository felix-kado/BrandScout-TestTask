package test

import (
	"net/http"
	"quote-api/internal/handler"
	"quote-api/internal/service"
	"quote-api/internal/store"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"quote-api/internal/testutil"
)

func TestAPI_FullFlow(t *testing.T) {
	st := store.NewInMemoryStore()
	svc := service.NewQuoteService(st)
	r := handler.NewRouter(svc)

	t.Run("CreateQuotes", func(t *testing.T) {
		cases := []struct{ Author, Quote string }{{"Alice", "First"}, {"Bob", "Second"}, {"Alice", "Third"}}
		for _, c := range cases {
			t.Run(c.Author+"_"+c.Quote, func(t *testing.T) {
				resp := testutil.DoRequest(r, "POST", "/quotes", map[string]string{"author": c.Author, "quote": c.Quote})
				testutil.AssertStatus(t, resp, http.StatusCreated)
				var out map[string]interface{}
				testutil.ParseJSON(t, resp, &out)
				assert.Equal(t, c.Author, out["author"])
			})
		}
	})

	t.Run("GetAllQuotes", func(t *testing.T) {
		resp := testutil.DoRequest(r, "GET", "/quotes", nil)
		testutil.AssertStatus(t, resp, http.StatusOK)
		var list []map[string]interface{}
		testutil.ParseJSON(t, resp, &list)
		assert.Len(t, list, 3)
	})

	t.Run("FilterByAuthor", func(t *testing.T) {
		resp := testutil.DoRequest(r, "GET", "/quotes?author=Alice", nil)
		testutil.AssertStatus(t, resp, http.StatusOK)
		var list []map[string]interface{}
		testutil.ParseJSON(t, resp, &list)
		assert.Len(t, list, 2)
	})

	t.Run("GetRandomQuote", func(t *testing.T) {
		resp := testutil.DoRequest(r, "GET", "/quotes/random", nil)
		testutil.AssertStatus(t, resp, http.StatusOK)
		var q map[string]interface{}
		testutil.ParseJSON(t, resp, &q)
		assert.NotEmpty(t, q["quote"])
	})

	t.Run("DeleteQuote", func(t *testing.T) {
		resp := testutil.DoRequest(r, "GET", "/quotes", nil)
		var all []map[string]interface{}
		testutil.ParseJSON(t, resp, &all)
		id := int(all[0]["id"].(float64))

		delPath := "/quotes/" + strconv.Itoa(id)
		dResp := testutil.DoRequest(r, "DELETE", delPath, nil)
		testutil.AssertStatus(t, dResp, http.StatusNoContent)

		// повторное удаление -> 404
		dResp = testutil.DoRequest(r, "DELETE", delPath, nil)
		testutil.AssertStatus(t, dResp, http.StatusNotFound)
	})
}

func TestAPI_InvalidInput(t *testing.T) {
	st := store.NewInMemoryStore()
	svc := service.NewQuoteService(st)
	r := handler.NewRouter(svc)

	t.Run("BadPOST", func(t *testing.T) {
		resp := testutil.DoRequest(r, "POST", "/quotes", nil)
		testutil.AssertStatus(t, resp, http.StatusBadRequest)
	})

	t.Run("BadDELETEID", func(t *testing.T) {
		resp := testutil.DoRequest(r, "DELETE", "/quotes/abc", nil)
		testutil.AssertStatus(t, resp, http.StatusBadRequest)
	})

	t.Run("RandomOnEmptyStore", func(t *testing.T) {
		resp := testutil.DoRequest(r, "GET", "/quotes/random", nil)
		testutil.AssertStatus(t, resp, http.StatusNotFound)
	})
}
