ALTER TABLE `users`
ADD COLUMN `avatar_url` VARCHAR(255) NULL DEFAULT NULL AFTER `password`,
ADD COLUMN `image_status` VARCHAR(20) NOT NULL DEFAULT 'default' AFTER `avatar_url`;