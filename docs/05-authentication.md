# 5. Authentication 🔐

## JWT Authentication

This project uses **JWT (JSON Web Tokens)** for authentication, similar to what you'd use in Node.js.

## How Authentication Works

```
1. User registers/logs in
2. Server creates JWT token
3. Client stores token (localStorage/cookies)
4. Client sends token in Authorization header
5. Server validates token on protected routes
```

## JWT Flow Diagram

```
┌────────┐              ┌────────┐              ┌──────────┐
│ Client │              │ Server │              │ Database │
└───┬────┘              └───┬────┘              └────┬─────┘
    │                       │                        │
    │ POST /register        │                        │
    ├──────────────────────>│                        │
    │                       │ Hash Password          │
    │                       │ Save User              │
    │                       ├───────────────────────>│
    │                       │<───────────────────────┤
    │                       │ Generate JWT           │
    │<──────────────────────┤                        │
    │ {token, user}         │                        │
    │                       │                        │
    │ GET /me               │                        │
    │ Authorization: Bearer │                        │
    ├──────────────────────>│                        │
    │                       │ Verify JWT             │
    │                       │ Get User               │
    │                       ├───────────────────────>│
    │                       │<───────────────────────┤
    │<──────────────────────┤                        │
    │ {user data}           │                        │
```

## Registration Flow

### Location: `handlers/auth_handler.go`

🔍 **Compare with Node.js:**
```javascript
// Node.js (Express + bcrypt + jwt)
router.post('/register', async (req, res) => {
    const { email, password, fullName } = req.body;
    
    // Hash password
    const hash = await bcrypt.hash(password, 10);
    
    // Create user
    const user = await User.create({
        email,
        password: hash,
        fullName
    });
    
    // Generate JWT
    const token = jwt.sign(
        { userId: user._id },
        process.env.JWT_SECRET,
        { expiresIn: '7d' }
    );
    
    res.json({ user, token });
});
```

```go
// Go (Gin + bcrypt + jwt)
func (h *AuthHandler) Register(c *gin.Context) {
    var req RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": "invalid request"})
        return
    }
    
    // Hash password
    hashedPassword, err := bcrypt.GenerateFromPassword(
        []byte(req.Password), 
        h.bcryptCost,
    )
    
    // Create user
    user := &models.User{
        Email:        req.Email,
        PasswordHash: string(hashedPassword),
        FullName:     req.FullName,
    }
    
    if err := h.db.Create(user).Error; err != nil {
        c.JSON(400, gin.H{"error": "user already exists"})
        return
    }
    
    // Generate JWT
    token, _ := h.generateToken(user.ID)
    
    c.JSON(201, RegisterResponse{
        User:  *user,
        Token: token,
    })
}
```

### Request/Response Structs

```go
// Input validation
type RegisterRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8"`
    FullName string `json:"full_name" binding:"required"`
}

// Output format
type RegisterResponse struct {
    User  models.User `json:"user"`
    Token string      `json:"token"`
}
```

💡 **Tip:** `binding` tags provide automatic validation

## Password Hashing

### bcrypt Implementation

```go
// Hash password
hashedPassword, err := bcrypt.GenerateFromPassword(
    []byte(password),
    10,  // Cost factor (10 = 2^10 iterations)
)

// Verify password
err := bcrypt.CompareHashAndPassword(
    []byte(hashedPassword),
    []byte(plainPassword),
)
if err != nil {
    // Password doesn't match
}
```

### Cost Factor

- **Cost 10**: ~100ms (recommended for production)
- **Cost 12**: ~400ms (higher security)
- **Cost 8**: ~25ms (development only)

⚠️ **Warning:** Never store plain passwords!

```go
// ❌ NEVER DO THIS
user.Password = req.Password

// ✅ Always hash
user.PasswordHash = string(hashedPassword)
```

## JWT Token Generation

### Creating Tokens

```go
func (h *AuthHandler) generateToken(userID uuid.UUID) (string, error) {
    // Create claims
    claims := jwt.MapClaims{
        "user_id": userID.String(),
        "exp":     time.Now().Add(h.jwtExpires).Unix(),
        "iat":     time.Now().Unix(),
    }
    
    // Create token
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    
    // Sign with secret
    return token.SignedString([]byte(h.jwtSecret))
}
```

### JWT Structure

```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMTIzNCIsImV4cCI6MTYzOTM5MjAwMCwiaWF0IjoxNjM4NzkyMDAwfQ.signature
│                                      │                                                            │
│          HEADER                      │                    PAYLOAD                                 │  SIGNATURE
```

**Header:**
```json
{
  "alg": "HS256",
  "typ": "JWT"
}
```

**Payload (Claims):**
```json
{
  "user_id": "123e4567-e89b-12d3-a456-426614174000",
  "exp": 1639392000,
  "iat": 1638792000
}
```

## Login Flow

```go
func (h *AuthHandler) Login(c *gin.Context) {
    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": "invalid request"})
        return
    }
    
    // Find user by email
    var user models.User
    if err := h.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
        c.JSON(401, gin.H{"error": "invalid credentials"})
        return
    }
    
    // Verify password
    if err := bcrypt.CompareHashAndPassword(
        []byte(user.PasswordHash),
        []byte(req.Password),
    ); err != nil {
        c.JSON(401, gin.H{"error": "invalid credentials"})
        return
    }
    
    // Generate token
    token, _ := h.generateToken(user.ID)
    
    c.JSON(200, LoginResponse{
        User:  user,
        Token: token,
    })
}
```

💡 **Tip:** Return same error message for "user not found" and "wrong password" to prevent user enumeration

## Authentication Middleware

### Location: `middleware/auth.go`

```go
func AuthMiddleware(db *gorm.DB, jwtSecret string) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Get Authorization header
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(401, gin.H{"error": "authorization required"})
            c.Abort()
            return
        }
        
        // Extract token (Bearer <token>)
        parts := strings.SplitN(authHeader, " ", 2)
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.JSON(401, gin.H{"error": "invalid header format"})
            c.Abort()
            return
        }
        
        tokenString := parts[1]
        
        // Parse and validate token
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("unexpected signing method")
            }
            return []byte(jwtSecret), nil
        })
        
        if err != nil || !token.Valid {
            c.JSON(401, gin.H{"error": "invalid token"})
            c.Abort()
            return
        }
        
        // Extract user ID from claims
        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            c.JSON(401, gin.H{"error": "invalid claims"})
            c.Abort()
            return
        }
        
        userID, _ := uuid.Parse(claims["user_id"].(string))
        
        // Get user from database
        var user models.User
        if err := db.First(&user, userID).Error; err != nil {
            c.JSON(401, gin.H{"error": "user not found"})
            c.Abort()
            return
        }
        
        // Store user in context
        c.Set("user", &user)
        c.Set("user_id", user.ID)
        
        // Continue to next handler
        c.Next()
    }
}
```

### Using Middleware

```go
// In server/server.go
protected := v1.Group("")
protected.Use(middleware.AuthMiddleware(db, jwtSecret))
{
    protected.GET("/me", authHandler.GetMe)
    protected.POST("/orders", orderHandler.Create)
}
```

### Getting User from Context

```go
func (h *AuthHandler) GetMe(c *gin.Context) {
    user, err := middleware.GetUserFromContext(c)
    if err != nil {
        c.JSON(401, gin.H{"error": "unauthorized"})
        return
    }
    
    c.JSON(200, user)
}

