# ☕ POS Kopitiam

Sistem Point of Sale (POS) modern berbasis web untuk kedai kopi, dibangun dengan Golang + Next.js dan MySQL via Docker.

---

## 🧱 Tech Stack

| Layer      | Teknologi                                              |
|------------|--------------------------------------------------------|
| Backend    | Go 1.25 · Gin · GORM · JWT · Asynq (worker)           |
| Frontend   | Next.js 16 (App Router) · React 19 · TypeScript        |
| Styling    | Tailwind CSS · Shadcn UI · Recharts                    |
| State      | Zustand (cart & auth) · TanStack Query (data fetching) |
| Database   | MySQL 8.0 via Docker                                   |
| Cache/Queue| Redis 7 via Docker                                     |
| Container  | Docker · Docker Compose                                |

---

## 📁 Struktur Folder

```
pos-kopitiam/
│
├── cmd/                          ← Entry points aplikasi Go
│   ├── api/main.go               ← HTTP API server (port 8080)
│   ├── migrate/main.go           ← Database migration runner
│   └── worker/main.go            ← Background worker (Asynq)
│
├── internal/                     ← Kode inti backend (Clean Architecture)
│   ├── config/
│   │   ├── config.go             ← Load env variables
│   │   ├── database.go           ← Inisialisasi MySQL (GORM)
│   │   └── redis.go              ← Inisialisasi Redis
│   ├── dto/
│   │   ├── request.go            ← Struct request (Login, Order, dll)
│   │   └── response.go           ← Struct response standar
│   ├── handler/                  ← HTTP handler per domain
│   │   ├── analytics_handler.go
│   │   ├── auth_handler.go
│   │   ├── employee_handler.go
│   │   ├── inventory_handler.go
│   │   ├── menu_handler.go
│   │   ├── order_handler.go
│   │   ├── printer_handler.go
│   │   ├── return_handler.go
│   │   └── shift_handler.go
│   ├── middleware/
│   │   ├── auth.go               ← JWT authentication middleware
│   │   └── rbac.go               ← Role-based access control
│   ├── models/
│   │   ├── domain.go             ← Semua GORM model (tabel DB)
│   │   └── enums.go              ← Konstanta role, status, dll
│   ├── repository/               ← Data access layer (query DB)
│   │   ├── analytics_repository.go
│   │   ├── base.go
│   │   ├── inventory_repository.go
│   │   ├── menu_repository.go
│   │   ├── order_repo.go
│   │   ├── return_repository.go
│   │   ├── shift_repository.go
│   │   └── user_repo.go
│   ├── routes/
│   │   └── routes.go             ← Definisi semua route API
│   ├── service/                  ← Business logic layer
│   │   ├── analytics_service.go
│   │   ├── auth_service.go
│   │   ├── employee_service.go
│   │   ├── inventory_service.go
│   │   ├── menu_service.go
│   │   ├── order_service.go
│   │   ├── printer_service.go
│   │   ├── return_service.go
│   │   └── shift_service.go
│   └── utils/
│       ├── jwt.go                ← Generate & validasi JWT
│       ├── password.go           ← Hash & verify password (bcrypt)
│       └── response.go           ← Helper JSON response
│
├── migrations/                   ← Skema migrasi database manual
│   └── schema.sql
│
├── pkg/                          ← Eksternal library/paket
│   └── bigcapital/               ← Integrasi BigCapital (third-party)
│
├── seeds/
│   └── seed.go                   ← Data awal: branch, user, menu, dll
│
├── fe/frontend/                  ← Next.js frontend
│   ├── src/
│   │   ├── app/                  ← Next.js App Router
│   │   │   ├── (auth)/login/     ← Halaman login
│   │   │   ├── admin/            ← Dashboard Super Admin (CRUD karyawan)
│   │   │   ├── kasir/            ← Dashboard Kasir (POS + refund)
│   │   │   ├── kitchen/          ← Dashboard Dapur (order tickets)
│   │   │   ├── manager/          ← Dashboard Manager (analitik)
│   │   │   ├── layout.tsx        ← Root layout + Providers
│   │   │   ├── page.tsx          ← Root redirect berdasarkan role
│   │   │   └── globals.css       ← Global styles + Tailwind
│   │   ├── components/
│   │   │   ├── auth/             ← LoginForm
│   │   │   ├── analytics/        ← AnalyticsDashboard
│   │   │   ├── employees/        ← EmployeePanel, EditEmployeeDialog
│   │   │   ├── layout/           ← DashboardLayout (sidebar + topbar)
│   │   │   ├── pos/              ← MenuCard, CheckoutDialog, HistoryDialog,
│   │   │   │                        RefundByOrderNumberDialog, dll
│   │   │   ├── ui/               ← Shadcn UI components
│   │   │   └── Providers.tsx     ← QueryClientProvider wrapper
│   │   ├── lib/
│   │   │   ├── api/axios.ts      ← Axios instance + interceptors (JWT inject)
│   │   │   └── utils.ts          ← Helper utilities
│   │   ├── store/
│   │   │   ├── useAuthStore.ts   ← Zustand auth store (user + token)
│   │   │   └── useCartStore.ts   ← Zustand cart store (items + total)
│   │   └── types/
│   │       └── index.ts          ← TypeScript types (Menu, Order, User, dll)
│   ├── .env.local                ← Environment frontend
│   ├── next.config.ts
│   └── package.json
│
├── backend/.env                  ← Env khusus folder backend (legacy)
├── pos_kopitiam_full.sql         ← Database dump (legacy SQL)
├── pos_kopitiam_full.zip         ← Backup / arsip project
├── .env                          ← Environment utama backend (root)
├── .env.example                  ← Template env
├── docker-compose.yml            ← MySQL + Redis containers
├── go.mod                        ← Go module dependencies
├── go.sum
└── README.md
```

