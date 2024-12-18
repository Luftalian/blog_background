-- +goose Up
ALTER TABLE `analysis` MODIFY COLUMN `articleId` CHAR(36);