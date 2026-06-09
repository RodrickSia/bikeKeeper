-- +goose Up
INSERT INTO cards (card_uid, card_type, status)
VALUES ('CARD-MOCK', 'casual', 'active')
ON CONFLICT (card_uid) DO NOTHING;

-- +goose Down
DELETE FROM cards WHERE card_uid = 'CARD-MOCK';
