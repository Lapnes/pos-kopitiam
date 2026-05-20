#!/bin/bash

echo "═══════════════════════════════════════════════════"
echo "   BENCHMARK: Email Login vs PIN Login"
echo "═══════════════════════════════════════════════════"
echo ""

# 1. Email Login (Manager) - seharusnya ~2.5s karena 1 bcrypt
echo "[1] Email Login (bcrypt cost=14, 1 hash)"
time curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"manager@kopitiam.id","password":"password123"}' \
  -o /dev/null -w "\nHTTP %{http_code} | Time: %{time_total}s | DNS: %{time_namelookup}s | Connect: %{time_connect}s | TTFB: %{time_starttransfer}s\n"

echo ""

# 2. PIN Login Cashier - seharusnya ~6s karena scan + bcrypt
echo "[2] PIN Login Cashier (bcrypt cost=14, scan 3 users)"
time curl -s -X POST http://localhost:8080/api/v1/auth/login-pin \
  -H "Content-Type: application/json" \
  -d '{"pin":"222222"}' \
  -o /dev/null -w "\nHTTP %{http_code} | Time: %{time_total}s | DNS: %{time_namelookup}s | Connect: %{time_connect}s | TTFB: %{time_starttransfer}s\n"

echo ""

# 3. PIN Login Kitchen - user ke-3, harus fail 2x dulu
echo "[3] PIN Login Kitchen (bcrypt cost=14, fail 2 + match 1)"
time curl -s -X POST http://localhost:8080/api/v1/auth/login-pin \
  -H "Content-Type: application/json" \
  -d '{"pin":"333333"}' \
  -o /dev/null -w "\nHTTP %{http_code} | Time: %{time_total}s | DNS: %{time_namelookup}s | Connect: %{time_connect}s | TTFB: %{time_starttransfer}s\n"

echo ""

# 4. Create Employee - 2 bcrypt hash (password + PIN)
echo "[4] Create Employee (bcrypt cost=14 × 2 = password + PIN)"
MANAGER_TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"manager@kopitiam.id","password":"password123"}' | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)

time curl -s -X POST http://localhost:8080/api/v1/employees \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $MANAGER_TOKEN" \
  -d '{"branch_id":"ef4ca53f-6dc4-45c7-8a52-f972dc9b3293","name":"Benchmark Cashier","email":"bench@kopitiam.id","password":"password123","pin_code":"444444","role":"cashier","is_active":true}' \
  -o /dev/null -w "\nHTTP %{http_code} | Time: %{time_total}s | DNS: %{time_namelookup}s | Connect: %{time_connect}s | TTFB: %{time_starttransfer}s\n"

echo ""
echo "═══════════════════════════════════════════════════"
echo "   ANALISIS:"
echo "═══════════════════════════════════════════════════"
echo "Email login  : ~2.5s (1 bcrypt hash)"
echo "PIN login    : ~4-7s (N bcrypt compares, N = user count)"
echo "Create emp   : ~4-5s (2 bcrypt hashes: password + PIN)"
echo ""
echo "Bcrypt cost=14 = ~250ms per hash/compare"
echo "PIN login scan semua user = N × 250ms"
echo ""
echo "SOLUSI: Turunkan bcrypt cost PIN ke 10 (~50ms)"
echo "        atau gunakan Redis cache untuk PIN"
