# POS KopiTiam Project

POS KopiTiam adalah sistem Point of Sale (POS) profesional untuk kafe KopiTiam, dibangun menggunakan bahasa pemrograman Go dengan arsitektur Clean Architecture. Sistem ini menyediakan API RESTful terstruktur untuk mengelola pesanan, data master, shift kasir, dan laporan penjualan.

## Fitur & Struktur Backend

Aplikasi ini mengadopsi prinsip **Clean Architecture**, memisahkan secara tegas layer `Handler`, `Service`, `Repository`, dan `Model`. Berikut adalah pemetaan fungsionalitas dan aliran sistem backend yang telah diimplementasikan:

### 1. Authentication & Authorization
Fitur otentikasi diurus melalui middleware khusus (`internal/middleware/auth.go`) menggunakan mekanisme **JSON Web Tokens (JWT)** dan **Role-Based Access Control (RBAC)**.
- **Login**: `POST /api/v1/auth/login` (Menerima `email` & `password`, divalidasi via Bcrypt. Mengembalikan pasangan `access_token` & `refresh_token`).
- **Login PIN**: `POST /api/v1/auth/login-pin` (Metode otentikasi cepat untuk mesin kasir menggunakan 6-digit PIN unik).
- **Refresh Token**: `POST /api/v1/auth/refresh` (Memperbarui sesi token akses yang telah kedaluwarsa).
- **Logout**: `POST /api/v1/auth/logout` (Menutup sesi otentikasi dan *blacklist* token).

### 2. Order Management
Alur utama pemrosesan transaksi yang dienkapsulasi pada layer *business logic* (`internal/service/order_service.go`).
- **Create Order**: `POST /api/v1/orders` (Validasi keranjang pesanan, menghitung total bayar, PPN otomatis 11%, service charge 5%, serta memverifikasi ketersediaan stok).
- **Read / Get Orders**:
  - `GET /api/v1/orders` (Mengambil seluruh daftar pesanan dengan dukungan Pagination dan Filter kompleks).
  - `GET /api/v1/orders/:id` (Mengambil detail penuh dari satu pesanan, termasuk daftar struk item pesanan).
- **Update Order**: `PUT /api/v1/orders/:id` (Memperbarui isi pesanan hanya jika status belum tereksekusi ke tahap selanjutnya).
- **Confirm Order**: `PUT /api/v1/orders/:id/confirm` (Validasi final pesanan masuk dan meneruskannya ke tim Kitchen / Bar).
- **Cancel Order**: `PUT /api/v1/orders/:id/cancel` (Membatalkan seluruh pesanan yang belum diproses dan mengembalikan perhitungan *daily stock*).

### 3. Void Operations
Manajemen risiko dan koreksi struk pada `internal/handler/order_handler.go`, dilindungi sangat ketat oleh middleware RBAC (Hanya peran `Manager` dan `Superadmin` yang memiliki izin eksekusi).
- **Void Specific Items**: `POST /api/v1/orders/:id/void-item` (Membatalkan item tunggal dalam pesanan yang sudah berjalan, melacak pengguna yang membatalkan dan log alasan).
- **Void Whole Order**: `POST /api/v1/orders/:id/void` (Membatalkan seluruh transaksi yang sudah terkonfirmasi secara aman).

### 4. Returns & Refunds System
Sistem retur barang dan pengembalian uang (*Refund*) bagi pelanggan:
- **Max 80% Rule**: Batas maksimal pengembalian uang adalah 80% dari total pesanan. Jika permintaan *refund* melebihi batas ini, sistem akan menolak secara otomatis.
- **Data Audit Trail**: Sama seperti Void, *Refund* tidak melakukan *Hard Delete* pada transaksi asli. Transaksi yang diretur akan ditandai (`HasReturns = true`), memastikan catatan penjualan tetap aman untuk keperluan rekonsiliasi.

