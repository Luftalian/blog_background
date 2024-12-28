package model

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	ID        uuid.UUID      `db:"id"`
	ArticleID uuid.UUID      `db:"article_id"`
	AuthorID  uuid.UUID      `db:"author_id"`
	Author    sql.NullString `db:"username"` // このフィールドはDBには存在しない。Userテーブルから取得してくる
	Content   string         `db:"content"`
	CreatedAt time.Time      `db:"created_at"`
}

func (repo *Repository) GetCommentByID(ctx context.Context, id uuid.UUID) (Comment, error) {
	var comment Comment
	err := repo.db.GetContext(ctx, &comment, "SELECT c.*, u.username FROM comments c JOIN users u ON c.author_id = u.id WHERE c.id = ?", id)
	return comment, err
}

func (repo *Repository) GetComments(ctx context.Context, limitNumber *int) ([]Comment, error) {
	var comments []Comment
	if limitNumber != nil {
		err := repo.db.SelectContext(ctx, &comments, "SELECT c.*, u.username FROM comments c JOIN users u ON c.author_id = u.id LIMIT ?", limitNumber)
		return comments, err
	} else {
		err := repo.db.SelectContext(ctx, &comments, "SELECT c.*, u.username FROM comments c JOIN users u ON c.author_id = u.id")
		return comments, err
	}
}

func (repo *Repository) CreateComment(ctx context.Context, comment Comment) error {
	_, err := repo.db.NamedExecContext(ctx, "INSERT INTO comments (id, article_id, author_id, content, created_at) VALUES (:id, :article_id, :author_id, :content, :created_at)", comment)
	if err != nil {
		return err
	}
	_, err = repo.db.ExecContext(ctx, "UPDATE users SET username = ? WHERE id = ?", comment.Author, comment.AuthorID)
	return err
}

func (repo *Repository) UpdateComment(ctx context.Context, comment Comment) error {
	_, err := repo.db.NamedExecContext(ctx, "UPDATE comments SET article_id = :article_id, author_id = :author_id, content = :content, created_at = :created_at WHERE id = :id", comment)
	return err
}

func (repo *Repository) DeleteComment(ctx context.Context, id uuid.UUID) error {
	_, err := repo.db.ExecContext(ctx, "DELETE FROM comments WHERE id = ?", id)
	return err
}

func (repo *Repository) GetCommentsByArticle(ctx context.Context, article_id uuid.UUID, limitNumber *int) ([]Comment, error) {
	var comments []Comment
	if limitNumber != nil {
		err := repo.db.SelectContext(ctx, &comments, "SELECT c.*, u.username FROM comments c JOIN users u ON c.author_id = u.id WHERE article_id = ? LIMIT ?", article_id, limitNumber)
		return comments, err
	} else {
		err := repo.db.SelectContext(ctx, &comments, "SELECT c.*, u.username FROM comments c JOIN users u ON c.author_id = u.id WHERE article_id = ?", article_id)
		return comments, err
	}
}

func (repo *Repository) GetCommentsByAuthor(ctx context.Context, author_id uuid.UUID, limitNumber *int) ([]Comment, error) {
	var comments []Comment
	if limitNumber != nil {
		err := repo.db.SelectContext(ctx, &comments, "SELECT c.*, u.username FROM comments c JOIN users u ON c.author_id = u.id WHERE author_id = ? LIMIT ?", author_id, limitNumber)
		return comments, err
	} else {
		err := repo.db.SelectContext(ctx, &comments, "SELECT c.*, u.username FROM comments c JOIN users u ON c.author_id = u.id WHERE author_id = ?", author_id)
		return comments, err
	}
}
