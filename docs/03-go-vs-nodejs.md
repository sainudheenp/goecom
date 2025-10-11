# 3. Go vs Node.js - Key Differences 🔄

## Paradigm Shift

Coming from Node.js, you'll notice some fundamental differences in how Go approaches problems.

## Language Differences

### 1. Type System

**Node.js (JavaScript):**
```javascript
// Dynamic typing
let user = { name: "John", age: 30 };
user.age = "thirty"; // No error! 😱
```

**Go:**
```go
// Static typing
type User struct {
    Name string
    Age  int
}

user := User{Name: "John", Age: 30}
user.Age = "thirty" // ❌ Compile error!
```

💡 **Tip:** Go catches type errors at **compile time**, not runtime

### 2. Error Handling

**Node.js:**
```javascript
// Try-catch
try {
    const user = await User.findById(id);
    return user;
} catch (error) {
    console.error(error);
    throw error;
}

// Or callbacks
fs.readFile('file.txt', (err, data) => {
    if (err) throw err;
    console.log(data);
});
```

**Go:**
```go
// Explicit error return
user, err := getUserById(id)
if err != nil {
    log.Println(err)
    return nil, err
}
return user, nil
```

📌 **Remember:** 
- Go doesn't have `try-catch`
- Errors are **values**, not exceptions
- Always check `if err != nil`

### 3. Async/Await vs Goroutines

**Node.js:**
```javascript
// Async/await
async function getUser(id) {
    const user = await db.users.findOne({ id });
    const posts = await db.posts.find({ userId: id });
    return { user, posts };
}

// Promise.all for parallel
const [user, posts] = await Promise.all([
    db.users.findOne({ id }),
    db.posts.find({ userId: id })
]);
```

**Go:**
```go
// Synchronous (looks blocking but isn't)
func getUser(id string) (*User, error) {
    user, err := db.FindUser(id)
    if err != nil {
        return nil, err
    }
    posts, err := db.FindPosts(id)
    if err != nil {
        return nil, err
    }
    return user, nil
}

// Goroutines for parallel (like workers)
func getUserParallel(id string) (*User, error) {
    userChan := make(chan *User)
    postsChan := make(chan []Post)
    
    go func() {
        user, _ := db.FindUser(id)
        userChan <- user
    }()
    
    go func() {
        posts, _ := db.FindPosts(id)
        postsChan <- posts
    }()
    
    user := <-userChan
    posts := <-postsChan
    return user, nil
}
```

💡 **Tip:** Go is **concurrent** by default. Goroutines are lightweight threads.

### 4. null vs nil

**Node.js:**
```javascript
let user = null;
if (user === null) {
    console.log("No user");
}

// Also undefined
let count;  // undefined
```

**Go:**
```go
var user *User  // nil (pointer)
if user == nil {
    fmt.Println("No user")
}

// Go has zero values, not undefined
var count int  // 0
var str string // ""
var slice []int // nil
```

### 5. Functions

**Node.js:**
```javascript
// Arrow function
const add = (a, b) => a + b;

// Regular function
function multiply(a, b) {
    return a * b;
}

// Default parameters
function greet(name = "World") {
    return `Hello ${name}`;
}
```

**Go:**
```go
// Function declaration
func add(a int, b int) int {
    return a + b
}

// Shortened parameter syntax
func multiply(a, b int) int {
    return a * b
}

// No default parameters! Use overloading pattern
func greet(name string) string {
    if name == "" {
        name = "World"
    }
    return "Hello " + name
}

// Multiple return values
func divide(a, b int) (int, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}
```

## Web Framework Differences

### Express vs Gin

**Node.js (Express):**
```javascript
const express = require('express');
const app = express();

// Middleware
app.use(express.json());
app.use(cors());

// Routes
app.get('/users/:id', (req, res) => {
    const id = req.params.id;
    const user = await User.findById(id);
    res.json(user);
});

app.post('/users', (req, res) => {
    const user = req.body;
    await User.create(user);
    res.status(201).json(user);
});

app.listen(3000);
```

