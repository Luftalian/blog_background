package model

import (
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
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

func (r *Repository) GetTagList(ctx echo.Context) ([]TagItem, error) {
	var tags []TagItem
	err := r.db.SelectContext(ctx.Request().Context(), &tags, "SELECT id, name FROM tags")
	return tags, err
}

func (r *Repository) AddTagItem(ctx echo.Context, tag TagItem) (TagItem, error) {
	// すでにtagが存在するか確認し、存在しない場合のみ新規作成
	var t TagItem
	err := r.db.GetContext(ctx.Request().Context(), &t, "SELECT id, name FROM tags WHERE name = ?", tag.Name)
	if err != nil && err.Error() != "sql: no rows in result set" {
		return t, err
	}
	if err == nil {
		return t, nil
	}
	_, err = r.db.NamedExecContext(ctx.Request().Context(), "INSERT INTO tags (id, name) VALUES (:id, :name)", tag)
	return tag, err
}

func (r *Repository) AddTagItemNames(ctx echo.Context, tags []TagItem) error {
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
	_, err = r.db.NamedExecContext(ctx.Request().Context(), "INSERT INTO tags (id, name) VALUES (:id, :name)", newTags)
	return err
}

func (r *Repository) AddTagPair(ctx echo.Context, tagPair TagPair) error {
	_, err := r.db.NamedExecContext(ctx.Request().Context(), "INSERT INTO article_tags (article_id, tag_id) VALUES (:article_id, :tag_id)", tagPair)
	return err
}

func (r *Repository) AddTagPairs(ctx echo.Context, tagPairs []TagPair) error {
	_, err := r.db.NamedExecContext(ctx.Request().Context(), "INSERT INTO article_tags (article_id, tag_id) VALUES (:article_id, :tag_id)", tagPairs)
	return err
}

func (r *Repository) AddTagPairsByArticle(ctx echo.Context, articleID uuid.UUID, tagItems []TagItem) error {
	var tagPairs []TagPair
	for _, tagItem := range tagItems {
		tagPairs = append(tagPairs, TagPair{TagID: tagItem.ID, ArticleID: articleID})
	}
	return r.AddTagPairs(ctx, tagPairs)
}

func (r *Repository) GetTags(ctx echo.Context, limitNumber *int) ([]Tag, error) {
	var tags []Tag
	if limitNumber != nil {
		err := r.db.SelectContext(ctx.Request().Context(), &tags, `SELECT t.id, at.article_id, t.name FROM tags t JOIN article_tags at ON t.id = at.tag_id LIMIT ?`, limitNumber)
		return tags, err
	} else {
		err := r.db.SelectContext(ctx.Request().Context(), &tags, "SELECT t.id, at.article_id, t.name FROM tags t JOIN article_tags at ON t.id = at.tag_id")
		return tags, err
	}
}

func (r *Repository) AddTag(ctx echo.Context, tag Tag) error {
	// すでにtagが存在するか確認し、存在しない場合のみ新規作成
	var t Tag
	err := r.db.GetContext(ctx.Request().Context(), &t, "SELECT t.id, at.article_id, t.name FROM tags t JOIN article_tags at ON t.id = at.tag_id WHERE t.name = ?", tag.Name)
	if err != nil && err.Error() != "sql: no rows in result set" {
		return err
	}
	if err == nil {
		return nil
	}
	tx, err := r.db.BeginTxx(ctx.Request().Context(), nil)
	if err != nil {
		return err
	}
	_, err = tx.NamedExecContext(ctx.Request().Context(), "INSERT INTO tags (id, article_id, user_id, name) VALUES (:id, :article_id, :user_id, :name)", tag)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.NamedExecContext(ctx.Request().Context(), "INSERT INTO article_tags (article_id, tag_id) VALUES (:article_id, :id)", tag)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	return err
}

func (r *Repository) AddTags(ctx echo.Context, tags []Tag) error {
	return errors.New("Not implemented")

	// 存在しないtagのリストを作り、それを新規作成
	var newTags []Tag
	tagsFromDB, err := r.GetTags(ctx, nil)
	if err != nil {
		log.Println("GetTags error:", err)
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
	tx, err := r.db.BeginTxx(ctx.Request().Context(), nil)
	if err != nil {
		log.Println("BeginTxx error:", err)
		return err
	}
	_, err = tx.NamedExecContext(ctx.Request().Context(), "INSERT INTO tags (id, name) VALUES (:id, :name)", newTags)
	if err != nil {
		tx.Rollback()
		log.Println("NamedExecContext error:", err)
		return err
	}
	_, err = tx.NamedExecContext(ctx.Request().Context(), "INSERT INTO article_tags (article_id, tag_id) VALUES (:article_id, :id)", newTags)
	if err != nil {
		tx.Rollback()
		log.Println("NamedExecContext error2:", err)
		return err
	}
	err = tx.Commit()
	log.Println("Commit error:", err)
	return err
}

func (r *Repository) UpdateTag(ctx echo.Context, tag Tag) error {
	tx, err := r.db.BeginTxx(ctx.Request().Context(), nil)
	if err != nil {
		return err
	}
	_, err = tx.NamedExecContext(ctx.Request().Context(), "UPDATE tags SET article_id = :article_id, user_id = :user_id, name = :name WHERE id = :id", tag)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.NamedExecContext(ctx.Request().Context(), "UPDATE article_tags SET article_id = :article_id WHERE tag_id = :id", tag)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	return err
}

func (r *Repository) DeleteTag(ctx echo.Context, id uuid.UUID) error {
	tx, err := r.db.BeginTxx(ctx.Request().Context(), nil)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx.Request().Context(), "DELETE FROM tags WHERE id = ?", id)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.ExecContext(ctx.Request().Context(), "DELETE FROM article_tags WHERE tag_id = ?", id)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	return err
}

func (r *Repository) GetTagNameByID(ctx echo.Context, id uuid.UUID) (Tag, error) {
	return Tag{}, errors.New("Not implemented")
	// 記事に使われたことがないタグがある場合、そのタグは取得できない
	var tag Tag
	err := r.db.GetContext(ctx.Request().Context(), &tag, "SELECT t.id, at.article_id, t.name FROM tags t JOIN article_tags at ON t.id = at.tag_id WHERE t.id = ?", id)
	return tag, err
}

func (r *Repository) GetTagIDByName(ctx echo.Context, name string) (Tag, error) {
	return Tag{}, errors.New("Not implemented")
	// 記事に使われたことがないタグがある場合、そのタグは取得できない
	var tag Tag
	err := r.db.GetContext(ctx.Request().Context(), &tag, "SELECT t.id, at.article_id, t.name FROM tags t JOIN article_tags at ON t.id = at.tag_id WHERE t.name = ?", name)
	return tag, err
}

func (r *Repository) GetTagsByArticle(ctx echo.Context, article_id uuid.UUID, limitNumber *int) ([]Tag, error) {
	var tags []Tag
	if limitNumber != nil {
		err := r.db.SelectContext(ctx.Request().Context(), &tags, "SELECT t.id, at.article_id, t.name FROM tags t JOIN article_tags at ON t.id = at.tag_id WHERE at.article_id = ? LIMIT ?", article_id, limitNumber)
		return tags, err
	} else {
		err := r.db.SelectContext(ctx.Request().Context(), &tags, "SELECT t.id, at.article_id, t.name FROM tags t JOIN article_tags at ON t.id = at.tag_id WHERE at.article_id = ?", article_id)
		return tags, err
	}
}

func (r *Repository) GetTagsByUser(ctx echo.Context, user_id uuid.UUID, limitNumber *int) ([]Tag, error) {
	var tags []Tag
	if limitNumber != nil {
		err := r.db.SelectContext(ctx.Request().Context(), &tags, "SELECT t.id, at.article_id, t.name FROM tags t JOIN article_tags at ON t.id = at.tag_id WHERE at.user_id = ? LIMIT ?", user_id, limitNumber)
		return tags, err
	} else {
		err := r.db.SelectContext(ctx.Request().Context(), &tags, "SELECT t.id, at.article_id, t.name FROM tags t JOIN article_tags at ON t.id = at.tag_id WHERE at.user_id = ?", user_id)
		return tags, err
	}
}

func (r *Repository) GetTagItemsByID(ctx echo.Context, id uuid.UUID) (TagItem, error) {
	var tag TagItem
	err := r.db.GetContext(ctx.Request().Context(), &tag, "SELECT id, name FROM tags WHERE id = ?", id)
	return tag, err
}