---

## ⚙️ Environment Setup

### 1. Backend — `.env` (root)

```env
APP_ENV=dev
APP_PORT=8080
APP_SECRET=pos_kopitiam_app_secret_2026

DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=root
DB_PASSWORD=password
DB_NAME=pos_kopitiam

REDIS_HOST=127.0.0.1
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

JWT_SECRET=pos_kopitiam_super_secret_key_2026
JWT_EXPIRY=8h
REFRESH_TOKEN_EXPIRY=168h
```

### 2. Frontend — `fe/frontend/.env.local`

```env
NEXT_PUBLIC_API_URL=http://localhost:8080
```

---

## 🚀 Cara Menjalankan (Local Development)

### Prasyarat

- Go ≥ 1.21
- Node.js ≥ 20 + npm
- Docker & Docker Compose

> ⚠️ **PENTING**: Anda membutuhkan **minimal 2 terminal/tab terpisah** untuk menjalankan Backend (API) dan Frontend secara bersamaan. Tanpa menjalankan API Backend, Frontend akan gagal memuat data dari database.

---

### Step 1 — Jalankan Database & Redis (Terminal 1)

```bash
docker compose up -d
```

Tunggu MySQL siap (sekitar 15–30 detik):

```bash
docker exec pos_kopitiam_db mysqladmin ping -h localhost -u root -ppassword --silent
# Output: mysqld is alive
```

---

### Step 2 — Migrasi Skema Database (Terminal 1)

> ⚠️ Perintah ini akan **drop semua tabel lama** lalu buat ulang. Hanya jalankan sekali atau saat schema berubah.

```bash
go run cmd/migrate/main.go
```

Output sukses:
```
Auto Migration Completed Successfully!
```

---

### Step 3 — Isi Data Awal / Seed (Terminal 1)

```bash
go run seeds/seed.go
```

Data yang di-seed:
| Data      | Detail                                                      |
|-----------|-------------------------------------------------------------|
| Branch    | Main Branch                                                 |
| Karyawan  | Super Admin, Manager, Cashier 1, Kitchen 1                  |
| Area      | Indoor                                                      |
| Meja      | Table 1, Table 2                                            |
| Kategori  | Coffee, Tea, Iced, Foods, Snacks, Desserts, Juices          |
| Menu      | 10 item (Milk Coffee, Black Coffee, Pulled Tea, Fried Rice, dll) |

Output sukses:
```
Seeding Completed!
```

---

### Step 4 — Jalankan Backend API (Terminal 1)

**Ini Wajib Dijalankan agar frontend bisa mengambil data!**

```bash
go run cmd/api/main.go
```

Server berjalan di: `http://localhost:8080`

Verifikasi:
```bash
curl http://localhost:8080/health
# {"status":"up","redis":true,"db":true}
```

> **Biarkan terminal ini tetap menyala/running.**

---

### Step 5 — Jalankan Background Worker (Terminal 2 - Opsional)

Buka tab terminal baru. Worker digunakan untuk task async (sinkronisasi BigCapital, dll):

```bash
go run cmd/worker/main.go
```

---

### Step 6 — Jalankan Frontend (Terminal 3)

Buka tab terminal baru (atau terminal kedua jika step 5 dilewati), lalu masuk ke folder frontend:

```bash
cd fe/frontend
npm install      # Pertama kali saja
npm run dev      # Development server
```

Frontend berjalan di: `http://localhost:3000`

