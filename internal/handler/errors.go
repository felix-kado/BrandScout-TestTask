package handler

import (
	"errors"
	"log"
	"net/http"
	"quote-api/internal/model"

	"quote-api/internal/render"
)

func errorResponse(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, model.ErrEmptyQuoteAuthor):
		render.JSON(w, http.StatusBadRequest, map[string]string{"error": "Пользователь не должен быть пустым"})
	case errors.Is(err, model.ErrEmptyQuoteText):
		render.JSON(w, http.StatusBadRequest, map[string]string{"error": "Текст не должен быть пустым"})
	case errors.Is(err, model.ErrQuoteNotFound):
		render.JSON(w, http.StatusNotFound, map[string]string{"error": "Цитата не найдена"})
	case errors.Is(err, model.ErrInvalidQuoteID):
		render.JSON(w, http.StatusBadRequest, map[string]string{"error": "Некорректный ID"})
	default:
		log.Printf("handler error: %+v", err)
		render.JSON(w, http.StatusInternalServerError, map[string]string{"error": "внутренняя ошибка сервера"})
	}
}
