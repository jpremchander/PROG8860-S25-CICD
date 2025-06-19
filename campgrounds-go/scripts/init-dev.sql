-- Development database initialization
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create sample users
INSERT INTO users (username, email, password, created_at, updated_at) VALUES
('devuser1', 'dev1@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', NOW(), NOW()),
('devuser2', 'dev2@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', NOW(), NOW())
ON CONFLICT (email) DO NOTHING;

-- Create sample campgrounds
INSERT INTO campgrounds (title, description, location, price, author_id, created_at, updated_at) VALUES
('Dev Campground 1', 'A beautiful development campground for testing purposes', 'Test Valley', 25.99, 1, NOW(), NOW()),
('Dev Campground 2', 'Another test campground with amazing views', 'Debug Mountains', 35.50, 2, NOW(), NOW())
ON CONFLICT DO NOTHING;
