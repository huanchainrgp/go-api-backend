package handlers

import (
	"log"
	"net/http"
	"strconv"

	"go-api-test1/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// TransactionHandler handles transaction-related HTTP requests
type TransactionHandler struct {
	db *gorm.DB
}

// NewTransactionHandler creates a new TransactionHandler
func NewTransactionHandler(db *gorm.DB) *TransactionHandler {
	return &TransactionHandler{db: db}
}

// GetTransactions retrieves all transactions
// @Summary      Get all transactions
// @Description  Get a list of all transactions
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}  models.Transaction
// @Failure      401  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /transactions [get]
func (h *TransactionHandler) GetTransactions(c *gin.Context) {
	log.Printf("Transaction: GetTransactions request from %s", c.ClientIP())
	
	var transactions []models.Transaction
	if err := h.db.Preload("User").Preload("Asset").Find(&transactions).Error; err != nil {
		log.Printf("Transaction: Database error retrieving transactions: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Database error",
			Message: "Failed to retrieve transactions",
		})
		return
	}

	log.Printf("Transaction: Successfully retrieved %d transactions", len(transactions))
	c.JSON(http.StatusOK, transactions)
}

// GetTransaction retrieves a specific transaction by ID
// @Summary      Get transaction by ID
// @Description  Get a specific transaction by its ID
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Transaction ID"
// @Success      200  {object}  models.Transaction
// @Failure      400  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /transactions/{id} [get]
func (h *TransactionHandler) GetTransaction(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		log.Printf("Transaction: Invalid transaction ID format: %s from %s", c.Param("id"), c.ClientIP())
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid ID",
			Message: "Transaction ID must be a valid number",
		})
		return
	}

	log.Printf("Transaction: GetTransaction request for ID: %d from %s", id, c.ClientIP())

	var transaction models.Transaction
	if err := h.db.Preload("User").Preload("Asset").First(&transaction, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("Transaction: Transaction not found with ID: %d", id)
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error:   "Transaction not found",
				Message: "The requested transaction does not exist",
			})
			return
		}
		log.Printf("Transaction: Database error retrieving transaction ID: %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Database error",
			Message: "Failed to retrieve transaction",
		})
		return
	}

	log.Printf("Transaction: Successfully retrieved transaction ID: %d, type: %s, amount: %.2f", transaction.ID, transaction.Type, transaction.Amount)
	c.JSON(http.StatusOK, transaction)
}

// CreateTransaction creates a new transaction
// @Summary      Create transaction
// @Description  Create a new transaction
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        transaction body      models.CreateTransactionRequest  true  "Transaction data"
// @Success      201  {object}  models.Transaction
// @Failure      400  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /transactions [post]
func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	log.Printf("Transaction: CreateTransaction request from %s", c.ClientIP())
	
	var createReq models.CreateTransactionRequest
	if err := c.ShouldBindJSON(&createReq); err != nil {
		log.Printf("Transaction: Invalid create request from %s: %v", c.ClientIP(), err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

	// Get user ID from JWT token
	userID, exists := c.Get("user_id")
	if !exists {
		log.Printf("Transaction: User ID not found in token from %s", c.ClientIP())
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "Unauthorized",
			Message: "User ID not found in token",
		})
		return
	}

	log.Printf("Transaction: Creating transaction for user ID: %v, asset ID: %d, type: %s, amount: %.2f", 
		userID, createReq.AssetID, createReq.Type, createReq.Amount)

	// Verify asset exists
	var asset models.Asset
	if err := h.db.First(&asset, createReq.AssetID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("Transaction: Asset not found with ID: %d", createReq.AssetID)
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error:   "Asset not found",
				Message: "The specified asset does not exist",
			})
			return
		}
		log.Printf("Transaction: Database error verifying asset ID: %d: %v", createReq.AssetID, err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Database error",
			Message: "Failed to verify asset",
		})
		return
	}

	log.Printf("Transaction: Asset verified - ID: %d, name: %s", asset.ID, asset.Name)

	// Calculate total value
	totalValue := createReq.Amount * createReq.Price
	log.Printf("Transaction: Calculated total value: %.2f (amount: %.2f * price: %.2f)", totalValue, createReq.Amount, createReq.Price)

	transaction := models.Transaction{
		UserID:      userID.(uint),
		AssetID:     createReq.AssetID,
		Type:        createReq.Type,
		Amount:      createReq.Amount,
		Price:       createReq.Price,
		TotalValue:  totalValue,
		Status:      "pending",
		Description: createReq.Description,
	}

	if err := h.db.Create(&transaction).Error; err != nil {
		log.Printf("Transaction: Database error creating transaction: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Database error",
			Message: "Failed to create transaction",
		})
		return
	}

	log.Printf("Transaction: Successfully created transaction ID: %d for user ID: %v", transaction.ID, userID)

	// Load the created transaction with relationships
	h.db.Preload("User").Preload("Asset").First(&transaction, transaction.ID)

	c.JSON(http.StatusCreated, transaction)
}

