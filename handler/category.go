package handler

import (
	"blog-backend/api"
	"blog-backend/model"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// Get a list of categories
// (GET /categories)
func (h *Handler) GetCategories(ctx echo.Context) error {
	categories, err := h.Repo.GetCategories(ctx.Request().Context(), nil)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	if len(categories) == 0 {
		return ctx.JSON(http.StatusNotFound, "No categories found")
	}
	var apiCategories []api.Category
	for _, category := range categories {
		id := category.ID.String()
		apiCategories = append(apiCategories, api.Category{
			Id:   &id,
			Name: &category.Name,
		})
	}
	return ctx.JSON(http.StatusOK, apiCategories)
}

// Create a new category
// (POST /categories)
func (h *Handler) PostCategories(ctx echo.Context) error {
	var req api.PostCategoriesJSONRequestBody
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}
	categoryId := uuid.New()
	err := h.Repo.CreateCategory(ctx.Request().Context(), model.Category{
		ID:   categoryId,
		Name: *req.Name,
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	categoryIdStr := categoryId.String()
	return ctx.JSON(http.StatusCreated, api.Category{
		Id:   &categoryIdStr,
		Name: req.Name,
	})
}
