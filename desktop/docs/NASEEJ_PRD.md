# Product Requirements Document (PRD)

## Naseej — Omni-Channel Retail Operating Suite for Clothing Merchants

**Version:** 1.0  
**Last Updated:** 2026-05-18  
**Status:** Active Development

---

## 1. Executive Summary

Naseej is a production-grade, multi-tenant omni-channel retail operating suite purpose-built for clothing merchants. The platform bridges digital e-commerce and brick-and-mortar operations through a unified backend, ensuring real-time inventory synchronization across all sales channels with zero drift.

The system uses a **single shared Go GraphQL server** and **single shared PostgreSQL database** that serves all merchants (row-level isolation via `merchant_id`). Each merchant receives their **own dedicated repository** for the web storefront and POS application, deployed on their **own custom domain**. Onboarding a new merchant means creating a merchant row in the DB, then building out their frontend repos from templates.

---

## 2. Target Users & Personas

| Persona | Role | Primary Need |
|---------|------|--------------|
| **Cashier** | In-store POS operator | Fast, keyboard/barcode-driven checkout with hardware integration |
| **Store Manager** | Multi-location supervisor | Real-time dashboards, inventory oversight, staff management |
| **E-Commerce Manager** | Online operations lead | Catalog management, order fulfillment, promotions |
| **Inventory Manager** | Stock controller | Cross-location stock visibility, movement tracking, reorder alerts |
| **End Consumer** | Online shopper | Responsive, SEO-optimized storefront with real-time availability |
| **Merchant Owner** | Business stakeholder | Unified reporting, multi-channel analytics, cost reduction |

---

