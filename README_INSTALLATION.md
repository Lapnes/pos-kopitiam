# 🛠️ POS KopiTiam v3NF — Installation & Troubleshooting Guide

This guide provides step-by-step instructions for setting up your environment, running migrations and seeders, running the application services, and troubleshooting common local/network errors.

---

## 📋 Prerequisites

Before starting, ensure you have installed:
* **Go 1.23+**
* **MySQL 8.4** (or MariaDB)
* **Caddy** (to host and proxy frontend assets)
* **Redis** (optional; used for legacy compatibility)

---

## 🚀 Step-by-Step Installation

### Step 1: Clone & Configure Environment
1. Enter the project root directory.
2. Copy the example environment template:
   ```bash
   cp .env.example .env
   ```
3. Open `.env` and configure your credentials. Specifically, set up the **Midtrans Snap Keys** and **Database credentials**:
   ```env
   APP_ENV=dev
   APP_PORT=8080
   JWT_SECRET=super-secret-jwt-key

   # Database settings
   DB_HOST=127.0.0.1
   DB_PORT=3306
   DB_USER=root
   DB_PASSWORD=popoi
   DB_NAME=pos_kopitiam

   # Midtrans Keys
   MIDTRANS_CLIENT_KEY=
   MIDTRANS_SERVER_KEY=
   MIDTRANS_ENV=sandbox
   ```

### Step 2: Database Initialization & Schema Views
We separate the initialization into GORM AutoMigration, native SQL migration constraints, and views generation.

1. **Run the GORM AutoMigrate script:**
   This creates all the raw tables and basic columns.
   ```bash
   go run cmd/migrate/main.go
   ```

2. **Apply the custom indexes and foreign keys:**
   Since GORM does not generate all specific foreign key `ON DELETE` restrict rules, run:
   ```bash
   mysql -h 127.0.0.1 -u root -p pos_kopitiam < migrations/001_add_indexes_and_fks.sql
   ```

3. **Apply the analytics views (Dynamic Standard Views):**
   Run this script to register the 8 dynamic analytics views used by the Manager Dashboard:
   ```bash
   mysql -h 127.0.0.1 -u root -p pos_kopitiam < migrations/002_create_analytics_views.sql
   ```

4. **Seed the database:**
   Populate clean branches, employees, menu items, categories, taxes, and raw material recipe ingredients:
   ```bash
   go run seeds/seed.go
   ```

### Step 3: Run the Application
1. **Start the Go REST API Server:**
   ```bash
   go run cmd/api/main.go
   ```
   *The API will run on `http://localhost:8080` (or the configured `APP_PORT`).*

2. **Start the Caddy Frontend Server:**
   Open a separate terminal window:
   ```bash
   cd FE/kopitiam-pos-fixed
   caddy run
   ```
   *Caddy will serve static assets on `http://localhost:8090` and proxy backend requests under `/api/*` and `/catalog/*`.*

---

## 🔍 Validation & Tests

To automatically test and verify all endpoint handlers (including authentication, orders, payments, void/return, and analytics):
```bash
cd docs/API
go run api_runner.go
```
*Make sure your database is in a clean state before running tests to prevent stock level depletion conflicts.*

---

## 🔧 Troubleshooting Common Errors

### 1. Employee Creation Fails (`Error 1452 (23000): Cannot add or update a child row: a foreign key constraint fails`)
* **Symptom:** Adding an employee from the Manager Dashboard fails with a `500 Server Error`. The logs show an insert fail on `fk_users_branch`.
* **Cause:** If you drop and re-seed the database, the seeder generates a **new UUID for the branch**. However, if you are still logged into the dashboard, your browser holds the old session JWT containing the **old branch ID**. When creating an employee, the dashboard submits this cached old branch ID, causing the database to reject the row.
* **Solution:** **Log Out** of the Manager Dashboard and **Log Back In** to refresh your browser's cached JWT token.
  * *Manager Email:* `manager@kopitiam.id`
  * *Password:* `password123`

### 2. QRIS/Online Payment Fails (`Cannot read properties of undefined (reading 'pay')`)
* **Symptom:** When accessing the cashier or self-service catalog from another device on the local network (e.g. `http://192.168.101.67:8090/`), clicking online checkout fails.
* **Cause:** The Midtrans Snap SDK script (`snap.js`) is loaded from Midtrans' external CDN. If the device accessing the catalog does not have active **internet access**, it cannot download the script, making `window.snap` undefined.
* **Solution:** 
  1. Verify the client device has internet access.
  2. Open the browser's developer console (F12) on the device to check for network/CORS blocks on `https://app.sandbox.midtrans.com/snap/snap.js`.
  3. Try exposing the local environment through an HTTPS tunnel (e.g., using `ngrok http 8090`) to avoid browser restrictions on serving payment scripts inside insecure HTTP IP contexts.

### 3. Database Connection Fails in Sandboxed/Isolated Environments
* **Symptom:** Running the migration/app throws `Can't connect to server on '127.0.0.1:3306'`.
* **Cause:** In isolated or namespace-restricted pipeline environments, TCP loopbacks are blocked.
* **Solution:** The database driver natively supports UNIX socket connections. Update your `.env` configuration to use the socket file path directly as the host:
  ```env
  DB_HOST=/tmp/mysql.sock
  ```
  *(Make sure the database socket path is correct and accessible).*
