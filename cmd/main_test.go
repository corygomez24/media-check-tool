package main

import (
	"github.com/labstack/echo"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/stretchr/testify/assert"
)


func TestHealthCheck(t *testing.T){
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/healthz", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, checkHealth(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "200", rec.Body.String())
	}
}
