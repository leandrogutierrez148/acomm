-- Insert 3 product categories
INSERT INTO product_categories (code, name) VALUES
('CLTH', 'Clothing'),
('SHS', 'Shoes'),
('ACC', 'Accessories');

-- Insert Brands
INSERT INTO brands (name, meta_tag_title, meta_tag_description, sort_order) VALUES
('Nike', 'Nike Shoes and Apparel', 'Buy Nike products online', 1),
('Adidas', 'Adidas Gear', 'Premium Adidas clothing', 2),
('Puma', 'Puma Store', 'Puma athletic wear', 3);

-- Insert 8 products (now referencing brand_id)
INSERT INTO products (brand_id, category_id) VALUES
((SELECT id FROM brands WHERE name = 'Nike'),   (SELECT id FROM product_categories WHERE code = 'CLTH')),
((SELECT id FROM brands WHERE name = 'Adidas'), (SELECT id FROM product_categories WHERE code = 'SHS')),
((SELECT id FROM brands WHERE name = 'Puma'),   (SELECT id FROM product_categories WHERE code = 'ACC')),
((SELECT id FROM brands WHERE name = 'Nike'),   (SELECT id FROM product_categories WHERE code = 'CLTH')),
((SELECT id FROM brands WHERE name = 'Adidas'), (SELECT id FROM product_categories WHERE code = 'ACC')),
((SELECT id FROM brands WHERE name = 'Puma'),   (SELECT id FROM product_categories WHERE code = 'SHS')),
((SELECT id FROM brands WHERE name = 'Nike'),   (SELECT id FROM product_categories WHERE code = 'CLTH')),
((SELECT id FROM brands WHERE name = 'Adidas'), (SELECT id FROM product_categories WHERE code = 'ACC'));

-- Insert Items (previously variants) for each product

-- Product 1: 3 items
INSERT INTO items (product_id, sku, price) VALUES
(1, 'SKU001A', 11.99),
(1, 'SKU001B', 15.99),
(1, 'SKU001C', 20.00);

-- Product 2: 2 items
INSERT INTO items (product_id, sku, price) VALUES
(2, 'SKU002A', 49.99),
(2, 'SKU002B', 59.99);

-- Product 3: 1 item
INSERT INTO items (product_id, sku, price) VALUES
(3, 'SKU003A', 8.99);

-- Product 4: 4 items
INSERT INTO items (product_id, sku, price) VALUES
(4, 'SKU004A', 15.50),
(4, 'SKU004B', 16.00),
(4, 'SKU004C', 18.00),
(4, 'SKU004D', 16.99);

-- Product 5: 6 items
INSERT INTO items (product_id, sku, price) VALUES
(5, 'SKU005A', 23.99),
(5, 'SKU005B', 25.00),
(5, 'SKU005C', 26.50),
(5, 'SKU005D', 22.99),
(5, 'SKU005E', 23.49),
(5, 'SKU005F', 28.00);

-- Product 6: No items

-- Product 7: 5 items
INSERT INTO items (product_id, sku, price) VALUES
(7, 'SKU007A', 15.00),
(7, 'SKU007B', 16.00),
(7, 'SKU007C', 17.50),
(7, 'SKU007D', 18.00),
(7, 'SKU007E', 18.75);

-- Product 8: 1 item
INSERT INTO items (product_id, sku, price) VALUES
(8, 'SKU008A', 10.49);

-- Insert sample specifications for product 1
INSERT INTO specifications (product_id, key, value) VALUES
(1, 'Material', 'Cotton'),
(1, 'Color', 'Red'),
(1, 'Size', 'M');
