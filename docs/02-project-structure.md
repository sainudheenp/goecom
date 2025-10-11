# 2. Project Structure 📁

## Overview

This project follows a **flat, intermediate-level structure** - simpler than complex enterprise patterns but more organized than a single-file application.

```
goecom/
├── cmd/                    # Application entry points
│   └── server/
│       └── main.go        # Main entry point (like index.js)
├── config/                 # Configuration management
│   └── config.go          # Load env variables
├── db/                     # Database layer
│   └── database.go        # DB connection setup
├── handlers/               # HTTP request handlers (like Express routes)
│   ├── auth_handler.go    # Register, Login
│   └── product_handler.go # Product CRUD
├── middleware/             # HTTP middleware
│   ├── auth.go            # JWT authentication
│   ├── logger.go          # Request logging
│   ├── ratelimit.go       # Rate limiting
│   └── recovery.go        # Panic recovery
├── models/                 # Data models (like Mongoose models)
│   ├── models.go          # User, Product, Order, etc.
│   └── json_types.go      # Custom JSON types
├── server/                 # Server setup and routing
│   └── server.go          # Router configuration
├── migrations/             # Database migrations
│   ├── 001_create_users_table.up.sql
│   ├── 001_create_users_table.down.sql
│   └── ...
├── scripts/                # Utility scripts
│   ├── migrate.sh         # Run migrations
│   └── seed.sh            # Seed data
├── docs/                   # Documentation (you are here!)
├── .env                    # Environment variables (gitignored)
├── .env.example           # Example environment file
├── go.mod                  # Dependencies (like package.json)
├── go.sum                  # Lock file (like package-lock.json)
└── README.md              # Project readme
```

## 🔍 Comparing with Node.js Structure

### Node.js Express (typical structure):
```
my-api/
├── src/
│   ├── controllers/       # Request handlers
│   ├── models/            # Database models
│   ├── routes/            # Route definitions
│   ├── middlewares/       # Express middleware
│   ├── services/          # Business logic
│   └── utils/             # Helper functions
├── config/
│   └── database.js
├── index.js               # Entry point
└── package.json
```

### This Go Project:
```
goecom/
├── cmd/server/main.go     # Entry point
├── handlers/              # Controllers + Routes combined
├── models/                # Database models
├── middleware/            # Middleware
├── config/                # Configuration
├── db/                    # Database setup
└── go.mod                 # Dependencies
```

## Key Differences

| Aspect | Node.js | Go (This Project) |
|--------|---------|-------------------|
| Entry Point | `index.js` | `cmd/server/main.go` |
| Dependencies | `package.json` | `go.mod` |
| Lock File | `package-lock.json` | `go.sum` |
| Routes | Separate `routes/` folder | In `server/server.go` |
| Controllers | `controllers/` | `handlers/` |
| Business Logic | `services/` | Inline in handlers (simplified) |
| Models | Mongoose/Sequelize | GORM structs |

## Detailed Breakdown

### 1. `cmd/server/main.go` - Entry Point

🔍 **Compare:**
```javascript
// Node.js: index.js
const express = require('express');
const app = express();

// Load config
require('dotenv').config();

// Setup routes
app.use('/api', routes);

// Start server
app.listen(3000, () => {
  console.log('Server running on port 3000');
});
```

```go
// Go: cmd/server/main.go
package main

import (
    "github.com/sainudheenp/goecom/config"
    "github.com/sainudheenp/goecom/server"
)

func main() {
    // Load configuration
    cfg, err := config.Load()
    
    // Create server
    srv, err := server.NewServer(cfg)
    
    // Run server
    srv.Run()
}
```

**Purpose:** Bootstrap the application

### 2. `config/config.go` - Configuration

📌 **Remember:** Go doesn't have a built-in way to load `.env` files (unlike Node.js's `dotenv`)

```go
// Loads from .env file
type Config struct {
    Server    ServerConfig
    Database  DatabaseConfig
    JWT       JWTConfig
}

func Load() (*Config, error) {
    // Load .env file
    godotenv.Load()
    
    // Read environment variables
    return &Config{
        Server: ServerConfig{
            Port: os.Getenv("SERVER_PORT"),
        },
        // ...
    }
}
```

### 3. `db/database.go` - Database Connection

🔍 **Compare:**
```javascript
// Node.js with Mongoose
const mongoose = require('mongoose');
mongoose.connect(process.env.DATABASE_URL);
```

```go
// Go with GORM
func NewDB(databaseURL string) (*DB, error) {
    db, err := gorm.Open(postgres.Open(databaseURL))
    // Configure connection pool
    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetMaxOpenConns(100)
    return &DB{db}, nil
}
```

### 4. `models/models.go` - Data Models

🔍 **Compare:**
```javascript
// Node.js (Mongoose)
const userSchema = new mongoose.Schema({
  email: { type: String, required: true, unique: true },
  password: { type: String, required: true },
  fullName: String,
  role: { type: String, default: 'user' }
});

module.exports = mongoose.model('User', userSchema);
```

