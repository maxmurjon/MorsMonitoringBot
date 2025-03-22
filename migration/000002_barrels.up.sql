CREATE TABLE barrels (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL, -- Barrel nomi (ID yoki unikal kod)
    volume_liters DECIMAL(10,2) NOT NULL, -- Barrel hajmi (litrda)
    current_volume DECIMAL(10,2) DEFAULT 0, -- Hozirgi mavjud mors miqdori
    location_name VARCHAR(100) NOT NULL, -- Barrel joylashgan joy nomi
    latitude DECIMAL(10,6) NOT NULL, -- Barrel joylashuvi (GPS)
    longitude DECIMAL(10,6) NOT NULL,
    assigned_seller_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
