CREATE TABLE IF NOT EXISTS users (
  id SERIAL NOT NULL PRIMARY KEY,
  first_name VARCHAR(255) NOT NULL,
  last_name VARCHAR(255),
  email VARCHAR(255) NOT NULL UNIQUE,
  mobile VARCHAR(20),
  country VARCHAR(100),
  state VARCHAR(100),
  city VARCHAR(100),
  address TEXT,
  password TEXT,  
  created_at TIMESTAMP DEFAULT (NOW() AT TIME ZONE 'UTC')
); 