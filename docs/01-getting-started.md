# 1. Getting Started 🚀

## Prerequisites

### What You Need to Install

1. **Go (Golang)** - Version 1.21 or higher
   ```bash
   # Check if Go is installed
   go version
   
   # Download from: https://golang.org/dl/
   ```

2. **PostgreSQL** - Version 15 or higher
   ```bash
   # Check if PostgreSQL is installed
   psql --version
   
   # Ubuntu/Debian
   sudo apt-get install postgresql postgresql-contrib
   ```

3. **Git** (probably already installed)
   ```bash
   git --version
   ```

4. **VS Code** (recommended) or any text editor
   - Install Go extension: `ms-vscode.go`

## Project Setup

### Step 1: Clone the Repository

```bash
git clone https://github.com/sainudheenp/goecom.git
cd goecom
```

### Step 2: Install Go Dependencies

🔍 **Compare with Node.js:**
- Node.js: `npm install` or `yarn install`
- Go: `go mod download`

```bash
# Download all dependencies listed in go.mod
go mod download

# Tidy up (like npm prune)
go mod tidy
```

💡 **Tip:** Go uses `go.mod` (like `package.json`) and `go.sum` (like `package-lock.json`)

### Step 3: Setup Environment Variables

```bash
# Copy example environment file
cp .env.example .env

# Edit .env with your favorite editor
nano .env  # or: vim .env, code .env, etc.
```

**Important variables to set:**

```env
# Database connection
DATABASE_URL=postgresql://username:password@localhost:5432/dbname?sslmode=disable

# Server configuration
SERVER_PORT=8080
SERVER_ENV=development

# JWT Secret (MUST be a strong random string)
JWT_SECRET=your-super-secret-key-min-32-characters-long
JWT_EXPIRES_HOURS=168

# Security
BCRYPT_COST=10

# CORS (for frontend)
CORS_ORIGINS=http://localhost:3000,http://localhost:5173
```

💡 **Tip:** Generate a strong JWT secret:
```bash
openssl rand -base64 32
```

### Step 4: Create Database

```bash
# Connect to PostgreSQL
sudo -u postgres psql

# Create database
CREATE DATABASE ecommerce;

# Create user (optional)
CREATE USER ecomuser WITH PASSWORD 'your-password';
GRANT ALL PRIVILEGES ON DATABASE ecommerce TO ecomuser;

# Exit
\q
```

### Step 5: Run Database Migrations

```bash
# Apply all migrations
./scripts/migrate.sh up

# If you see "migrations completed successfully" - you're good!
```

🔍 **Compare with Node.js:**
- Node.js: Sequelize migrations, Prisma migrate, etc.
- Go: SQL migration files (similar to raw SQL migrations)

### Step 6: Seed Sample Data

```bash
# Populate database with test data
./scripts/seed.sh
```

This creates:
- 2 test users (admin and regular user)
- 10 sample products

**Test Credentials:**
- Admin: `admin@example.com` / `admin123`
- User: `user@example.com` / `admin123`

### Step 7: Run the Server

```bash
# Method 1: Using go run (development)
go run cmd/server/main.go

# Method 2: Build and run (production-like)
go build -o server cmd/server/main.go
./server
```

🔍 **Compare with Node.js:**
```javascript
// Node.js
npm run dev        // or: node server.js
```

```go
// Go
go run cmd/server/main.go
```

You should see:
```
2024/10/11 10:30:00 Running database migrations...
2024/10/11 10:30:01 Database connection established
2024/10/11 10:30:01 Starting server in development mode
2024/10/11 10:30:01 Starting server on :8080
```

## Testing the API

### 1. Health Check

```bash
curl http://localhost:8080/health
```

Expected response:
```json
{
  "status": "ok",
  "time": "2024-10-11T10:30:00Z"
}
```

### 2. Register a New User

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "full_name": "Test User"
  }'
```

### 3. Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

Save the `token` from the response!

### 4. Get Products

```bash
curl http://localhost:8080/api/v1/products
```

### 5. Access Protected Route

```bash
curl http://localhost:8080/api/v1/me \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

## Common Issues & Solutions

### Issue: "go: command not found"
**Solution:** Go is not installed or not in PATH
```bash
# Check Go installation
which go

# Add to PATH (add to ~/.bashrc or ~/.zshrc)
export PATH=$PATH:/usr/local/go/bin
```

### Issue: "Database connection failed"
**Solution:** Check your DATABASE_URL
- Is PostgreSQL running? `sudo systemctl status postgresql`
- Is the database created? `psql -l`
- Are credentials correct?

### Issue: "Port already in use"
**Solution:** Another process is using port 8080
```bash
# Find process using port 8080
sudo lsof -i :8080

# Kill the process
kill -9 <PID>

# Or change port in .env
SERVER_PORT=8081
```

### Issue: "Migration failed"
**Solution:** 
```bash
# Check if database exists
psql -U postgres -l

# Drop and recreate database
psql -U postgres -c "DROP DATABASE ecommerce;"
psql -U postgres -c "CREATE DATABASE ecommerce;"

# Run migrations again
./scripts/migrate.sh up
```

## Next Steps

Now that your server is running:

1. 📖 Read **Project Structure** to understand the codebase
2. 🔍 Check **Go vs Node.js** to understand key differences
3. 🎯 Explore **API Endpoints** to see what's available
4. 💻 Dive into **Code Walkthrough** to understand how it works

## Development Workflow

```bash
# 1. Make changes to code
# 2. Run the server
go run cmd/server/main.go

# 3. Test your changes
curl http://localhost:8080/api/v1/...

# 4. Format code (like prettier)
go fmt ./...

# 5. Check for errors (like eslint)
go vet ./...

# 6. Build for production
go build -o server cmd/server/main.go
```

💡 **Tip:** Install `air` for hot-reloading (like nodemon):
```bash
go install github.com/cosmtrek/air@latest
air
```

---

**Ready?** Let's move to [Project Structure →](./02-project-structure.md)
