CREATE TABLE IF NOT EXISTS productimages (
	product_id int NOT NULL,
	url text NOT NULL,
	FOREIGN KEY(product_id) 
	REFERENCES products(id) ON DELETE CASCADE ON UPDATE CASCADE
);
