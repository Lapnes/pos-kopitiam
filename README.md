# POS KopiTiam Project

POS KopiTiam adalah sistem Point of Sale (POS) sederhana untuk kafe KopiTiam, dibangun menggunakan bahasa pemrograman Go. Sistem ini menyediakan API RESTful untuk mengelola pesanan, data master, dan laporan penjualan.

## Fitur Utama

### 1. Manajemen Pesanan (Orders)
- **Buat Pesanan Baru**: Membuat pesanan dengan daftar item menu, quantity, dan otomatis menghitung total harga.
- **Lihat Semua Pesanan**: Mengambil daftar semua pesanan yang ada.
- **Lihat Pesanan Berdasarkan ID**: Mengambil detail pesanan tertentu berdasarkan ID unik.
- **Laporan Pesanan**: Menampilkan laporan penjualan dengan detail item per pesanan.

### 2. Data Master
- **Menu**: Daftar menu makanan dan minuman dengan harga dan stok harian.
- **Kategori**: Kategori untuk mengelompokkan menu (misalnya: Makanan, Minuman).
- **Karyawan**: Daftar karyawan yang dapat ditugaskan untuk pesanan.

### 3. Fitur Tambahan
- **Health Check**: Endpoint untuk memeriksa status kesehatan API.
- **Database Migration**: Migrasi skema database otomatis.
- **Seeding Data**: Pengisian data awal untuk development (hanya di environment dev).
- **Validasi Stok**: Memastikan stok harian cukup sebelum membuat pesanan.
- **Transaksi Database**: Menggunakan transaksi untuk memastikan konsistensi data saat membuat pesanan.

## Teknologi yang Digunakan

- **Bahasa Pemrograman**: Go 1.25
- **Web Framework**: Gin
- **ORM**: GORM
- **Database**: MySQL 8.4 (via Docker)
- **Frontend**: Next.js 16, React 19, TypeScript, Tailwind CSS v4
- **UUID**: Untuk generate ID unik pesanan
- **Environment Config**: Menggunakan .env file

## Struktur Proyek

```
pos_kopitiam/
├── cmd/
│   ├── api/v1/
│   │   └── main.go          # Entry point API
│   └── db/
│       └── runMigration.go  # Script migrasi database
├── internal/
│   ├── config/
│   │   ├── config.go        # Konfigurasi aplikasi
│   │   ├── database.go      # Koneksi database
│   │   ├── db_logic.go      # Logika database
│   │   ├── fk_constraint.go # Foreign key constraints
│   │   └── seeder.go        # Data seeding
│   ├── dto/
│   │   └── report.go        # Data Transfer Objects
│   ├── models/
│   │   └── domain.go        # Model domain
│   └── service/
│       └── order_service.go # Business logic untuk order
├── app/                     # Next.js app directory
│   ├── globals.css          # Global styles & CSS variables
│   ├── layout.tsx           # Root layout
│   └── page.tsx             # Halaman utama POS
├── components/              # React components
│   ├── CartPanel.tsx
│   ├── MenuCard.tsx
│   ├── Modal.tsx
│   ├── OrdersTab.tsx
│   └── Toast.tsx
├── lib/
│   └── api.ts               # API helper & utilities
├── types/
│   └── index.ts             # TypeScript types
├── docker-compose.yml       # MySQL via Docker
├── go.mod                   # Go module file
├── package.json             # Node.js dependencies
├── pos_kopitiam_full.sql    # Schema database
└── README.md                # Dokumentasi ini
```

---

## Setup dan Instalasi Lengkap

### Prasyarat

Pastikan semua tool berikut sudah terinstall di sistem kamu:

| Tool | Versi Minimum | Cek Versi |
|------|---------------|-----------|
| Go | 1.21+ | `go version` |
| Node.js | 18+ | `node --version` |
| npm | 9+ | `npm --version` |
| Docker | 24+ | `docker --version` |
| Docker Compose | v2+ | `docker compose version` |
| Git | - | `git --version` |

---

## Langkah 1 — Clone Repository

```bash
git clone https://github.com/Lapnes/pos-kopitiam.git
cd pos-kopitiam
```

---

## Langkah 2 — Konfigurasi Environment (.env)

Buat file `.env` di root directory project:

```bash
# Buat file .env
touch .env
```

Isi file `.env` dengan konfigurasi berikut:

```env
# Database
DB_HOST=localhost
DB_PORT=4406
DB_USER=root
DB_PASSWORD=rahasia123
DB_NAME=pos_kopitiam

# App environment: "dev" akan auto-seed data awal, "prod" tidak
APP_ENV=dev

# Frontend — URL backend (digunakan oleh Next.js)
NEXT_PUBLIC_API_URL=http://localhost:3400/api/v1
```

