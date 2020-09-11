-- change datatype of description to text and remove not null
CREATE TABLE IF NOT EXISTS category (
	id 		SERIAL PRIMARY KEY,
	name	varchar(50)   NOT NULL,
	description		varchar(200) NOT NULL
);