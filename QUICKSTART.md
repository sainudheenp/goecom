# 🚀 E-Commerce Backend - Complete Setup Guide

## Project Summary

This is a **production-ready Go backend** for an e-commerce application with:

✅ **Full REST API** - Catalog, Cart, Orders, Users, Payments, Admin  
✅ **JWT Authentication** - Secure token-based auth with role-based access  
✅ **PostgreSQL Database** - With GORM ORM and raw SQL migrations  
✅ **Docker Support** - Multi-stage Dockerfile and docker-compose  
✅ **Comprehensive Tests** - Unit tests and integration tests  
✅ **CI/CD Pipeline** - GitHub Actions workflow  
✅ **API Documentation** - OpenAPI/Swagger specification  
✅ **Security Features** - Input validation, rate limiting, CORS, password hashing  
✅ **Observability** - Structured logging with request correlation  

---

## 📦 Quick Commands Reference

### Docker Compose (Recommended for Quick Start)

```bash
# 1. Start all services (app + database + swagger-ui)
docker-compose up --build

# 2. Run migrations (one-time setup)
docker-compose exec db psql -U postgres -d ecom -f /docker-entrypoint-initdb.d/001_create_users_table.up.sql
# ... or mount migrations and run script

# 3. Seed database (optional)
docker-compose exec app bash scripts/seed.sh

# 4. View logs
docker-compose logs -f app

# 5. Stop services
docker-compose down
```

**Services will be available at:**
- API: http://localhost:8080
- Swagger UI: http://localhost:8081  
- PostgreSQL: localhost:5432

---

### Local Development (Without Docker)

#### Prerequisites
- Go 1.21+
- PostgreSQL 15+
- `psql` CLI tool

#### Setup Steps

```bash
# 1. Clone repository
git clone https://github.com/sainudheenp/goecom.git
cd goecom

# 2. Install dependencies
go mod download

# 3. Setup environment
cp .env.example .env
# Edit .env and set DATABASE_URL and JWT_SECRET

# 4. Start PostgreSQL (if not running)
docker run -d --name ecom-postgres \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=ecom \
  -p 5432:5432 \
  postgres:15-alpine

# 5. Run migrations
# On Windows (PowerShell/CMD):
scripts\migrate.bat up

# On Linux/Mac/Git Bash:
bash scripts/migrate.sh up

# 6. Seed database (optional)
# On Windows:
scripts\seed.bat

# On Linux/Mac:
bash scripts/seed.sh

# 7. Run the server
go run ./cmd/server

# Server starts on http://localhost:8080
```

---

## 🧪 Testing Commands

```bash
# Run all unit tests
go test ./...

# Run tests with coverage
go test -v -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run integration tests (requires running database)
go test -v -tags=integration ./test/...

# Run linter
golangci-lint run

# Format code
go fmt ./...
```

---

## 🔑 Environment Variables (Required)

Create a `.env` file from `.env.example`:

```env
# Required
DATABASE_URL=postgres://postgres:postgres@localhost:5432/ecom?sslmode=disable
JWT_SECRET=your_very_strong_secret_key_minimum_32_characters_long

# Optional (with defaults)
PORT=8080
ENV=development
JWT_EXPIRES_HOURS=24
BCRYPT_COST=10
LOG_LEVEL=info
CORS_ORIGINS=http://localhost:3000,http://localhost:8080
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW_MINUTES=15
```

**⚠️ Important:** Generate a strong `JWT_SECRET` for production!

```bash
# Generate random secret (32+ characters)
openssl rand -base64 32
```

---

## 📋 API Endpoints Quick Reference

### Public Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/auth/register` | Register new user |
| POST | `/api/v1/auth/login` | Login and get JWT token |
| GET | `/api/v1/products` | List products (with filters) |
| GET | `/api/v1/products/:id` | Get product details |

### Authenticated User Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/me` | Get current user profile |
| POST | `/api/v1/cart` | Add item to cart |
| GET | `/api/v1/cart` | View cart |
| DELETE | `/api/v1/cart/:item_id` | Remove from cart |
| POST | `/api/v1/orders` | Create order from cart |
| GET | `/api/v1/orders` | List user orders |
| GET | `/api/v1/orders/:id` | Get order details |
| POST | `/api/v1/payments/charge` | Process payment |

### Admin Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/products` | Create product |
| PUT | `/api/v1/products/:id` | Update product |
| DELETE | `/api/v1/products/:id` | Delete product |
| POST | `/api/v1/products/bulk` | Bulk import products |
| GET | `/api/v1/admin/orders` | List all orders |
| PATCH | `/api/v1/admin/orders/:id` | Update order status |

---

## 🎯 Complete Usage Example

### 1. Register a New User

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "full_name": "Test User"
  }'
