package testutil

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// DoRequest выполняет HTTP-запрос к переданному handler и возвращает ResponseRecorder
func DoRequest(h http.Handler, method, path string, body interface{}) *httptest.ResponseRecorder {
	var buf bytes.Buffer
	if body != nil {
		_ = json.NewEncoder(&buf).Encode(body)
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	return w
}

// AssertStatus проверяет, что код ответа совпадает с ожиданием
func AssertStatus(t *testing.T, w *httptest.ResponseRecorder, expected int) {
	assert.Equal(t, expected, w.Code)
}

// ParseJSON декодирует JSON-ответ в dest и проверяет отсутствие ошибки
func ParseJSON(t *testing.T, w *httptest.ResponseRecorder, dest interface{}) {
	err := json.NewDecoder(w.Body).Decode(dest)
	assert.NoError(t, err)
}
