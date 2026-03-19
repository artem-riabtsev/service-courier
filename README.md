# course-go-avito-artem-riabtsev
# Service Courier

Go microservice for courier management.

[![Go CI Pipeline](https://github.com/yourusername/service-courier/actions/workflows/ci.yml/badge.svg)](https://github.com/yourusername/service-courier/actions/workflows/ci.yml)

## Features

- **Courier Management**: CRUD operations for couriers with status tracking
- **Order Assignment**: Automatic courier assignment to incoming orders
- **Event-Driven Architecture**: Kafka integration for order status updates
- **Metrics & Monitoring**: Prometheus metrics and Grafana dashboards
- **Structured Logging**: JSON-formatted logs for better observability
- **Background Processing**: Automatic handling of overdue deliveries
- **Health Checks**: Readiness and liveness probes

## Architecture

- **HTTP API**: RESTful endpoints for courier and delivery management
- **Database**: PostgreSQL for data persistence
- **Message Queue**: Kafka for order event processing
- **Monitoring**: Prometheus for metrics collection, Grafana for visualization
- **Service Discovery**: Integration with external order service via gRPC/HTTP

## Quick Start

```bash
# 1. Copy environment template (if not exists)
cp .env.example .env

# 2. Install goose globally (one-time installation)
go install github.com/pressly/goose/v3/cmd/goose@latest

# 3. Setup Goose environment variables (one-time setup)
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bashrc
echo 'export GOOSE_DRIVER=postgres' >> ~/.bashrc
echo 'export GOOSE_DBSTRING="postgres://myuser:mypassword@localhost:5432/test_db"' >> ~/.bashrc
echo 'export GOOSE_MIGRATION_DIR=./migrations' >> ~/.bashrc
source ~/.bashrc

# 4. Start PostgreSQL in Docker
docker-compose up -d

# 5. Run database migrations
goose up

# 6. Run server (uses PORT from .env or default 8080)
go run cmd/main.go

# 7. Run with custom port flag (overrides .env)
go run cmd/main.go --port 3000
```

## API Endpoints
### Courier Management
- GET /couriers - List all couriers

- GET /courier/{id} - Get courier by ID

- POST /courier - Create new courier

- PUT /courier - Update courier

### Delivery Management
- POST /delivery/assign - Assign courier to order

- POST /delivery/unassign - Unassign courier from order

### Health & Monitoring
- GET /ping - Service health check

- HEAD /healthcheck - Readiness probe

- GET /metrics - Prometheus metrics endpoint

### Event Processing
- The service consumes Kafka events from the orders topic:

- created - Assigns available courier to new order

- cancelled - Unassigns courier and removes delivery

- completed - Frees courier while keeping delivery record

### Monitoring Stack
#### Prometheus
- Scrapes metrics from /metrics endpoint

- Available at: http://localhost:9090

- Configuration: prometheus.yml

#### Grafana
- Dashboard for service metrics

- Available at: http://localhost:3000 (admin/admin)

- Pre-configured data source: Prometheus