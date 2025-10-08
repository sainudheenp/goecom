package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/sainudheenp/goecom/internal/store"
	"gorm.io/gorm"
)

// ProductService handles product business logic
type ProductService struct {
	productRepo *store.ProductRepository
}

// NewProductService creates a new product service
func NewProductService(productRepo *store.ProductRepository) *ProductService {
	return &ProductService{
		productRepo: productRepo,
	}
}

// CreateProductRequest represents product creation input
type CreateProductRequest struct {
	SKU         string   `json:"sku" binding:"required"`
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description"`
	PriceCents  int      `json:"price_cents" binding:"required,min=0"`
	Currency    string   `json:"currency" binding:"required"`
	Stock       int      `json:"stock" binding:"required,min=0"`
	Images      []string `json:"images"`
}

// UpdateProductRequest represents product update input
type UpdateProductRequest struct {
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	PriceCents  *int     `json:"price_cents" binding:"omitempty,min=0"`
	Currency    *string  `json:"currency"`
	Stock       *int     `json:"stock" binding:"omitempty,min=0"`
	Images      []string `json:"images"`
}

// CreateProduct creates a new product
func (s *ProductService) CreateProduct(ctx context.Context, req CreateProductRequest) (*store.Product, error) {
	// Check if SKU already exists
	existing, err := s.productRepo.GetBySKU(ctx, req.SKU)
	if err == nil && existing != nil {
		return nil, errors.New("product with this SKU already exists")
	}

	product := &store.Product{
		SKU:         req.SKU,
		Name:        req.Name,
		Description: req.Description,
		PriceCents:  req.PriceCents,
		Currency:    req.Currency,
		Stock:       req.Stock,
		Images:      req.Images,
	}

	if err := s.productRepo.Create(ctx, product); err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return product, nil
}

// GetProduct retrieves a product by ID
func (s *ProductService) GetProduct(ctx context.Context, id uuid.UUID) (*store.Product, error) {
	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, fmt.Errorf("failed to get product: %w", err)
	}
	return product, nil
}

// ListProducts retrieves products with filtering and pagination
func (s *ProductService) ListProducts(ctx context.Context, filter store.ProductFilter) (*store.ProductListResult, error) {
	return s.productRepo.List(ctx, filter)
}

// UpdateProduct updates a product
func (s *ProductService) UpdateProduct(ctx context.Context, id uuid.UUID, req UpdateProductRequest) (*store.Product, error) {
	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	// Update fields
	if req.Name != nil {
		product.Name = *req.Name
	}
	if req.Description != nil {
		product.Description = *req.Description
	}
	if req.PriceCents != nil {
		product.PriceCents = *req.PriceCents
	}
	if req.Currency != nil {
		product.Currency = *req.Currency
	}
	if req.Stock != nil {
		product.Stock = *req.Stock
	}
	if req.Images != nil {
		product.Images = req.Images
	}

	if err := s.productRepo.Update(ctx, product); err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	return product, nil
}

// DeleteProduct deletes a product
func (s *ProductService) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	_, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("product not found")
		}
		return fmt.Errorf("failed to get product: %w", err)
	}

	if err := s.productRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	return nil
}

// BulkImportProducts imports multiple products
func (s *ProductService) BulkImportProducts(ctx context.Context, requests []CreateProductRequest) ([]store.Product, error) {
	products := make([]store.Product, 0, len(requests))

	for _, req := range requests {
		product := store.Product{
			SKU:         req.SKU,
			Name:        req.Name,
			Description: req.Description,
			PriceCents:  req.PriceCents,
			Currency:    req.Currency,
			Stock:       req.Stock,
			Images:      req.Images,
		}
		products = append(products, product)
	}

	if err := s.productRepo.BulkCreate(ctx, products); err != nil {
		return nil, fmt.Errorf("failed to bulk import products: %w", err)
	}

	return products, nil
}
