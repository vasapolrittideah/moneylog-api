package contract

import "time"

// APIResponse is the response structure for the API.
type APIResponse struct {
	Data      any       `json:"data,omitempty"`
	Error     *APIError `json:"error,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

// APIError is the error structure for the API.
type APIError struct {
	Code    string               `json:"code"`
	Message string               `json:"message"`
	Details []APIValidationError `json:"details,omitempty"`
}

// APIValidationError is the validation error structure for the API.
type APIValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   any    `json:"value,omitempty"`
}

const (
	ErrorCodeValidation   = "VALIDATION_ERROR"
	ErrorCodeNotFound     = "NOT_FOUND"
	ErrorCodeUnauthorized = "UNAUTHORIZED"
	ErrorCodeForbidden    = "FORBIDDEN"
	ErrorCodeInternal     = "INTERNAL_ERROR"
	ErrorCodeBadRequest   = "BAD_REQUEST"
	ErrorCodeConflict     = "CONFLICT"
	ErrorCodeRateLimit    = "RATE_LIMIT_EXCEEDED"
)

func NewSuccessResponse(data any) APIResponse {
	return APIResponse{
		Data:      data,
		Timestamp: time.Now(),
	}
}

func NewErrorResponse(code, message string) APIResponse {
	return APIResponse{
		Error: &APIError{
			Code:    code,
			Message: message,
		},
		Timestamp: time.Now(),
	}
}

func NewValidationErrorResponse(details []APIValidationError) APIResponse {
	return APIResponse{
		Error: &APIError{
			Code:    ErrorCodeValidation,
			Message: "Validation failed",
			Details: details,
		},
		Timestamp: time.Now(),
	}
}
