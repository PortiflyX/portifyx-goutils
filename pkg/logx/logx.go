package logx

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type contextKey string

const (
	TraceIDKey contextKey = "trace-id"
	SpanIDKey  contextKey = "span-id"
)

var (
	globalLogger *zap.Logger
	globalSugar  *zap.SugaredLogger
	once         sync.Once
)

func init() {
	InitNop()
}

// Init configures and replaces the global logger.
func Init(serviceName, serviceVersion, env string) (*zap.Logger, error) {
	var cfg zap.Config
	if env == "local" || env == "dev" || env == "development" {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		cfg = zap.NewProductionConfig()
		cfg.EncoderConfig.TimeKey = "timestamp"
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		cfg.EncoderConfig.StacktraceKey = "" // omit stacktrace in prod logs
		cfg.InitialFields = map[string]interface{}{
			"service-name":    serviceName,
			"service-version": serviceVersion,
			"env":             env,
		}
	}

	logger, err := cfg.Build(zap.AddCallerSkip(1))
	if err != nil {
		return nil, fmt.Errorf("build logger: %w", err)
	}

	globalLogger = logger
	globalSugar = logger.Sugar()
	zap.ReplaceGlobals(logger)

	return logger, nil
}

// InitNop installs a no-op logger for tests.
func InitNop() *zap.Logger {
	nop := zap.NewNop()
	globalLogger = nop
	globalSugar = nop.Sugar()
	zap.ReplaceGlobals(nop)
	return nop
}

// Sync flushes buffered log entries.
func Sync() error {
	if globalLogger != nil {
		return globalLogger.Sync()
	}
	return nil
}

// ContextWithTrace attaches trace and span IDs to the context.
func ContextWithTrace(ctx context.Context, traceID, spanID string) context.Context {
	ctx = context.WithValue(ctx, TraceIDKey, traceID)
	if spanID != "" {
		ctx = context.WithValue(ctx, SpanIDKey, spanID)
	}
	return ctx
}

func extractTraceFields(ctx context.Context) []zap.Field {
	if ctx == nil {
		return nil
	}
	fields := make([]zap.Field, 0, 4)

	if tid, ok := ctx.Value(TraceIDKey).(string); ok && tid != "" {
		fields = append(fields, zap.String("trace-id", tid))
		fields = append(fields, zap.String("dd.trace_id", tid))
	} else {
		fields = append(fields, zap.String("trace-id", "0"))
		fields = append(fields, zap.String("dd.trace_id", "0"))
	}

	if sid, ok := ctx.Value(SpanIDKey).(string); ok && sid != "" {
		fields = append(fields, zap.String("span-id", sid))
		fields = append(fields, zap.String("dd.span_id", sid))
	} else {
		fields = append(fields, zap.String("span-id", "0"))
		fields = append(fields, zap.String("dd.span_id", "0"))
	}

	return fields
}

// WithContext returns a logger tagged with trace context fields.
func WithContext(ctx context.Context) *zap.Logger {
	fields := extractTraceFields(ctx)
	return globalLogger.WithOptions(zap.AddCallerSkip(-1)).With(fields...)
}

// Log returns the global sugared logger without context.
func Log() *zap.SugaredLogger {
	return globalSugar
}

// JSONString serializes an object into JSON string, safe against errors.
func JSONString(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(b)
}

// Standard logging methods with trace context
func Info(ctx context.Context, msg string, fields ...zap.Field) {
	traceFields := extractTraceFields(ctx)
	allFields := append(traceFields, fields...)
	globalLogger.Info(msg, allFields...)
}

func Infof(ctx context.Context, format string, args ...any) {
	traceFields := extractTraceFields(ctx)
	globalLogger.With(traceFields...).Sugar().Infof(format, args...)
}

func Infow(ctx context.Context, msg string, keysAndValues ...any) {
	traceFields := extractTraceFields(ctx)
	sugar := globalLogger.With(traceFields...).Sugar()
	sugar.Infow(msg, keysAndValues...)
}

func Debug(ctx context.Context, msg string, fields ...zap.Field) {
	traceFields := extractTraceFields(ctx)
	allFields := append(traceFields, fields...)
	globalLogger.Debug(msg, allFields...)
}

func Debugf(ctx context.Context, format string, args ...any) {
	traceFields := extractTraceFields(ctx)
	globalLogger.With(traceFields...).Sugar().Debugf(format, args...)
}

func Debugw(ctx context.Context, msg string, keysAndValues ...any) {
	traceFields := extractTraceFields(ctx)
	globalLogger.With(traceFields...).Sugar().Debugw(msg, keysAndValues...)
}

func Warn(ctx context.Context, msg string, fields ...zap.Field) {
	traceFields := extractTraceFields(ctx)
	allFields := append(traceFields, fields...)
	globalLogger.Warn(msg, allFields...)
}

func Warnf(ctx context.Context, format string, args ...any) {
	traceFields := extractTraceFields(ctx)
	globalLogger.With(traceFields...).Sugar().Warnf(format, args...)
}

func Warnw(ctx context.Context, msg string, keysAndValues ...any) {
	traceFields := extractTraceFields(ctx)
	globalLogger.With(traceFields...).Sugar().Warnw(msg, keysAndValues...)
}

func Error(ctx context.Context, cause error, msg string, fields ...zap.Field) {
	traceFields := extractTraceFields(ctx)
	allFields := append(traceFields, zap.NamedError("cause", cause))
	allFields = append(allFields, fields...)
	globalLogger.Error(msg, allFields...)
}

func Errorf(ctx context.Context, cause error, format string, args ...any) {
	traceFields := extractTraceFields(ctx)
	fields := append(traceFields, zap.NamedError("cause", cause))
	msg := fmt.Sprintf(format, args...)
	globalLogger.Error(msg, fields...)
}

func Errorw(ctx context.Context, msg string, keysAndValues ...any) {
	traceFields := extractTraceFields(ctx)
	globalLogger.With(traceFields...).Sugar().Errorw(msg, keysAndValues...)
}

func Fatal(ctx context.Context, cause error, msg string, fields ...zap.Field) {
	traceFields := extractTraceFields(ctx)
	allFields := append(traceFields, zap.NamedError("cause", cause))
	allFields = append(allFields, fields...)
	globalLogger.Fatal(msg, allFields...)
}

func Fatalf(ctx context.Context, cause error, format string, args ...any) {
	traceFields := extractTraceFields(ctx)
	fields := append(traceFields, zap.NamedError("cause", cause))
	msg := fmt.Sprintf(format, args...)
	globalLogger.Fatal(msg, fields...)
}

func Fatalw(ctx context.Context, msg string, keysAndValues ...any) {
	traceFields := extractTraceFields(ctx)
	globalLogger.With(traceFields...).Sugar().Fatalw(msg, keysAndValues...)
}

// Ensure clean exit in tests
func _() {
	_ = time.Second
	_ = os.Stdout
}
