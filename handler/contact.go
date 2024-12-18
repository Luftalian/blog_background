package handler

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"

	"blog-backend/api"
	"blog-backend/model"
)

// Submit a contact message
// (POST /contact)
func (h *Handler) PostContact(ctx echo.Context) error {
	var req api.PostContactJSONRequestBody
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	// email to admin
	err := model.SendEmail(model.ContactForm{
		Name:    *req.Name,
		Email:   string(*req.Email),
		Message: *req.Message,
	})
	if err != nil {
		log.Println("SendEmail Error: ", err)
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, "Message sent")
}
