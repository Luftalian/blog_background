-- +goose Up

UPDATE `articles`
SET `author_id` = (
  SELECT `id`
  FROM `users`
  WHERE `is_admin` = TRUE
  LIMIT 1
)
WHERE (SELECT COUNT(*) FROM `users` WHERE `is_admin` = TRUE) > 0;