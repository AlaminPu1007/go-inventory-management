
-- DROP TRIGGERS
DROP TRIGGER IF EXISTS set_users_updated_at ON users;
DROP TRIGGER IF EXISTS set_orders_updated_at ON orders;
DROP TRIGGER IF EXISTS set_categories_updated_at ON categories;
DROP TRIGGER IF EXISTS set_products_updated_at ON products;
DROP TRIGGER IF EXISTS set_order_items_updated_at ON order_items;

-- DROP FUNCTION
DROP FUNCTION IF EXISTS set_updated_at;

