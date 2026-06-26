-- +goose Up
CREATE TABLE monthly_passes (
    id             UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id        UUID         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    vehicle_id     VARCHAR(50)  NOT NULL,
    vehicle_plate  VARCHAR(20)  NOT NULL,
    vehicle_brand  VARCHAR(100) NOT NULL,
    month          VARCHAR(7)   NOT NULL,
    start_date     DATE         NOT NULL,
    end_date       DATE         NOT NULL,
    price          DECIMAL(10,2) NOT NULL,
    status         VARCHAR(20)  NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'expired', 'pending')),
    is_auto_renew  BOOLEAN      NOT NULL DEFAULT false,
    purchased_at   TIMESTAMP    NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS monthly_passes;
