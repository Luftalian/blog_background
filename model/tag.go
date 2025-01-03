package model

import (
	"blog-backend/logger"
	"context"
	"errors"

	"github.com/google/uuid"
)

type Tag struct {
	ID        uuid.UUID `db:"id"`         // tagsテーブルとarticle_tagsテーブルに存在
	ArticleID uuid.UUID `db:"article_id"` // article_tagsテーブルに存在
	Name      string    `db:"name"`       // tagsテーブルに存在
}

type TagItem struct {
	ID   uuid.UUID `db:"id"`
	Name string    `db:"name"`
}

type TagPair struct {
	TagID     uuid.UUID `db:"tag_id"`
	ArticleID uuid.UUID `db:"article_id"`
}

func (r *Repository) GetTagList(ctx context.Context) ([]TagItem, error) {
	var tags []TagItem
	err := r.db.SelectContext(ctx, &tags, "SELECT id, name FROM tags")
	return tags, err
}

func (r *Repository) AddTagItem(ctx context.Context, tag TagItem) (TagItem, error) {
	// すでにtagが存在するか確認し、存在しない場合のみ新規作成
	var t TagItem
	err := r.db.GetContext(ctx, &t, "SELECT id, name FROM tags WHERE name = ?", tag.Name)
	if err != nil && err.Error() != "sql: no rows in result set" {
		return t, err
	}
	if err == nil {
		return t, nil
	}
	_, err = r.db.NamedExecContext(ctx, "INSERT INTO tags (id, name) VALUES (:id, :name)", tag)
	return tag, err
}

func (r *Repository) AddTagItemNames(ctx context.Context, tags []TagItem) error {
	// 存在しないtagのリストを作り、それを新規作成
	var newTags []TagItem
	tagsFromDB, err := r.GetTagList(ctx)
	if err != nil {
		return err
	}
	for _, tag := range tags {
		var exists bool
		for _, t := range tagsFromDB {
			if tag.Name == t.Name {
				exists = true
				break
			}
		}
		if !exists {
			tag.ID = uuid.New()
			newTags = append(newTags, tag)
		}
	}
	if len(newTags) == 0 {
		return nil
	}
	_, err = r.db.NamedExecContext(ctx, "INSERT INTO tags (id, name) VALUES (:id, :name)", newTags)
	return err
}

func (r *Repository) AddTagPair(ctx context.Context, tagPair TagPair) error {
	_, err := r.db.NamedExecContext(ctx, "INSERT INTO article_tags (article_id, tag_id) VALUES (:article_id, :tag_id)", tagPair)
	return err
}

func (r *Repository) AddTagPairs(ctx context.Context, tagPairs []TagPair) error {
	_, err := r.db.NamedExecContext(ctx, "INSERT INTO article_tags (article_id, tag_id) VALUES (:article_id, :tag_id)", tagPairs)
	return err
}

func (r *Repository) AddTagPairsByArticle(ctx context.Context, articleID uuid.UUID, tagItems []TagItem) error {
	var tagPairs []TagPair
	for _, tagItem := range tagItems {
		tagPairs = append(tagPairs, TagPair{TagID: tagItem.ID, ArticleID: articleID})
	}
	return r.AddTagPairs(ctx, tagPairs)
}

func (r *Repository) GetTags(ctx context.Context, limitNumber *int) ([]Tag, error) {
	var tags []Tag
	if limitNumber != nil {
		err := r.db.SelectContext(ctx, &tags, `SELECT t.id, at.article_id, t.name FROM tags t JOIN article_tags at ON t.id = at.tag_id LIMIT ?`, limitNumber)
		return tags, err
	} else {
		err := r.db.SelectContext(ctx, &tags, "SELECT t.id, at.article_id, t.name FROM tags t JOIN article_tags at ON t.id = at.tag_id")
		return tags, err
	}
}

func (r *Repository) AddTag(ctx context.Context, tag Tag) error {
	// すでにtagが存在するか確認し、存在しない場合のみ新規作成
	var t Tag
	err := r.db.GetContext(ctx, &t, "SELECT t.id, at.article_id, t.name FROM tags t JOIN article_tags at ON t.id = at.tag_id WHERE t.name = ?", tag.Name)
	if err != nil && err.Error() != "sql: no rows in result set" {
		return err
	}
	if err == nil {
		return nil
	}
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	_, err = tx.NamedExecContext(ctx, "INSERT INTO tags (id, article_id, user_id, name) VALUES (:id, :article_id, :user_id, :name)", tag)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.NamedExecContext(ctx, "INSERT INTO article_tags (article_id, tag_id) VALUES (:article_id, :id)", tag)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	return err
}

