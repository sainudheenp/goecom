@echo off
REM Seed script for Windows

IF "%DATABASE_URL%"=="" (
    echo Error: DATABASE_URL environment variable is not set
    exit /b 1
)

echo Seeding database with sample data...

psql "%DATABASE_URL%" -c "INSERT INTO users (email, password_hash, full_name, role) VALUES ('admin@example.com', '$2a$10$rH5Z5VKz5VKz5VKz5VKz5uxWNYXkJ5kJ5kJ5kJ5kJ5kJ5kJ5kJ5ka', 'Admin User', 'admin') ON CONFLICT (email) DO NOTHING;"

psql "%DATABASE_URL%" -c "INSERT INTO users (email, password_hash, full_name, role) VALUES ('user@example.com', '$2a$10$rH5Z5VKz5VKz5VKz5VKz5uxWNYXkJ5kJ5kJ5kJ5kJ5kJ5kJ5kJ5ka', 'Regular User', 'user') ON CONFLICT (email) DO NOTHING;"

psql "%DATABASE_URL%" -c "INSERT INTO products (sku, name, description, price_cents, currency, stock, images) VALUES ('LAPTOP-001', 'MacBook Pro 14\"', 'Powerful laptop for developers', 199900, 'USD', 10, '[\"https://example.com/macbook.jpg\"]') ON CONFLICT (sku) DO NOTHING;"

psql "%DATABASE_URL%" -c "INSERT INTO products (sku, name, description, price_cents, currency, stock, images) VALUES ('LAPTOP-002', 'Dell XPS 13', 'Compact and powerful ultrabook', 129900, 'USD', 15, '[\"https://example.com/dell-xps.jpg\"]') ON CONFLICT (sku) DO NOTHING;"

psql "%DATABASE_URL%" -c "INSERT INTO products (sku, name, description, price_cents, currency, stock, images) VALUES ('PHONE-001', 'iPhone 15 Pro', 'Latest Apple smartphone', 99900, 'USD', 25, '[\"https://example.com/iphone15.jpg\"]') ON CONFLICT (sku) DO NOTHING;"

echo Database seeded successfully!
echo.
echo Sample credentials:
echo   Admin: admin@example.com / admin123
echo   User:  user@example.com / admin123
