ALTER TABLE order_items 
ADD COLUMN original_price INT DEFAULT 0 COMMENT 'Giá gốc của sản phẩm khi tạo đơn hàng (VND)'; 