// UpdateTransaction updates a specific transaction
// @Summary      Update transaction
// @Description  Update a specific transaction by its ID
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Transaction ID"
// @Param        transaction body      models.UpdateTransactionRequest  true  "Transaction update data"
// @Success      200  {object}  models.Transaction
// @Failure      400  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /transactions/{id} [put]
func (h *TransactionHandler) UpdateTransaction(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		log.Printf("Transaction: Invalid transaction ID format for update: %s from %s", c.Param("id"), c.ClientIP())
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid ID",
			Message: "Transaction ID must be a valid number",
		})
		return
	}

	log.Printf("Transaction: UpdateTransaction request for ID: %d from %s", id, c.ClientIP())

	var transaction models.Transaction
	if err := h.db.First(&transaction, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("Transaction: Transaction not found for update with ID: %d", id)
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error:   "Transaction not found",
				Message: "The requested transaction does not exist",
			})
			return
		}
		log.Printf("Transaction: Database error retrieving transaction for update ID: %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Database error",
			Message: "Failed to retrieve transaction",
		})
		return
	}

	var updateReq models.UpdateTransactionRequest
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		log.Printf("Transaction: Invalid update request for transaction ID: %d: %v", id, err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

	log.Printf("Transaction: Updating transaction ID: %d with fields: type=%s, amount=%.2f, price=%.2f, status=%s", 
		id, updateReq.Type, updateReq.Amount, updateReq.Price, updateReq.Status)

	// Update fields if provided
	if updateReq.Type != "" {
		transaction.Type = updateReq.Type
	}
	if updateReq.Amount > 0 {
		transaction.Amount = updateReq.Amount
	}
	if updateReq.Price > 0 {
		transaction.Price = updateReq.Price
	}
	if updateReq.Status != "" {
		transaction.Status = updateReq.Status
	}
	if updateReq.Description != "" {
		transaction.Description = updateReq.Description
	}

	// Recalculate total value if amount or price changed
	if updateReq.Amount > 0 || updateReq.Price > 0 {
		transaction.TotalValue = transaction.Amount * transaction.Price
		log.Printf("Transaction: Recalculated total value: %.2f", transaction.TotalValue)
	}

	if err := h.db.Save(&transaction).Error; err != nil {
		log.Printf("Transaction: Database error updating transaction ID: %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Database error",
			Message: "Failed to update transaction",
		})
		return
	}

	log.Printf("Transaction: Successfully updated transaction ID: %d", transaction.ID)

	// Load the updated transaction with relationships
	h.db.Preload("User").Preload("Asset").First(&transaction, transaction.ID)

	c.JSON(http.StatusOK, transaction)
}

// DeleteTransaction deletes a specific transaction
// @Summary      Delete transaction
// @Description  Delete a specific transaction by its ID
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Transaction ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /transactions/{id} [delete]
func (h *TransactionHandler) DeleteTransaction(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		log.Printf("Transaction: Invalid transaction ID format for delete: %s from %s", c.Param("id"), c.ClientIP())
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid ID",
			Message: "Transaction ID must be a valid number",
		})
		return
	}

	log.Printf("Transaction: DeleteTransaction request for ID: %d from %s", id, c.ClientIP())

	var transaction models.Transaction
	if err := h.db.First(&transaction, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("Transaction: Transaction not found for delete with ID: %d", id)
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error:   "Transaction not found",
				Message: "The requested transaction does not exist",
			})
			return
		}
		log.Printf("Transaction: Database error retrieving transaction for delete ID: %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Database error",
			Message: "Failed to retrieve transaction",
		})
		return
	}

	log.Printf("Transaction: Deleting transaction ID: %d, type: %s, amount: %.2f", transaction.ID, transaction.Type, transaction.Amount)

	if err := h.db.Delete(&transaction).Error; err != nil {
		log.Printf("Transaction: Database error deleting transaction ID: %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Database error",
			Message: "Failed to delete transaction",
		})
		return
	}

	log.Printf("Transaction: Successfully deleted transaction ID: %d", transaction.ID)
	c.JSON(http.StatusOK, gin.H{"message": "Transaction deleted successfully"})
}
