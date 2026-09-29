# Naseej — Implementation Strategy & Phased Plan

## Overview

This document outlines the phased implementation strategy for the Naseej retail operating suite. Each phase builds on the previous, delivering incremental value while maintaining a shippable product at each milestone.

**Total Estimated Phases:** 5  
**Architecture:** Go/Gin shared API server → Multiple independent frontends per merchant

---

## Phase 1: Foundation & Core Backend

**Goal:** Establish the shared backend infrastructure — database, authentication, multi-tenant middleware, and core data models.

### Deliverables

| # | Task | Details |
|---|------|---------|
| 1.1 | Project scaffolding | Go module, directory structure, Docker Compose (PostgreSQL), Makefile |
| 1.2 | GORM models & AutoMigrate | Merchant, User, Role, Permission, Location — server auto-creates tables on startup |
| 1.3 | Multi-tenant middleware | `X-Merchant-ID` header extraction → context injection → `Scoped(ctx)` base repository |
| 1.4 | Authentication | JWT issuance, refresh tokens, password hashing, role-based access guards |
| 1.5 | GraphQL skeleton | gqlgen setup, schema split by domain, generated resolvers |
| 1.6 | Seed command | `cmd/seed/main.go` — creates test merchant, admin user, default roles & permissions |
| 1.7 | CI pipeline | Lint + test + build on every push |

### Exit Criteria
- Server starts, creates all tables via AutoMigrate
- Login/logout works, JWT validated on protected queries
- Scoped queries return only data for the authenticated merchant
- Seed creates a fully usable test merchant

---

## Phase 2: Product Catalog & Inventory

**Goal:** Full product management with the SizeGroup system, categories, variants, and real-time stock tracking.

### Deliverables

| # | Task | Details |
|---|------|---------|
| 2.1 | SizeGroup & Size models | CRUD for SizeGroups, admin assigns SizeGroups to categories |
| 2.2 | Category tree | Hierarchical categories with `parent_id`, depth limit of 3 |
| 2.3 | Product & ProductVariant | Products linked to categories, variants auto-populated from category's SizeGroup |
| 2.4 | Inventory management | Stock levels per variant per location, stock movement log |
| 2.5 | Image upload | Product image upload to object storage (S3/R2), position ordering |
| 2.6 | GraphQL API complete | All queries/mutations for catalog browsing and admin management |
| 2.7 | Admin seed data | Sample SizeGroups, categories, products with variants and stock |

### Exit Criteria
- Admin creates a SizeGroup → assigns to category → creates product → variant form auto-populates sizes
- Stock decrements on sale, increments on restock
- Catalog browsable via GraphQL with filtering and pagination

---

## Phase 3: Web Storefront (First Merchant)

**Goal:** Launch-ready customer-facing e-commerce website for the first merchant.

### Deliverables

| # | Task | Details |
|---|------|---------|
| 3.1 | Next.js project setup | App Router, Tailwind, shadcn/ui, `merchant.ts` configuration |
| 3.2 | GraphQL client | Typed client with `X-Merchant-ID` header, ISR/SSR data fetching |
| 3.3 | Homepage & collections | Hero, featured products, collection grid, category filtering |
| 3.4 | Product detail page | Image gallery, variant selector (size/color), stock status, add-to-cart |
| 3.5 | Cart & checkout flow | Cart state, multi-step checkout (info → shipping → payment → confirmation) |
| 3.6 | Customer accounts | Registration, login, order history, profile management |
| 3.7 | Search | Full-text product search with filters |
| 3.8 | Responsive design | Mobile-first, tested across breakpoints |
| 3.9 | Deployment | Vercel/custom hosting on merchant's domain, production environment variables |

### Exit Criteria
- End-to-end customer journey: browse → add to cart → checkout → receive order confirmation
- Site loads on merchant's custom domain with their branding
- Lighthouse score > 90 (performance), SEO meta tags generated per page

