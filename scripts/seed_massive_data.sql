-- =============================================================
-- Massive Test Data Seed
-- Run: docker compose exec -T postgres psql -U bikekeeper < scripts/seed_massive_data.sql
-- =============================================================

BEGIN;

-- =============================================================
-- 1. MEMBERS (50 students)
-- =============================================================
INSERT INTO members (student_id, full_name, phone) VALUES
  ('SV20230001', 'Lê Hoàng Nam',     '0901111111'),
  ('SV20230002', 'Trần Minh Quân',   '0901111112'),
  ('SV20230003', 'Phạm Thị Hương',   '0901111113'),
  ('SV20230004', 'Nguyễn Văn Tùng',  '0901111114'),
  ('SV20230005', 'Đặng Thùy Trang',  '0901111115'),
  ('SV20230006', 'Võ Thành Đạt',     '0901111116'),
  ('SV20230007', 'Bùi Anh Tuấn',     '0901111117'),
  ('SV20230008', 'Đỗ Minh Hoàng',    '0901111118'),
  ('SV20230009', 'Trương Thị Ngọc',  '0901111119'),
  ('SV20230010', 'Hoàng Văn Long',   '0901111120'),
  ('SV20230011', 'Ngô Thị Mai',      '0901111121'),
  ('SV20230012', 'Lý Văn Phúc',      '0901111122'),
  ('SV20230013', 'Trịnh Thị Hà',     '0901111123'),
  ('SV20230014', 'Phan Văn Huy',     '0901111124'),
  ('SV20230015', 'Hồ Thị Yến',       '0901111125'),
  ('SV20230016', 'Dương Văn Khánh',  '0901111126'),
  ('SV20230017', 'Mai Thị Lan',      '0901111127'),
  ('SV20230018', 'Đinh Văn Sơn',     '0901111128'),
  ('SV20230019', 'Lâm Thị Thảo',     '0901111129'),
  ('SV20230020', 'Châu Văn Đức',     '0901111130'),
  ('SV20230021', 'Nguyễn Thị Kim',   '0901111131'),
  ('SV20230022', 'Trần Văn Tài',     '0901111132'),
  ('SV20230023', 'Phạm Thị Nhung',   '0901111133'),
  ('SV20230024', 'Lê Văn Phước',     '0901111134'),
  ('SV20230025', 'Vũ Thị Hoa',       '0901111135'),
  ('SV20230026', 'Huỳnh Văn Đông',   '0901111136'),
  ('SV20230027', 'Cao Thị Hạnh',     '0901111137'),
  ('SV20230028', 'Tạ Văn Nghĩa',     '0901111138'),
  ('SV20230029', 'Lương Thị Phượng', '0901111139'),
  ('SV20230030', 'Hà Văn Tiến',      '0901111140')
ON CONFLICT (student_id) DO NOTHING;

-- =============================================================
-- 2. USERS (students + 3 guards + 1 faculty)
-- bcrypt hash of "demo123" (same as existing seed)
-- =============================================================
INSERT INTO users (email, password_hash, role, status, member_id)
SELECT 'student3@university.edu.vn', '$2a$10$MzJZqWhKv7/lGKYnb1iV.ONhnv0.UyUj78cxeQlJ2uMEx34LPOng.', 'student', 'active', id
FROM members WHERE student_id = 'SV20230001' AND NOT EXISTS (SELECT 1 FROM users WHERE email = 'student3@university.edu.vn');

-- Batch create students + faculty + guard accounts
DO $$
DECLARE
  s RECORD;
  idx INT := 3;
  v_email TEXT;
BEGIN
  FOR s IN SELECT student_id, full_name FROM members WHERE student_id LIKE 'SV2023%' ORDER BY student_id
  LOOP
    idx := idx + 1;
    v_email := 'student' || idx || '@university.edu.vn';
    IF NOT EXISTS (SELECT 1 FROM users u WHERE u.email = v_email) THEN
      INSERT INTO users (email, password_hash, role, status, member_id)
      VALUES (v_email, '$2a$10$MzJZqWhKv7/lGKYnb1iV.ONhnv0.UyUj78cxeQlJ2uMEx34LPOng.', 'student', 'active', (SELECT id FROM members WHERE student_id = s.student_id));
    END IF;
  END LOOP;
