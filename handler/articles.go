package handler

import (
	"blog-backend/api"
	"blog-backend/model"
	"database/sql"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// articles per page
const articlesPerPage = 10

// Get a list of articles
// (GET /articles)
func (h *Handler) GetArticles(ctx echo.Context, params api.GetArticlesParams) error {
	var articlesLength int = 0
	if params.Page != nil {
		articlesLength = articlesPerPage * *params.Page
	}
	category_id := uuid.Nil
	categoryName := ""
	tag_id := uuid.Nil
	tagName := ""
	if params.Tag != nil {
		tag_id = uuid.MustParse(*params.Tag)
		tag, err := h.Repo.GetTagItemsByID(ctx, tag_id)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, err)
		}
		tagName = tag.Name
	}
	if params.Category != nil {
		category_id = uuid.MustParse(*params.Category)
		category, err := h.Repo.GetCategoryNameByID(ctx, category_id)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, err)
		}
		categoryName = category.Name
	}
	orderby := "created_at"
	order := "desc"
	if params.Orderby != nil && *params.Orderby == api.ViewCount {
		orderby = "view_count"
	}
	if params.Order != nil && *params.Order == api.Asc {
		order = "asc"
	}

	articles, err := h.Repo.GetArticlesByCategoryTagSearch(ctx, &category_id, &tag_id, params.Search, &articlesLength, orderby, order)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	if len(articles) == 0 {
		return ctx.JSON(http.StatusNotFound, "No articles found")
	}

	apiArticles, err := convertArticlesToAPIArticles(ctx, articles, h.Repo)
	if err != nil {
		log.Println("Convert articles error: ", err)
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, api.ArticleResponse{
		Articles: &apiArticles,
		Condition: &struct {
			Category *string `json:"category,omitempty"`
			Order    *string `json:"order,omitempty"`
			Orderby  *string `json:"orderby,omitempty"`
			Search   *string `json:"search,omitempty"`
			Tag      *string `json:"tag,omitempty"`
		}{
			Category: &categoryName,
			Order:    &orderby,
			Orderby:  &order,
			Tag:      &tagName,
			Search:   params.Search,
		},
	})
}

