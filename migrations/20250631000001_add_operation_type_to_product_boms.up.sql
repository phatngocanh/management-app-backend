-- Add operation type to product_boms to distinguish manufacturing vs packaging
ALTER TABLE product_boms 
ADD COLUMN operation_type ENUM('MANUFACTURING', 'PACKAGING', 'PURCHASE') NOT NULL DEFAULT 'MANUFACTURING' 
COMMENT 'Loại công việc: MANUFACTURING (sản xuất) hoặc PACKAGING (đóng gói) hoặc PURCHASE (mua nguyên liệu)'; 