END $$;

-- Faculty account
INSERT INTO users (email, password_hash, role, status)
VALUES ('faculty@university.edu.vn', '$2a$10$MzJZqWhKv7/lGKYnb1iV.ONhnv0.UyUj78cxeQlJ2uMEx34LPOng.', 'faculty', 'active')
ON CONFLICT (email) DO NOTHING;

-- Extra guards
INSERT INTO users (email, password_hash, role, status) VALUES
  ('guard2@parksmart.vn', '$2a$10$MzJZqWhKv7/lGKYnb1iV.ONhnv0.UyUj78cxeQlJ2uMEx34LPOng.', 'staff', 'active'),
  ('guard3@parksmart.vn', '$2a$10$MzJZqWhKv7/lGKYnb1iV.ONhnv0.UyUj78cxeQlJ2uMEx34LPOng.', 'staff', 'active'),
  ('guard4@parksmart.vn', '$2a$10$MzJZqWhKv7/lGKYnb1iV.ONhnv0.UyUj78cxeQlJ2uMEx34LPOng.', 'staff', 'active')
ON CONFLICT (email) DO NOTHING;

-- =============================================================
-- 3. REGISTERED VEHICLES (2 per student member)
-- =============================================================
INSERT INTO registered_vehicles (plate_number, member_id, description)
SELECT '59F1-' || LPAD(seq::TEXT, 5, '0'), m.id, 'Xe máy ' || m.full_name
FROM members m CROSS JOIN (SELECT generate_series(1, 2) AS seq)
WHERE m.student_id LIKE 'SV2023%'
ON CONFLICT DO NOTHING;

-- =============================================================
-- 4. VEHICLES (newer table, 1-2 per user)
-- =============================================================
INSERT INTO vehicles (license_plate, brand, model, color, owner_id, is_active)
SELECT
  '59A-' || LPAD(u.id::TEXT, 5, '0'),
  CASE (row_number() OVER ()) % 5
    WHEN 0 THEN 'Honda' WHEN 1 THEN 'Yamaha' WHEN 2 THEN 'Piaggio' WHEN 3 THEN 'SYM' ELSE 'Suzuki'
  END,
  CASE (row_number() OVER ()) % 4
    WHEN 0 THEN 'Vision' WHEN 1 THEN 'Air Blade' WHEN 2 THEN 'Lead' ELSE 'Janus'
  END,
  CASE (row_number() OVER ()) % 6
    WHEN 0 THEN 'Đen' WHEN 1 THEN 'Trắng' WHEN 2 THEN 'Đỏ' WHEN 3 THEN 'Xanh' WHEN 4 THEN 'Bạc' ELSE 'Xám'
  END,
  u.id, TRUE
FROM users u WHERE u.role IN ('student', 'staff') AND u.email LIKE '%university.edu.vn'
ON CONFLICT (license_plate) DO NOTHING;

-- =============================================================
-- 5. CARDS (1-2 per member + casual cards)
-- =============================================================
INSERT INTO cards (card_uid, card_type, member_id, status, balance)
SELECT
  'NFC-MTH-' || LPAD(seq::TEXT, 4, '0'),
  'monthly',
  m.id,
  CASE (row_number() OVER ()) % 10
    WHEN 0 THEN 'blocked' WHEN 1 THEN 'lost' ELSE 'active'
  END,
  (random() * 500000)::DECIMAL(10,2)
FROM members m CROSS JOIN (SELECT generate_series(1, 2) AS seq)
WHERE m.student_id LIKE 'SV2023%'
ON CONFLICT (card_uid) DO NOTHING;

