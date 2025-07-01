CREATE TABLE `inventory` (
  `id` int NOT NULL AUTO_INCREMENT,
  `product_id` int NOT NULL,
  `quantity` int NOT NULL DEFAULT '0',
  `version` varchar(36) NOT NULL COMMENT 'UUID version của inventory, thay đổi mỗi khi có thay đổi quantity để đảm bảo consistency khi FE tạo đơn hàng',
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_product_inventory` (`product_id`),
  CONSTRAINT `inventory_ibfk_1` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`) ON DELETE CASCADE,
  CONSTRAINT `check_inventory_quantity` CHECK ((`quantity` >= 0))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci; 