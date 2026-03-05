CREATE TABLE movements (
    id            VARCHAR(36) PRIMARY KEY,
    product_ean13 CHAR(13) NOT NULL REFERENCES products(ean13),
    type          VARCHAR(3) NOT NULL CHECK (type IN ('in', 'out')),
    quantity      DECIMAL(10,3) NOT NULL CHECK (quantity > 0),
    reason        VARCHAR(200),
    notes         TEXT,
    created_by    VARCHAR(100),
    created_at    TIMESTAMP NOT NULL DEFAULT NOW()
);
