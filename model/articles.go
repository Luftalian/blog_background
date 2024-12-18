package model

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Article struct {
	ID         uuid.UUID      `db:"id"`
	Title      string         `db:"title"`
	Content    string         `db:"content"`
	AuthorID   uuid.UUID      `db:"author_id"`
	CategoryID uuid.UUID      `db:"category_id"`
	CreatedAt  time.Time      `db:"created_at"`
	UpdatedAt  time.Time      `db:"updated_at"`
	ViewCount  sql.NullInt64  `db:"view_count"`
	ImageURL   sql.NullString `db:"image_url"`
}

func (repo *Repository) GetArticleByID(ctx echo.Context, id uuid.UUID) (Article, error) {
	var article Article
	err := repo.db.GetContext(ctx.Request().Context(), &article, "SELECT * FROM articles WHERE id = ?", id)
	return article, err
}

func (repo *Repository) GetArticles(ctx echo.Context, limitNumber *int) ([]Article, error) {
	var articles []Article
	if limitNumber != nil {
		err := repo.db.SelectContext(ctx.Request().Context(), &articles, "SELECT * FROM articles LIMIT ?", limitNumber)
		return articles, err
	} else {
		err := repo.db.SelectContext(ctx.Request().Context(), &articles, "SELECT * FROM articles")
		return articles, err
	}
}

func (repo *Repository) CreateArticle(ctx echo.Context, article Article) (Article, error) {
	_, err := repo.db.NamedExecContext(ctx.Request().Context(), "INSERT INTO articles (id, title, content, author_id, category_id, created_at, updated_at) VALUES (:id, :title, :content, :author_id, :category_id, :created_at, :updated_at)", article)
	return article, err
}

func (repo *Repository) UpdateArticle(ctx echo.Context, article Article) (Article, error) {
	_, err := repo.db.NamedExecContext(ctx.Request().Context(), "UPDATE articles SET title = :title, content = :content, author_id = :author_id, category_id = :category_id, updated_at = :updated_at WHERE id = :id", article)
	return article, err
}

func (repo *Repository) UpdateArticleImageURL(ctx echo.Context, id uuid.UUID, imageURL string) error {
	_, err := repo.db.ExecContext(ctx.Request().Context(), "UPDATE articles SET image_url = ? WHERE id = ?", imageURL, id)
	return err
}

func (repo *Repository) DeleteArticle(ctx echo.Context, id uuid.UUID) error {
	_, err := repo.db.ExecContext(ctx.Request().Context(), "DELETE FROM articles WHERE id = ?", id)
	return err
}

func (repo *Repository) GetArticlesByCategory(ctx echo.Context, category_id uuid.UUID, limitNumber *int) ([]Article, error) {
	var articles []Article
	if limitNumber != nil {
		err := repo.db.SelectContext(ctx.Request().Context(), &articles, "SELECT * FROM articles WHERE category_id = ? LIMIT ?", category_id, limitNumber)
		return articles, err
	} else {
		err := repo.db.SelectContext(ctx.Request().Context(), &articles, "SELECT * FROM articles WHERE category_id = ?", category_id)
		return articles, err
	}
}

func (repo *Repository) GetArticlesByAuthor(ctx echo.Context, author uuid.UUID, limitNumber *int) ([]Article, error) {
	var articles []Article
	if limitNumber != nil {
		err := repo.db.SelectContext(ctx.Request().Context(), &articles, "SELECT * FROM articles WHERE author_id = ? LIMIT ?", author, limitNumber)
		return articles, err
	} else {
		err := repo.db.SelectContext(ctx.Request().Context(), &articles, "SELECT * FROM articles WHERE author_id = ?", author)
		return articles, err
	}
}

func (repo *Repository) GetArticlesByDate(ctx echo.Context, start time.Time, end time.Time, limitNumber *int) ([]Article, error) {
	var articles []Article
	if limitNumber != nil {
		err := repo.db.SelectContext(ctx.Request().Context(), &articles, "SELECT * FROM articles WHERE created_at BETWEEN ? AND ? LIMIT ?", start, end, limitNumber)
		return articles, err
	} else {
		err := repo.db.SelectContext(ctx.Request().Context(), &articles, "SELECT * FROM articles WHERE created_at BETWEEN ? AND ?", start, end)
		return articles, err
	}
}

