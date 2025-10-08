package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sainudheenp/goecom/internal/middleware"
	"github.com/sainudheenp/goecom/internal/service"
)

// PaymentHandler handles payment endpoints
type PaymentHandler struct {
	paymentService *service.PaymentService
}

// NewPaymentHandler creates a new payment handler
func NewPaymentHandler(paymentService *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
	}
}

// ProcessCharge processes a payment for an order
// @Summary Process payment
// @Tags payments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body service.ChargeRequest true "Payment details"
// @Success 200 {object} service.ChargeResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/v1/payments/charge [post]
func (h *PaymentHandler) ProcessCharge(c *gin.Context) {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	var req service.ChargeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request",
			"details": err.Error(),
		})
		return
	}

	resp, err := h.paymentService.ProcessCharge(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "payment processing failed",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}
