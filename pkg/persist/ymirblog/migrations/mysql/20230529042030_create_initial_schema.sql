-- Create "users" table
CREATE TABLE `users` (`id` bigint NOT NULL AUTO_INCREMENT, `name` varchar(255) NOT NULL, `email` varchar(255) NOT NULL, PRIMARY KEY (`id`)) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "ymirs" table
CREATE TABLE `ymirs` (`id` bigint NOT NULL AUTO_INCREMENT, `version` varchar(255) NOT NULL DEFAULT 'alpha-test-dev1', PRIMARY KEY (`id`)) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "articles" table
CREATE TABLE `articles` (`id` bigint NOT NULL AUTO_INCREMENT, `title` varchar(255) NOT NULL DEFAULT 'untitled', `body` varchar(255) NULL, `user_articles` bigint NULL, PRIMARY KEY (`id`), INDEX `articles_users_articles` (`user_articles`), CONSTRAINT `articles_users_articles` FOREIGN KEY (`user_articles`) REFERENCES `users` (`id`) ON UPDATE RESTRICT ON DELETE SET NULL) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "tags" table
CREATE TABLE `tags` (`id` bigint NOT NULL AUTO_INCREMENT, `name` varchar(255) NOT NULL, PRIMARY KEY (`id`), UNIQUE INDEX `name` (`name`)) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "tag_articles" table
CREATE TABLE `tag_articles` (`tag_id` bigint NOT NULL, `article_id` bigint NOT NULL, PRIMARY KEY (`tag_id`, `article_id`), INDEX `tag_articles_article_id` (`article_id`), CONSTRAINT `tag_articles_article_id` FOREIGN KEY (`article_id`) REFERENCES `articles` (`id`) ON UPDATE RESTRICT ON DELETE CASCADE, CONSTRAINT `tag_articles_tag_id` FOREIGN KEY (`tag_id`) REFERENCES `tags` (`id`) ON UPDATE RESTRICT ON DELETE CASCADE) CHARSET utf8mb4 COLLATE utf8mb4_bin;
