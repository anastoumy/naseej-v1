# Developer Guide & Code Standards

## Naseej — Setup, Conventions, and Development Workflow

**Applies to:** All contributors across backend, web, and POS codebases

---

## 1. Engineering Principles

### 1.1 Architecture Rules

| Principle | Description |
|-----------|-------------|
| **Contract-First** | GraphQL schemas are the source of truth. Backend and frontend develop against the same `.graphqls` contracts. |
| **Three-Layer Separation** | Backend follows Resolver → Service → Repository. No layer may skip another. |
| **Composition Over Inheritance** | Build isolated, atomic components. Pass styles and behavior via props, not class hierarchies. |
| **Strict Boundary Enforcement** | Frontend applications (Next.js, Electron) never access the database directly. All data flows through the GraphQL API. |
| **Tenant Isolation by Default** | Every repository query uses `Scoped(ctx)`. There is no "unscoped" query path for tenant data. |

### 1.2 Code Conventions

| Area | Convention |
|------|-----------|
| Go naming | `camelCase` for unexported, `PascalCase` for exported, acronyms uppercase (`ID`, `URL`, `HTTP`) |
| Go errors | Return `error` as last return value. Wrap with context: `fmt.Errorf("create order: %w", err)` |
| TypeScript | Strict mode enabled. No `any` types in committed code. |
| React components | Functional components only. Hooks for state and side effects. |
| CSS | Tailwind utility classes. No custom CSS files except for CSS variable declarations. |
| GraphQL | `camelCase` for fields, `PascalCase` for types, `SCREAMING_SNAKE` for enums |
| File naming | `snake_case` for Go, `kebab-case` for TypeScript/React components |
| Commits | Conventional Commits format: `feat:`, `fix:`, `refactor:`, `docs:`, `chore:` |

---

## 2. Prerequisites

| Tool | Version | Purpose |
|------|---------|---------|
| Go | 1.22+ | Backend server |
| PostgreSQL | 15+ | Primary database |
| Node.js | 20 LTS | Frontend tooling |
| Docker & Docker Compose | Latest | Local development environment |
| gqlgen | v0.17+ | GraphQL code generation |

---

## 3. Initial Setup

### 3.1 Clone & Environment

```bash
git clone git@github.com:your-org/naseej.git
cd naseej

# Copy environment template
cp server/.env.example server/.env
```

**Required environment variables (`server/.env`):**

```env
# Database
DATABASE_URL=postgres://naseej:naseej@localhost:5432/naseej?sslmode=disable

# Server
PORT=8080
GIN_MODE=debug

# Auth
JWT_SECRET=your-dev-secret-key-change-in-production
JWT_EXPIRY=24h
REFRESH_TOKEN_EXPIRY=168h

# CORS
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173
```

### 3.2 Start Infrastructure

```bash
# Launch PostgreSQL via Docker
cd server
docker-compose up -d

# Verify database is accepting connections
docker-compose exec postgres pg_isready
```

### 3.3 Run the Server (AutoMigrate handles schema)

The server uses GORM AutoMigrate — on first startup it creates all tables automatically. No manual migration step needed.

```bash
# Start the server (creates tables on first run)
go run cmd/server/main.go
```

### 3.4 Seed Development Data

```bash
go run cmd/seed/main.go
```

This creates a test merchant, size groups, sample products with variants, and staff accounts for local development.

---

## 4. Backend Development (Go)

### 4.1 Generate GraphQL Code

After modifying any `.graphqls` schema file:

```bash
cd server
go run github.com/99designs/gqlgen generate
```

This regenerates:
- `internal/graph/generated/generated.go` — runtime execution engine
- `internal/graph/model/models_gen.go` — Go structs from GraphQL types
- Resolver stubs for any new queries/mutations

### 4.2 Run the Server

```bash
# Development (with hot reload via air)
air

# Or without hot reload
go run cmd/server/main.go
```

The server starts at `http://localhost:8080`:
- GraphQL Playground: `http://localhost:8080/graphql`
- WebSocket endpoint: `ws://localhost:8080/ws`

### 4.3 Adding a New Domain Feature

Follow this sequence when adding a new feature (e.g., "Gift Cards"):

```
1. Define schema:      internal/graph/schema/gift_card.graphqls
2. Run codegen:        go run github.com/99designs/gqlgen generate
3. Define model:       internal/models/gift_card.go (add to AutoMigrate list in main.go)
4. Add repository:     internal/repository/gift_card_repo.go
5. Add service:        internal/service/gift_card_service.go
6. Implement resolver: internal/graph/gift_card.resolvers.go
7. Write tests:        internal/service/gift_card_service_test.go
8. Restart server:     AutoMigrate creates the new table automatically
```

### 4.4 Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific package
go test ./internal/service/...
```

---

## 5. Web Storefront Development (Next.js)

### 5.1 Setup & Run

```bash
cd web
npm install
npm run dev
```

Development server at `http://localhost:3000`.

### 5.2 GraphQL Client

The web app uses a typed GraphQL client. After schema changes:

```bash
npm run codegen
```

This generates TypeScript types from the backend schemas.

### 5.3 Merchant Configuration

