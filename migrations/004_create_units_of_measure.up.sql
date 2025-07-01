CREATE TABLE `units_of_measure` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL COMMENT 'Tên đơn vị (VD: Thùng, Cái, ML)',
  `code` varchar(20) NOT NULL COMMENT 'Mã đơn vị (VD: THUNG, CAI, ML)',
  `description` text COMMENT 'Mô tả đơn vị',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci; 