package store

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ProductRepository handles product data operations
type ProductRepository struct {
	db *DB
}

// NewProductRepository creates a new product repository
func NewProductRepository(db *DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// ProductFilter holds filter criteria for listing products
type ProductFilter struct {
	Query    string
	MinPrice *int
	MaxPrice *int
	Sort     string // price_asc, price_desc, name_asc, name_desc, created_desc
	Page     int
	Size     int
}

// ProductListResult holds paginated product results
type ProductListResult struct {
	Items []Product `json:"items"`
	Total int64     `json:"total"`
	Page  int       `json:"page"`
	Size  int       `json:"size"`
}

// Create creates a new product
func (r *ProductRepository) Create(ctx context.Context, product *Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

// GetByID retrieves a product by ID
func (r *ProductRepository) GetByID(ctx context.Context, id uuid.UUID) (*Product, error) {
	var product Product
	err := r.db.WithContext(ctx).First(&product, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// GetBySKU retrieves a product by SKU
func (r *ProductRepository) GetBySKU(ctx context.Context, sku string) (*Product, error) {
	var product Product
	err := r.db.WithContext(ctx).First(&product, "sku = ?", sku).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// List retrieves products with filtering and pagination
func (r *ProductRepository) List(ctx context.Context, filter ProductFilter) (*ProductListResult, error) {
	query := r.db.WithContext(ctx).Model(&Product{})

	// Apply filters
	if filter.Query != "" {
		searchPattern := fmt.Sprintf("%%%s%%", filter.Query)
		query = query.Where("name ILIKE ? OR description ILIKE ?", searchPattern, searchPattern)
	}

	if filter.MinPrice != nil {
		query = query.Where("price_cents >= ?", *filter.MinPrice)
	}

	if filter.MaxPrice != nil {
		query = query.Where("price_cents <= ?", *filter.MaxPrice)
	}

	// Count total
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// Apply sorting
	switch filter.Sort {
	case "price_asc":
		query = query.Order("price_cents ASC")
	case "price_desc":
		query = query.Order("price_cents DESC")
	case "name_asc":
		query = query.Order("name ASC")
	case "name_desc":
		query = query.Order("name DESC")
	case "created_desc":
		query = query.Order("created_at DESC")
	default:
		query = query.Order("created_at DESC")
	}

	// Apply pagination
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Size < 1 {
		filter.Size = 20
	}
	if filter.Size > 100 {
		filter.Size = 100
	}

	offset := (filter.Page - 1) * filter.Size
	query = query.Offset(offset).Limit(filter.Size)

	// Execute query
	var products []Product
	if err := query.Find(&products).Error; err != nil {
		return nil, err
	}

	return &ProductListResult{
		Items: products,
		Total: total,
		Page:  filter.Page,
		Size:  filter.Size,
	}, nil
}

// Update updates a product
func (r *ProductRepository) Update(ctx context.Context, product *Product) error {
	return r.db.WithContext(ctx).Save(product).Error
}

// Delete deletes a product
func (r *ProductRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&Product{}, "id = ?", id).Error
}

// DecrementStock decrements product stock atomically
func (r *ProductRepository) DecrementStock(ctx context.Context, tx *gorm.DB, productID uuid.UUID, quantity int) error {
	db := r.db.DB
	if tx != nil {
		db = tx
	}

	result := db.WithContext(ctx).Model(&Product{}).
		Where("id = ? AND stock >= ?", productID, quantity).
		Update("stock", gorm.Expr("stock - ?", quantity))

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("insufficient stock for product %s", productID)
	}

	return nil
}

// BulkCreate creates multiple products
func (r *ProductRepository) BulkCreate(ctx context.Context, products []Product) error {
	return r.db.WithContext(ctx).Create(&products).Error
}