### 5. Master Data Management
Mengatur integritas entitas krusial perusahaan yang didefinisikan pada `internal/models/domain.go`. Seluruh *primary key* dikelola dalam format **UUID v4**.
- **Categories**: Pengelompokan hirarkis dari tipe menu (Makanan, Minuman, Pastry).
- **Menus**: Komponen penjualan utama berisi Harga Jual, Harga Pokok Penjualan (Cost Price/HPP), validasi *Daily Stock*, dan instruksi routing ke Stasiun Dapur.
- **Employees**: Manajemen pengguna beserta penetapan Peran / Jabatan sistem.
- **Tables**: Manajemen ketersediaan meja terpusat yang memantau status langsung (Available / Occupied).
- **Areas**: Manajemen zona ruang (Misal: Lantai 1, Smoking Area) per *Branch* / Cabang kafe.
- **Inventory & BOM (Bill of Materials)**: Melacak `RawMaterial` secara *real-time*. Terintegrasi dengan `Recipe` pada setiap menu sehingga saat transaksi terkonfirmasi, bahan baku otomatis berkurang dengan akurat.

### 6. Advanced Analytics
Sistem agregasi data untuk pelaporan manajemen tingkat tinggi:
- **Sales Summary**: Melacak Gross Revenue, Net Revenue (Gross dikurangi Refunds), dan Total Refunded dalam jangka waktu tertentu.
- **Best Sellers**: Menganalisis pergerakan produk dengan penjualan tertinggi, mengabaikan data dari transaksi yang telah di-Void.
- **Return Impact**: Merekap dan memetakan alasan retur terbanyak untuk mengidentifikasi kebocoran kualitas atau layanan.

### 7. Shift & Cash Drawer Management
Modul terintegrasi bagi kasir untuk memvalidasi aliran uang tunai (*petty cash* dan *sales*) selama satu giliran kerja.
- **Open Shift**: `POST /api/v1/shifts/open` (Memasukkan modal awal kasir).
- **Get Current Shift**: `GET /api/v1/shifts/current` (Melacak performa shift aktif dan akumulasi transaksi).
- **Close Shift**: `POST /api/v1/shifts/close` (Menginput sisa saldo aktual kasir dan sistem otomatis mengkalkulasi selisih/varian terhadap transaksi sistem).

### 8. Payload Printers untuk ESC/POS
Engine pemformatan payload termal yang dioptimalkan untuk perangkat keras Point-of-Sale:
- **Receipts**: Menghasilkan JSON komprehensif berisi *Store Header*, *Cashier*, rincian produk beserta harga, perpajakan, *Service Charge*, hingga detail metode pembayaran.
- **Kitchen Tickets**: Versi ringkas tanpa harga, ditujukan ke stasiun dapur/bar, menampilkan tipe pesanan, meja pesanan, kuantitas item, serta instruksi khusus pelanggan.

### 9. Database & Infrastructure
Sistem terdistribusi dan reliabel yang dipisahkan melalui folder skrip `cmd/`.
- **Migrations (`cmd/migrate/main.go`)**: Skrip terpisah berbasis *GORM AutoMigrate* yang bertugas meruntuhkan dan membangun ulang relasi basis data (Foreign Keys) untuk mencegah konflik tipe skema *legacy*.
- **Seeding (`seeds/seed.go`)**: Pre-populasi tabel basis data dengan *Master Data* struktural (Akun Superadmin, Cabang, Area) agar lingkungan siap pakai untuk pengembangan.
- **Health Checks**: `GET /health` (Menyediakan parameter kesiapan kontainer API, menguji koneksi *ping* ke MySQL, dan ketersediaan modul antrean di Redis).

## Teknologi yang Digunakan

- **Bahasa Pemrograman**: Go 1.23+
- **Web Framework**: Gin (Dipilih karena performa HTTP router yang sangat cepat)
- **ORM**: GORM
- **Identifier**: UUID v4 (Untuk *generate* ID unik pada master data dan pesanan untuk mencegah kebocoran *sequential ID*)

