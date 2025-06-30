CREATE TABLE product_images (
    id INT AUTO_INCREMENT PRIMARY KEY,
    product_id INT NOT NULL,
    image_url TEXT NOT NULL COMMENT 'Full S3 URL for the image',
    image_key VARCHAR(500) COMMENT 'S3 object key for the image',
    is_primary BOOLEAN DEFAULT FALSE COMMENT 'Is this the main product image',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
); 