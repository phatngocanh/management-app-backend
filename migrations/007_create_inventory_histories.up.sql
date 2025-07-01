CREATE TABLE `inventory_histories` (
  `id` int NOT NULL AUTO_INCREMENT,
  `product_id` int NOT NULL,
  `quantity` int NOT NULL COMMENT 'Số lượng nhập',
  `importer_name` varchar(255) NOT NULL COMMENT 'Tên người nhập',
  `imported_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Thời gian nhập',
  `note` text COMMENT 'Ghi chú về lần nhập hàng',
  `final_quantity` int NOT NULL DEFAULT '0',
  `reference_id` int DEFAULT NULL COMMENT 'ID tham chiếu đến đơn hàng hoặc nguồn khác',
  PRIMARY KEY (`id`),
  KEY `product_id` (`product_id`),
  CONSTRAINT `inventory_histories_ibfk_1` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci; 