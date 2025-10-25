package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"go-api-test1/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// LoggerMiddleware logs HTTP requests with timing and status information
func LoggerMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("[%s] %s %s %d %s %s %s\n",
			param.TimeStamp.Format(time.RFC3339),
			param.Method,
			param.Path,
			param.StatusCode,
			param.Latency,
			param.ClientIP,
			param.ErrorMessage,
		)
	})
}

// CORS middleware for handling cross-origin requests
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("CORS: Processing %s request from %s to %s", c.Request.Method, c.ClientIP(), c.Request.URL.Path)
		
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			log.Printf("CORS: Handling OPTIONS preflight request from %s", c.ClientIP())
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// AuthMiddleware validates JWT tokens
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("Auth: Validating token for %s request to %s from %s", c.Request.Method, c.Request.URL.Path, c.ClientIP())
		
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Printf("Auth: Missing authorization header for %s from %s", c.Request.URL.Path, c.ClientIP())
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Check if the header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			log.Printf("Auth: Invalid authorization header format for %s from %s", c.Request.URL.Path, c.ClientIP())
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		// Extract the token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse and validate the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				log.Printf("Auth: Invalid signing method for token from %s", c.ClientIP())
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(config.Load().JWTSecret), nil
		})

		if err != nil || !token.Valid {
			log.Printf("Auth: Invalid or expired token for %s from %s: %v", c.Request.URL.Path, c.ClientIP(), err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Extract user ID from claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if userID, exists := claims["user_id"]; exists {
				log.Printf("Auth: Token validated successfully for user %v accessing %s", userID, c.Request.URL.Path)
				c.Set("user_id", userID)
			} else {
				log.Printf("Auth: Missing user_id in token claims for %s from %s", c.Request.URL.Path, c.ClientIP())
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
				c.Abort()
				return
			}
		} else {
			log.Printf("Auth: Invalid token claims format for %s from %s", c.Request.URL.Path, c.ClientIP())
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		c.Next()
	}
}
