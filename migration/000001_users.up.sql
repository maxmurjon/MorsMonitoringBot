CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    username VARCHAR(100),
    phone VARCHAR(20) UNIQUE NOT NULL,
    telegram_id BIGINT UNIQUE NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'seller', -- 'seller', 'courier', 'admin'
    is_verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);