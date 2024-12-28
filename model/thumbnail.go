package model

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fogleman/gg"
	"github.com/golang/freetype"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func generateThumbnail(article Article, tags []TagItem, category string, authorName string) (*image.RGBA, error) {
	// サムネイルのサイズ
	width := 300
	height := 300

	log.Println("authorName", authorName)
	log.Println("category", category)
	log.Println("tags", tagsToStringSlice(tags))

	// 背景色の設定
	bgColor := color.RGBA{173, 216, 230, 255} // LightBlue

	// テキストの色
	titleColor := color.RGBA{0, 0, 139, 255}    // DarkBlue
	infoColor := color.RGBA{105, 105, 105, 255} // DimGray

	// フォントファイルの読み込み
	fontBytes, err := ioutil.ReadFile("fonts/Roboto-Regular.ttf")
	if err != nil {
		log.Printf("フォントファイルの読み込みに失敗しました: %v", err)
		return nil, err
	}

	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Printf("フォントの解析に失敗しました: %v", err)
		return nil, err
	}

	// 新しいRGBA画像を作成
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// 背景を塗りつぶす
	draw.Draw(img, img.Bounds(), &image.Uniform{bgColor}, image.Point{}, draw.Src)

	// freetypeコンテキストの設定
	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(f)
	c.SetClip(img.Bounds())
	c.SetDst(img)

	// タイトル描画の設定
	c.SetFontSize(42)
	c.SetSrc(image.NewUniform(titleColor))
	titlePt := freetype.Pt(40, 50) // タイトルの描画位置
	_, err = c.DrawString(article.Title, titlePt)
	if err != nil {
		log.Printf("タイトルの描画に失敗しました: %v", err)
		return nil, err
	}

	// 情報描画の設定 (著者名、タグ、カテゴリ)
	c.SetFontSize(20)
	c.SetSrc(image.NewUniform(infoColor))

	infoPt := freetype.Pt(40, 120)
	_, err = c.DrawString("Author: "+authorName, infoPt)
	if err != nil {
		log.Printf("著者名の描画に失敗しました: %v", err)
		return nil, err
	}

	tagsPt := freetype.Pt(40, 200)
	_, err = c.DrawString("#"+strings.Join(tagsToStringSlice(tags), ", #"), tagsPt)
	if err != nil {
		log.Printf("タグの描画に失敗しました: %v", err)
		return nil, err
	}

	categoryPt := freetype.Pt(40, 250)
	_, err = c.DrawString("Category: "+category, categoryPt)
	if err != nil {
		log.Printf("カテゴリの描画に失敗しました: %v", err)
		return nil, err
	}

	return img, nil
}

func calculateFontSize(text string, maxWidth float64, maxHeight float64, maxFontSize, minFontSize float64, dc *gg.Context, fontFile string) float64 {
	fontSize := maxFontSize
	for fontSize >= minFontSize {
		dc.LoadFontFace(fontFile, fontSize)
		height := float64(len(WordWrapJapanese(text, maxWidth, dc))) * fontSize * 1.2
		if height <= maxHeight {
			return fontSize
		}
		fontSize -= 1.0
	}
	return minFontSize
}

