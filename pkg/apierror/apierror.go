package apierror

import "net/http"

const (
	NameInternalServerError = "InternalServerError"
	NameBadRequestError     = "BadRequestError"
)

var _ error = (*APIError)(nil)

type APIError struct {
	Message    string `json:"message"`
	Name       string `json:"name"`
	StatusCode int    `json:"statusCode"`
	BaseError  error  `json:"-"`
}

func (apiErr *APIError) Error() string {
	return apiErr.Message
}

func NewBadRequestError(message string, messageLocalizerKey string) *APIError {
	return &APIError{
		Message:    message,
		Name:       NameBadRequestError,
		StatusCode: http.StatusBadRequest,
	}
}

func NewInternalServerError() *APIError {
	return &APIError{
		Name:       NameInternalServerError,
		StatusCode: http.StatusInternalServerError,
	}
}
