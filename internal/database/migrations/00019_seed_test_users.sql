-- +goose Up
-- Seed test users matching frontend demo accounts (password: demo123 for all)

-- Student 1: Nguyen Van An
WITH m AS (
    INSERT INTO members (student_id, full_name, phone)
    VALUES ('SV20210001', 'Nguyễn Văn An', '0901234567')
    RETURNING id
)
INSERT INTO users (email, password_hash, role, status, member_id)
SELECT 'student@university.edu.vn', '$2a$10$MzJZqWhKv7/lGKYnb1iV.ONhnv0.UyUj78cxeQlJ2uMEx34LPOng.', 'student', 'active', id
FROM m;

-- Student 2: Pham Duc Huy
WITH m AS (
    INSERT INTO members (student_id, full_name, phone)
    VALUES ('SV20210042', 'Phạm Đức Huy', '0909876543')
    RETURNING id
)
INSERT INTO users (email, password_hash, role, status, member_id)
SELECT 'student2@university.edu.vn', '$2a$10$MzJZqWhKv7/lGKYnb1iV.ONhnv0.UyUj78cxeQlJ2uMEx34LPOng.', 'student', 'active', id
FROM m;

-- Guard (Bảo vệ)
INSERT INTO users (email, password_hash, role, status)
VALUES ('guard@parksmart.vn', '$2a$10$MzJZqWhKv7/lGKYnb1iV.ONhnv0.UyUj78cxeQlJ2uMEx34LPOng.', 'staff', 'active');

-- Admin
INSERT INTO users (email, password_hash, role, status)
VALUES ('admin@parksmart.vn', '$2a$10$MzJZqWhKv7/lGKYnb1iV.ONhnv0.UyUj78cxeQlJ2uMEx34LPOng.', 'admin', 'active');

-- +goose Down
DELETE FROM users WHERE email IN (
    'student@university.edu.vn',
    'student2@university.edu.vn',
    'guard@parksmart.vn',
    'admin@parksmart.vn'
);
DELETE FROM members WHERE student_id IN ('SV20210001', 'SV20210042');
