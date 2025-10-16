-- Create admin user
-- Password: admin123 (hashed with bcrypt)
INSERT INTO users (id, email, username, password_hash, display_name, role, created_at, updated_at)
VALUES (
  'a0000000-0000-0000-0000-000000000001',
  'admin@ecomate.com',
  'admin',
  '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', -- admin123
  'System Administrator',
  'admin',
  NOW(),
  NOW()
)
ON CONFLICT (email) DO UPDATE
SET role = 'admin',
    username = 'admin',
    display_name = 'System Administrator',
    password_hash = '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy';

-- Confirm creation
SELECT id, email, username, role, display_name FROM users WHERE role = 'admin';
