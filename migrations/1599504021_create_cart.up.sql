CREATE TABLE IF NOT EXISTS cart (
  id INTEGER NOT NULL,
  product_id INTEGER NOT NULL REFERENCES products(id),
  quantity INTEGER,
  PRIMARY KEY(id, product_id)
);