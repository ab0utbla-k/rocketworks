## 🚀 Assignment: OrderService Implementation

The order service manages orders for building spacecraft. All orders are stored in an internal in-memory `map` structure. The service interacts with `InventoryService` to calculate order prices and with `PaymentService` to process payments.

---

### 📌 Goals

- Master integration of multiple gRPC/HTTP services
- Manage internal application state (in-memory)
- Implement order lifecycle: creation → payment → cancellation

---

### 📁 Implementation Requirements

#### `POST /api/v1/orders` — Create Order

Creates a new order based on user-selected parts.

**Request Body:**

```json
{
  "user_uuid": "string",
  "part_uuids": [
    "string",
    "string"
  ]
}
```

**Response:**

```json
{
  "order_uuid": "string",
  "total_price": 123.45
}
```

**Behavior:**

- Retrieves parts via `InventoryService.ListParts`
- Validates that all parts exist. If any part is missing — returns error
- Calculates `total_price`
- Generates `order_uuid`
- Saves order with status `PENDING_PAYMENT`

---

#### `POST /api/v1/orders/{order_uuid}/pay` — Pay Order

Processes payment for a previously created order.

**Request Body:**

```json
{
  "payment_method": "DEBIT_CARD"
}
```

**Response:**

```json
{
  "transaction_uuid": "string"
}
```

**Behavior:**

- Finds order by `order_uuid`. If not found — returns 404 Not Found
- Calls `PaymentService.PayOrder`, passes `user_uuid`, `order_uuid` and `payment_method`. Receives `transaction_uuid`
- Updates order: status → `PAID`, saves `transaction_uuid`, `payment_method`

---

#### `GET /api/v1/orders/{order_uuid}` — Get Order by UUID

Returns order information.

**Response:**

```json
{
  "order_uuid": "string",
  "user_uuid": "string",
  "part_uuids": [
    "string"
  ],
  "total_price": 123.45,
  "transaction_uuid": "string",
  "payment_method": "DEBIT_CARD",
  "status": "PAID"
}
```

**Behavior:**

- Searches for order by UUID
- If found — returns order details
- If not found — 404 Not Found

---

#### `POST /api/v1/orders/{order_uuid}/cancel` — Cancel Order

Cancels an order.

**Responses:**

- `204 No Content` — order successfully cancelled
- `404 Not Found` — order not found
- `409 Conflict` — order already paid and cannot be cancelled

**Behavior:**

- Checks order status
- If `PENDING_PAYMENT` — changes status to `CANCELLED`
- If `PAID` — returns 409 error

---

### 💾 In-Memory Order Structure

| Field            | Type     | Description                                     |
|------------------|----------|-------------------------------------------------|
| order_uuid       | string   | Unique order identifier                         |
| user_uuid        | string   | User UUID                                       |
| part_uuids       | []string | List of part UUIDs                              |
| total_price      | double   | Total cost                                      |
| transaction_uuid | *string  | Transaction UUID (if paid)                      |
| payment_method   | *string  | Payment method (if paid)                        |
| status           | string   | Status (`PENDING_PAYMENT`, `PAID`, `CANCELLED`) |

---

### 📦 Enum: PaymentMethod

| Name          | Value | Description    |
|---------------|-------|----------------|
| UNKNOWN       | 0     | Unknown method |
| DEBIT_CARD    | 1     | Debit card     |
| BANK_TRANSFER | 2     | Bank transfer  |
| CREDIT_CARD   | 3     | Credit card    |
| CRYPTO        | 4     | Cryptocurrency |

---

### ✅ Acceptance Criteria

- All 4 methods implemented strictly according to OpenAPI contract
- Order status handling implemented correctly
- Orders properly saved and updated in `map` with race condition protection
- Request validation and error handling present
