package apperror

import (
	"errors"
	"fmt"
	"net/http"
)

type ErrorCode int
type ErrorMessage string

const (
	// 5xxx: Internal & Downstream Errors
	CodeInternalServerError ErrorCode = 5000
	CodeDownstreamFailed     ErrorCode = 5110
	CodeCircuitBreakerOpen   ErrorCode = 5030
	CodeDatabaseError        ErrorCode = 5200

	// 6xxx: Client Input & Validation Errors
	CodeInvalidJSON          ErrorCode = 6001
	CodeInvalidField         ErrorCode = 6010
	CodeMissingRequiredField ErrorCode = 6050
	CodeResourceNotFound     ErrorCode = 6040

	// 7xxx: Exchange Connector Errors
	CodeInvalidExchangeKey   ErrorCode = 7001
	CodeExchangeRateLimit    ErrorCode = 7002
	CodeExchangeOffline      ErrorCode = 7003

	// 9xxx: Infrastructure Failures
	CodeBrokerError          ErrorCode = 9001
)

type AppError struct {
	Code           ErrorCode `json:"code"`
	Message        string    `json:"message"`
	HTTPStatusCode int       `json:"-"`
	Cause          error     `json:"-"`
}

func (e AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

func (e AppError) Unwrap() error {
	return e.Cause
}

func (e AppError) WithReason(err error, reason string) AppError {
	copyErr := e
	copyErr.Cause = err
	if reason != "" {
		copyErr.Message = fmt.Sprintf("%s: %s", e.Message, reason)
	}
	return copyErr
}

func IsAppError(err error) (AppError, bool) {
	var appErr AppError
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return AppError{}, false
}

// Constructors
func NewInternalServerError(cause error) AppError {
	return AppError{
		Code:           CodeInternalServerError,
		Message:        "Internal server error",
		HTTPStatusCode: http.StatusInternalServerError,
		Cause:          cause,
	}
}

func NewDownstreamError(service string, cause error) AppError {
	return AppError{
		Code:           CodeDownstreamFailed,
		Message:        fmt.Sprintf("Downstream service call failed: %s", service),
		HTTPStatusCode: http.StatusBadGateway,
		Cause:          cause,
	}
}

func NewDatabaseError(cause error) AppError {
	return AppError{
		Code:           CodeDatabaseError,
		Message:        "Database operation failed",
		HTTPStatusCode: http.StatusInternalServerError,
		Cause:          cause,
	}
}

func NewInvalidInputType() AppError {
	return AppError{
		Code:           CodeInvalidJSON,
		Message:        "Invalid JSON format or request payload",
		HTTPStatusCode: http.StatusUnprocessableEntity,
	}
}

func NewValidationError(field, reason string) AppError {
	return AppError{
		Code:           CodeInvalidField,
		Message:        fmt.Sprintf("Validation failed for '%s': %s", field, reason),
		HTTPStatusCode: http.StatusUnprocessableEntity,
	}
}

func NewMissingRequiredField(field string) AppError {
	return AppError{
		Code:           CodeMissingRequiredField,
		Message:        fmt.Sprintf("Required field missing: %s", field),
		HTTPStatusCode: http.StatusUnprocessableEntity,
	}
}

func NewNotFoundError(resource string) AppError {
	return AppError{
		Code:           CodeResourceNotFound,
		Message:        fmt.Sprintf("%s not found", resource),
		HTTPStatusCode: http.StatusNotFound,
	}
}

func NewInvalidExchangeKey() AppError {
	return AppError{
		Code:           CodeInvalidExchangeKey,
		Message:        "Invalid exchange API key or secret",
		HTTPStatusCode: http.StatusUnauthorized,
	}
}

func NewExchangeRateLimit() AppError {
	return AppError{
		Code:           CodeExchangeRateLimit,
		Message:        "Exchange rate limit exceeded",
		HTTPStatusCode: http.StatusTooManyRequests,
	}
}

func NewExchangeOffline() AppError {
	return AppError{
		Code:           CodeExchangeOffline,
		Message:        "Upstream exchange is currently offline or unreachable",
		HTTPStatusCode: http.StatusServiceUnavailable,
	}
}

func NewBrokerError(cause error) AppError {
	return AppError{
		Code:           CodeBrokerError,
		Message:        "Message broker communication failure",
		HTTPStatusCode: http.StatusInternalServerError,
		Cause:          cause,
	}
}
