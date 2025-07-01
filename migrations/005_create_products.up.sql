CREATE TABLE `products` (
  `id` int NOT NULL AUTO_INCREMENT,
  `code` varchar(10) NOT NULL COMMENT 'Mã sản phẩm (SP00001)',
  `name` varchar(255) NOT NULL,
  `cost` decimal(10,3) NOT NULL DEFAULT '0.000' COMMENT 'Giá vốn của sản phẩm (VND)',
  `category_id` int DEFAULT NULL COMMENT 'Danh mục sản phẩm',
  `unit_id` int DEFAULT NULL COMMENT 'Đơn vị tính',
  `description` text COMMENT 'Mô tả chi tiết sản phẩm',
  `operation_type` enum('MANUFACTURING','PACKAGING','PURCHASE') NOT NULL DEFAULT 'MANUFACTURING' COMMENT 'Loại sản phẩm: MANUFACTURING (sản xuất) hoặc PACKAGING (đóng gói) hoặc PURCHASE (mua nguyên liệu)',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_products_code` (`code`),
  KEY `category_id` (`category_id`),
  KEY `unit_id` (`unit_id`),
  CONSTRAINT `products_ibfk_1` FOREIGN KEY (`category_id`) REFERENCES `product_categories` (`id`),
  CONSTRAINT `products_ibfk_2` FOREIGN KEY (`unit_id`) REFERENCES `units_of_measure` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci; 