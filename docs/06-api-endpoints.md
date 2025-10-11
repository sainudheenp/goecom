# 6. API Endpoints 🌐

## Base URL

```
http://localhost:8080
```

## API Overview

All API endpoints are under `/api/v1`:

```
/api/v1/auth/*        # Authentication
/api/v1/products/*    # Product catalog
/api/v1/me            # User profile
```

## Authentication Endpoints

### Register User

**POST** `/api/v1/auth/register`

Create a new user account.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123",
  "full_name": "John Doe"
}
```

**Response:** `201 Created`
```json
{
  "user": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "email": "user@example.com",
    "full_name": "John Doe",
    "role": "user",
    "created_at": "2024-10-11T10:30:00Z",
    "updated_at": "2024-10-11T10:30:00Z"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Errors:**
- `400` - Invalid request (validation failed)
- `400` - User already exists

**Example:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123",
    "full_name": "John Doe"
  }'
```

---

### Login

**POST** `/api/v1/auth/login`

Authenticate and get JWT token.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response:** `200 OK`
```json
{
  "user": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "email": "user@example.com",
    "full_name": "John Doe",
    "role": "user",
    "created_at": "2024-10-11T10:30:00Z",
    "updated_at": "2024-10-11T10:30:00Z"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Errors:**
- `400` - Invalid request
- `401` - Invalid credentials

**Example:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

---

## Product Endpoints

### List Products

**GET** `/api/v1/products`

Get paginated list of products with optional search.

**Query Parameters:**
- `q` - Search query (searches name and description)
- `page` - Page number (default: 1)
- `size` - Items per page (default: 20)

**Response:** `200 OK`
```json
{
  "products": [
    {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "sku": "LAPTOP-001",
      "name": "MacBook Pro 14\"",
      "description": "Powerful laptop for developers",
      "price_cents": 199900,
      "currency": "USD",
      "stock": 10,
      "images": ["https://example.com/macbook.jpg"],
      "created_at": "2024-10-11T10:30:00Z",
      "updated_at": "2024-10-11T10:30:00Z"
    }
  ],
  "total": 100,
  "page": 1,
  "size": 20
}
```

**Examples:**
```bash
# Get all products
curl http://localhost:8080/api/v1/products

# Search products
curl http://localhost:8080/api/v1/products?q=laptop

# Pagination
curl http://localhost:8080/api/v1/products?page=2&size=10
```

---

### Get Product by ID

**GET** `/api/v1/products/:id`

Get details of a specific product.

**Parameters:**
- `id` - Product UUID

**Response:** `200 OK`
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "sku": "LAPTOP-001",
  "name": "MacBook Pro 14\"",
  "description": "Powerful laptop for developers",
  "price_cents": 199900,
  "currency": "USD",
  "stock": 10,
  "images": ["https://example.com/macbook.jpg"],
  "created_at": "2024-10-11T10:30:00Z",
  "updated_at": "2024-10-11T10:30:00Z"
}
```

**Errors:**
- `400` - Invalid product ID
- `404` - Product not found

**Example:**
```bash
curl http://localhost:8080/api/v1/products/123e4567-e89b-12d3-a456-426614174000
```

---

## Protected Endpoints

These endpoints require authentication. Include JWT token in Authorization header.

### Get Current User

**GET** `/api/v1/me`

Get authenticated user's profile.

**Headers:**
```
Authorization: Bearer <token>
```

**Response:** `200 OK`
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "email": "user@example.com",
  "full_name": "John Doe",
  "role": "user",
  "created_at": "2024-10-11T10:30:00Z",
  "updated_at": "2024-10-11T10:30:00Z"
}
```

**Errors:**
- `401` - Unauthorized (no token or invalid token)

**Example:**
```bash
TOKEN="your-jwt-token-here"

curl http://localhost:8080/api/v1/me \
  -H "Authorization: Bearer $TOKEN"
```

---

## Health Check

### Server Health

**GET** `/health`

Check if server is running.

**Response:** `200 OK`
```json
{
  "status": "ok",
  "time": "2024-10-11T10:30:00Z"
}
```

**Example:**
```bash
curl http://localhost:8080/health
```

---

## Error Responses

All errors follow a consistent format:

```json
{
  "error": "error message",
  "details": "optional detailed error message"
}
```

### HTTP Status Codes

| Code | Meaning |
|------|---------|
| 200 | OK - Request successful |
| 201 | Created - Resource created successfully |
| 400 | Bad Request - Invalid input |
| 401 | Unauthorized - Authentication required |
| 403 | Forbidden - Insufficient permissions |
| 404 | Not Found - Resource doesn't exist |
| 500 | Internal Server Error - Server error |

---

## Making Requests

### Using curl

