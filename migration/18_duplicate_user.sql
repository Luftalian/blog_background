-- +goose Up
ALTER TABLE `users` DROP INDEX `uq_users_username`;
ALTER TABLE `users` ADD INDEX `idx_users_username` (`username`);