func (r *Repository) AddTags(ctx context.Context, tags []Tag) error {
	return errors.New("Not implemented")

	// 存在しないtagのリストを作り、それを新規作成
	var newTags []Tag
	tagsFromDB, err := r.GetTags(ctx, nil)
	if err != nil {
		logger.Println("GetTags error:", err)
		return err
	}
	for _, tag := range tags {
		var exists bool
		for _, t := range tagsFromDB {
			if tag.Name == t.Name {
				exists = true
				break
			}
		}
		if !exists {
			newTags = append(newTags, tag)
		}
	}
	if len(newTags) == 0 {
		return nil
	}
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		logger.Println("BeginTxx error:", err)
		return err
	}
	_, err = tx.NamedExecContext(ctx, "INSERT INTO tags (id, name) VALUES (:id, :name)", newTags)
	if err != nil {
		tx.Rollback()
		logger.Println("NamedExecContext error:", err)
		return err
	}
	_, err = tx.NamedExecContext(ctx, "INSERT INTO article_tags (article_id, tag_id) VALUES (:article_id, :id)", newTags)
	if err != nil {
		tx.Rollback()
		logger.Println("NamedExecContext error2:", err)
		return err
	}
	err = tx.Commit()
	logger.Println("Commit error:", err)
	return err
}

func (r *Repository) UpdateTag(ctx context.Context, tag Tag) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	_, err = tx.NamedExecContext(ctx, "UPDATE tags SET article_id = :article_id, user_id = :user_id, name = :name WHERE id = :id", tag)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.NamedExecContext(ctx, "UPDATE article_tags SET article_id = :article_id WHERE tag_id = :id", tag)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	return err
}

func (r *Repository) DeleteTag(ctx context.Context, id uuid.UUID) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, "DELETE FROM tags WHERE id = ?", id)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.ExecContext(ctx, "DELETE FROM article_tags WHERE tag_id = ?", id)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	return err
}

func (r *Repository) GetTagNameByID(ctx context.Context, id uuid.UUID) (Tag, error) {
	return Tag{}, errors.New("Not implemented")
	// 記事に使われたことがないタグがある場合、そのタグは取得できない
	var tag Tag
	err := r.db.GetContext(ctx, &tag, "SELECT t.id, at.article_id, t.name FROM tags t JOIN article_tags at ON t.id = at.tag_id WHERE t.id = ?", id)
	return tag, err
}

func (r *Repository) GetTagIDByName(ctx context.Context, name string) (Tag, error) {
	return Tag{}, errors.New("Not implemented")
	// 記事に使われたことがないタグがある場合、そのタグは取得できない
	var tag Tag
	err := r.db.GetContext(ctx, &tag, "SELECT t.id, at.article_id, t.name FROM tags t JOIN article_tags at ON t.id = at.tag_id WHERE t.name = ?", name)
	return tag, err
}

func (r *Repository) GetTagsByArticle(ctx context.Context, article_id uuid.UUID, limitNumber *int) ([]Tag, error) {
	var tags []Tag
	if limitNumber != nil {
		err := r.db.SelectContext(ctx, &tags, "SELECT t.id, at.article_id, t.name FROM tags t JOIN article_tags at ON t.id = at.tag_id WHERE at.article_id = ? LIMIT ?", article_id, limitNumber)
		return tags, err
	} else {
		err := r.db.SelectContext(ctx, &tags, "SELECT t.id, at.article_id, t.name FROM tags t JOIN article_tags at ON t.id = at.tag_id WHERE at.article_id = ?", article_id)
		return tags, err
	}
}

func (r *Repository) GetTagsByUser(ctx context.Context, user_id uuid.UUID, limitNumber *int) ([]Tag, error) {
	var tags []Tag
	if limitNumber != nil {
		err := r.db.SelectContext(ctx, &tags, "SELECT t.id, at.article_id, t.name FROM tags t JOIN article_tags at ON t.id = at.tag_id WHERE at.user_id = ? LIMIT ?", user_id, limitNumber)
		return tags, err
	} else {
		err := r.db.SelectContext(ctx, &tags, "SELECT t.id, at.article_id, t.name FROM tags t JOIN article_tags at ON t.id = at.tag_id WHERE at.user_id = ?", user_id)
		return tags, err
	}
}

func (r *Repository) GetTagItemsByID(ctx context.Context, id uuid.UUID) (TagItem, error) {
	var tag TagItem
	err := r.db.GetContext(ctx, &tag, "SELECT id, name FROM tags WHERE id = ?", id)
	return tag, err
}
