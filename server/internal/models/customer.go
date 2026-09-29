package models

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
)

// Customer represents an end customer (shopper).
type Customer struct {
	TenantModel
	FirstName          string         `gorm:"size:100;not null"`
	LastName           string         `gorm:"size:100;not null"`
	Email              string         `gorm:"size:255"`
	Phone              string         `gorm:"size:20"`
	LoyaltyPoints      int            `gorm:"default:0;check:loyalty_points >= 0"`
	LoyaltyCardBarcode string         `gorm:"size:50"`
	Tags               pq.StringArray `gorm:"type:text[]"`
	Notes              string
	TotalSpent         float64           `gorm:"type:decimal(12,2);default:0"`
	OrderCount         int               `gorm:"default:0"`
	Addresses          []CustomerAddress `gorm:"foreignKey:CustomerID"`
}

// CustomerAddress represents a stored address for a customer.
type CustomerAddress struct {
	BaseModel
	CustomerID   uuid.UUID `gorm:"type:uuid;not null;index"`
	Label        string    `gorm:"size:50"`
	AddressLine1 string    `gorm:"size:255;not null"`
	AddressLine2 string    `gorm:"size:255"`
	City         string    `gorm:"size:100;not null"`
	State        string    `gorm:"size:100"`
	PostalCode   string    `gorm:"size:20"`
	Country      string    `gorm:"size:2;not null"`
	IsDefault    bool      `gorm:"default:false"`
}

// LoyaltyTransaction records loyalty point changes.
type LoyaltyTransaction struct {
	TenantModel
	CustomerID   uuid.UUID  `gorm:"type:uuid;not null;index"`
	Customer     Customer   `gorm:"foreignKey:CustomerID"`
	OrderID      *uuid.UUID `gorm:"type:uuid"`
	Order        *Order     `gorm:"foreignKey:OrderID"`
	PointsChange int        `gorm:"not null"`
	Reason       string     `gorm:"size:100;not null"`
	BalanceAfter int        `gorm:"not null"`
}

// TableName overrides the default table name.
func (Customer) TableName() string {
	return "customers"
}

// TableName overrides the default table name.
func (CustomerAddress) TableName() string {
	return "customer_addresses"
}

// TableName overrides the default table name.
func (LoyaltyTransaction) TableName() string {
	return "loyalty_transactions"
}
