package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Like struct {
	ID        uuid.UUID `db:"id"`
	ArticleID uuid.UUID `db:"article_id"`
	UserID    uuid.UUID `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
}

func (repo *Repository) GetLikeByID(ctx echo.Context, id uuid.UUID) (Like, error) {
	var like Like
	err := repo.db.GetContext(ctx.Request().Context(), &like, "SELECT * FROM likes WHERE id = ?", id)
	return like, err
}

func (repo *Repository) GetLikes(ctx echo.Context, limitNumber *int) ([]Like, error) {
	var likes []Like
	if limitNumber != nil {
		err := repo.db.SelectContext(ctx.Request().Context(), &likes, "SELECT * FROM likes LIMIT ?", limitNumber)
		return likes, err
	} else {
		err := repo.db.SelectContext(ctx.Request().Context(), &likes, "SELECT * FROM likes")
		return likes, err
	}
}

func (repo *Repository) CreateLike(ctx echo.Context, like Like) error {
	_, err := repo.db.NamedExecContext(ctx.Request().Context(), "INSERT INTO likes (id, article_id, user_id, created_at) VALUES (:id, :article_id, :user_id, :created_at)", like)
	return err
}

func (repo *Repository) UpdateLike(ctx echo.Context, like Like) error {
	_, err := repo.db.NamedExecContext(ctx.Request().Context(), "UPDATE likes SET article_id = :article_id, user_id = :user_id, created_at = :created_at WHERE id = :id", like)
	return err
}

func (repo *Repository) DeleteLike(ctx echo.Context, id uuid.UUID) error {
	_, err := repo.db.ExecContext(ctx.Request().Context(), "DELETE FROM likes WHERE id = ?", id)
	return err
}

func (repo *Repository) GetLikesByArticle(ctx echo.Context, article_id uuid.UUID, limitNumber *int) ([]Like, error) {
	var likes []Like
	if limitNumber != nil {
		err := repo.db.SelectContext(ctx.Request().Context(), &likes, "SELECT * FROM likes WHERE article_id = ? LIMIT ?", article_id, limitNumber)
		return likes, err
	} else {
		err := repo.db.SelectContext(ctx.Request().Context(), &likes, "SELECT * FROM likes WHERE article_id = ?", article_id)
		return likes, err
	}
}

func (repo *Repository) GetLikesByUser(ctx echo.Context, user_id uuid.UUID, limitNumber *int) ([]Like, error) {
	var likes []Like
	if limitNumber != nil {
		err := repo.db.SelectContext(ctx.Request().Context(), &likes, "SELECT * FROM likes WHERE user_id = ? LIMIT ?", user_id, limitNumber)
		return likes, err
	} else {
		err := repo.db.SelectContext(ctx.Request().Context(), &likes, "SELECT * FROM likes WHERE user_id = ?", user_id)
		return likes, err
	}
}

func (repo *Repository) GetLikesCountByArticle(ctx echo.Context, article_id uuid.UUID) (int, error) {
	var count int
	err := repo.db.GetContext(ctx.Request().Context(), &count, "SELECT COUNT(*) FROM likes WHERE article_id = ?", article_id)
	return count, err
}
