package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHealth(t *testing.T) {
	w := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	ginContext, _ := gin.CreateTestContext(w)
	ginContext.Request = request

	Health(ginContext)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"status":"ok"}`, w.Body.String())
}
