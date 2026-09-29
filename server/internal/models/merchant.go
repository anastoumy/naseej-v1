package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// Merchant represents a tenant in the multi-tenant system.
type Merchant struct {
	BaseModel
	BusinessName   string `gorm:"size:150;not null"`
	Slug           string `gorm:"size:50;uniqueIndex;not null"`
	Domain         string `gorm:"size:255"`
	LogoURL        string
	PrimaryColor   string         `gorm:"size:7"`
	SecondaryColor string         `gorm:"size:7"`
	TaxRate        float64        `gorm:"type:decimal(5,4);default:0.15"`
	Currency       string         `gorm:"size:3;default:'SAR'"`
	Timezone       string         `gorm:"size:50;default:'Asia/Riyadh'"`
	IsActive       bool           `gorm:"default:true"`
	Settings       datatypes.JSON `gorm:"type:jsonb;default:'{}'"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// Location represents a physical store branch or warehouse.
type Location struct {
	TenantModel
	Name        string `gorm:"size:150;not null"`
	Address     string
	City        string `gorm:"size:100"`
	Country     string `gorm:"size:2"`
	Phone       string `gorm:"size:20"`
	IsWarehouse bool   `gorm:"default:false"`
	IsActive    bool   `gorm:"default:true"`
}

// MerchantSettings holds dynamic merchant configuration loaded on bootstrap.
type MerchantSettings struct {
	LoyaltyEnabled        bool     `json:"loyalty_enabled"`
	LoyaltyPointsPerUnit  int      `json:"loyalty_points_per_unit"`
	LoyaltyRedemptionRate float64  `json:"loyalty_redemption_rate"`
	PaymentGateways       []string `json:"payment_gateways"`
	ReceiptHeader         string   `json:"receipt_header"`
	ReceiptFooter         string   `json:"receipt_footer"`
}

// TableName overrides the default table name.
func (Merchant) TableName() string {
	return "merchants"
}

// TableName overrides the default table name.
func (Location) TableName() string {
	return "locations"
}

// BeforeCreate sets a new UUID if not already set.
func (m *Merchant) BeforeCreate(tx interface{}) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}
