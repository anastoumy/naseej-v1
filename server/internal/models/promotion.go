package models

import (
	"time"

	"github.com/lib/pq"
)

// DiscountType defines the type of discount applied by a coupon.
type DiscountType string

const (
	DiscountPercentage DiscountType = "PERCENTAGE"
	DiscountFixed      DiscountType = "FIXED_AMOUNT"
	DiscountBuyXGetY   DiscountType = "BUY_X_GET_Y"
)

// Coupon represents a promotional discount code.
type Coupon struct {
	TenantModel
	Code                 string `gorm:"size:50;not null"`
	Description          string
	DiscountType         DiscountType `gorm:"size:20;not null"`
	DiscountValue        float64      `gorm:"type:decimal(10,2);not null"`
	MinOrderValue        *float64     `gorm:"type:decimal(10,2)"`
	MaxDiscount          *float64     `gorm:"type:decimal(10,2)"`
	UsageLimit           *int
	UsedCount            int            `gorm:"default:0"`
	PerCustomerLimit     int            `gorm:"default:1"`
	ValidFrom            time.Time      `gorm:"not null"`
	ValidUntil           time.Time      `gorm:"not null"`
	IsActive             bool           `gorm:"default:true"`
	ApplicableCategories pq.StringArray `gorm:"type:uuid[]"`
	ApplicableProducts   pq.StringArray `gorm:"type:uuid[]"`
}

// TableName overrides the default table name.
func (Coupon) TableName() string {
	return "coupons"
}
