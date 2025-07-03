-- First drop the foreign key constraint
ALTER TABLE `order_items` DROP FOREIGN KEY `order_items_bom_parent_fk`;

-- Then drop the column
ALTER TABLE `order_items` DROP COLUMN `bom_parent_id`;