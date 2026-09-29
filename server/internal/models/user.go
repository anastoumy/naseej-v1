package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a staff member or admin within a merchant.
type User struct {
	TenantModel
	Email        string `gorm:"size:255;not null"`
	PasswordHash string `gorm:"size:255;not null"`
	FirstName    string `gorm:"size:100;not null"`
	LastName     string `gorm:"size:100;not null"`
	Phone        string `gorm:"size:20"`
	PinCode      string `gorm:"size:255"`
	AvatarURL    string
	IsActive     bool `gorm:"default:true"`
	LastLoginAt  *time.Time
	Roles        []Role     `gorm:"many2many:user_roles;"`
	Locations    []Location `gorm:"many2many:user_locations;"`
}

// TableName overrides the default table name.
func (User) TableName() string {
	return "users"
}

// BeforeCreate sets a new UUID if not already set.
func (u *User) BeforeCreate(tx interface{}) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}