## 3. System Architecture Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                        Naseej Ecosystem                          │
├─────────────────────────────────────────────────────────────────┤
│                                                                   │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐       │
│  │  Naseej Web  │    │  Naseej POS  │    │  Mobile App  │       │
│  │  (Next.js)   │    │  (Electron)  │    │  (Future)    │       │
│  └──────┬───────┘    └──────┬───────┘    └──────┬───────┘       │
│         │                    │                    │               │
│         └────────────┬───────┴────────────┬──────┘               │
│                      │                    │                       │
│              ┌───────▼────────┐   ┌──────▼───────┐              │
│              │  GraphQL API   │   │  WebSocket   │              │
│              │  (Queries &    │   │  (Real-time  │              │
│              │   Mutations)   │   │   Updates)   │              │
│              └───────┬────────┘   └──────┬───────┘              │
│                      │                    │                       │
│              ┌───────▼────────────────────▼───────┐              │
│              │         Naseej Core (Go)           │              │
│              │   Gin + gqlgen + GORM + WebSocket  │              │
│              └───────────────┬────────────────────┘              │
│                              │                                    │
│              ┌───────────────▼────────────────────┐              │
│              │     PostgreSQL (Multi-Tenant)       │              │
│              │  Single DB · merchant_id row-level   │              │
│              └────────────────────────────────────┘              │
│                                                                   │
└─────────────────────────────────────────────────────────────────┘
```

**Technology Stack:**

| Layer | Technology |
|-------|-----------|
| Backend Server | Go + Gin + gqlgen (GraphQL) + GORM |
| Database | PostgreSQL (single DB, `merchant_id` scopes all data) |
| Schema Management | GORM AutoMigrate (no manual SQL migrations) |
| Web Storefront | Next.js (App Router) — separate repo per merchant, own custom domain |
| Desktop POS | Electron (Vite + React) — separate repo per merchant |
| Real-time | WebSockets (inventory broadcast, per-merchant rooms) |

---

## 4. Functional Requirements

### 4.1 Product Catalog & Variant Matrix

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-001 | Products follow a base → variant hierarchy (Product → ProductVariant) | P0 |
| FR-002 | Each variant is the atomic sellable unit, identified by SKU and barcode | P0 |
| FR-003 | Variants are defined by a Size × Color matrix per base product | P0 |
| FR-003a | Categories are linked to **Size Groups** — a named set of sizes relevant to that product type (e.g. "Shoes Women 36–41", "Shirts S–XXL") | P0 |
| FR-003b | When creating a product under a category, the variant form auto-populates with the category's assigned size group sizes | P0 |
| FR-004 | Products belong to hierarchical categories (parent → child) | P1 |
| FR-005 | Brand attribution and filtering on product catalog | P1 |
| FR-006 | Product images with variant-level gallery support | P1 |

### 4.2 Unified Inventory Management

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-010 | Inventory tracked per variant per physical location (`InventoryLevel`) | P0 |
| FR-011 | Stock can never go negative — enforced at database constraint level | P0 |
| FR-012 | All stock changes recorded as `InventoryMovement` with reason codes | P0 |
| FR-013 | Real-time stock broadcast via WebSocket within < 1.5 seconds of change | P0 |
| FR-014 | Cross-location stock lookup from POS (check other branches) | P1 |
| FR-015 | Low-stock threshold alerts per variant per location | P1 |
| FR-016 | Stock adjustment workflows (damaged, lost, manual count) | P1 |

### 4.3 Order Management (Unified Model)

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-020 | Single `Order` model serves both POS and web channels | P0 |
| FR-021 | Channel enum distinguishes order origin (POS, WEB, MOBILE) | P0 |
| FR-022 | Prices snapshotted on `OrderItem` at time of sale (immutable) | P0 |
| FR-023 | Order lifecycle states: PENDING → CONFIRMED → PROCESSING → SHIPPED → DELIVERED → CANCELLED | P0 |
| FR-024 | Returns and exchanges with line-item level granularity | P1 |
| FR-025 | Order search and filtering by date, channel, status, customer | P1 |

### 4.4 Payment Processing

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-030 | Multiple payment methods per order (split payments) | P0 |
| FR-031 | Payment types: CASH, CARD, BANK_TRANSFER, LOYALTY_POINTS | P0 |
| FR-032 | Refund tracking linked to original payment | P1 |
| FR-033 | POS cash drawer reconciliation (expected vs. actual) | P1 |
| FR-034 | Integration with payment terminal hardware via IPC | P1 |

### 4.5 Customer Management

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-040 | Customer profiles with multiple addresses | P0 |
| FR-041 | Customer tagging for segmentation | P1 |
| FR-042 | Loyalty points system (earn on purchase, redeem at checkout) | P1 |
| FR-043 | Purchase history accessible from POS and web | P1 |
| FR-044 | Customer search by phone, email, or loyalty card barcode | P1 |

### 4.6 Promotions & Discounts

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-050 | Coupon codes with configurable conditions (min spend, date range, usage limits) | P1 |
| FR-051 | Line-item and cart-wide discount application | P1 |
| FR-052 | Loyalty point transactions (earn and burn) | P1 |
| FR-053 | Automatic promotions (buy X get Y, percentage off category) | P2 |

### 4.7 Staff & Access Control (RBAC)

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-060 | Permission model follows `resource:action` naming convention | P0 |
| FR-061 | Built-in role hierarchy: `super_admin` → `admin` → `store_manager` → `senior_cashier` → `cashier` → `ecommerce_manager` → `inventory_manager` → `viewer` | P0 |
| FR-062 | Roles are additive — a user may hold multiple roles | P0 |
| FR-063 | Location scoping is independent from permission grants | P0 |
| FR-064 | Sensitive POS actions require manager PIN override (logged in `ManagerOverride`) | P0 |
| FR-065 | Staff login via PIN (POS) or email/password (web admin) | P0 |

### 4.8 POS Hardware Integration

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-070 | Barcode scanner input via global keyboard listener (invisible interceptor) | P0 |
| FR-071 | ESC/POS thermal receipt printing via IPC bridge | P0 |
| FR-072 | Cash drawer open trigger (electronic signal) | P1 |
| FR-073 | Customer-facing display support | P2 |

### 4.9 Multi-Tenancy & Deployment Model

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-080 | Single shared server + single PostgreSQL database serves all merchants | P0 |
| FR-081 | Row-level tenant isolation via `merchant_id` on all data tables | P0 |
| FR-082 | Tenant resolution from `X-Merchant-ID` header in middleware | P0 |
| FR-083 | `BaseRepository.Scoped(ctx)` auto-applies merchant filter on every query | P0 |
| FR-084 | Merchant bootstrap query returns runtime config (tax, payment gateways, loyalty settings) | P0 |
| FR-085 | Each merchant’s web storefront is a separate repo deployed on its own custom domain | P0 |
| FR-086 | Each merchant’s POS app is a separate repo with hardcoded `merchant_id` in `merchant.ts` | P0 |
| FR-087 | Static branding (colors, fonts, logo) applied immediately from `merchant.ts` as CSS variables | P0 |

---

## 5. Non-Functional Requirements

| Category | Requirement | Target |
|----------|-------------|--------|
| **Performance** | GraphQL query response time (p95) | ≤ 120ms |
| **Performance** | WebSocket inventory broadcast latency | < 1.5 seconds |
| **Performance** | POS checkout transaction (scan → receipt) | < 3 seconds |
| **Scalability** | Concurrent merchants on shared server | 500+ |
| **Scalability** | Products per merchant | 100,000+ |
| **Availability** | API uptime SLA | 99.9% |
| **Security** | Data isolation between merchants | Row-level via `merchant_id` (enforced at ORM layer) |
| **Security** | Authentication | JWT with refresh tokens |
| **Security** | POS sensitive actions | Manager PIN override with audit log |
| **Compatibility** | POS display resolution | 1024×768 minimum |
| **Compatibility** | Web storefront | All modern browsers, mobile responsive |

---

## 6. Onboarding Flow (Per-Client Deployment)

The shared server is always running. Onboarding adds data to the existing DB and deploys new frontends:

```
1. INSERT merchant row into shared PostgreSQL  → receive merchant_id
2. Seed initial data (admin user, roles, permissions, locations)
3. Create new website repo for merchant         → configure merchant.ts (id + theme + API URL)
4. Create new POS repo for merchant             → configure merchant.ts (id + theme + API URL)
5. Deploy website on merchant’s custom domain   → www.merchantbrand.com
6. Package POS app for merchant                 → branded Electron installer
```

**The central API server is already running — no server redeployment needed.**  
**No subdomains. No slug-based routing. Each merchant’s website is a standalone app on its own domain.**

---

## 7. Success Metrics

| Metric | Target |
|--------|--------|
| Inventory drift rate (online vs. physical) | 0% |
| Average POS transaction time | < 15 seconds |
| Merchant onboarding time (technical) | < 30 minutes |
| System downtime per month | < 45 minutes |
| WebSocket message delivery success | > 99.5% |

---

## 8. Out of Scope (v1.0)

- Native mobile applications (iOS/Android)
- Multi-currency support
- Warehouse management (pick/pack/ship workflows)
- Marketplace integrations (Amazon, Noon, etc.)
- AI-powered demand forecasting
- Multi-language storefront (i18n)
