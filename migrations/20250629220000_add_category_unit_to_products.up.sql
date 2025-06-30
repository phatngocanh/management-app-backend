-- Add category and unit references to products table
ALTER TABLE products 
ADD COLUMN category_id INT COMMENT 'Danh mục sản phẩm',
ADD COLUMN unit_id INT COMMENT 'Đơn vị tính',
ADD COLUMN description TEXT COMMENT 'Mô tả chi tiết sản phẩm',
ADD FOREIGN KEY (category_id) REFERENCES product_categories(id),
ADD FOREIGN KEY (unit_id) REFERENCES units_of_measure(id); 