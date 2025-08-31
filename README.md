# Order Service API

API for creating orders for building spacecraft. Orders are managed in-memory with support for creation, payment, retrieval, and cancellation.

## Prerequisites

### Redocly CLI
```bash
npm install -g @redocly/cli
```

### OpenAPI Generator
```bash
# macOS
brew install openapi-generator

# Linux
npm install -g @openapitools/openapi-generator-cli
```

## Code Generation

### Generate Go Server Code

**macOS & Linux:**
```bash
openapi-generator generate -i api/openapi/order/v1/order.openapi.yaml -g go-server -o ./generated/server
```

### Generate Go Client SDK

**macOS & Linux:**
```bash
openapi-generator generate -i api/openapi/order/v1/order.openapi.yaml -g go -o ./generated/client
```

### Generate Documentation

**macOS & Linux:**
```bash
openapi-generator generate -i api/openapi/order/v1/order.openapi.yaml -g html2 -o ./docs
```

### Bundle OpenAPI Spec

```bash
redocly bundle order.openapi.yaml -o order.bundled.yaml
```

## API Endpoints

- `POST /api/v1/orders` - Create a new order
- `GET /api/v1/orders/{order_uuid}` - Get order details
- `POST /api/v1/orders/{order_uuid}/pay` - Pay for an order
- `POST /api/v1/orders/{order_uuid}/cancel` - Cancel an order

## Development

```bash
# Build the project
make build

# Run tests
make test

# Start the server
make run
```

## Payment Methods

- `UNKNOWN` - Unknown method
- `DEBIT_CARD` - Debit card
- `BANK_TRANSFER` - Bank transfer
- `CREDIT_CARD` - Credit card
- `CRYPTO` - Cryptocurrency

## Order Status

- `PENDING_PAYMENT` - Order is waiting for payment
- `PAID` - Order has been successfully paid
- `CANCELLED` - Order has been cancelled
