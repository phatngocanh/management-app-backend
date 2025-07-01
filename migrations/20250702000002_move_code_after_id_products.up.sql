-- Move code column to be positioned after id column in products table
ALTER TABLE products MODIFY COLUMN code VARCHAR(10) NOT NULL COMMENT 'Mã sản phẩm (SP00001)' AFTER id; 