```go
// Go (GORM)
type User struct {
    ID           uuid.UUID `gorm:"type:uuid;primary_key;"`
    Email        string    `gorm:"uniqueIndex;not null"`
    PasswordHash string    `gorm:"not null" json:"-"`
    FullName     string
    Role         string    `gorm:"not null;default:'user'"`
    CreatedAt    time.Time
    UpdatedAt    time.Time
}
```

**Key Differences:**
- Go uses **struct tags** for validation and DB mapping
- Go models are **structs** (like TypeScript interfaces)
- GORM is more explicit than Mongoose

### 5. `handlers/` - Request Handlers

🔍 **Compare:**
```javascript
// Node.js (Express)
router.post('/register', async (req, res) => {
  try {
    const { email, password, fullName } = req.body;
    // Hash password
    const hash = await bcrypt.hash(password, 10);
    // Create user
    const user = await User.create({ email, password: hash, fullName });
    // Return response
    res.status(201).json({ user });
  } catch (error) {
    res.status(400).json({ error: error.message });
  }
});
```

```go
// Go (Gin)
func (h *AuthHandler) Register(c *gin.Context) {
    var req RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": "invalid request"})
        return
    }
    
    // Hash password
    hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
    
    // Create user
    user := &models.User{
        Email:        req.Email,
        PasswordHash: string(hash),
        FullName:     req.FullName,
    }
    h.db.Create(user)
    
    c.JSON(201, gin.H{"user": user})
}
```

**Key Differences:**
- Go uses `c.JSON()` vs Express's `res.json()`
- Error handling is explicit (no try-catch)
- Go uses structs for request/response

### 6. `middleware/` - HTTP Middleware

🔍 **Compare:**
```javascript
// Node.js (Express)
const authMiddleware = (req, res, next) => {
  const token = req.headers.authorization?.split(' ')[1];
  if (!token) {
    return res.status(401).json({ error: 'No token' });
  }
  
  try {
    const decoded = jwt.verify(token, process.env.JWT_SECRET);
    req.user = decoded;
    next();
  } catch (error) {
    res.status(401).json({ error: 'Invalid token' });
  }
};
```

```go
// Go (Gin)
func AuthMiddleware(db *gorm.DB, secret string) gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(401, gin.H{"error": "No token"})
            c.Abort()
            return
        }
        
        token := strings.Split(authHeader, " ")[1]
        // Verify JWT...
        
        c.Set("user", user)
        c.Next()
    }
}
```

**Key Differences:**
- Go middleware returns a function
- Use `c.Abort()` instead of `return`
- Use `c.Next()` to continue (explicit)

### 7. `server/server.go` - Router Setup

This is where all routes are defined:

```go
func (s *Server) setupRoutes() {
    // Initialize handlers
    authHandler := handler.NewAuthHandler(...)
    productHandler := handler.NewProductHandler(...)
    
    // API v1 routes
    v1 := s.router.Group("/api/v1")
    {
        // Public routes
        auth := v1.Group("/auth")
        {
            auth.POST("/register", authHandler.Register)
            auth.POST("/login", authHandler.Login)
        }
        
        // Protected routes
        protected := v1.Group("")
        protected.Use(middleware.AuthMiddleware(...))
        {
            protected.GET("/me", authHandler.GetMe)
        }
    }
}
```

## File Naming Conventions

### Node.js Common Patterns:
- `userController.js`
- `user.controller.js`
- `user-controller.js`

### Go Conventions:
- `user_handler.go` ✅ (snake_case)
- `auth_middleware.go` ✅
- `config.go` ✅

📌 **Remember:** Go uses **snake_case** for file names, **PascalCase** for exported functions

## Package Organization

### Important Concept: Packages

In Node.js:
```javascript
// user.js
module.exports = { createUser, getUser };

// In another file
const { createUser } = require('./user');
```

In Go:
```go
// handlers/auth_handler.go
package handler  // Package name

// Exported (starts with capital)
func Register() { }

// In another file
import "github.com/sainudheenp/goecom/handlers"
handler.Register()
```

💡 **Tip:** 
- **Exported** (public): `RegisterUser` (capital first letter)
- **Unexported** (private): `hashPassword` (lowercase first letter)

## Why This Structure?

✅ **Simple** - Easy to understand for beginners
✅ **Practical** - Handlers directly interact with database
✅ **Scalable** - Can add more features easily
✅ **Clear** - Each folder has a single responsibility

## Common Questions

**Q: Why no `services/` folder?**
A: For intermediate learning, we keep business logic in handlers. In production, you'd separate into services.

**Q: Why `cmd/server/`?**
A: Go convention. `cmd/` contains entry points for different executables.

**Q: Can I rename folders?**
A: Yes, but update all imports! Go uses full import paths.

**Q: Where do I add new features?**
A: See [Common Tasks](./08-common-tasks.md) guide

---

**Next:** [Go vs Node.js Comparison →](./03-go-vs-nodejs.md)
