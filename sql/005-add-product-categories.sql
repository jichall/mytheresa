ALTER TABLE products ADD COLUMN category_id INTEGER REFERENCES product_categories(id);

INSERT INTO product_categories (code, name) VALUES
('clothing', 'Clothing'),
('shoes', 'Shoes'),
('accessories', 'Accessories');

UPDATE products SET category_id = (SELECT id FROM product_categories WHERE code = 'clothing') WHERE code IN ('PROD001', 'PROD004', 'PROD007');
UPDATE products SET category_id = (SELECT id FROM product_categories WHERE code = 'shoes') WHERE code IN ('PROD002', 'PROD006');
UPDATE products SET category_id = (SELECT id FROM product_categories WHERE code = 'accessories') WHERE code IN ('PROD003', 'PROD005', 'PROD008');

CREATE VIEW vw_products AS SELECT p.code, pc.name, p.price FROM products p INNER JOIN product_categories pc ON p.category_id = pc.id;
