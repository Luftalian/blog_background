package model

import (
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

// Configuration - 画像保存先やベースURLなどの設定を持つ構造体
type Configuration struct {
	ImageUploadPath string // 例: "./uploads/images"
	BaseURL         string // 例: "http://localhost:8080"
	MaxFileSize     int64  // 例: 5 * 1024 * 1024 (5MB)
}

func NewUploader(ImageUploadPath string, BaseURL string, MaxFileSize int64) *Configuration {
	return &Configuration{
		ImageUploadPath: ImageUploadPath,
		BaseURL:         BaseURL,
		MaxFileSize:     MaxFileSize,
	}
}
