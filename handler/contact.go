package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"blog-backend/api"
	"blog-backend/logger"
	"blog-backend/model"
)

// Submit a contact message
// (POST /contact)
func (h *Handler) PostContact(ctx echo.Context) error {
	var req api.PostContactJSONRequestBody
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	// check if the required fields are empty
	if req.Name == nil || req.Email == nil || req.Message == nil {
		return ctx.JSON(http.StatusBadRequest, "Name, Email, and Message are required")
	}

	// send Slack to admin
	err := model.SendSlack(ctx, *req.Name, string(*req.Email), *req.Message)
	if err != nil {
		logger.Println("Failed to send message to Slack: ", err)
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.JSON(http.StatusOK, "Message sent successfully")
}
