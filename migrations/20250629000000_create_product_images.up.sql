CREATE TABLE product_images (
    id INT AUTO_INCREMENT PRIMARY KEY,
    product_id INT NOT NULL,
    image_key VARCHAR(500) NOT NULL COMMENT 'S3 object key for the image',
    is_primary BOOLEAN DEFAULT FALSE COMMENT 'Is this the main product image',
    signed_url_expires_at TIMESTAMP NULL COMMENT 'Thời gian hết hạn của signed URL',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
); 