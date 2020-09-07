CREATE TABLE cart (
  id integer NOT NULL,
  product_id integer NOT NULL REFERENCES products(id),
  quantity integer
);