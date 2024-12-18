package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Get RSS feed
// (GET /rss)
func (h *Handler) GetRss(ctx echo.Context) error {
	// limit := 5
	// articles, err := h.Repo.GetArticles(ctx, &limit)
	// if err != nil {
	// 	// エラーログを出力
	// 	ctx.Logger().Errorf("Failed to fetch articles: %v", err)
	// 	return ctx.JSON(http.StatusInternalServerError, "Failed to fetch articles")
	// }

	// // RSSフィードの設定
	// err = h.Config.RSSmaker(ctx, articles)
	// if err != nil {
	// 	return ctx.JSON(http.StatusInternalServerError, err)
	// }

	// 処理が成功したことを示すレスポンスとしてRSSフィードのURLを返す
	return ctx.JSON(http.StatusOK, h.Config.BaseURL+"/rss/rss.xml")
}
