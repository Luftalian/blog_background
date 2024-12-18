package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"blog-backend/handlers"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetArticles(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/articles?page=1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// ハンドラ関数を呼び出す
	err := handlers.GetArticles(c)

	// アサーション
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "articles")
}
