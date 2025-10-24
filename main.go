package main

import (
	"log"
	"os"

	"go-api-test1/docs"
	"go-api-test1/internal/config"
	"go-api-test1/internal/database"
	"go-api-test1/internal/handlers"
	"go-api-test1/internal/middleware"
	"go-api-test1/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/files"
)

// @title           Go API Test1
// @version         1.0
// @description     A backend API with Users, Assets, and Transactions
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize Swagger docs
	docs.SwaggerInfo.Title = "Go API Test1"
	docs.SwaggerInfo.Description = "A backend API with Users, Assets, and Transactions"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http"}

	// Initialize configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.Initialize(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate the schema
	if err := db.AutoMigrate(&models.User{}, &models.Asset{}, &models.Transaction{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Initialize Gin router
	router := gin.Default()

	// Add CORS middleware
	router.Use(middleware.CORS())

	// Initialize handlers
	userHandler := handlers.NewUserHandler(db)
	assetHandler := handlers.NewAssetHandler(db)
	transactionHandler := handlers.NewTransactionHandler(db)
	authHandler := handlers.NewAuthHandler(db)

	// API routes
	v1 := router.Group("/api/v1")
	{
		// Authentication routes
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Protected routes
		protected := v1.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			// User routes
			users := protected.Group("/users")
			{
				users.GET("", userHandler.GetUsers)
				users.GET("/:id", userHandler.GetUser)
				users.PUT("/:id", userHandler.UpdateUser)
				users.DELETE("/:id", userHandler.DeleteUser)
			}

			// Asset routes
			assets := protected.Group("/assets")
			{
				assets.GET("", assetHandler.GetAssets)
				assets.GET("/:id", assetHandler.GetAsset)
				assets.POST("", assetHandler.CreateAsset)
				assets.PUT("/:id", assetHandler.UpdateAsset)
				assets.DELETE("/:id", assetHandler.DeleteAsset)
			}

			// Transaction routes
			transactions := protected.Group("/transactions")
			{
				transactions.GET("", transactionHandler.GetTransactions)
				transactions.GET("/:id", transactionHandler.GetTransaction)
				transactions.POST("", transactionHandler.CreateTransaction)
				transactions.PUT("/:id", transactionHandler.UpdateTransaction)
				transactions.DELETE("/:id", transactionHandler.DeleteTransaction)
			}
		}
	}

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Printf("Swagger UI available at: http://localhost:%s/swagger/index.html", port)
	
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
