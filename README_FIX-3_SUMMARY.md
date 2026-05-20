# POS KopiTiam v3NF — Fix Phase 3 Summary

This release contains the integration of cashless payments via Midtrans Snap and critical hotfixes for the terminal cashier and kitchen display panels.

---

## 1. Cashless Payment: Midtrans Snap Integration

### Root Cause / Requirement:
- The cashier terminal needed a secure sandbox payment gateway to handle Card and QRIS transactions, rather than processing them as simple manual cash splits.

### Fix / Implementation:
- **Environment & Configuration**:
  - Appended `MIDTRANS_CLIENT_KEY`, `MIDTRANS_SERVER_KEY`, and `MIDTRANS_ENV` credentials/settings to the `.env` configuration file.
  - Parsed and propagated keys through the `config.go` struct.
  - Dynamically switches endpoints between Sandbox (`https://app.sandbox.midtrans.com/snap/v1/transactions`) and Production (`https://app.midtrans.com/snap/v1/transactions`) depending on the `MIDTRANS_ENV` value.
- **Dynamic Frontend Script Loader (`midtrans-loader.js`)**:
  - Registered `/api/v1/config/midtrans` backend route in the Gin engine to expose `client_key` and `environment` config.
  - Created `midtrans-loader.js` dynamic script injector for frontend pages (`cashier.html`, `catalog.html`). This dynamically fetches keys from the backend and loads the correct version of the Midtrans Snap SDK (sandbox or production) without hardcoding client keys.
- **Transaction Dispatcher (`payment_service.go`)**:
  - Created a client helper function `requestMidtransSnap` using `net/http` to POST transaction requests to the configured sandbox or production API endpoint.
  - When a transaction with online payment splits (`qris`, `debit_card`, or `credit_card`) is processed, the system fetches a valid Snap Token and passes it back to the client.
- **Frontend overlay UI (`cashier.html` & `cashier.js`)**:
  - Integrated the dynamic loader in the HTML files.
  - Handled checkout responses inside the cashier script; if a `snap_token` is present, it launches the interactive `window.snap.pay(...)` modal. Cashier UI hooks handle successful, pending, cancelled, or closed checkout states accordingly.

---

## 2. Interactive Terminal UI Fixes

### 1. Stuck Checkout Button (Processing... Freeze)
- **Root Cause**: When the cashier clicked the checkout button, the content was replaced by a loader animation, which stripped the `#charge-text` node from the DOM. Clearing the cart triggered `updateTotals()`, which crashed on a `TypeError` when reading `textContent` from the missing node, causing the UI to freeze in a `Processing...` state.
- **Fix**: Added a defensive query selector verification to ensure `$('charge-text')` exists before updating its `textContent` in `cashier.js`.

### 2. KDS Card Price Rendering (Rp 0 Bug)
- **Root Cause**: The kitchen panel displayed a total of `Rp 0` for all order cards because the template was querying `order.total_amount`. Since GORM serializes the total as `"total"`, the frontend read `undefined` and defaulted to `Rp 0`.
- **Fix**: Refactored the card template renderer in `kitchen.js` to dynamically fall back: `order.total || order.total_amount || 0`.


### 3. Self-Service Catalog with Table Numbers
- **Requirement**: The customer-facing digital catalog must act as a self-service panel. Orders require table number inputs, but database schema modifications were restricted.
- **Fix**:
  - Appended a required `table_number` field inside `CatalogCheckoutRequest`.
  - Inside the Go API (`catalog_handler.go`), generated a deterministic SHA1 UUID based on the table number to populate the database's existing `table_id` field.
  - Stored the readable table designation directly inside the order's `Notes` field (e.g. `"Meja 5 | Self-ordered from catalog tablet"`) so it prints automatically on kitchen tickets, cashier screens, and KDS cards.
  - Added a "Nomor Meja" text input field on the catalog frontend layout (`catalog.html` / `catalog.js`) with input validation.

### 4. Unified Online / Midtrans Payment Option
- **Root Cause / Requirement**: Forcing granular payment filters (QRIS, VA Bank, Cards) led to "No payment channels available" errors if specific options were unconfigured on sandbox merchant credentials. The user requested simplifying Card payments to "Online/Midtrans" allowing customers to select payment options dynamically inside Midtrans' secure Snap overlay.
- **Fix**:
  - Removed explicit `enabled_payments` restrictions from the backend Snap payload generator (`payment_service.go`). Midtrans now automatically shows all payment options configured on the MAP dashboard.
  - Replaced the distinct QRIS, VA, and GoPay buttons on the catalog page with a single clean "Online / Midtrans" toggle.
  - Changed the Cashier card selector button label to "Online/Midtrans".
  - Configured Caddy file to proxy `/catalog/*` requests directly to GIN (port 8080) to support seamless Basic Auth login across both ports.

### 5. Automated Audit Log System
- **Requirement**: Track database mutations (CREATE, UPDATE, DELETE) dynamically and display them inside a stylized scrollable terminal box in the Manager Dashboard without causing layout shifts.
- **Fix**:
  - Implemented an asynchronous **Audit Middleware** (`audit.go`) registered globally under the `/api/v1` router group.
  - Automatically captures HTTP method types, intercepts request body streams, maps endpoints to target entities (e.g. `menus`, `employees`, `discounts`), checks response status validity (only logs successful 2xx changes), parses the active JWT `userID` context, and inserts records into the `audit_logs` table in a background goroutine.
  - Manually logs catalog self-service checkout mutations within `catalog_handler.go` utilizing the context-fallback cashier ID.
  - Added a GET `/api/v1/analytics/audit-logs` endpoint under the `/analytics` group, protected by manager role constraints.
  - Placed a high-fidelity Unix terminal UI widget inside the Manager Dashboard overview (`analytics.html`). Features custom styling mimicking developer console headers (three-dot chrome headers), internal scroll constraints (`max-h-80 overflow-y-auto`), color-coded terminal text mapping actions to different color categories, and manual refresh controls.
  - **Dynamic Range Slider**: Implemented a responsive custom range slider on the terminal chrome header (ranging from 100 to 1000 items with a step size of 100) allowing managers to fetch and view larger chunks of system history instantly, complete with automatic data reload on value changes.
  - **Robust JSON Handling**: Integrated default `"null"` JSON string fallbacks to avoid MySQL JSON column validation errors (`Error 3140`) on empty payload insertions.

---

## 3. Sandboxed Infrastructure Improvements

### UNIX Socket Database Connections (`database.go`)
- **Root Cause**: Test pipelines and sandboxed environments isolate network namespaces, blocking connection attempts to loopback TCP ports (e.g. `127.0.0.1:3306`).
- **Fix**: Upgraded database configuration string builder (`DSN()`) to support unix socket file paths natively. When `DB_HOST` in `.env` starts with `/` or ends with `.sock`, it automatically wraps the database connection using GORM's `unix()` protocol connector, allowing the application to connect via `/tmp/mysql.sock`.

---

## 4. Operational Commands & Setup Reference

For full details on configuring `.env` variables, running database migrations, and executing seeders, please refer to the main **[README_INSTALLATION.md](file:///home/tenzly/Documents/00_GOLang/02_KopiTiam_Project-1/pos-kopitiam-v3nf%20%282%29/README_INSTALLATION.md)**.

### Quick Start Services (Go API + Caddy Frontend)
Once the database is initialized and seeded:
```bash
# 1. Build and launch Go API server (will connect via /tmp/mysql.sock or TCP port 3306)
go run cmd/api/main.go

# 2. Launch Caddy server to serve frontend assets (on port 8090)
cd FE/kopitiam-pos-fixed
caddy run
```
