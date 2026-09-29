package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// PaymentMethod identifies how a payment was made.
type PaymentMethod string

const (
	PayCash         PaymentMethod = "CASH"
	PayCard         PaymentMethod = "CARD"
	PayBankTransfer PaymentMethod = "BANK_TRANSFER"
	PayLoyalty      PaymentMethod = "LOYALTY_POINTS"
)

// PaymentStatus tracks payment lifecycle.
type PaymentStatus string

const (
	PayStatusPending   PaymentStatus = "PENDING"
	PayStatusCompleted PaymentStatus = "COMPLETED"
	PayStatusFailed    PaymentStatus = "FAILED"
	PayStatusRefunded  PaymentStatus = "REFUNDED"
)

// Payment represents a single payment against an order.
type Payment struct {
	TenantModel
	OrderID     uuid.UUID      `gorm:"type:uuid;not null;index"`
	Order       Order          `gorm:"foreignKey:OrderID"`
	Method      PaymentMethod  `gorm:"size:20;not null"`
	Status      PaymentStatus  `gorm:"size:20;not null;default:'PENDING'"`
	Amount      float64        `gorm:"type:decimal(10,2);not null"`
	Reference   string         `gorm:"size:255"`
	Metadata    datatypes.JSON `gorm:"type:jsonb;default:'{}'"`
	ProcessedAt *time.Time
}

// Refund represents a refund issued against a payment.
type Refund struct {
	TenantModel
	PaymentID   uuid.UUID `gorm:"type:uuid;not null"`
	Payment     Payment   `gorm:"foreignKey:PaymentID"`
	OrderID     uuid.UUID `gorm:"type:uuid;not null"`
	Order       Order     `gorm:"foreignKey:OrderID"`
	Amount      float64   `gorm:"type:decimal(10,2);not null"`
	Reason      string
	ProcessedBy *uuid.UUID `gorm:"type:uuid"`
	Processor   *User      `gorm:"foreignKey:ProcessedBy"`
}

// TableName overrides the default table name.
func (Payment) TableName() string {
	return "payments"
}

// TableName overrides the default table name.
func (Refund) TableName() string {
	return "refunds"
}
