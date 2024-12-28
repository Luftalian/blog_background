package model

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Like struct {
	ID        uuid.UUID `db:"id"`
	ArticleID uuid.UUID `db:"article_id"`
	UserID    uuid.UUID `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
}

func (repo *Repository) GetLikeByID(ctx context.Context, id uuid.UUID) (Like, error) {
	var like Like
	err := repo.db.GetContext(ctx, &like, "SELECT * FROM likes WHERE id = ?", id)
	return like, err
}

func (repo *Repository) GetLikes(ctx context.Context, limitNumber *int) ([]Like, error) {
	var likes []Like
	if limitNumber != nil {
		err := repo.db.SelectContext(ctx, &likes, "SELECT * FROM likes LIMIT ?", limitNumber)
		return likes, err
	} else {
		err := repo.db.SelectContext(ctx, &likes, "SELECT * FROM likes")
		return likes, err
	}
}

func (repo *Repository) CreateLike(ctx context.Context, like Like) error {
	_, err := repo.db.NamedExecContext(ctx, "INSERT INTO likes (id, article_id, user_id, created_at) VALUES (:id, :article_id, :user_id, :created_at)", like)
	return err
}

func (repo *Repository) UpdateLike(ctx context.Context, like Like) error {
	_, err := repo.db.NamedExecContext(ctx, "UPDATE likes SET article_id = :article_id, user_id = :user_id, created_at = :created_at WHERE id = :id", like)
	return err
}

func (repo *Repository) DeleteLike(ctx context.Context, id uuid.UUID) error {
	_, err := repo.db.ExecContext(ctx, "DELETE FROM likes WHERE id = ?", id)
	return err
}

func (repo *Repository) GetLikesByArticle(ctx context.Context, article_id uuid.UUID, limitNumber *int) ([]Like, error) {
	var likes []Like
	if limitNumber != nil {
		err := repo.db.SelectContext(ctx, &likes, "SELECT * FROM likes WHERE article_id = ? LIMIT ?", article_id, limitNumber)
		return likes, err
	} else {
		err := repo.db.SelectContext(ctx, &likes, "SELECT * FROM likes WHERE article_id = ?", article_id)
		return likes, err
	}
}

func (repo *Repository) GetLikesByUser(ctx context.Context, user_id uuid.UUID, limitNumber *int) ([]Like, error) {
	var likes []Like
	if limitNumber != nil {
		err := repo.db.SelectContext(ctx, &likes, "SELECT * FROM likes WHERE user_id = ? LIMIT ?", user_id, limitNumber)
		return likes, err
	} else {
		err := repo.db.SelectContext(ctx, &likes, "SELECT * FROM likes WHERE user_id = ?", user_id)
		return likes, err
	}
}

func (repo *Repository) GetLikesCountByArticle(ctx context.Context, article_id uuid.UUID) (int, error) {
	var count int
	err := repo.db.GetContext(ctx, &count, "SELECT COUNT(*) FROM likes WHERE article_id = ?", article_id)
	return count, err
}
