-- +goose Up
CREATE TABLE card_requests (
    id              UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    member_id       UUID        NOT NULL REFERENCES members(id) ON DELETE CASCADE,
    vehicle_plate   VARCHAR(20) NOT NULL,
    vehicle_brand   VARCHAR(50) NOT NULL,
    vehicle_model   VARCHAR(50) NOT NULL,
    vehicle_color   VARCHAR(30) NOT NULL,
    id_card_number  VARCHAR(30) NOT NULL,
    note            TEXT,
    status          VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending','approved','rejected','blocked')),
    card_uid        VARCHAR(50) REFERENCES cards(card_uid) ON DELETE SET NULL,
    rejected_reason TEXT,
    submitted_at    TIMESTAMP   NOT NULL DEFAULT NOW(),
    reviewed_at     TIMESTAMP,
    reviewed_by     UUID        REFERENCES users(id) ON DELETE SET NULL
);

-- +goose Down
DROP TABLE IF EXISTS card_requests;
