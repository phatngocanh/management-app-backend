-- Bill of Materials: defines which products are made from other products
CREATE TABLE product_boms (
    id INT AUTO_INCREMENT PRIMARY KEY,
    parent_product_id INT NOT NULL COMMENT 'The finished product that is built',
    component_product_id INT NOT NULL COMMENT 'The component/material needed',
    quantity DECIMAL(10,3) NOT NULL COMMENT 'How many units of component needed to make 1 unit of parent',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (parent_product_id) REFERENCES products(id) ON DELETE CASCADE,
    FOREIGN KEY (component_product_id) REFERENCES products(id) ON DELETE CASCADE,
    UNIQUE KEY unique_parent_component (parent_product_id, component_product_id)
); 