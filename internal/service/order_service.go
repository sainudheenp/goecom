package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/sainudheenp/goecom/internal/store"
	"gorm.io/gorm"
)

// OrderService handles order business logic
type OrderService struct {
	orderRepo   *store.OrderRepository
	cartRepo    *store.CartRepository
	productRepo *store.ProductRepository
	db          *store.DB
}

// NewOrderService creates a new order service
func NewOrderService(
	orderRepo *store.OrderRepository,
	cartRepo *store.CartRepository,
	productRepo *store.ProductRepository,
	db *store.DB,
) *OrderService {
	return &OrderService{
		orderRepo:   orderRepo,
		cartRepo:    cartRepo,
		productRepo: productRepo,
		db:          db,
	}
}

// CreateOrderRequest represents order creation input
type CreateOrderRequest struct {
	ShippingAddress map[string]interface{} `json:"shipping_address" binding:"required"`
}

// CreateOrder creates an order from the user's cart
func (s *OrderService) CreateOrder(ctx context.Context, userID uuid.UUID, req CreateOrderRequest) (*store.Order, error) {
	// Get cart items
	cartItems, err := s.cartRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	if len(cartItems) == 0 {
		return nil, errors.New("cart is empty")
	}

	var order *store.Order
	var orderItems []store.OrderItem

	// Create order in a transaction
	err = s.db.WithTransaction(ctx, func(tx *gorm.DB) error {
		// Calculate total and prepare order items
		var totalCents int
		var currency string = "USD"

		for _, cartItem := range cartItems {
			if cartItem.Product == nil {
				return fmt.Errorf("product not found for cart item %s", cartItem.ID)
			}

			// Check stock availability
			if cartItem.Product.Stock < cartItem.Quantity {
				return fmt.Errorf("insufficient stock for product %s", cartItem.Product.Name)
			}

			// Decrement stock
			if err := s.productRepo.DecrementStock(ctx, tx, cartItem.ProductID, cartItem.Quantity); err != nil {
				return err
			}

			subtotal := cartItem.Product.PriceCents * cartItem.Quantity
			totalCents += subtotal
			currency = cartItem.Product.Currency

			orderItems = append(orderItems, store.OrderItem{
				ProductID:  cartItem.ProductID,
				PriceCents: cartItem.Product.PriceCents,
				Quantity:   cartItem.Quantity,
			})
		}

		// Create order
		order = &store.Order{
			UserID:          userID,
			TotalCents:      totalCents,
			Currency:        currency,
			Status:          "pending",
			ShippingAddress: req.ShippingAddress,
		}

		if err := tx.Create(order).Error; err != nil {
			return fmt.Errorf("failed to create order: %w", err)
		}

		// Create order items
		for i := range orderItems {
			orderItems[i].OrderID = order.ID
		}

		if err := tx.Create(&orderItems).Error; err != nil {
			return fmt.Errorf("failed to create order items: %w", err)
		}

		// Clear cart
		if err := tx.Where("user_id = ?", userID).Delete(&store.CartItem{}).Error; err != nil {
			return fmt.Errorf("failed to clear cart: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Load order with items
	return s.orderRepo.GetByID(ctx, order.ID)
}

// GetOrder retrieves an order by ID
func (s *OrderService) GetOrder(ctx context.Context, orderID uuid.UUID, userID uuid.UUID) (*store.Order, error) {
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("order not found")
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	// Verify ownership
	if order.UserID != userID {
		return nil, errors.New("unauthorized to view this order")
	}

	return order, nil
}

// ListUserOrders retrieves orders for a user
func (s *OrderService) ListUserOrders(ctx context.Context, userID uuid.UUID, page, size int) ([]store.Order, int64, error) {
	return s.orderRepo.GetByUserID(ctx, userID, page, size)
}

// ListAllOrders retrieves all orders (admin)
func (s *OrderService) ListAllOrders(ctx context.Context, page, size int) ([]store.Order, int64, error) {
	return s.orderRepo.List(ctx, page, size)
}

// UpdateOrderStatus updates an order status (admin)
func (s *OrderService) UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status string) error {
	// Validate status
	validStatuses := map[string]bool{
		"pending":   true,
		"paid":      true,
		"shipped":   true,
		"cancelled": true,
	}

	if !validStatuses[status] {
		return errors.New("invalid order status")
	}

	// Check if order exists
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("order not found")
		}
		return fmt.Errorf("failed to get order: %w", err)
	}

	order.Status = status
	return s.orderRepo.Update(ctx, order)
}
