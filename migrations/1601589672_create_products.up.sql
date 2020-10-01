CREATE TABLE IF NOT EXISTS products (
    id 			SERIAL PRIMARY KEY,
	name 		varchar(50)  UNIQUE NOT NULL ,
	description varchar(200) NOT NULL,
	price		float NOT NULL,
	discount   	float DEFAULT 0,
	quantity	int NOT NULL ,
	tax 		float DEFAULT 0,
	category_id	BIGINT NOT NULL ,
	brand 		varchar(50) NOT NULL ,
	color 		varchar(50) NOT NULL DEFAULT '',
	size 		varchar(50) NOT NULL DEFAULT '',
	image_urls  text[],
	FOREIGN KEY(category_id) 
	REFERENCES Category(id) ON DELETE CASCADE ON UPDATE CASCADE
);
