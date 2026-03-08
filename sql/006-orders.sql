CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    customer_email VARCHAR(256) NOT NULL,
    customer_name VARCHAR(256) NOT NULL,
    customer_address VARCHAR(256) NOT NULL,
    customer_phone VARCHAR(256) NOT NULL,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
CREATE TABLE IF NOT EXISTS item_orders (
    item_id INTEGER NOT NULL REFERENCES items(id) ON DELETE CASCADE,
    order_id INTEGER NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    quantity INTEGER NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    PRIMARY KEY (item_id, order_id)
);