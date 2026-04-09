-- +goose Up
CREATE TABLE cards (
    card_uid  VARCHAR(50) PRIMARY KEY,
    card_type VARCHAR(20) NOT NULL CHECK (card_type IN ('monthly', 'casual')),
    member_id UUID        REFERENCES members(id) ON DELETE SET NULL,
    is_inside BOOLEAN     NOT NULL DEFAULT FALSE,
    status    VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'blocked', 'lost'))
);

-- +goose Down
DROP TABLE IF EXISTS cards;
