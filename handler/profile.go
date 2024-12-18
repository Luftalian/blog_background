package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Get profile information
// (GET /profile)
func (h *Handler) GetProfile(ctx echo.Context) error {
	return ctx.JSON(http.StatusNotImplemented, "not implemented")
}
