package logx

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogxInit(t *testing.T) {
	logger, err := Init("test-service", "v1.0.0", "local")
	assert.NoError(t, err)
	assert.NotNil(t, logger)

	ctx := ContextWithTrace(context.Background(), "trace-123", "span-456")
	Info(ctx, "hello test info")
	Debug(ctx, "hello test debug")
	Warn(ctx, "hello test warn")
	Error(ctx, assert.AnError, "hello test error")
	Infow(ctx, "structured info", "key1", "val1")

	_ = Sync()
}

func TestJSONString(t *testing.T) {
	data := map[string]string{"foo": "bar"}
	res := JSONString(data)
	assert.Equal(t, `{"foo":"bar"}`, res)
}
