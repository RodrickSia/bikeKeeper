-- +goose Up
CREATE TABLE parking_sessions (
    id                  BIGSERIAL    PRIMARY KEY,
    card_uid            VARCHAR(50)  NOT NULL REFERENCES cards(card_uid) ON DELETE RESTRICT,
    plate_in            VARCHAR(20),
    img_plate_in_path   TEXT,
    img_person_in_path  TEXT,
    check_in_time       TIMESTAMP    NOT NULL DEFAULT NOW(),
    plate_out           VARCHAR(20),
    img_plate_out_path  TEXT,
    img_person_out_path TEXT,
    check_out_time      TIMESTAMP,
    cost                DECIMAL(10, 2) NOT NULL DEFAULT 0,
    is_warning          BOOLEAN      NOT NULL DEFAULT FALSE,
    status              VARCHAR(20)  NOT NULL DEFAULT 'ongoing' CHECK (status IN ('ongoing', 'completed'))
);

CREATE INDEX idx_parking_sessions_card_uid      ON parking_sessions(card_uid);
CREATE INDEX idx_parking_sessions_check_in_time ON parking_sessions(check_in_time);
CREATE INDEX idx_parking_sessions_status        ON parking_sessions(status);

-- +goose Down
DROP INDEX IF EXISTS idx_parking_sessions_status;
DROP INDEX IF EXISTS idx_parking_sessions_check_in_time;
DROP INDEX IF EXISTS idx_parking_sessions_card_uid;
DROP TABLE IF EXISTS parking_sessions;
