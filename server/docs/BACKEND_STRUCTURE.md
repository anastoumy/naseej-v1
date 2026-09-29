# Backend Directory Structure

## Naseej Core — Go GraphQL Multi-Tenant API Server

**Stack:** Go · Gin · gqlgen · GORM · PostgreSQL · WebSocket  
**Architecture:** Resolver → Service → Repository → GORM → Database  
**Multi-Tenancy:** Single server + single database, row-level isolation via `merchant_id`  
**Schema Management:** GORM AutoMigrate (no manual SQL migrations)

---

## Directory Layout

```
server/
├── cmd/
│   ├── server/
│   │   └── main.go                     # Entry point: init DB, AutoMigrate, middleware, routes, start server
│   └── seed/
│       └── main.go                     # Seed initial merchant config + admin user
│
├── config/
│   ├── config.go                       # Environment variable parsing & validation
│   └── config.yaml                     # Default configuration (overridden by env vars)
│
├── internal/
│   ├── models/                         # GORM model definitions (source of truth for DB schema)
│   │   ├── base.go                     # TenantModel embed (ID, MerchantID, timestamps)
│   │   ├── merchant.go                 # Merchant entity (tenant config, API keys)
│   │   ├── user.go                     # Staff & customer users
│   │   ├── product.go                  # Product + ProductVariant
│   │   ├── category.go                 # Hierarchical categories + SizeGroup link
│   │   ├── size_group.go               # SizeGroup + associated sizes
│   │   ├── inventory.go                # InventoryLevel + InventoryMovement
│   │   ├── order.go                    # Order + OrderItem
│   │   ├── payment.go                  # Payment + Refund
│   │   ├── customer.go                 # Customer + Address + LoyaltyTransaction
│   │   ├── promotion.go               # Coupon + PromotionRule
│   │   └── staff.go                    # Role + Permission + ManagerOverride
│   │
│   ├── repository/                     # Data access layer (GORM queries)
│   │   ├── base.go                     # BaseRepository with Scoped(ctx) method
│   │   ├── merchant_repo.go
│   │   ├── product_repo.go
│   │   ├── inventory_repo.go
│   │   ├── order_repo.go
│   │   ├── customer_repo.go
│   │   ├── payment_repo.go
│   │   └── staff_repo.go
│   │
│   ├── service/                        # Business logic layer
│   │   ├── product_service.go
│   │   ├── inventory_service.go        # Stock updates, movement logging, broadcast
│   │   ├── order_service.go            # Order creation, state transitions
│   │   ├── payment_service.go          # Payment processing, split payments, refunds
│   │   ├── customer_service.go         # Profile management, loyalty points
│   │   ├── promotion_service.go        # Coupon validation, discount application
│   │   ├── auth_service.go             # Login, JWT issuance, PIN verification
│   │   └── staff_service.go            # RBAC enforcement, role management
│   │
│   ├── graph/                          # GraphQL layer (gqlgen)
│   │   ├── schema/                     # Schema definition files (.graphqls)
│   │   │   ├── schema.graphqls         # Root Query, Mutation, Subscription types
│   │   │   ├── merchant.graphqls       # Merchant types & bootstrap query
│   │   │   ├── product.graphqls        # Product catalog types & operations
│   │   │   ├── inventory.graphqls      # Stock level types & mutations
│   │   │   ├── order.graphqls          # Order types & operations
│   │   │   ├── customer.graphqls       # Customer types & operations
│   │   │   ├── payment.graphqls        # Payment types
│   │   │   ├── promotion.graphqls      # Coupon & discount types
│   │   │   └── staff.graphqls          # RBAC types & operations
│   │   ├── generated/
│   │   │   └── generated.go            # Auto-generated gqlgen runtime (DO NOT EDIT)
│   │   ├── model/
│   │   │   └── models_gen.go           # Auto-generated GraphQL model structs
│   │   ├── resolver.go                 # Root resolver (dependency injection container)
│   │   ├── merchant.resolvers.go
│   │   ├── product.resolvers.go
│   │   ├── inventory.resolvers.go
│   │   ├── order.resolvers.go
│   │   ├── customer.resolvers.go
│   │   ├── payment.resolvers.go
│   │   ├── promotion.resolvers.go
│   │   └── staff.resolvers.go
│   │
│   ├── middleware/                     # Gin middleware stack
│   │   ├── tenant.go                   # X-Merchant-ID header → context injection
│   │   ├── auth.go                     # JWT validation → user context
│   │   ├── rbac.go                     # Permission enforcement per resolver
│   │   ├── cors.go                     # CORS configuration
│   │   └── logger.go                   # Structured request logging
│   │
│   └── websocket/                      # Real-time inventory broadcast
│       ├── hub.go                      # Connection manager (per-merchant rooms)
│       ├── client.go                   # Individual WebSocket connection handler
│       └── events.go                   # Event type definitions (stock_update, order_placed)
│
├── pkg/                                # Shared utilities (importable by other packages)
│   ├── database/
│   │   └── postgres.go                 # Connection pool setup, GORM initialization
│   ├── errors/
│   │   └── errors.go                   # Standardized error types & codes
│   └── utils/
│       ├── pagination.go               # Cursor-based pagination helpers
│       └── validation.go               # Input validation utilities
│
├── gqlgen.yml                          # gqlgen configuration (schema paths, resolver mapping)
├── Dockerfile                          # Multi-stage production build
├── docker-compose.yml                  # Local development (API + PostgreSQL)
├── Makefile                            # Common commands (generate, run, test)
├── go.mod
└── go.sum
```

---

## Architectural Patterns

### Three-Layer Architecture

Every data operation flows through three clean layers:

