-- +goose Up

-- 1) コメントをすべて削除
TRUNCATE TABLE `comments`;

-- 2) 記事を書いたことがあるadminユーザー以外を削除
DELETE FROM `users`
WHERE NOT (
  `is_admin` = TRUE
  AND `id` IN (SELECT `author_id` FROM `articles`)
);

-- 3) 記事を書いたadminユーザーの名前をLuftalianに変更
UPDATE `users`
SET `username` = 'Luftalian'
WHERE `is_admin` = TRUE
  AND `id` IN (SELECT `author_id` FROM `articles`);