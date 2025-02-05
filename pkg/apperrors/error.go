package apperrors

import "errors"

var (
	NotFound error = errors.New("not found")
	Quit     error = errors.New("quit")
)
