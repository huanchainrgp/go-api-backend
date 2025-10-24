package main

import (
	"log"
	"os"

	"go-api-test1/internal/config"
	"go-api-test1/internal/database"
	"go-api-test1/internal/models"

	"gorm.io/gorm"
)

func main() {
	// Load environment variables
	if err := os.Setenv("DATABASE_URL", os.Getenv("DATABASE_URL")); err != nil {
		log.Fatal("Failed to set DATABASE_URL:", err)
	}

	// Initialize configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.Initialize(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Check if users table exists and has data
	var userCount int64
	if err := db.Model(&models.User{}).Count(&userCount).Error; err != nil {
		log.Printf("Users table doesn't exist or error counting: %v", err)
		userCount = 0
	}

	if userCount > 0 {
		log.Printf("Found %d existing users. Checking for null usernames...", userCount)
		
		// Check for users with null usernames
		var nullUsernameCount int64
		if err := db.Model(&models.User{}).Where("username IS NULL OR username = ''").Count(&nullUsernameCount).Error; err != nil {
			log.Printf("Error checking for null usernames: %v", err)
		}

		if nullUsernameCount > 0 {
			log.Printf("Found %d users with null/empty usernames. Updating them...", nullUsernameCount)
			
			// Update users with null usernames to have a default username based on their email
			result := db.Model(&models.User{}).
				Where("username IS NULL OR username = ''").
				Update("username", gorm.Expr("SUBSTRING(email FROM 1 FOR POSITION('@' IN email) - 1) || '_' || id"))
			
			if result.Error != nil {
				log.Printf("Error updating usernames: %v", result.Error)
			} else {
				log.Printf("Updated %d users with generated usernames", result.RowsAffected)
			}
		}
	}

	// Now run the migration
	log.Println("Running database migration...")
	if err := db.AutoMigrate(&models.User{}, &models.Asset{}, &models.Transaction{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database migration completed successfully!")
}
