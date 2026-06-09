-- +goose Up
CREATE TABLE support_tickets (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category    VARCHAR(30) NOT NULL CHECK (category IN ('wallet_issue','card_issue','staff_attitude','other')),
    subject     VARCHAR(200) NOT NULL,
    description TEXT        NOT NULL,
    status      VARCHAR(20) NOT NULL DEFAULT 'open' CHECK (status IN ('open','in_progress','resolved','closed')),
    created_at  TIMESTAMP   NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP   NOT NULL DEFAULT NOW()
);

CREATE TABLE ticket_responses (
    id         UUID      PRIMARY KEY DEFAULT gen_random_uuid(),
    ticket_id  UUID      NOT NULL REFERENCES support_tickets(id) ON DELETE CASCADE,
    sender_id  UUID      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    message    TEXT      NOT NULL,
    is_admin   BOOLEAN   NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS ticket_responses;
DROP TABLE IF EXISTS support_tickets;
