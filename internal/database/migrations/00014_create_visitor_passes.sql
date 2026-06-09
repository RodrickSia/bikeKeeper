-- +goose Up
CREATE TABLE visitor_passes (
    id             UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id        UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    visitor_name   VARCHAR(100) NOT NULL,
    visitor_phone  VARCHAR(20),
    vehicle_plate  VARCHAR(20) NOT NULL,
    valid_date     DATE        NOT NULL,
    status         VARCHAR(20) NOT NULL DEFAULT 'valid' CHECK (status IN ('valid','used','expired','cancelled')),
    qr_code_data   VARCHAR(200) NOT NULL,
    created_at     TIMESTAMP   NOT NULL DEFAULT NOW(),
    used_at        TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS visitor_passes;
