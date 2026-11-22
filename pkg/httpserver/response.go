package httpserver

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/portifyx/portifyx-goutils/pkg/apperror"
	"github.com/portifyx/portifyx-goutils/pkg/logx"
)

type SuccessEnvelope struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type ErrorEnvelope struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// NewSuccessResponse formats a standard successful JSON response.
func NewSuccessResponse(c echo.Context, data any) error {
	return c.JSON(http.StatusOK, SuccessEnvelope{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// NewSuccessResponseWithStatus formats a successful JSON response with a custom status.
func NewSuccessResponseWithStatus(c echo.Context, status int, data any) error {
	return c.JSON(status, SuccessEnvelope{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// NewErrorResponse formats a standard typed error JSON response according to the Trade Platform standard.
func NewErrorResponse(c echo.Context, err error) error {
	ctx := c.Request().Context()

	if appErr, ok := apperror.IsAppError(err); ok {
		// Log 5xx errors at Error level, 4xx errors at Warn level
		if appErr.HTTPStatusCode >= 500 {
			logx.Error(ctx, appErr.Cause, "[!] internal downstream error in handler")
		} else {
			logx.Warnf(ctx, "[!] client validation/input error: [%d] %s", appErr.Code, appErr.Message)
		}

		return c.JSON(appErr.HTTPStatusCode, ErrorEnvelope{
			Code:    int(appErr.Code),
			Message: appErr.Message,
		})
	}

	// Unknown/unhandled internal error - never leak internal stack/details to caller
	logx.Error(ctx, err, "[!] unhandled server error")
	return c.JSON(http.StatusInternalServerError, ErrorEnvelope{
		Code:    int(apperror.CodeInternalServerError),
		Message: "Internal server error",
	})
}
