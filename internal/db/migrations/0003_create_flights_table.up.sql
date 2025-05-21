CREATE TABLE flights (
    id SERIAL PRIMARY KEY,
    flight_number VARCHAR(20) NOT NULL UNIQUE,
    publisher_username VARCHAR(255),
    publisher_tg_id BIGINT NOT NULL,
    published_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    description TEXT,
    origin VARCHAR(255),
    destination VARCHAR(255),
    map_url TEXT,  -- 🔥 Новое поле для ссылки на карту маршрута
    status VARCHAR(30) NOT NULL CHECK (
        status IN (
            'pending',
            'approved',
            'rejected',
            'deleted'
        )
    )
);