package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"

	"github.com/google/uuid"
	"github.com/sainudheenp/goecom/internal/store"
)

// PaymentService handles payment processing
type PaymentService struct {
	orderRepo *store.OrderRepository
}

// NewPaymentService creates a new payment service
func NewPaymentService(orderRepo *store.OrderRepository) *PaymentService {
	return &PaymentService{
		orderRepo: orderRepo,
	}
}

// ChargeRequest represents payment charge input
type ChargeRequest struct {
	OrderID        uuid.UUID              `json:"order_id" binding:"required"`
	PaymentMethod  string                 `json:"payment_method" binding:"required"` // card, upi, wallet
	PaymentDetails map[string]interface{} `json:"payment_details"`
}

// ChargeResponse represents payment charge output
type ChargeResponse struct {
	OrderID       uuid.UUID `json:"order_id"`
	Status        string    `json:"status"` // success, failed
	TransactionID string    `json:"transaction_id"`
	Message       string    `json:"message"`
}

// ProcessCharge processes a payment (stub implementation)
func (s *PaymentService) ProcessCharge(ctx context.Context, userID uuid.UUID, req ChargeRequest) (*ChargeResponse, error) {
	// Get order
	order, err := s.orderRepo.GetByID(ctx, req.OrderID)
	if err != nil {
		return nil, fmt.Errorf("order not found: %w", err)
	}

	// Verify ownership
	if order.UserID != userID {
		return nil, errors.New("unauthorized to process payment for this order")
	}

	// Check if order is already paid
	if order.Status == "paid" {
		return nil, errors.New("order is already paid")
	}

	// Stub payment processing - simulate success/failure
	success := s.simulatePayment(req.PaymentMethod)

	if success {
		// Update order status to paid
		order.Status = "paid"
		order.PaymentInfo = map[string]interface{}{
			"method":         req.PaymentMethod,
			"transaction_id": s.generateTransactionID(),
		}

		if err := s.orderRepo.Update(ctx, order); err != nil {
			return nil, fmt.Errorf("failed to update order: %w", err)
		}

		return &ChargeResponse{
			OrderID:       order.ID,
			Status:        "success",
			TransactionID: order.PaymentInfo["transaction_id"].(string),
			Message:       "Payment processed successfully",
		}, nil
	}

	return &ChargeResponse{
		OrderID: order.ID,
		Status:  "failed",
		Message: "Payment failed. Please try again.",
	}, nil
}

// simulatePayment simulates payment processing (stub)
func (s *PaymentService) simulatePayment(method string) bool {
	// 90% success rate
	return rand.Float32() < 0.9
}

// generateTransactionID generates a mock transaction ID
func (s *PaymentService) generateTransactionID() string {
	return fmt.Sprintf("TXN_%s", uuid.New().String()[:8])
}

// Note: In production, this would integrate with real payment providers:
// - Stripe: use stripe-go SDK
// - Razorpay: use razorpay-go SDK
// - PayPal: use PayPal REST API
//
// Example Stripe integration pattern:
//
// import "github.com/stripe/stripe-go/v76"
// import "github.com/stripe/stripe-go/v76/paymentintent"
//
// func (s *PaymentService) processStripePayment(amount int64, currency string) (*stripe.PaymentIntent, error) {
//     params := &stripe.PaymentIntentParams{
//         Amount:   stripe.Int64(amount),
//         Currency: stripe.String(currency),
//     }
//     return paymentintent.New(params)
// }
