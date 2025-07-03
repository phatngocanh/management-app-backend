-- Drop the existing foreign key constraint
ALTER TABLE `order_items` DROP FOREIGN KEY `order_items_ibfk_1`;

-- Add the foreign key constraint back with ON DELETE CASCADE
ALTER TABLE `order_items` ADD CONSTRAINT `order_items_ibfk_1` 
    FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`) ON DELETE CASCADE; 