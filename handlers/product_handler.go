package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sainudheenp/goecom/internal/middleware"
	"github.com/sainudheenp/goecom/internal/service"
	"github.com/sainudheenp/goecom/internal/store"
)

// ProductHandler handles product endpoints
type ProductHandler struct {
	productService *service.ProductService
}

// NewProductHandler creates a new product handler
func NewProductHandler(productService *service.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

// ListProducts lists products with filtering and pagination
// @Summary List products
// @Tags products
// @Produce json
// @Param q query string false "Search query"
// @Param page query int false "Page number" default(1)
// @Param size query int false "Page size" default(20)
// @Param min_price query int false "Minimum price in cents"
// @Param max_price query int false "Maximum price in cents"
// @Param sort query string false "Sort by: price_asc, price_desc, name_asc, name_desc, created_desc"
// @Success 200 {object} store.ProductListResult
// @Router /api/v1/products [get]
func (h *ProductHandler) ListProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	query := c.Query("q")
	sort := c.Query("sort")

	var minPrice, maxPrice *int
	if minPriceStr := c.Query("min_price"); minPriceStr != "" {
		if mp, err := strconv.Atoi(minPriceStr); err == nil {
			minPrice = &mp
		}
	}
	if maxPriceStr := c.Query("max_price"); maxPriceStr != "" {
		if mp, err := strconv.Atoi(maxPriceStr); err == nil {
			maxPrice = &mp
		}
	}

	filter := store.ProductFilter{
		Query:    query,
		MinPrice: minPrice,
		MaxPrice: maxPrice,
		Sort:     sort,
		Page:     page,
		Size:     size,
	}

	result, err := h.productService.ListProducts(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to list products",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetProduct retrieves a product by ID
// @Summary Get product by ID
// @Tags products
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} store.Product
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/products/{id} [get]
func (h *ProductHandler) GetProduct(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid product ID",
		})
		return
	}

	product, err := h.productService.GetProduct(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "product not found",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, product)
}

// CreateProduct creates a new product (admin only)
// @Summary Create product
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body service.CreateProductRequest true "Product details"
// @Success 201 {object} store.Product
// @Failure 400 {object} ErrorResponse
// @Router /api/v1/products [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req service.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request",
			"details": err.Error(),
		})
		return
	}

	product, err := h.productService.CreateProduct(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "failed to create product",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, product)
}

// UpdateProduct updates a product (admin only)
// @Summary Update product
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Product ID"
// @Param request body service.UpdateProductRequest true "Product update details"
// @Success 200 {object} store.Product
// @Failure 400 {object} ErrorResponse
// @Router /api/v1/products/{id} [put]
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid product ID",
		})
		return
	}

	var req service.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request",
			"details": err.Error(),
		})
		return
	}

	product, err := h.productService.UpdateProduct(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "failed to update product",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, product)
}

// DeleteProduct deletes a product (admin only)
// @Summary Delete product
// @Tags products
// @Security BearerAuth
// @Param id path string true "Product ID"
// @Success 204
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/products/{id} [delete]
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid product ID",
		})
		return
	}

	if err := h.productService.DeleteProduct(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "failed to delete product",
			"details": err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// BulkImportProducts imports multiple products (admin only)
// @Summary Bulk import products
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body []service.CreateProductRequest true "Products to import"
// @Success 201 {object} []store.Product
// @Failure 400 {object} ErrorResponse
// @Router /api/v1/products/bulk [post]
func (h *ProductHandler) BulkImportProducts(c *gin.Context) {
	user, _ := middleware.GetUserFromContext(c)
	if user.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "admin access required",
		})
		return
	}

	var req []service.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request",
			"details": err.Error(),
		})
		return
	}

	products, err := h.productService.BulkImportProducts(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "failed to import products",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, products)
}
