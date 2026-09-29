# Business Requirements Document (BRD)

## Naseej — Omni-Channel Retail Operating Suite for Clothing Merchants

**Version:** 1.0  
**Last Updated:** 2026-05-18  
**Status:** Active Development

---

## 1. Executive Vision

Naseej democratizes modern omni-channel retail technology for clothing merchants. Traditional platforms either neglect the physical checkout experience or introduce delayed, error-prone inventory synchronization. Naseej eliminates inventory drift entirely by unifying all sales channels under a single real-time backend.

**Mission:** Deliver a turnkey, white-labeled retail operating suite that any clothing merchant can adopt within 30 minutes — no enterprise budgets, no custom development, no inventory discrepancies.

---

## 2. Problem Statement

| Challenge | Impact | Current Market Gap |
|-----------|--------|--------------------|
| Inventory drift between channels | Overselling, customer dissatisfaction, lost revenue | Most platforms sync on intervals (5–15 min), allowing ghost sales |
| High total cost of ownership | Small merchants priced out of modern retail tech | Separate POS, e-commerce, and inventory tools cost $500–2000/mo combined |
| Slow merchant onboarding | Weeks of setup, custom integrations | Competitors require custom development per merchant |
| Disconnected customer data | Fragmented loyalty, no unified purchase history | POS and web operate as isolated silos |

---

## 3. Core Value Propositions

### 3.1 Zero Drift Guarantee

Real-time WebSocket broadcasting ensures that the moment a variant sells (online or in-store), all connected systems reflect the change within 1.5 seconds. Database-level constraints prevent negative stock, making overselling architecturally impossible.

### 3.2 Unified Platform, Single Subscription

Merchants receive a complete retail stack — POS terminal software, e-commerce storefront, inventory management, customer CRM, and analytics — under one subscription. No piecing together third-party tools.

### 3.3 Effortless Merchant Onboarding

A single central server handles all merchants. Adding a new client requires only a database insert and deploying pre-built frontend repos on the merchant’s own custom domain. No server changes. No subdomain routing. Each merchant’s storefront and POS are standalone projects connected to the shared API.

### 3.4 Clothing-First Design

The product matrix (Size × Color variants), visual merchandising tools, and category hierarchy are purpose-built for apparel. Not a generic platform retrofitted for fashion.

---

## 4. Target Market

### 4.1 Primary Segment

- **Mid-market clothing retailers** operating 1–50 physical locations
- Annual revenue: $500K – $50M
- Currently using disconnected POS + basic e-commerce (Shopify, WooCommerce)
- Pain point: Manual inventory reconciliation, overselling, fragmented reporting

### 4.2 Secondary Segment

- **Emerging fashion brands** launching both online and physical retail simultaneously
- Need a unified platform from day one
- Price-sensitive, require predictable subscription costs

### 4.3 Geographic Focus (Initial)

- MENA region (Saudi Arabia, UAE, Egypt, Jordan)
- Localization: RTL support, local payment gateways, VAT compliance

---

## 5. Revenue Architecture

### 5.1 SaaS Subscription Tiers

| Tier | Monthly Price | Locations | Variants | Users |
|------|--------------|-----------|----------|-------|
| **Starter** | $99 | 1 | 5,000 | 5 |
| **Growth** | $299 | 5 | 25,000 | 20 |
| **Enterprise** | $799 | Unlimited | Unlimited | Unlimited |

### 5.2 Additional Revenue Streams

| Stream | Model |
|--------|-------|
| **Transaction fees** | 0.5% on integrated payment processing |
| **Premium add-ons** | Advanced analytics, AI recommendations, multi-warehouse |
| **Onboarding services** | White-glove setup, data migration, custom theme design |
| **API access** | Third-party integrations (accounting, shipping providers) |

---

## 6. Competitive Landscape

| Competitor | Weakness Naseej Exploits |
|------------|--------------------------|
| Shopify POS | Delayed inventory sync (5+ min), not clothing-specialized |
| Square | Limited variant matrix, no dedicated deployment model |
| Lightspeed | High cost, complex setup, generic retail focus |
| Custom solutions | Months of development, ongoing maintenance burden |

**Naseej Differentiators:**
- Sub-second inventory sync (vs. minutes)
- Clothing-native variant matrix (Size Groups × Color) with category-aware sizing
- Template-based onboarding — single shared server, per-merchant frontend repos
- Each merchant gets their own custom domain and branded apps

---

## 7. Key Business Metrics (KPIs)

| Metric | Year 1 Target | Year 3 Target |
|--------|---------------|---------------|
| Active merchants | 50 | 500 |
| Monthly Recurring Revenue (MRR) | $15K | $200K |
| Merchant churn rate | < 5% monthly | < 3% monthly |
| Average onboarding time | < 1 hour | < 15 minutes (self-serve) |
| Inventory drift incidents | 0 | 0 |
| System uptime | 99.9% | 99.95% |

---

## 8. Risk Assessment

| Risk | Probability | Impact | Mitigation |
|------|-------------|--------|------------|
| Payment gateway integration delays | Medium | High | Start with cash + manual card; integrate progressively |
| Merchant data migration complexity | High | Medium | Build CSV import tools; offer migration service |
| Hardware compatibility issues (printers/scanners) | Medium | Medium | Support top 3 ESC/POS brands; maintain compatibility matrix |
| Multi-tenant data leak | Low | Critical | Row-level isolation enforced at ORM layer (`Scoped(ctx)`) + integration tests |
| Regional payment regulation changes | Medium | Medium | Abstract payment layer; gateway-agnostic architecture |

---

## 9. Stakeholder Sign-Off Requirements

| Stakeholder | Approval Needed For |
|-------------|---------------------|
| Product Owner | Feature prioritization, MVP scope |
| Engineering Lead | Architecture decisions, technology choices |
| Security Officer | Multi-tenant isolation, PCI compliance approach |
| Business Development | Pricing model, go-to-market strategy |

---

## 10. Success Criteria (MVP Launch)

The minimum viable product is considered launch-ready when:

1. A merchant can complete a full POS transaction (scan → pay → receipt) in under 15 seconds
2. The same sale reflects on the web storefront inventory within 1.5 seconds
3. A new merchant can be onboarded (template clone + deploy) in under 2 hours
4. Zero inventory oversell incidents during a 1000-transaction stress test
5. Three pilot merchants are running in production with positive feedback
