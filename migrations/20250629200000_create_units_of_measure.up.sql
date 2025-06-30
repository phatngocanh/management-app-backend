-- Units of measure for chemical company
CREATE TABLE units_of_measure (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL COMMENT 'Tên đơn vị (VD: Thùng, Cái, ML)',
    code VARCHAR(20) NOT NULL UNIQUE COMMENT 'Mã đơn vị (VD: THUNG, CAI, ML)',
    description TEXT COMMENT 'Mô tả đơn vị',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert basic units for chemical company
INSERT INTO units_of_measure (name, code, description) VALUES
('Thùng', 'THUNG', 'Đơn vị tính cho thùng chứa'),
('Cái', 'CAI', 'Đơn vị tính cho nhãn, nắp, chai'),
('ML', 'ML', 'Đơn vị tính cho hóa chất (đã chuẩn hóa từ L/ML)'),
('Kg', 'KG', 'Đơn vị tính theo khối lượng'),
('Mét', 'M', 'Đơn vị tính cho màng co'); 