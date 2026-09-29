# Data Models & Schemas

## Naseej Core — GORM Models & GraphQL Type Definitions

**Database:** PostgreSQL 15+  
**ORM:** GORM (Go) — AutoMigrate on server startup  
**API:** GraphQL via gqlgen  
**Multi-Tenancy:** Row-level isolation (`merchant_id` on all tenant-scoped tables)

---

## 1. GORM Model Definitions

All models are defined as Go structs. GORM AutoMigrate creates and updates tables automatically on server startup. No manual SQL migration files.

### 1.1 Base Models

```go
// internal/models/base.go
package models

import (
    "time"
    "github.com/google/uuid"
)

// BaseModel is embedded in non-tenant entities (like Merchant itself)
type BaseModel struct {
    ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    CreatedAt time.Time
    UpdatedAt time.Time
}

// TenantModel is embedded in all merchant-scoped entities
type TenantModel struct {
    ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    MerchantID uuid.UUID `gorm:"type:uuid;not null;index"`
    CreatedAt  time.Time
    UpdatedAt  time.Time
}
```

### 1.2 Merchant

```go
// internal/models/merchant.go
type Merchant struct {
    BaseModel
    BusinessName    string         `gorm:"size:150;not null"`
    Slug            string         `gorm:"size:50;uniqueIndex;not null"`
    Domain          string         `gorm:"size:255"`
    LogoURL         string
    PrimaryColor    string         `gorm:"size:7"`
    SecondaryColor  string         `gorm:"size:7"`
    TaxRate         float64        `gorm:"type:decimal(5,4);default:0.15"`
    Currency        string         `gorm:"size:3;default:'SAR'"`
    Timezone        string         `gorm:"size:50;default:'Asia/Riyadh'"`
    IsActive        bool           `gorm:"default:true"`
    Settings        datatypes.JSON `gorm:"type:jsonb;default:'{}'"`
    // Settings JSON holds: payment gateways, loyalty config, receipt header/footer
}
```

### 1.3 Users & Staff (RBAC)

```go
// internal/models/user.go
type User struct {
    TenantModel
    Email        string     `gorm:"size:255;not null"`
    PasswordHash string     `gorm:"size:255;not null"`
    FirstName    string     `gorm:"size:100;not null"`
    LastName     string     `gorm:"size:100;not null"`
    Phone        string     `gorm:"size:20"`
    PinCode      string     `gorm:"size:255"`          // hashed PIN for POS quick-login
    AvatarURL    string
    IsActive     bool       `gorm:"default:true"`
    LastLoginAt  *time.Time
    Roles        []Role     `gorm:"many2many:user_roles;"`
    Locations    []Location `gorm:"many2many:user_locations;"`
}
// Unique constraint: (merchant_id, email) — added via AutoMigrate hook or composite index

// internal/models/staff.go
type Role struct {
    TenantModel
    Name        string       `gorm:"size:50;not null"`
    Description string
    IsSystem    bool         `gorm:"default:false"`
    Permissions []Permission `gorm:"many2many:role_permissions;"`
}

type Permission struct {
    BaseModel
    Resource    string `gorm:"size:50;not null"`
    Action      string `gorm:"size:50;not null"`
    Description string
}
// Unique constraint: (resource, action)

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
```

### 1.4 Locations

```go
// internal/models/location.go (defined within user.go or standalone)
type Location struct {
    TenantModel
    Name        string `gorm:"size:150;not null"`
    Address     string
    City        string `gorm:"size:100"`
    Country     string `gorm:"size:2"` // ISO 3166-1 alpha-2
    Phone       string `gorm:"size:20"`
    IsWarehouse bool   `gorm:"default:false"`
    IsActive    bool   `gorm:"default:true"`
}
```

### 1.5 Product Catalog (with SizeGroup)