-- Casual cards
INSERT INTO cards (card_uid, card_type, status, balance)
SELECT
  'NFC-CSL-' || LPAD(seq::TEXT, 4, '0'),
  'casual',
  'active',
  0
FROM generate_series(1, 20) AS seq
ON CONFLICT (card_uid) DO NOTHING;

-- =============================================================
-- 6. PARKING LOTS (2 lots)
-- =============================================================
INSERT INTO parking_lots (name, address, type, total_capacity, current_occupancy, open_time, close_time, manager_name, contact_phone)
VALUES
  ('Bãi xe A - Khu Giảng đường', 'Khu A, Đại học UIT, KP6, Linh Trung, Thủ Đức', 'outdoor', 500, 0, '06:00', '22:00', 'Nguyễn Văn An', '02837252001'),
  ('Bãi xe B - Ký túc xá', 'Khu B, Đại học UIT, KP6, Linh Trung, Thủ Đức', 'multi_level', 800, 0, '05:30', '23:00', 'Trần Thị Bình', '02837252002')
ON CONFLICT (name) DO NOTHING;

-- =============================================================
-- 7. PARKING SESSIONS (100+ sessions — mix ongoing/completed, past 45 days)
-- =============================================================
DO $$
DECLARE
  card_rec RECORD;
  session_count INT := 0;
  card_status TEXT;
  r_plate VARCHAR(20);
  check_in TIMESTAMP;
  check_out TIMESTAMP;
  duration INT;
  lot_id UUID;
BEGIN
  SELECT id INTO lot_id FROM parking_lots WHERE name = 'Bãi xe A - Khu Giảng đường' LIMIT 1;
  IF lot_id IS NULL THEN
    lot_id := gen_random_uuid();
  END IF;

  FOR card_rec IN SELECT card_uid, card_type, status FROM cards WHERE card_type = 'monthly' AND status = 'active' LIMIT 30
  LOOP
    -- 3-5 completed sessions per card over the past 45 days
    FOR i IN 1..(3 + (random() * 2)::INT)
    LOOP
      duration := (30 + (random() * 480)::INT); -- 30 min to 8 hours
      check_in := NOW() - ((random() * 45)::INT || ' days')::INTERVAL - ((random() * 12)::INT || ' hours')::INTERVAL;
      check_out := check_in + (duration || ' minutes')::INTERVAL;
      r_plate := '59F1-' || LPAD((10000 + (random() * 50000)::INT)::TEXT, 5, '0');

      INSERT INTO parking_sessions (card_uid, plate_in, check_in_time, plate_out, check_out_time, status)
      VALUES (card_rec.card_uid, r_plate, check_in, r_plate, check_out, 'completed');
      session_count := session_count + 1;
    END LOOP;

    -- 0-1 ongoing session per card
    IF random() < 0.3 THEN
      r_plate := '59F1-' || LPAD((10000 + (random() * 50000)::INT)::TEXT, 5, '0');
      INSERT INTO parking_sessions (card_uid, plate_in, check_in_time, status)
      VALUES (card_rec.card_uid, r_plate, NOW() - ((random() * 4)::INT || ' hours')::INTERVAL, 'ongoing');
      session_count := session_count + 1;
    END IF;
  END LOOP;

  -- Sessions from casual cards (fewer)
  FOR card_rec IN SELECT card_uid FROM cards WHERE card_type = 'casual' AND status = 'active' LIMIT 10
  LOOP
    FOR i IN 1..(1 + (random() * 3)::INT)
    LOOP
      duration := (10 + (random() * 60)::INT);
      check_in := NOW() - ((random() * 30)::INT || ' days')::INTERVAL - ((random() * 8)::INT || ' hours')::INTERVAL;
      check_out := check_in + (duration || ' minutes')::INTERVAL;
      r_plate := '51H-' || LPAD((10000 + (random() * 50000)::INT)::TEXT, 5, '0');

      INSERT INTO parking_sessions (card_uid, plate_in, check_in_time, plate_out, check_out_time, status)
      VALUES (card_rec.card_uid, r_plate, check_in, r_plate, check_out, 'completed');
      session_count := session_count + 1;
    END LOOP;
  END LOOP;

  RAISE NOTICE 'Created % parking sessions', session_count;
