package response

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

// BaseApiResponse represents a standardized API response structure.
type BaseApiResponse struct {
	Success bool        `json:"success"`
	Code    string      `json:"code"` // HTTP status code as string (e.g., "200", "400", "500")
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

// JSON sends a JSON response with the given status code using BaseApiResponse.
func JSON(c echo.Context, statusCode int, success bool, message string, data interface{}, errs interface{}) error {
	resp := BaseApiResponse{
		Success: success,
		Code:    fmt.Sprintf("%d", statusCode), // Convert status code to string
		Message: message,
		Data:    data,
		Errors:  errs,
	}

	return c.JSON(statusCode, resp)
}

// Success sends a success response.
func Success(c echo.Context, message string, data interface{}) error {
	return JSON(c, http.StatusOK, true, message, data, nil)
}

// Created sends a 201 Created response.
func Created(c echo.Context, message string, data interface{}) error {
	return JSON(c, http.StatusCreated, true, message, data, nil)
}

// Error sends an error response.
func Error(c echo.Context, statusCode int, message string, errs interface{}) error {
	return JSON(c, statusCode, false, message, nil, errs)
}

// BadRequest sends a 400 Bad Request response.
func BadRequest(c echo.Context, message string) error {
	return Error(c, http.StatusBadRequest, message, nil)
}

// Unauthorized sends a 401 Unauthorized response.
func Unauthorized(c echo.Context, message string) error {
	return Error(c, http.StatusUnauthorized, message, nil)
}

// Forbidden sends a 403 Forbidden response.
func Forbidden(c echo.Context, message string) error {
	return Error(c, http.StatusForbidden, message, nil)
}

// NotFound sends a 404 Not Found response.
func NotFound(c echo.Context, message string) error {
	return Error(c, http.StatusNotFound, message, nil)
}

// InternalServerError sends a 500 Internal Server Error response, sanitizing database/driver errors to protect server internals.
func InternalServerError(c echo.Context, message string) error {
	slog.Error("internal server error occurred", "error", message)

	sanitizedMsg := sanitizeErrorMessage(message)

	return Error(c, http.StatusInternalServerError, sanitizedMsg, nil)
}

// sanitizeErrorMessage hides database/system-level details in HTTP 500 responses.
func sanitizeErrorMessage(message string) string {
	lowerMsg := strings.ToLower(message)
	dbKeywords := []string{
		"sql:", "gorm", "postgres", "pq:", "pgconn", "connection", "dial tcp",
		"database", "rows", "scan", "tx", "transaction", "driver", "query",
		"column", "table", "relation", "syntax", "violation", "duplicate key",
		"constraint", "schema", "db client", "record not found",
		"sqlstate", "invalid input value", "enum",
	}
	for _, kw := range dbKeywords {
		if strings.Contains(lowerMsg, kw) {
			return "Something went wrong. Please try again later."
		}
	}

	return message
}
