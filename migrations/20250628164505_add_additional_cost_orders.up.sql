ALTER TABLE orders 
ADD COLUMN additional_cost DECIMAL(15,2) DEFAULT 0.00 COMMENT 'Chi phí phát sinh (VND)'; 