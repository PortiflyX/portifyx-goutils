package apperror

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppErrorCodes(t *testing.T) {
	err := NewInternalServerError(errors.New("db down"))
	assert.Equal(t, CodeInternalServerError, err.Code)
	assert.Equal(t, http.StatusInternalServerError, err.HTTPStatusCode)
	assert.Contains(t, err.Error(), "[5000]")

	vErr := NewValidationError("amount", "must be positive")
	assert.Equal(t, CodeInvalidField, vErr.Code)
	assert.Equal(t, http.StatusUnprocessableEntity, vErr.HTTPStatusCode)

	appErr, ok := IsAppError(vErr)
	assert.True(t, ok)
	assert.Equal(t, CodeInvalidField, appErr.Code)

	_, ok = IsAppError(errors.New("standard error"))
	assert.False(t, ok)
}
