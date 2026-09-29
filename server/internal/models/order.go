package models

import "github.com/google/uuid"

// OrderChannel identifies the sales channel.
type OrderChannel string

const (
	ChannelPOS    OrderChannel = "POS"
	ChannelWeb    OrderChannel = "WEB"
	ChannelMobile OrderChannel = "MOBILE"
)

// OrderStatus tracks order lifecycle.
type OrderStatus string

const (
	StatusPending    OrderStatus = "PENDING"
	StatusConfirmed  OrderStatus = "CONFIRMED"
	StatusProcessing OrderStatus = "PROCESSING"
	StatusShipped    OrderStatus = "SHIPPED"
	StatusDelivered  OrderStatus = "DELIVERED"
	StatusCancelled  OrderStatus = "CANCELLED"
	StatusRefunded   OrderStatus = "REFUNDED"
)

// Order represents a unified order from any channel (POS, Web, Mobile).
type Order struct {
	TenantModel
	OrderNumber    string       `gorm:"size:20;not null"`
	CustomerID     *uuid.UUID   `gorm:"type:uuid"`
	Customer       *Customer    `gorm:"foreignKey:CustomerID"`
	LocationID     *uuid.UUID   `gorm:"type:uuid"`
	Location       *Location    `gorm:"foreignKey:LocationID"`
	StaffID        *uuid.UUID   `gorm:"type:uuid"`
	Staff          *User        `gorm:"foreignKey:StaffID"`
	Channel        OrderChannel `gorm:"size:10;not null"`
	Status         OrderStatus  `gorm:"size:20;not null;default:'PENDING'"`
	Subtotal       float64      `gorm:"type:decimal(10,2);not null"`
	TaxAmount      float64      `gorm:"type:decimal(10,2);not null;default:0"`
	DiscountAmount float64      `gorm:"type:decimal(10,2);not null;default:0"`
	Total          float64      `gorm:"type:decimal(10,2);not null"`
	Notes          string
	CouponID       *uuid.UUID  `gorm:"type:uuid"`
	Coupon         *Coupon     `gorm:"foreignKey:CouponID"`
	Items          []OrderItem `gorm:"foreignKey:OrderID"`
	Payments       []Payment   `gorm:"foreignKey:OrderID"`
}

// OrderItem represents a line item in an order (snapshotted at time of sale).
type OrderItem struct {
	BaseModel
	OrderID        uuid.UUID      `gorm:"type:uuid;not null;index"`
	VariantID      uuid.UUID      `gorm:"type:uuid;not null"`
	Variant        ProductVariant `gorm:"foreignKey:VariantID"`
	ProductName    string         `gorm:"size:255;not null"`
	VariantLabel   string         `gorm:"size:100;not null"`
	SKU            string         `gorm:"size:50;not null"`
	Quantity       int            `gorm:"not null;check:quantity > 0"`
	UnitPrice      float64        `gorm:"type:decimal(10,2);not null"`
	DiscountAmount float64        `gorm:"type:decimal(10,2);default:0"`
	Total          float64        `gorm:"type:decimal(10,2);not null"`
}

// TableName overrides the default table name.
func (Order) TableName() string {
	return "orders"
}

// TableName overrides the default table name.
func (OrderItem) TableName() string {
	return "order_items"
}
