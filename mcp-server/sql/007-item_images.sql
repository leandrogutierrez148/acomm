CREATE TABLE IF NOT EXISTS item_images (
    id SERIAL PRIMARY KEY,
    item_id INTEGER NOT NULL REFERENCES items(id) ON DELETE CASCADE,
    archive_id INTEGER,
    sku_id INTEGER,
    name VARCHAR(255),
    url TEXT NOT NULL,
    file_location TEXT,
    label VARCHAR(255),
    text TEXT,
    is_main BOOLEAN NOT NULL DEFAULT FALSE,
    position INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);