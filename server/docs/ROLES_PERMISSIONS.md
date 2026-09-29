# Naseej — Roles & Permissions

## Overview

Naseej uses a **Role-Based Access Control (RBAC)** system. Every staff member is assigned one or more roles, and each role carries a set of permissions. Permissions are **additive** — a staff member with multiple roles gets the union of all permissions across those roles.

There are two orthogonal access controls:

- **Permission** — what actions a staff member can perform
- **Location Scope** — which store branches they can see data from

---

## Permission Naming Convention

All permissions follow the pattern `resource:action`.

```
orders:refund
pos:open_drawer
inventory:transfer
```

---

## Full Permission List

### Products

| Permission | Description |
|---|---|
| `products:view` | View products and variants |
| `products:create` | Create new products |
| `products:edit` | Edit existing products |
| `products:delete` | Delete products |
| `products:publish` | Toggle active / featured status |

### Inventory

| Permission | Description |
|---|---|
| `inventory:view` | View stock levels across locations |
| `inventory:edit` | Adjust stock counts manually |
| `inventory:transfer` | Move stock between locations |
| `inventory:price` | Edit cost price, retail price, compare-at price |

### Orders

| Permission | Description |
|---|---|
| `orders:view` | View orders |
| `orders:create` | Create manual orders |
| `orders:edit` | Edit order items, address, notes |
| `orders:cancel` | Cancel an order |
| `orders:refund` | Issue full or partial refunds |
| `orders:discount` | Apply manual discounts to an order |
| `orders:export` | Export orders to CSV |

### POS

| Permission | Description |
|---|---|
| `pos:access` | Open and use the POS at all |
| `pos:open_drawer` | Open the cash drawer |
| `pos:void_item` | Remove a line item from an active sale |
| `pos:override_price` | Manually override a variant's price mid-sale |
| `pos:apply_coupon` | Apply a coupon code to a sale |
| `pos:park_sale` | Hold / park a transaction |
| `pos:cash_in` | Record cash added to the drawer |
| `pos:cash_out` | Record cash removed from the drawer |
| `pos:end_shift` | Reconcile and close own register shift |
| `pos:end_shift_other` | Close another cashier's open shift |

### Customers

| Permission | Description |
|---|---|
| `customers:view` | View customer profiles and history |
| `customers:create` | Create new customer records |
| `customers:edit` | Edit customer details and addresses |
| `customers:delete` | Delete customer records |
| `customers:export` | Export customer data to CSV |
| `customers:loyalty` | Manually adjust loyalty points balance |

### Coupons & Promotions

| Permission | Description |
|---|---|
| `coupons:view` | View coupons and promotions |
| `coupons:create` | Create new coupons |
| `coupons:edit` | Edit existing coupons |
| `coupons:delete` | Delete coupons |

### Reports & Analytics

| Permission | Description |
|---|---|
| `reports:sales` | View sales reports |
| `reports:inventory` | View inventory reports |
| `reports:customers` | View customer reports |
| `reports:financial` | View financial reports (margins, costs) — sensitive |

### Staff Management

| Permission | Description |
|---|---|
| `staff:view` | View staff list and profiles |
| `staff:create` | Create new staff accounts |
| `staff:edit` | Edit staff details and PIN |
| `staff:delete` | Delete staff accounts |
| `staff:assign_role` | Assign or remove roles from staff |

### Settings

| Permission | Description |
|---|---|
| `settings:view` | View store settings |
| `settings:edit` | Edit taxes, shipping, payment methods |
| `settings:locations` | Add and manage store branches |

---

## Built-in Roles

### Super Admin
> Full access to everything. Typically the business owner.

All permissions granted. Cannot be restricted.

---

### Admin
> Manages the full operation. Cannot manage other admins or assign roles.

All permissions **except**: `staff:delete`, `staff:assign_role`.

---

### Store Manager
> Runs day-to-day operations at one or more locations. Can process refunds, override prices, end shifts.

```
products:view, products:edit, products:publish
inventory:view, inventory:edit, inventory:transfer
orders:view, orders:create, orders:edit, orders:cancel, orders:refund, orders:discount, orders:export
pos:access, pos:open_drawer, pos:void_item, pos:override_price, pos:apply_coupon
pos:park_sale, pos:cash_in, pos:cash_out, pos:end_shift, pos:end_shift_other
customers:view, customers:create, customers:edit, customers:loyalty
coupons:view, coupons:create, coupons:edit
reports:sales, reports:inventory, reports:customers
staff:view
settings:view
```

