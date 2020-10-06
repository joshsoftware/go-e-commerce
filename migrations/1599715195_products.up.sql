CREATE TABLE IF NOT EXISTS products (
	id 			SERIAL PRIMARY KEY,
	name 		varchar(50)  UNIQUE NOT NULL ,
	description 	varchar(200) NOT NULL,
	price		float NOT NULL, CHECK (price > 0),
	discount   	float NOT NULL DEFAULT 0  CHECK (discount >= 0 AND discount <= 100),
	quantity	int NOT NULL    CHECK (quantity >= 0 AND quantity <= 1000),
	tax 		float NOT NULL DEFAULT 0  CHECK (tax >= 0 AND tax <= 100),
	cid			int NOT NULL ,
	brand 		varchar(50) NOT NULL ,
	color 		varchar(50) NOT NULL DEFAULT '',
	size 		varchar(50) NOT NULL DEFAULT '',
	image_urls text[] DEFAULT NULL,
	FOREIGN KEY(cid) 
	REFERENCES Category(cid) ON DELETE CASCADE ON UPDATE CASCADE
);
