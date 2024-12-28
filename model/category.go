package model

import (
	"context"

	"github.com/google/uuid"
)

type Category struct {
	ID   uuid.UUID `db:"id"`
	Name string    `db:"name"`
}

func (r *Repository) GetCategories(ctx context.Context, limitNumber *int) ([]Category, error) {
	var categories []Category
	if limitNumber != nil {
		err := r.db.SelectContext(ctx, &categories, "SELECT * FROM categories LIMIT ?", limitNumber)
		return categories, err
	} else {
		err := r.db.SelectContext(ctx, &categories, "SELECT * FROM categories")
		return categories, err
	}
}

func (r *Repository) CreateCategory(ctx context.Context, category Category) error {
	_, err := r.db.NamedExecContext(ctx, "INSERT INTO categories (id, name) VALUES (:id, :name)", category)
	return err
}

func (r *Repository) UpdateCategory(ctx context.Context, category Category) error {
	_, err := r.db.NamedExecContext(ctx, "UPDATE categories SET name = :name WHERE id = :id", category)
	return err
}

func (r *Repository) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM categories WHERE id = ?", id)
	return err
}

func (r *Repository) GetCategoryNameByID(ctx context.Context, id uuid.UUID) (Category, error) {
	var category Category
	err := r.db.GetContext(ctx, &category, "SELECT * FROM categories WHERE id = ?", id)
	return category, err
}

func (r *Repository) GetCategoryIDByName(ctx context.Context, name string) (Category, error) {
	var category Category
	err := r.db.GetContext(ctx, &category, "SELECT * FROM categories WHERE name = ?", name)
	return category, err
}

func (r *Repository) AddCategory(ctx context.Context, name string) (Category, error) {
	// Check if the category already exists
	var category Category
	err := r.db.GetContext(ctx, &category, "SELECT * FROM categories WHERE name = ?", name)
	if err != nil && err.Error() != "sql: no rows in result set" {
		return category, err
	}
	if err == nil {
		return category, nil
	}
	// Create a new category
	category = Category{
		ID:   uuid.New(),
		Name: name,
	}
	err = r.CreateCategory(ctx, category)
	return category, err
}

func (r *Repository) AddCategoryID(ctx context.Context, id uuid.UUID) (Category, error) {
	// Check if the category already exists
	var category Category
	err := r.db.GetContext(ctx, &category, "SELECT * FROM categories WHERE id = ?", id)
	if err != nil {
		return category, err
	}
	return category, nil
}