func (repo *Repository) GetArticlesByCategoryTagSearch(ctx echo.Context, category_id *uuid.UUID, tag_id *uuid.UUID, search *string, limitNumber *int, orderby string, order string) ([]Article, error) {
	var articles []Article
	fmt.Println(*limitNumber != 0)
	fmt.Println(*category_id != uuid.Nil)
	fmt.Println(*tag_id != uuid.Nil)
	fmt.Println(search != nil)
	fmt.Println("nil conditions:", *limitNumber != 0, *category_id != uuid.Nil, *tag_id != uuid.Nil, search != nil)
	if *limitNumber != 0 {
		log.Println("limitNumber:", *limitNumber)
		if *category_id != uuid.Nil && *tag_id != uuid.Nil && search != nil {
			if order == "asc" {
				log.Println("condition 1-1")
				err := repo.db.SelectContext(ctx.Request().Context(), &articles, "SELECT * FROM articles WHERE category_id = ? AND id IN (SELECT article_id FROM article_tags WHERE tag_id = ?) AND (title LIKE ? OR content LIKE ?) ORDER BY ? ASC LIMIT ?", category_id, tag_id, "%"+*search+"%", "%"+*search+"%", orderby, limitNumber)
				return articles, err
			}
			log.Println("condition 1")
			err := repo.db.SelectContext(ctx.Request().Context(), &articles, "SELECT * FROM articles WHERE category_id = ? AND id IN (SELECT article_id FROM article_tags WHERE tag_id = ?) AND (title LIKE ? OR content LIKE ?) ORDER BY ? DESC LIMIT ?", category_id, tag_id, "%"+*search+"%", "%"+*search+"%", orderby, limitNumber)
			return articles, err
		} else if *category_id != uuid.Nil && *tag_id != uuid.Nil {
			if order == "asc" {
				log.Println("condition 2-1")
				err := repo.db.SelectContext(ctx.Request().Context(), &articles, "SELECT * FROM articles WHERE category_id = ? AND id IN (SELECT article_id FROM article_tags WHERE tag_id = ?) ORDER BY ? ASC LIMIT ?", category_id, tag_id, orderby, limitNumber)
				return articles, err
			}
			log.Println("condition 2")
			err := repo.db.SelectContext(ctx.Request().Context(), &articles, "SELECT * FROM articles WHERE category_id = ? AND id IN (SELECT article_id FROM article_tags WHERE tag_id = ?) ORDER BY ? DESC LIMIT ?", category_id, tag_id, orderby, limitNumber)
			return articles, err
		} else if *category_id != uuid.Nil && search != nil {
			if order == "asc" {
				log.Println("condition 3-1")
				err := repo.db.SelectContext(ctx.Request().Context(), &articles, "SELECT * FROM articles WHERE category_id = ? AND (title LIKE ? OR content LIKE ?) ORDER BY ? ASC LIMIT ?", category_id, "%"+*search+"%", "%"+*search+"%", orderby, limitNumber)
				return articles, err
			}
			log.Println("condition 3")
			err := repo.db.SelectContext(ctx.Request().Context(), &articles, "SELECT * FROM articles WHERE category_id = ? AND (title LIKE ? OR content LIKE ?) ORDER BY ? DESC LIMIT ?", category_id, "%"+*search+"%", "%"+*search+"%", orderby, limitNumber)
			return articles, err
		} else if *tag_id != uuid.Nil && search != nil {
			if order == "asc" {
				log.Println("condition 4-1")
				err := repo.db.SelectContext(ctx.Request().Context(), &articles, "SELECT * FROM articles WHERE id IN (SELECT article_id FROM article_tags WHERE tag_id = ?) AND (title LIKE ? OR content LIKE ?) ORDER BY ? ASC LIMIT ?", tag_id, "%"+*search+"%", "%"+*search+"%", orderby, limitNumber)
				return articles, err
			}
			log.Println("condition 4")
			err := repo.db.SelectContext(ctx.Request().Context(), &articles, "SELECT * FROM articles WHERE id IN (SELECT article_id FROM article_tags WHERE tag_id = ?) AND (title LIKE ? OR content LIKE ?) ORDER BY ? DESC LIMIT ?", tag_id, "%"+*search+"%", "%"+*search+"%", orderby, limitNumber)
			return articles, err
		} else if *category_id != uuid.Nil {
			if order == "asc" {
				log.Println("condition 5-1")
				err := repo.db.SelectContext(ctx.Request().Context(), &articles, "SELECT * FROM articles WHERE category_id = ? ORDER BY ? ASC LIMIT ?", category_id, orderby, limitNumber)
				return articles, err
			}
			log.Println("condition 5")
			err := repo.db.SelectContext(ctx.Request().Context(), &articles, "SELECT * FROM articles WHERE category_id = ? ORDER BY ? DESC LIMIT ?", category_id, orderby, limitNumber)
			return articles, err
		} else if *tag_id != uuid.Nil {
			if order == "asc" {
				log.Println("condition 6-1")
				err := repo.db.SelectContext(ctx.Request().Context(), &articles, "SELECT * FROM articles WHERE id IN (SELECT article_id FROM article_tags WHERE tag_id = ?) ORDER BY ? ASC LIMIT ?", tag_id, orderby, limitNumber)
				return articles, err
			}
			log.Println("condition 6")
			err := repo.db.SelectContext(ctx.Request().Context(), &articles, "SELECT * FROM articles WHERE id IN (SELECT article_id FROM article_tags WHERE tag_id = ?) ORDER BY ? DESC LIMIT ?", tag_id, orderby, limitNumber)
			return articles, err
		} else if search != nil {
			if order == "asc" {
				log.Println("condition 7-1")
				err := repo.db.SelectContext(ctx.Request().Context(), &articles, "SELECT * FROM articles WHERE title LIKE ? OR content LIKE ? ORDER BY ? ASC LIMIT ?", "%"+*search+"%", "%"+*search+"%", orderby, limitNumber)
				return articles, err
			}
			log.Println("condition 7")
			err := repo.db.SelectContext(ctx.Request().Context(), &articles, "SELECT * FROM articles WHERE title LIKE ? OR content LIKE ? ORDER BY ? DESC LIMIT ?", "%"+*search+"%", "%"+*search+"%", orderby, limitNumber)
			return articles, err
		} else {
			if order == "asc" {
				log.Println("condition 8-1")
				err := repo.db.SelectContext(ctx.Request().Context(), &articles, "SELECT * FROM articles ORDER BY ? ASC LIMIT ?", orderby, limitNumber)
				return articles, err
			}
			log.Println("condition 8")
			err := repo.db.SelectContext(ctx.Request().Context(), &articles, "SELECT * FROM articles ORDER BY ? DESC LIMIT ?", orderby, limitNumber)
			log.Println("err:", err)
			return articles, err
		}
	} else {
		if *category_id != uuid.Nil && *tag_id != uuid.Nil && search != nil {
			if order == "asc" {
				log.Println("condition 9-1")
				err := repo.db.SelectContext(ctx.Request().Context(), &articles, "SELECT * FROM articles WHERE category_id = ? AND id IN (SELECT article_id FROM article_tags WHERE tag_id = ?) AND (title LIKE ? OR content LIKE ?) ORDER BY ? ASC", category_id, tag_id, "%"+*search+"%", "%"+*search+"%", orderby)
				return articles, err
			}
			log.Println("condition 9")
			err := repo.db.SelectContext(ctx.Request().Context(), &articles, "SELECT * FROM articles WHERE category_id = ? AND id IN (SELECT article_id FROM article_tags WHERE tag_id = ?) AND (title LIKE ? OR content LIKE ?) ORDER BY ? DESC", category_id, tag_id, "%"+*search+"%", "%"+*search+"%", orderby)
			return articles, err
		} else if *category_id != uuid.Nil && *tag_id != uuid.Nil {
			if order == "asc" {
				log.Println("condition 10-1")
				err := repo.db.SelectContext(ctx.Request().Context(), &articles, "SELECT * FROM articles WHERE category_id = ? AND id IN (SELECT article_id FROM article_tags WHERE tag_id = ?) ORDER BY ? ASC", category_id, tag_id, orderby)
				return articles, err
			}
			log.Println("condition 10")
			err := repo.db.SelectContext(ctx.Request().Context(), &articles, "SELECT * FROM articles WHERE category_id = ? AND id IN (SELECT article_id FROM article_tags WHERE tag_id = ?) ORDER BY ? DESC", category_id, tag_id, orderby)
			return articles, err
		}
		if *category_id != uuid.Nil && search != nil {
			if order == "asc" {
				log.Println("condition 11-1")
				err := repo.db.SelectContext(ctx.Request().Context(), &articles, "SELECT * FROM articles WHERE category_id = ? AND (title LIKE ? OR content LIKE ?) ORDER BY ? ASC", category_id, "%"+*search+"%", "%"+*search+"%", orderby)
				return articles, err
			}
			log.Println("condition 11")
			err := repo.db.SelectContext(ctx.Request().Context(), &articles, "SELECT * FROM articles WHERE category_id = ? AND (title LIKE ? OR content LIKE ?) ORDER BY ? DESC", category_id, "%"+*search+"%", "%"+*search+"%", orderby)
			return articles, err
		}
		if *tag_id != uuid.Nil && search != nil {
			if order == "asc" {
				log.Println("condition 12-1")
				err := repo.db.SelectContext(ctx.Request().Context(), &articles, "SELECT * FROM articles WHERE id IN (SELECT article_id FROM article_tags WHERE tag_id = ?) AND (title LIKE ? OR content LIKE ?) ORDER BY ? ASC", tag_id, "%"+*search+"%", "%"+*search+"%", orderby)
				return articles, err
			}
			log.Println("condition 12")
			err := repo.db.SelectContext(ctx.Request().Context(), &articles, "SELECT * FROM articles WHERE id IN (SELECT article_id FROM article_tags WHERE tag_id = ?) AND (title LIKE ? OR content LIKE ?) ORDER BY ? DESC", tag_id, "%"+*search+"%", "%"+*search+"%", orderby)
			return articles, err
		}
		if *category_id != uuid.Nil {
			if order == "asc" {
				log.Println("condition 13-1")
				err := repo.db.SelectContext(ctx.Request().Context(), &articles, "SELECT * FROM articles WHERE category_id = ? ORDER BY ? ASC", category_id, orderby)
				return articles, err
			}
			log.Println("condition 13")
			err := repo.db.SelectContext(ctx.Request().Context(), &articles, "SELECT * FROM articles WHERE category_id = ? ORDER BY ? DESC", category_id, orderby)
			return articles, err
		}
		if *tag_id != uuid.Nil {
			if order == "asc" {
				log.Println("condition 14-1")
				err := repo.db.SelectContext(ctx.Request().Context(), &articles, "SELECT * FROM articles WHERE id IN (SELECT article_id FROM article_tags WHERE tag_id = ?) ORDER BY ? ASC", tag_id, orderby)
				return articles, err
			}
			log.Println("condition 14")
			err := repo.db.SelectContext(ctx.Request().Context(), &articles, "SELECT * FROM articles WHERE id IN (SELECT article_id FROM article_tags WHERE tag_id = ?) ORDER BY ? DESC", tag_id, orderby)
			return articles, err
		}
		if search != nil {
			if order == "asc" {
				log.Println("condition 15-1")
				err := repo.db.SelectContext(ctx.Request().Context(), &articles, "SELECT * FROM articles WHERE title LIKE ? OR content LIKE ? ORDER BY ? ASC", "%"+*search+"%", "%"+*search+"%", orderby)
				return articles, err
			}
			log.Println("condition 15")
			err := repo.db.SelectContext(ctx.Request().Context(), &articles, "SELECT * FROM articles WHERE title LIKE ? OR content LIKE ? ORDER BY ? DESC", "%"+*search+"%", "%"+*search+"%", orderby)
			return articles, err
		}
		if order == "asc" {
			log.Println("condition 16-1")
			err := repo.db.SelectContext(ctx.Request().Context(), &articles, "SELECT * FROM articles ORDER BY ? ASC", orderby)
			return articles, err
		}
		err := repo.db.SelectContext(ctx.Request().Context(), &articles, "SELECT * FROM articles ORDER BY ? DESC", orderby)
		return articles, err
	}
}

func (repo *Repository) SaveViewCount(ctx echo.Context, id uuid.UUID) error {
	_, err := repo.db.ExecContext(ctx.Request().Context(), "UPDATE articles SET view_count = view_count + 1 WHERE id = ?", id)
	log.Println("view count error:", err)
	return err
}
