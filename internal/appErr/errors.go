package appErr

import "errors"

var (
	ErrAlreadyExists = errors.New("already exists")
	ErrInvalidInput  = errors.New("invalid input")
	ErrRequiredField = errors.New("required field")
	ErrInvalidRange  = errors.New("invalid time range")
	ErrNotFound      = errors.New("not found")
	ErrConflict      = errors.New("conflict")
)