END $$;

-- =============================================================
-- 8. TRANSACTIONS (for completed sessions)
-- =============================================================
INSERT INTO transactions (card_uid, amount, type, session_id)
SELECT
  s.card_uid,
  (random() * 15000 + 2000)::DECIMAL(10,2),
  'parking_fee',
  s.id
FROM parking_sessions s
WHERE s.status = 'completed'
  AND NOT EXISTS (SELECT 1 FROM transactions t WHERE t.session_id = s.id)
LIMIT 80;

-- =============================================================
-- 9. SHIFTS (next 14 days, 3 shifts per day = 42 shifts)
-- =============================================================
DO $$
DECLARE
  d DATE;
  shift_id UUID;
  shift_count INT := 0;
  guard_user RECORD;
BEGIN
  d := CURRENT_DATE;
  FOR day_offset IN 0..13 LOOP
    -- Morning shift 06:00-12:00
    INSERT INTO shifts (name, type, start_time, end_time, date, status, notes)
    VALUES ('Ca sáng ' || (d + day_offset)::TEXT, 'morning', '06:00', '12:00', d + day_offset, 'scheduled',
            'Ca sáng - Bãi xe A')
    RETURNING id INTO shift_id;

    -- Afternoon shift 12:00-18:00
    INSERT INTO shifts (name, type, start_time, end_time, date, status, notes)
    VALUES ('Ca chiều ' || (d + day_offset)::TEXT, 'afternoon', '12:00', '18:00', d + day_offset, 'scheduled',
            'Ca chiều - Bãi xe A');

    -- Night shift 18:00-00:00
    INSERT INTO shifts (name, type, start_time, end_time, date, status, notes)
    VALUES ('Ca tối ' || (d + day_offset)::TEXT, 'night', '18:00', '00:00', d + day_offset, 'scheduled',
            'Ca tối - Bãi xe B');

    shift_count := shift_count + 3;

    -- Mark past shifts as completed
    IF (d + day_offset) < CURRENT_DATE THEN
      UPDATE shifts SET status = 'completed' WHERE date = (d + day_offset);
    END IF;
  END LOOP;
  RAISE NOTICE 'Created % shifts', shift_count;
END $$;

-- =============================================================
-- 10. SHIFT ASSIGNMENTS (assign guards to shifts)
-- =============================================================
DO $$
DECLARE
  s RECORD;
  guard_ids UUID[];
  g_id UUID;
  assigned INT := 0;
BEGIN
  SELECT ARRAY_AGG(id) INTO guard_ids FROM users WHERE role = 'staff' AND status = 'active';
  IF guard_ids IS NOT NULL THEN
    FOR s IN SELECT id FROM shifts WHERE status IN ('scheduled', 'completed') ORDER BY date
    LOOP
      -- Pick a random guard
      g_id := guard_ids[1 + (random() * (array_length(guard_ids, 1) - 1))::INT];
      INSERT INTO shift_assignments (shift_id, user_id)
      VALUES (s.id, g_id)
      ON CONFLICT (shift_id, user_id) DO NOTHING;
      assigned := assigned + 1;
    END LOOP;
  END IF;
  RAISE NOTICE 'Created % shift assignments', assigned;
END $$;

