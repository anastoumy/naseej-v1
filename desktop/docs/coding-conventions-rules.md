# Coding Conventions & Rules

> This document defines the standard architecture and rules for all API interactions in this project.
> Every developer and every AI agent (GitHub Copilot) must follow this strictly.

---

## Core Architecture: 3-Layer Data Flow

```
API Client  →  Controllers  →  Hooks  →  Components
(raw HTTP)     (business logic)  (React state)   (UI)
```

**The rule is simple: each layer only talks to the layer directly above it.**
A component never calls a Controller directly. A Controller never uses a Hook. A Hook never touches `fetch`.

---

## Layer 1 — API Client (`src/lib/api/`)

**Responsibility:** Raw HTTP communication only. No business logic. No TypeScript domain types. No auth decisions beyond attaching headers.

### Files

```
src/lib/api/
  api-client.ts        ← Base classes + singleton exports
```

### What lives here

- `BaseApiClient` — abstract class with `buildURL`, `handleResponse`, `request`
- `PublicApiClient` — unauthenticated JSON requests
- `AuthApiClient` — authenticated JSON requests (reads token from Zustand store)
- `MultipartApiClient` — authenticated `multipart/form-data` requests (file uploads)
- Token refresh logic (`ensureValidToken`) — shared utility, lives here only
- Singleton exports: `publicApi`, `authApi`, `multipartApi`
- Utility guards: `isApiSuccess`, `isUnauthorizedError`, `isNotFoundError`, etc.

### Rules

- ✅ Accepts `endpoint: string`, `data: any`, `config?: RequestConfig`
- ✅ Returns `Promise<ApiResponse<T>>`
- ❌ Never imports domain types (e.g. `Product`, `User`)
- ❌ Never contains `if (product.status === ...)` — no business logic
- ❌ Never called directly from a component or hook

### Example

```ts
// ✅ Correct — thin HTTP wrapper
export const authApi = new AuthApiClient();
```

---

## Layer 2 — Controllers (`src/lib/controllers/`)

**Responsibility:** Domain-specific API calls. Maps HTTP responses to typed domain objects. Owns validation (e.g. "productId is required"). Groups all API calls for one domain in one class.

### Files

```
src/lib/controllers/
  product.controller.ts
  order.controller.ts
  auth.controller.ts
  user.controller.ts
  ...
```

### Naming Convention

| Pattern    | Example                         |
| ---------- | ------------------------------- |
| File name  | `product.controller.ts`         |
| Class name | `ProductController`             |
| Methods    | `static async getProducts(...)` |

### What lives here

- Static class methods — one method per API endpoint
- Auth verification (`verifyAuth()` private guard)
- Response unwrapping — throws on `!response.success`
- Domain type mapping — returns typed objects, not raw `ApiResponse`
- Parameter validation — throws `Error` for missing required IDs

### Rules

- ✅ Imports and uses `authApi`, `publicApi`, `multipartApi` from `api-client.ts`
- ✅ Imports domain types (e.g. `Product`, `ProductFilterParams`)
- ✅ Every method is `static async`
- ✅ Throws a descriptive `Error` on failure — never returns `null` silently
- ❌ Never imports `useQuery`, `useMutation`, or any React hook
- ❌ Never uses `useState`, `useEffect`, or any React API
- ❌ Never imports from another Controller

### Example

```ts
// src/lib/controllers/product.controller.ts

export class ProductController {
  private static verifyAuth(): void {
    const { token } = useAuthStore.getState();
    if (!token) throw new Error("Authentication required.");
  }

  static async getProducts(
    filters: ProductFilterParams = {},
  ): Promise<ProductListResponse> {
    this.verifyAuth();
    const response = await authApi.get<ProductListResponse>("/products", {
      params: filters,
    });
    if (!response.success || !response.data) {
      throw new Error(response.error?.message || "Failed to fetch products");
    }
    return response.data;
  }

  static async createProduct(data: CreateProductRequest): Promise<Product> {
    this.verifyAuth();
    const response = await authApi.post<Product>("/products", data);
    if (!response.success || !response.data) {
      throw new Error(response.error?.message || "Failed to create product");
    }
    return response.data;
  }

  static async deleteProduct(productId: number): Promise<void> {
    if (!productId) throw new Error("Product ID is required");
    this.verifyAuth();
    await authApi.delete(`/products/${productId}`);
  }
}
```

---

## Layer 3 — Hooks (`src/hooks/` or `src/features/[domain]/hooks/`)

**Responsibility:** React state management for server data. Wraps Controller methods with `react-query`. Handles loading, error, caching, optimistic updates, and toast notifications.

### Files

```
src/features/shop/hooks/
  useProducts.ts
  useOrders.ts
  useAuth.ts
  ...
```

### Naming Convention

| Pattern       | Example                               |
| ------------- | ------------------------------------- |
| File name     | `useProducts.ts`                      |
| Query hook    | `useProductList(filters, pagination)` |
| Single item   | `useProduct(productId)`               |
| Mutation hook | `useCreateProduct()`                  |
| Mutation hook | `useDeleteProduct()`                  |

