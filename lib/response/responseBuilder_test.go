package response_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"template/lib/response"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type resp struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func TestOK(t *testing.T) {

	req, err := http.NewRequestWithContext(context.TODO(), echo.GET, "/", strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	ctx := echo.New().NewContext(req, rec)

	expectResult := ctx.JSON(http.StatusOK, resp{
		Code:    http.StatusOK,
		Message: "Test OK",
		Data:    "Test",
	})

	result := response.ToJson(ctx).OK("Test", "Test OK")

	assert.Equal(t, expectResult, result)
}

func TestInternalServerError(t *testing.T) {

	req, err := http.NewRequestWithContext(context.TODO(), echo.GET, "/", strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	ctx := echo.New().NewContext(req, rec)

	expectResult := ctx.JSON(http.StatusOK, resp{
		Code:    http.StatusInternalServerError,
		Message: "Test Internal Server Error",
		Data:    nil,
	})

	result := response.ToJson(ctx).InternalServerError("Test Internal Server Error")

	assert.Equal(t, expectResult, result)
}

func TestBadRequest(t *testing.T) {

	req, err := http.NewRequestWithContext(context.TODO(), echo.GET, "/", strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	ctx := echo.New().NewContext(req, rec)

	expectResult := ctx.JSON(http.StatusOK, resp{
		Code:    http.StatusBadRequest,
		Message: "Test Bad Request",
		Data:    nil,
	})

	result := response.ToJson(ctx).BadRequest("Test Bad Request")

	assert.Equal(t, expectResult, result)
}
