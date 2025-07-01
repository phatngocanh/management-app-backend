-- Add code field to products table
ALTER TABLE products ADD COLUMN code VARCHAR(10) NOT NULL DEFAULT '' COMMENT 'Mã sản phẩm (SP00001)';

-- Generate codes for existing products using SP prefix + 5-digit productId
UPDATE products SET code = CONCAT('SP', LPAD(id, 5, '0'));

-- Add unique constraint for code field
ALTER TABLE products ADD CONSTRAINT uk_products_code UNIQUE (code); 