package main

import (
	"log"
	"os"
	"strconv"
	"strings"

	"blog-backend/api"
	"blog-backend/handler"
	"blog-backend/migration"
	"blog-backend/model"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	allowOrigins := strings.Split(os.Getenv("ALLOW_ORIGINS"), ",")

	// middlewares
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: allowOrigins,
	}))

	e.Static("/uploads/images", "uploads/images")
	// RSSフィードを静的ファイルとして配信
	e.Static("/rss", "rss")

	dev, err := strconv.ParseBool(os.Getenv("DEVELOPMENT"))
	if err != nil {
		dev = false
	}
	log.Println("development mode:", dev)

	// connect to database
	db, err := sqlx.Connect("mysql", model.MySQL().FormatDSN())
	if err != nil {
		e.Logger.Fatal(err)
	}
	defer db.Close()

	// migrate tables
	if err := migration.MigrateTables(db.DB); err != nil {
		e.Logger.Fatal(err)
	}

	// setup repository
	repo := model.New(db)

	// setup configuration (5MBの上限など)
	config := model.NewUploader("/app/uploads/images", os.Getenv("BASE_URL"), 5*1024*1024)

	// Google Driveサービスセットアップ & ファイルダウンロード
	driveService := model.SetupGoogleDrive()

	// ハンドラーにGoogle Driveサービスを渡す
	h := handler.New(repo, config, driveService)

	// ルーティング
	api.RegisterHandlersWithBaseURL(e, h, "/api/v1")

	e.Logger.Fatal(e.Start(":8080"))
}
