-- Move code column to be positioned after id column in customers table
ALTER TABLE customers MODIFY COLUMN code VARCHAR(10) NOT NULL COMMENT 'Mã khách hàng (KH00001)' AFTER id; 
