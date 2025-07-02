CREATE TABLE `order_items` (
  `id` int NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `order_id` int NOT NULL COMMENT 'Đơn hàng',
  -- một order item chỉ có thể là 1 product hoặc 1 bom | không thể là cả 2
  `product_id` int NULL COMMENT 'Sản phẩm',
  `bom_id` int NULL COMMENT 'BOM',
  `quantity` int NOT NULL COMMENT 'Số lượng',
  `selling_price` int NOT NULL COMMENT 'Giá bán',
  `original_price` int NOT NULL COMMENT 'Giá vốn',
  `discount_percent` int NOT NULL COMMENT 'Chiết khấu',
  `final_amount` int NOT NULL COMMENT 'Tổng tiền',
  CONSTRAINT `order_items_ibfk_1` FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`),
  CONSTRAINT `order_items_ibfk_2` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`),
  CONSTRAINT `order_items_ibfk_3` FOREIGN KEY (`bom_id`) REFERENCES `product_boms` (`id`)
);