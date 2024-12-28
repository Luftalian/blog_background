package model

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// CREATE TABLE `analysis` (
//     `id` CHAR(36) NOT NULL,
//     `timestamp` TIMESTAMP,
//     `articleId` CHAR(36),
//     `ipaddress` VARCHAR(45),
//     `search_word` VARCHAR(255),
//     `api` VARCHAR(2083),
//     `is_error` BOOLEAN
// );

type Analysis struct {
	ID         uuid.UUID `db:"id"`
	Timestamp  time.Time `db:"timestamp"`
	ArticleID  uuid.UUID `db:"articleId"`
	IpAddress  string    `db:"ipaddress"`
	SearchWord string    `db:"search_word"`
	API        string    `db:"api"`
	IsError    bool      `db:"is_error"`
}

func (repo *Repository) GetAnalysisByID(ctx context.Context, id string) (Analysis, error) {
	var analysis Analysis
	err := repo.db.GetContext(ctx, &analysis, "SELECT * FROM analysis WHERE id = ?", id)
	return analysis, err
}

func (repo *Repository) GetAnalysis(ctx context.Context, limitNumber *int) ([]Analysis, error) {
	var analysis []Analysis
	if limitNumber != nil {
		err := repo.db.SelectContext(ctx, &analysis, "SELECT * FROM analysis LIMIT ?", limitNumber)
		return analysis, err
	} else {
		err := repo.db.SelectContext(ctx, &analysis, "SELECT * FROM analysis")
		return analysis, err
	}
}

func (repo *Repository) CreateAnalysis(ctx context.Context, analysis Analysis) error {
	_, err := repo.db.NamedExecContext(ctx, "INSERT INTO analysis (id, timestamp, articleId, ipaddress, search_word, api, is_error) VALUES (:id, :timestamp, :articleId, :ipaddress, :search_word, :api, :is_error)", analysis)
	return err
}

func (repo *Repository) UpdateAnalysis(ctx context.Context, analysis Analysis) error {
	_, err := repo.db.NamedExecContext(ctx, "UPDATE analysis SET timestamp = :timestamp, articleId = :articleId, ipaddress = :ipaddress, search_word = :search_word, api = :api, is_error = :is_error WHERE id = :id", analysis)
	return err
}

func (repo *Repository) DeleteAnalysis(ctx context.Context, id string) error {
	_, err := repo.db.ExecContext(ctx, "DELETE FROM analysis WHERE id = ?", id)
	return err

}

func (repo *Repository) GetAnalysisByArticle(ctx context.Context, articleId int, limitNumber *int) ([]Analysis, error) {
	var analysis []Analysis
	if limitNumber != nil {
		err := repo.db.SelectContext(ctx, &analysis, "SELECT * FROM analysis WHERE articleId = ? LIMIT ?", articleId, limitNumber)
		return analysis, err
	} else {
		err := repo.db.SelectContext(ctx, &analysis, "SELECT * FROM analysis WHERE articleId = ?", articleId)
		return analysis, err
	}
}

func (repo *Repository) GetAnalysisByIpAddress(ctx context.Context, ipaddress string, limitNumber *int) ([]Analysis, error) {
	var analysis []Analysis
	if limitNumber != nil {
		err := repo.db.SelectContext(ctx, &analysis, "SELECT * FROM analysis WHERE ipaddress = ? LIMIT ?", ipaddress, limitNumber)
		return analysis, err
	} else {
		err := repo.db.SelectContext(ctx, &analysis, "SELECT * FROM analysis WHERE ipaddress = ?", ipaddress)
		return analysis, err
	}
}

func (repo *Repository) GetAnalysisBySearchWord(ctx context.Context, searchWord string, limitNumber *int) ([]Analysis, error) {
	var analysis []Analysis
	if limitNumber != nil {
		err := repo.db.SelectContext(ctx, &analysis, "SELECT * FROM analysis WHERE search_word = ? LIMIT ?", searchWord, limitNumber)
		return analysis, err
	} else {
		err := repo.db.SelectContext(ctx, &analysis, "SELECT * FROM analysis WHERE search_word = ?", searchWord)
		return analysis, err
	}
}

func (repo *Repository) GetAnalysisByAPI(ctx context.Context, api string, limitNumber *int) ([]Analysis, error) {
	var analysis []Analysis
	if limitNumber != nil {
		err := repo.db.SelectContext(ctx, &analysis, "SELECT * FROM analysis WHERE api = ? LIMIT ?", api, limitNumber)
		return analysis, err
	} else {
		err := repo.db.SelectContext(ctx, &analysis, "SELECT * FROM analysis WHERE api = ?", api)
		return analysis, err
	}
}

func (repo *Repository) GetAnalysisByIsError(ctx context.Context, isError bool, limitNumber *int) ([]Analysis, error) {
	var analysis []Analysis
	if limitNumber != nil {
		err := repo.db.SelectContext(ctx, &analysis, "SELECT * FROM analysis WHERE is_error = ? LIMIT ?", isError, limitNumber)
		return analysis, err
	} else {
		err := repo.db.SelectContext(ctx, &analysis, "SELECT * FROM analysis WHERE is_error = ?", isError)
		return analysis, err
	}
}

func (repo *Repository) GetAnalysisByDate(ctx context.Context, start string, end string, limitNumber *int) ([]Analysis, error) {
	var analysis []Analysis
	if limitNumber != nil {
		err := repo.db.SelectContext(ctx, &analysis, "SELECT * FROM analysis WHERE timestamp BETWEEN ? AND ? LIMIT ?", start, end, limitNumber)
		return analysis, err
	} else {
		err := repo.db.SelectContext(ctx, &analysis, "SELECT * FROM analysis WHERE timestamp BETWEEN ? AND ?", start, end)
		return analysis, err
	}
}