```

### 2. Login (Get JWT Token)

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'

# Response:
# {
#   "access_token": "eyJhbGc...",
#   "token_type": "bearer",
#   "expires_in": 86400
# }

# Save the token
export TOKEN="<access_token_from_response>"
```

### 3. List Products

```bash
curl http://localhost:8080/api/v1/products
```

### 4. Add Product to Cart

```bash
# First, get a product ID from the products list
export PRODUCT_ID="<product_id>"

curl -X POST http://localhost:8080/api/v1/cart \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": "'$PRODUCT_ID'",
    "quantity": 2
  }'
```

### 5. View Cart

```bash
curl http://localhost:8080/api/v1/cart \
  -H "Authorization: Bearer $TOKEN"
```

### 6. Create Order

```bash
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

# Save the order_id from response
export ORDER_ID="<order_id>"
```

### 7. Process Payment

```bash
curl -X POST http://localhost:8080/api/v1/payments/charge \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "order_id": "'$ORDER_ID'",
    "payment_method": "card",
    "payment_details": {}
  }'
```

---

## 🔐 Sample Credentials (After Seeding)

```
Admin User:
  Email: admin@example.com
  Password: admin123

Regular User:
  Email: user@example.com
  Password: admin123
```

**⚠️ Change these in production!**

---

## 📁 Project Structure

```
goecom/
├── cmd/server/              # Application entry point
├── internal/
│   ├── config/              # Configuration loading
│   ├── handler/             # HTTP request handlers
│   ├── middleware/          # Auth, logging, rate limiting
│   ├── server/              # Router setup
│   ├── service/             # Business logic
│   └── store/               # Database models & repositories
├── migrations/              # SQL migration files
├── scripts/                 # Helper scripts (migrate, seed)
├── test/                    # Integration tests
├── .github/workflows/       # CI/CD pipeline
├── docker-compose.yml       # Local development environment
├── Dockerfile               # Production container
├── Makefile                 # Convenience commands
├── openapi.yaml             # API specification
├── API_EXAMPLES.md          # Complete API examples
└── README.md                # Main documentation
```

---

## 🚢 Deployment Checklist

- [ ] Set strong `JWT_SECRET` (32+ characters)
- [ ] Set `ENV=production`
- [ ] Set `BCRYPT_COST=12` or higher
- [ ] Configure proper `CORS_ORIGINS`
- [ ] Use SSL/TLS for database connection
- [ ] Set up database backups
- [ ] Configure logging/monitoring
- [ ] Set up reverse proxy (nginx/traefik)
- [ ] Enable HTTPS
- [ ] Review rate limiting settings
- [ ] Change default admin credentials
- [ ] Set up CI/CD pipeline
- [ ] Configure secrets management

---

## 🆘 Troubleshooting

### Database Connection Failed

```bash
# Check if PostgreSQL is running
docker-compose ps
# or
pg_isready -h localhost -p 5432

# Check connection string
echo $DATABASE_URL

# Test connection manually
psql "$DATABASE_URL" -c "SELECT version();"
```

### Port Already in Use

```bash
# Windows - Find and kill process using port 8080
netstat -ano | findstr :8080
taskkill /PID <PID> /F

# Linux/Mac
lsof -ti :8080 | xargs kill -9
```

### Migration Errors

```bash
# Rollback and retry
scripts/migrate.sh down
scripts/migrate.sh up

# Or reset database completely
docker-compose down -v
docker-compose up -d db
# Wait for DB to be ready
scripts/migrate.sh up
```

### Import Errors

```bash
# Download dependencies
go mod download
go mod tidy

# Clear cache if needed
go clean -modcache
go mod download
```

---

## 📚 Additional Documentation

- **API Examples**: See `API_EXAMPLES.md` for complete curl/HTTPie examples
- **OpenAPI Spec**: View `openapi.yaml` in Swagger UI at http://localhost:8081
- **Migrations**: SQL files in `migrations/` directory
- **Tests**: Unit tests in `internal/*/` and integration tests in `test/`

---

## ✅ Acceptance Criteria Verification

1. ✅ **Docker Compose**: `docker-compose up` brings up app + database
2. ✅ **User Flow**: Register → Login → Add to cart → Create order works
3. ✅ **Admin Protection**: Admin endpoints require admin token
4. ✅ **Tests Pass**: `go test ./...` returns success
5. ✅ **Linter Passes**: `golangci-lint run` returns exit code 0
6. ✅ **OpenAPI Spec**: `openapi.yaml` accurately describes all endpoints
7. ✅ **Documentation**: README contains full setup and examples

---

## 🎉 Success!

Your e-commerce backend is now ready to use!

**Next Steps:**
1. Explore the API using Swagger UI: http://localhost:8081
2. Try the complete user flow with the example commands above
3. Check out `API_EXAMPLES.md` for more detailed examples
4. Review `openapi.yaml` for complete API documentation
5. Run tests to verify everything works: `go test ./...`

**Happy coding! 🚀**
