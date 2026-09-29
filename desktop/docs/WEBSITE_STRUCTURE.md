# Web Storefront Directory Structure

## Naseej Web — Next.js E-Commerce Storefront

**Stack:** Next.js 14+ (App Router) · TypeScript · Tailwind CSS · shadcn/ui  
**Deployment:** Each merchant gets their own repo + custom domain (no slug/subdomain routing)  
**Rendering:** SSR + ISR for product pages, static for marketing pages  
**State:** React Context + lightweight client state (no heavy state library)  
**API:** Connects to shared central server via `X-Merchant-ID` header

---

## Directory Layout

```
web/
├── src/
│   ├── app/
│   │   ├── layout.tsx                        # Root layout (theme provider, fonts, global providers)
│   │   ├── page.tsx                          # Homepage (hero, collections, featured)
│   │   │
│   │   ├── collections/
│   │   │   ├── page.tsx                      # All collections grid
│   │   │   └── [collection_slug]/
│   │   │       └── page.tsx              # Category product grid (filterable)
│   │   │
│   │   ├── products/
│   │   │   └── [product_slug]/
│   │   │       └── page.tsx              # PDP: variant selector, stock, images
│   │   │
│   │   ├── cart/
│   │   │   └── page.tsx                  # Cart overview, stock validation
│   │   │
│   │   ├── checkout/
│   │   │   ├── page.tsx                  # Checkout entry (redirect to step 1)
│   │   │   ├── information/
│   │   │   │   └── page.tsx              # Shipping & customer details
│   │   │   ├── shipping/
│   │   │   │   └── page.tsx              # Delivery method selection
│   │   │   ├── payment/
│   │   │   │   └── page.tsx              # Payment gateway integration
│   │   │   └── confirmed/
│   │   │       └── page.tsx              # Order success page
│   │   │
│   │   ├── account/
│   │   │   ├── page.tsx                  # Account dashboard redirect
│   │   │   ├── login/
│   │   │   │   └── page.tsx              # Login form
│   │   │   ├── register/
│   │   │   │   └── page.tsx              # Registration form
│   │   │   ├── dashboard/
│   │   │   │   └── page.tsx              # Profile, loyalty points summary
│   │   │   └── orders/
│   │   │       ├── page.tsx              # Order history list
│   │   │       └── [order_id]/
│   │   │           └── page.tsx          # Order detail + tracking
│   │   │
│   │   ├── not-found.tsx                     # 404 page
│   │   └── globals.css                       # Tailwind directives + CSS variables
│   │
│   ├── components/
│   │   ├── ui/                               # Base primitives (shadcn/ui)
│   │   │   ├── button.tsx
│   │   │   ├── input.tsx
│   │   │   ├── select.tsx
│   │   │   ├── dialog.tsx
│   │   │   ├── sheet.tsx
│   │   │   ├── skeleton.tsx
│   │   │   ├── badge.tsx
│   │   │   └── ...
│   │   │
│   │   ├── layout/                           # Structural layout components
│   │   │   ├── header.tsx                    # Navigation bar with cart icon
│   │   │   ├── footer.tsx                    # Site footer
│   │   │   ├── mobile-nav.tsx                # Responsive hamburger menu
│   │   │   └── announcement-bar.tsx          # Promotional banner
│   │   │
│   │   ├── product/                          # Product domain components
│   │   │   ├── product-card.tsx              # Grid card (image, name, price)
│   │   │   ├── product-grid.tsx              # Responsive grid container
│   │   │   ├── variant-selector.tsx          # Size × Color matrix picker
│   │   │   ├── price-display.tsx             # Price with sale/compare formatting
│   │   │   ├── stock-indicator.tsx           # Real-time availability badge
│   │   │   ├── image-gallery.tsx             # Product image carousel/zoom
│   │   │   └── size-guide.tsx                # Size chart modal
│   │   │
│   │   ├── cart/                             # Cart domain components
│   │   │   ├── cart-drawer.tsx               # Slide-out cart panel
│   │   │   ├── cart-item.tsx                 # Individual line item
│   │   │   ├── cart-summary.tsx              # Subtotal, tax, total display
│   │   │   └── quantity-selector.tsx         # +/- quantity input
│   │   │
│   │   ├── checkout/                         # Checkout domain components
│   │   │   ├── checkout-steps.tsx            # Progress indicator
│   │   │   ├── address-form.tsx              # Shipping address form
│   │   │   └── payment-form.tsx              # Payment method selection
│   │   │
│   │   └── shared/                           # Cross-cutting components
│   │       ├── search-bar.tsx                # Product search with autocomplete
│   │       ├── filters.tsx                   # Category/size/color/price filters
│   │       ├── pagination.tsx                # Cursor-based pagination
│   │       └── empty-state.tsx               # Generic empty state display
│   │
│   ├── lib/                                  # Utilities & configuration
│   │   ├── graphql-client.ts                 # Apollo/urql client setup
│   │   ├── api-client.ts                     # Typed GraphQL request wrapper
│   │   ├── utils.ts                          # General utilities (cn, formatPrice, etc.)
│   │   └── constants.ts                      # App-wide constants
│   │
│   ├── config/
│   │   └── merchant.ts                       # Per-deployment merchant configuration
│   │
│   ├── graphql/                              # GraphQL operations
│   │   ├── queries/
│   │   │   ├── products.graphql
│   │   │   ├── collections.graphql
│   │   │   ├── cart.graphql
│   │   │   ├── customer.graphql
│   │   │   └── merchant.graphql              # MerchantBootstrap query
│   │   ├── mutations/
│   │   │   ├── cart.graphql
│   │   │   ├── checkout.graphql
│   │   │   ├── auth.graphql
│   │   │   └── customer.graphql
│   │   └── generated/                        # Auto-generated typed hooks/operations
│   │       └── graphql.ts
│   │
│   ├── hooks/                                # Custom React hooks
│   │   ├── use-cart.ts                       # Cart state & operations
│   │   ├── use-merchant.ts                   # Merchant context consumer
│   │   ├── use-auth.ts                       # Auth state & token management
│   │   ├── use-real-time.ts                  # WebSocket subscription handler
│   │   └── use-filters.ts                    # URL-synced filter state
│   │
│   ├── providers/                            # React context providers
│   │   ├── merchant-provider.tsx             # Merchant config + bootstrap data
│   │   ├── cart-provider.tsx                 # Cart state (persisted to localStorage)
│   │   ├── auth-provider.tsx                 # Auth state + token refresh
│   │   └── realtime-provider.tsx             # WebSocket connection manager
│   │
│   └── types/                                # TypeScript type definitions
│       ├── product.ts
│       ├── order.ts
│       ├── customer.ts
│       └── merchant.ts
│
├── public/
│   ├── fonts/                                # Self-hosted web fonts
│   └── images/                               # Static assets (fallbacks)
│
├── next.config.ts                            # Next.js configuration
├── tailwind.config.ts                        # Tailwind + theme token mapping
├── tsconfig.json
├── package.json
├── codegen.ts                                # GraphQL code generation config
└── .env.local                                # Local environment variables
```

