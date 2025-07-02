CREATE TABLE `orders` (
  `id` int NOT NULL AUTO_INCREMENT,
  `code` varchar(10) NOT NULL COMMENT 'Mã đơn hàng (DH00001)',
  `customer_id` int NOT NULL COMMENT 'Khách hàng',
  `order_date` datetime NOT NULL COMMENT 'Ngày đặt hàng',
  `note` text COMMENT 'Ghi chú',
  `total_original_cost` int NOT NULL COMMENT 'Tổng giá vốn sản phẩm',
  `total_sales_revenue` int NOT NULL COMMENT 'Tổng doanh thu',
  `additional_cost` int NOT NULL COMMENT 'Chi phí phát sinh',
  `additional_cost_note` text COMMENT 'Ghi chú chi phí phát sinh',
  `tax_percent` int NOT NULL COMMENT 'Thuế suất',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_orders_code` (`code`),
  CONSTRAINT `orders_ibfk_1` FOREIGN KEY (`customer_id`) REFERENCES `customers` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci; 