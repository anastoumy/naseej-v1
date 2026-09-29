# Web Storefront Site Map

## Naseej Web — Route Architecture & Page Definitions

**Framework:** Next.js App Router  
**Deployment:** Each merchant has their own repo deployed on their own custom domain  
**Rendering Strategy:** SSR for dynamic pages, ISR for product catalog, CSR for cart/checkout

---

## Route Tree

Each storefront is a standalone Next.js app on the merchant's own domain. Routes are flat — no slug or subdomain prefix needed.

```
/                                           (Root Layout — Theme & Context Provider)
│
├── /                                       Homepage
│
├── /collections                            Collections Hub
│   └── /[collection_slug]                  Category Product Grid
│
├── /products
│   └── /[product_slug]                     Product Detail Page
│
├── /cart                                   Shopping Cart
│
├── /checkout                               Checkout Flow
│   ├── /information                        Customer & Shipping Details
│   ├── /shipping                           Delivery Method Selection
│   ├── /payment                            Payment Processing
│   └── /confirmed                          Order Confirmation
│
├── /account                                Customer Portal
│   ├── /login                              Authentication
│   ├── /register                           Registration
│   ├── /dashboard                          Profile & Loyalty Summary
│   └── /orders                             Order History
│       └── /[order_id]                     Order Detail & Tracking
│
└── /search                                 Search Results Page
```

---

## Page Specifications

### Homepage (`/`)

| Aspect | Detail |
| -------- | -------- |
| **Purpose** | Brand landing page, discovery, and conversion entry |
| **Rendering** | ISR (revalidate every 60s) |
| **Data** | Featured products, curated collections, hero banners |
| **Components** | Hero carousel, collection cards, product grid (featured), announcement bar |

### Collections Hub (`/collections`)

| Aspect | Detail |
| -------- | -------- |
| **Purpose** | Browse all product categories |
| **Rendering** | ISR (revalidate every 5 min) |
| **Data** | All active categories with product counts and images |
| **Components** | Category grid cards with images |

### Category Grid (`/collections/[collection_slug]`)

| Aspect | Detail |
| -------- | -------- |
| **Purpose** | Filterable product listing within a category |
| **Rendering** | SSR with URL-based filter state |
| **Data** | Products in category, available sizes/colors for filters |
| **Filters** | Size, Color, Price range, Availability, Brand, Sort order |
| **Components** | Product grid, filter sidebar, active filter badges, pagination |

### Product Detail Page (`/products/[product_slug]`)

| Aspect | Detail |
| -------- | -------- |
| **Purpose** | Full product information, variant selection, add-to-cart |
| **Rendering** | ISR (revalidate on inventory change via webhook) |
| **Data** | Product, all variants with stock levels, images, related products |
| **Components** | Image gallery, variant matrix selector (Size x Color), stock indicator, price display, size guide, add-to-cart button |
| **Real-time** | WebSocket subscription updates stock badges without refresh |
| **SEO** | Dynamic `generateMetadata` with product name, description, OG image |

### Shopping Cart (`/cart`)

| Aspect | Detail |
| -------- | -------- |
| **Purpose** | Review selected items, validate stock, proceed to checkout |
| **Rendering** | Client-side (cart state in localStorage + context) |
| **Data** | Cart items with real-time stock validation against API |
| **Components** | Cart item list, quantity selectors, stock warnings, summary totals, coupon input |
| **Validation** | On load: verify all items still in stock at required quantities |

### Checkout — Information (`/checkout/information`)

| Aspect | Detail |
| -------- | -------- |
| **Purpose** | Collect shipping address and customer contact details |
| **Rendering** | CSR (protected — requires items in cart) |
| **Components** | Address form, saved addresses (for logged-in users), contact fields |

### Checkout — Shipping (`/checkout/shipping`)

| Aspect | Detail |
| -------- | -------- |
| **Purpose** | Select delivery method and calculate shipping cost |
| **Rendering** | CSR |
| **Components** | Shipping method options with pricing, estimated delivery dates |

### Checkout — Payment (`/checkout/payment`)

| Aspect | Detail |
| -------- | -------- |
| **Purpose** | Process payment through configured gateway |
| **Rendering** | CSR |
| **Integrations** | Stripe, local payment gateways (configured per merchant) |
| **Components** | Payment method selector, card form, loyalty points redemption option |
| **Security** | PCI-compliant — card data never touches Naseej servers |

### Order Confirmed (`/checkout/confirmed`)

| Aspect | Detail |
| -------- | -------- |
| **Purpose** | Success confirmation, order summary, next steps |
| **Rendering** | CSR (receives order ID from payment step) |
| **Actions** | Clears cart, displays order number, shows estimated delivery |
| **Components** | Success state, order summary, "continue shopping" CTA |

### Customer Login (`/account/login`)

| Aspect | Detail |
|--------|--------|
| **Purpose** | Authenticate returning customers |
| **Components** | Email/password form, "forgot password" link, register CTA |

### Customer Dashboard (`/account/dashboard`)

| Aspect | Detail |
| -------- | -------- |
| **Purpose** | Profile overview, loyalty points balance, quick actions |
| **Rendering** | CSR (protected — requires authentication) |
| **Data** | Customer profile, loyalty balance, recent orders |
| **Components** | Profile card, loyalty points display, recent orders list, edit profile form |

### Order History (`/account/orders`)

| Aspect | Detail |
| -------- | -------- |
| **Purpose** | View all past orders with status |
| **Data** | Paginated order list with status badges |
| **Components** | Order list with status, date, total, and channel indicators |

### Order Detail (`/account/orders/[order_id]`)

| Aspect | Detail |
| -------- | -------- |
| **Purpose** | Full order details, tracking, and digital receipt |
| **Data** | Order items, payments, shipping status, tracking info |
| **Components** | Item list with images, payment breakdown, status timeline, receipt download |

### Search Results (`/search`)

| Aspect | Detail |
| -------- | -------- |
| **Purpose** | Full-text product search with filtering |
| **Rendering** | CSR with debounced API calls |
| **Components** | Search input, results grid, filters, "no results" state |

---

## SEO & Performance Strategy

| Technique | Application |
| ----------- | ------------ |
| **Dynamic Metadata** | Every page generates title, description, and OG tags from product/category data |
| **Structured Data** | Product pages include JSON-LD (Product schema with offers and availability) |
| **Sitemap Generation** | Dynamic sitemap.xml per tenant (products, collections, static pages) |
| **Image Optimization** | Next.js `<Image>` with CDN-backed responsive srcsets |
| **ISR Revalidation** | Product pages revalidate on inventory webhook; collections every 5 min |
| **Core Web Vitals** | LCP < 2.5s, FID < 100ms, CLS < 0.1 (targets) |
