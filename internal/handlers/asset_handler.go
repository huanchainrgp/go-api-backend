package handlers

import (
	"log"
	"net/http"
	"strconv"

	"go-api-test1/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AssetHandler handles asset-related HTTP requests
type AssetHandler struct {
	db *gorm.DB
}

// NewAssetHandler creates a new AssetHandler
func NewAssetHandler(db *gorm.DB) *AssetHandler {
	return &AssetHandler{db: db}
}

// GetAssets retrieves all assets
// @Summary      Get all assets
// @Description  Get a list of all assets
// @Tags         assets
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}  models.Asset
// @Failure      401  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /assets [get]
func (h *AssetHandler) GetAssets(c *gin.Context) {
	log.Printf("Asset: GetAssets request from %s", c.ClientIP())
	
	var assets []models.Asset
	if err := h.db.Find(&assets).Error; err != nil {
		log.Printf("Asset: Database error retrieving assets: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Database error",
			Message: "Failed to retrieve assets",
		})
		return
	}

	log.Printf("Asset: Successfully retrieved %d assets", len(assets))
	c.JSON(http.StatusOK, assets)
}

// GetAsset retrieves a specific asset by ID
// @Summary      Get asset by ID
// @Description  Get a specific asset by its ID
// @Tags         assets
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Asset ID"
// @Success      200  {object}  models.Asset
// @Failure      400  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /assets/{id} [get]
func (h *AssetHandler) GetAsset(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		log.Printf("Asset: Invalid asset ID format: %s from %s", c.Param("id"), c.ClientIP())
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid ID",
			Message: "Asset ID must be a valid number",
		})
		return
	}

	log.Printf("Asset: GetAsset request for ID: %d from %s", id, c.ClientIP())

	var asset models.Asset
	if err := h.db.First(&asset, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("Asset: Asset not found with ID: %d", id)
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error:   "Asset not found",
				Message: "The requested asset does not exist",
			})
			return
		}
		log.Printf("Asset: Database error retrieving asset ID: %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Database error",
			Message: "Failed to retrieve asset",
		})
		return
	}

	log.Printf("Asset: Successfully retrieved asset ID: %d, name: %s", asset.ID, asset.Name)
	c.JSON(http.StatusOK, asset)
}

// CreateAsset creates a new asset
// @Summary      Create asset
// @Description  Create a new asset
// @Tags         assets
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        asset body      models.CreateAssetRequest  true  "Asset data"
// @Success      201  {object}  models.Asset
// @Failure      400  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /assets [post]
func (h *AssetHandler) CreateAsset(c *gin.Context) {
	log.Printf("Asset: CreateAsset request from %s", c.ClientIP())
	
	var createReq models.CreateAssetRequest
	if err := c.ShouldBindJSON(&createReq); err != nil {
		log.Printf("Asset: Invalid create request from %s: %v", c.ClientIP(), err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

	log.Printf("Asset: Creating asset with name: %s, symbol: %s, type: %s", createReq.Name, createReq.Symbol, createReq.Type)

	asset := models.Asset{
		Name:        createReq.Name,
		Symbol:      createReq.Symbol,
		Type:        createReq.Type,
		Description: createReq.Description,
		Price:       createReq.Price,
		IsActive:    true,
	}

	if err := h.db.Create(&asset).Error; err != nil {
		log.Printf("Asset: Database error creating asset: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Database error",
			Message: "Failed to create asset",
		})
		return
	}

	log.Printf("Asset: Successfully created asset ID: %d, name: %s", asset.ID, asset.Name)
	c.JSON(http.StatusCreated, asset)
}

// UpdateAsset updates a specific asset
// @Summary      Update asset
// @Description  Update a specific asset by its ID
// @Tags         assets
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Asset ID"
// @Param        asset body      models.UpdateAssetRequest  true  "Asset update data"
// @Success      200  {object}  models.Asset
// @Failure      400  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /assets/{id} [put]
func (h *AssetHandler) UpdateAsset(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		log.Printf("Asset: Invalid asset ID format for update: %s from %s", c.Param("id"), c.ClientIP())
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid ID",
			Message: "Asset ID must be a valid number",
		})
		return
	}

	log.Printf("Asset: UpdateAsset request for ID: %d from %s", id, c.ClientIP())

	var asset models.Asset
	if err := h.db.First(&asset, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("Asset: Asset not found for update with ID: %d", id)
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error:   "Asset not found",
				Message: "The requested asset does not exist",
			})
			return
		}
		log.Printf("Asset: Database error retrieving asset for update ID: %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Database error",
			Message: "Failed to retrieve asset",
		})
		return
	}

	var updateReq models.UpdateAssetRequest
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		log.Printf("Asset: Invalid update request for asset ID: %d: %v", id, err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

	log.Printf("Asset: Updating asset ID: %d with fields: name=%s, symbol=%s, type=%s, price=%.2f", 
		id, updateReq.Name, updateReq.Symbol, updateReq.Type, updateReq.Price)

	// Update fields if provided
	if updateReq.Name != "" {
		asset.Name = updateReq.Name
	}
	if updateReq.Symbol != "" {
		asset.Symbol = updateReq.Symbol
	}
	if updateReq.Type != "" {
		asset.Type = updateReq.Type
	}
	if updateReq.Description != "" {
		asset.Description = updateReq.Description
	}
	if updateReq.Price > 0 {
		asset.Price = updateReq.Price
	}
	if updateReq.IsActive != nil {
		asset.IsActive = *updateReq.IsActive
	}

	if err := h.db.Save(&asset).Error; err != nil {
		log.Printf("Asset: Database error updating asset ID: %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Database error",
			Message: "Failed to update asset",
		})
		return
	}

	log.Printf("Asset: Successfully updated asset ID: %d, name: %s", asset.ID, asset.Name)
	c.JSON(http.StatusOK, asset)
}

// DeleteAsset deletes a specific asset
// @Summary      Delete asset
// @Description  Delete a specific asset by its ID
// @Tags         assets
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Asset ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /assets/{id} [delete]
func (h *AssetHandler) DeleteAsset(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		log.Printf("Asset: Invalid asset ID format for delete: %s from %s", c.Param("id"), c.ClientIP())
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid ID",
			Message: "Asset ID must be a valid number",
		})
		return
	}

	log.Printf("Asset: DeleteAsset request for ID: %d from %s", id, c.ClientIP())

	var asset models.Asset
	if err := h.db.First(&asset, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("Asset: Asset not found for delete with ID: %d", id)
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error:   "Asset not found",
				Message: "The requested asset does not exist",
			})
			return
		}
		log.Printf("Asset: Database error retrieving asset for delete ID: %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Database error",
			Message: "Failed to retrieve asset",
		})
		return
	}

	log.Printf("Asset: Deleting asset ID: %d, name: %s", asset.ID, asset.Name)

	if err := h.db.Delete(&asset).Error; err != nil {
		log.Printf("Asset: Database error deleting asset ID: %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Database error",
			Message: "Failed to delete asset",
		})
		return
	}

	log.Printf("Asset: Successfully deleted asset ID: %d, name: %s", asset.ID, asset.Name)
	c.JSON(http.StatusOK, gin.H{"message": "Asset deleted successfully"})
}
