CREATE TABLE IF NOT EXISTS category (
	id 		SERIAL PRIMARY KEY,
	name	varchar(50) UNIQUE NOT NULL,
	description		varchar(200) NOT NULL
);