```
GraphQL Resolver (input validation, response mapping)
        │
        ▼
   Service Layer (business rules, orchestration, event emission)
        │
        ▼
  Repository Layer (GORM queries, scoped to tenant)
        │
        ▼
    PostgreSQL (data persistence)
```

**Rules:**
- Resolvers never call repositories directly
- Services never import `graph/` packages
- Repositories never contain business logic
- All layers receive `context.Context` carrying tenant & user info

### GORM AutoMigrate (Schema Management)

The server uses GORM's `AutoMigrate` to manage the database schema. On startup, `main.go` calls `AutoMigrate` with all model structs. GORM creates tables, adds missing columns, and creates indexes automatically. No manual SQL migration files are needed.

```go
// cmd/server/main.go (during initialization)
func autoMigrate(db *gorm.DB) error {
    return db.AutoMigrate(
        &models.Merchant{},
        &models.User{},
        &models.Role{},
        &models.Permission{},
        &models.Location{},
        &models.SizeGroup{},
        &models.Size{},
        &models.Color{},
        &models.Category{},
        &models.Brand{},
        &models.Product{},
        &models.ProductImage{},
        &models.ProductVariant{},
        &models.InventoryLevel{},
        &models.InventoryMovement{},
        &models.Customer{},
        &models.CustomerAddress{},
        &models.Order{},
        &models.OrderItem{},
        &models.Payment{},
        &models.Refund{},
        &models.Coupon{},
        &models.LoyaltyTransaction{},
        &models.ManagerOverride{},
    )
}
```

**Why AutoMigrate for now:**
- The app is in early development — schema changes are frequent
- No production data to protect yet
- When the app matures and has production merchants, switch to versioned migrations (golang-migrate)

### Merchant Table (Configuration & Tenant Identity)

The `merchants` table is the core tenant identity. Each merchant row holds runtime configuration that the frontends load on boot via `merchantBootstrap` query:

```go
// internal/models/merchant.go
type Merchant struct {
    ID              uuid.UUID       `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    BusinessName    string          `gorm:"size:150;not null"`
    Slug            string          `gorm:"size:50;uniqueIndex;not null"`
    Domain          string          `gorm:"size:255"`
    LogoURL         string
    TaxRate         float64         `gorm:"type:decimal(5,4);default:0.15"`
    Currency        string          `gorm:"size:3;default:'SAR'"`
    Timezone        string          `gorm:"size:50;default:'Asia/Riyadh'"`
    IsActive        bool            `gorm:"default:true"`
    Settings        datatypes.JSON  `gorm:"type:jsonb;default:'{}'"`  // payment gateways, loyalty, receipt text
    CreatedAt       time.Time
    UpdatedAt       time.Time
}
```

### Multi-Tenancy: The Scoped Context Pattern

All tenant-owned tables embed `TenantModel` which includes `MerchantID`. The `BaseRepository.Scoped(ctx)` method automatically filters every query by the current merchant.

```go
// internal/models/base.go
type TenantModel struct {
    ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    MerchantID uuid.UUID `gorm:"type:uuid;not null;index"`
    CreatedAt  time.Time
    UpdatedAt  time.Time
}

// internal/repository/base.go
type BaseRepository struct {
    db *gorm.DB
}

// Scoped returns a GORM session pre-filtered to the current tenant.
// This method MUST be used for every query to prevent cross-tenant data leaks.
func (r *BaseRepository) Scoped(ctx context.Context) *gorm.DB {
    merchantID := middleware.GetMerchantID(ctx)
    return r.db.Where("merchant_id = ?", merchantID)
}
```

**Tenant resolution flow:**
1. Each frontend (web/POS) sends `X-Merchant-ID` header with every request (configured in `merchant.ts`)
2. `middleware/tenant.go` validates the merchant exists and is active
3. Merchant ID injected into `context.Context`
4. Every repository call uses `Scoped(ctx)` — tenant filter is automatic and unforgettable

**Note:** Since each merchant’s website/POS is a separate deployment with a hardcoded `merchant_id` in `merchant.ts`, there is no dynamic tenant discovery. The frontend always knows which merchant it belongs to.

### RBAC Enforcement

```go
// Permissions follow resource:action naming
const (
    PermOrdersCreate  = "orders:create"
    PermOrdersRefund  = "orders:refund"
    PermInventoryEdit = "inventory:edit"
    PermPOSOpenDrawer = "pos:open_drawer"
    PermStaffManage   = "staff:manage"
)

// Built-in role hierarchy (each inherits permissions of roles below it):
// super_admin → admin → store_manager → senior_cashier → cashier
//                     → ecommerce_manager
//                     → inventory_manager
//                     → viewer
```

### WebSocket Real-Time Broadcast

```go
// When inventory changes, the service emits to the hub:
func (s *InventoryService) UpdateStock(ctx context.Context, variantID uuid.UUID, delta int) error {
    // ... update DB ...
    
    merchantID := middleware.GetMerchantID(ctx)
    // Broadcast to all connected clients for this merchant
    s.wsHub.Broadcast(merchantID, events.StockUpdate{
        VariantID: variantID,
        NewQty:    newLevel.Quantity,
        Location:  locationID,
    })
    return nil
}
```

---

## Key Design Decisions

| Decision | Rationale |
|----------|-----------|
| GORM over raw SQL | Rapid development, AutoMigrate for early-stage schema management, struct-based models |
| gqlgen over graphql-go | Code-first generation, type safety, schema-driven development |
| Single DB multi-tenant | Simpler operations, one server to maintain, row-level isolation sufficient |
| GORM AutoMigrate | Early stage — schema evolving rapidly; will switch to versioned migrations at maturity |
| Gin over net/http | Middleware ecosystem, performance, request context helpers |
| Domain-split schemas | Independent domain evolution, smaller resolver files |
| WebSocket over SSE | Bidirectional communication needed for POS acknowledgments |
