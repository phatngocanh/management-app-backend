-- Product categories for chemical company materials
CREATE TABLE product_categories (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL COMMENT 'Tên danh mục',
    code VARCHAR(50) NOT NULL UNIQUE COMMENT 'Mã danh mục',
    description TEXT COMMENT 'Mô tả danh mục',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert categories for chemical company
INSERT INTO product_categories (name, code, description) VALUES
('Thùng chứa', 'THUNG', 'Các loại thùng dùng để chứa hóa chất'),
('Nhãn mác', 'NHAN', 'Các loại nhãn dán sản phẩm'),
('Hóa chất', 'HOA_CHAT', 'Các loại hóa chất nguyên liệu'),
('Màng co', 'MANG_CO', 'Màng co để đóng gói sản phẩm'),
('Nắp chai', 'NAP', 'Các loại nắp đậy chai'),
('Chai lọ', 'CHAI', 'Các loại chai đựng sản phẩm'),
('Chai thành phẩm', 'CHAI_THANH_PHAM', 'Chai thành phẩm đã hoàn thiện'),
('Thùng thành phẩm', 'THUNG_THANH_PHAM', 'Thùng thành phẩm đã hoàn thiện'); 