ALTER TABLE `products`
ADD COLUMN `type` JSON NOT NULL AFTER `image_status`;