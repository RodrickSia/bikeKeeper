-- +goose Up
ALTER TABLE users ADD COLUMN status VARCHAR(20) NOT NULL DEFAULT 'pending_approval' CHECK (status IN ('pending_approval', 'active', 'rejected', 'suspended'));
UPDATE users SET status = 'active' WHERE email = 'admin@bikekeeper.local';

-- +goose Down
ALTER TABLE users DROP COLUMN IF EXISTS status;
