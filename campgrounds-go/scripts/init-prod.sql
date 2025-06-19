-- Production database initialization
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_campgrounds_author_id ON campgrounds(author_id);
CREATE INDEX IF NOT EXISTS idx_reviews_campground_id ON reviews(campground_id);
CREATE INDEX IF NOT EXISTS idx_reviews_author_id ON reviews(author_id);
CREATE INDEX IF NOT EXISTS idx_images_campground_id ON images(campground_id);

-- Create admin user (password should be changed immediately)
INSERT INTO users (username, email, password, created_at, updated_at) VALUES
('admin', 'admin@campgrounds.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', NOW(), NOW())
ON CONFLICT (email) DO NOTHING;
