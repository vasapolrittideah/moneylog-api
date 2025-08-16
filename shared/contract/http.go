package contract

import (
	"net/http"
	"time"

	"google.golang.org/grpc/codes"
)

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

func ErrorCodeFromGRPCCode(code codes.Code) string {
	switch code {
	case codes.OK:
		return ""
	case codes.Canceled:
		return ErrorCodeInternal
	case codes.Unknown:
		return ErrorCodeInternal
	case codes.InvalidArgument:
		return ErrorCodeBadRequest
	case codes.DeadlineExceeded:
		return ErrorCodeInternal
	case codes.NotFound:
		return ErrorCodeNotFound
	case codes.AlreadyExists:
		return ErrorCodeConflict
	case codes.PermissionDenied:
		return ErrorCodeForbidden
	case codes.ResourceExhausted:
		return ErrorCodeRateLimit
	default:
		return ErrorCodeInternal
	}
}

func HTTPStatusFromGRPCCode(code codes.Code) int {
	switch code {
	case codes.OK:
		return http.StatusOK
	case codes.Canceled:
		return http.StatusRequestTimeout
	case codes.Unknown:
		return http.StatusInternalServerError
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.DeadlineExceeded:
		return http.StatusGatewayTimeout
	case codes.NotFound:
		return http.StatusNotFound
	case codes.AlreadyExists:
		return http.StatusConflict
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.ResourceExhausted:
		return http.StatusTooManyRequests
	case codes.FailedPrecondition:
		return http.StatusPreconditionFailed
	case codes.Aborted:
		return http.StatusConflict
	case codes.OutOfRange:
		return http.StatusBadRequest
	case codes.Unimplemented:
		return http.StatusNotImplemented
	case codes.Internal:
		return http.StatusInternalServerError
	case codes.Unavailable:
		return http.StatusServiceUnavailable
	case codes.DataLoss:
		return http.StatusInternalServerError
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}
