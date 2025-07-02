CREATE TABLE `product_boms` (
  `id` int NOT NULL AUTO_INCREMENT,
  `parent_product_id` int NOT NULL COMMENT 'Sản phẩm thành phẩm được tạo ra',
  `component_product_id` int NOT NULL COMMENT 'Nguyên liệu/linh kiện cần thiết',
  `quantity` decimal(10,3) NOT NULL COMMENT 'Số lượng đơn vị nguyên liệu cần để tạo 1 đơn vị sản phẩm thành phẩm',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_parent_component` (`parent_product_id`,`component_product_id`),
  KEY `component_product_id` (`component_product_id`),
  CONSTRAINT `product_boms_ibfk_1` FOREIGN KEY (`parent_product_id`) REFERENCES `products` (`id`) ON DELETE CASCADE,
  CONSTRAINT `product_boms_ibfk_2` FOREIGN KEY (`component_product_id`) REFERENCES `products` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci; 