CREATE TABLE IF NOT EXISTS user_blacklisted_tokens(
  id SERIAL NOT NULL PRIMARY KEY,
  user_id BIGINT REFERENCES users(id),
  token TEXT,
  expiration_date TIMESTAMP
); 