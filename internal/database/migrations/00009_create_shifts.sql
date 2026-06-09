-- +goose Up
CREATE TABLE shifts (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(100) NOT NULL,
    type        VARCHAR(20)  NOT NULL CHECK (type IN ('morning','afternoon','evening','night')),
    start_time  VARCHAR(5)   NOT NULL,
    end_time    VARCHAR(5)   NOT NULL,
    date        DATE         NOT NULL,
    status      VARCHAR(20)  NOT NULL DEFAULT 'scheduled' CHECK (status IN ('scheduled','active','completed','cancelled')),
    notes       TEXT,
    created_at  TIMESTAMP    NOT NULL DEFAULT NOW()
);

CREATE TABLE shift_assignments (
    shift_id    UUID NOT NULL REFERENCES shifts(id) ON DELETE CASCADE,
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    assigned_at TIMESTAMP NOT NULL DEFAULT NOW(),
    assigned_by UUID REFERENCES users(id) ON DELETE SET NULL,
    PRIMARY KEY (shift_id, user_id)
);

-- +goose Down
DROP TABLE IF EXISTS shift_assignments;
DROP TABLE IF EXISTS shifts;
