package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Auth middleware
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		claims, err := validateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Next()
	}
}

func setupRoutes(router *gin.Engine, db *Database) {
	// Public routes
	auth := router.Group("/auth")
	{
		auth.POST("/register", registerHandler(db))
		auth.POST("/login", loginHandler(db))
	}

	// Public product routes
	router.GET("/products", getProductsHandler(db))
	router.GET("/products/:id", getProductHandler(db))

	// Protected routes
	protected := router.Group("/")
	protected.Use(authMiddleware())
	{
		// Cart routes
		protected.GET("/cart", getCartHandler(db))
		protected.POST("/cart", addToCartHandler(db))
		protected.DELETE("/cart/:id", removeFromCartHandler(db))

		// Order routes
		protected.POST("/orders", createOrderHandler(db))
		protected.GET("/orders", getOrdersHandler(db))
	}
}

// Register handler
func registerHandler(db *Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Hash password
		hashedPassword, err := hashPassword(user.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		user.Password = hashedPassword

		// Create user
		if err := db.Create(&user).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "user_id": user.ID})
	}
}

// Login handler
func loginHandler(db *Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginData struct {
			Email    string `json:"email" binding:"required"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&loginData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var user User
		if err := db.Where("email = ?", loginData.Email).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		if !checkPasswordHash(loginData.Password, user.Password) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		token, err := generateToken(user.ID, user.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": token, "user": gin.H{"id": user.ID, "email": user.Email, "name": user.Name}})
	}
}

// Get products handler
func getProductsHandler(db *Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		var products []Product
		if err := db.Find(&products).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
			return
		}
		c.JSON(http.StatusOK, products)
	}
}

// Get single product handler
func getProductHandler(db *Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var product Product
		if err := db.First(&product, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		c.JSON(http.StatusOK, product)
	}
}

// Get cart handler
func getCartHandler(db *Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")
		var cartItems []CartItem
		if err := db.Preload("Product").Where("user_id = ?", userID).Find(&cartItems).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cart"})
			return
		}
		c.JSON(http.StatusOK, cartItems)
	}
}

// Add to cart handler
func addToCartHandler(db *Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")
		var cartData struct {
			ProductID uint `json:"product_id" binding:"required"`
			Quantity  int  `json:"quantity" binding:"required,min=1"`
		}

		if err := c.ShouldBindJSON(&cartData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check if product exists
		var product Product
		if err := db.First(&product, cartData.ProductID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}

		// Check if item already in cart
		var existingItem CartItem
		if err := db.Where("user_id = ? AND product_id = ?", userID, cartData.ProductID).First(&existingItem).Error; err == nil {
			// Update quantity
			existingItem.Quantity += cartData.Quantity
			db.Save(&existingItem)
			c.JSON(http.StatusOK, gin.H{"message": "Cart updated"})
			return
		}

		// Create new cart item
		cartItem := CartItem{
			UserID:    userID,
			ProductID: cartData.ProductID,
			Quantity:  cartData.Quantity,
		}

		if err := db.Create(&cartItem).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add to cart"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Added to cart"})
	}
}

// Remove from cart handler
func removeFromCartHandler(db *Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")
		cartItemID := c.Param("id")

		if err := db.Where("id = ? AND user_id = ?", cartItemID, userID).Delete(&CartItem{}).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cart item not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Item removed from cart"})
	}
}

// Create order handler
func createOrderHandler(db *Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")

		// Get cart items
		var cartItems []CartItem
		if err := db.Preload("Product").Where("user_id = ?", userID).Find(&cartItems).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cart"})
			return
		}

		if len(cartItems) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cart is empty"})
			return
		}

		// Calculate total
		var total float64
		for _, item := range cartItems {
			total += item.Product.Price * float64(item.Quantity)
		}

		// Create order
		order := Order{
			UserID: userID,
			Total:  total,
			Status: "pending",
		}

		if err := db.Create(&order).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
			return
		}

		// Clear cart
		db.Where("user_id = ?", userID).Delete(&CartItem{})

		c.JSON(http.StatusCreated, gin.H{"message": "Order created", "order_id": order.ID, "total": total})
	}
}

// Get orders handler
func getOrdersHandler(db *Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")
		var orders []Order
		if err := db.Where("user_id = ?", userID).Find(&orders).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
			return
		}
		c.JSON(http.StatusOK, orders)
	}
}