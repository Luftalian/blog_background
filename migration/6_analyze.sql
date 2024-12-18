-- +goose Up
CREATE TABLE `analysis` (
    `id` CHAR(36) NOT NULL,
    `timestamp` TIMESTAMP,
    `articleId` INT,
    `ipaddress` VARCHAR(45),
    `search_word` VARCHAR(255),
    `api` VARCHAR(2083),
    `is_error` BOOLEAN
);

ALTER TABLE `articles` ADD COLUMN `view_count` INT;