// Helper function
func GetUserFromContext(c *gin.Context) (*models.User, error) {
    userInterface, exists := c.Get("user")
    if !exists {
        return nil, errors.New("user not found in context")
    }
    
    user, ok := userInterface.(*models.User)
    if !ok {
        return nil, errors.New("invalid user type")
    }
    
    return user, nil
}
```

## Making Authenticated Requests

### With curl

```bash
# 1. Login to get token
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'

# Response: {"token": "eyJhbGci...", "user": {...}}

# 2. Use token in protected routes
curl http://localhost:8080/api/v1/me \
  -H "Authorization: Bearer eyJhbGci..."
```

### With JavaScript (Frontend)

```javascript
// Login
const response = await fetch('http://localhost:8080/api/v1/auth/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
        email: 'user@example.com',
        password: 'password123'
    })
});

const { token, user } = await response.json();

// Store token
localStorage.setItem('token', token);

// Use token
const profileResponse = await fetch('http://localhost:8080/api/v1/me', {
    headers: {
        'Authorization': `Bearer ${token}`
    }
});
```

## Role-Based Access Control (RBAC)

### RequireRole Middleware

```go
func RequireRole(roles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        user, err := GetUserFromContext(c)
        if err != nil {
            c.JSON(401, gin.H{"error": "unauthorized"})
            c.Abort()
            return
        }
        
        // Check if user has required role
        hasRole := false
        for _, role := range roles {
            if user.Role == role {
                hasRole = true
                break
            }
        }
        
        if !hasRole {
            c.JSON(403, gin.H{"error": "insufficient permissions"})
            c.Abort()
            return
        }
        
        c.Next()
    }
}
```

### Usage

```go
// Admin only routes
admin := v1.Group("/admin")
admin.Use(middleware.AuthMiddleware(db, secret))
admin.Use(middleware.RequireRole("admin"))
{
    admin.POST("/products", productHandler.Create)
    admin.DELETE("/products/:id", productHandler.Delete)
}
```

## Security Best Practices

### 1. Strong JWT Secret

```bash
# Generate secure secret
openssl rand -base64 32
```

⚠️ **Warning:** Never commit JWT_SECRET to git!

### 2. Token Expiration

```go
// Set reasonable expiration
JWT_EXPIRES_HOURS=168  // 7 days
```

### 3. HTTPS Only

```go
// In production, enforce HTTPS
if env == "production" {
    router.Use(secureMiddleware())
}
```

### 4. Rate Limiting

```go
// Prevent brute force attacks
rateLimiter := middleware.NewRateLimiter(100, 1)  // 100 req/min
router.Use(rateLimiter.Middleware())
```

### 5. Password Requirements

```go
type RegisterRequest struct {
    Password string `binding:"required,min=8,max=72"`
}
```

## Common Auth Issues

### Issue: "invalid token"
**Solution:** Token expired or tampered with. Login again.

### Issue: "authorization required"
**Solution:** Missing Authorization header or Bearer prefix

### Issue: "user not found"
**Solution:** User was deleted or token has old user_id

### Issue: "jwt malformed"
**Solution:** Token format is incorrect. Check Bearer prefix.

## Testing Authentication

### Register Test

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "newuser@example.com",
    "password": "securepass123",
    "full_name": "New User"
  }'
```

### Login Test

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "newuser@example.com",
    "password": "securepass123"
  }'
```

### Protected Route Test

```bash
TOKEN="your-jwt-token-here"

curl http://localhost:8080/api/v1/me \
  -H "Authorization: Bearer $TOKEN"
```

## Practice Exercises

1. **Add email verification** - Require users to verify email before login
2. **Implement refresh tokens** - Allow token renewal without re-login
3. **Add password reset** - Email-based password recovery
4. **Two-factor authentication** - Add TOTP/SMS verification
5. **Session management** - Track active sessions per user

---

**Next:** [API Endpoints →](./06-api-endpoints.md)
