# WMS — Warehouse Management System

A full-stack Warehouse Management System that integrates with marketplace platforms to manage order fulfillment — from order sync, picking, packing, through shipping.

| Layer    | Stack                                      |
| -------- | ------------------------------------------ |
| Backend  | Go, Fiber, Bun ORM, PostgreSQL             |
| Frontend | React, TypeScript, TanStack Query, Zustand |

## Evidence of Working Flow

Screenshots demonstrating the working flow can be found in the `working-flow/` folder.

### How to run :

```bash
Backend
cd backend/ -> make run-dev

Register first using :
you can choose role between :
ADMIN,STAFF,PICKER,PACKER
ADMIN : See list order, detail Order and ship order
STAFF : See list Order, detail Order
PICKER : See list Order, detail Order and Pick Oder
PACKER : See list Order, detail Order and Pack order

curl -X POST http://localhost:1903/auth/register -H "Content-Type: application/json" -d '{"email":"tedyysyyach@gmail.com","password":"1234567890","role":"ADMIN"}'

Migration
cd backend/ -> make migrate-init -> make migarte-up

Frontend
cd frontend/ -> npm run install -> npm run dev

```

### How to build :

```bash
Backend
cd backend/ -> GOOS=linux GOARCH=amd64 go build -o app -> chmod +x app
// you can run it with pm2 or using systemd

Frontend
cd frontend/ -> npm run build
// you can also run it with pm2
```

---

## Table of Contents

