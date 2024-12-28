package handler

import (
	"blog-backend/api"
	"blog-backend/logger"
	"blog-backend/model"
	"database/sql"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func convertNullStringToStringPoint(s sql.NullString) *string {
	if s.Valid {
		return &s.String
	}
	return nil
}

func convertNullInt64ToIntPoint(i sql.NullInt64, defaltValue *int) *int {
	if i.Valid {
		val := int(i.Int64)
		return &val
	}
	return defaltValue
}

func convertArticlesToAPIArticles(ctx echo.Context, articles []model.Article, repo *model.Repository) ([]api.Article, error) {
	var returnArticles []api.Article
	authorNameIdMap := make(map[uuid.UUID]string)
	catergoryNameIdMap := make(map[uuid.UUID]api.Category)
	for _, article := range articles {

		id := article.ID.String()

		if _, ok := authorNameIdMap[article.AuthorID]; !ok {
			author, err := repo.GetUserNameById(ctx, article.AuthorID)
			if err != nil {
				logger.Println("error getting author name", err)
				return nil, err
			}
			authorNameIdMap[article.AuthorID] = author.Username.String
		}
		author := authorNameIdMap[article.AuthorID]

		if _, ok := catergoryNameIdMap[article.CategoryID]; !ok {
			category, err := repo.GetCategoryNameByID(ctx, article.CategoryID)
			if err != nil {
				logger.Println("error getting category name", err)
				return nil, err
			}
			category_id := category.ID.String()
			catergoryNameIdMap[article.CategoryID] = api.Category{
				Id:   &category_id,
				Name: &category.Name,
			}
		}
		category := catergoryNameIdMap[article.CategoryID]

		tags, err := repo.GetTagsByArticle(ctx, article.ID, nil)
		if err != nil {
			logger.Println("error getting tags", err)
			return nil, err
		}
		var tagList []api.Tag
		for _, tag := range tags {
			tag_id := tag.ID.String()
			tagList = append(tagList, api.Tag{
				Id:   &tag_id,
				Name: &tag.Name,
			})
		}

		likeCount, err := repo.GetLikesCountByArticle(ctx, article.ID)
		if err != nil {
			logger.Println("error getting like count", err)
			return nil, err
		}

		zeroViewCount := 0

		authorIdStr := article.AuthorID.String()
		returnArticles = append(returnArticles, api.Article{
			Author:    &author,
			AuthorId:  &authorIdStr,
			Category:  &category,
			Content:   &article.Content,
			ImageUrl:  convertNullStringToStringPoint(article.ImageURL),
			ViewCount: convertNullInt64ToIntPoint(article.ViewCount, &zeroViewCount),
			CreatedAt: &article.CreatedAt,
			Id:        &id,
			LikeCount: &likeCount,
			Tags:      &tagList,
			Title:     &article.Title,
			UpdatedAt: &article.UpdatedAt,
		})
	}
	return returnArticles, nil
}
