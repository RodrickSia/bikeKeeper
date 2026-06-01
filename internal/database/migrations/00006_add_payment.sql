-- +goose Up
ALTER TABLE cards ADD COLUMN balance DECIMAL(10,2) NOT NULL DEFAULT 0;

CREATE TABLE transactions (
    id             UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    card_uid       VARCHAR(50)  NOT NULL REFERENCES cards(card_uid) ON DELETE CASCADE,
    amount         DECIMAL(10,2) NOT NULL,
    type           VARCHAR(20)  NOT NULL CHECK (type IN ('deposit', 'parking_fee')),
    payment_method VARCHAR(20)  CHECK (payment_method IN ('card_balance', 'cash')),
    session_id     BIGINT       REFERENCES parking_sessions(id) ON DELETE SET NULL,
    created_at     TIMESTAMP    NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS transactions;
ALTER TABLE cards DROP COLUMN IF EXISTS balance;