**Infrastruktur & Konfigurasi Khusus:**
- **Database**: 
  - **MySQL 8.0**: Bertindak sebagai *Primary Relational Database* untuk menyimpan seluruh data master (Menu, Karyawan, Cabang) dan transaksi persisten (Pesanan, Log Pembayaran) secara konsisten menggunakan ACID properties.
  - **Redis 7**: Bertindak sebagai *In-Memory Data Store* berkecepatan sangat tinggi. Pada arsitektur ini, Redis dikhususkan sebagai *Message Broker* untuk menampung data antrean sebelum diproses oleh *worker*.
- **Background Job (Asynq)**: *Library* tangguh dari ekosistem Go yang bertugas memproses antrean tugas di latar belakang (*background task processing*). Asynq memantau Redis secara konstan. Di proyek KopiTiam, ini digunakan khusus untuk menyinkronkan data pesanan (Jurnal Keuangan) ke *Bigcapital API* secara asinkron agar *response* struk kasir tetap instan tanpa harus menunggu balasan dari pihak ke-3.
- **Environment Config (`.env`)**: Konfigurasi tingkat sistem operasi yang di-*load* via `godotenv`. Hal ini memusatkan seluruh *credential* rahasia (seperti Password MySQL, Host Redis, Secret Key JWT, dan API Key Bigcapital) sehingga *source code* aman untuk diunggah ke publik (GitHub) dan sistem menjadi sangat *portable* ketika dipindahkan antar *Environment* (Dev, Staging, maupun Prod).

## Struktur Proyek

```text
pos_kopitiam_db_baru/
├── cmd/
│   ├── api/
│   │   └── main.go          # Entry point API Server (Gin)
│   ├── migrate/
│   │   └── main.go          # Script auto-migrasi database GORM
│   └── worker/
│       └── main.go          # Entry point Asynq Background Worker
├── internal/
│   ├── config/              # Inisialisasi DB, Redis, dan Environment
│   ├── dto/                 # Data Transfer Objects (Request & Response Payload)
│   ├── handler/             # HTTP Handlers / Controllers
│   ├── middleware/          # Gin Middlewares (Auth, RBAC Checker)
│   ├── models/              # Struktur Domain dan Enum (Struktur Tabel)
│   ├── repository/          # Data Access Layer (Komunikasi Database)
│   ├── routes/              # Konfigurasi Routing API
│   ├── service/             # Business Logic Layer
│   └── utils/               # Hash Password, JWT Helpers, Format Response
├── pkg/
│   └── bigcapital/          # Integrasi Eksternal ke API Bigcapital
├── seeds/
│   └── seed.go              # Script pengisian data awal
├── docker-compose.yml       # Konfigurasi infrastruktur (MySQL, Redis)
├── go.mod                   # Go module file
└── README.md                # Dokumentasi ini
```

## Setup dan Instalasi

### Prasyarat
- Go 1.23 atau lebih baru
- Docker & Docker Compose
- Git

### Langkah Instalasi

1. **Clone Repository**
   ```bash
   git clone https://github.com/Lapnes/pos-kopitiam.git
   cd pos-kopitiam
   ```

2. **Install Dependencies**
   ```bash
   go mod tidy
   ```

3. **Jalankan Database via Docker**
   Jalankan container MySQL dan Redis di latar belakang secara terisolasi:
   ```bash
   docker compose up -d --build
   ```

4. **Konfigurasi Environment**
   Buat file `.env` di root directory. Anda bisa menyalin file contoh yang telah disediakan:
   ```bash
   cp .env.example .env
   ```
   *Secara default, nilainya sudah cocok dengan konfigurasi Docker Compose.*

5. **Jalankan Migrasi Database**
   Pastikan MySQL sudah siap (kurang lebih butuh beberapa detik). Setelah siap, buat skema tabelnya:
   ```bash
   go run cmd/migrate/main.go
   ```