```go
// internal/models/size_group.go

// SizeGroup represents a named collection of sizes applicable to a product category.
// Examples: "Shoes Women (36-41)", "Shoes Men (40-45)", "Shirts (S-XXL)", "Kids (2Y-12Y)"
// When a user creates a product under a category, the form auto-populates with
// the sizes from that category's assigned SizeGroup.
type SizeGroup struct {
    TenantModel
    Name  string `gorm:"size:100;not null"` // e.g. "Women's Shoes", "Men's Shirts"
    Sizes []Size `gorm:"foreignKey:SizeGroupID"`
}

type Size struct {
    TenantModel
    SizeGroupID uuid.UUID `gorm:"type:uuid;not null;index"`
    SizeGroup   SizeGroup `gorm:"foreignKey:SizeGroupID"`
    Name        string    `gorm:"size:20;not null"` // e.g. 'S', 'M', 'L', '38', '42'
    SortOrder   int       `gorm:"default:0"`
}

// internal/models/category.go
type Category struct {
    TenantModel
    ParentID    *uuid.UUID `gorm:"type:uuid;index"`
    Parent      *Category  `gorm:"foreignKey:ParentID"`
    Children    []Category `gorm:"foreignKey:ParentID"`
    SizeGroupID *uuid.UUID `gorm:"type:uuid"`          // links category to its size system
    SizeGroup   *SizeGroup `gorm:"foreignKey:SizeGroupID"`
    Name        string     `gorm:"size:100;not null"`
    Slug        string     `gorm:"size:100;not null"`
    Description string
    ImageURL    string
    SortOrder   int        `gorm:"default:0"`
}

type Color struct {
    TenantModel
    Name    string `gorm:"size:50;not null"`
    HexCode string `gorm:"size:7"`
}

type Brand struct {
    TenantModel
    Name    string `gorm:"size:100;not null"`
    Slug    string `gorm:"size:100;not null"`
    LogoURL string
}

// internal/models/product.go
type Product struct {
    TenantModel
    BrandID      *uuid.UUID     `gorm:"type:uuid"`
    Brand        *Brand         `gorm:"foreignKey:BrandID"`
    CategoryID   *uuid.UUID     `gorm:"type:uuid"`
    Category     *Category      `gorm:"foreignKey:CategoryID"`
    Name         string         `gorm:"size:255;not null"`
    Slug         string         `gorm:"size:255;not null"`
    Description  string
    BasePrice    float64        `gorm:"type:decimal(10,2);not null"`
    ComparePrice *float64       `gorm:"type:decimal(10,2)"`
    CostPrice    *float64       `gorm:"type:decimal(10,2)"`
    IsPublished  bool           `gorm:"default:false"`
    IsActive     bool           `gorm:"default:true"`
    Tags         pq.StringArray `gorm:"type:text[]"`
    Images       []ProductImage `gorm:"foreignKey:ProductID"`
    Variants     []ProductVariant `gorm:"foreignKey:ProductID"`
}

type ProductImage struct {
    BaseModel
    ProductID uuid.UUID `gorm:"type:uuid;not null;index"`
    URL       string    `gorm:"not null"`
    AltText   string    `gorm:"size:255"`
    SortOrder int       `gorm:"default:0"`
    IsPrimary bool      `gorm:"default:false"`
}

type ProductVariant struct {
    BaseModel
    ProductID     uuid.UUID `gorm:"type:uuid;not null;index"`
    Product       Product   `gorm:"foreignKey:ProductID"`
    SizeID        uuid.UUID `gorm:"type:uuid;not null"`
    Size          Size      `gorm:"foreignKey:SizeID"`
    ColorID       uuid.UUID `gorm:"type:uuid;not null"`
    Color         Color     `gorm:"foreignKey:ColorID"`
    SKU           string    `gorm:"size:50;uniqueIndex;not null"`
    Barcode       string    `gorm:"size:50;uniqueIndex"`
    PriceOverride *float64  `gorm:"type:decimal(10,2)"`
    CostOverride  *float64  `gorm:"type:decimal(10,2)"`
    WeightGrams   *int
    IsActive      bool      `gorm:"default:true"`
}
// Unique constraint: (product_id, size_id, color_id)
```

