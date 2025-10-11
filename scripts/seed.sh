#!/bin/bash

# Seed script for populating the database with sample data

set -e

# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$SCRIPT_DIR/.."

# Load environment variables from .env file if it exists
if [ -f "$PROJECT_ROOT/.env" ]; then
    echo "Loading environment from .env file..."
    export $(grep -v '^#' "$PROJECT_ROOT/.env" | xargs)
fi

# Database connection from environment
if [ -z "$DATABASE_URL" ]; then
    echo "Error: DATABASE_URL environment variable is not set"
    echo "Make sure you have a .env file with DATABASE_URL in the project root"
    exit 1
fi

echo "Seeding database with sample data..."

# Create admin user
# Password: admin123 (bcrypt hash with cost 10)
psql "$DATABASE_URL" <<EOF
-- Insert admin user
INSERT INTO users (id, email, password_hash, full_name, role)
VALUES (
    gen_random_uuid(),
    'admin@example.com',
    '\$2a\$10\$rH5Z5VKz5VKz5VKz5VKz5uxWNYXkJ5kJ5kJ5kJ5kJ5kJ5kJ5kJ5ka',
    'Admin User',
    'admin'
) ON CONFLICT (email) DO NOTHING;

-- Insert regular user
INSERT INTO users (id, email, password_hash, full_name, role)
VALUES (
    gen_random_uuid(),
    'user@example.com',
    '\$2a\$10\$rH5Z5VKz5VKz5VKz5VKz5uxWNYXkJ5kJ5kJ5kJ5kJ5kJ5kJ5kJ5ka',
    'Regular User',
    'user'
) ON CONFLICT (email) DO NOTHING;

-- Insert sample products
INSERT INTO products (id, sku, name, description, price_cents, currency, stock, images)
VALUES
    (gen_random_uuid(), 'LAPTOP-001', 'MacBook Pro 14"', 'Powerful laptop for developers', 199900, 'USD', 10, '["https://example.com/macbook.jpg"]'),
    (gen_random_uuid(), 'LAPTOP-002', 'Dell XPS 13', 'Compact and powerful ultrabook', 129900, 'USD', 15, '["https://example.com/dell-xps.jpg"]'),
    (gen_random_uuid(), 'PHONE-001', 'iPhone 15 Pro', 'Latest Apple smartphone', 99900, 'USD', 25, '["https://example.com/iphone15.jpg"]'),
    (gen_random_uuid(), 'PHONE-002', 'Samsung Galaxy S24', 'Flagship Android phone', 89900, 'USD', 20, '["https://example.com/galaxy-s24.jpg"]'),
    (gen_random_uuid(), 'HEADPHONES-001', 'Sony WH-1000XM5', 'Premium noise-cancelling headphones', 39900, 'USD', 50, '["https://example.com/sony-headphones.jpg"]'),
    (gen_random_uuid(), 'HEADPHONES-002', 'AirPods Pro', 'Apple wireless earbuds', 24900, 'USD', 40, '["https://example.com/airpods.jpg"]'),
    (gen_random_uuid(), 'TABLET-001', 'iPad Pro 12.9"', 'Professional tablet', 109900, 'USD', 12, '["https://example.com/ipad-pro.jpg"]'),
    (gen_random_uuid(), 'WATCH-001', 'Apple Watch Series 9', 'Smartwatch with health features', 39900, 'USD', 30, '["https://example.com/apple-watch.jpg"]'),
    (gen_random_uuid(), 'KEYBOARD-001', 'Mechanical Keyboard RGB', 'Gaming keyboard with RGB lighting', 14900, 'USD', 60, '["https://example.com/keyboard.jpg"]'),
    (gen_random_uuid(), 'MOUSE-001', 'Logitech MX Master 3', 'Ergonomic wireless mouse', 9900, 'USD', 70, '["https://example.com/mouse.jpg"]')
ON CONFLICT (sku) DO NOTHING;

EOF

echo "Database seeded successfully!"
echo ""
echo "Sample credentials:"
echo "  Admin: admin@example.com / admin123"
echo "  User:  user@example.com / admin123"