-- =============================================================
-- 11. CARD REQUESTS (some pending, some approved)
-- =============================================================
INSERT INTO card_requests (member_id, vehicle_plate, vehicle_brand, vehicle_model, vehicle_color, id_card_number, note, status, card_uid)
SELECT
  m.id,
  '59F1-' || LPAD((10000 + (random() * 50000)::INT)::TEXT, 5, '0'),
  CASE (row_number() OVER ()) % 4 WHEN 0 THEN 'Honda' WHEN 1 THEN 'Yamaha' WHEN 2 THEN 'Piaggio' ELSE 'SYM' END,
  CASE (row_number() OVER ()) % 3 WHEN 0 THEN 'Vision' WHEN 1 THEN 'Air Blade' ELSE 'Future' END,
  CASE (row_number() OVER ()) % 5 WHEN 0 THEN 'Đen' WHEN 1 THEN 'Trắng' WHEN 2 THEN 'Đỏ' ELSE 'Xanh' END,
  '079' || LPAD((100000000 + (random() * 899999999)::BIGINT)::TEXT, 9, '0'),
  'Yêu cầu cấp thẻ sinh viên',
  CASE (row_number() OVER ()) % 3 WHEN 0 THEN 'approved' WHEN 1 THEN 'pending' ELSE 'blocked' END,
  (SELECT card_uid FROM cards WHERE card_type = 'monthly' OFFSET (random() * (SELECT count(*) FROM cards WHERE card_type = 'monthly'))::INT LIMIT 1)
FROM members m
WHERE m.student_id LIKE 'SV2023%'
LIMIT 15;

-- =============================================================
-- 12. INCIDENTS
-- =============================================================
INSERT INTO incidents (reported_by, vehicle_plate, type, description, location, status)
SELECT
  (SELECT id FROM users WHERE email = 'guard@parksmart.vn' LIMIT 1),
  '59F1-' || LPAD((10000 + (random() * 50000)::INT)::TEXT, 5, '0'),
  CASE (row_number() OVER ()) % 5
    WHEN 0 THEN 'wrong_parking' WHEN 1 THEN 'damaged_vehicle' WHEN 2 THEN 'suspicious' WHEN 3 THEN 'unregistered' ELSE 'other'
  END,
  CASE (row_number() OVER ()) % 4
    WHEN 0 THEN 'Xe đậu sai vị trí quy định'
    WHEN 1 THEN 'Xe bị trầy xước phần đuôi'
    WHEN 2 THEN 'Người lạ vào khu vực để xe sau giờ đóng cửa'
    ELSE 'Xe không có thẻ ra vào'
  END,
  CASE (row_number() OVER ()) % 3
    WHEN 0 THEN 'Bãi xe A - Tầng 1' WHEN 1 THEN 'Bãi xe B - Lối vào' ELSE 'Bãi xe A - Tầng 2'
  END,
  CASE (row_number() OVER ()) % 3 WHEN 0 THEN 'open' WHEN 1 THEN 'resolved' ELSE 'escalated' END
FROM generate_series(1, 12);

-- =============================================================
-- 13. VISITOR PASSES
-- =============================================================
INSERT INTO visitor_passes (user_id, visitor_name, visitor_phone, vehicle_plate, valid_date, status, qr_code_data)
SELECT
  u.id,
  CASE (row_number() OVER ()) % 4
    WHEN 0 THEN 'Nguyễn Thị Khách' WHEN 1 THEN 'Trần Văn Thăm' WHEN 2 THEN 'Lê Thị Viếng' ELSE 'Phạm Văn Bạn'
  END,
  '090' || LPAD((1000000 + (random() * 8999999)::INT)::TEXT, 7, '0'),
  '59F1-' || LPAD((10000 + (random() * 50000)::INT)::TEXT, 5, '0'),
  CURRENT_DATE + (row_number() OVER ())::INT,
  'valid',
  'QR-' || md5(random()::text)
FROM users u WHERE u.role = 'student' AND u.email LIKE '%university.edu.vn'
LIMIT 20;

