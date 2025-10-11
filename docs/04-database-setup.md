# 4. Database Setup 🗄️

## PostgreSQL + GORM

This project uses **PostgreSQL** (like MySQL but better) and **GORM** (like Sequelize or TypeORM for Go).

## Database Connection

### Location: `db/database.go`

🔍 **Compare with Node.js:**
```javascript
// Node.js (Mongoose)
mongoose.connect(process.env.DATABASE_URL, {
    useNewUrlParser: true,
    useUnifiedTopology: true
});
```

```go
// Go (GORM)
func NewDB(databaseURL string, logLevel logger.LogLevel) (*DB, error) {
    gormConfig := &gorm.Config{
        Logger: logger.Default.LogMode(logLevel),
    }
    
    db, err := gorm.Open(postgres.Open(databaseURL), gormConfig)
    if err != nil {
        return nil, fmt.Errorf("failed to connect: %w", err)
    }
    
    // Connection pool settings
    sqlDB, _ := db.DB()
    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetMaxOpenConns(100)
    sqlDB.SetConnMaxLifetime(time.Hour)
    
    return &DB{db}, nil
}
```

### Connection String Format

```
postgresql://username:password@host:port/database?sslmode=disable
```

**Example:**
```
postgresql://postgres:mypassword@localhost:5432/ecommerce?sslmode=disable
```

**For Supabase/Cloud:**
```
postgresql://user:pass@aws-0-region.pooler.supabase.com:5432/postgres
```

## Models (Data Schema)

### Location: `models/models.go`

### User Model

```go
type User struct {
    ID           uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
    Email        string    `gorm:"uniqueIndex;not null" json:"email"`
    PasswordHash string    `gorm:"not null" json:"-"`
    FullName     string    `json:"full_name"`
    Role         string    `gorm:"not null;default:'user'" json:"role"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}
```

### GORM Tags Explained

```go
type User struct {
    // UUID as primary key
    ID uuid.UUID `gorm:"type:uuid;primary_key;"`
    
    // Unique email with index
    Email string `gorm:"uniqueIndex;not null"`
    
    // Not included in JSON response (password security)
    PasswordHash string `gorm:"not null" json:"-"`
    
    // Default value
    Role string `gorm:"default:'user'"`
    
    // Auto timestamps
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

### Common GORM Tags

| Tag | Purpose | Example |
|-----|---------|---------|
| `primary_key` | Primary key | `gorm:"primary_key"` |
| `not null` | Required field | `gorm:"not null"` |
| `unique` | Unique value | `gorm:"unique"` |
| `uniqueIndex` | Unique index | `gorm:"uniqueIndex"` |
| `default:` | Default value | `gorm:"default:'user'"` |
| `size:` | Column size | `gorm:"size:255"` |
| `type:` | Column type | `gorm:"type:uuid"` |

### JSON Tags

```go
// Field name in JSON
Name string `json:"name"`

// Different name in JSON
FullName string `json:"full_name"`

// Omit from JSON response
Password string `json:"-"`

// Omit if empty
Phone string `json:"phone,omitempty"`
```

### Product Model

```go
type Product struct {
    ID          uuid.UUID       `gorm:"type:uuid;primary_key;"`
    SKU         string          `gorm:"uniqueIndex;not null"`
    Name        string          `gorm:"not null"`
    Description string
    PriceCents  int             `gorm:"not null"`
    Currency    string          `gorm:"default:'USD'"`
    Stock       int             `gorm:"default:0"`
    Images      JSONStringSlice `gorm:"type:jsonb"`  // PostgreSQL JSON
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

### Auto UUID Generation

```go
// BeforeCreate hook
func (u *User) BeforeCreate(tx *gorm.DB) error {
    if u.ID == uuid.Nil {
        u.ID = uuid.New()
    }
    return nil
}
```

💡 **Tip:** This automatically generates UUID before saving

## Migrations

### What are Migrations?

Migrations are like **version control for your database**. They help you:
- Track database changes
- Share schema with team
- Deploy to production safely
- Rollback if needed

🔍 **Compare with Node.js:**
```javascript
// Sequelize migration
module.exports = {
  up: async (queryInterface, Sequelize) => {
    await queryInterface.createTable('users', {
      id: { type: Sequelize.UUID, primaryKey: true },
      email: { type: Sequelize.STRING, unique: true }
    });
  },
  down: async (queryInterface) => {
    await queryInterface.dropTable('users');
  }
};
```

### Migration Files

Located in `migrations/` folder:

```
migrations/
├── 001_create_users_table.up.sql       # Create
├── 001_create_users_table.down.sql     # Rollback
├── 002_create_products_table.up.sql
├── 002_create_products_table.down.sql
└── ...
```

### Example Migration

**`001_create_users_table.up.sql`:**
```sql
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255),
    role VARCHAR(50) NOT NULL DEFAULT 'user',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