**Go (Gin):**
```go
import "github.com/gin-gonic/gin"

router := gin.Default()  // Includes logger & recovery

// Middleware
router.Use(cors.Default())

// Routes
router.GET("/users/:id", func(c *gin.Context) {
    id := c.Param("id")
    user, err := getUserById(id)
    if err != nil {
        c.JSON(404, gin.H{"error": "not found"})
        return
    }
    c.JSON(200, user)
})

router.POST("/users", func(c *gin.Context) {
    var user User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    createUser(&user)
    c.JSON(201, user)
})

router.Run(":8080")
```

### Key Differences:

| Feature | Express | Gin |
|---------|---------|-----|
| JSON Parsing | `express.json()` | `c.ShouldBindJSON()` |
| Response | `res.json()` | `c.JSON()` |
| Status Code | `res.status(200)` | `c.JSON(200, ...)` |
| Params | `req.params.id` | `c.Param("id")` |
| Query | `req.query.page` | `c.Query("page")` |
| Body | `req.body` | Bind to struct |

## Database - Mongoose vs GORM

### Defining Models

**Node.js (Mongoose):**
```javascript
const userSchema = new mongoose.Schema({
    email: { 
        type: String, 
        required: true, 
        unique: true 
    },
    password: { 
        type: String, 
        required: true 
    },
    createdAt: { 
        type: Date, 
        default: Date.now 
    }
});

const User = mongoose.model('User', userSchema);
```

**Go (GORM):**
```go
type User struct {
    ID        uuid.UUID `gorm:"type:uuid;primary_key"`
    Email     string    `gorm:"uniqueIndex;not null"`
    Password  string    `gorm:"not null"`
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

### CRUD Operations

**Create:**
```javascript
// Node.js
const user = await User.create({
    email: "test@example.com",
    password: "hashed"
});
```

```go
// Go
user := &User{
    Email:    "test@example.com",
    Password: "hashed",
}
db.Create(user)
```

**Read:**
```javascript
// Node.js
const user = await User.findOne({ email: "test@example.com" });
const users = await User.find({ role: "admin" });
```

```go
// Go
var user User
db.Where("email = ?", "test@example.com").First(&user)

var users []User
db.Where("role = ?", "admin").Find(&users)
```

**Update:**
```javascript
// Node.js
await User.updateOne(
    { _id: userId },
    { $set: { name: "New Name" } }
);
```

```go
// Go
db.Model(&User{}).
    Where("id = ?", userId).
    Update("name", "New Name")
```

**Delete:**
```javascript
// Node.js
await User.deleteOne({ _id: userId });
```

```go
// Go
db.Delete(&User{}, userId)
```

## JSON Handling

**Node.js:**
```javascript
// Automatic
const data = { name: "John", age: 30 };
res.json(data);  // Automatically serialized

const parsed = JSON.parse(jsonString);
```

**Go:**
```go
// Explicit with struct tags
type User struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

// Marshal (struct to JSON)
data := User{Name: "John", Age: 30}
jsonBytes, _ := json.Marshal(data)

// Unmarshal (JSON to struct)
var user User
json.Unmarshal(jsonBytes, &user)
```

## Environment Variables

**Node.js:**
```javascript
require('dotenv').config();
const port = process.env.PORT || 3000;
```

**Go:**
```go
import (
    "os"
    "github.com/joho/godotenv"
)