**SizeGroup Flow — How It Works:**

```
1. Admin creates SizeGroups:
   - "Women's Shoes" → sizes: 36, 37, 38, 39, 40, 41
   - "Men's Shoes"   → sizes: 40, 41, 42, 43, 44, 45
   - "Shirts"        → sizes: S, M, L, XL, XXL
   - "Kids Clothes"  → sizes: 2Y, 4Y, 6Y, 8Y, 10Y, 12Y

2. Admin assigns SizeGroup to Category:
   - Category "Women's Footwear" → SizeGroup "Women's Shoes"
   - Category "Men's Shirts"     → SizeGroup "Shirts"

3. When creating a product under "Women's Footwear":
   - The variant creation form auto-populates size options: 36, 37, 38...41
   - User only needs to pick colors and confirm which sizes to stock
   - SKU/barcode auto-generated or entered per variant
```

### 1.6 Inventory

```go
// internal/models/inventory.go
type InventoryLevel struct {
    BaseModel
    VariantID    uuid.UUID      `gorm:"type:uuid;not null;index"`
    Variant      ProductVariant `gorm:"foreignKey:VariantID"`
    LocationID   uuid.UUID      `gorm:"type:uuid;not null;index"`
    Location     Location       `gorm:"foreignKey:LocationID"`
    Quantity     int            `gorm:"not null;default:0;check:quantity >= 0"`
    Reserved     int            `gorm:"not null;default:0;check:reserved >= 0"`
    ReorderPoint int            `gorm:"default:0"`
}
// Unique constraint: (variant_id, location_id)

type MovementReason string
const (
    MovementSale        MovementReason = "sale"
    MovementReturn      MovementReason = "return"
    MovementAdjustment  MovementReason = "adjustment"
    MovementTransferIn  MovementReason = "transfer_in"
    MovementTransferOut MovementReason = "transfer_out"
    MovementRestock     MovementReason = "restock"
    MovementDamaged     MovementReason = "damaged"
    MovementLost        MovementReason = "lost"
    MovementManualCount MovementReason = "manual_count"
)

type InventoryMovement struct {
    TenantModel
    VariantID      uuid.UUID      `gorm:"type:uuid;not null;index"`
    Variant        ProductVariant `gorm:"foreignKey:VariantID"`
    LocationID     uuid.UUID      `gorm:"type:uuid;not null"`
    Location       Location       `gorm:"foreignKey:LocationID"`
    QuantityChange int            `gorm:"not null"` // positive = in, negative = out
    Reason         MovementReason `gorm:"size:20;not null"`
    ReferenceID    *uuid.UUID     `gorm:"type:uuid"`
    ReferenceType  string         `gorm:"size:50"` // 'order', 'transfer', 'adjustment'
    Notes          string
    PerformedBy    *uuid.UUID     `gorm:"type:uuid"`
    Performer      *User          `gorm:"foreignKey:PerformedBy"`
}
```

### 1.7 Orders

```go
// internal/models/order.go
type OrderChannel string
const (
    ChannelPOS    OrderChannel = "POS"
    ChannelWeb    OrderChannel = "WEB"
    ChannelMobile OrderChannel = "MOBILE"
)

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
    CouponID       *uuid.UUID   `gorm:"type:uuid"`
    Coupon         *Coupon      `gorm:"foreignKey:CouponID"`
    Items          []OrderItem  `gorm:"foreignKey:OrderID"`
    Payments       []Payment    `gorm:"foreignKey:OrderID"`
}
// Unique constraint: (merchant_id, order_number)

type OrderItem struct {
    BaseModel
    OrderID        uuid.UUID `gorm:"type:uuid;not null;index"`
    VariantID      uuid.UUID `gorm:"type:uuid;not null"`
    Variant        ProductVariant `gorm:"foreignKey:VariantID"`
    ProductName    string    `gorm:"size:255;not null"`    // snapshotted
    VariantLabel   string    `gorm:"size:100;not null"`    // e.g. "M / Navy Blue"
    SKU            string    `gorm:"size:50;not null"`     // snapshotted
    Quantity       int       `gorm:"not null;check:quantity > 0"`
    UnitPrice      float64   `gorm:"type:decimal(10,2);not null"`
    DiscountAmount float64   `gorm:"type:decimal(10,2);default:0"`
    Total          float64   `gorm:"type:decimal(10,2);not null"`
}
```