```bash
# GET request
curl http://localhost:8080/api/v1/products

# POST request with JSON body
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "pass123"}'

# With authentication
curl http://localhost:8080/api/v1/me \
  -H "Authorization: Bearer your-token-here"

# Pretty print JSON response
curl http://localhost:8080/api/v1/products | jq
```

### Using JavaScript (Fetch API)

```javascript
// GET request
const response = await fetch('http://localhost:8080/api/v1/products');
const data = await response.json();

// POST request
const response = await fetch('http://localhost:8080/api/v1/auth/login', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    email: 'user@example.com',
    password: 'password123'
  })
});
const data = await response.json();

// With authentication
const token = localStorage.getItem('token');
const response = await fetch('http://localhost:8080/api/v1/me', {
  headers: {
    'Authorization': `Bearer ${token}`
  }
});
```

### Using axios (Node.js/React)

```javascript
import axios from 'axios';

const API = axios.create({
  baseURL: 'http://localhost:8080/api/v1'
});

// Add token to all requests
API.interceptors.request.use(config => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Usage
const { data } = await API.get('/products');
const { data } = await API.post('/auth/login', { email, password });
const { data } = await API.get('/me');
```

---

## Testing with Postman

### Setup

1. Create new Collection: "E-Commerce API"
2. Set Collection Variable: `baseUrl` = `http://localhost:8080`
3. Set Collection Variable: `token` = (leave empty for now)

### Authentication Flow

1. **Register:**
   - POST `{{baseUrl}}/api/v1/auth/register`
   - Body: `{"email": "test@example.com", "password": "pass123", "full_name": "Test User"}`
   - Save `token` from response

2. **Login:**
   - POST `{{baseUrl}}/api/v1/auth/login`
   - Body: `{"email": "test@example.com", "password": "pass123"}`
   - Update Collection Variable `token` with response

3. **Get Profile:**
   - GET `{{baseUrl}}/api/v1/me`
   - Header: `Authorization: Bearer {{token}}`

---

## Rate Limiting

The API implements rate limiting to prevent abuse:

- **Limit:** 100 requests per minute per IP
- **Headers:**
  - `X-RateLimit-Limit` - Maximum requests allowed
  - `X-RateLimit-Remaining` - Remaining requests
  - `X-RateLimit-Reset` - Time when limit resets

**Response when limit exceeded:**
```json
{
  "error": "rate limit exceeded"
}
```

---

## CORS

Cross-Origin Resource Sharing is enabled for these origins (configurable in `.env`):

```env
CORS_ORIGINS=http://localhost:3000,http://localhost:5173
```

Allowed methods: `GET`, `POST`, `PUT`, `PATCH`, `DELETE`, `OPTIONS`

---

## Request/Response Examples

### Complete Registration + Login Flow

```bash
# 1. Register
REGISTER_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "demo@example.com",
    "password": "demo123456",
    "full_name": "Demo User"
  }')

echo $REGISTER_RESPONSE | jq

# 2. Extract token
TOKEN=$(echo $REGISTER_RESPONSE | jq -r '.token')

# 3. Get profile
curl http://localhost:8080/api/v1/me \
  -H "Authorization: Bearer $TOKEN" | jq

# 4. Get products
curl http://localhost:8080/api/v1/products | jq
```

### Search and Pagination

```bash
# Search for laptops
curl 'http://localhost:8080/api/v1/products?q=laptop' | jq

# Get page 2 with 5 items
curl 'http://localhost:8080/api/v1/products?page=2&size=5' | jq

# Search + pagination
curl 'http://localhost:8080/api/v1/products?q=phone&page=1&size=10' | jq
```

---

## API Client Libraries

### Go Client Example

```go
package main

import (
    "bytes"
    "encoding/json"
    "net/http"
)

type Client struct {
    baseURL string
    token   string
}

func (c *Client) Login(email, password string) error {
    body := map[string]string{
        "email":    email,
        "password": password,
    }
    
    jsonData, _ := json.Marshal(body)
    
    resp, err := http.Post(
        c.baseURL+"/api/v1/auth/login",
        "application/json",
        bytes.NewBuffer(jsonData),
    )
    
    var result struct {
        Token string `json:"token"`
    }
    json.NewDecoder(resp.Body).Decode(&result)
    
    c.token = result.Token
    return nil
}

func (c *Client) GetProducts() ([]Product, error) {
    req, _ := http.NewRequest("GET", c.baseURL+"/api/v1/products", nil)
    req.Header.Set("Authorization", "Bearer "+c.token)
    
    resp, err := http.DefaultClient.Do(req)
    // Parse response...
}
```

---

## Next Steps

Now that you know the API endpoints:

1. 📖 Read **Code Walkthrough** to understand implementation
2. 🎯 Try **Common Tasks** to add new features
3. 🧪 Explore **Testing & Debugging** for development tips

---

**Next:** [Code Walkthrough →](./07-code-walkthrough.md)
