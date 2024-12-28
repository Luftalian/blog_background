package handler

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"blog-backend/api"
	"blog-backend/logger"
	"blog-backend/model"
)

// Add a like to an article
// (POST /likes)
func (h *Handler) PostLikes(ctx echo.Context) error {
	var req api.PostLikesJSONRequestBody
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, "invalid request")
	}
	// user idが存在しない場合はIPアドレスでユーザーを特定する
	userID := uuid.Nil
	if req.UserId != nil {
		userID = uuid.MustParse(*req.UserId)
	} else {
		userIDFromDB, err := h.Repo.CheckIPAddressAndReturnUserID(ctx.Request().Context(), ctx.RealIP())
		if err != nil {
			logger.Println("CheckIPAddressAndReturnUserID Error: ", err)
			return ctx.JSON(http.StatusInternalServerError, err)
		}
		userID = userIDFromDB
	}
	// add like
	err := h.Repo.CreateLike(ctx.Request().Context(), model.Like{
		ID:        uuid.New(),
		ArticleID: uuid.MustParse(req.ArticleId),
		UserID:    userID,
		CreatedAt: time.Now(),
	})
	if err != nil {
		logger.Println("CreateLike Error: ", err)
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.JSON(http.StatusCreated, "like added")
}

// Get likes for an article
// (GET /likes)
func (h *Handler) GetLikes(ctx echo.Context, params api.GetLikesParams) error {
	// find likes by article id
	likes, err := h.Repo.GetLikesByArticle(ctx.Request().Context(), uuid.MustParse(params.ArticleId), nil)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	// CheckIPAddressAndReturnUserID and check if the user has liked the article
	userID, err := h.Repo.CheckIPAddressAndReturnUserID(ctx.Request().Context(), ctx.RealIP())
	if err != nil {
		logger.Println("CheckIPAddressAndReturnUserID Error: ", err)
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	liked := false
	for _, like := range likes {
		if like.UserID == userID {
			liked = true
		}
	}
	likeCount := len(likes)
	userIDStr := userID.String()
	return ctx.JSON(http.StatusOK, api.LikeReturn{
		ArticleId: &params.ArticleId,
		LikeCount: &likeCount,
		Liked:     &liked,
		UserId:    &userIDStr,
	})

}
