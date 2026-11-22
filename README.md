# portifyx-goutils

Shared Go utility module for the PortifyX microservice ecosystem, modeled directly after the Bitkub Trade Platform engineering standards (`tsl-goutils`).

## Features & Packages

- **`pkg/logx`**: Context-first Zap logging with DataDog / W3C trace ID and span ID correlation, ISO8601 timestamps, JSON output for staging/production, and pretty console logs for local development.
- **`pkg/httpx`**: Resty client wrapper with automated outbound trace header injection (`X-Trace-ID`), latency tracking, status-aware logging (>=500 Error, >=400 Warn, else Info), and URL credential redaction.
- **`pkg/apperror`**: Typed application error system with numeric error bands:
  - `5xxx`: Internal & downstream failures (5000 generic internal, 5110 downstream HTTP failure, 5200 database error).
  - `6xxx`: Client input & validation errors (6001 invalid JSON, 6010 bad field, 6040 not found, 6050 missing required field).
  - `7xxx`: Exchange connector rejections & rate limits.
  - `9xxx`: Broker / queue infrastructure failures.
- **`pkg/httpserver`**: Standard response envelopes (`{ "code": 0, "message": "success", "data": ... }`) and secure error formatting (safely maps `AppError` to HTTP status codes without leaking internal stack traces).
- **`pkg/validators`**: Struct validation powered by `go-playground/validator/v10` mapping to standardized `AppError` responses.
- **`pkg/timeutil`**: Precision-agnostic Unix timestamp conversion (seconds, milliseconds, microseconds, nanoseconds).
- **`pkg/decimal`**: Safe financial math and parsing utilities using `shopspring/decimal`.

## Testing

```bash
make test
```
