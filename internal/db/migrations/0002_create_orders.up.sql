CREATE TYPE order_type_enum AS ENUM ('personal', 'store');
CREATE TYPE order_moderation_status AS ENUM ('PENDING', 'APPROVED', 'REJECTED');

CREATE TABLE orders (
  id SERIAL PRIMARY KEY,
  order_number VARCHAR(30) NOT NULL UNIQUE,
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  publisher_username VARCHAR(255),
  publisher_tg_id BIGINT NOT NULL,
  published_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

  origin_city VARCHAR(255) NOT NULL,
  destination_city VARCHAR(255) NOT NULL,
  start_date DATE NOT NULL,
  end_date DATE NOT NULL,

  title TEXT NOT NULL,
  description TEXT NOT NULL,

  reward NUMERIC(10,2) NOT NULL,
  deposit NUMERIC(10,2),
  cost NUMERIC(10,2),
  store_link TEXT,

  media_urls TEXT[] DEFAULT '{}',
  type order_type_enum NOT NULL DEFAULT 'personal',
  status order_moderation_status NOT NULL DEFAULT 'PENDING'
);