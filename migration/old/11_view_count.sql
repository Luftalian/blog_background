-- +goose Up

UPDATE `articles` 
SET `view_count` = 0 
WHERE `view_count` IS NULL;