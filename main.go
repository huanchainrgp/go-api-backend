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
	"gorm.io/gorm"
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

	// Auto-migrate the schema with error handling for existing data
	if err := migrateDatabase(db); err != nil {
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

// migrateDatabase handles database migration with proper error handling for existing data
func migrateDatabase(db *gorm.DB) error {
	// First, try to migrate without handling existing data
	if err := db.AutoMigrate(&models.User{}, &models.Asset{}, &models.Transaction{}); err != nil {
		log.Printf("Initial migration failed: %v", err)
		
		// Check if the error is related to username constraint
		if contains(err.Error(), "username") && contains(err.Error(), "null values") {
			log.Println("Detected username constraint issue. Attempting to fix existing data...")
			
			// First, add the username column as nullable
			log.Println("Adding username column as nullable first...")
			if err := db.Exec("ALTER TABLE users ADD COLUMN username text").Error; err != nil {
				log.Printf("Error adding username column: %v", err)
				// Column might already exist, continue
			}
			
			// Check for users with null usernames using raw SQL
			var nullUsernameCount int64
			if err := db.Raw("SELECT COUNT(*) FROM users WHERE username IS NULL OR username = ''").Scan(&nullUsernameCount).Error; err != nil {
				log.Printf("Error checking for null usernames: %v", err)
				return err
			}

			if nullUsernameCount > 0 {
				log.Printf("Found %d users with null/empty usernames. Updating them...", nullUsernameCount)
				
				// Update users with null usernames using raw SQL
				result := db.Exec(`
					UPDATE users 
					SET username = SUBSTRING(email FROM 1 FOR POSITION('@' IN email) - 1) || '_' || id::text
					WHERE username IS NULL OR username = ''
				`)
				
				if result.Error != nil {
					log.Printf("Error updating usernames: %v", result.Error)
					return result.Error
				}
				
				log.Printf("Updated %d users with generated usernames", result.RowsAffected)
			}
			
			// Now add the NOT NULL and UNIQUE constraints
			log.Println("Adding NOT NULL and UNIQUE constraints to username column...")
			if err := db.Exec("ALTER TABLE users ALTER COLUMN username SET NOT NULL").Error; err != nil {
				log.Printf("Error setting NOT NULL constraint: %v", err)
				return err
			}
			
			if err := db.Exec("CREATE UNIQUE INDEX idx_users_username ON users(username)").Error; err != nil {
				log.Printf("Error creating unique index: %v", err)
				// Index might already exist, continue
			}
			
			log.Println("Username column migration completed successfully!")
		} else if contains(err.Error(), "password") && contains(err.Error(), "null values") {
			log.Println("Detected password constraint issue. Attempting to fix existing data...")
			
			// First, add the password column as nullable
			log.Println("Adding password column as nullable first...")
			if err := db.Exec("ALTER TABLE users ADD COLUMN password text").Error; err != nil {
				log.Printf("Error adding password column: %v", err)
				// Column might already exist, continue
			}
			
			// Check for users with null passwords using raw SQL
			var nullPasswordCount int64
			if err := db.Raw("SELECT COUNT(*) FROM users WHERE password IS NULL OR password = ''").Scan(&nullPasswordCount).Error; err != nil {
				log.Printf("Error checking for null passwords: %v", err)
				return err
			}

			if nullPasswordCount > 0 {
				log.Printf("Found %d users with null/empty passwords. Updating them with default password...", nullPasswordCount)
				
				// Update users with null passwords using raw SQL
				// Note: In production, you should use proper password hashing
				result := db.Exec(`
					UPDATE users 
					SET password = 'default_password_' || id::text
					WHERE password IS NULL OR password = ''
				`)
				
				if result.Error != nil {
					log.Printf("Error updating passwords: %v", result.Error)
					return result.Error
				}
				
				log.Printf("Updated %d users with default passwords", result.RowsAffected)
				log.Println("WARNING: Users with default passwords should change their passwords immediately!")
			}
			
			// Now add the NOT NULL constraint
			log.Println("Adding NOT NULL constraint to password column...")
			if err := db.Exec("ALTER TABLE users ALTER COLUMN password SET NOT NULL").Error; err != nil {
				log.Printf("Error setting NOT NULL constraint: %v", err)
				return err
			}
			
			log.Println("Password column migration completed successfully!")
		} else {
			// If it's not a known constraint issue, return the original error
			return err
		}
	}
	
	log.Println("Database migration completed successfully!")
	return nil
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || (len(substr) > 0 && containsHelper(s, substr)))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