- [Architecture](#architecture)
- [Database Design](#database-design)
- [Order Lifecycle](#order-lifecycle)
- [Marketplace Integration](#marketplace-integration)
- [Error Handling](#error-handling)
- [API Endpoints](#api-endpoints)
- [Frontend](#frontend)

---

## Architecture

### Backend

The backend follows a layered architecture with clear separation of concerns:

```
cmd/
├── api/            → entrypoint: setup fiber, routes, DI
└── migrate/        → database migration runner

core/
├── base/           → generic bun transaction & query helpers
├── config/         → app, db, fiber, logger, validator, viper config
├── errors/         → AppError types & constructors
├── lib/marketplace → marketplace HTTP client (auth, orders, logistics)
├── middlewares/    → JWT auth middleware
├── routes/         → route registration (auth, webhook & order)
└── utils/          → response wrapper, token, password, formatter

internal/
├── controller/     → HTTP handlers (parse request, call service, return response)
├── converter/      → entity ↔ model mapping
├── entity/         → bun ORM structs (DB schema)
├── model/          → request/response structs
├── repository/     → DB queries (CRUD, upsert, bulk ops)
└── service/        → business logic

migrations/         → SQL migration files (up/down)
```

#### Backend Request Flow

```
HTTP Request
    │
    ├── /auth/*        → no auth required
    ├── /webhook/*     → no auth required (called by marketplace server)
    │
    └── /order/*
            │
            ▼
      Auth Middleware (JWT validation)
            │
            ▼
      Controller (parse & validate request)
            │
            ▼
      Service (business logic)
            │
            ├──▶ Repository (DB via Bun ORM)
            │
            └──▶ Marketplace Client (external HTTP API)
                        │
                        └──▶ Session Store (token management)
```

### Frontend

```
src/
├── features/
│   ├── auth/
│   │   ├── components/     → LoginPage, ProtectedRoute
│   │   ├── hooks/          → useAuth (login/logout logic)
│   │   ├── store/          → authStore (Zustand + persist)
│   │   └── types/
│   └── dashboard/
│       └── outbound/
│           ├── components/ → OrderPage, ActionButton, StatsCard
│           ├── hooks/      → useOrders, usePickOrder, usePackOrder, useShipOrder
│           └── types/
├── components/             → shared UI: Button, Input, Label, Toast
├── lib/
│   ├── axios.ts            → Axios instance + interceptors
│   └── date.ts             → date formatters
├── router/                 → React Router config
└── types/
    └── response.ts         → ApiResponse<T> generic type
```

#### Frontend Request Flow

```
User Action
    │
    ▼
TanStack Query (useQuery / useMutation)
    │
    ▼
Axios Instance (lib/axios.ts)
    │
    ├── Request interceptor → attach Bearer token from Zustand store
    │
    ▼
Backend API
    │
    ├── Response interceptor → on 401, logout + redirect to /login
    │
    ▼
UI re-render via query invalidation
```

---

## Database Design

### Tables

#### `users`

Stores internal user accounts for system access.

| Column          | Type   | Constraint       |
| --------------- | ------ | ---------------- |
| `id`            | string | PK               |
| `email`         | string | UNIQUE, NOT NULL |
| `password_hash` | string | NOT NULL         |

#### `orders`

Main order table. Synced from marketplace and tracked through the WMS fulfillment pipeline.

| Column                    | Type          | Constraint       | Notes                      |
| ------------------------- | ------------- | ---------------- | -------------------------- |
| `id`                      | string        | PK               | Internal UUID              |
| `order_sn`                | string        | UNIQUE, NOT NULL | Marketplace business key   |
| `shop_id`                 | string        | NOT NULL         |                            |
| `marketplace_status`      | string        | NOT NULL         | Status from marketplace    |
| `shipping_status`         | string        | NOT NULL         | Carrier shipping status    |
| `wms_status`              | string (enum) | NOT NULL         | Internal WMS status        |
| `tracking_number`         | string        | NULLABLE         | Available after shipping   |
| `total_amount`            | decimal(15,2) | NOT NULL         |                            |
| `raw_marketplace_payload` | jsonb         |                  | Raw API response for audit |
| `created_at`              | timestamp     | NOT NULL         |                            |
| `updated_at`              | timestamp     | NOT NULL         |                            |

**`wms_status` enum values:**

| Value           | Description                        |
| --------------- | ---------------------------------- |
| `READY_TO_PICK` | Order received, waiting for picker |
| `PICKING`       | Picker is collecting items         |
| `PACKED`        | Items packed, ready to ship        |
| `SHIPED`        | Handed over to carrier             |

#### `order_items`

Line items for each order. Joined via `order_sn`.

| Column       | Type          | Constraint                     | Notes |
| ------------ | ------------- | ------------------------------ | ----- |
| `id`         | string        | PK                             |       |
| `order_sn`   | string        | FK → orders.order_sn, NOT NULL |       |
| `sku`        | string        | NOT NULL                       |       |
| `quantity`   | int           | NOT NULL                       |       |
| `price`      | decimal(15,2) | NOT NULL                       |       |
| `created_at` | timestamp     | NOT NULL                       |       |

### Relationships

```
users       (standalone — auth only)

orders      one to many{ order_items  (via order_sn) }
```

> **Note:** The join key between `orders` and `order_items` is `order_sn` (the marketplace business key), not the internal `id`. This keeps the data consistent with how the marketplace identifies orders.

---

## Order Lifecycle

Orders move through a strict sequential WMS status pipeline. Each transition is enforced — an order cannot skip a step.

```
Marketplace
    │
    │  SyncOrders (on GET /order/)
    ▼
[READY_TO_PICK]  ──── POST /:order_sn/pick ────▶  [PICKING]
                                                        │
                                           POST /:order_sn/pack
                                                        │
                                                        ▼
                                                    [PACKED]
                                                        │
                                           POST /:order_sn/ship
                                           (validates channel,
                                            calls marketplace,
                                            gets tracking no.)
                                                        │
                                                        ▼
                                                    [SHIPED]
```

### Stage Details

**1. Sync — `READY_TO_PICK`**

Triggered automatically when `GET /order/` is called. The service fetches the current order list from the marketplace, then upserts orders into the database. `wms_status` is mapped from the marketplace's `shipping_status`:

| Marketplace `shipping_status`       | WMS Status      |
| ----------------------------------- | --------------- |
| `awaiting_pickup`, `label_created`  | `READY_TO_PICK` |
| `shipped`, `delivered`, `cancelled` | `SHIPED`        |
| _(default)_                         | `READY_TO_PICK` |

Order items are fully replaced on each sync (delete + bulk insert) to keep them up to date.

**2. Pick — `PICKING`**

`POST /:order_sn/pick`

- Guard: order must be in `READY_TO_PICK`
- Updates `wms_status` to `PICKING`

**3. Pack — `PACKED`**

`POST /:order_sn/pack`

- Guard: order must be in `PICKING`
- Updates `wms_status` to `PACKED`

**4. Ship — `SHIPED`**

`POST /:order_sn/ship`

- Guard: order must be in `PACKED`
- Validates that the requested `channel_id` exists via `LogisticChannel` API
- Calls `ShipOrder` on the marketplace client
- Updates `wms_status` to `SHIPED`, sets `shipping_status` and `tracking_number` from marketplace response

> ⚠️ Tracking numbers are always generated by the Marketplace — the WMS never generates its own.

### Webhook — Async Status Updates

In addition to the sync-on-demand flow above, the marketplace can push status changes at any time via webhooks. These are handled independently from the WMS lifecycle and only update the relevant status fields without touching `wms_status`.

| Webhook Endpoint                | Updates              |
| ------------------------------- | -------------------- |
| `POST /webhook/order-status`    | `marketplace_status` |
| `POST /webhook/shipping-status` | `shipping_status`    |

---

## Marketplace Integration

The marketplace client (`core/lib/marketplace`) is a self-contained HTTP client that handles authentication, token lifecycle, order fetching, and logistics.

### Client Configuration

The client is initialized via `NewClient()` using Viper config keys:

| Config Key        | Description               |
| ----------------- | ------------------------- |
| `mock.url`        | Marketplace base URL      |
| `mock.partnerId`  | Partner ID for signing    |
| `mock.partnerKey` | Partner key (HMAC secret) |

Additional options are injected via the functional options pattern (`MarketplaceClientOption`):

```go
client, _ := marketplace.NewClient(viper,
    marketplace.WithSessionStore(store, "session-key"),
    marketplace.WithShopID("shop-123"),
    marketplace.WithLogger(logger),
    marketplace.WithTimeout(30 * time.Second),
)
```

### Request Signing

All OAuth requests are signed using HMAC-SHA256:

```
base_string = partnerID + path + timestamp + extra
signature   = HMAC-SHA256(partnerKey, base_string)
```

The `extra` field varies per endpoint: `shopID` for authorize, `code` for token exchange, `accessToken` for refresh.

### Token Lifecycle — `EnsureToken`

`EnsureToken` is called automatically before every API request. It handles three cases:

```
EnsureToken()
    │
    ├── No session found
    │       └──▶ Authorize() → GetToken() → save session
    │
    ├── Session expired
    │       ├──▶ RefreshToken()
    │       │       ├── success → update session
    │       │       └── 403 Forbidden → re-Authorize() → GetToken()
    │       └──▶ update session
    │
    └── Token still valid
            └──▶ use existing access token
```

### Session Store

Token sessions are stored in-memory via `MemorySessionStore` (implements `SessionStore` interface). The interface allows swapping to a persistent store (Redis, DB) without changing the client:

```go
type SessionStore interface {
    Get(ctx, key) (*Session, error)
    Set(ctx, key, session) error
    Delete(ctx, key) error
}
```

`MemorySessionStore` is safe for concurrent use via `sync.RWMutex`.

### Retry & Error Handling

`DoRequest` automatically retries up to **3 attempts** with exponential backoff for transient errors:

| Status Code | Behaviour                                         | Backoff         |
| ----------- | ------------------------------------------------- | --------------- |
| `401`       | Call `EnsureToken()` to refresh token, then retry | immediate       |
| `429`       | Rate limited — wait and retry                     | 500ms → 1s → 2s |
| `500`       | Random marketplace failure — retry                | 500ms → 1s → 2s |
| Other 4xx   | Client error — return immediately, no retry       | —               |

Backoff is context-aware: if the request context is cancelled (e.g. client disconnects), the retry loop exits immediately instead of waiting for the delay.

If all 3 attempts fail, the last error is returned wrapped with the attempt count.

### Available API Methods

| Method              | HTTP | Endpoint             | Description                       |
| ------------------- | ---- | -------------------- | --------------------------------- |
| `Authorize()`       | GET  | `/oauth/authorize`   | Initiate OAuth flow               |
| `GetToken()`        | POST | `/oauth/token`       | Exchange auth code for token      |
| `RefreshToken()`    | POST | `/oauth/token`       | Refresh expired access token      |
| `OrderList()`       | GET  | `/order/list`        | Fetch all orders from marketplace |
| `OrderDetail()`     | GET  | `/order/detail`      | Fetch single order by `order_sn`  |
| `LogisticChannel()` | GET  | `/logistic/channels` | Get available shipping channels   |
| `ShipOrder()`       | POST | `/logistic/ship`     | Submit shipment to marketplace    |

---

## Error Handling

All errors are represented as `AppError`, a structured type that carries HTTP status code, a user-facing message, optional details, and an optional internal error (not exposed in response).

### `AppError` Structure

```go
type AppError struct {
    Code       int    // HTTP status code
    Message    string // user-facing message
    Details    any    // optional: validation errors, context map
    Internal   error  // not serialized — for logging only
    StackTrace string // not serialized
}
```

### Error Response Format

```json
{
  "status": false,
  "message": "Validation failed",
  "errors": [
    {
      "field": "channelId",
      "message": "channelId is required",
      "tag": "required",
      "value": ""
    }
  ],
  "trace_id": "optional-trace-id"
}
```

### Error Constructors

| Constructor                            | HTTP Status | Use Case                                  |
| -------------------------------------- | ----------- | ----------------------------------------- |
| `NewBadRequestError(message, details)` | 400         | Invalid input, forbidden state transition |
| `NewValidationError(err)`              | 400         | `validator` struct tag failures           |
| `NewNotFoundError(resource)`           | 404         | Entity not found in DB                    |
| `NewUnauthorizedError(message)`        | 401         | Invalid or missing JWT token              |
| `NewConflictError(message)`            | 409         | Duplicate resource                        |
| `NewInternalError(err)`                | 500         | Unexpected errors (DB, marshaling)        |

### Frontend Error Handling

The Axios instance handles two error cases automatically via response interceptor:

| Scenario              | Behaviour                                      |
| --------------------- | ---------------------------------------------- |
| `401` on any request  | Clear Zustand auth state, redirect to `/login` |
| `401` on login itself | Rejected normally — displayed as toast error   |

Mutation errors (pick, pack, ship) are surfaced to the user via `sonner` toast notifications in `ActionButton`.

### Validation Error Detail

When `NewValidationError` is triggered, each field failure is parsed into a `ValidationError`:

```go
type ValidationError struct {
    Field   string // snake_case field name
    Message string // human-readable message
    Tag     string // validator tag (e.g. "required", "email")
    Value   string // submitted value
}
```

Supported validation tags with human-readable messages: `required`, `email`, `min`, `max`, `len`, `gte`, `lte`, `gt`, `lt`, `uuid`, `url`, `oneof`, `alphanum`, `numeric`.

---

## API Endpoints

### Auth

| Method | Endpoint         | Auth Required | Description       |
| ------ | ---------------- | ------------- | ----------------- |
| POST   | `/auth/register` | No            | Register new user |
| POST   | `/auth/login`    | No            | Login, get JWT    |

**Register / Login request:**

```json
{ "email": "user@example.com", "password": "secret" }
```

**Login response:**

```json
{
  "status": true,
  "message": "Login success",
  "data": { "accessToken": "<jwt>" }
}
```

### Webhooks

Webhook endpoints are **public** — no JWT required. They are called by the marketplace server to push status updates asynchronously.

| Method | Endpoint                   | Description                          |
| ------ | -------------------------- | ------------------------------------ |
| POST   | `/webhook/order-status`    | Update `marketplace_status` of order |
| POST   | `/webhook/shipping-status` | Update `shipping_status` of order    |

**Webhook order-status payload:**

```json
{
  "data": {
    "order_sn": "SHP001",
    "status": "delivered"
  }
}
```

**Webhook shipping-status payload:**

```json
{
  "data": {
    "order_sn": "SHP001",
    "shipping_state": "shipped"
  }
}
```

### Orders

All order endpoints require a valid `Authorization: Bearer <token>` header.

| Method | Endpoint                | Description                                       |
| ------ | ----------------------- | ------------------------------------------------- |
| GET    | `/order/`               | Sync from marketplace, return orders (filterable) |
| GET    | `/order/:order_sn`      | Get order detail with items                       |
| POST   | `/order/:order_sn/pick` | Transition: `READY_TO_PICK` → `PICKING`           |
| POST   | `/order/:order_sn/pack` | Transition: `PICKING` → `PACKED`                  |
| POST   | `/order/:order_sn/ship` | Transition: `PACKED` → `SHIPED`, get tracking no. |

**Query Parameters — `GET /order/`:**

| Param       | Type   | Example                       | Description                        |
| ----------- | ------ | ----------------------------- | ---------------------------------- |
| `wmsStatus` | string | `PICKING` or `PICKING,PACKED` | Filter by one or more WMS statuses |

**Ship order request:**
currently channelId hardcoded in FE using JNE

```json
{ "channelId": "channel-code" }
```

**Ship order response:**

```json
{
  "status": true,
  "message": "Success",
  "data": {
    "order_sn": "SN-001",
    "wms_status": "SHIPED",
    "shipping_status": "shipped",
    "tracking_number": "JNE-123456"
  }
}
```

### Standard Response Envelope

All endpoints return a consistent wrapper:

```json
{
  "status": true,
  "message": "Success",
  "data": {}
}
```

---

## Frontend

### Pages & Routes

| Path                  | Component   | Auth Required | Description             |
| --------------------- | ----------- | ------------- | ----------------------- |
| `/login`              | `LoginPage` | No            | Login form              |
| `/dashboard/outbound` | `OrderPage` | Yes           | Order list & management |

Protected routes are wrapped by `ProtectedRoute` which checks `isAuthenticated` from Zustand store. Unauthenticated users are redirected to `/login`.

### Auth Flow

1. User submits email + password on `LoginPage`
2. `useAuth.login()` calls `POST /auth/login`
3. On success, `accessToken` is saved to Zustand store via `setAuth()`
4. Zustand `persist` middleware saves the token to `localStorage` under key `auth-storage`
5. User is redirected to `/dashboard/outbound`
6. On subsequent requests, the Axios request interceptor reads the token from Zustand and attaches it as `Authorization: Bearer <token>`
7. If any request returns `401` (token expired), the response interceptor automatically calls `logout()` and redirects to `/login`

### State Management

| Concern        | Library           | Notes                                            |
| -------------- | ----------------- | ------------------------------------------------ |
| Auth token     | Zustand + persist | Persisted to localStorage, survives page refresh |
| Server data    | TanStack Query    | Cached, auto-invalidated after mutations         |
| Table UI state | Local useState    | Sorting, filtering, pagination — component-local |

### Order Management — `OrderPage`

The order list is powered by **TanStack Table** with the following features:

- Global search across all columns
- Column sorting (updated_at, wms_status, etc.)
- Client-side pagination (10 / 20 / 50 rows per page) with smart page number display
- Color-coded status badges for `marketplace_status`, `shipping_status`, and `wms_status`
- Stats cards showing total orders and cancelled count
- Per-row action buttons (pick / pack / ship) rendered by `ActionButton`

### Data Fetching — TanStack Query

| Hook             | Type     | Endpoint               | Invalidates            |
| ---------------- | -------- | ---------------------- | ---------------------- |
| `useOrders`      | query    | `GET /order/`          | —                      |
| `useOrderDetail` | query    | `GET /order/:order_sn` | —                      |
| `usePickOrder`   | mutation | `POST /order/:sn/pick` | `orders`, `orders/:sn` |
| `usePackOrder`   | mutation | `POST /order/:sn/pack` | `orders`, `orders/:sn` |
| `useShipOrder`   | mutation | `POST /order/:sn/ship` | `orders`, `orders/:sn` |

After every mutation, both the order list and the relevant order detail are invalidated, triggering a background refetch to keep the UI consistent with the backend state.
