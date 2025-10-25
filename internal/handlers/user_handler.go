package handlers

import (
	"log"
	"net/http"
	"strconv"

	"go-api-test1/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	db *gorm.DB
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{db: db}
}

// GetUsers retrieves all users
// @Summary      Get all users
// @Description  Get a list of all users
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}  models.User
// @Failure      401  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /users [get]
func (h *UserHandler) GetUsers(c *gin.Context) {
	log.Printf("User: GetUsers request from %s", c.ClientIP())
	
	var users []models.User
	if err := h.db.Find(&users).Error; err != nil {
		log.Printf("User: Database error retrieving users: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Database error",
			Message: "Failed to retrieve users",
		})
		return
	}

	log.Printf("User: Successfully retrieved %d users", len(users))
	c.JSON(http.StatusOK, users)
}

// GetUser retrieves a specific user by ID
// @Summary      Get user by ID
// @Description  Get a specific user by their ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  models.User
// @Failure      400  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		log.Printf("User: Invalid user ID format: %s from %s", c.Param("id"), c.ClientIP())
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid ID",
			Message: "User ID must be a valid number",
		})
		return
	}

	log.Printf("User: GetUser request for ID: %d from %s", id, c.ClientIP())

	var user models.User
	if err := h.db.First(&user, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("User: User not found with ID: %d", id)
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error:   "User not found",
				Message: "The requested user does not exist",
			})
			return
		}
		log.Printf("User: Database error retrieving user ID: %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Database error",
			Message: "Failed to retrieve user",
		})
		return
	}

	log.Printf("User: Successfully retrieved user ID: %d, email: %s", user.ID, user.Email)
	c.JSON(http.StatusOK, user)
}

// UpdateUser updates a specific user
// @Summary      Update user
// @Description  Update a specific user by their ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "User ID"
// @Param        user body      models.UpdateUserRequest  true  "User update data"
// @Success      200  {object}  models.User
// @Failure      400  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		log.Printf("User: Invalid user ID format for update: %s from %s", c.Param("id"), c.ClientIP())
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid ID",
			Message: "User ID must be a valid number",
		})
		return
	}

	log.Printf("User: UpdateUser request for ID: %d from %s", id, c.ClientIP())

	var user models.User
	if err := h.db.First(&user, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("User: User not found for update with ID: %d", id)
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error:   "User not found",
				Message: "The requested user does not exist",
			})
			return
		}
		log.Printf("User: Database error retrieving user for update ID: %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Database error",
			Message: "Failed to retrieve user",
		})
		return
	}

	var updateReq models.UpdateUserRequest
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		log.Printf("User: Invalid update request for user ID: %d: %v", id, err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

	log.Printf("User: Updating user ID: %d with fields: email=%s, username=%s, firstName=%s, lastName=%s", 
		id, updateReq.Email, updateReq.Username, updateReq.FirstName, updateReq.LastName)

	// Update fields if provided
	if updateReq.Email != "" {
		user.Email = updateReq.Email
	}
	if updateReq.Username != "" {
		user.Username = updateReq.Username
	}
	if updateReq.FirstName != "" {
		user.FirstName = updateReq.FirstName
	}
	if updateReq.LastName != "" {
		user.LastName = updateReq.LastName
	}
	if updateReq.IsActive != nil {
		user.IsActive = *updateReq.IsActive
	}

	if err := h.db.Save(&user).Error; err != nil {
		log.Printf("User: Database error updating user ID: %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Database error",
			Message: "Failed to update user",
		})
		return
	}

	log.Printf("User: Successfully updated user ID: %d, email: %s", user.ID, user.Email)
	c.JSON(http.StatusOK, user)
}

// DeleteUser deletes a specific user
// @Summary      Delete user
// @Description  Delete a specific user by their ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		log.Printf("User: Invalid user ID format for delete: %s from %s", c.Param("id"), c.ClientIP())
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid ID",
			Message: "User ID must be a valid number",
		})
		return
	}

	log.Printf("User: DeleteUser request for ID: %d from %s", id, c.ClientIP())

	var user models.User
	if err := h.db.First(&user, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("User: User not found for delete with ID: %d", id)
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error:   "User not found",
				Message: "The requested user does not exist",
			})
			return
		}
		log.Printf("User: Database error retrieving user for delete ID: %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Database error",
			Message: "Failed to retrieve user",
		})
		return
	}

	log.Printf("User: Deleting user ID: %d, email: %s", user.ID, user.Email)

	if err := h.db.Delete(&user).Error; err != nil {
		log.Printf("User: Database error deleting user ID: %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Database error",
			Message: "Failed to delete user",
		})
		return
	}

	log.Printf("User: Successfully deleted user ID: %d, email: %s", user.ID, user.Email)
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
