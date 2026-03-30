CREATE TABLE
    IF NOT EXISTS order_items (
        id VARCHAR(36) PRIMARY KEY,
        order_sn VARCHAR(100) NOT NULL,
        sku VARCHAR(100) NOT NULL,
        quantity INT NOT NULL,
        price DECIMAL(15, 2) NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        CONSTRAINT fk_order_items_order_sn FOREIGN KEY (order_sn) REFERENCES orders (order_sn) ON DELETE CASCADE ON UPDATE CASCADE,
        CONSTRAINT chk_order_items_quantity CHECK (quantity > 0),
        CONSTRAINT chk_order_items_price CHECK (price >= 0)
    );

CREATE INDEX idx_order_items_order_sn ON order_items (order_sn);

CREATE INDEX idx_order_items_sku ON order_items (sku);