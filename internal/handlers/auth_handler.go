package handlers

import (
	"log"
	"net/http"
	"time"

	"go-api-test1/internal/config"
	"go-api-test1/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	db *gorm.DB
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{db: db}
}

// Register registers a new user
// @Summary      Register a new user
// @Description  Register a new user account
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user body      models.RegisterRequest  true  "User registration data"
// @Success      201  {object}  models.AuthResponse
// @Failure      400  {object}  models.ErrorResponse
// @Failure      409  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	log.Printf("Auth: Registration attempt from %s", c.ClientIP())
	
	var registerReq models.RegisterRequest
	if err := c.ShouldBindJSON(&registerReq); err != nil {
		log.Printf("Auth: Invalid registration request from %s: %v", c.ClientIP(), err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

	log.Printf("Auth: Processing registration for email: %s, username: %s", registerReq.Email, registerReq.Username)

	// Check if user already exists
	var existingUser models.User
	if err := h.db.Where("email = ? OR username = ?", registerReq.Email, registerReq.Username).First(&existingUser).Error; err == nil {
		log.Printf("Auth: Registration failed - user already exists with email: %s or username: %s", registerReq.Email, registerReq.Username)
		c.JSON(http.StatusConflict, models.ErrorResponse{
			Error:   "User already exists",
			Message: "A user with this email or username already exists",
		})
		return
	}

	// Hash password
	log.Printf("Auth: Hashing password for user: %s", registerReq.Email)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerReq.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Auth: Password hashing failed for user: %s: %v", registerReq.Email, err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Password hashing error",
			Message: "Failed to hash password",
		})
		return
	}

	// Create user
	user := models.User{
		Email:     registerReq.Email,
		Username:  registerReq.Username,
		Password:  string(hashedPassword),
		FirstName: registerReq.FirstName,
		LastName:  registerReq.LastName,
		IsActive:  true,
	}

	log.Printf("Auth: Creating user in database: %s", registerReq.Email)
	if err := h.db.Create(&user).Error; err != nil {
		log.Printf("Auth: Database error creating user: %s: %v", registerReq.Email, err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Database error",
			Message: "Failed to create user",
		})
		return
	}

	log.Printf("Auth: User created successfully with ID: %d, email: %s", user.ID, user.Email)

	// Generate JWT token
	log.Printf("Auth: Generating JWT token for user ID: %d", user.ID)
	token, err := h.generateToken(user.ID)
	if err != nil {
		log.Printf("Auth: Token generation failed for user ID: %d: %v", user.ID, err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Token generation error",
			Message: "Failed to generate authentication token",
		})
		return
	}

	log.Printf("Auth: Registration successful for user ID: %d, email: %s", user.ID, user.Email)
	c.JSON(http.StatusCreated, models.AuthResponse{
		Token: token,
		User:  user,
	})
}

// Login authenticates a user
// @Summary      Login user
// @Description  Authenticate a user and return a JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        credentials body      models.LoginRequest  true  "User login credentials"
// @Success      200  {object}  models.AuthResponse
// @Failure      400  {object}  models.ErrorResponse
// @Failure      401  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	log.Printf("Auth: Login attempt from %s", c.ClientIP())
	
	var loginReq models.LoginRequest
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		log.Printf("Auth: Invalid login request from %s: %v", c.ClientIP(), err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

	log.Printf("Auth: Processing login for email: %s", loginReq.Email)

	// Find user by email
	var user models.User
	if err := h.db.Where("email = ?", loginReq.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("Auth: Login failed - user not found with email: %s", loginReq.Email)
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Error:   "Invalid credentials",
				Message: "Email or password is incorrect",
			})
			return
		}
		log.Printf("Auth: Database error retrieving user: %s: %v", loginReq.Email, err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Database error",
			Message: "Failed to retrieve user",
		})
		return
	}

	log.Printf("Auth: User found with ID: %d, email: %s", user.ID, user.Email)

	// Check if user is active
	if !user.IsActive {
		log.Printf("Auth: Login failed - account disabled for user ID: %d, email: %s", user.ID, user.Email)
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "Account disabled",
			Message: "Your account has been disabled",
		})
		return
	}

	// Verify password
	log.Printf("Auth: Verifying password for user ID: %d", user.ID)
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password)); err != nil {
		log.Printf("Auth: Login failed - invalid password for user ID: %d, email: %s", user.ID, user.Email)
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "Invalid credentials",
			Message: "Email or password is incorrect",
		})
		return
	}

	log.Printf("Auth: Password verified successfully for user ID: %d", user.ID)

	// Generate JWT token
	log.Printf("Auth: Generating JWT token for user ID: %d", user.ID)
	token, err := h.generateToken(user.ID)
	if err != nil {
		log.Printf("Auth: Token generation failed for user ID: %d: %v", user.ID, err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Token generation error",
			Message: "Failed to generate authentication token",
		})
		return
	}

	log.Printf("Auth: Login successful for user ID: %d, email: %s", user.ID, user.Email)
	c.JSON(http.StatusOK, models.AuthResponse{
		Token: token,
		User:  user,
	})
}

// generateToken generates a JWT token for the user
func (h *AuthHandler) generateToken(userID uint) (string, error) {
	log.Printf("Auth: Creating JWT token for user ID: %d", userID)
	
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.Load().JWTSecret))
	
	if err != nil {
		log.Printf("Auth: Failed to sign JWT token for user ID: %d: %v", userID, err)
		return "", err
	}
	
	log.Printf("Auth: JWT token created successfully for user ID: %d", userID)
	return tokenString, nil
}
