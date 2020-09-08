CREATE TABLE cart (
  id integer NOT NULL PRIMARY KEY,
  product_id integer REFERENCES products(id),
  quantity integer
);