---

## Phase 4: POS Application

**Goal:** Desktop point-of-sale app for in-store transactions.

### Deliverables

| # | Task | Details |
|---|------|---------|
| 4.1 | Electron + Vite + React setup | Project scaffolding, `merchant.ts` config, IPC bridge |
| 4.2 | Authentication | Staff login with PIN or email/password, role-based UI |
| 4.3 | Product lookup | Barcode scan, search, category browse — quick add to cart |
| 4.4 | Cart & sale flow | Line items, discounts, tax calculation, split payment |
| 4.5 | Payment integration | Cash, card terminal (mock/actual), change calculation |
| 4.6 | Receipt printing | Thermal printer support (ESC/POS), configurable receipt template |
| 4.7 | Offline mode | Local SQLite cache, queue transactions, sync when online |
| 4.8 | End-of-day | Register close, cash count, daily summary report |
| 4.9 | Real-time stock sync | WebSocket connection for live inventory updates across locations |

### Exit Criteria
- Complete sale cycle: scan → pay → print receipt → stock decrements
- Works offline, syncs when connectivity restored
- Runs on Windows & macOS

---

## Phase 5: Advanced Features & Scale

**Goal:** Analytics, loyalty, multi-merchant onboarding automation, and operational maturity.

### Deliverables

| # | Task | Details |
|---|------|---------|
| 5.1 | Analytics dashboard | Sales by period, top products, revenue per location, trend graphs |
| 5.2 | Loyalty program | Points earning/redemption, tier system, customer lifetime value |
| 5.3 | Promotions engine | Discount rules (% off, BOGO, minimum spend), scheduled campaigns |
| 5.4 | Merchant onboarding CLI | Automated: insert merchant → seed → scaffold frontend repos → deploy |
| 5.5 | Notifications | Order status (email/SMS), low-stock alerts, promotion broadcasts |
| 5.6 | Returns & exchanges | Return flow in POS and web, stock reversal, refund processing |
| 5.7 | Multi-location transfers | Stock transfer requests between locations, approval workflow |
| 5.8 | Audit logging | Complete audit trail for sensitive operations |
| 5.9 | Performance & monitoring | APM integration, error tracking, uptime monitoring |

### Exit Criteria
- Second merchant fully onboarded using automated tooling
- Analytics dashboard shows real-time and historical data
- Loyalty points earn and redeem across web and POS

---

## Cross-Cutting Concerns (All Phases)

| Concern | Approach |
|---------|----------|
| **Testing** | Unit tests for services, integration tests for resolvers, E2E for critical flows |
| **Security** | Input validation, SQL injection prevention (GORM parameterized), XSS protection, rate limiting |
| **Observability** | Structured logging (zerolog), request tracing, health endpoints |
| **Documentation** | API schema is self-documenting (GraphQL introspection), README per package |
| **Code Review** | All changes via PR, require 1 approval minimum |

---

## Dependency Graph

```
Phase 1 (Foundation)
    ↓
Phase 2 (Catalog & Inventory)
    ↓
  ┌─────────────┐
  ↓             ↓
Phase 3       Phase 4
(Web)         (POS)
  └─────────────┘
        ↓
    Phase 5
  (Advanced)
```

Phases 3 and 4 can run **in parallel** once Phase 2 is complete.

---

## Risk Mitigation

| Risk | Impact | Mitigation |
|------|--------|------------|
| Payment gateway delays | Blocks checkout/POS | Mock gateway for dev, integrate real gateway as late as possible |
| Offline sync conflicts | Data corruption in POS | Last-write-wins with conflict log for manual resolution |
| Performance at scale | Slow queries with many merchants | Database indexing on `merchant_id`, query analysis, connection pooling |
| Scope creep | Delayed delivery | Strict phase boundaries, defer features to Phase 5 |
| Single point of failure | API downtime affects all merchants | POS offline mode, CDN for static assets, health-check-based routing |