-- =============================================================
-- 14. DEVICES
-- =============================================================
INSERT INTO devices (name, type, location_label, status, ip_address, firmware_version)
VALUES
  ('Camera cổng A - 01',  'camera',     'Cổng vào Bãi A',       'online',      '192.168.1.10', 'v2.3.1'),
  ('Camera cổng A - 02',  'camera',     'Cổng ra Bãi A',        'online',      '192.168.1.11', 'v2.3.1'),
  ('Camera cổng B - 01',  'camera',     'Cổng vào Bãi B',       'online',      '192.168.1.12', 'v2.3.1'),
  ('Camera cổng B - 02',  'camera',     'Cổng ra Bãi B',        'warning',     '192.168.1.13', 'v2.3.0'),
  ('Camera trong Bãi A',  'camera',     'Tầng 1 Bãi A',         'online',      '192.168.1.14', 'v2.3.1'),
  ('Camera trong Bãi B',  'camera',     'Tầng 2 Bãi B',         'maintenance', '192.168.1.15', 'v2.2.9'),
  ('Barie cổng A vào',    'barrier',    'Cổng vào Bãi A',       'online',      '192.168.1.20', 'v1.9.2'),
  ('Barie cổng A ra',     'barrier',    'Cổng ra Bãi A',        'online',      '192.168.1.21', 'v1.9.2'),
  ('Barie cổng B vào',    'barrier',    'Cổng vào Bãi B',       'offline',     '192.168.1.22', 'v1.9.1'),
  ('Barie cổng B ra',     'barrier',    'Cổng ra Bãi B',        'online',      '192.168.1.23', 'v1.9.2'),
  ('RFID Đầu đọc A1',     'rfid_reader','Cổng vào Bãi A - Làn 1','online',     '192.168.1.30', 'v3.0.1'),
  ('RFID Đầu đọc A2',     'rfid_reader','Cổng vào Bãi A - Làn 2','online',     '192.168.1.31', 'v3.0.1'),
  ('RFID Đầu đọc B1',     'rfid_reader','Cổng vào Bãi B',       'online',      '192.168.1.32', 'v3.0.0'),
  ('RFID Đầu đọc B2',     'rfid_reader','Cổng ra Bãi B',        'warning',     '192.168.1.33', 'v3.0.0'),
  ('Cảm biến nhiệt A',    'sensor',     'Bãi A - Tầng 1',       'online',      '192.168.1.40', 'v1.0.3')
ON CONFLICT DO NOTHING;

-- =============================================================
-- 15. DEVICE ALERTS
-- =============================================================
INSERT INTO device_alerts (device_id, message, severity)
SELECT d.id, CASE d.type
  WHEN 'camera' THEN 'Mất kết nối camera - kiểm tra cáp mạng'
  WHEN 'barrier' THEN 'Barie không đóng/mở được - cần kiểm tra motor'
  WHEN 'rfid_reader' THEN 'Đầu đọc RFID chậm phản hồi'
  ELSE 'Cảm biến gửi tín hiệu bất thường'
END, 'medium'
FROM devices d WHERE d.status IN ('warning', 'offline', 'maintenance');

-- =============================================================
-- 16. NOTIFICATIONS (for each student user)
-- =============================================================
INSERT INTO notifications (user_id, title, message, type)
SELECT
  u.id,
  CASE (row_number() OVER ()) % 3
    WHEN 0 THEN 'Nhắc nhở: Thẻ xe sắp hết hạn'
    WHEN 1 THEN 'Thông báo: Lịch bảo trì bãi xe'
    ELSE 'Cập nhật: Gói tháng mới'
  END,
  CASE (row_number() OVER ()) % 3
    WHEN 0 THEN 'Thẻ xe của bạn sẽ hết hạn trong 7 ngày tới. Vui lòng gia hạn để tiếp tục sử dụng.'
    WHEN 1 THEN 'Bãi xe B sẽ bảo trì vào Chủ Nhật tuần này từ 8h-17h. Vui lòng gửi xe tại bãi A.'
    ELSE 'Gói tháng mới với nhiều ưu đãi đã có mặt. Đăng ký ngay để nhận ưu đãi!'
  END,
  CASE (row_number() OVER ()) % 4 WHEN 0 THEN 'info' WHEN 1 THEN 'warning' WHEN 2 THEN 'success' ELSE 'error' END