func generateThumbnailNew(article Article, tags []TagItem, category string, authorName string) (*image.RGBA, error) {
	// サムネイルのサイズ
	width := 300
	height := 300

	// 背景色とテキストの色
	bgColor := color.White
	textColor := color.Black

	// 新しい画像作成
	dc := gg.NewContext(width, height)

	// 背景色を塗りつぶす
	dc.SetColor(bgColor)
	dc.Clear()

	// 上下左右のマージン
	marginTop := 40.0
	marginLeft := 20.0
	marginRight := 20.0
	marginBottom := 20.0

	// 各要素間のスペース
	spaceBetweenElements := 20.0
	spaceBetweenTagElements := 7.0

	// 使用可能な幅と高さ
	availableWidth := float64(width) - marginLeft - marginRight
	availableHeight := float64(height) - marginTop - marginBottom

	availableHeightForTitle := 110.0
	availableHeightForAuthor := availableHeight / 5
	availableHeightForTags := availableHeight / 2
	availableHeightForDate := 0.0
	availableHeightForCategory := availableHeight - availableHeightForTitle - availableHeightForAuthor - availableHeightForTags - availableHeightForDate - spaceBetweenElements*3 - spaceBetweenTagElements*float64(len(tags)-1)

	fmt.Printf("availableHeightForTitle: %f\n", availableHeightForTitle)
	fmt.Printf("availableHeightForAuthor: %f\n", availableHeightForAuthor)
	fmt.Printf("availableHeightForTags: %f\n", availableHeightForTags)
	fmt.Printf("availableHeightForCategory: %f\n", availableHeightForCategory)

	dataStartY := float64(height) - marginBottom - availableHeightForDate

	tempY := marginTop

	// タイトルの描画
	fontFile := "fonts/GenShinGothic-Bold.ttf"
	titleFontSize := calculateFontSize(article.Title, availableWidth, min(availableHeight, availableHeightForTitle), 30, 0, dc, fontFile)
	dc.LoadFontFace(fontFile, titleFontSize)
	titleWrapped := WordWrapJapanese(article.Title, availableWidth, dc)
	titleHeight := float64(len(titleWrapped)) * titleFontSize * 1.2

	if titleHeight < availableHeightForTitle {
		tempY += (availableHeightForTitle - titleHeight) / 2
	}

	log.Println("titleHeight: ", titleHeight)
	log.Println("maxHeightForTitle: ", availableHeightForTitle)

	for _, line := range titleWrapped {
		dc.SetColor(textColor)
		dc.DrawString(line, marginLeft, tempY)
		tempY += titleFontSize * 1.2
	}
	tempY += -titleFontSize*1.2 + spaceBetweenElements

	// 著者名の描画
	fontFile = "fonts/GenShinGothic-Regular.ttf"
	authorFontSize := calculateFontSize("Author: "+authorName, availableWidth, min(availableHeight+marginTop-tempY, availableHeightForAuthor), 17, 0, dc, fontFile)
	dc.LoadFontFace(fontFile, authorFontSize)
	authorWrapped := WordWrapJapanese("Author: "+authorName, availableWidth, dc)
	authorHeight := float64(len(authorWrapped)) * authorFontSize * 1.2

	fmt.Printf("tempY: %f\n", tempY)
	fmt.Printf("plan: %f\n", marginTop+availableHeightForTitle+spaceBetweenElements+spaceBetweenElements)

	if authorHeight < availableHeightForAuthor {
		tempY += (availableHeightForAuthor - authorHeight) / 2
	} else {
		log.Fatalf("authorHeight is too large")
	}

	for _, line := range authorWrapped {
		dc.SetColor(textColor)
		dc.DrawString(line, marginLeft, tempY)
		tempY += authorFontSize * 1.2
	}
	tempY += -authorFontSize*1.2 + spaceBetweenElements*1.5

	// カテゴリの描画
	fontFile = "fonts/GenShinGothic-Regular.ttf"
	categoryFontSize := calculateFontSize("Category: "+category, availableWidth, min(availableHeight+marginTop-tempY, 100), min(17, authorFontSize-2), 0, dc, fontFile)
	dc.LoadFontFace(fontFile, categoryFontSize)
	categoryWrapped := WordWrapJapanese("Category: "+category, availableWidth, dc)
	categoryHeight := float64(len(categoryWrapped)) * categoryFontSize * 1.2
	if categoryHeight < availableHeightForCategory {
		tempY += (availableHeightForCategory - categoryHeight) / 2
	}
	for _, line := range categoryWrapped {
		dc.SetColor(textColor)
		dc.DrawString(line, marginLeft, tempY)
		tempY += categoryFontSize * 1.2
	}

	tempY += -categoryFontSize*1.2 + spaceBetweenElements

	// タグの描画
	fontFile = "fonts/GenShinGothic-Regular.ttf"
	tagFontSize := 14.0
	log.Println(availableHeight+marginTop-tempY, availableHeightForTags, dataStartY, tempY)
	padding := 10.0
	maxHeightForTags := min(availableHeight+marginTop-tempY, availableHeightForTags, dataStartY-tempY)
	maxHeightForOneTag := (maxHeightForTags - spaceBetweenElements - (spaceBetweenTagElements+padding)*float64(len(tags)-1)) / float64(len(tags))
	tagFontSize = calculateFontSize("あああ", availableWidth-20, maxHeightForOneTag, 14, 0, dc, fontFile)
	for _, tag := range tags {
		tag.Name = "#" + tag.Name
		dc.LoadFontFace(fontFile, tagFontSize)
		// wrappedTag := WordWrapJapanese(tag.Name, availableWidth-20, dc)
		wrappedTag := make([]string, 1)
		wrappedTag[0] = tag.Name
		maxLineWidth := 0.0
		for _, line := range wrappedTag {
			lineWidth, _ := dc.MeasureString(line)
			if lineWidth > maxLineWidth {
				maxLineWidth = lineWidth
			}
		}
		padding := 10.0
		rectX := marginLeft
		rectY := tempY
		rectWidth := maxLineWidth + 2*padding
		rectHeight := float64(len(wrappedTag))*tagFontSize + padding

		// タグの矩形を描画
		dc.SetColor(color.RGBA{200, 200, 200, 255})
		dc.DrawRoundedRectangle(rectX, rectY, rectWidth, rectHeight, 5)
		dc.Fill()

		// タグテキストの描画
		dc.SetColor(textColor)
		lineY := rectY + tagFontSize
		for _, line := range wrappedTag {
			dc.DrawString(line, rectX+padding, lineY)
			lineY += tagFontSize * 1.2
		}

		tempY += rectHeight + spaceBetweenTagElements
	}

	// tempY += -spaceBetweenTagElements + spaceBetweenElements
	tempY = dataStartY

	// 日付を描画
	fontFile = "fonts/GenShinGothic-Regular.ttf"
	dateFontSize := 14.0
	dateFontSize = calculateFontSize(article.CreatedAt.Format("2006-01-02"), availableWidth, min(availableHeight+marginTop-tempY, availableHeightForDate), 14, 0, dc, fontFile)
	dc.LoadFontFace(fontFile, dateFontSize)
	dateWrapped := WordWrapJapanese(article.CreatedAt.Format("2006-01-02"), availableWidth, dc)
	// dateHeight := float64(len(dateWrapped)) * dateFontSize * 1.2
	// if dateHeight < availableHeight+marginTop-tempY {
	// 	tempY += (availableHeight + marginTop - tempY - dateHeight) / 2
	// }
	for _, line := range dateWrapped {
		dc.SetColor(textColor)
		dc.DrawString(line, marginLeft, tempY)
		tempY += dateFontSize * 1.2
	}

	// 画像を生成
	img := dc.Image().(*image.RGBA)
	return img, nil
}

