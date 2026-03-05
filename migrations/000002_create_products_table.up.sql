CREATE TABLE products (
    ean13         CHAR(13) PRIMARY KEY,
    name          VARCHAR(200) NOT NULL,
    description   TEXT,
    unit          VARCHAR(30) NOT NULL,
    min_stock     DECIMAL(10,3) NOT NULL DEFAULT 0,
    current_stock DECIMAL(10,3) NOT NULL DEFAULT 0,
    category_id   VARCHAR(36) REFERENCES categories(id),
    active        BOOLEAN NOT NULL DEFAULT TRUE,
    created_at    TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMP NOT NULL DEFAULT NOW()
);
