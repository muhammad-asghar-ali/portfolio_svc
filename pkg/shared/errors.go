package errors

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type (
	APIError struct {
		Message error `json:"message,omitempty"`
		Code    int   `json:"code"`
	}
)

func (e *APIError) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

func NewHttpError(code int, message error) *APIError {
	return &APIError{
		Code:    code,
		Message: message,
	}
}

// NewBadRequestError returns APIError with status code 400.
func NewBadRequestError(message error) *APIError {
	return NewHttpError(http.StatusBadRequest, message)
}

// NewUnauthorizedError returns APIError with status code 401.
func NewUnauthorizedError(message error) *APIError {
	return NewHttpError(http.StatusUnauthorized, message)
}

// NewNotFoundError returns APIError with status code 404.
func NewNotFoundError(message error) *APIError {
	return NewHttpError(http.StatusNotFound, message)
}

// NewForbiddenError returns APIError with status code 403.
func NewForbiddenError(message error) *APIError {
	return NewHttpError(http.StatusForbidden, message)
}

// NewInternalServerError returns APIError with status code 500.
func NewInternalServerError(message error) *APIError {
	return NewHttpError(http.StatusInternalServerError, message)
}

func HandleHttpError(c *gin.Context, err error) {
	c.Error(err)
}
