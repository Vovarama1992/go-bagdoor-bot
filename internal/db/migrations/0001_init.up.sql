CREATE TABLE IF NOT EXISTS users (
  id SERIAL PRIMARY KEY,
  tg_username VARCHAR(255),
  tg_id BIGINT NOT NULL UNIQUE,
  first_name VARCHAR(100),
  last_name VARCHAR(100),
  phone_number VARCHAR(20),
  registered_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);