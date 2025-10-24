package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-api-test1/internal/config"
	"go-api-test1/internal/database"
	"go-api-test1/internal/handlers"
	"go-api-test1/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&models.User{}, &models.Asset{}, &models.Transaction{})
	return db
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	db := setupTestDB()
	
	router := gin.New()
	
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
		
		// User routes
		users := v1.Group("/users")
		{
			users.GET("", userHandler.GetUsers)
			users.GET("/:id", userHandler.GetUser)
			users.PUT("/:id", userHandler.UpdateUser)
			users.DELETE("/:id", userHandler.DeleteUser)
		}
		
		// Asset routes
		assets := v1.Group("/assets")
		{
			assets.GET("", assetHandler.GetAssets)
			assets.GET("/:id", assetHandler.GetAsset)
			assets.POST("", assetHandler.CreateAsset)
			assets.PUT("/:id", assetHandler.UpdateAsset)
			assets.DELETE("/:id", assetHandler.DeleteAsset)
		}
		
		// Transaction routes
		transactions := v1.Group("/transactions")
		{
			transactions.GET("", transactionHandler.GetTransactions)
			transactions.GET("/:id", transactionHandler.GetTransaction)
			transactions.POST("", transactionHandler.CreateTransaction)
			transactions.PUT("/:id", transactionHandler.UpdateTransaction)
			transactions.DELETE("/:id", transactionHandler.DeleteTransaction)
		}
	}
	
	return router
}

func TestUserRegistration(t *testing.T) {
	router := setupTestRouter()
	
	userData := models.RegisterRequest{
		Email:     "test@example.com",
		Username:  "testuser",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
	}
	
	jsonData, _ := json.Marshal(userData)
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusCreated, w.Code)
	
	var response models.AuthResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.NotEmpty(t, response.Token)
	assert.Equal(t, userData.Email, response.User.Email)
}

func TestAssetCreation(t *testing.T) {
	router := setupTestRouter()
	
	assetData := models.CreateAssetRequest{
		Name:        "Bitcoin",
		Symbol:      "BTC",
		Type:        "cryptocurrency",
		Description: "Digital currency",
		Price:       50000.00,
	}
	
	jsonData, _ := json.Marshal(assetData)
	req, _ := http.NewRequest("POST", "/api/v1/assets", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusCreated, w.Code)
	
	var asset models.Asset
	json.Unmarshal(w.Body.Bytes(), &asset)
	assert.Equal(t, assetData.Name, asset.Name)
	assert.Equal(t, assetData.Symbol, asset.Symbol)
}

func TestGetAssets(t *testing.T) {
	router := setupTestRouter()
	
	req, _ := http.NewRequest("GET", "/api/v1/assets", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var assets []models.Asset
	json.Unmarshal(w.Body.Bytes(), &assets)
	assert.IsType(t, []models.Asset{}, assets)
}
