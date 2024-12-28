package handler

import (
	"blog-backend/api"
	"blog-backend/model"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// Get a list of tags
// (GET /tags)
func (h *Handler) GetTags(ctx echo.Context) error {
	tags, err := h.Repo.GetTagList(ctx.Request().Context())
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	if len(tags) == 0 {
		return ctx.JSON(http.StatusNotFound, "No tags found")
	}
	apiTags := make([]api.Tag, 0)
	for _, tag := range tags {
		id := tag.ID.String()
		apiTags = append(apiTags, api.Tag{
			Id:   &id,
			Name: &tag.Name,
		})
	}
	return ctx.JSON(http.StatusOK, apiTags)
}

// Create a new tag
// (POST /tags)
func (h *Handler) PostTags(ctx echo.Context) error {
	var req api.PostTagsJSONRequestBody
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}
	tagId := uuid.New()
	tag, err := h.Repo.AddTagItem(ctx.Request().Context(), model.TagItem{
		ID:   tagId,
		Name: *req.Name,
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	tagIdStr := tag.ID.String()
	return ctx.JSON(http.StatusCreated, api.Tag{
		Id:   &tagIdStr,
		Name: req.Name,
	})
}

// Add tags to an article
// (POST /tags/{articleId})
func (h *Handler) PostTagsArticleId(ctx echo.Context, articleId string) error {
	var req api.PostTagsArticleIdJSONRequestBody
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}
	article_id := uuid.MustParse(articleId)
	err := h.Repo.AddTag(ctx.Request().Context(), model.Tag{
		ID:        uuid.New(),
		ArticleID: article_id,
		Name:      *req.Tag.Name,
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.JSON(http.StatusCreated, "Tag added")
}
