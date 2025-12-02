# Payment Gateway

A unified payment integration service.

## Features

- Multiple payment providers (Adyen, Stripe)
- Idempotency support
- Webhook handling
- Refund management

## Quick Start

### Using Docker Compose

\`\`\`bash
docker-compose up
\`\`\`

### Local Development

\`\`\`bash
# Install dependencies
go mod download

# Run database
docker-compose up postgres -d

# Run application
go run cmd/server/main.go
\`\`\`

## API Endpoints

### Create Payment
\`\`\`bash
POST /api/v1/payments
Headers:
  X-API-Key: your-api-key
  Idempotency-Key: unique-key (optional)
Body:
{
  "order_id": "ORDER123",
  "amount": 100.50,
  "currency": "USD",
  "provider": "mock",
  "customer_email": "customer@example.com"
}
\`\`\`

### Get Payment
\`\`\`bash
GET /api/v1/payments/:id
Headers:
  X-API-Key: your-api-key
\`\`\`

### List Payments
\`\`\`bash
GET /api/v1/payments?limit=20&offset=0
Headers:
  X-API-Key: your-api-key
\`\`\`

### Cancel Payment
\`\`\`bash
POST /api/v1/payments/:id/cancel
Headers:
  X-API-Key: your-api-key
\`\`\`

### Refund Payment
\`\`\`bash
POST /api/v1/payments/:id/refund
Headers:
  X-API-Key: your-api-key
Body:
{
  "amount": 50.00
}
\`\`\`

## Configuration

Edit `config/config.yaml` to configure providers and database.

## License

MIT