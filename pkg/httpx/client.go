package httpx

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/portifyx/portifyx-goutils/pkg/logx"
)

type Option func(*Client)

type Client struct {
	client           *resty.Client
	baseURL          string
	redactPath       bool
	logBodies        bool
	redactedPatterns []string
}

// NewClient initializes an httpx Resty wrapper with default timeouts and logging hooks.
func NewClient(opts ...Option) *Client {
	r := resty.New().SetTimeout(30 * time.Second)
	c := &Client{
		client:           r,
		logBodies:        true,
		redactedPatterns: []string{"password", "secret", "token", "key"},
	}

	// Register trace propagation middleware
	r.OnBeforeRequest(func(cli *resty.Client, req *resty.Request) error {
		ctx := req.Context()
		if ctx != nil {
			if tid, ok := ctx.Value(logx.TraceIDKey).(string); ok && tid != "" {
				req.SetHeader("X-Trace-ID", tid)
				req.SetHeader("X-B3-TraceId", tid)
			}
			if sid, ok := ctx.Value(logx.SpanIDKey).(string); ok && sid != "" {
				req.SetHeader("X-Span-ID", sid)
				req.SetHeader("X-B3-SpanId", sid)
			}
		}

		path := req.URL
		if c.redactPath {
			path = "[REDACTED]"
		}

		logx.Debugf(ctx, "[*] outgoing http request: %s %s", req.Method, path)
		return nil
	})

	// Register status-aware logging after response
	r.OnAfterResponse(func(cli *resty.Client, resp *resty.Response) error {
		req := resp.Request
		ctx := req.Context()
		status := resp.StatusCode()
		path := req.URL
		if c.redactPath {
			path = "[REDACTED]"
		}

		duration := resp.Time()
		msg := fmt.Sprintf("[%s] %s %s -> %d (%s)", req.Method, path, duration, status, duration)

		if status >= 500 {
			logx.Errorf(ctx, nil, "[!] outgoing http error: %s", msg)
		} else if status >= 400 {
			logx.Warnf(ctx, "[!] outgoing http client error: %s", msg)
		} else {
			logx.Infof(ctx, "[*] outgoing http success: %s", msg)
		}

		return nil
	})

	// Transport error hook
	r.OnError(func(req *resty.Request, err error) {
		ctx := req.Context()
		path := req.URL
		if c.redactPath {
			path = "[REDACTED]"
		}
		logx.Errorf(ctx, err, "[!] outgoing http transport failure: %s %s", req.Method, path)
	})

	// Apply options
	for _, opt := range opts {
		opt(c)
	}

	return c
}

func WithBaseURL(url string) Option {
	return func(c *Client) {
		c.baseURL = url
		c.client.SetBaseURL(url)
	}
}

func WithTimeout(d time.Duration) Option {
	return func(c *Client) {
		c.client.SetTimeout(d)
	}
}

func WithHeader(key, val string) Option {
	return func(c *Client) {
		c.client.SetHeader(key, val)
	}
}

func WithBeforeRequest(mw resty.RequestMiddleware) Option {
	return func(c *Client) {
		c.client.OnBeforeRequest(mw)
	}
}

func WithRestyHook(fn func(*resty.Client)) Option {
	return func(c *Client) {
		fn(c.client)
	}
}

func WithRedactedPath() Option {
	return func(c *Client) {
		c.redactPath = true
	}
}

func WithoutBodyLogging() Option {
	return func(c *Client) {
		c.logBodies = false
	}
}

func (c *Client) BaseURL() string {
	return c.baseURL
}

func (c *Client) R() *resty.Request {
	return c.client.R()
}

func (c *Client) Get(ctx context.Context, url string) (*resty.Response, error) {
	return c.client.R().SetContext(ctx).Get(url)
}

func (c *Client) Post(ctx context.Context, url string, body any) (*resty.Response, error) {
	req := c.client.R().SetContext(ctx)
	if body != nil {
		req.SetBody(body)
	}
	return req.Post(url)
}

func (c *Client) Put(ctx context.Context, url string, body any) (*resty.Response, error) {
	req := c.client.R().SetContext(ctx)
	if body != nil {
		req.SetBody(body)
	}
	return req.Put(url)
}

func (c *Client) Patch(ctx context.Context, url string, body any) (*resty.Response, error) {
	req := c.client.R().SetContext(ctx)
	if body != nil {
		req.SetBody(body)
	}
	return req.Patch(url)
}

func (c *Client) Delete(ctx context.Context, url string) (*resty.Response, error) {
	return c.client.R().SetContext(ctx).Delete(url)
}

// Ensure interface compliance
var _ = http.StatusOK
var _ = strings.Contains