> **Catatan port:**
> - `DB_PORT=4406` adalah port host Docker (MySQL container berjalan di 4406:3306)
> - Backend Go berjalan di `:3400`
> - Frontend Next.js berjalan di `:3000`

---

## Langkah 3 — Jalankan Database (MySQL via Docker)

```bash
# Jalankan MySQL container di background
docker compose up -d

# Verifikasi container berjalan
docker compose ps
```

Output yang diharapkan:
```
NAME                   STATUS          PORTS
pos_kopitiam_full      running         0.0.0.0:4406->3306/tcp
```

Tunggu sekitar 10–15 detik agar MySQL selesai inisialisasi sebelum lanjut ke langkah berikutnya.

> **Cek koneksi database (opsional):**
> ```bash
> docker exec -it pos_kopitiam_full mysql -uroot -p${DB_PASSWORD} -e "SHOW DATABASES;"
> ```

---

## Langkah 4 — Setup Backend (Go)

### 4a. Install Go Dependencies

```bash
go mod tidy
```

Perintah ini akan mengunduh semua package yang diperlukan (Gin, GORM, uuid, godotenv, dll).

### 4b. (Opsional) Jalankan Migrasi Manual

Jika kamu ingin menjalankan migrasi schema secara manual:

```bash
go run cmd/db/runMigration.go
```

> **Catatan:** Saat `APP_ENV=dev`, backend akan otomatis menjalankan seed data saat pertama kali dijalankan. Jadi langkah ini opsional.

### 4c. Jalankan Backend Server

```bash
go run cmd/api/v1/main.go
```

Output yang diharapkan:
```
[GIN-debug] [WARNING] Creating an engine instance with the Logger and Recovery middleware already attached.
[GIN-debug] GET    /health
[GIN-debug] POST   /api/v1/orders
[GIN-debug] GET    /api/v1/orders
...
[GIN-debug] Listening and serving HTTP on :3400
```

Verifikasi backend berjalan:

```bash
curl http://localhost:3400/health
# Response: {"status":"ok","message":"POS KopiTiam API is running"}
```

> **Biarkan terminal ini tetap berjalan.** Buka terminal baru untuk langkah berikutnya.

---

## Langkah 5 — Setup Frontend (Next.js)

Buka terminal baru, masuk ke folder project yang sama:

```bash
cd pos-kopitiam
```

### 5a. Install Node Dependencies

```bash
npm install
```

Ini akan menginstall Next.js, React, Tailwind CSS, lucide-react, dan semua dependensi frontend.

### 5b. Jalankan Development Server

```bash
npm run dev
```

Output yang diharapkan:
```
  ▲ Next.js 16.2.5
  - Local:        http://localhost:3000
  - Network:      http://192.168.x.x:3000

 ✓ Starting...
 ✓ Ready in 1.2s
```

---

## Langkah 6 — Buka di Browser

Buka browser dan akses:

```
http://localhost:3000
```

Kamu akan melihat tampilan POS KopiTiam dengan tema **Hijau Sage + Krem Hangat**.

---

## Ringkasan Urutan Menjalankan

Selalu ikuti urutan ini setiap kali ingin menjalankan project:

```
1. docker compose up -d          ← Database (MySQL)
        ↓ tunggu ~10 detik
2. go run cmd/api/v1/main.go     ← Backend API (:3400)
        ↓ buka terminal baru
3. npm run dev                   ← Frontend Next.js (:3000)
        ↓
4. Buka http://localhost:3000
```

> **Tips:** Gunakan tmux, atau buka 3 tab terminal terpisah agar semua proses bisa berjalan bersamaan.

---

## Menghentikan Semua Service

```bash
# Stop frontend: Ctrl+C di terminal npm run dev

# Stop backend: Ctrl+C di terminal go run

# Stop database
docker compose down

# Stop database + hapus data volume (reset total)
docker compose down -v
```

---

## Troubleshooting

### ❌ Error: `dial tcp 127.0.0.1:4406: connect: connection refused`

MySQL belum siap. Tunggu 15–20 detik setelah `docker compose up -d`, lalu coba lagi.

### ❌ Error: `Error 1049: Unknown database 'pos_kopitiam'`

Database belum dibuat. Cek file `.env`, pastikan `DB_NAME=pos_kopitiam` dan jalankan ulang:

```bash
docker compose down -v
docker compose up -d
```

### ❌ Error: `port 3400 already in use`

Ada proses lain di port 3400. Cari dan matikan:

```bash
# Linux/macOS
lsof -i :3400 | grep LISTEN
kill -9 <PID>

# Windows
netstat -ano | findstr :3400
taskkill /PID <PID> /F
```

### ❌ Frontend tidak bisa fetch data (menu kosong)

Pastikan:
1. Backend berjalan di `:3400` → `curl http://localhost:3400/health`
2. File `.env` punya `NEXT_PUBLIC_API_URL=http://localhost:3400/api/v1`
3. Restart `npm run dev` setelah mengubah `.env`

