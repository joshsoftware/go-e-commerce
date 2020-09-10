CREATE TABLE cart (
  id integer NOT NULL ,
  product_id integer REFERENCES products(id),
  quantity integer
);
