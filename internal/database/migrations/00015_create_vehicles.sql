-- +goose Up
CREATE TABLE vehicles (
    id            UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    license_plate VARCHAR(20) NOT NULL UNIQUE,
    brand         VARCHAR(50) NOT NULL,
    model         VARCHAR(50) NOT NULL,
    color         VARCHAR(30) NOT NULL,
    owner_id      UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    is_active     BOOLEAN     NOT NULL DEFAULT TRUE,
    registered_at TIMESTAMP   NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS vehicles;
