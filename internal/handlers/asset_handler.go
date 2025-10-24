package handlers

import (
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
	var assets []models.Asset
	if err := h.db.Find(&assets).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Database error",
			Message: "Failed to retrieve assets",
		})
		return
	}

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
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid ID",
			Message: "Asset ID must be a valid number",
		})
		return
	}

	var asset models.Asset
	if err := h.db.First(&asset, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error:   "Asset not found",
				Message: "The requested asset does not exist",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Database error",
			Message: "Failed to retrieve asset",
		})
		return
	}

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
	var createReq models.CreateAssetRequest
	if err := c.ShouldBindJSON(&createReq); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

	asset := models.Asset{
		Name:        createReq.Name,
		Symbol:      createReq.Symbol,
		Type:        createReq.Type,
		Description: createReq.Description,
		Price:       createReq.Price,
		IsActive:    true,
	}

	if err := h.db.Create(&asset).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Database error",
			Message: "Failed to create asset",
		})
		return
	}

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
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid ID",
			Message: "Asset ID must be a valid number",
		})
		return
	}

	var asset models.Asset
	if err := h.db.First(&asset, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error:   "Asset not found",
				Message: "The requested asset does not exist",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Database error",
			Message: "Failed to retrieve asset",
		})
		return
	}

	var updateReq models.UpdateAssetRequest
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

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
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Database error",
			Message: "Failed to update asset",
		})
		return
	}

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
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid ID",
			Message: "Asset ID must be a valid number",
		})
		return
	}

	var asset models.Asset
	if err := h.db.First(&asset, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error:   "Asset not found",
				Message: "The requested asset does not exist",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Database error",
			Message: "Failed to retrieve asset",
		})
		return
	}

	if err := h.db.Delete(&asset).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Database error",
			Message: "Failed to delete asset",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Asset deleted successfully"})
}