### What lives here

- `useQuery` wrappers — fetching, caching, pagination
- `useMutation` wrappers — create, update, delete with toast feedback
- Optimistic updates via `queryClient.setQueryData`
- Cache invalidation via `queryClient.invalidateQueries`
- Toast notifications (`useToast`)
- Query key definitions (centralized per hook file)

### Rules

- ✅ Calls only Controller static methods — never `authApi` directly
- ✅ Every query has a `queryKey` that includes all filter/pagination params
- ✅ Every mutation invalidates or updates the relevant query cache on `onSettled`
- ✅ Mutations use `onError` to rollback optimistic updates
- ❌ Never contains raw `fetch` or any HTTP logic
- ❌ Never throws — handles errors via `onError` callback
- ❌ Never imports from `api-client.ts` directly

### Example

```ts
// src/features/shop/hooks/useProducts.ts

export const useProductList = (
  filters: Partial<ProductFilterParams> = {},
  pagination: PaginationParams = { page: 1, perPage: 10 },
) => {
  return useQuery({
    queryKey: ["products", filters, pagination],
    queryFn: () =>
      ProductController.getProducts({
        ...filters,
        page: pagination.page,
        per_page: pagination.perPage,
      }),
    placeholderData: {
      products: [],
      pagination: { page: 1, perPage: 10, total: 0, totalPages: 0 },
    },
    staleTime: 5 * 60 * 1000,
  });
};

export const useDeleteProduct = () => {
  const queryClient = useQueryClient();
  const { toast } = useToast();

  return useMutation({
    mutationFn: (productId: number) =>
      ProductController.deleteProduct(productId),
    onSuccess: (_, productId) => {
      queryClient.setQueryData(
        ["products"],
        (old: ProductListResponse | undefined) =>
          old
            ? {
                ...old,
                products: old.products.filter((p) => p.id !== productId),
              }
            : old,
      );
      toast({ title: "Product Deleted" });
    },
    onError: (error) => {
      toast({
        title: "Delete Failed",
        description:
          error instanceof Error ? error.message : "Unexpected error",
        variant: "destructive",
      });
    },
    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: ["products"] });
    },
  });
};
```

---

## Layer 4 — Components (`src/components/` or `src/features/[domain]/components/`)

**Responsibility:** UI rendering only. Consumes hooks, renders data, triggers mutations.

### Rules

- ✅ Imports and calls hooks only
- ✅ Handles loading and error states from hook return values
- ❌ Never imports a Controller directly
- ❌ Never imports `authApi`, `publicApi`, or `multipartApi`
- ❌ Never calls `fetch` or any HTTP method

### Example

```tsx
// ✅ Correct
const { data, isLoading } = useProductList(filters);
const { mutate: deleteProduct } = useDeleteProduct();

// ❌ Wrong — never do this in a component
const products = await ProductController.getProducts();
const products = await authApi.get("/products");
```

---

## File Upload Pattern

For endpoints that require `multipart/form-data`:

```ts
// Controller — builds FormData and delegates to multipartApi
static async createProductWithMedia(request: CreateProductWithMediaRequest): Promise<Product> {
  this.verifyAuth();
  const formData = new FormData();
  formData.append("productData", JSON.stringify(request.productData));
  request.imageFiles?.forEach((file) => formData.append("images", file, file.name));

  const response = await multipartApi.post<Product>("/products/with-files", formData);
  if (!response.success || !response.data) throw new Error("Failed to create product");
  return response.data;
}
```

---

## Query Key Convention

All query keys follow this structure:

```ts
["domain"][("domain", filters)][("domain", "detail", id)][ // all records // filtered list // single record
  ("domain", "summary", id, period)
]; // sub-resource
```

**Examples:**

```ts
["products"][("products", filters, pagination)][
  ("products", "detail", productId)
][("products", "sales", productId, "30d")]["orders"][
  ("orders", "detail", orderId)
];
```

---

## Error Handling Convention

| Layer      | How errors are handled                                            |
| ---------- | ----------------------------------------------------------------- |
| API Client | Throws `ApiClientError` with `code`, `status`, `message`          |
| Controller | Unwraps response, re-throws with readable `Error` message         |
| Hook       | Catches via `onError`, shows toast, rolls back optimistic updates |
| Component  | Reads `isError` / `error` from hook, renders error UI             |

---

## Summary Cheat Sheet

```
DO ✅                                        DON'T ❌
─────────────────────────────────────────────────────────────────────
Component calls useProductList()             Component calls ProductController
Hook calls ProductController.getProducts()   Hook calls authApi.get()
Controller calls authApi.get()               Controller imports useQuery
API Client returns ApiResponse<T>            API Client imports Product type
Controller throws on !response.success       Controller returns null silently
Hook invalidates cache on onSettled          Hook skips cache invalidation
```