> Jika Anda hanya ingin tes kasir, minimal Step 4 (Backend API) dan Step 6 (Frontend) **harus aktif secara bersamaan di terminal yang berbeda**.

Untuk production build:
```bash
npm run build
npx next start   # Jalankan production build
```

---

## 🔑 Akun Default (Hasil Seed)

| Role        | Email                       | Password      | PIN    |
|-------------|-----------------------------|---------------|--------|
| Super Admin | admin@kopitiam.com          | password123   | 111111 |
| Manager     | manager@kopitiam.com        | password123   | 222222 |
| Kasir       | cashier1@kopitiam.com       | password123   | 333333 |
| Kitchen     | kitchen1@kopitiam.com       | password123   | 444444 |

---

## 🌐 API Endpoints Utama

```
POST   /api/v1/auth/login              ← Login dengan email + password
POST   /api/v1/auth/login-pin          ← Login dengan PIN (kasir)
POST   /api/v1/auth/logout

GET    /api/v1/menus                   ← Daftar menu aktif
GET    /api/v1/orders                  ← Daftar order
POST   /api/v1/orders                  ← Buat order baru
PUT    /api/v1/orders/:id/confirm      ← Konfirmasi order (dapur selesai)
PUT    /api/v1/orders/:id/cancel       ← Batalkan order
POST   /api/v1/orders/:id/void         ← Void order (manager)
POST   /api/v1/orders/:id/returns      ← Proses retur/refund
GET    /api/v1/orders/by-number/:num   ← Cari order by nomor struk

GET    /api/v1/employees               ← Daftar karyawan (admin only)
POST   /api/v1/employees               ← Tambah karyawan
PUT    /api/v1/employees/:id           ← Edit karyawan
DELETE /api/v1/employees/:id           ← Hapus karyawan

GET    /api/v1/analytics/sales-summary ← Ringkasan penjualan
GET    /api/v1/analytics/best-sellers  ← Menu terlaris
GET    /api/v1/analytics/return-impact ← Dampak retur

GET    /api/v1/shifts/current          ← Shift aktif
POST   /api/v1/shifts/open             ← Buka shift
POST   /api/v1/shifts/close            ← Tutup shift

GET    /health                         ← Health check
```

---

## 🎭 Role & Akses Halaman

| Role        | Halaman yang Dapat Diakses                              |
|-------------|----------------------------------------------------------|
| `superadmin`| `/admin`, `/manager`, `/kasir`, `/kitchen`              |
| `manager`   | `/manager`                                               |
| `cashier`   | `/kasir`                                                 |
| `kitchen`   | `/kitchen`                                               |

---

## 🐳 Docker Services

```bash
docker compose up -d          # Jalankan semua container
docker compose down           # Hentikan container (data tetap ada)
docker compose down -v        # Hentikan + hapus semua data volume
docker compose ps             # Cek status container
docker compose logs -f mysql  # Log MySQL real-time
```

| Container          | Image         | Port Host → Container |
|--------------------|---------------|-----------------------|
| `pos_kopitiam_db`  | mysql:8.0     | 3306 → 3306           |
| `pos_kopitiam_redis`| redis:7-alpine| 6379 → 6379          |

---

## 🔍 Troubleshooting

### Backend gagal konek ke database

```bash
# Cek container aktif
docker compose ps

# Ping MySQL
docker exec pos_kopitiam_db mysqladmin ping -h localhost -u root -ppassword

# Pastikan .env di root sudah benar
cat .env | grep DB_
```

### Port 8080 sudah dipakai

```bash
lsof -i :8080
kill -9 <PID>
```

### Port 3000 sudah dipakai (Next.js fallback ke 3001)

```bash
lsof -i :3000
kill -9 <PID>
# Lalu jalankan ulang: npm run dev
```

### Login frontend gagal (401 Unauthorized)

Pastikan:
1. Backend berjalan di port 8080
2. `fe/frontend/.env.local` berisi `NEXT_PUBLIC_API_URL=http://localhost:8080`
3. Login menggunakan **email** (bukan nama), contoh: `admin@kopitiam.com`

### Analytics error 500

```bash
# Jalankan ulang migrate untuk buat tabel order_returns
go run cmd/migrate/main.go
go run seeds/seed.go
```

---

## 🏗️ Arsitektur

```
Browser (Next.js :3000)
        │
        │  HTTP + Bearer JWT
        ▼
Backend API (Gin :8080)
        │
        ├── MySQL :3306   ← Data utama (orders, menus, employees)
        ├── Redis :6379   ← Cache + job queue (Asynq)
        └── Worker        ← Background tasks (async processing)
```

---

## 📝 Lisensi

MIT License — Dibuat untuk keperluan akademis MBD (Manajemen Basis Data).
