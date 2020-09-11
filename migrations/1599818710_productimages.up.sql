CREATE TABLE IF NOT EXISTS productImages (
    id SERIAL NOT NULL PRIMARY KEY, 
	product_id BIGINT NOT NULL,
	url TEXT NOT NULL,
	FOREIGN KEY(product_Id) 
    REFERENCES products(id) ON DELETE CASCADE ON UPDATE CASCADE
);