---

## Key Architectural Patterns

### Theme Application

Each merchant's storefront is its own deployed app on its own domain. The theme from `merchant.ts` is applied immediately as CSS variables on the root layout — preventing any flash of unstyled content:

```typescript
// src/app/layout.tsx
import { merchant } from '@/config/merchant';

export default function RootLayout({ children }) {
  return (
    <html>
      <body
        style={{
          '--color-primary': merchant.theme.primaryColor,
          '--color-secondary': merchant.theme.secondaryColor,
          '--font-family': merchant.theme.fontFamily,
        } as React.CSSProperties}
      >
        <MerchantProvider>
          <CartProvider>
            {children}
          </CartProvider>
        </MerchantProvider>
      </body>
    </html>
  );
}
```

### Dynamic SEO Metadata

Every page generates merchant-specific metadata for search engine optimization:

```typescript
// src/app/products/[product_slug]/page.tsx
export async function generateMetadata({ params }): Promise<Metadata> {
  const product = await fetchProduct(params.product_slug);
  return {
    title: `${product.name} | ${merchant.name}`,
    description: product.description,
    openGraph: {
      images: [{ url: product.images[0]?.url }],
    },
  };
}
```

### Merchant Bootstrap Flow

```
App Start
    │
    ├── 1. Apply static theme (merchant.ts → CSS variables)  ← instant
    │
    ├── 2. Call MerchantBootstrap query                       ← async
    │       └── Returns: tax rate, loyalty config, payment gateways, locations
    │
    └── 3. Hydrate MerchantProvider context                   ← ready
```

### Real-Time Stock Updates

The storefront subscribes to inventory changes via WebSocket, updating product availability without page refresh:

```typescript
// src/hooks/use-real-time.ts
const { data } = useSubscription(INVENTORY_UPDATED);

useEffect(() => {
  if (data?.inventoryUpdated) {
    // Update local cache with new stock quantity
    updateVariantStock(data.inventoryUpdated.variantId, data.inventoryUpdated.newAvailable);
  }
}, [data]);
```
