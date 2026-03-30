CREATE TABLE
    IF NOT EXISTS orders (
        id VARCHAR(36) PRIMARY KEY,
        order_sn VARCHAR(100) NOT NULL,
        shop_id VARCHAR(100) NOT NULL,
        marketplace_status VARCHAR(50) NOT NULL,
        shipping_status VARCHAR(50) NOT NULL,
        wms_status VARCHAR(50) NOT NULL,
        tracking_number VARCHAR(100) NULL,
        total_amount DECIMAL(15, 2) NOT NULL DEFAULT 0,
        raw_marketplace_payload JSONB NOT NULL DEFAULT '{}',
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        CONSTRAINT uq_orders_order_sn UNIQUE (order_sn),
        CONSTRAINT chk_orders_total_amount CHECK (total_amount >= 0)
    );

CREATE INDEX idx_orders_shop_id ON orders (shop_id);

CREATE INDEX idx_orders_marketplace_status ON orders (marketplace_status);

CREATE INDEX idx_orders_shipping_status ON orders (shipping_status);

CREATE INDEX idx_orders_wms_status ON orders (wms_status);

CREATE INDEX idx_orders_created_at ON orders (created_at DESC);