6. **Jalankan Seeder (Opsional)**
   Populasikan database dengan akun karyawan (Admin, Kasir, dsb.), area, meja, kategori, dan menu:
   ```bash
   go run seeds/seed.go
   ```

7. **Jalankan Aplikasi**
   Buka dua jendela terminal untuk menjalankan subsistem secara bersamaan.

   **Terminal 1 (Menjalankan API Server):**
   ```bash
   go run cmd/api/main.go
   ```
   *API akan berjalan di http://localhost:8080*

   **Terminal 2 (Menjalankan Background Worker):**
   ```bash
   go run cmd/worker/main.go
   ```

## Dokumentasi API

Base URL: `http://localhost:8080`

### Health Check
- `GET /health`
Memeriksa status kesehatan API, Redis, dan MySQL.
```json
{
  "status": "up",
  "redis": true,
  "db": true
}
```

### Authentication
Akses Endpoint untuk autentikasi sistem KopiTiam:

#### 1. Login dengan Email & Password
- **Endpoint**: `POST /api/v1/auth/login`
- **Request Body**:
```json
{
  "email": "admin@kopitiam.com",
  "password": "password123"
}
```
- **Response**:
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "70c6ade5-a4a2-476c-9522-2e6f78e06105",
      "role": "superadmin",
      "branch_id": "f763ab0b-192f-4e69-a77d-7f87ad4e3448"
    }
  }
}
```

#### 2. Login Cepat dengan PIN
- **Endpoint**: `POST /api/v1/auth/login-pin`
- **Request Body**:
```json
{
  "pin": "111111"
}
```

#### 3. Refresh Token & Logout
- `POST /api/v1/auth/refresh`
- `POST /api/v1/auth/logout`

### Orders
*Catatan: Endpoint di bawah ini membutuhkan Bearer Token.*

#### 1. Membuat Pesanan Baru
- **Endpoint**: `POST /api/v1/orders`
- **Request Body**:
```json
{
  "table_id": "9bd36e63-0304-446c-adaf-940d50dece96",
  "customer_name": "Budi",
  "order_type": "dine_in",
  "notes": "Tolong cepat",
  "items": [
    {
      "menu_id": "d05fab71-7f41-4023-8938-041572e7b5a6",
      "menu_name": "Nasi Goreng Spesial",
      "quantity": 2,
      "price": 35000.00,
      "notes": "Tidak pakai pedas"
    }
  ]
}
```
- **Response**:
```json
{
  "success": true,
  "message": "Order created successfully",
  "data": {
    "id": "c12a8400-e29b-41d4-a716-446655440000",
    "order_number": "ORD-20260508-1234",
    "status": "pending",
    "subtotal": 70000.00,
    "tax_amount": 7700.00,
    "service_charge": 3500.00,
    "total": 81200.00
  }
}
```

#### 2. Konfirmasi / Batal / Update
- `GET /api/v1/orders` - Mengambil daftar semua pesanan.
- `GET /api/v1/orders/:id` - Mengambil detail informasi order secara spesifik.
- `PUT /api/v1/orders/:id` - Memperbarui (update) isi pesanan (Hanya jika status masih Pending/Confirmed).
- `PUT /api/v1/orders/:id/confirm` - Memverifikasi pesanan masuk.
- `PUT /api/v1/orders/:id/cancel` - Membatalkan pesanan masuk.
- `POST /api/v1/orders/:id/void-item` - Menghapus item menu tertentu (Akses: Manager/Superadmin).
- `POST /api/v1/orders/:id/void` - Membatalkan seluruh pesanan yang sudah berjalan (Akses: Manager/Superadmin).

### Shift Management
*Catatan: Endpoint di bawah ini membutuhkan Bearer Token (Peran: Cashier, Manager, Superadmin).*
- `POST /api/v1/shifts/open` - Membuka shift kerja baru.
- `GET /api/v1/shifts/current` - Mendapatkan informasi shift yang sedang berjalan.
- `POST /api/v1/shifts/close` - Menutup shift dan menghitung kecocokan kas (Variance).

### Returns & Refunds
*Catatan: Endpoint di bawah ini membutuhkan Bearer Token (Peran: Manager, Superadmin).*
- `POST /api/v1/orders/:id/returns` - Memproses pengembalian pesanan (maksimal 80% dari total transaksi).

### Master Data Management
*Catatan: Endpoint di bawah ini membutuhkan Bearer Token (Peran: Admin, Superadmin, Manager).*
- `POST, PUT, GET, DELETE /api/v1/master/raw-materials` - Manajemen bahan baku mentah.
- `POST, PUT, GET, DELETE /api/v1/master/recipes` - Manajemen resep (BOM) untuk setiap produk.

### Advanced Analytics
*Catatan: Endpoint di bawah ini membutuhkan Bearer Token (Peran: Manager, Superadmin).*
- `GET /api/v1/analytics/sales-summary` - Mengambil laporan penjualan ringkas.
- `GET /api/v1/analytics/best-sellers` - Melihat produk paling laku (Best Sellers).
- `GET /api/v1/analytics/return-impact` - Melihat dampak retur pada pendapatan.

### Payload Printers
*Catatan: Endpoint di bawah ini membutuhkan Bearer Token.*
- `GET /api/v1/printers/receipt/:id` - Generate JSON payload untuk struk pelanggan.
- `GET /api/v1/printers/kitchen/:id` - Generate JSON payload untuk struk dapur.

## Model Data

Semua Entity Database pada arsitektur terbaru diformat menggunakan `UUID`.

### Category
```json
{
  "id": "1ca393ff-74f7-4f44-90cf-7194fca0c15b",
  "name": "Makanan",
  "station": "kitchen",
  "sort_order": 0
}
```

### Menu
```json
{
  "id": "d05fab71-7f41-4023-8938-041572e7b5a6",
  "category_id": "1ca393ff-74f7-4f44-90cf-7194fca0c15b",
  "name": "Nasi Goreng Spesial",
  "price": 35000.00,
  "cost_price": 20000.00,
  "daily_stock": 50,
  "station": "kitchen",
  "is_active": true
}
```

### Employee
```json
{
  "id": "70c6ade5-a4a2-476c-9522-2e6f78e06105",
  "branch_id": "f763ab0b-192f-4e69-a77d-7f87ad4e3448",
  "name": "Super Admin",
  "email": "admin@kopitiam.com",
  "role": "superadmin",
  "is_active": true
}
```

### Order
```json
{
  "id": "c12a8400-e29b-41d4-a716-446655440000",
  "order_number": "ORD-20260508-1234",
  "employee_id": "70c6ade5-a4a2-476c-9522-2e6f78e06105",
  "order_type": "dine_in",
  "status": "pending",
  "subtotal": 35000.00,
  "tax_amount": 3850.00,
  "service_charge": 1750.00,
  "total": 40600.00
}
```

### OrderDetail
```json
{
  "id": "e44c8400-e29b-41d4-a716-446655440001",
  "order_id": "c12a8400-e29b-41d4-a716-446655440000",
  "menu_id": "d05fab71-7f41-4023-8938-041572e7b5a6",
  "menu_name": "Nasi Goreng Spesial",
  "quantity": 1,
  "price": 35000.00,
  "subtotal": 35000.00,
  "station": "kitchen"
}
```

## Kontribusi

1. Fork repository
2. Buat branch fitur baru (`git checkout -b feature/AmazingFeature`)
3. Commit perubahan (`git commit -m 'Add some AmazingFeature'`)
4. Push ke branch (`git push origin feature/AmazingFeature`)
5. Buat Pull Request

## Lisensi

Distributed under the MIT License. See `LICENSE` for more information.
