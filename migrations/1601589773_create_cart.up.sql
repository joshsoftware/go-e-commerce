CREATE TABLE IF NOT EXISTS cart (
    id integer NOT NULL ,
    product_id BIGINT NOT NULL,
    quantity integer,
    FOREIGN KEY(product_id) 
    REFERENCES products(id) ON DELETE CASCADE,
    PRIMARY KEY(id, product_id)
);