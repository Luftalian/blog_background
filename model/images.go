package model

import (
	"blog-backend/logger"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"google.golang.org/api/drive/v3"
)

// isAllowedContentType は指定されたMIMEタイプが許可されているかを確認します
func IsAllowedContentType(contentType string, allowedTypes []string) bool {
	for _, t := range allowedTypes {
		if contentType == t {
			return true
		}
	}
	return false
}

// mimeExtension はMIMEタイプからファイル拡張子を返します
func MimeExtension(contentType string) string {
	switch contentType {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/gif":
		return ".gif"
	default:
		return ""
	}
}

func GenerateUniqueFileName(ext string) string {
	return uuid.New().String() + ext
}

func SaveImageToLocal(src io.ReadSeeker, dstPath string) error {
	// ディレクトリがなければ作成
	dir := filepath.Dir(dstPath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if mkdirErr := os.MkdirAll(dir, os.ModePerm); mkdirErr != nil {
			return fmt.Errorf("failed to create image upload directory: %w", mkdirErr)
		}
	}
	// 書き込み先作成
	dst, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	// ファイルポインタを先頭に戻してからコピー
	if _, err = src.Seek(0, 0); err != nil {
		return fmt.Errorf("failed to reset file pointer: %w", err)
	}
	if _, err = io.Copy(dst, src); err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}
	return nil
}

func UploadAsyncToDrive(service *drive.Service, path, name, folderID string) {
	go func() {
		fileToUpload, e := os.Open(path)
		if e != nil {
			logger.Println("Failed to open uploaded file:", e)
			return
		}
		defer fileToUpload.Close()
		driveFile := &drive.File{
			Name:    name,
			Parents: []string{folderID},
		}
		if _, e := service.Files.Create(driveFile).Media(fileToUpload).Do(); e != nil {
			logger.Println("Failed to upload to Drive:", e)
		}
	}()
}
