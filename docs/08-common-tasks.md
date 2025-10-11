# 8. Common Tasks - How-To Guide 🛠️

## Adding a New Feature

This section walks you through common development tasks with step-by-step examples.

## Task 1: Add a New Field to User Model

**Goal:** Add a `phone` field to the User model.

### Step 1: Update the Model

**File:** `models/models.go`

```go
type User struct {
    ID           uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
    Email        string    `gorm:"uniqueIndex;not null" json:"email"`
    PasswordHash string    `gorm:"not null" json:"-"`
    FullName     string    `json:"full_name"`
    Phone        string    `json:"phone"`  // ← ADD THIS
    Role         string    `gorm:"not null;default:'user'" json:"role"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}
```

### Step 2: Create Migration

**File:** `migrations/006_add_phone_to_users.up.sql`

```sql
ALTER TABLE users ADD COLUMN IF NOT EXISTS phone VARCHAR(20);
CREATE INDEX IF NOT EXISTS idx_users_phone ON users(phone);
```

**File:** `migrations/006_add_phone_to_users.down.sql`

```sql
DROP INDEX IF EXISTS idx_users_phone;
ALTER TABLE users DROP COLUMN IF EXISTS phone;
```

### Step 3: Update Register Request

**File:** `handlers/auth_handler.go`

```go
type RegisterRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8"`
    FullName string `json:"full_name" binding:"required"`
    Phone    string `json:"phone"`  // ← ADD THIS
}

// In Register function
user := &models.User{
    Email:        req.Email,
    PasswordHash: string(hashedPassword),
    FullName:     req.FullName,
    Phone:        req.Phone,  // ← ADD THIS
}
```

### Step 4: Run Migration

```bash
./scripts/migrate.sh up
```

### Step 5: Test

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "full_name": "Test User",
    "phone": "+1234567890"
  }'
```

---

## Task 2: Add a New Endpoint

**Goal:** Add an endpoint to update user profile.

### Step 1: Create Request Struct

**File:** `handlers/auth_handler.go`

```go
type UpdateProfileRequest struct {
    FullName string `json:"full_name"`
    Phone    string `json:"phone"`
}
```

### Step 2: Create Handler Function

```go
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
    // Get current user
    user, err := middleware.GetUserFromContext(c)
    if err != nil {
        c.JSON(401, gin.H{"error": "unauthorized"})
        return
    }
    
    // Parse request
    var req UpdateProfileRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": "invalid request"})
        return
    }
    
    // Update user
    updates := map[string]interface{}{}
    if req.FullName != "" {
        updates["full_name"] = req.FullName
    }
    if req.Phone != "" {
        updates["phone"] = req.Phone
    }
    
    if err := h.db.Model(user).Updates(updates).Error; err != nil {
        c.JSON(500, gin.H{"error": "update failed"})
        return
    }
    
    // Return updated user
    c.JSON(200, user)
}
```

### Step 3: Add Route

**File:** `server/server.go`

```go
// In setupRoutes()
protected := v1.Group("")
protected.Use(middleware.AuthMiddleware(s.db.DB, s.config.JWT.Secret))
{
    protected.GET("/me", authHandler.GetMe)
    protected.PUT("/me", authHandler.UpdateProfile)  // ← ADD THIS
}
```

### Step 4: Test

```bash
TOKEN="your-jwt-token"

curl -X PUT http://localhost:8080/api/v1/me \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "Updated Name",
    "phone": "+1987654321"
  }'
```

---

## Task 3: Add Search Filters

**Goal:** Add price range filtering to products.

### Step 1: Update Handler

**File:** `handlers/product_handler.go`

