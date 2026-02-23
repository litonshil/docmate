package response

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// BaseApiResponse represents a standardized API response structure
type BaseApiResponse struct {
	Success bool        `json:"success"`
	Code    string      `json:"code"` // HTTP status code as string (e.g., "200", "400", "500")
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

// JSON sends a JSON response with the given status code using BaseApiResponse
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

// Success sends a success response
func Success(c echo.Context, message string, data interface{}) error {
	return JSON(c, http.StatusOK, true, message, data, nil)
}

// Error sends an error response
func Error(c echo.Context, statusCode int, message string, errs interface{}) error {
	return JSON(c, statusCode, false, message, nil, errs)
}

// BadRequest sends a 400 Bad Request response
func BadRequest(c echo.Context, message string) error {
	return Error(c, http.StatusBadRequest, message, nil)
}

// Unauthorized sends a 401 Unauthorized response
func Unauthorized(c echo.Context, message string) error {
	return Error(c, http.StatusUnauthorized, message, nil)
}

// Forbidden sends a 403 Forbidden response
func Forbidden(c echo.Context, message string) error {
	return Error(c, http.StatusForbidden, message, nil)
}

// NotFound sends a 404 Not Found response
func NotFound(c echo.Context, message string) error {
	return Error(c, http.StatusNotFound, message, nil)
}

// InternalServerError sends a 500 Internal Server Error response
func InternalServerError(c echo.Context, message string) error {
	return Error(c, http.StatusInternalServerError, message, nil)
}
