package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID        uint           `json:"id" gorm:"primaryKey" example:"1"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null" example:"user@example.com"`
	Username  string         `json:"username" gorm:"uniqueIndex;not null" example:"johndoe"`
	Password  string         `json:"-" gorm:"not null"` // Hidden from JSON
	FirstName string         `json:"first_name" example:"John"`
	LastName  string         `json:"last_name" example:"Doe"`
	IsActive  bool           `json:"is_active" gorm:"default:true" example:"true"`
	CreatedAt time.Time      `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt time.Time      `json:"updated_at" example:"2023-01-01T00:00:00Z"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// Asset represents an asset in the system
type Asset struct {
	ID          uint           `json:"id" gorm:"primaryKey" example:"1"`
	Name        string         `json:"name" gorm:"not null" example:"Bitcoin"`
	Symbol      string         `json:"symbol" gorm:"uniqueIndex;not null" example:"BTC"`
	Type        string         `json:"type" gorm:"not null" example:"cryptocurrency"`
	Description string         `json:"description" example:"Digital currency"`
	Price       float64        `json:"price" gorm:"not null" example:"50000.00"`
	IsActive    bool           `json:"is_active" gorm:"default:true" example:"true"`
	CreatedAt   time.Time      `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt   time.Time      `json:"updated_at" example:"2023-01-01T00:00:00Z"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// Transaction represents a transaction between users and assets
type Transaction struct {
	ID          uint           `json:"id" gorm:"primaryKey" example:"1"`
	UserID      uint           `json:"user_id" gorm:"not null" example:"1"`
	AssetID     uint           `json:"asset_id" gorm:"not null" example:"1"`
	Type        string         `json:"type" gorm:"not null" example:"buy"` // buy, sell, transfer
	Amount      float64        `json:"amount" gorm:"not null" example:"0.5"`
	Price       float64        `json:"price" gorm:"not null" example:"50000.00"`
	TotalValue  float64        `json:"total_value" gorm:"not null" example:"25000.00"`
	Status      string         `json:"status" gorm:"default:'pending'" example:"completed"` // pending, completed, failed, cancelled
	Description string         `json:"description" example:"Buying Bitcoin"`
	CreatedAt   time.Time      `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt   time.Time      `json:"updated_at" example:"2023-01-01T00:00:00Z"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	
	// Relationships
	User  User  `json:"user" gorm:"foreignKey:UserID"`
	Asset Asset `json:"asset" gorm:"foreignKey:AssetID"`
}

// CreateUserRequest represents the request payload for creating a user
type CreateUserRequest struct {
	Email     string `json:"email" binding:"required,email" example:"user@example.com"`
	Username  string `json:"username" binding:"required,min=3,max=20" example:"johndoe"`
	Password  string `json:"password" binding:"required,min=6" example:"password123"`
	FirstName string `json:"first_name" example:"John"`
	LastName  string `json:"last_name" example:"Doe"`
}

// UpdateUserRequest represents the request payload for updating a user
type UpdateUserRequest struct {
	Email     string `json:"email" example:"user@example.com"`
	Username  string `json:"username" example:"johndoe"`
	FirstName string `json:"first_name" example:"John"`
	LastName  string `json:"last_name" example:"Doe"`
	IsActive  *bool  `json:"is_active" example:"true"`
}

// CreateAssetRequest represents the request payload for creating an asset
type CreateAssetRequest struct {
	Name        string  `json:"name" binding:"required" example:"Bitcoin"`
	Symbol      string  `json:"symbol" binding:"required" example:"BTC"`
	Type        string  `json:"type" binding:"required" example:"cryptocurrency"`
	Description string  `json:"description" example:"Digital currency"`
	Price       float64 `json:"price" binding:"required,min=0" example:"50000.00"`
}

// UpdateAssetRequest represents the request payload for updating an asset
type UpdateAssetRequest struct {
	Name        string  `json:"name" example:"Bitcoin"`
	Symbol      string  `json:"symbol" example:"BTC"`
	Type        string  `json:"type" example:"cryptocurrency"`
	Description string  `json:"description" example:"Digital currency"`
	Price       float64 `json:"price" example:"50000.00"`
	IsActive    *bool   `json:"is_active" example:"true"`
}

// CreateTransactionRequest represents the request payload for creating a transaction
type CreateTransactionRequest struct {
	AssetID     uint    `json:"asset_id" binding:"required" example:"1"`
	Type        string  `json:"type" binding:"required,oneof=buy sell transfer" example:"buy"`
	Amount      float64 `json:"amount" binding:"required,min=0" example:"0.5"`
	Price       float64 `json:"price" binding:"required,min=0" example:"50000.00"`
	Description string  `json:"description" example:"Buying Bitcoin"`
}

// UpdateTransactionRequest represents the request payload for updating a transaction
type UpdateTransactionRequest struct {
	Type        string  `json:"type" example:"buy"`
	Amount      float64 `json:"amount" example:"0.5"`
	Price       float64 `json:"price" example:"50000.00"`
	Status      string  `json:"status" example:"completed"`
	Description string  `json:"description" example:"Buying Bitcoin"`
}

// LoginRequest represents the request payload for user login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required" example:"password123"`
}

// RegisterRequest represents the request payload for user registration
type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email" example:"user@example.com"`
	Username  string `json:"username" binding:"required,min=3,max=20" example:"johndoe"`
	Password  string `json:"password" binding:"required,min=6" example:"password123"`
	FirstName string `json:"first_name" example:"John"`
	LastName  string `json:"last_name" example:"Doe"`
}

// AuthResponse represents the response payload for authentication
type AuthResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User  User   `json:"user"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error" example:"Invalid request"`
	Message string `json:"message" example:"The request body is invalid"`
}
