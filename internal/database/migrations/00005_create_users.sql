-- +goose Up
CREATE TABLE users (
    id            UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    email         VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role          VARCHAR(20)  NOT NULL CHECK (role IN ('student', 'staff', 'faculty', 'admin')),
    member_id     UUID         REFERENCES members(id) ON DELETE SET NULL,
    created_at    TIMESTAMP    NOT NULL DEFAULT NOW()
);

-- Seed default admin user (password: admin123 — change immediately after first login)
-- bcrypt hash of "admin123"
INSERT INTO users (email, password_hash, role)
VALUES ('admin@bikekeeper.local', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'admin');

-- +goose Down
DROP TABLE IF EXISTS users;
