-- +goose Up
CREATE TABLE parking_lots (
    id               UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    name             VARCHAR(100) NOT NULL UNIQUE,
    address          VARCHAR(200) NOT NULL,
    type             VARCHAR(20)  NOT NULL CHECK (type IN ('indoor','outdoor','multi_level')),
    status           VARCHAR(20)  NOT NULL DEFAULT 'active' CHECK (status IN ('active','inactive','maintenance')),
    total_capacity   INTEGER      NOT NULL DEFAULT 0,
    current_occupancy INTEGER     NOT NULL DEFAULT 0,
    open_time        VARCHAR(5)   NOT NULL DEFAULT '06:00',
    close_time       VARCHAR(5)   NOT NULL DEFAULT '22:00',
    contact_phone    VARCHAR(20),
    manager_name     VARCHAR(100),
    description      TEXT,
    created_at       TIMESTAMP    NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMP    NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS parking_lots;
