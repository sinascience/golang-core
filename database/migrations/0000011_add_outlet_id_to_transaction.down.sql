ALTER TABLE `transactions`
DROP COLUMN `outlet_id`;
DROP INDEX `outlet_id` ON `transactions`;