CREATE TABLE `inventory_receipt_items` (
  `id` int NOT NULL AUTO_INCREMENT,
  `inventory_receipt_id` int NOT NULL,
  `product_id` int NOT NULL,
  `quantity` int NOT NULL,
  `unit_cost` decimal(10,3) DEFAULT NULL,
  `notes` text,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `inventory_receipt_id` (`inventory_receipt_id`),
  KEY `product_id` (`product_id`),
  CONSTRAINT `inventory_receipt_items_ibfk_1` FOREIGN KEY (`inventory_receipt_id`) REFERENCES `inventory_receipts` (`id`) ON DELETE CASCADE,
  CONSTRAINT `inventory_receipt_items_ibfk_2` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci; 