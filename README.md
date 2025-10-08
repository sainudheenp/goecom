# Simple E-Commerce Backend - Go

A basic Go backend for an e-commerce application with essential features.

## 🚀 Basic Features

- **User Registration & Login**: Simple JWT authentication
- **Product Catalog**: View and search products
- **Shopping Cart**: Add/remove items from cart
- **Basic Orders**: Create orders from cart
- **Database**: PostgreSQL with basic migrations

## 📋 Prerequisites

- **Go**: 1.21 or higher
- **PostgreSQL**: 15 or higher
- **Docker & Docker Compose** (optional)

## 🏗️ Simple Architecture

```
.
├── cmd/server/          # Application entry point
├── handlers/            # HTTP handlers
├── models/              # Database models
├── database/            # Database connection
├── auth/                # Simple authentication
├── migrations/          # SQL migrations
└── main.go              # Main application
```

## 🚀 Quick Start

### 1. Clone and Setup

```bash
git clone https://github.com/sainudheenp/goecom.git
cd goecom

# Copy environment file
cp .env.example .env

# Edit .env and set JWT_SECRET to a strong random string (min 32 chars)
```

### 2. Start Services

```bash
# Build and start all services (app + postgres + swagger-ui)
docker-compose up --build

# The app will be available at:
# - API: http://localhost:8080
# - Swagger UI: http://localhost:8081
```

### 3. Seed Database

In a new terminal:

```bash
# Wait for services to be healthy, then seed data
docker-compose exec app sh -c "apt-get update && apt-get install -y postgresql-client"
docker-compose exec app bash /root/scripts/seed.sh
```

Or manually seed using:

```bash
# Connect to the database
docker-compose exec db psql -U postgres -d ecom

# Run seed script (see scripts/seed.sh for SQL)
```

### 4. Test the API

```bash
# Health check
curl http://localhost:8080/health

# Register a user
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "full_name": "Test User"
  }'

# Login (use seeded admin)
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "admin123"
  }'

# Save the access_token from the response
export TOKEN="<access_token>"

# Get current user
curl http://localhost:8080/api/v1/me \
  -H "Authorization: Bearer $TOKEN"

# List products
curl http://localhost:8080/api/v1/products

# Get specific product
curl http://localhost:8080/api/v1/products/<product_id>

# Add to cart
curl -X POST http://localhost:8080/api/v1/cart \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": "<product_id>",
    "quantity": 2
  }'

# View cart
curl http://localhost:8080/api/v1/cart \
  -H "Authorization: Bearer $TOKEN"

# Create order
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "shipping_address": {
      "line1": "123 Main St",
      "city": "Bengaluru",
      "state": "KA",
      "country": "IN",
      "postcode": "560001"
    }
  }'

# List orders
curl http://localhost:8080/api/v1/orders \
  -H "Authorization: Bearer $TOKEN"

# Process payment
curl -X POST http://localhost:8080/api/v1/payments/charge \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "order_id": "<order_id>",
    "payment_method": "card",
    "payment_details": {}
  }'
```

## 🛠️ Local Development (Without Docker)

### 1. Install Dependencies

```bash
# Install Go dependencies
go mod download

# Install PostgreSQL locally (or use Docker for just the DB)
docker run -d \
  --name ecom-postgres \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=ecom \
  -p 5432:5432 \
  postgres:15-alpine
```

### 2. Configure Environment

```bash
cp .env.example .env

# Edit .env
DATABASE_URL=postgres://postgres:postgres@localhost:5432/ecom?sslmode=disable
JWT_SECRET=your_very_strong_secret_key_minimum_32_characters_long
```

### 3. Run Migrations