Each deployed storefront is a **separate repo** for each merchant, deployed on their **own custom domain** (e.g., `www.brandname.com`). There is no subdomain routing or slug-based multi-tenancy in the URL.

The storefront connects to the shared central API and identifies itself via the `merchant_id` in `merchant.ts`:

```typescript
// web/src/config/merchant.ts
export const merchant = {
  id: 'uuid-from-database',
  name: 'Brand Name',
  theme: {
    primaryColor: '#1B2A4A',
    secondaryColor: '#D4A574',
    fontFamily: 'Inter',
    logoUrl: '/logo.svg',
  },
  api: {
    url: 'https://api.naseej.com/graphql',
    wsUrl: 'wss://api.naseej.com/ws',
  },
} as const;
```

On boot, the app:
1. Applies static theme immediately as CSS variables (no FOUC)
2. Calls `merchantBootstrap` query to load dynamic settings (tax, loyalty, gateways)

---

## 6. POS Desktop Development (Electron)

### 6.1 Setup & Run

```bash
cd desktop
npm install
npm run dev
```

This launches the Electron app in development mode with Vite HMR.

### 6.2 POS Merchant Configuration

Identical structure to web `merchant.ts`:

```typescript
// desktop/src/config/merchant.ts
export const merchant = {
  id: 'uuid-from-database',
  name: 'Brand Name',
  theme: {
    primaryColor: '#1B2A4A',
    secondaryColor: '#D4A574',
    fontFamily: 'Inter',
    logoUrl: '/logo.svg',
  },
  api: {
    url: 'https://api.naseej.com/graphql',
    wsUrl: 'wss://api.naseej.com/ws',
  },
  pos: {
    defaultLocationId: 'uuid-of-store-location',
    printerModel: 'epson-tm-t88v',
    receiptWidth: 80,  // mm
  },
} as const;
```

### 6.3 Hardware Testing

For local development without physical hardware:

```bash
# Mock printer output to console
MOCK_HARDWARE=true npm run dev
```

---

## 7. New Merchant Onboarding

The shared server is always running. Onboarding a new merchant requires **no server changes**:

```bash
# Step 1: Create merchant in the shared database
psql $DATABASE_URL -c "
  INSERT INTO merchants (business_name, slug, domain, tax_rate, currency)
  VALUES ('New Brand', 'new-brand', 'www.newbrand.com', 0.15, 'SAR')
  RETURNING id;
"

# Step 2: Seed initial data for the merchant
go run cmd/seed/main.go --merchant-id=<returned-uuid>
# Creates: admin user, default roles, permissions, locations, size groups

# Step 3: Create a new website repo for this merchant
# (clone from the web template, configure merchant.ts with the merchant ID)

# Step 4: Create a new POS repo for this merchant
# (clone from the desktop template, configure merchant.ts with the merchant ID)

# Step 5: Deploy website on merchant's custom domain
# Step 6: Package and deliver POS installer
```

**Key points:**
- The API server at `api.naseej.com` is already running and serving all merchants
- No subdomain routing — each merchant's website is its own standalone Next.js app on its own domain
- The frontend sends `X-Merchant-ID` header with every API request
- The server uses `Scoped(ctx)` to isolate data per merchant

---

## 8. Project Scripts Reference

### Backend (`server/Makefile`)

| Command | Description |
|---------|-------------|
| `make run` | Start the Go server (AutoMigrate runs on startup) |
| `make generate` | Run gqlgen code generation |
| `make seed` | Seed development data |
| `make test` | Run all tests |
| `make lint` | Run golangci-lint |
| `make build` | Build production binary |

### Web (`web/package.json`)

| Command | Description |
|---------|-------------|
| `npm run dev` | Start Next.js dev server |
| `npm run build` | Production build |
| `npm run start` | Start production server |
| `npm run codegen` | Generate GraphQL types |
| `npm run lint` | Run ESLint |

### POS (`desktop/package.json`)

| Command | Description |
|---------|-------------|
| `npm run dev` | Start Electron + Vite dev mode |
| `npm run build` | Build for production |
| `npm run package` | Package as distributable (.dmg / .exe / .AppImage) |
| `npm run codegen` | Generate GraphQL types |

---

## 9. Deployment Architecture

```
                    ┌─────────────────────────┐
                    │   api.naseej.com         │
                    │   (Single Go Server)     │
                    │   Serves ALL merchants   │
                    └──────────┬──────────────┘
                               │
               ┌───────────────┼───────────────┐
               │               │               │
    ┌──────────▼──┐  ┌────────▼────┐  ┌──────▼───────┐
    │ brand-a.com │  │ brand-b.com │  │ brand-c.com  │
    │ (Next.js)   │  │ (Next.js)   │  │ (Next.js)    │
    └─────────────┘  └─────────────┘  └──────────────┘
               │               │               │
    ┌──────────▼──┐  ┌────────▼────┐  ┌──────▼───────┐
    │ POS App A   │  │ POS App B   │  │ POS App C    │
    │ (Electron)  │  │ (Electron)  │  │ (Electron)   │
    └─────────────┘  └─────────────┘  └──────────────┘
```

Each frontend is a cloned template with only `merchant.ts` changed. All communicate with the same central API server.