```

**`001_create_users_table.down.sql`:**
```sql
DROP INDEX IF EXISTS idx_users_role;
DROP INDEX IF EXISTS idx_users_email;
DROP TABLE IF EXISTS users;
```

### Running Migrations

```bash
# Apply all migrations
./scripts/migrate.sh up

# Rollback all migrations
./scripts/migrate.sh down
```

### Migration Script (`scripts/migrate.sh`)

```bash
#!/bin/bash
DIRECTION=${1:-up}

# Load .env file
export $(grep -v '^#' .env | xargs)

# Run migrations
for file in migrations/*."$DIRECTION".sql; do
    echo "Running: $(basename "$file")"
    psql "$DATABASE_URL" -f "$file"
done
```

## Database Operations

### Create

```go
// Single record
user := &User{
    Email:        "test@example.com",
    PasswordHash: "hashed",
    FullName:     "Test User",
}
db.Create(user)

// Multiple records
users := []User{
    {Email: "user1@example.com"},
    {Email: "user2@example.com"},
}
db.Create(&users)
```

### Read

```go
// Find by ID
var user User
db.First(&user, "id = ?", userId)

// Find by email
db.Where("email = ?", "test@example.com").First(&user)

// Find all with condition
var users []User
db.Where("role = ?", "admin").Find(&users)

// Find with pagination
db.Limit(10).Offset(20).Find(&users)

// Count
var count int64
db.Model(&User{}).Where("role = ?", "admin").Count(&count)
```

### Update

```go
// Update single field
db.Model(&user).Update("full_name", "New Name")

// Update multiple fields
db.Model(&user).Updates(User{
    FullName: "New Name",
    Role:     "admin",
})

// Update with map
db.Model(&user).Updates(map[string]interface{}{
    "full_name": "New Name",
    "role":      "admin",
})
```

### Delete

```go
// Soft delete (sets deleted_at)
db.Delete(&user)

// Permanent delete
db.Unscoped().Delete(&user)

// Delete by condition
db.Where("role = ?", "guest").Delete(&User{})
```

### Advanced Queries

**Joins:**
```go
type Result struct {
    UserName    string
    ProductName string
}

var results []Result
db.Table("orders").
    Select("users.full_name as user_name, products.name as product_name").
    Joins("LEFT JOIN users ON orders.user_id = users.id").
    Joins("LEFT JOIN products ON orders.product_id = products.id").
    Scan(&results)
```

**Search:**
```go
// LIKE query
db.Where("name LIKE ?", "%laptop%").Find(&products)

// ILIKE (case-insensitive PostgreSQL)
db.Where("name ILIKE ?", "%laptop%").Find(&products)

// Multiple conditions
db.Where("price_cents > ? AND stock > 0", 50000).Find(&products)
```

**Order By:**
```go
// Ascending
db.Order("created_at asc").Find(&products)

// Descending
db.Order("price_cents desc").Find(&products)

// Multiple
db.Order("price_cents desc, name asc").Find(&products)
```

## Seeding Data

### Location: `scripts/seed.sh`

```bash
#!/bin/bash

psql "$DATABASE_URL" <<EOF
-- Insert users
INSERT INTO users (id, email, password_hash, full_name, role)
VALUES 
    (gen_random_uuid(), 'admin@example.com', '\$2a\$10\$...', 'Admin', 'admin'),
    (gen_random_uuid(), 'user@example.com', '\$2a\$10\$...', 'User', 'user')
ON CONFLICT (email) DO NOTHING;

-- Insert products
INSERT INTO products (id, sku, name, price_cents, currency, stock)
VALUES
    (gen_random_uuid(), 'LAPTOP-001', 'MacBook Pro', 199900, 'USD', 10),
    (gen_random_uuid(), 'PHONE-001', 'iPhone 15', 99900, 'USD', 25)
ON CONFLICT (sku) DO NOTHING;
EOF
```

Run with:
```bash
./scripts/seed.sh
```

## Connection Pool

### Why Connection Pooling?

Instead of creating a new database connection for each request:
- ✅ Reuse existing connections
- ✅ Faster response times
- ✅ Handle more concurrent requests

```go
sqlDB, _ := db.DB()

// Maximum idle connections
sqlDB.SetMaxIdleConns(10)

// Maximum open connections
sqlDB.SetMaxOpenConns(100)

// Connection lifetime
sqlDB.SetConnMaxLifetime(time.Hour)
```

## Transactions

For operations that must all succeed or all fail:

```go
// Begin transaction
tx := db.Begin()

// Create user
if err := tx.Create(&user).Error; err != nil {
    tx.Rollback()
    return err
}

// Create profile
if err := tx.Create(&profile).Error; err != nil {
    tx.Rollback()
    return err
}

// Commit
tx.Commit()
```

## Database Best Practices

### 1. Always Check Errors

```go
// ❌ Bad
db.Create(&user)

// ✅ Good
if err := db.Create(&user).Error; err != nil {
    log.Printf("Failed to create user: %v", err)
    return err
}
```

### 2. Use Prepared Statements

```go
// ❌ Bad (SQL injection risk)
db.Where(fmt.Sprintf("email = '%s'", email)).First(&user)

// ✅ Good
db.Where("email = ?", email).First(&user)
```

### 3. Use Indexes

```sql
-- Speed up queries
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_products_sku ON products(sku);
```

### 4. Pagination

```go
// Don't load all records
page := 1
pageSize := 20
offset := (page - 1) * pageSize

db.Limit(pageSize).Offset(offset).Find(&products)
```

### 5. Select Specific Fields

```go
// ❌ Bad (loads everything)
db.Find(&users)

// ✅ Good (only needed fields)
db.Select("id, email, full_name").Find(&users)
```

## Common Issues & Solutions

### Issue: "relation does not exist"
**Solution:** Run migrations
```bash
./scripts/migrate.sh up
```

### Issue: "duplicate key value"
**Solution:** Record already exists, use `ON CONFLICT` or check before insert

### Issue: "connection refused"
**Solution:** PostgreSQL not running
```bash
sudo systemctl start postgresql
```

### Issue: "too many connections"
**Solution:** Reduce connection pool size
```go
sqlDB.SetMaxOpenConns(50)  // Lower number
```

## Database Tools

### psql (PostgreSQL CLI)

```bash
# Connect to database
psql postgresql://user:pass@localhost:5432/dbname

# List databases
\l

# List tables
\dt

# Describe table
\d users

# Run query
SELECT * FROM users LIMIT 5;

# Exit
\q
```

### GUI Tools

- **pgAdmin** - Full-featured PostgreSQL admin
- **DBeaver** - Universal database tool
- **TablePlus** - Modern, native GUI

## Practice Exercises

1. **Add a new field** to User model (e.g., phone number)
2. **Create a migration** for the new field
3. **Query users** by role
4. **Update a user's** email
5. **Delete users** created before a certain date

---

**Next:** [Authentication →](./05-authentication.md)