// Create a new article
// (POST /articles)
func (h *Handler) PostArticles(ctx echo.Context) error {
	var req api.PostArticlesJSONRequestBody
	if err := ctx.Bind(&req); err != nil {
		log.Println("Bind Error: ", err)
		return ctx.JSON(http.StatusBadRequest, err)
	}

	if !req.IsAdmin {
		return ctx.JSON(http.StatusUnauthorized, "Unauthorized")
	}

	// add category
	categoryID := uuid.MustParse(req.Category)
	category, err := h.Repo.AddCategoryID(ctx, categoryID)
	if err != nil {
		log.Println("AddCategoryID Error: ", err)
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	// add tags
	tags := make([]model.TagItem, 0)
	for _, tag := range *req.Tags {
		tags = append(tags, model.TagItem{
			ID: uuid.MustParse(tag),
		})
	}

	// get tags name
	for i, tag := range tags {
		log.Println("TagID: ", tag.ID)
		tagName, err := h.Repo.GetTagItemsByID(ctx, tag.ID)
		if err != nil {
			log.Println("GetTagNameByID Error: ", err)
			return ctx.JSON(http.StatusInternalServerError, err)
		}
		tags[i].Name = tagName.Name
	}

	articleId := uuid.New()

	// // req.AuthorIdが存在しない場合はIPアドレスでユーザーを特定する
	// userId := uuid.Nil
	// if req.AuthorId == nil {
	// 	newUserId, err := h.Repo.CheckIPAddressAndReturnUserIDWithUserName(ctx, req.Author)
	// 	if err != nil {
	// 		log.Println("CheckIPAddressAndReturnUserID Error: ", err)
	// 		return ctx.JSON(http.StatusInternalServerError, err)
	// 	}
	// 	userId = newUserId
	// } else {
	// 	userId = uuid.MustParse(*req.AuthorId)
	// 	// update username
	// 	user, err := h.Repo.GetUserByID(ctx, userId)
	// 	if err != nil {
	// 		log.Println("GetUserByID Error: ", err)
	// 		return ctx.JSON(http.StatusInternalServerError, err)
	// 	}
	// 	if user.Username.String != req.Author {
	// 		user.Username.String = req.Author
	// 		err = h.Repo.UpdateUser(ctx, user)
	// 		if err != nil {
	// 			log.Println("UpdateUser Error: ", err)
	// 			return ctx.JSON(http.StatusInternalServerError, err)
	// 		}
	// 	}
	// }

	// adminのuserIdを取得し、usernameをreq.Authorに設定
	users, err := h.Repo.GetAdminUsers(ctx)
	userId := uuid.Nil
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			// errが存在して、かつエラーが"sql: no rows in result set"でない場合はエラーを返す
			log.Println("GetAdminUserId Error: ", err)
			return ctx.JSON(http.StatusInternalServerError, err)
		} else {
			// admin userが存在しない場合は新規作成
			newUserId, err := h.Repo.CheckIPAddressAndReturnUserIDWithUserName(ctx, req.Author)
			if err != nil {
				log.Println("CheckIPAddressAndReturnUserID Error: ", err)
				return ctx.JSON(http.StatusInternalServerError, err)
			}
			userId = newUserId
		}
	} else if len(users) > 1 {
		// admin userが複数いる場合はエラーを返す
		log.Println("Multiple admin users found")
		return ctx.JSON(http.StatusInternalServerError, "Multiple admin users found")
	} else if len(users) == 0 {
		// admin userが存在しない場合は新規作成
		newUserId, err := h.Repo.CheckIPAddressAndReturnUserIDWithUserName(ctx, req.Author)
		if err != nil {
			log.Println("CheckIPAddressAndReturnUserID Error: ", err)
			return ctx.JSON(http.StatusInternalServerError, err)
		}
		userId = newUserId
	} else {
		// admin userが1つの場合はそのuserIdを取得し、usernameをreq.Authorに設定
		adminUser := users[0]
		if adminUser.Username.String != req.Author {
			adminUser.Username.String = req.Author
			err = h.Repo.UpdateUser(ctx, adminUser)
			if err != nil {
				log.Println("UpdateUser Error: ", err)
				return ctx.JSON(http.StatusInternalServerError, err)
			}
		}
		userId = adminUser.ID
	}

	newArticle := model.Article{
		ID:         articleId,
		Title:      req.Title,
		Content:    req.Content,
		AuthorID:   userId,
		CategoryID: category.ID,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		ViewCount:  sql.NullInt64{Int64: 0, Valid: true},
	}

	article, err := h.Repo.CreateArticle(ctx, newArticle)
	if err != nil {
		log.Println("CreateArticle Error: ", err)
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	err = h.Repo.AddTagPairsByArticle(ctx, articleId, tags)
	if err != nil {
		log.Println("AddTagPairsByArticle Error: ", err)
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	authorName, err := h.Repo.GetUserNameById(ctx, userId)
	if err != nil {
		log.Println("GetUserNameById Error: ", err)
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	authorNameForThumbnail := authorName.Username.String
	if authorNameForThumbnail == "" {
		authorNameForThumbnail = "Luftalian"
	}
	idStr := article.ID.String()
	imageUrl, err := h.Config.HandleThumbnailGeneration(ctx, newArticle, tags, category.Name, authorName.Username.String)
	if err != nil {
		log.Println("HandleThumbnailGeneration Error: ", err)
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	if imageUrl != "" {
		err = h.Repo.UpdateArticleImageURL(ctx, articleId, imageUrl)
		if err != nil {
			log.Println("UpdateArticleImage Error: ", err)
			return ctx.JSON(http.StatusInternalServerError, err)
		}
	}

	{
		limit := 5
		articles, err := h.Repo.GetArticlesByCategoryTagSearch(ctx, &uuid.Nil, &uuid.Nil, nil, &limit, "created_at", "desc")
		if err != nil {
			// エラーログを出力
			ctx.Logger().Errorf("Failed to fetch articles: %v", err)
			// return ctx.JSON(http.StatusInternalServerError, "Failed to fetch articles")
		}

		// articlesのsortを行う。created_atの降順でソートする。
		sort.Slice(articles, func(i, j int) bool {
			return articles[i].CreatedAt.After(articles[j].CreatedAt)
		})

		// RSSフィードの設定
		err = h.Config.RSSmaker(ctx, articles)
		if err != nil {
			ctx.Logger().Errorf("Failed to generate RSS feed: %v", err)
			// return ctx.JSON(http.StatusInternalServerError, err)
		}
	}

	return ctx.JSON(http.StatusCreated, api.Article{
		Id: &idStr,
	})
}

// Delete an article
// (DELETE /articles/{id})
func (h *Handler) DeleteArticlesId(ctx echo.Context, id string) error {
	// convert id to uuid
	articleId, err := uuid.Parse(id)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}
	err = h.Repo.DeleteArticle(ctx, articleId)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.JSON(http.StatusNoContent, nil)
}

// Get article details
// (GET /articles/{id})
func (h *Handler) GetArticlesId(ctx echo.Context, id string) error {
	saveAnalysis := func(ctx echo.Context, articleId uuid.UUID, api string, isError bool) error {
		err := h.Repo.CreateAnalysis(ctx, model.Analysis{
			ID:         uuid.New(),
			Timestamp:  time.Now(),
			ArticleID:  articleId,
			IpAddress:  ctx.RealIP(),
			SearchWord: "",
			API:        api,
			IsError:    isError,
		})
		if err != nil {
			return err
		}
		if !isError {
			log.Println("SaveAnalysis Error: ")
			err = h.Repo.SaveViewCount(ctx, articleId)
			if err != nil {
				return err
			}
			// get article view count
			article, err := h.Repo.GetArticleByID(ctx, articleId)
			if err != nil {
				return err
			}
			log.Println("ViewCount: ", article.ViewCount)
		}
		return nil
	}

	articleId, err := uuid.Parse(id)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	err = h.Repo.SaveViewCount(ctx, articleId)
	if err != nil {
		errSaveAnalysis := saveAnalysis(ctx, articleId, "GetArticlesId", true)
		if errSaveAnalysis != nil {
			log.Println("SaveAnalysis Error: ", errSaveAnalysis)
		}
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	article, err := h.Repo.GetArticleByID(ctx, articleId)
	if err != nil {
		errSaveAnalysis := saveAnalysis(ctx, articleId, "GetArticlesId", true)
		if errSaveAnalysis != nil {
			log.Println("SaveAnalysis Error: ", errSaveAnalysis)
		}
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	articles := []model.Article{article}
	apiArticle, err := convertArticlesToAPIArticles(ctx, articles, h.Repo)
	if err != nil {
		errSaveAnalysis := saveAnalysis(ctx, articleId, "GetArticlesId", true)
		if errSaveAnalysis != nil {
			log.Println("SaveAnalysis Error: ", errSaveAnalysis)
		}
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	errSaveAnalysis := saveAnalysis(ctx, articleId, "GetArticlesId", false)
	if errSaveAnalysis != nil {
		log.Println("SaveAnalysis Error: ", errSaveAnalysis)
	}
	return ctx.JSON(http.StatusOK, apiArticle[0])
}

// Update an article
// (PATCH /articles/{id})
func (h *Handler) PatchArticlesId(ctx echo.Context, id string) error {
	articleId, err := uuid.Parse(id)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	var req api.PatchArticlesIdJSONRequestBody
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	// add category
	categoryId, err := h.Repo.AddCategory(ctx, *req.Category)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	// add tags
	tags := make([]model.Tag, 0)
	for _, tag := range *req.Tags {
		tags = append(tags, model.Tag{
			ID: uuid.MustParse(tag),
		})
	}
	err = h.Repo.AddTags(ctx, tags)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	article, err := h.Repo.UpdateArticle(ctx, model.Article{
		ID:         articleId,
		Title:      *req.Title,
		Content:    *req.Content,
		AuthorID:   uuid.MustParse(*req.AuthorId),
		CategoryID: categoryId.ID,
		UpdatedAt:  time.Now(),
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.JSON(http.StatusOK, article)
}

// Get archive of articles
// (GET /articles/archive)
func (h *Handler) GetArticlesArchive(ctx echo.Context) error {
	return ctx.JSON(http.StatusNotImplemented, "Not implemented")
}

// Get articles by author
// (GET /articles/author/{authorId})
func (h *Handler) GetArticlesAuthorAuthorId(ctx echo.Context, authorId string) error {
	author, err := h.Repo.GetUserNameById(ctx, uuid.MustParse(authorId))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	articles, err := h.Repo.GetArticlesByAuthor(ctx, uuid.MustParse(authorId), nil)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	apiArticles, err := convertArticlesToAPIArticles(ctx, articles, h.Repo)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	authorName := author.Username.String
	return ctx.JSON(http.StatusOK, api.ArticleByAuthor{
		Articles: &apiArticles,
		Author:   &authorName,
		AuthorId: &authorId,
	})
}
