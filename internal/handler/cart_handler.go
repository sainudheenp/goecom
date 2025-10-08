package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sainudheenp/goecom/internal/middleware"
	"github.com/sainudheenp/goecom/internal/service"
)

// CartHandler handles cart endpoints
type CartHandler struct {
	cartService *service.CartService
}

// NewCartHandler creates a new cart handler
func NewCartHandler(cartService *service.CartService) *CartHandler {
	return &CartHandler{
		cartService: cartService,
	}
}

// AddToCart adds or updates an item in the cart
// @Summary Add to cart
// @Tags cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body service.AddToCartRequest true "Cart item"
// @Success 200 {object} service.CartResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/v1/cart [post]
func (h *CartHandler) AddToCart(c *gin.Context) {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	var req service.AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request",
			"details": err.Error(),
		})
		return
	}

	cart, err := h.cartService.AddToCart(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "failed to add to cart",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, cart)
}

// GetCart retrieves the user's cart
// @Summary Get cart
// @Tags cart
// @Produce json
// @Security BearerAuth
// @Success 200 {object} service.CartResponse
// @Router /api/v1/cart [get]
func (h *CartHandler) GetCart(c *gin.Context) {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	cart, err := h.cartService.GetCart(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to get cart",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, cart)
}

// RemoveFromCart removes an item from the cart
// @Summary Remove from cart
// @Tags cart
// @Security BearerAuth
// @Param item_id path string true "Cart item ID"
// @Success 204
// @Router /api/v1/cart/{item_id} [delete]
func (h *CartHandler) RemoveFromCart(c *gin.Context) {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	itemID, err := uuid.Parse(c.Param("item_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid item ID",
		})
		return
	}

	if err := h.cartService.RemoveFromCart(c.Request.Context(), userID, itemID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "failed to remove from cart",
			"details": err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}
