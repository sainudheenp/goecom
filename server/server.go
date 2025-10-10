package server

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sainudheenp/goecom/config"
	store "github.com/sainudheenp/goecom/db"
	handler "github.com/sainudheenp/goecom/handlers"
	"github.com/sainudheenp/goecom/middleware"
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

	database, err := store.NewDB(cfg.Database.URL, logLevel)
	if err != nil {
		return nil, err
	}

	// Run migrations
	log.Println("Running database migrations...")
	if err := database.AutoMigrate(); err != nil {
		return nil, err
	}

	// Create router
	router := gin.New()

	s := &Server{
		router: router,
		config: cfg,
		db:     database,
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
	// Initialize handlers
	authHandler := handler.NewAuthHandler(s.db.DB, s.config.JWT.Secret, s.config.JWT.ExpiresHours, s.config.Security.BcryptCost)
	productHandler := handler.NewProductHandler(s.db.DB)

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
		protected.Use(middleware.AuthMiddleware(s.db.DB, s.config.JWT.Secret))
		{
			// User routes
			protected.GET("/me", authHandler.GetMe)
		}
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

// GetRouter returns the Gin router (for testing)
func (s *Server) GetRouter() *gin.Engine {
	return s.router
}