func saveImage(img *image.RGBA, path string) error {
	// 画像をファイルに保存
	out, err := os.Create(path)
	if err != nil {
		log.Printf("failed to create thumbnail file: %v", err)
		return err
	}
	defer out.Close()

	err = png.Encode(out, img)
	if err != nil {
		log.Printf("failed to encode thumbnail: %v", err)
		return err
	}

	log.Printf("thumbnail saved: %s", path)
	return nil
}

func WordWrapJapanese(text string, maxWidth float64, dc *gg.Context) []string {
	var lines []string
	var currentLine string

	for _, char := range text {
		testLine := currentLine + string(char)
		width, _ := dc.MeasureString(testLine)
		if width > maxWidth {
			if currentLine == "" {
				// 単一の文字がmaxWidthを超える場合、そのまま追加 // 一つの文字でmaxWidthを超えるもののこと
				lines = append(lines, testLine)
				currentLine = ""
			} else {
				lines = append(lines, currentLine)
				currentLine = string(char)
			}
		} else {
			currentLine = testLine
		}
	}

	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return lines
}

func (c *Configuration) HandleThumbnailGeneration(ctx echo.Context, article Article, tags []TagItem, category string, authorName string) (string, string, string, error) {
	// 画像を保存するパスを生成
	ext := ".png"
	thumbnailFileName := uuid.New().String() + "_thumb" + ext
	thumbnailPath := filepath.Join(c.ImageUploadPath, thumbnailFileName)

	// 画像保存先ディレクトリが存在しない場合は作成
	if _, err := os.Stat(c.ImageUploadPath); os.IsNotExist(err) {
		err = os.MkdirAll(c.ImageUploadPath, os.ModePerm)
		if err != nil {
			log.Println("Failed to create image upload directory:", err)
			return "", "", "", err
		}
	}

	thumbnail, err := generateThumbnailNew(article, tags, category, authorName)
	if err != nil {
		log.Println("Failed to generate thumbnail:", err)
		return "", "", "", err
	}

	err = saveImage(thumbnail, thumbnailPath)
	if err != nil {
		log.Println("Failed to save thumbnail:", err)
		return "", "", "", err
	}

	thumbnailURL := fmt.Sprintf("%s/uploads/images/%s", strings.TrimRight(c.BaseURL, "/"), thumbnailFileName)

	return thumbnailURL, thumbnailPath, thumbnailFileName, nil
}

func tagsToStringSlice(tags []TagItem) []string {
	var tagStrings []string
	for _, tag := range tags {
		tagStrings = append(tagStrings, tag.Name)
	}
	return tagStrings
}