```go
func (h *ProductHandler) ListProducts(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
    q := c.Query("q")
    minPrice := c.Query("min_price")  // ← ADD THIS
    maxPrice := c.Query("max_price")  // ← ADD THIS
    
    var products []models.Product
    dbQuery := h.db.Model(&models.Product{})
    
    // Search query
    if q != "" {
        dbQuery = dbQuery.Where(
            "name ILIKE ? OR description ILIKE ?",
            "%"+q+"%",
            "%"+q+"%",
        )
    }
    
    // Price filter
    if minPrice != "" {
        min, _ := strconv.Atoi(minPrice)
        dbQuery = dbQuery.Where("price_cents >= ?", min)
    }
    if maxPrice != "" {
        max, _ := strconv.Atoi(maxPrice)
        dbQuery = dbQuery.Where("price_cents <= ?", max)
    }
    
    // Count total
    var total int64
    dbQuery.Count(&total)
    
    // Pagination
    offset := (page - 1) * size
    dbQuery.Limit(size).Offset(offset).Find(&products)
    
    c.JSON(200, gin.H{
        "products": products,
        "total":    total,
        "page":     page,
        "size":     size,
    })
}
```

### Step 2: Test

```bash
# Products between $500 and $2000
curl 'http://localhost:8080/api/v1/products?min_price=50000&max_price=200000'

# Laptops under $1500
curl 'http://localhost:8080/api/v1/products?q=laptop&max_price=150000'
```

---

## Task 4: Add Email Validation

**Goal:** Validate email domain before registration.

### Step 1: Create Validation Function

**File:** `handlers/auth_handler.go`

```go
func isValidEmailDomain(email string) bool {
    // Block temporary email domains
    blockedDomains := []string{
        "tempmail.com",
        "throwaway.email",
        "guerrillamail.com",
    }
    
    parts := strings.Split(email, "@")
    if len(parts) != 2 {
        return false
    }
    
    domain := parts[1]
    for _, blocked := range blockedDomains {
        if domain == blocked {
            return false
        }
    }
    
    return true
}
```

### Step 2: Use in Register

```go
func (h *AuthHandler) Register(c *gin.Context) {
    var req RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": "invalid request"})
        return
    }
    
    // Validate email domain
    if !isValidEmailDomain(req.Email) {
        c.JSON(400, gin.H{"error": "invalid email domain"})
        return
    }
    
    // Continue with registration...
}
```

---

## Task 5: Add Logging

**Goal:** Log all API requests.

### Middleware Already Exists!

**File:** `middleware/logger.go`

```go
func Logger() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        path := c.Request.URL.Path
        
        // Process request
        c.Next()
        
        // Log after processing
        duration := time.Since(start)
        log.Printf(
            "%s %s %d %s",
            c.Request.Method,
            path,
            c.Writer.Status(),
            duration,
        )
    }
}
```

Already enabled in `server/server.go`:
```go
router.Use(middleware.Logger())
```

### Custom Logging

Add structured logging:

```bash
go get github.com/sirupsen/logrus
```

**File:** `middleware/logger.go`

```go
import "github.com/sirupsen/logrus"

var log = logrus.New()

func Logger() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        
        c.Next()
        
        log.WithFields(logrus.Fields{
            "method":   c.Request.Method,
            "path":     c.Request.URL.Path,
            "status":   c.Writer.Status(),
            "duration": time.Since(start),
            "ip":       c.ClientIP(),
        }).Info("API Request")
    }
}
```

---

## Task 6: Add Input Sanitization

**Goal:** Sanitize user input to prevent XSS.

### Step 1: Install Package

```bash
go get github.com/microcosm-cc/bluemonday
```

### Step 2: Create Sanitization Function

**File:** `handlers/auth_handler.go`

```go
import "github.com/microcosm-cc/bluemonday"

var sanitizer = bluemonday.StrictPolicy()

func sanitizeString(input string) string {
    return sanitizer.Sanitize(input)
}
```

### Step 3: Use in Handlers

```go
func (h *AuthHandler) Register(c *gin.Context) {
    var req RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": "invalid request"})
        return
    }
    
    // Sanitize inputs
    req.Email = sanitizeString(req.Email)
    req.FullName = sanitizeString(req.FullName)
    
    // Continue...
}
```

---

## Task 7: Add Pagination Helper

**Goal:** Reusable pagination logic.

### Create Helper Function

**File:** `models/pagination.go`

