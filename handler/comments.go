package handler

import (
	"blog-backend/api"
	"blog-backend/model"
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// Get comments for an article
// (GET /comments)
func (h *Handler) GetComments(ctx echo.Context, params api.GetCommentsParams) error {
	// find comments by article id
	comments, err := h.Repo.GetCommentsByArticle(ctx, uuid.MustParse(params.ArticleId), nil)
	if err != nil {
		log.Println("GetCommentsByArticle Error: ", err)
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	if len(comments) == 0 {
		log.Println("No comments found")
		return ctx.JSON(http.StatusNotFound, "No comments found")
	}
	return ctx.JSON(http.StatusOK, comments)
}

// Post a comment
// (POST /comments)
func (h *Handler) PostComments(ctx echo.Context) error {
	log.Println("PostComments")
	var req api.PostCommentsJSONRequestBody
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}
	// UserIdが存在しない場合はIPアドレスでユーザーを特定する
	userId := uuid.Nil
	if req.UserId != nil {
		userId = uuid.MustParse(*req.UserId)
	} else {
		userIdFromDB, err := h.Repo.CheckIPAddressAndReturnUserIDWithUserName(ctx, req.Username)
		if err != nil {
			log.Println("CheckIPAddressAndReturnUserIDWithUserName Error: ", err)
			return ctx.JSON(http.StatusInternalServerError, err)
		}
		userId = userIdFromDB
	}
	log.Println("UserId: ", userId)
	// add comment
	err := h.Repo.CreateComment(ctx, model.Comment{
		ID:        uuid.New(),
		ArticleID: uuid.MustParse(req.ArticleId),
		AuthorID:  userId,
		Content:   req.Content,
		CreatedAt: time.Now(),
		Author:    sql.NullString{String: req.Username, Valid: req.Username != ""},
	})
	if err != nil {
		log.Println("CreateComment Error: ", err)
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.JSON(http.StatusCreated, "Comment created")
}

// Delete a comment
// (DELETE /comments/{id})
func (h *Handler) DeleteCommentsId(ctx echo.Context, id string) error {
	err := h.Repo.DeleteComment(ctx, uuid.MustParse(id))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.JSON(http.StatusOK, "Comment deleted")
}

// Edit a comment
// (PATCH /comments/{id})
func (h *Handler) PatchCommentsId(ctx echo.Context, id string) error {
	var req api.PatchCommentsIdJSONRequestBody
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}
	// update comment
	err := h.Repo.UpdateComment(ctx, model.Comment{
		ID:      uuid.MustParse(id),
		Content: *req.Content,
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.JSON(http.StatusOK, "Comment updated")
}
