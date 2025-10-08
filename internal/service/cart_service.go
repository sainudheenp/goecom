package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/sainudheenp/goecom/internal/store"
	"gorm.io/gorm"
)

// CartService handles cart business logic
type CartService struct {
	cartRepo    *store.CartRepository
	productRepo *store.ProductRepository
}

// NewCartService creates a new cart service
func NewCartService(cartRepo *store.CartRepository, productRepo *store.ProductRepository) *CartService {
	return &CartService{
		cartRepo:    cartRepo,
		productRepo: productRepo,
	}
}

// AddToCartRequest represents add to cart input
type AddToCartRequest struct {
	ProductID uuid.UUID `json:"product_id" binding:"required"`
	Quantity  int       `json:"quantity" binding:"required,min=1"`
}

// CartResponse represents cart output
type CartResponse struct {
	Items      []CartItemResponse `json:"items"`
	TotalCents int                `json:"total_cents"`
	Currency   string             `json:"currency"`
}

// CartItemResponse represents a cart item output
type CartItemResponse struct {
	ID        uuid.UUID      `json:"id"`
	ProductID uuid.UUID      `json:"product_id"`
	Product   *store.Product `json:"product"`
	Quantity  int            `json:"quantity"`
	Subtotal  int            `json:"subtotal_cents"`
}

// AddToCart adds or updates an item in the cart
func (s *CartService) AddToCart(ctx context.Context, userID uuid.UUID, req AddToCartRequest) (*CartResponse, error) {
	// Verify product exists and has sufficient stock
	product, err := s.productRepo.GetByID(ctx, req.ProductID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	if product.Stock < req.Quantity {
		return nil, errors.New("insufficient stock")
	}

	// Add or update cart item
	cartItem := &store.CartItem{
		UserID:    userID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
	}

	if err := s.cartRepo.AddOrUpdate(ctx, cartItem); err != nil {
		return nil, fmt.Errorf("failed to add to cart: %w", err)
	}

	// Return updated cart
	return s.GetCart(ctx, userID)
}

// GetCart retrieves the user's cart
func (s *CartService) GetCart(ctx context.Context, userID uuid.UUID) (*CartResponse, error) {
	items, err := s.cartRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	var totalCents int
	var currency string = "USD"
	cartItems := make([]CartItemResponse, 0, len(items))

	for _, item := range items {
		if item.Product == nil {
			continue
		}

		subtotal := item.Product.PriceCents * item.Quantity
		totalCents += subtotal
		currency = item.Product.Currency

		cartItems = append(cartItems, CartItemResponse{
			ID:        item.ID,
			ProductID: item.ProductID,
			Product:   item.Product,
			Quantity:  item.Quantity,
			Subtotal:  subtotal,
		})
	}

	return &CartResponse{
		Items:      cartItems,
		TotalCents: totalCents,
		Currency:   currency,
	}, nil
}

// RemoveFromCart removes an item from the cart
func (s *CartService) RemoveFromCart(ctx context.Context, userID uuid.UUID, itemID uuid.UUID) error {
	return s.cartRepo.Delete(ctx, itemID, userID)
}

// ClearCart clears all items from the cart
func (s *CartService) ClearCart(ctx context.Context, userID uuid.UUID) error {
	return s.cartRepo.Clear(ctx, userID)
}
