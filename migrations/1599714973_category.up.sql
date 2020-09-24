CREATE TABLE IF NOT EXISTS category (
	cid 		SERIAL PRIMARY KEY,
	cname	varchar(50) UNIQUE NOT NULL,
	description		varchar(200) NOT NULL
);