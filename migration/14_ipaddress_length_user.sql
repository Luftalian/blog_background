-- +goose Up
ALTER TABLE `users` MODIFY COLUMN `ipaddress` VARCHAR(500);