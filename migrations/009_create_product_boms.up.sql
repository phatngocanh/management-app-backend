CREATE TABLE `product_boms` (
  `id` int NOT NULL AUTO_INCREMENT,
  `parent_product_id` int NOT NULL COMMENT 'The finished product that is built',
  `component_product_id` int NOT NULL COMMENT 'The component/material needed',
  `quantity` decimal(10,3) NOT NULL COMMENT 'How many units of component needed to make 1 unit of parent',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_parent_component` (`parent_product_id`,`component_product_id`),
  KEY `component_product_id` (`component_product_id`),
  CONSTRAINT `product_boms_ibfk_1` FOREIGN KEY (`parent_product_id`) REFERENCES `products` (`id`) ON DELETE CASCADE,
  CONSTRAINT `product_boms_ibfk_2` FOREIGN KEY (`component_product_id`) REFERENCES `products` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci; 