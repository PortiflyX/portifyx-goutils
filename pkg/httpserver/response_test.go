package httpserver

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/portifyx/portifyx-goutils/pkg/apperror"
	"github.com/stretchr/testify/assert"
)

func TestSuccessResponse(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := NewSuccessResponse(c, map[string]string{"foo": "bar"})
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp SuccessEnvelope
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, 0, resp.Code)
	assert.Equal(t, "success", resp.Message)
}

func TestErrorResponse(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	appErr := apperror.NewMissingRequiredField("user_id")
	err := NewErrorResponse(c, appErr)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)

	var resp ErrorEnvelope
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, int(apperror.CodeMissingRequiredField), resp.Code)
}
