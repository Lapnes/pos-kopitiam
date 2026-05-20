# POS KopiTiam v3NF — Migration, Views & Infrastructure

## Ringkasan Perubahan Tambahan

Setelah N+1 fix dan transaction boundaries, ini adalah layer infrastruktur database yang **wajib** untuk production.

---

## 1. Database Migration (`migrations/`)

### `001_add_indexes_and_fks.sql`
**Jalankan SETELAH `AutoMigrate` berhasil.**

**Isi:**
- **Foreign Key Constraints** — 20+ FK constraints untuk referential integrity
- **Single Column Indexes** — 40+ indexes untuk foreign keys, filter columns, search columns
- **Composite Indexes** — 4 composite indexes untuk query pattern yang sering dipakai:
  - `idx_orders_dashboard` — `(branch_id, status, created_at DESC)`
  - `idx_details_analytics` — `(order_id, is_voided, subtotal)`
  - `idx_cashmovements_recon` — `(shift_id, type, amount)`
  - `idx_ingredients_profit` — `(recipe_id, raw_material_id, quantity)`

**Kenapa perlu migration terpisah?**
> GORM `AutoMigrate` membuat index dari struct tags, tapi **tidak membuat foreign key constraints** dengan `ON DELETE` rules yang spesifik. Migration ini menjamin FK integrity dan performance index yang lengkap di MySQL.

### `002_create_analytics_views.sql`
**Jalankan SETELAH migration 001.**

**Isi:**
- **8 Standard Views** untuk dashboard analytics yang akan ter-update otomatis secara live seiring perubahan data.

---

## 2. Analytics Views (8 Views)

| View | Fungsi | Sifat |
|------|--------|---------|
| `v_daily_sales_summary` | Grafik dasbor harian manajer | Live (Dynamic) |
| `v_weekly_sales_summary` | Grafik dasbor mingguan manajer | Live (Dynamic) |
| `v_payment_method_summary` | Rekonsiliasi QRIS/Transfer vs uang fisik | Live (Dynamic) |
| `v_menu_profitability_live` | Margin keuntungan real-time per menu | Live (Dynamic) |
| `v_best_selling_items` | Best sellers (endpoint langsung) | Live (Dynamic) |
| `v_shift_reconciliation` | Rapor keandalan kasir | Live (Dynamic) |
| `v_void_and_return_logs` | Pantau anomali void/return | Live (Dynamic) |
| `v_cashier_performance` | Karyawan Terbaik / bonus system | Live (Dynamic) |

### Kenapa menggunakan Standard Views?

Berbeda dengan PostgreSQL yang memiliki fitur *Materialized Views*, arsitektur terbaru ini (MySQL 8.4) menggunakan standar *Dynamic Views* yang langsung mengeksekusi aggregasi data terhadap tabel utama secara *real-time*.

**Benefit:**
- **Selalu Up-to-Date:** Tidak ada *stale data* atau jeda waktu (delay), manajer akan selalu melihat laporan paling aktual detik ini.
- **Mudah Dikelola:** Tidak membutuhkan background scheduler khusus seperti *pg_cron* atau implementasi fungsi refresh terpisah.
- **Performa Andal:** Selama *underlying tables* telah diindeks dengan optimal (sesuai *migration 001*), waktu query tetap akan berada di bawah batas kewajaran untuk operasional UMKM / chain kecil KopiTiam.

---

## 3. Database Config (`internal/config/database.go`)

### Connection Pool

```go
MaxOpenConns:    25        // Max koneksi aktif ke MySQL
MaxIdleConns:    10        // Koneksi idle yang dipertahankan
ConnMaxLifetime: 1 hour    // Recycle koneksi (hindari memory leak)
ConnMaxIdleTime: 10 min    // Tutup koneksi idle yang terlalu lama
```

**Kenapa penting?**
- Tanpa connection pool, setiap request HTTP membuka koneksi baru → MySQL overwhelmed
- `MaxOpenConns` = 25 adalah sweet spot untuk aplikasi POS dengan 10-50 concurrent users
- `ConnMaxLifetime` mencegah koneksi "zombie" yang hang

