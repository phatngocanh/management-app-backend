-- Keep existing customers table but add ERP enhancements
CREATE TABLE customers (
    id INT AUTO_INCREMENT PRIMARY KEY,
    company_id INT NOT NULL DEFAULT 1 COMMENT 'For multi-tenant support, defaulting to 1 for existing data',
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(30), -- Increased size as requested in previous migration
    address TEXT, -- Changed to TEXT for longer addresses
    email VARCHAR(100),
    customer_code VARCHAR(50) COMMENT 'Unique customer identifier',
    tax_id VARCHAR(50) COMMENT 'Customer tax ID',
    credit_limit DECIMAL(15,2) DEFAULT 0.00 COMMENT 'Credit limit for this customer',
    payment_terms VARCHAR(100) COMMENT 'Payment terms (e.g., NET 30)',
    customer_type ENUM('RETAIL', 'WHOLESALE', 'DISTRIBUTOR') DEFAULT 'RETAIL',
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE,
    UNIQUE KEY unique_customer_code_company (customer_code, company_id)
); 