godotenv.Load()  // Load .env
port := os.Getenv("PORT")
if port == "" {
    port = "8080"
}
```

## Package Management

**Node.js:**
```bash
npm install express
npm install --save-dev nodemon
```

**Go:**
```bash
go get github.com/gin-gonic/gin
go mod tidy  # Clean up
```

## Testing

**Node.js (Jest):**
```javascript
describe('User API', () => {
    test('should create user', async () => {
        const response = await request(app)
            .post('/api/users')
            .send({ email: 'test@example.com' });
        
        expect(response.status).toBe(201);
        expect(response.body.email).toBe('test@example.com');
    });
});
```

**Go (testing package):**
```go
func TestCreateUser(t *testing.T) {
    user := &User{Email: "test@example.com"}
    err := createUser(user)
    
    if err != nil {
        t.Errorf("Expected no error, got %v", err)
    }
    
    if user.Email != "test@example.com" {
        t.Errorf("Expected email test@example.com, got %s", user.Email)
    }
}
```

## Common Gotchas for Node.js Developers

### 1. Capitalization Matters

```go
// Public (exported) - starts with capital
func CreateUser() { }

// Private (unexported) - starts with lowercase
func hashPassword() { }
```

### 2. Pointers

```javascript
// Node.js - everything is reference by default
let user1 = { name: "John" };
let user2 = user1;
user2.name = "Jane";
console.log(user1.name);  // "Jane"
```

```go
// Go - need explicit pointers
user1 := User{Name: "John"}
user2 := user1  // Copy!
user2.Name = "Jane"
fmt.Println(user1.Name)  // "John"

// Use pointer for reference
user2 := &user1  // Pointer
user2.Name = "Jane"
fmt.Println(user1.Name)  // "Jane"
```

### 3. No Classes

```javascript
// Node.js - classes
class User {
    constructor(name) {
        this.name = name;
    }
    
    greet() {
        return `Hello ${this.name}`;
    }
}
```

```go
// Go - structs + methods
type User struct {
    Name string
}

func (u *User) Greet() string {
    return "Hello " + u.Name
}
```

### 4. Slice vs Array

```javascript
// Node.js - arrays are dynamic
let numbers = [1, 2, 3];
numbers.push(4);  // [1, 2, 3, 4]
```

```go
// Go - arrays are fixed size
var numbers [3]int = [3]int{1, 2, 3}

// Use slices for dynamic arrays
numbers := []int{1, 2, 3}
numbers = append(numbers, 4)  // [1, 2, 3, 4]
```

### 5. Maps

```javascript
// Node.js - objects as maps
let user = { name: "John", age: 30 };
user.email = "john@example.com";
```

```go
// Go - explicit maps
user := map[string]interface{}{
    "name": "John",
    "age":  30,
}
user["email"] = "john@example.com"

// Or use structs (preferred)
type User struct {
    Name  string
    Age   int
    Email string
}
```

## Performance Comparison

| Aspect | Node.js | Go |
|--------|---------|-----|
| Concurrency | Event loop, async/await | Goroutines, channels |
| Memory | Higher (V8 engine) | Lower (compiled) |
| Speed | Fast | Faster (compiled) |
| Startup | Slower | Instant |
| CPU-bound | Single-threaded | Multi-threaded |

## Mental Model Shift

### Node.js Developer Mindset:
- "Everything is async"
- "Callbacks and promises"
- "Dynamic and flexible"

### Go Developer Mindset:
- "Explicit is better than implicit"
- "Errors are values"
- "Compile-time safety"

## Quick Reference Cheat Sheet

| Task | Node.js | Go |
|------|---------|-----|
| Import | `require()` / `import` | `import` |
| Export | `module.exports` | Capital letter |
| Async | `async/await` | Goroutines |
| Error | `try-catch` | `if err != nil` |
| Null | `null` / `undefined` | `nil` |
| String interpolation | `` `Hello ${name}` `` | `"Hello " + name` |
| Array | `[]` | `[]Type` |
| Object | `{}` | `map` or `struct` |

## Tips for Transition

1. ✅ Embrace explicit error handling
2. ✅ Think in types and structs
3. ✅ Use `go fmt` religiously
4. ✅ Read Go standard library docs
5. ✅ Don't fight the language

---

**Next:** [Database Setup →](./04-database-setup.md)
