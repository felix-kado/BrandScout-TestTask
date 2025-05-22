package model

import "errors"

var (
	ErrEmptyAuthor   = errors.New("author cannot be empty")
	ErrEmptyText     = errors.New("text cannot be empty")
	ErrQuoteNotFound = errors.New("quote not found")
	ErrInvalidID     = errors.New("invalid ID")
)