### Context Propagation

```go
// Handler menerima context dari HTTP request
func (h *Handler) SomeHandler(c *gin.Context) {
    ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
    defer cancel()

    result, err := h.service.SomeMethod(ctx, ...)
    // ...
}
```

---

## 4. Cara Menjalankan Migration

### Development

```bash
# 1. Pastikan MySQL running via Docker (docker-compose up -d mysql)
# 2. Jalankan AutoMigrate (dari main.go atau cmd/migrate)
go run cmd/migrate/main.go

# 3. Jalankan migration SQL via golang-migrate atau eksekusi langsung ke MySQL
mysql -h 127.0.0.1 -u root -p pos_kopitiam < migrations/001_add_indexes_and_fks.sql
mysql -h 127.0.0.1 -u root -p pos_kopitiam < migrations/002_create_analytics_views.sql
```

### Production

```bash
# Gunakan migration tool seperti golang-migrate
# Sangat disarankan untuk memisahkan file .up.sql dan .down.sql agar state bisa di roll-back secara clean.

# Contoh dengan mysql CLI langsung:
mysql -h $DB_HOST -u $DB_USER -p$DB_PASSWORD $DB_NAME < migrations/001_add_indexes_and_fks.sql
mysql -h $DB_HOST -u $DB_USER -p$DB_PASSWORD $DB_NAME < migrations/002_create_analytics_views.sql
```

---

## 5. Environment Variables (Database)

| Variable | Default | Deskripsi |
|----------|---------|-----------|
| `DB_HOST` | localhost | MySQL host (atau path socket UNIX seperti `/tmp/mysql.sock`) |
| `DB_PORT` | 3306 | MySQL port |
| `DB_USER` | root | Username |
| `DB_PASSWORD` | "" | Password |
| `DB_NAME` | pos_kopitiam | Database name |
| `DB_SSLMODE` | disable | SSL mode (Opsional) |
| `DB_MAX_OPEN_CONNS` | 25 | Max open connections |
| `DB_MAX_IDLE_CONNS` | 10 | Max idle connections |
| `DB_CONN_MAX_LIFETIME` | 1h | Connection max lifetime |
| `DB_CONN_MAX_IDLE_TIME` | 10m | Connection max idle time |
| `DB_SLOW_THRESHOLD` | 200ms | Log slow queries |

---

## 6. API Endpoints Baru (Analytics Views)

| Endpoint | Method | View | Deskripsi |
|----------|--------|------|-----------|
| `/api/v1/analytics/daily-sales` | GET | v_daily_sales_summary | Grafik harian |
| `/api/v1/analytics/weekly-sales` | GET | v_weekly_sales_summary | Grafik mingguan |
| `/api/v1/analytics/payment-summary` | GET | v_payment_method_summary | Rekonsiliasi pembayaran |
| `/api/v1/analytics/menu-profitability` | GET | v_menu_profitability_live | Margin menu |
| `/api/v1/analytics/shift-reconciliation` | GET | v_shift_reconciliation | Rapor kasir |
| `/api/v1/analytics/void-return-logs` | GET | v_void_and_return_logs | Log anomali |
| `/api/v1/analytics/cashier-performance` | GET | v_cashier_performance | Performance kasir |

---

## 7. Referensi Integrasi UUID & GORM di MySQL

Aplikasi ini menggunakan integrasi `google/uuid` yang dikelola secara eksklusif oleh *Application Layer* lewat GORM Hook `BeforeCreate`. Hal ini digunakan karena MySQL tidak memiliki tipe data `uuid` native layaknya Postgres. 

```go
func (base *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	if base.ID == uuid.Nil {
		base.ID = uuid.New()
	}
	return
}
```
Field UUID pada GORM akan diterjemahkan secara dinamis menjadi tipe data `varchar(36)` atau `char(36)` oleh driver MySQL.
