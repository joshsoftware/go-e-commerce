-- use Serial id for auto increment
CREATE TABLE IF NOT EXISTS products (
	id SERIAL PRIMARY KEY,
	name varchar(50)   NOT NULL,
	description  varchar(200)   NOT NULL,
	price	float NOT NULL,
	discount   float NOT NULL,
	quantity	int NOT NULL ,
	category_id	int NOT NULL ,
	brand varchar(50) NOT NULL,
	color varchar(50) DEFAULT '',
	size varchar(50) DEFAULT '',
	FOREIGN KEY(category_id) 
	REFERENCES Category(id) ON DELETE CASCADE ON UPDATE CASCADE
);