```bash
# On Windows (Git Bash or WSL recommended for shell scripts)
bash scripts/migrate.sh up

# Or use psql directly
psql $DATABASE_URL -f migrations/001_create_users_table.up.sql
psql $DATABASE_URL -f migrations/002_create_products_table.up.sql
psql $DATABASE_URL -f migrations/003_create_cart_items_table.up.sql
psql $DATABASE_URL -f migrations/004_create_orders_table.up.sql
psql $DATABASE_URL -f migrations/005_create_order_items_table.up.sql
```

### 4. Seed Database

```bash
bash scripts/seed.sh

# Or manually with the SQL from scripts/seed.sh
```

### 5. Run the Server

```bash
# Using go run
go run ./cmd/server

# Or build and run
go build -o bin/server.exe ./cmd/server
./bin/server.exe

# Or using Make (requires Make on Windows: choco install make)
make run
```

Server will start on `http://localhost:8080`

## 🧪 Testing

### Run Unit Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -v -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Or using Make
make test
make test-coverage
```

### Run Integration Tests

```bash
# Ensure database is running and seeded
export DATABASE_URL=postgres://postgres:postgres@localhost:5432/ecom?sslmode=disable
export JWT_SECRET=test_jwt_secret_key_for_testing_purposes_minimum_32_chars

# Run integration tests
go test -v -tags=integration ./test/...

# Or using Make
make integration-test
```

## 📊 Code Quality

### Linting

```bash
# Install golangci-lint (one-time)
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linter
golangci-lint run

# Or using Make
make lint
```

### Formatting

```bash
# Format code
go fmt ./...

# Or using Make
make fmt
```

## 🔑 Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `PORT` | Server port | `8080` | No |
| `ENV` | Environment (development/production) | `development` | No |
| `DATABASE_URL` | PostgreSQL connection string | - | **Yes** |
| `JWT_SECRET` | Secret for JWT signing (min 32 chars) | - | **Yes** |
| `JWT_EXPIRES_HOURS` | JWT expiration time in hours | `24` | No |
| `BCRYPT_COST` | Bcrypt hashing cost | `10` | No |
| `LOG_LEVEL` | Logging level | `info` | No |
| `CORS_ORIGINS` | Allowed CORS origins (comma-separated) | `*` | No |
| `RATE_LIMIT_REQUESTS` | Max requests per window | `100` | No |
| `RATE_LIMIT_WINDOW_MINUTES` | Rate limit window | `15` | No |

## 📖 API Documentation

### View Swagger UI

```bash
# Start swagger-ui with docker-compose
docker-compose up swagger-ui

# Or run standalone
docker run -p 8081:8080 \
  -e SWAGGER_JSON=/openapi.yaml \
  -v $(pwd)/openapi.yaml:/openapi.yaml \
  swaggerapi/swagger-ui

