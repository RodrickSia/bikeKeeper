-- +goose Up
CREATE TABLE notifications (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title      VARCHAR(200) NOT NULL,
    message    TEXT        NOT NULL,
    type       VARCHAR(20) NOT NULL CHECK (type IN ('info','warning','success','error')),
    is_read    BOOLEAN     NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP   NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS notifications;
