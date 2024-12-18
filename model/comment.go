package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Comment struct {
	ID        uuid.UUID      `db:"id"`
	ArticleID uuid.UUID      `db:"article_id"`
	AuthorID  uuid.UUID      `db:"author_id"`
	Author    sql.NullString `db:"username"` // このフィールドはDBには存在しない。Userテーブルから取得してくる
	Content   string         `db:"content"`
	CreatedAt time.Time      `db:"created_at"`
}

func (repo *Repository) GetCommentByID(ctx echo.Context, id uuid.UUID) (Comment, error) {
	var comment Comment
	err := repo.db.GetContext(ctx.Request().Context(), &comment, "SELECT c.*, u.username FROM comments c JOIN users u ON c.author_id = u.id WHERE c.id = ?", id)
	return comment, err
}

func (repo *Repository) GetComments(ctx echo.Context, limitNumber *int) ([]Comment, error) {
	var comments []Comment
	if limitNumber != nil {
		err := repo.db.SelectContext(ctx.Request().Context(), &comments, "SELECT c.*, u.username FROM comments c JOIN users u ON c.author_id = u.id LIMIT ?", limitNumber)
		return comments, err
	} else {
		err := repo.db.SelectContext(ctx.Request().Context(), &comments, "SELECT c.*, u.username FROM comments c JOIN users u ON c.author_id = u.id")
		return comments, err
	}
}

func (repo *Repository) CreateComment(ctx echo.Context, comment Comment) error {
	_, err := repo.db.NamedExecContext(ctx.Request().Context(), "INSERT INTO comments (id, article_id, author_id, content, created_at) VALUES (:id, :article_id, :author_id, :content, :created_at)", comment)
	if err != nil {
		return err
	}
	_, err = repo.db.ExecContext(ctx.Request().Context(), "UPDATE users SET username = ? WHERE id = ?", comment.Author, comment.AuthorID)
	return err
}

func (repo *Repository) UpdateComment(ctx echo.Context, comment Comment) error {
	_, err := repo.db.NamedExecContext(ctx.Request().Context(), "UPDATE comments SET article_id = :article_id, author_id = :author_id, content = :content, created_at = :created_at WHERE id = :id", comment)
	return err
}

func (repo *Repository) DeleteComment(ctx echo.Context, id uuid.UUID) error {
	_, err := repo.db.ExecContext(ctx.Request().Context(), "DELETE FROM comments WHERE id = ?", id)
	return err
}

func (repo *Repository) GetCommentsByArticle(ctx echo.Context, article_id uuid.UUID, limitNumber *int) ([]Comment, error) {
	var comments []Comment
	if limitNumber != nil {
		err := repo.db.SelectContext(ctx.Request().Context(), &comments, "SELECT c.*, u.username FROM comments c JOIN users u ON c.author_id = u.id WHERE article_id = ? LIMIT ?", article_id, limitNumber)
		return comments, err
	} else {
		err := repo.db.SelectContext(ctx.Request().Context(), &comments, "SELECT c.*, u.username FROM comments c JOIN users u ON c.author_id = u.id WHERE article_id = ?", article_id)
		return comments, err
	}
}

func (repo *Repository) GetCommentsByAuthor(ctx echo.Context, author_id uuid.UUID, limitNumber *int) ([]Comment, error) {
	var comments []Comment
	if limitNumber != nil {
		err := repo.db.SelectContext(ctx.Request().Context(), &comments, "SELECT c.*, u.username FROM comments c JOIN users u ON c.author_id = u.id WHERE author_id = ? LIMIT ?", author_id, limitNumber)
		return comments, err
	} else {
		err := repo.db.SelectContext(ctx.Request().Context(), &comments, "SELECT c.*, u.username FROM comments c JOIN users u ON c.author_id = u.id WHERE author_id = ?", author_id)
		return comments, err
	}
}
