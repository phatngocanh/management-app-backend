ALTER TABLE orders 
ADD COLUMN status_transitioned_at TIMESTAMP NULL COMMENT 'Thời gian chuyển trạng thái'; 