# Access at http://localhost:8081
```

### API Endpoints Summary

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | `/api/v1/auth/register` | Public | Register new user |
| POST | `/api/v1/auth/login` | Public | Login user |
| GET | `/api/v1/me` | User | Get current user |
| GET | `/api/v1/products` | Public | List products (with filters) |
| GET | `/api/v1/products/:id` | Public | Get product by ID |
| POST | `/api/v1/products` | Admin | Create product |
| PUT | `/api/v1/products/:id` | Admin | Update product |
| DELETE | `/api/v1/products/:id` | Admin | Delete product |
| POST | `/api/v1/cart` | User | Add to cart |
| GET | `/api/v1/cart` | User | Get cart |
| DELETE | `/api/v1/cart/:item_id` | User | Remove from cart |
| POST | `/api/v1/orders` | User | Create order |
| GET | `/api/v1/orders` | User | List user orders |
| GET | `/api/v1/orders/:id` | User | Get order by ID |
| POST | `/api/v1/payments/charge` | User | Process payment |
| GET | `/api/v1/admin/orders` | Admin | List all orders |
| PATCH | `/api/v1/admin/orders/:id` | Admin | Update order status |

## 🔒 Security Features

- **JWT Authentication**: Secure token-based auth with configurable expiration
- **Password Hashing**: bcrypt with configurable cost factor
- **Role-Based Access Control**: User and admin roles
- **Input Validation**: Request validation using Gin binding
- **SQL Injection Prevention**: Parameterized queries via GORM
- **CORS**: Configurable cross-origin resource sharing
- **Rate Limiting**: Token bucket algorithm per IP
- **Request Correlation**: X-Request-ID for request tracing

## 🗄️ Database Schema

See `migrations/*.sql` files for complete schema.

Key tables:
- **users**: User accounts with email, password hash, role
- **products**: Product catalog with SKU, pricing, stock
- **cart_items**: Shopping cart items per user
- **orders**: Customer orders with status and shipping
- **order_items**: Line items for each order

All tables use UUID primary keys and include `created_at`/`updated_at` timestamps in UTC.

## 🐳 Docker Commands

```bash
# Build image
docker build -t goecom:latest .

# Run container
docker run -p 8080:8080 \
  -e DATABASE_URL=<db_url> \
  -e JWT_SECRET=<secret> \
  goecom:latest

# Using docker-compose
docker-compose up -d          # Start in background
docker-compose logs -f app    # View logs
docker-compose down           # Stop services
docker-compose down -v        # Stop and remove volumes
```

## 📝 Make Commands

```bash
make help              # Show all commands
make run               # Run the server
make build             # Build binary
make test              # Run unit tests
make test-coverage     # Run tests with coverage
make lint              # Run linter
make fmt               # Format code
make migrate-up        # Run migrations
make migrate-down      # Rollback migrations
make seed              # Seed database
make docker-up         # Start docker-compose
make docker-down       # Stop docker-compose
make clean             # Clean build artifacts
```

## 🚢 Deployment

### Build Production Image

```bash
docker build -t goecom:v1.0.0 .
docker tag goecom:v1.0.0 your-registry.com/goecom:v1.0.0
docker push your-registry.com/goecom:v1.0.0
```

### Environment Setup

Ensure production environment has:
- PostgreSQL database with migrations applied
- Secure `JWT_SECRET` (32+ characters, randomly generated)
- `BCRYPT_COST` set to 12 or higher for production
- `ENV=production` to run Gin in release mode
- Proper `CORS_ORIGINS` configured
- Rate limiting tuned for expected traffic

### CI/CD

GitHub Actions workflow (`.github/workflows/ci.yml`) runs on every push:
1. Lints code with `golangci-lint`
2. Runs unit tests
3. Runs integration tests
4. Builds Docker image
5. (Optional) Push to registry and deploy

## 🤝 Sample Credentials (Seeded)

After running the seed script:

- **Admin**: `admin@example.com` / `admin123`
- **User**: `user@example.com` / `admin123`

**⚠️ Change these credentials in production!**

## 📚 Additional Resources

- [Go Documentation](https://golang.org/doc/)
- [Gin Framework](https://gin-gonic.com/docs/)
- [GORM Documentation](https://gorm.io/docs/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [OpenAPI Specification](https://swagger.io/specification/)

## 🐛 Troubleshooting

### Database connection issues

```bash
# Check if PostgreSQL is running
docker-compose ps

# Check logs
docker-compose logs db

# Test connection
psql $DATABASE_URL -c "SELECT version();"
```

### Port already in use

```bash
# On Windows, find process using port 8080
netstat -ano | findstr :8080

# Kill process
taskkill /PID <PID> /F

# Or change port in .env
PORT=8081
```

### Migrations failed

```bash
# Rollback and retry
bash scripts/migrate.sh down
bash scripts/migrate.sh up

# Or reset database
docker-compose down -v
docker-compose up -d db
# Wait for DB to be ready
bash scripts/migrate.sh up
bash scripts/seed.sh
```

## 📄 License

This project is licensed under the MIT License.

## 👨‍💻 Author

Created by GitHub Copilot

---

**Happy Coding! 🚀**