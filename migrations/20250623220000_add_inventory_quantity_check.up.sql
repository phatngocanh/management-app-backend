ALTER TABLE inventory
ADD CONSTRAINT check_inventory_quantity CHECK (quantity >= 0); 