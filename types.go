package jev

import "fmt"

type Logger interface {
	Debug(msg string, args ...any)
}

type APIError struct {
	StatusCode int
	Body       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("jev: api returned %d: %s", e.StatusCode, e.Body)
}
