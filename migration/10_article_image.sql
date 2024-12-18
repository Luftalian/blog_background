-- +goose Up

ALTER TABLE `articles` ADD COLUMN `image_url` VARCHAR(2083);