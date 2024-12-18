-- +goose Up
-- usersテーブルの`email`と`username`と`password_hash`はNULLを許容するように変更
ALTER TABLE `users` MODIFY COLUMN `email` VARCHAR(255) DEFAULT NULL,
    MODIFY COLUMN `username` VARCHAR(50) DEFAULT NULL,
    MODIFY COLUMN `password_hash` VARCHAR(255) DEFAULT NULL;