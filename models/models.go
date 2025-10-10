package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a user account
type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	Email        string    `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string    `gorm:"not null" json:"-"`
	FullName     string    `json:"full_name"`
	Role         string    `gorm:"not null;default:'user'" json:"role"` // user, admin
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// BeforeCreate hook to generate UUID before creating
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

// Product represents a product in the catalog
type Product struct {
	ID          uuid.UUID       `gorm:"type:uuid;primary_key;" json:"id"`
	SKU         string          `gorm:"uniqueIndex;not null" json:"sku"`
	Name        string          `gorm:"not null" json:"name"`
	Description string          `json:"description"`
	PriceCents  int             `gorm:"not null" json:"price_cents"`
	Currency    string          `gorm:"not null;default:'USD'" json:"currency"`
	Stock       int             `gorm:"not null;default:0" json:"stock"`
	Images      JSONStringSlice `gorm:"type:jsonb" json:"images"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

// BeforeCreate hook to generate UUID before creating
func (p *Product) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

// CartItem represents an item in a user's shopping cart
type CartItem struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index:idx_cart_user_product" json:"user_id"`
	User      *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	ProductID uuid.UUID `gorm:"type:uuid;not null;index:idx_cart_user_product" json:"product_id"`
	Product   *Product  `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Quantity  int       `gorm:"not null" json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BeforeCreate hook to generate UUID before creating
func (c *CartItem) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

// Order represents a customer order
type Order struct {
	ID              uuid.UUID   `gorm:"type:uuid;primary_key;" json:"id"`
	UserID          uuid.UUID   `gorm:"type:uuid;not null;index" json:"user_id"`
	User            *User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	TotalCents      int         `gorm:"not null" json:"total_cents"`
	Currency        string      `gorm:"not null" json:"currency"`
	Status          string      `gorm:"not null;default:'pending'" json:"status"` // pending, paid, shipped, cancelled
	ShippingAddress JSONMap     `gorm:"type:jsonb" json:"shipping_address"`
	PaymentInfo     JSONMap     `gorm:"type:jsonb" json:"payment_info,omitempty"`
	Items           []OrderItem `gorm:"foreignKey:OrderID" json:"items,omitempty"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
}

// BeforeCreate hook to generate UUID before creating
func (o *Order) BeforeCreate(tx *gorm.DB) error {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	return nil
}

// OrderItem represents a line item in an order
type OrderItem struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	OrderID    uuid.UUID `gorm:"type:uuid;not null;index" json:"order_id"`
	Order      *Order    `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	ProductID  uuid.UUID `gorm:"type:uuid;not null" json:"product_id"`
	Product    *Product  `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	PriceCents int       `gorm:"not null" json:"price_cents"`
	Quantity   int       `gorm:"not null" json:"quantity"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// BeforeCreate hook to generate UUID before creating
func (oi *OrderItem) BeforeCreate(tx *gorm.DB) error {
	if oi.ID == uuid.Nil {
		oi.ID = uuid.New()
	}
	return nil
}
