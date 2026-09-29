package models

import "github.com/google/uuid"

// Role represents an RBAC role assigned to users.
type Role struct {
	TenantModel
	Name        string `gorm:"size:50;not null"`
	Description string
	IsSystem    bool         `gorm:"default:false"`
	Permissions []Permission `gorm:"many2many:role_permissions;"`
}

// Permission represents a single resource:action permission.
type Permission struct {
	BaseModel
	Resource    string `gorm:"size:50;not null"`
	Action      string `gorm:"size:50;not null"`
	Description string
}

// ManagerOverride logs when a manager authorizes a sensitive POS action.
type ManagerOverride struct {
	TenantModel
	ManagerID  uuid.UUID `gorm:"type:uuid;not null"`
	Manager    User      `gorm:"foreignKey:ManagerID"`
	CashierID  uuid.UUID `gorm:"type:uuid;not null"`
	Cashier    User      `gorm:"foreignKey:CashierID"`
	Action     string    `gorm:"size:100;not null"`
	Reason     string
	LocationID *uuid.UUID `gorm:"type:uuid"`
	Location   *Location  `gorm:"foreignKey:LocationID"`
}

// TableName overrides the default table name.
func (Role) TableName() string {
	return "roles"
}

// TableName overrides the default table name.
func (Permission) TableName() string {
	return "permissions"
}

// TableName overrides the default table name.
func (ManagerOverride) TableName() string {
	return "manager_overrides"
}

// Permission constants following resource:action naming.
const (
	PermProductsView      = "products:view"
	PermProductsCreate    = "products:create"
	PermProductsEdit      = "products:edit"
	PermProductsDelete    = "products:delete"
	PermProductsPublish   = "products:publish"
	PermInventoryView     = "inventory:view"
	PermInventoryEdit     = "inventory:edit"
	PermInventoryTransfer = "inventory:transfer"
	PermInventoryPrice    = "inventory:price"
	PermOrdersView        = "orders:view"
	PermOrdersCreate      = "orders:create"
	PermOrdersEdit        = "orders:edit"
	PermOrdersCancel      = "orders:cancel"
	PermOrdersRefund      = "orders:refund"
	PermOrdersDiscount    = "orders:discount"
	PermOrdersExport      = "orders:export"
	PermPOSAccess         = "pos:access"
	PermPOSOpenDrawer     = "pos:open_drawer"
	PermPOSVoidItem       = "pos:void_item"
	PermPOSOverridePrice  = "pos:override_price"
	PermPOSApplyCoupon    = "pos:apply_coupon"
	PermPOSParkSale       = "pos:park_sale"
	PermPOSCashIn         = "pos:cash_in"
	PermPOSCashOut        = "pos:cash_out"
	PermPOSEndShift       = "pos:end_shift"
	PermPOSEndShiftOther  = "pos:end_shift_other"
	PermCustomersView     = "customers:view"
	PermCustomersCreate   = "customers:create"
	PermCustomersEdit     = "customers:edit"
	PermCustomersDelete   = "customers:delete"
	PermCustomersExport   = "customers:export"
	PermCustomersLoyalty  = "customers:loyalty"
	PermCouponsView       = "coupons:view"
	PermCouponsCreate     = "coupons:create"
	PermCouponsEdit       = "coupons:edit"
	PermCouponsDelete     = "coupons:delete"
	PermReportsSales      = "reports:sales"
	PermReportsInventory  = "reports:inventory"
	PermReportsCustomers  = "reports:customers"
	PermReportsFinancial  = "reports:financial"
	PermStaffView         = "staff:view"
	PermStaffCreate       = "staff:create"
	PermStaffEdit         = "staff:edit"
	PermStaffDelete       = "staff:delete"
	PermStaffAssignRole   = "staff:assign_role"
	PermSettingsView      = "settings:view"
	PermSettingsEdit      = "settings:edit"
	PermSettingsLocations = "settings:locations"
)
