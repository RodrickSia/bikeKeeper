-- +goose Up
CREATE TABLE members (
    id         UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id VARCHAR(20)  NOT NULL UNIQUE,
    full_name  VARCHAR(100) NOT NULL,
    phone      VARCHAR(15),
    created_at TIMESTAMP    NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS members;
