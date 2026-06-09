-- +goose Up
CREATE TABLE incidents (
    id            UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    reported_by   UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    vehicle_plate VARCHAR(20),
    type          VARCHAR(30) NOT NULL CHECK (type IN ('wrong_parking','damaged_vehicle','suspicious','unregistered','other')),
    description   TEXT        NOT NULL,
    location      VARCHAR(200),
    status        VARCHAR(20) NOT NULL DEFAULT 'open' CHECK (status IN ('open','resolved','escalated')),
    resolved_at   TIMESTAMP,
    resolved_note TEXT,
    created_at    TIMESTAMP   NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS incidents;
