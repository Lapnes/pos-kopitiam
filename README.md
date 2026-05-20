# POS KopiTiam v3NF

Enterprise POS system for coffee shop chain with 3NF database design.

## Tech Stack

- **Go 1.23** with Gin framework
- **MySQL 8.4** with GORM v2 ORM
- **JWT** authentication with role-based access control (RBAC)
- **Redis** (optional; retained for compatibility)
- **Swagger** API documentation
- **Midtrans Snap Sandbox Integration** for cashless payments (QRIS & Cards)

## Quick Start

### Prerequisites

- Go 1.23+
- MySQL 8.4 (or MariaDB)
- Redis 7+

### Setup

1. Clone and enter directory:
```bash
cd pos-kopitiam-v3nf
```

2. Copy environment file:
```bash
cp .env.example .env
```
Ensure you set the following Midtrans keys in `.env`:
```env
MIDTRANS_CLIENT_KEY=
MIDTRANS_SERVER_KEY=
```

3. Database Socket Connection (Optional):
The database connection dynamically supports both TCP host and UNIX socket files. If your environment isolates TCP loopbacks, update your `.env` to:
```env
DB_HOST=/tmp/mysql.sock
```

4. Start dependencies with Docker (or start local database):
```bash
docker-compose up -d mysql redis
```
*(Alternatively, run `./start_mysql.sh` to launch MariaDB locally).*

5. Run migrations:
```bash
go run cmd/migrate/main.go
```

6. Seed data:
```bash
go run seeds/seed.go
```

7. Start API server:
```bash
go run cmd/api/main.go
```

8. Start Caddy web server (for frontend assets):
```bash
cd FE/kopitiam-pos-fixed
caddy run
```

### API Documentation

Swagger UI available at: http://localhost:8080/swagger/index.html

### Default Login Credentials

| Client / Interface | Auth Type | Username / PIN | Password |
|---|---|---|---|
| Manager Panel | JWT (Email) | `manager@kopitiam.id` | `password123` |
| Cashier Panel | JWT (PIN) | `222222` | - |
| Kitchen Display (KDS) | JWT (PIN) | `333333` | - |
| Customer Self-Order Catalog | HTTP Basic Auth | `catalog` | `kopitiam123` |

### Catalog Page Entry Point
The customer self-order catalog runs under Basic Auth at:
* **Catalog URL:** `http://localhost:8080/catalog/`
* **Catalog API Endpoints:** `/catalog/api/menus` (GET), `/catalog/api/checkout` (POST)

## Project Structure

```
pos-kopitiam-v3nf/
├── cmd/
│   ├── api/          # HTTP server entry point
│   ├── migrate/      # Database migration
│   └── worker/       # Background job worker
├── FE/
│   └── kopitiam-pos-fixed/   # Frontend files (Caddy-served)
├── internal/
│   ├── config/       # App, DB, Redis config
│   ├── dto/          # Request/Response DTOs
│   ├── handler/      # HTTP handlers
│   ├── middleware/   # Auth & RBAC middleware
│   ├── models/       # GORM models
│   ├── repository/   # Data access layer
│   ├── routes/       # Route wiring
│   ├── service/      # Business logic
│   └── utils/        # JWT, password, response helpers
├── seeds/            # Database seeder
├── docs/             # Swagger docs
├── docker-compose.yml
├── go.mod
└── README.md
```

## Key Features

- **3NF Database**: Fully normalized with derived fields removed (rendered via VIEWs)
- **Transaction Safety**: Order confirmation, payment, and returns use atomic DB transactions
- **Anti N+1**: JOINs and subqueries for list endpoints
- **Dynamic Tax**: TaxConfig lookup per branch/global at calculation time
- **Stock Management**: Recipe-based BOM deduction + daily stock tracking
- **Shift Management**: Opening/closing with cash movement tracking
- **Split Payments**: Multiple payment methods per order
- **Midtrans Snap Integration**: Automatic token generation and payment popup overlay for Card and QRIS transactions.
- **Access Control Boundaries**: Custom KDS (Kitchen Display System) routing and page navigation tailored strictly to user roles (Manager, Cashier, Kitchen).
- **Audit Logs**: Immutable change tracking
- **Receipt & Kitchen Ticket Generation**: Formatted ticket builder ready for thermal printing hardware integration

## Documentation Reference

For details on installation, applied fixes, and history of modifications, refer to:
- [README_INSTALLATION.md](file:///home/tenzly/Documents/00_GOLang/02_KopiTiam_Project-1/pos-kopitiam-v3nf%20%282%29/README_INSTALLATION.md) — Step-by-step setup and troubleshooting.
- [README_N1_FIXES.md](file:///home/tenzly/Documents/00_GOLang/02_KopiTiam_Project-1/pos-kopitiam-v3nf%20%282%29/README_N1_FIXES.md) — Anti N+1 queries.
- [README_MIGRATION_VIEWS.md](file:///home/tenzly/Documents/00_GOLang/02_KopiTiam_Project-1/pos-kopitiam-v3nf%20%282%29/README_MIGRATION_VIEWS.md) — Schema views.

## License

MIT
