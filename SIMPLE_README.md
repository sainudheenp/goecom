# Simple E-Commerce API

## Overview

This is a simplified version of the e-commerce backend with basic features suitable for intermediate developers.

## Features Removed/Simplified

- ❌ Complex layered architecture (services, repositories, interfaces)
- ❌ Payment processing
- ❌ Admin panel and role-based access
- ❌ Rate limiting and advanced middleware  
- ❌ Complex testing with mocks
- ❌ CI/CD pipelines
- ❌ Swagger documentation
- ❌ Migration scripts
- ❌ Docker setup
- ❌ Configuration management

## Basic Features Kept

- ✅ User registration and login with JWT
- ✅ Product listing
- ✅ Shopping cart (add/remove items)
- ✅ Basic order creation
- ✅ Simple database models with GORM

## Files Structure

```
├── simple_main.go     # Main application entry point
├── database.go        # Database setup and models
├── auth.go           # Simple JWT authentication
├── handlers.go       # HTTP handlers for all endpoints
├── simple_go.mod     # Go dependencies
└── .env              # Environment variables
```

## Quick Start

1. **Setup Database**
```bash
# Create PostgreSQL database
createdb simple_ecom
```

2. **Environment Variables**
Create `.env` file:
```
DATABASE_URL=postgres://postgres:password@localhost:5432/simple_ecom?sslmode=disable
PORT=8080
```

3. **Run Application**
```bash
# Install dependencies
go mod tidy

# Run the application
go run simple_main.go database.go auth.go handlers.go
```

## API Endpoints

### Authentication
- `POST /auth/register` - Register new user
- `POST /auth/login` - Login user

### Products
- `GET /products` - Get all products
- `GET /products/:id` - Get single product

### Cart (Protected)
- `GET /cart` - Get user's cart
- `POST /cart` - Add item to cart
- `DELETE /cart/:id` - Remove item from cart

### Orders (Protected)
- `POST /orders` - Create order from cart
- `GET /orders` - Get user's orders

## Example Usage

### Register User
```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123","name":"John Doe"}'
```

### Login
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'
```

### Add to Cart (with token)
```bash
curl -X POST http://localhost:8080/cart \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{"product_id":1,"quantity":2}'
```

## What Makes This "Intermediate Level"

1. **Single Package**: Everything in one package, no complex module structure
2. **Direct Database Access**: No repository pattern or interfaces
3. **Simple Handlers**: Business logic directly in HTTP handlers
4. **Basic Models**: Simple GORM models without complex relationships
5. **Minimal Dependencies**: Only essential packages
6. **No Testing**: Focus on learning core concepts first
7. **No Advanced Features**: No caching, complex validation, etc.

This version is perfect for developers who want to understand the basics of building a REST API in Go without getting overwhelmed by enterprise patterns and advanced features.