-- Remove operation type from product_boms table (moving to products table instead)
ALTER TABLE product_boms DROP COLUMN operation_type; 