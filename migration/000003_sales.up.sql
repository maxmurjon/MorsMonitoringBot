CREATE TABLE sales (
    id BIGSERIAL PRIMARY KEY,
    seller_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    barrel_id BIGINT REFERENCES barrels(id) ON DELETE CASCADE,
    quantity DECIMAL(10,2) NOT NULL, -- Sotilgan hajm (litrda)
    sold_at TIMESTAMP DEFAULT NOW()
);
