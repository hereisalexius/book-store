package handler

// ErrorResponse is the standard error body returned on failures.
type ErrorResponse struct {
	Error string `json:"error"`
}