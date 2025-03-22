CREATE TABLE courier_assignments (
    id BIGSERIAL PRIMARY KEY,
    courier_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    barrel_id BIGINT REFERENCES barrels(id) ON DELETE CASCADE,
    status VARCHAR(20) CHECK (status IN ('pending', 'completed')) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT NOW(),
    completed_at TIMESTAMP NULL
);