### 1.8 Payments

```go
// internal/models/payment.go
type PaymentMethod string
const (
    PayCash         PaymentMethod = "CASH"
    PayCard         PaymentMethod = "CARD"
    PayBankTransfer PaymentMethod = "BANK_TRANSFER"
    PayLoyalty      PaymentMethod = "LOYALTY_POINTS"
)

type PaymentStatus string
const (
    PayStatusPending   PaymentStatus = "PENDING"
    PayStatusCompleted PaymentStatus = "COMPLETED"
    PayStatusFailed    PaymentStatus = "FAILED"
    PayStatusRefunded  PaymentStatus = "REFUNDED"
)

type Payment struct {
    TenantModel
    OrderID     uuid.UUID     `gorm:"type:uuid;not null;index"`
    Order       Order         `gorm:"foreignKey:OrderID"`
    Method      PaymentMethod `gorm:"size:20;not null"`
    Status      PaymentStatus `gorm:"size:20;not null;default:'PENDING'"`
    Amount      float64       `gorm:"type:decimal(10,2);not null"`
    Reference   string        `gorm:"size:255"`
    Metadata    datatypes.JSON `gorm:"type:jsonb;default:'{}'"`
    ProcessedAt *time.Time
}

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
```

### 1.9 Customers

```go
// internal/models/customer.go
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
    TotalSpent         float64        `gorm:"type:decimal(12,2);default:0"`
    OrderCount         int            `gorm:"default:0"`
    Addresses          []CustomerAddress `gorm:"foreignKey:CustomerID"`
}
// Unique constraint: (merchant_id, email) where email is not null

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
```

### 1.10 Promotions & Coupons

```go
// internal/models/promotion.go
type DiscountType string
const (
    DiscountPercentage DiscountType = "PERCENTAGE"
    DiscountFixed      DiscountType = "FIXED_AMOUNT"
    DiscountBuyXGetY   DiscountType = "BUY_X_GET_Y"
)

type Coupon struct {
    TenantModel
    Code               string         `gorm:"size:50;not null"`
    Description        string
    DiscountType       DiscountType   `gorm:"size:20;not null"`
    DiscountValue      float64        `gorm:"type:decimal(10,2);not null"`
    MinOrderValue      *float64       `gorm:"type:decimal(10,2)"`
    MaxDiscount        *float64       `gorm:"type:decimal(10,2)"`
    UsageLimit         *int
    UsedCount          int            `gorm:"default:0"`
    PerCustomerLimit   int            `gorm:"default:1"`
    ValidFrom          time.Time      `gorm:"not null"`
    ValidUntil         time.Time      `gorm:"not null"`
    IsActive           bool           `gorm:"default:true"`
    ApplicableCategories pq.StringArray `gorm:"type:uuid[]"`
    ApplicableProducts   pq.StringArray `gorm:"type:uuid[]"`
}
// Unique constraint: (merchant_id, code)
```

---

## 2. GraphQL Schema Definitions

### 2.1 Root Schema (`schema.graphqls`)