### ❌ Error: `go: module not found`

```bash
go mod tidy
go mod download
```

### ❌ Error: `npm: command not found`

Install Node.js dari [nodejs.org](https://nodejs.org) (pilih versi LTS).

---

## Dokumentasi API

Base URL: `http://localhost:3400`

### Health Check

#### GET /health
Memeriksa status kesehatan API.

**Response:**
```json
{
  "status": "ok",
  "message": "POS KopiTiam API is running"
}
```

### Orders

#### POST /api/v1/orders
Membuat pesanan baru.

**Request Body:**
```json
{
  "employee_id": 1,
  "items": [
    {
      "menu_id": 1,
      "qty": 2
    },
    {
      "menu_id": 2,
      "qty": 1
    }
  ]
}
```

**Response:**
```json
{
  "order_id": "550e8400-e29b-41d4-a716-446655440000",
  "employee_id": 1,
  "total_price": 50000,
  "items": [
    {
      "menu_id": 1,
      "menu_name": "Nasi Goreng",
      "quantity": 2,
      "unit_price": 20000,
      "subtotal": 40000
    },
    {
      "menu_id": 2,
      "menu_name": "Teh Manis",
      "quantity": 1,
      "unit_price": 10000,
      "subtotal": 10000
    }
  ]
}
```

#### GET /api/v1/orders
Mengambil semua pesanan.

**Response:**
```json
[
  {
    "order_id": "550e8400-e29b-41d4-a716-446655440000",
    "employee_id": 1,
    "total_price": 50000,
    "employee": {
      "employee_id": 1,
      "employee_name": "John Doe",
      "phone_number": "08123456789"
    },
    "order_details": [
      {
        "order_detail_id": "550e8400-e29b-41d4-a716-446655440001",
        "menu_id": 1,
        "quantity": 2,
        "unit_price": 20000,
        "subtotal": 40000,
        "menu": {
          "menu_id": 1,
          "menu_name": "Nasi Goreng",
          "price": 20000,
          "daily_stock": 48
        }
      }
    ]
  }
]
```

#### GET /api/v1/orders/{id}
Mengambil pesanan berdasarkan ID.

**Response:** Sama dengan response GET /api/v1/orders untuk satu item.

#### GET /api/v1/orders/report
Mengambil laporan penjualan.

**Response:**
```json
[
  {
    "order_id": "550e8400-e29b-41d4-a716-446655440000",
    "employee_name": "John Doe",
    "menu_name": "Nasi Goreng",
    "quantity": 2,
    "unit_price": 20000,
    "subtotal": 40000,
    "total_price": 50000
  }
]
```

### Master Data

#### GET /api/v1/menus
Mengambil semua menu.

**Response:**
```json
[
  {
    "menu_id": 1,
    "category_id": 1,
    "menu_name": "Nasi Goreng",
    "price": 20000,
    "daily_stock": 50,
    "category": {
      "category_id": 1,
      "category_name": "Makanan"
    }
  }
]
```

#### GET /api/v1/categories
Mengambil semua kategori.

**Response:**
```json
[
  {
    "category_id": 1,
    "category_name": "Makanan",
    "menus": [
      {
        "menu_id": 1,
        "menu_name": "Nasi Goreng",
        "price": 20000,
        "daily_stock": 50
      }
    ]
  }
]
```

#### GET /api/v1/employees
Mengambil semua karyawan.

**Response:**
```json
[
  {
    "employee_id": 1,
    "employee_name": "John Doe",
    "phone_number": "08123456789"
  }
]
```

---

## Model Data

### Category
```json
{
  "category_id": 1,
  "category_name": "Makanan"
}
```

### Menu
```json
{
  "menu_id": 1,
  "category_id": 1,
  "menu_name": "Nasi Goreng",
  "price": 20000,
  "daily_stock": 50
}
```

### Employee
```json
{
  "employee_id": 1,
  "employee_name": "John Doe",
  "phone_number": "08123456789"
}
```

### Order
```json
{
  "order_id": "550e8400-e29b-41d4-a716-446655440000",
  "employee_id": 1,
  "total_price": 50000
}
```

### OrderDetail
```json
{
  "order_detail_id": "550e8400-e29b-41d4-a716-446655440001",
  "order_id": "550e8400-e29b-41d4-a716-446655440000",
  "menu_id": 1,
  "quantity": 2,
  "unit_price": 20000,
  "subtotal": 40000
}
```

---

## Kontribusi

1. Fork repository
2. Buat branch fitur baru (`git checkout -b feature/AmazingFeature`)
3. Commit perubahan (`git commit -m 'Add some AmazingFeature'`)
4. Push ke branch (`git push origin feature/AmazingFeature`)
5. Buat Pull Request

## Lisensi

Distributed under the MIT License. See `LICENSE` for more information.