---

### Senior Cashier
> Trusted front-line staff. Can handle returns, open the drawer, void items, and close their own shift.

```
pos:access, pos:open_drawer, pos:void_item, pos:apply_coupon
pos:park_sale, pos:cash_in, pos:cash_out, pos:end_shift
orders:view, orders:create, orders:cancel
customers:view, customers:create, customers:edit
inventory:view
```

---

### Cashier
> Standard front-line POS staff. Can sell and apply coupons. Cannot override prices or open the drawer without a manager.

```
pos:access, pos:apply_coupon, pos:park_sale
orders:view, orders:create
customers:view, customers:create
inventory:view
```

---

### E-commerce Manager
> Manages the online storefront. No POS access.

```
products:view, products:create, products:edit, products:publish
inventory:view
orders:view, orders:edit, orders:cancel, orders:refund, orders:export
customers:view, customers:create, customers:edit, customers:export
coupons:view, coupons:create, coupons:edit, coupons:delete
reports:sales, reports:customers
```

---

### Inventory Manager
> Warehouse or stockroom staff. Manages stock levels and transfers. No financials, no POS.

```
products:view, products:create, products:edit
inventory:view, inventory:edit, inventory:transfer
reports:inventory
settings:view
```

---

### Viewer / Auditor
> Read-only access. For external accountants, silent partners, or auditors.

```
products:view
inventory:view
orders:view, orders:export
customers:view
reports:sales, reports:inventory, reports:customers, reports:financial
settings:view
```

---

## Role Matrix

| Permission | Super Admin | Admin | Store Manager | Inventory Mgr | Cashier | Senior Cashier | Ecom Manager | Viewer |
|---|:---:|:---:|:---:|:---:|:---:|:---:|:---:|:---:|
| `products:view` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `products:create` | ✅ | ✅ | ❌ | ✅ | ❌ | ❌ | ✅ | ❌ |
| `products:edit` | ✅ | ✅ | ✅ | ✅ | ❌ | ❌ | ✅ | ❌ |
| `products:delete` | ✅ | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
| `products:publish` | ✅ | ✅ | ✅ | ❌ | ❌ | ❌ | ✅ | ❌ |
| `inventory:view` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `inventory:edit` | ✅ | ✅ | ✅ | ✅ | ❌ | ❌ | ❌ | ❌ |
| `inventory:transfer` | ✅ | ✅ | ✅ | ✅ | ❌ | ❌ | ❌ | ❌ |
| `inventory:price` | ✅ | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
| `orders:view` | ✅ | ✅ | ✅ | ❌ | ✅ | ✅ | ✅ | ✅ |
| `orders:create` | ✅ | ✅ | ✅ | ❌ | ✅ | ✅ | ✅ | ❌ |
| `orders:edit` | ✅ | ✅ | ✅ | ❌ | ❌ | ❌ | ✅ | ❌ |
| `orders:cancel` | ✅ | ✅ | ✅ | ❌ | ❌ | ✅ | ✅ | ❌ |
| `orders:refund` | ✅ | ✅ | ✅ | ❌ | ❌ | ❌ | ✅ | ❌ |
| `orders:discount` | ✅ | ✅ | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ |
| `orders:export` | ✅ | ✅ | ✅ | ❌ | ❌ | ❌ | ✅ | ✅ |
| `pos:access` | ✅ | ✅ | ✅ | ❌ | ✅ | ✅ | ❌ | ❌ |
| `pos:open_drawer` | ✅ | ✅ | ✅ | ❌ | ❌ | ✅ | ❌ | ❌ |
| `pos:void_item` | ✅ | ✅ | ✅ | ❌ | ❌ | ✅ | ❌ | ❌ |
| `pos:override_price` | ✅ | ✅ | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ |
| `pos:apply_coupon` | ✅ | ✅ | ✅ | ❌ | ✅ | ✅ | ❌ | ❌ |
| `pos:park_sale` | ✅ | ✅ | ✅ | ❌ | ✅ | ✅ | ❌ | ❌ |
| `pos:cash_in` | ✅ | ✅ | ✅ | ❌ | ❌ | ✅ | ❌ | ❌ |
| `pos:cash_out` | ✅ | ✅ | ✅ | ❌ | ❌ | ✅ | ❌ | ❌ |
| `pos:end_shift` | ✅ | ✅ | ✅ | ❌ | ❌ | ✅ | ❌ | ❌ |
| `pos:end_shift_other` | ✅ | ✅ | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ |
| `customers:view` | ✅ | ✅ | ✅ | ❌ | ✅ | ✅ | ✅ | ✅ |
| `customers:create` | ✅ | ✅ | ✅ | ❌ | ✅ | ✅ | ✅ | ❌ |
| `customers:edit` | ✅ | ✅ | ✅ | ❌ | ❌ | ✅ | ✅ | ❌ |
| `customers:delete` | ✅ | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
| `customers:export` | ✅ | ✅ | ✅ | ❌ | ❌ | ❌ | ✅ | ✅ |
| `customers:loyalty` | ✅ | ✅ | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ |
| `coupons:view` | ✅ | ✅ | ✅ | ❌ | ❌ | ❌ | ✅ | ❌ |
| `coupons:create` | ✅ | ✅ | ✅ | ❌ | ❌ | ❌ | ✅ | ❌ |
| `coupons:edit` | ✅ | ✅ | ✅ | ❌ | ❌ | ❌ | ✅ | ❌ |
| `coupons:delete` | ✅ | ✅ | ❌ | ❌ | ❌ | ❌ | ✅ | ❌ |
| `reports:sales` | ✅ | ✅ | ✅ | ❌ | ❌ | ❌ | ✅ | ✅ |
| `reports:inventory` | ✅ | ✅ | ✅ | ✅ | ❌ | ❌ | ❌ | ✅ |
| `reports:customers` | ✅ | ✅ | ✅ | ❌ | ❌ | ❌ | ✅ | ✅ |
| `reports:financial` | ✅ | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ✅ |
| `staff:view` | ✅ | ✅ | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ |
| `staff:create` | ✅ | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
| `staff:edit` | ✅ | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
| `staff:delete` | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
| `staff:assign_role` | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
| `settings:view` | ✅ | ✅ | ✅ | ✅ | ❌ | ❌ | ❌ | ✅ |
| `settings:edit` | ✅ | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
| `settings:locations` | ✅ | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |

