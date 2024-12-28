package handler

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"

	"blog-backend/api"
	"blog-backend/model"
)

// Upload an image
// (POST /images/upload)
func (h *Handler) UploadImage(ctx echo.Context) error {
	// リクエストからファイルを取得
	file, err := ctx.FormFile("image")
	if err != nil {
		log.Println("Image file is required")
		return ctx.JSON(http.StatusBadRequest, api.ErrorResponse{
			Message: "Image file is required",
			Code:    http.StatusBadRequest,
		})
	}

	// ファイルサイズの検証
	if file.Size > h.Config.MaxFileSize {
		log.Println("File size exceeds the maximum limit of ", h.Config.MaxFileSize, " bytes")
		return ctx.JSON(http.StatusBadRequest, api.ErrorResponse{
			Message: fmt.Sprintf("File size exceeds the maximum limit of %d bytes", h.Config.MaxFileSize),
			Code:    http.StatusBadRequest,
		})
	}

	// ファイルタイプの検証
	src, err := file.Open()
	if err != nil {
		log.Println("Failed to open uploaded file: ", err)
		return ctx.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Message: "Failed to open uploaded file",
			Code:    http.StatusInternalServerError,
		})
	}
	defer src.Close()

	// ファイルヘッダーからMIMEタイプを取得
	buffer := make([]byte, 512)
	_, err = src.Read(buffer)
	if err != nil {
		log.Println("Failed to read file: ", err)
		return ctx.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Message: "Failed to read file",
			Code:    http.StatusInternalServerError,
		})
	}
	contentType := http.DetectContentType(buffer)

	// 許可する画像フォーマット
	allowedTypes := []string{"image/jpeg", "image/png", "image/gif"}

	if !model.IsAllowedContentType(contentType, allowedTypes) {
		log.Println("Unsupported image format. Only JPEG, PNG, and GIF are allowed.")
		return ctx.JSON(http.StatusBadRequest, api.ErrorResponse{
			Message: "Unsupported image format. Only JPEG, PNG, and GIF are allowed.",
			Code:    http.StatusBadRequest,
		})
	}

	// 元のファイル名から拡張子を取得
	ext := filepath.Ext(file.Filename)
	if ext == "" {
		// 拡張子がない場合、MIMEタイプから推測
		ext = model.MimeExtension(contentType)
		if ext == "" {
			log.Println("Cannot determine file extension")
			return ctx.JSON(http.StatusBadRequest, api.ErrorResponse{
				Message: "Cannot determine file extension",
				Code:    http.StatusBadRequest,
			})
		}
	}

	// 一意なファイル名を生成
	newFileName := model.GenerateUniqueFileName(ext)

	// 画像保存先ディレクトリが存在しない場合は作成
	if _, err := os.Stat(h.Config.ImageUploadPath); os.IsNotExist(err) {
		err = os.MkdirAll(h.Config.ImageUploadPath, os.ModePerm)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, api.ErrorResponse{
				Message: "Failed to create image upload directory",
				Code:    http.StatusInternalServerError,
			})
		}
	}

	// 画像を保存するパスを生成
	dstPath := filepath.Join(h.Config.ImageUploadPath, newFileName)

	// ファイルを保存（手動で保存する方法）
	if err := model.SaveImageToLocal(src, dstPath); err != nil {
		log.Println("Failed to create file: ", err)
		return ctx.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Message: "アップロードされた画像を保存できませんでした",
			Code:    http.StatusInternalServerError,
		})
	}

	// 画像のURLを生成
	imageURL := fmt.Sprintf("%s/uploads/images/%s", strings.TrimRight(h.Config.BaseURL, "/"), newFileName)

	// レスポンスを返す
	err = ctx.JSON(http.StatusOK, api.ImageUploadResponse{
		Url: imageURL,
	})
	if err != nil {
		log.Println("JSONレスポンスの送信に失敗しました: ", err)
	}

	// 非同期でGoogle Driveへアップロード
	if h.DriveService != nil {
		model.UploadAsyncToDrive(h.DriveService, dstPath, newFileName, os.Getenv("DRIVE_FOLDER_ID"))
	} else {
		log.Println("Drive service is not set")
	}
	return nil
}
