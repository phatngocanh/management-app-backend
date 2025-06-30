-- Add operation type to products table to distinguish manufacturing vs packaging products
ALTER TABLE products 
ADD COLUMN operation_type ENUM('MANUFACTURING', 'PACKAGING', 'PURCHASE') NOT NULL DEFAULT 'MANUFACTURING' 
COMMENT 'Loại sản phẩm: MANUFACTURING (sản xuất) hoặc PACKAGING (đóng gói) hoặc PURCHASE (mua nguyên liệu)'; 