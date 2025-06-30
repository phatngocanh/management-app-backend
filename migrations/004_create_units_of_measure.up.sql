CREATE TABLE units_of_measure (
    id INT AUTO_INCREMENT PRIMARY KEY,
    company_id INT NOT NULL,
    name VARCHAR(100) NOT NULL COMMENT 'Unit name (e.g., Kilogram, Piece, Meter)',
    code VARCHAR(20) NOT NULL COMMENT 'Unit code (e.g., KG, PCS, M)',
    base_unit_id INT NULL COMMENT 'Reference to base unit for conversion',
    conversion_factor DECIMAL(15,6) DEFAULT 1.0 COMMENT 'Factor to convert to base unit',
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE,
    FOREIGN KEY (base_unit_id) REFERENCES units_of_measure(id) ON DELETE SET NULL,
    UNIQUE KEY unique_code_company (code, company_id)
); 