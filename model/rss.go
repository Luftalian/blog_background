package model

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/feeds"
	"github.com/labstack/echo/v4"
)

func (c *Configuration) RSSmaker(ctx echo.Context, articles []Article) error {
	// RSSフィードの基本情報を設定
	feed := &feeds.Feed{
		Title:       "My Blog",
		Link:        &feeds.Link{Href: getEnv("PAGE_LINK", "http://localhost:5173")},
		Description: "This is my personal blog",
		Author:      &feeds.Author{Name: getEnv("AUTHOR_NAME", ""), Email: getEnv("AUTHOR_EMAIL", "")},
		Created:     time.Now(),
	}

	// フィードに記事を追加
	for _, article := range articles {
		item := &feeds.Item{
			Title:       article.Title,
			Link:        &feeds.Link{Href: getEnv("PAGE_LINK", "http://localhost:5173") + "/article/" + article.ID.String()},
			Description: article.Content,                                                                   // 必要に応じて要約を使用
			Author:      &feeds.Author{Name: getEnv("AUTHOR_NAME", ""), Email: getEnv("AUTHOR_EMAIL", "")}, // 著者情報を適宜設定
			Created:     article.CreatedAt,
			Id:          article.ID.String(), // GUIDに使用
		}
		feed.Items = append(feed.Items, item)
	}

	// RSS XMLを生成
	rss, err := feed.ToRss()
	if err != nil {
		// エラーログを出力
		log.Printf("Failed to generate RSS feed: %v", err)
		return err
	}

	// rss.xmlファイルへのパスを指定
	rssFilePath := filepath.Join("/app/rss", "rss.xml")

	// ディレクトリが存在しない場合は作成
	err = os.MkdirAll(filepath.Dir(rssFilePath), os.ModePerm)
	if err != nil {
		log.Printf("Failed to create directory for rss.xml: %v", err)
		return err
	}

	// rss.xmlファイルに書き込み
	err = os.WriteFile(rssFilePath, []byte(rss), 0644)
	if err != nil {
		log.Printf("Failed to write rss.xml: %v", err)
		return err
	}

	return nil
}