```graphql
scalar Time
scalar UUID
scalar Decimal
scalar JSON

directive @auth(permissions: [String!]!) on FIELD_DEFINITION

type Query {
    # Merchant Bootstrap (called on app startup by frontends)
    merchantBootstrap: MerchantConfig!

    # Products
    products(filter: ProductFilter, pagination: PaginationInput): ProductConnection!
    product(slug: String!): Product
    productByBarcode(barcode: String!): ProductVariant

    # Size Groups (for product creation forms)
    sizeGroups: [SizeGroup!]!
    sizeGroup(id: UUID!): SizeGroup

    # Inventory
    inventoryLevels(locationId: UUID, variantId: UUID): [InventoryLevel!]!

    # Orders
    orders(filter: OrderFilter, pagination: PaginationInput): OrderConnection!
    order(id: UUID!): Order

    # Customers
    customers(search: String, pagination: PaginationInput): CustomerConnection!
    customer(id: UUID!): Customer

    # Staff
    staff(pagination: PaginationInput): [User!]!
    roles: [Role!]!

    # Analytics
    dailySalesSummary(date: Time!, locationId: UUID): SalesSummary!
}

type Mutation {
    # Auth
    login(input: LoginInput!): AuthPayload!
    loginWithPin(pin: String!): AuthPayload!
    refreshToken(token: String!): AuthPayload!

    # Products
    createProduct(input: CreateProductInput!): Product! @auth(permissions: ["products:create"])
    updateProduct(id: UUID!, input: UpdateProductInput!): Product! @auth(permissions: ["products:edit"])
    createVariant(input: CreateVariantInput!): ProductVariant! @auth(permissions: ["products:create"])

    # Size Groups
    createSizeGroup(input: CreateSizeGroupInput!): SizeGroup! @auth(permissions: ["products:create"])
    assignSizeGroupToCategory(categoryId: UUID!, sizeGroupId: UUID!): Category! @auth(permissions: ["products:edit"])

    # Inventory
    adjustStock(input: AdjustStockInput!): InventoryLevel! @auth(permissions: ["inventory:edit"])
    transferStock(input: TransferStockInput!): Boolean! @auth(permissions: ["inventory:transfer"])

    # Orders
    createOrder(input: CreateOrderInput!): Order! @auth(permissions: ["orders:create"])
    updateOrderStatus(id: UUID!, status: OrderStatus!): Order! @auth(permissions: ["orders:edit"])
    refundOrder(input: RefundInput!): Refund! @auth(permissions: ["orders:refund"])

    # Customers
    createCustomer(input: CreateCustomerInput!): Customer! @auth(permissions: ["customers:create"])
    updateCustomer(id: UUID!, input: UpdateCustomerInput!): Customer! @auth(permissions: ["customers:edit"])

    # Promotions
    createCoupon(input: CreateCouponInput!): Coupon! @auth(permissions: ["promotions:create"])
    applyCoupon(orderId: UUID!, code: String!): Order!

    # Staff & RBAC
    createStaff(input: CreateStaffInput!): User! @auth(permissions: ["staff:manage"])
    assignRole(userId: UUID!, roleId: UUID!): User! @auth(permissions: ["staff:manage"])
    managerOverride(input: ManagerOverrideInput!): ManagerOverride! @auth(permissions: ["pos:override"])

    # POS
    openCashDrawer(locationId: UUID!): Boolean! @auth(permissions: ["pos:open_drawer"])
    endOfDayReport(locationId: UUID!): EODReport! @auth(permissions: ["pos:end_of_day"])
}

type Subscription {
    inventoryUpdated(locationId: UUID): InventoryEvent!
    orderPlaced: Order!
}
```

### 2.2 Product & SizeGroup Types (`product.graphqls`)

