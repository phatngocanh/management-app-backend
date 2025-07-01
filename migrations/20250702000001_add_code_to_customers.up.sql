-- Add code field to customers table  
ALTER TABLE customers ADD COLUMN code VARCHAR(10) NOT NULL DEFAULT '' COMMENT 'Mã khách hàng (KH00001)';

-- Generate codes for existing customers using KH prefix + 5-digit customerId
UPDATE customers SET code = CONCAT('KH', LPAD(id, 5, '0'));

-- Add unique constraint for code field
ALTER TABLE customers ADD CONSTRAINT uk_customers_code UNIQUE (code); 