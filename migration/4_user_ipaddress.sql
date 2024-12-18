-- +goose Up
-- usersテーブルにipaddressカラムを追加
ALTER TABLE `users` ADD COLUMN `ipaddress` VARCHAR(15) AFTER `email`;