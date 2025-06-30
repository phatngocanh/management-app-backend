CREATE TABLE companies (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50) UNIQUE NOT NULL COMMENT 'Company code for identification',
    address TEXT,
    phone VARCHAR(20),
    email VARCHAR(100),
    tax_id VARCHAR(50) COMMENT 'Tax identification number',
    logo_url TEXT COMMENT 'Company logo URL in S3',
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
); 