# Campaign Targeting Engine

A high-performance, read-optimized microservice for campaign delivery, built with Go, Go kit, and PostgreSQL.

---

## Features

- **Go kit** architecture for modular, testable microservices.
- **PostgreSQL** as the persistent data store for campaigns and targeting rules.
- **In-memory caching** of campaigns and rules for ultra-low latency delivery.
- **Hot reload**: Automatically refreshes in-memory data from the database at regular intervals (default: every 1 minute).
- **Prometheus metrics** endpoint for monitoring.
- **Docker Compose** setup for Prometheus and Grafana (optional).

---

## Prerequisites

- Go 1.18+
- PostgreSQL database
- Docker (for Prometheus/Grafana monitoring, optional)

---

## Setup

### 1. Clone the repository

```sh
git clone https://github.com/glitchdawg/campaign-targeting-engine.git
cd campaign-targeting-engine
```

### 2. Configure Environment Variables

Create a `.env` file in the project root:

```
DB_USER=youruser
DB_PASSWORD=yourpassword
DB_HOST=localhost
DB_PORT=5432
DB_NAME=yourdb
```

### 3. Seed the Database

Edit the connection details in `.env` file, then run:

```sh
go run cmd/seed/main.go
```

This will create the necessary tables and insert sample campaigns and targeting rules.

### 4. Run the Service

```sh
go run cmd/service/main.go
```

The service will start on `http://localhost:8080`.

---

## API

### `GET /v1/delivery?app=<app_id>&country=<country>&os=<os>`

Returns a list of matching campaigns for the given request parameters.

- **200 OK**: Returns matching campaigns.
- **204 No Content**: No campaigns match.
- **400 Bad Request**: Missing required parameters.

---
<!-- Awaiting impementatio 
## Monitoring (Optional)

### Prometheus & Grafana

1. Ensure Docker is installed and running.
2. Create `docker-compose.yml` and `prometheus.yml` as described in the documentation.
3. Start monitoring stack:

   ```sh
   docker-compose up
   ```

4. Access:
   - Prometheus: [http://localhost:9090](http://localhost:9090)
   - Grafana: [http://localhost:3000](http://localhost:3000)

---
--->
## Implementation Details

- **Go kit** is used for clear separation of service, endpoint, and transport layers.
- **PostgreSQL** stores all campaign and targeting rule data.
- **In-memory cache** ensures all delivery requests are served with minimal latency.
- **Hot reload** mechanism updates the in-memory cache from the database every minute, ensuring fresh data without restarts.
- **Prometheus metrics** are exposed at `/metrics` for observability.
#### Note

Implementing hot-reload almost halved the test execution time:

From:
```
PS C:\targeting-engine> go test ./internal/service
ok      github.com/glitchdawg/campaign-targeting-engine/internal/service        0.421s
```
To:
```
PS C:\targeting-engine> go test ./internal/service
ok      github.com/glitchdawg/campaign-targeting-engine/internal/service        0.286s
```

## Development & Testing

- Unit tests for business logic are in `internal/service/delivery_test.go` and can be run with:

  ```sh
  go test ./internal/service
  ```

---

## Further planned improvements

- Implement monitoring dashboards with Grafana and Prometheus
- Move from an in-memory cache to Redis for distributed caching
