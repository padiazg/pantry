CREATE INDEX idx_movements_product ON movements(product_ean13);
CREATE INDEX idx_movements_created ON movements(created_at);
CREATE INDEX idx_products_category ON products(category_id);
CREATE INDEX idx_products_active   ON products(active);
