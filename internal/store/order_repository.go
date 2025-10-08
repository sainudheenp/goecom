package store

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// OrderRepository handles order data operations
type OrderRepository struct {
	db *DB
}

// NewOrderRepository creates a new order repository
func NewOrderRepository(db *DB) *OrderRepository {
	return &OrderRepository{db: db}
}

// Create creates a new order
func (r *OrderRepository) Create(ctx context.Context, order *Order) error {
	return r.db.WithContext(ctx).Create(order).Error
}

// CreateWithItems creates an order with items in a transaction
func (r *OrderRepository) CreateWithItems(ctx context.Context, order *Order, items []OrderItem) error {
	return r.db.WithTransaction(ctx, func(tx *gorm.DB) error {
		// Create order
		if err := tx.Create(order).Error; err != nil {
			return err
		}

		// Set order_id for all items
		for i := range items {
			items[i].OrderID = order.ID
		}

		// Create order items
		if err := tx.Create(&items).Error; err != nil {
			return err
		}

		return nil
	})
}

// GetByID retrieves an order by ID
func (r *OrderRepository) GetByID(ctx context.Context, id uuid.UUID) (*Order, error) {
	var order Order
	err := r.db.WithContext(ctx).
		Preload("Items").
		Preload("Items.Product").
		First(&order, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// GetByUserID retrieves orders for a specific user
func (r *OrderRepository) GetByUserID(ctx context.Context, userID uuid.UUID, page, size int) ([]Order, int64, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 20
	}
	if size > 100 {
		size = 100
	}

	var total int64
	query := r.db.WithContext(ctx).Model(&Order{}).Where("user_id = ?", userID)
	
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var orders []Order
	offset := (page - 1) * size
	err := query.
		Preload("Items").
		Preload("Items.Product").
		Order("created_at DESC").
		Offset(offset).
		Limit(size).
		Find(&orders).Error

	return orders, total, err
}

// List retrieves all orders (admin)
func (r *OrderRepository) List(ctx context.Context, page, size int) ([]Order, int64, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 20
	}
	if size > 100 {
		size = 100
	}

	var total int64
	query := r.db.WithContext(ctx).Model(&Order{})
	
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var orders []Order
	offset := (page - 1) * size
	err := query.
		Preload("Items").
		Preload("Items.Product").
		Preload("User").
		Order("created_at DESC").
		Offset(offset).
		Limit(size).
		Find(&orders).Error

	return orders, total, err
}

// UpdateStatus updates an order status
func (r *OrderRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	return r.db.WithContext(ctx).
		Model(&Order{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// Update updates an order
func (r *OrderRepository) Update(ctx context.Context, order *Order) error {
	return r.db.WithContext(ctx).Save(order).Error
}
