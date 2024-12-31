-- +goose Up
ALTER TABLE analysis MODIFY COLUMN ipaddress VARCHAR(500);