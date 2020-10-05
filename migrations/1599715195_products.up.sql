CREATE TABLE IF NOT EXISTS products (
	id 			SERIAL PRIMARY KEY,
	name 		varchar(50)  UNIQUE NOT NULL ,
	description 	varchar(200) NOT NULL,
	price		float NOT NULL, CHECK (price > 0),
	discount   	float DEFAULT 0  CHECK (discount >= 0 AND discount <= 100 AND price > discount),
	quantity	int NOT NULL ,   CHECK (quantity >= 0 AND quantity <= 1000),
	tax 		float DEFAULT 0,  CHECK (tax >= 0 AND tax <= 30),
	cid			int NOT NULL ,
	brand 		varchar(50) NOT NULL ,
	color 		varchar(50) NOT NULL DEFAULT '',
	size 		varchar(50) NOT NULL DEFAULT '',
	image_urls text[],
	FOREIGN KEY(cid) 
	REFERENCES Category(cid) ON DELETE CASCADE ON UPDATE CASCADE
);
