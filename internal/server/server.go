package server

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sainudheenp/goecom/internal/config"
	"github.com/sainudheenp/goecom/internal/handler"
	"github.com/sainudheenp/goecom/internal/middleware"
	"github.com/sainudheenp/goecom/internal/service"
	"github.com/sainudheenp/goecom/internal/store"
	"gorm.io/gorm/logger"
)

// Server represents the HTTP server
type Server struct {
	router *gin.Engine
	config *config.Config
	db     *store.DB
}

// NewServer creates a new server instance
func NewServer(cfg *config.Config) (*Server, error) {
	// Set Gin mode
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize database
	logLevel := logger.Info
	if cfg.IsDevelopment() {
		logLevel = logger.Info
	}

	db, err := store.NewDB(cfg.Database.URL, logLevel)
	if err != nil {
		return nil, err
	}

	// Run migrations
	log.Println("Running database migrations...")
	if err := db.AutoMigrate(); err != nil {
		return nil, err
	}

	// Create router
	router := gin.New()

	s := &Server{
		router: router,
		config: cfg,
		db:     db,
	}

	s.setupMiddleware()
	s.setupRoutes()

	return s, nil
}

// setupMiddleware configures middleware
func (s *Server) setupMiddleware() {
	// Recovery middleware
	s.router.Use(middleware.Recovery())

	// Request ID middleware
	s.router.Use(middleware.RequestID())

	// Logger middleware
	s.router.Use(middleware.Logger())

	// CORS middleware
	corsConfig := cors.Config{
		AllowOrigins:     s.config.CORS.Origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Request-ID"},
		ExposeHeaders:    []string{"X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	s.router.Use(cors.New(corsConfig))

	// Rate limiting middleware
	rateLimiter := middleware.NewRateLimiter(
		s.config.RateLimit.Requests,
		s.config.RateLimit.WindowMinutes,
	)
	s.router.Use(rateLimiter.Middleware())
}

// setupRoutes configures routes
func (s *Server) setupRoutes() {
	// Initialize repositories
	userRepo := store.NewUserRepository(s.db)
	productRepo := store.NewProductRepository(s.db)
	cartRepo := store.NewCartRepository(s.db)
	orderRepo := store.NewOrderRepository(s.db)

	// Initialize services
	authService := service.NewAuthService(
		userRepo,
		s.config.JWT.Secret,
		s.config.JWT.ExpiresHours,
		s.config.Security.BcryptCost,
	)
	productService := service.NewProductService(productRepo)
	cartService := service.NewCartService(cartRepo, productRepo)
	orderService := service.NewOrderService(orderRepo, cartRepo, productRepo, s.db)
	paymentService := service.NewPaymentService(orderRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	productHandler := handler.NewProductHandler(productService)
	cartHandler := handler.NewCartHandler(cartService)
	orderHandler := handler.NewOrderHandler(orderService)
	paymentHandler := handler.NewPaymentHandler(paymentService)

	// Health check
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"time":   time.Now().UTC(),
		})
	})

	// API v1 routes
	v1 := s.router.Group("/api/v1")
	{
		// Public routes
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Public product routes
		v1.GET("/products", productHandler.ListProducts)
		v1.GET("/products/:id", productHandler.GetProduct)

		// Protected routes
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(authService))
		{
			// User routes
			protected.GET("/me", authHandler.GetMe)

			// Cart routes
			protected.POST("/cart", cartHandler.AddToCart)
			protected.GET("/cart", cartHandler.GetCart)
			protected.DELETE("/cart/:item_id", cartHandler.RemoveFromCart)

			// Order routes
			protected.POST("/orders", orderHandler.CreateOrder)
			protected.GET("/orders", orderHandler.ListUserOrders)
			protected.GET("/orders/:id", orderHandler.GetOrder)

			// Payment routes
			protected.POST("/payments/charge", paymentHandler.ProcessCharge)
		}

		// Admin routes
		admin := v1.Group("/admin")
		admin.Use(middleware.AuthMiddleware(authService))
		admin.Use(middleware.RequireRole("admin"))
		{
			// Admin product routes
			admin.POST("/products", productHandler.CreateProduct)
			admin.PUT("/products/:id", productHandler.UpdateProduct)
			admin.DELETE("/products/:id", productHandler.DeleteProduct)
			admin.POST("/products/bulk", productHandler.BulkImportProducts)

			// Admin order routes
			admin.GET("/orders", orderHandler.ListAllOrders)
			admin.PATCH("/orders/:id", orderHandler.UpdateOrderStatus)
		}

		// Admin product routes at root level (alternative)
		v1.POST("/products", middleware.AuthMiddleware(authService), middleware.RequireRole("admin"), productHandler.CreateProduct)
		v1.PUT("/products/:id", middleware.AuthMiddleware(authService), middleware.RequireRole("admin"), productHandler.UpdateProduct)
		v1.DELETE("/products/:id", middleware.AuthMiddleware(authService), middleware.RequireRole("admin"), productHandler.DeleteProduct)
	}
}

// Run starts the HTTP server
func (s *Server) Run() error {
	addr := ":" + s.config.Server.Port
	log.Printf("Starting server on %s", addr)
	return s.router.Run(addr)
}

// Close closes the server and its resources
func (s *Server) Close() error {
	return s.db.Close()
}
