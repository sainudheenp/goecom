package store

import (
	"context"

	"github.com/google/uuid"
)

// CartRepository handles cart data operations
type CartRepository struct {
	db *DB
}

// NewCartRepository creates a new cart repository
func NewCartRepository(db *DB) *CartRepository {
	return &CartRepository{db: db}
}

// AddOrUpdate adds or updates a cart item
func (r *CartRepository) AddOrUpdate(ctx context.Context, item *CartItem) error {
	// Check if item already exists
	var existing CartItem
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND product_id = ?", item.UserID, item.ProductID).
		First(&existing).Error

	if err == nil {
		// Update existing item
		existing.Quantity = item.Quantity
		return r.db.WithContext(ctx).Save(&existing).Error
	}

	// Create new item
	return r.db.WithContext(ctx).Create(item).Error
}

// GetByUserID retrieves all cart items for a user
func (r *CartRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]CartItem, error) {
	var items []CartItem
	err := r.db.WithContext(ctx).
		Preload("Product").
		Where("user_id = ?", userID).
		Find(&items).Error
	return items, err
}

// GetByID retrieves a cart item by ID
func (r *CartRepository) GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*CartItem, error) {
	var item CartItem
	err := r.db.WithContext(ctx).
		Preload("Product").
		Where("id = ? AND user_id = ?", id, userID).
		First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// Delete deletes a cart item
func (r *CartRepository) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&CartItem{}).Error
}

// Clear clears all cart items for a user
func (r *CartRepository) Clear(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&CartItem{}).Error
}
