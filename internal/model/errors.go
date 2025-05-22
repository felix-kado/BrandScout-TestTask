package model

import "errors"

var (
	ErrEmptyQuoteAuthor = errors.New("author cannot be empty")
	ErrEmptyQuoteText   = errors.New("text cannot be empty")
	ErrQuoteNotFound    = errors.New("quote not found")
	ErrInvalidQuoteID   = errors.New("invalid ID")
)