---

## Location Scoping

Location scoping is **separate** from permissions. A staff member may have `orders:refund` but only see orders from their assigned branches.

Each staff member has a `staff_locations` join table:

```sql
staff_locations
  staff_id    → FK to staff
  location_id → FK to locations
  role_slug   → optional per-location role override
```

This allows a staff member to be a **Manager at Branch A** but only a **Cashier at Branch B**.

All repository queries are automatically filtered:

```go
func (r *OrderRepository) List(ctx context.Context) ([]Order, error) {
    allowedLocations := r.getStaffLocations(ctx)
    return r.Scoped(ctx).
        Where("location_id IN ?", allowedLocations).
        Find(&orders)
}
```

---

## POS Manager Override

Sensitive POS actions (price override, refund, void) can require a **manager PIN confirmation** without logging the cashier out. Every override is logged.

```go
type ManagerOverride struct {
    gorm.Model
    MerchantID    uint
    RequestedByID uint      // cashier who triggered it
    ApprovedByID  uint      // manager who entered PIN
    Action        string    // "price_override" | "void_item" | "refund"
    ReferenceID   uint      // OrderID or OrderItemID
    Value         float64   // e.g. override price
    ApprovedAt    time.Time
}
```

---

## Implementation Notes

**Roles are additive, never subtractive.** A staff member with two roles gets the union of both permission sets. Do not design roles that remove permissions — split into narrower roles instead.

**Permissions are embedded in the JWT** at login time, avoiding a DB hit on every request.

**Audit log everything sensitive** — refunds, price overrides, stock adjustments, manager overrides, and role changes should all write to an `AuditLog` table with `staff_id`, `action`, `before`, `after`, `ip_address`, and `merchant_id`.

**POS PIN ≠ login password.** The 4–6 digit PIN is only for in-session manager override escalation. Full authentication uses email + password.