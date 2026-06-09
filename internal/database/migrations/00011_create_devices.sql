-- +goose Up
CREATE TABLE devices (
    id               UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    name             VARCHAR(100) NOT NULL,
    type             VARCHAR(20)  NOT NULL CHECK (type IN ('camera','barrier','rfid_reader','sensor')),
    location_label   VARCHAR(200) NOT NULL,
    status           VARCHAR(20)  NOT NULL DEFAULT 'online' CHECK (status IN ('online','offline','warning','maintenance')),
    ip_address       VARCHAR(45),
    firmware_version VARCHAR(20),
    notes            TEXT,
    installed_at     DATE         NOT NULL DEFAULT NOW(),
    last_seen        TIMESTAMP    NOT NULL DEFAULT NOW()
);

CREATE TABLE device_alerts (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    device_id   UUID        NOT NULL REFERENCES devices(id) ON DELETE CASCADE,
    message     TEXT        NOT NULL,
    severity    VARCHAR(10) NOT NULL CHECK (severity IN ('low','medium','high')),
    created_at  TIMESTAMP   NOT NULL DEFAULT NOW(),
    resolved_at TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS device_alerts;
DROP TABLE IF EXISTS devices;
