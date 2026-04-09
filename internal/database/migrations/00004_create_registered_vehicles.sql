-- +goose Up
CREATE TABLE registered_vehicles (
    id           SERIAL      PRIMARY KEY,
    plate_number VARCHAR(20) NOT NULL,
    member_id    UUID        NOT NULL REFERENCES members(id) ON DELETE CASCADE,
    description  TEXT
);

CREATE INDEX idx_registered_vehicles_member_id    ON registered_vehicles(member_id);
CREATE INDEX idx_registered_vehicles_plate_number ON registered_vehicles(plate_number);

-- +goose Down
DROP INDEX IF EXISTS idx_registered_vehicles_plate_number;
DROP INDEX IF EXISTS idx_registered_vehicles_member_id;
DROP TABLE IF EXISTS registered_vehicles;
