package validators

import (
	"testing"

	"github.com/portifyx/portifyx-goutils/pkg/apperror"
	"github.com/stretchr/testify/assert"
)

type DummyPayload struct {
	Symbol string `validate:"required,symbol_format"`
	Amount string `validate:"required"`
}

func TestValidator(t *testing.T) {
	v, err := NewValidator()
	assert.NoError(t, err)

	valid := DummyPayload{
		Symbol: "BTC_USDT",
		Amount: "1.5",
	}
	assert.NoError(t, v.ValidateStruct(valid))

	invalid := DummyPayload{
		Symbol: "",
		Amount: "1.5",
	}
	err = v.ValidateStruct(invalid)
	assert.Error(t, err)
	appErr, ok := apperror.IsAppError(err)
	assert.True(t, ok)
	assert.Equal(t, apperror.CodeMissingRequiredField, appErr.Code)
}
