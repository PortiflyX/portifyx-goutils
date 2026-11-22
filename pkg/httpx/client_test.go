package httpx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/portifyx/portifyx-goutils/pkg/logx"
	"github.com/stretchr/testify/assert"
)

func TestHttpxClient(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID := r.Header.Get("X-Trace-ID")
		assert.Equal(t, "trace-abc-123", traceID)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	}))
	defer ts.Close()

	client := NewClient(
		WithBaseURL(ts.URL),
		WithTimeout(5*time.Second),
	)

	ctx := logx.ContextWithTrace(context.Background(), "trace-abc-123", "span-xyz")
	resp, err := client.Get(ctx, "/test")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())
	assert.Equal(t, ts.URL, client.BaseURL())
}
