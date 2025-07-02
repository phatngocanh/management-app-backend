-- Drop the existing foreign key constraint that references product_boms(id)
ALTER TABLE `order_items` DROP FOREIGN KEY `order_items_ibfk_3`;

-- Rename bom_id to bom_parent_id for clarity
ALTER TABLE `order_items` 
CHANGE COLUMN `bom_id` `bom_parent_id` int NULL COMMENT 'Parent Product ID - sản phẩm cha cần sản xuất';

-- Add new foreign key constraint to reference products(id) instead of product_boms(id)
ALTER TABLE `order_items` 
ADD CONSTRAINT `order_items_bom_parent_fk` FOREIGN KEY (`bom_parent_id`) REFERENCES `products` (`id`); 