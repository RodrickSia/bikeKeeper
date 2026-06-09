-- +goose Up
-- Fix admin password hash (password: admin123)
-- Generated bcrypt hash with cost 10
UPDATE users SET password_hash = '$2a$10$qvNJYvV7U3v0l4ZvWKP1E.VJdGZMuqVlMzxLVPb7KVe.1L.MhOwzK' 
WHERE email = 'admin@bikekeeper.local';

-- +goose Down
-- Revert to old (incorrect) hash if needed
UPDATE users SET password_hash = '$2a$10$NivP3Eyc1dXEt7yUp3kdPOUjJma8pONGFKE/1gcLu3kL7pDnCLN5G' 
WHERE email = 'admin@bikekeeper.local';
