CREATE TABLE IF NOT EXISTS productimages (
	product_id int NOT NULL,
	url varchar(50) NOT NULL,
	description text NOT NULL,
	FOREIGN KEY(product_id) 
	REFERENCES products(id) ON DELETE CASCADE ON UPDATE CASCADE
);
