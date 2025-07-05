ALTER TABLE `orders` 
ADD COLUMN `delivery_status` VARCHAR(20) 
CHECK (delivery_status IN ('PENDING', 'DELIVERED', 'UNPAID', 'COMPLETED'))
COMMENT 'Trạng thái giao hàng'; 