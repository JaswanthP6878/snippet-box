package models

import "errors"

var (
	ErrNoRecord = errors.New("models: No matching records found")

	ErrInvalidCredentials = errors.New("models:invalid credentials")

	ErrDuplicateEmail = errors.New("models: duplicate email")
)
