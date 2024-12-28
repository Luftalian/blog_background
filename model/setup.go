package model

import (
	"blog-backend/logger"
	"context"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
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

func SetupFirstRss(repo *Repository, config *Configuration) error {
	ctx := context.Background()
	limit := 5

	articles, err := repo.GetArticlesByCategoryTagSearch(ctx, &uuid.Nil, &uuid.Nil, nil, &limit, "created_at", "desc")
	if err != nil {
		logger.Printf("Failed to fetch articles: %v", err)
		return err
	}

	sort.Slice(articles, func(i, j int) bool {
		return articles[i].CreatedAt.After(articles[j].CreatedAt)
	})

	for i, article := range articles {
		logger.Printf("Article %d: %s", i, article.Title)
		logger.Printf("Article %d: %s", i, article.CreatedAt)
	}

	if err := config.RSSmaker(ctx, articles); err != nil {
		logger.Printf("Failed to generate RSS feed: %v", err)
		return err
	}

	return nil
}

// SetupGoogleDrive - Google Driveサービスをセットアップし、ファイルをダウンロードする
func SetupGoogleDrive() *drive.Service {
	// 環境変数からサービスアカウントのJSONを取得
	serviceAccountJSON := os.Getenv("SERVICE_ACCOUNT_JSON")
	if serviceAccountJSON == "" {
		logger.Fatal("No service account JSON found in SERVICE_ACCOUNT_JSON environment variable.")
	}

	// JWTConfig作成 (DriveFileScope = Driveへのフルアクセス)
	config, err := google.JWTConfigFromJSON([]byte(serviceAccountJSON), drive.DriveFileScope)
	if err != nil {
		logger.Fatalln("JWTConfigの作成に失敗:", err)
	}
	client := config.Client(context.Background())

	// Driveサービス生成
	srv, err := drive.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		logger.Fatalln("Driveサービスの作成に失敗:", err)
	}

	// フォルダID
	folderID := os.Getenv("DRIVE_FOLDER_ID")
	if folderID == "" {
		logger.Fatal("No folder ID specified in FOLDER_ID environment variable.")
	}

	// ファイル一覧を取得
	r, err := srv.Files.List().
		Q(fmt.Sprintf("'%s' in parents", folderID)).
		Fields("files(id, name)").Do()
	if err != nil {
		logger.Fatalf("Failed to retrieve files from folder %s: %v\n", folderID, err)
	}
	if len(r.Files) == 0 {
		logger.Printf("No files found in folderID=%s\n", folderID)
	} else {
		logger.Printf("Files in folderID=%s:\n", folderID)
	}

	// ディレクトリ作成
	if _, err := os.Stat("uploads/images"); os.IsNotExist(err) {
		if err := os.MkdirAll("uploads/images", 0755); err != nil {
			logger.Fatalf("Failed to create uploads/images directory: %v\n", err)
		}
	}

	// ファイルをローカルにダウンロード
	for _, f := range r.Files {
		path := fmt.Sprintf("uploads/images/%s", f.Name)
		if _, err := os.Stat(path); err == nil {
			logger.Printf("File %s already exists, skipping.", path)
			continue
		}

		logger.Printf("  Name=%s, ID=%s\n", f.Name, f.Id)
		resp, err := srv.Files.Get(f.Id).Download()
		if err != nil {
			logger.Fatalf("Failed to download file %s: %v\n", f.Name, err)
		}
		defer resp.Body.Close()

		outFile, err := os.Create(path)
		if err != nil {
			logger.Fatalf("Failed to create file %s: %v\n", f.Name, err)
		}
		defer outFile.Close()

		_, err = io.Copy(outFile, resp.Body)
		if err != nil {
			logger.Fatalf("Failed to save file %s: %v\n", f.Name, err)
		}
		logger.Printf("Downloaded file %s to %s\n", f.Name, path)
	}

	return srv
}