FROM users u WHERE u.role = 'student' AND u.status = 'active'
LIMIT 40;

-- =============================================================
-- 17. SUPPORT TICKETS + RESPONSES
-- =============================================================
INSERT INTO support_tickets (user_id, category, subject, description, status)
SELECT
  u.id,
  CASE (row_number() OVER ()) % 4
    WHEN 0 THEN 'wallet_issue' WHEN 1 THEN 'card_issue' WHEN 2 THEN 'staff_attitude' ELSE 'other'
  END,
  CASE (row_number() OVER ()) % 4
    WHEN 0 THEN 'Nạp tiền không thành công'
    WHEN 1 THEN 'Thẻ bị mất / hỏng'
    WHEN 2 THEN 'Nhân viên thiếu chuyên nghiệp'
    ELSE 'Góp ý cải thiện dịch vụ'
  END,
  CASE (row_number() OVER ()) % 4
    WHEN 0 THEN 'Tôi đã nạp 50.000đ qua Momo nhưng số dư không tăng.'
    WHEN 1 THEN 'Thẻ của tôi bị mất vào ngày 15/6. Vui lòng khóa thẻ và cấp lại.'
    WHEN 2 THEN 'Nhân viên bảo vệ ca tối ngày 10/6 có thái độ không tốt với sinh viên.'
    ELSE 'Nên bổ sung thêm máy bán nước tự động ở bãi xe.'
  END,
  CASE (row_number() OVER ()) % 3 WHEN 0 THEN 'open' WHEN 1 THEN 'resolved' ELSE 'closed' END
FROM users u WHERE u.role = 'student' AND u.status = 'active'
LIMIT 15;

-- Responses for resolved/closed tickets
INSERT INTO ticket_responses (ticket_id, sender_id, message, is_admin)
SELECT
  t.id,
  (SELECT id FROM users WHERE email = 'admin@parksmart.vn' LIMIT 1),
  'Cảm ơn bạn đã phản ánh. Chúng tôi đã xử lý vấn đề này. Vui lòng kiểm tra lại hoặc liên hệ hotline nếu cần hỗ trợ thêm.',
  TRUE
FROM support_tickets t
WHERE t.status IN ('resolved', 'closed');

-- =============================================================
-- 18. UPDATE parking lot occupancy (rough count of ongoing sessions)
-- =============================================================
UPDATE parking_lots SET current_occupancy = (SELECT count(*) FROM parking_sessions WHERE status = 'ongoing');

COMMIT;

-- Summary
SELECT 'SEED COMPLETE' AS result;
SELECT 'Members' AS entity, count(*) AS count FROM members
UNION ALL SELECT 'Users', count(*) FROM users
UNION ALL SELECT 'Cards', count(*) FROM cards
UNION ALL SELECT 'Registered Vehicles', count(*) FROM registered_vehicles
UNION ALL SELECT 'Vehicles', count(*) FROM vehicles
UNION ALL SELECT 'Parking Sessions', count(*) FROM parking_sessions
UNION ALL SELECT 'Transactions', count(*) FROM transactions
UNION ALL SELECT 'Shifts', count(*) FROM shifts
UNION ALL SELECT 'Shift Assignments', count(*) FROM shift_assignments
UNION ALL SELECT 'Card Requests', count(*) FROM card_requests
UNION ALL SELECT 'Parking Lots', count(*) FROM parking_lots
UNION ALL SELECT 'Devices', count(*) FROM devices
UNION ALL SELECT 'Device Alerts', count(*) FROM device_alerts
UNION ALL SELECT 'Notifications', count(*) FROM notifications
UNION ALL SELECT 'Support Tickets', count(*) FROM support_tickets
UNION ALL SELECT 'Ticket Responses', count(*) FROM ticket_responses
UNION ALL SELECT 'Incidents', count(*) FROM incidents
UNION ALL SELECT 'Visitor Passes', count(*) FROM visitor_passes
ORDER BY entity;
