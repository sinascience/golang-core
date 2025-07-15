ALTER TABLE `transactions`
ADD COLUMN `outlet_id` CHAR(36) NOT NULL AFTER `user_id`,
