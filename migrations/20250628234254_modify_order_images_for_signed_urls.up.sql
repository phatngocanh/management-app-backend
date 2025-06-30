ALTER TABLE order_images 
ADD COLUMN image_key VARCHAR(500) COMMENT 'S3 object key for the image',
ADD COLUMN signed_url_expires_at TIMESTAMP NULL COMMENT 'Thời gian hết hạn của signed URL'; 