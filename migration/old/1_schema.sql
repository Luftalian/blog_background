-- +goose Up

-- usersテーブル
CREATE TABLE `users` (
    `id` CHAR(36) NOT NULL,
    `email` VARCHAR(255) NOT NULL,
    `username` VARCHAR(50) NOT NULL,
    `password_hash` VARCHAR(255) NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uq_users_email` (`email`),
    UNIQUE KEY `uq_users_username` (`username`),
    KEY `idx_users_email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- categoriesテーブル
CREATE TABLE `categories` (
    `id` CHAR(36) NOT NULL,
    `name` VARCHAR(100) NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uq_categories_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- tagsテーブル
CREATE TABLE `tags` (
    `id` CHAR(36) NOT NULL,
    `name` VARCHAR(100) NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uq_tags_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- articlesテーブル
CREATE TABLE `articles` (
    `id` CHAR(36) NOT NULL,
    `title` VARCHAR(255) NOT NULL,
    `content` TEXT NOT NULL,
    `author_id` CHAR(36) NOT NULL,
    `category_id` CHAR(36) NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_articles_title` (`title`),
    KEY `idx_articles_category_id` (`category_id`),
    CONSTRAINT `fk_articles_users` FOREIGN KEY (`author_id`) REFERENCES `users`(`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_articles_categories` FOREIGN KEY (`category_id`) REFERENCES `categories`(`id`) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- article_tagsテーブル (多対多中間テーブル)
CREATE TABLE `article_tags` (
    `article_id` CHAR(36) NOT NULL,
    `tag_id` CHAR(36) NOT NULL,
    PRIMARY KEY (`article_id`,`tag_id`),
    KEY `idx_article_tags_article_id` (`article_id`),
    KEY `idx_article_tags_tag_id` (`tag_id`),
    CONSTRAINT `fk_article_tags_articles` FOREIGN KEY (`article_id`) REFERENCES `articles`(`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_article_tags_tags` FOREIGN KEY (`tag_id`) REFERENCES `tags`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- commentsテーブル
CREATE TABLE `comments` (
    `id` CHAR(36) NOT NULL,
    `article_id` CHAR(36) NOT NULL,
    `author_id` CHAR(36) NOT NULL,
    `content` TEXT NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_comments_article_id` (`article_id`),
    CONSTRAINT `fk_comments_articles` FOREIGN KEY (`article_id`) REFERENCES `articles`(`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_comments_users` FOREIGN KEY (`author_id`) REFERENCES `users`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- likesテーブル
CREATE TABLE `likes` (
    `id` CHAR(36) NOT NULL,
    `article_id` CHAR(36) NOT NULL,
    `user_id` CHAR(36) NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uq_likes_article_user` (`article_id`,`user_id`),
    KEY `idx_likes_article_id` (`article_id`),
    KEY `idx_likes_user_id` (`user_id`),
    CONSTRAINT `fk_likes_articles` FOREIGN KEY (`article_id`) REFERENCES `articles`(`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_likes_users` FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
