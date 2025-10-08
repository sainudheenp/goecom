package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sainudheenp/goecom/models"
	"gorm.io/gorm"
)

// ProductHandler handles product endpoints
type ProductHandler struct {
	db *gorm.DB
}

// NewProductHandler creates a new product handler
func NewProductHandler(db *gorm.DB) *ProductHandler {
	return &ProductHandler{
		db: db,
	}
}

// ListProducts lists products with filtering and pagination
func (h *ProductHandler) ListProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	q := c.Query("q")

	var products []models.Product
	dbQuery := h.db.Model(&models.Product{})

	if q != "" {
		dbQuery = dbQuery.Where("name ILIKE ? OR description ILIKE ?", "%"+q+"%", "%"+q+"%")
	}

	var total int64
	if err := dbQuery.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to count products",
		})
		return
	}

	offset := (page - 1) * size
	if err := dbQuery.Limit(size).Offset(offset).Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to list products",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"products": products,
		"total":    total,
		"page":     page,
		"size":     size,
	})
}

// GetProduct retrieves a product by ID
// @Summary Get product by ID
// @Tags products
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} store.Product
// @Failure 404 {object} ErrorResponse
// GetProduct retrieves a product by ID
func (h *ProductHandler) GetProduct(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid product ID",
		})
		return
	}

	var product models.Product
	if err := h.db.First(&product, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "product not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get product",
		})
		return
	}

	c.JSON(http.StatusOK, product)
}