```go
package models

type Pagination struct {
    Page  int   `json:"page"`
    Size  int   `json:"size"`
    Total int64 `json:"total"`
}

func (p *Pagination) Offset() int {
    return (p.Page - 1) * p.Size
}

func (p *Pagination) TotalPages() int {
    return int((p.Total + int64(p.Size) - 1) / int64(p.Size))
}

func NewPagination(page, size int) *Pagination {
    if page < 1 {
        page = 1
    }
    if size < 1 || size > 100 {
        size = 20
    }
    return &Pagination{
        Page: page,
        Size: size,
    }
}
```

### Use in Handler

```go
func (h *ProductHandler) ListProducts(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
    
    pagination := models.NewPagination(page, size)
    
    var products []models.Product
    var total int64
    
    h.db.Model(&models.Product{}).Count(&total)
    pagination.Total = total
    
    h.db.Limit(pagination.Size).
        Offset(pagination.Offset()).
        Find(&products)
    
    c.JSON(200, gin.H{
        "products":   products,
        "pagination": pagination,
    })
}
```

---

## Task 8: Add Error Handling Wrapper

**Goal:** Consistent error responses.

### Create Error Helper

**File:** `handlers/errors.go`

```go
package handler

import "github.com/gin-gonic/gin"

type ErrorResponse struct {
    Error   string `json:"error"`
    Details string `json:"details,omitempty"`
}

func BadRequest(c *gin.Context, message string) {
    c.JSON(400, ErrorResponse{Error: message})
}

func Unauthorized(c *gin.Context, message string) {
    c.JSON(401, ErrorResponse{Error: message})
}

func NotFound(c *gin.Context, message string) {
    c.JSON(404, ErrorResponse{Error: message})
}

func InternalError(c *gin.Context, message string) {
    c.JSON(500, ErrorResponse{Error: message})
}
```

### Use in Handlers

```go
func (h *AuthHandler) Register(c *gin.Context) {
    var req RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        BadRequest(c, "invalid request")
        return
    }
    
    if err := h.db.Create(user).Error; err != nil {
        InternalError(c, "failed to create user")
        return
    }
    
    c.JSON(201, RegisterResponse{User: *user, Token: token})
}
```

---

## Task 9: Add Environment-Based Config

**Goal:** Different settings for dev/prod.

### Update Config

**File:** `config/config.go`

```go
func (c *Config) IsDevelopment() bool {
    return c.Server.Env == "development"
}

func (c *Config) IsProduction() bool {
    return c.Server.Env == "production"
}

func (c *Config) GetLogLevel() string {
    if c.IsProduction() {
        return "error"
    }
    return "debug"
}
```

### Use in Code

```go
if cfg.IsDevelopment() {
    log.Println("Running in development mode")
    // Enable debug features
}

if cfg.IsProduction() {
    // Enforce HTTPS
    // Disable debug endpoints
}
```

---

## Task 10: Add Request Timeout

**Goal:** Prevent slow requests from hanging.

### Add Middleware

**File:** `middleware/timeout.go`

```go
package middleware

import (
    "context"
    "time"
    
    "github.com/gin-gonic/gin"
)

func Timeout(duration time.Duration) gin.HandlerFunc {
    return func(c *gin.Context) {
        ctx, cancel := context.WithTimeout(c.Request.Context(), duration)
        defer cancel()
        
        c.Request = c.Request.WithContext(ctx)
        
        done := make(chan struct{})
        go func() {
            c.Next()
            done <- struct{}{}
        }()
        
        select {
        case <-done:
            return
        case <-ctx.Done():
            c.JSON(408, gin.H{"error": "request timeout"})
            c.Abort()
        }
    }
}
```

### Use in Server

```go
router.Use(middleware.Timeout(30 * time.Second))
```

---

## Quick Reference Checklist

When adding a new feature:

- [ ] Update model (if needed)
- [ ] Create migration (if database change)
- [ ] Add handler function
- [ ] Update routes
- [ ] Add validation
- [ ] Handle errors
- [ ] Test with curl
- [ ] Update documentation
- [ ] Commit changes

---

**Next:** [Testing & Debugging →](./09-testing-debugging.md)