```graphql
type Product {
    id: UUID!
    name: String!
    slug: String!
    description: String
    basePrice: Decimal!
    comparePrice: Decimal
    brand: Brand
    category: Category
    images: [ProductImage!]!
    variants: [ProductVariant!]!
    tags: [String!]
    isPublished: Boolean!
    createdAt: Time!
}

type ProductVariant {
    id: UUID!
    product: Product!
    sku: String!
    barcode: String
    size: Size!
    color: Color!
    priceOverride: Decimal
    effectivePrice: Decimal!
    inventoryLevels: [InventoryLevel!]!
    totalStock: Int!
    isActive: Boolean!
}

type SizeGroup {
    id: UUID!
    name: String!
    sizes: [Size!]!
}

type Size {
    id: UUID!
    name: String!
    sortOrder: Int!
    sizeGroup: SizeGroup!
}

type Color {
    id: UUID!
    name: String!
    hexCode: String
}

type Category {
    id: UUID!
    name: String!
    slug: String!
    parent: Category
    children: [Category!]!
    sizeGroup: SizeGroup          # the sizes available when creating products in this category
}

type Brand {
    id: UUID!
    name: String!
    slug: String!
    logoUrl: String
}

input CreateSizeGroupInput {
    name: String!
    sizes: [String!]!             # e.g. ["S", "M", "L", "XL", "XXL"]
}

input CreateProductInput {
    categoryId: UUID
    brandId: UUID
    name: String!
    description: String
    basePrice: Decimal!
    comparePrice: Decimal
    costPrice: Decimal
    tags: [String!]
    variants: [CreateVariantInput!]!
}

input CreateVariantInput {
    productId: UUID               # optional if nested in CreateProductInput
    sizeId: UUID!
    colorId: UUID!
    sku: String!
    barcode: String
    priceOverride: Decimal
}
```

### 2.3 Merchant Types (`merchant.graphqls`)

```graphql
type MerchantConfig {
    id: UUID!
    businessName: String!
    slug: String!
    logoUrl: String
    primaryColor: String
    secondaryColor: String
    currency: String!
    timezone: String!
    taxRate: Decimal!
    settings: MerchantSettings!
    locations: [Location!]!
    sizeGroups: [SizeGroup!]!
}

type MerchantSettings {
    loyaltyEnabled: Boolean!
    loyaltyPointsPerUnit: Int!
    loyaltyRedemptionRate: Decimal!
    paymentGateways: [String!]!
    receiptHeader: String
    receiptFooter: String
}

type Location {
    id: UUID!
    name: String!
    address: String
    city: String
    isWarehouse: Boolean!
    isActive: Boolean!
}
```

### 2.4 Order, Payment, Inventory, Customer Types

*(These remain the same as defined in the previous version — Order, OrderItem, Payment, Refund, InventoryLevel, InventoryMovement, Customer, etc. See the full GraphQL schema files for details.)*

---

## 3. Entity Relationship Summary

```
merchants (1 row per tenant — runtime config store)
    │
    ├── users (staff)
    │   ├── user_roles → roles → role_permissions → permissions
    │   └── user_locations → locations
    │
    ├── locations
    │
    ├── size_groups
    │   └── sizes
    │
    ├── categories (self-referencing hierarchy)
    │   └── → size_group (FK: which sizes this category uses)
    │
    ├── brands
    │
    ├── colors
    │
    ├── products
    │   ├── product_images
    │   └── product_variants
    │       ├── → sizes (from the category's size_group)
    │       ├── → colors
    │       └── inventory_levels → locations
    │
    ├── customers
    │   ├── customer_addresses
    │   └── loyalty_transactions
    │
    ├── orders
    │   ├── order_items → product_variants
    │   └── payments → refunds
    │
    ├── coupons
    ├── inventory_movements
    └── manager_overrides
```

---

## 4. Key Constraints & Invariants

| Rule | Enforcement |
|------|-------------|
| Stock cannot go negative | GORM check constraint: `check:quantity >= 0` |
| SKU globally unique | `uniqueIndex` on `product_variants.sku` |
| One variant per size×color per product | Composite unique index `(product_id, size_id, color_id)` |
| Email unique per merchant | Composite unique index `(merchant_id, email)` on users and customers |
| Order prices are immutable after creation | Snapshotted into `order_items` columns |
| Loyalty points cannot go negative | GORM check constraint: `check:loyalty_points >= 0` |
| All tenant queries are scoped | `BaseRepository.Scoped(ctx)` applies `merchant_id` filter |
| Categories inherit sizing from SizeGroup | `category.size_group_id` FK determines available sizes |
