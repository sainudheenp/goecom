# E-Commerce API Examples

This document provides example requests for all API endpoints using `curl` and HTTPie.

## Prerequisites

Set your authentication token:

```bash
# After login, save your token
export TOKEN="your_jwt_token_here"

# Or on Windows CMD
set TOKEN=your_jwt_token_here

# Or on Windows PowerShell
$TOKEN="your_jwt_token_here"
```

## Authentication

### Register a New User

**cURL:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "newuser@example.com",
    "password": "securepass123",
    "full_name": "New User"
  }'
```

**HTTPie:**
```bash
http POST http://localhost:8080/api/v1/auth/register \
  email=newuser@example.com \
  password=securepass123 \
  full_name="New User"
```

### Login

**cURL:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "admin123"
  }'
```

**HTTPie:**
```bash
http POST http://localhost:8080/api/v1/auth/login \
  email=admin@example.com \
  password=admin123
```

### Get Current User Profile

**cURL:**
```bash
curl http://localhost:8080/api/v1/me \
  -H "Authorization: Bearer $TOKEN"
```

**HTTPie:**
```bash
http GET http://localhost:8080/api/v1/me \
  "Authorization: Bearer $TOKEN"
```

## Products

### List All Products

**cURL:**
```bash
curl http://localhost:8080/api/v1/products
```

**HTTPie:**
```bash
http GET http://localhost:8080/api/v1/products
```

### List Products with Filters

**cURL:**
```bash
curl "http://localhost:8080/api/v1/products?q=laptop&min_price=100000&max_price=200000&sort=price_asc&page=1&size=10"
```

**HTTPie:**
```bash
http GET http://localhost:8080/api/v1/products \
  q==laptop \
  min_price==100000 \
  max_price==200000 \
  sort==price_asc \
  page==1 \
  size==10
```

### Get Product by ID

**cURL:**
```bash
curl http://localhost:8080/api/v1/products/PRODUCT_ID
```

**HTTPie:**
```bash
http GET http://localhost:8080/api/v1/products/PRODUCT_ID
```

### Create Product (Admin Only)

**cURL:**
```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "sku": "NEW-PRODUCT-001",
    "name": "New Product",
    "description": "A great new product",
    "price_cents": 49900,
    "currency": "USD",
    "stock": 100,
    "images": ["https://example.com/image.jpg"]
  }'
```

**HTTPie:**
```bash
http POST http://localhost:8080/api/v1/products \
  "Authorization: Bearer $TOKEN" \
  sku=NEW-PRODUCT-001 \
  name="New Product" \
  description="A great new product" \
  price_cents:=49900 \
  currency=USD \
  stock:=100 \
  images:='["https://example.com/image.jpg"]'
```

### Update Product (Admin Only)

**cURL:**
```bash
curl -X PUT http://localhost:8080/api/v1/products/PRODUCT_ID \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated Product Name",
    "price_cents": 59900,
    "stock": 150
  }'
```

**HTTPie:**
```bash
http PUT http://localhost:8080/api/v1/products/PRODUCT_ID \
  "Authorization: Bearer $TOKEN" \
  name="Updated Product Name" \
  price_cents:=59900 \
  stock:=150
```

### Delete Product (Admin Only)

**cURL:**
```bash
curl -X DELETE http://localhost:8080/api/v1/products/PRODUCT_ID \
  -H "Authorization: Bearer $TOKEN"
```

**HTTPie:**
```bash
http DELETE http://localhost:8080/api/v1/products/PRODUCT_ID \
  "Authorization: Bearer $TOKEN"
```

## Shopping Cart

### Add Item to Cart

**cURL:**
```bash
curl -X POST http://localhost:8080/api/v1/cart \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": "PRODUCT_ID",
    "quantity": 2
  }'
```

**HTTPie:**
```bash
http POST http://localhost:8080/api/v1/cart \
  "Authorization: Bearer $TOKEN" \
  product_id=PRODUCT_ID \
  quantity:=2
```

### Get Cart

**cURL:**
```bash
curl http://localhost:8080/api/v1/cart \
  -H "Authorization: Bearer $TOKEN"
```

**HTTPie:**
```bash
http GET http://localhost:8080/api/v1/cart \
  "Authorization: Bearer $TOKEN"
```

### Remove Item from Cart

**cURL:**
```bash
curl -X DELETE http://localhost:8080/api/v1/cart/CART_ITEM_ID \
  -H "Authorization: Bearer $TOKEN"
```

**HTTPie:**
```bash
http DELETE http://localhost:8080/api/v1/cart/CART_ITEM_ID \
  "Authorization: Bearer $TOKEN"
```

## Orders

### Create Order

**cURL:**
```bash
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "shipping_address": {
      "line1": "123 Main St",
      "line2": "Apt 4B",
      "city": "Bengaluru",
      "state": "KA",
      "country": "IN",
      "postcode": "560001"
    }
  }'
```

**HTTPie:**
```bash
http POST http://localhost:8080/api/v1/orders \
  "Authorization: Bearer $TOKEN" \
  shipping_address:='{
    "line1": "123 Main St",
    "line2": "Apt 4B",
    "city": "Bengaluru",
    "state": "KA",
    "country": "IN",
    "postcode": "560001"
  }'
```

### List User Orders

**cURL:**
```bash
curl "http://localhost:8080/api/v1/orders?page=1&size=10" \
  -H "Authorization: Bearer $TOKEN"
```

**HTTPie:**
```bash
http GET http://localhost:8080/api/v1/orders \
  "Authorization: Bearer $TOKEN" \
  page==1 \
  size==10
```

### Get Order by ID

**cURL:**
```bash
curl http://localhost:8080/api/v1/orders/ORDER_ID \
  -H "Authorization: Bearer $TOKEN"
```

**HTTPie:**
```bash
http GET http://localhost:8080/api/v1/orders/ORDER_ID \
  "Authorization: Bearer $TOKEN"
```

## Payments

### Process Payment

**cURL:**
```bash
curl -X POST http://localhost:8080/api/v1/payments/charge \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "order_id": "ORDER_ID",
    "payment_method": "card",
    "payment_details": {
      "card_last4": "4242"
    }
  }'
```

**HTTPie:**
```bash
http POST http://localhost:8080/api/v1/payments/charge \
  "Authorization: Bearer $TOKEN" \
  order_id=ORDER_ID \
  payment_method=card \
  payment_details:='{"card_last4":"4242"}'
```

## Admin Endpoints

### List All Orders (Admin)

**cURL:**
```bash
curl "http://localhost:8080/api/v1/admin/orders?page=1&size=20" \
  -H "Authorization: Bearer $TOKEN"
```

**HTTPie:**
```bash
http GET http://localhost:8080/api/v1/admin/orders \
  "Authorization: Bearer $TOKEN" \
  page==1 \
  size==20
```

### Update Order Status (Admin)

**cURL:**
```bash
curl -X PATCH http://localhost:8080/api/v1/admin/orders/ORDER_ID \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "shipped"
  }'
```

**HTTPie:**
```bash
http PATCH http://localhost:8080/api/v1/admin/orders/ORDER_ID \
  "Authorization: Bearer $TOKEN" \
  status=shipped
```

## Health Check

**cURL:**
```bash
curl http://localhost:8080/health
```

**HTTPie:**
```bash
http GET http://localhost:8080/health
```
