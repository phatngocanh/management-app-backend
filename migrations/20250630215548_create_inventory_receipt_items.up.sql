CREATE TABLE inventory_receipt_items (
    id INT AUTO_INCREMENT PRIMARY KEY,
    inventory_receipt_id INT NOT NULL,
    product_id INT NOT NULL,
    quantity INT NOT NULL,
    unit_cost DECIMAL(10,2),
    notes TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (inventory_receipt_id) REFERENCES inventory_receipt(id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(id),
);
