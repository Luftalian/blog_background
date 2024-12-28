-- +goose Up
ALTER TABLE `articles`
  ADD INDEX `idx_articles_created_at_desc` (`created_at` DESC);