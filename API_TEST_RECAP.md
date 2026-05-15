# 📊 POS KopiTiam — API Test Recap

**Date:** 2026-05-15 15:52:50

## Summary

| Metric | Value |
|--------|-------|
| Total | 42 |
| ✅ Passed | 0 |
| ❌ Failed | 42 |
| Pass Rate | 0.0% |

## Results

| Phase | Name | Method | Status | Duration | Error |
|-------|------|--------|--------|----------|-------|
| 1.1 | Health Check | GET | ❌ 0 | 1.838161ms | Get "http://localhost:8080/health": dial tcp [::1]:8080: con... |
| 1.2 | Login Manager | POST | ❌ 0 | 911.711µs | Post "http://localhost:8080/api/v1/auth/login": dial tcp [::... |
| 1.3 | Login Cashier (PIN) | POST | ❌ 0 | 1.159876ms | Post "http://localhost:8080/api/v1/auth/login-pin": dial tcp... |
| 1.4 | Login Kitchen (PIN) | POST | ❌ 0 | 916.334µs | Post "http://localhost:8080/api/v1/auth/login-pin": dial tcp... |
| 2.1 | Get Active Menus (M2M categories) | GET | ❌ 0 | 1.086872ms | Get "http://localhost:8080/api/v1/menus": dial tcp [::1]:808... |
| 3.1 | Create Order (Dine In) | POST | ❌ 0 | 933.894µs | Post "http://localhost:8080/api/v1/orders": dial tcp [::1]:8... |
| 3.2 | Get All Orders — SKIPPED (order creation failed) | SKIP | ❌ 0 |  | skipped: dependency failed |
| 3.3 | Get Order by ID — SKIPPED | SKIP | ❌ 0 |  | skipped: dependency failed |
| 3.4 | Get Order by Number — SKIPPED | SKIP | ❌ 0 |  | skipped: dependency failed |
| 3.5 | Confirm Order — SKIPPED | SKIP | ❌ 0 |  | skipped: dependency failed |
| 3.6 | Update Order Notes — SKIPPED | SKIP | ❌ 0 |  | skipped: dependency failed |
| 4.1 | Process Payment (Cash) — SKIPPED | SKIP | ❌ 0 |  | skipped: dependency failed |
| 6.1 | Process Return — SKIPPED | SKIP | ❌ 0 |  | skipped: dependency failed |
| 12.1 | Generate Receipt — SKIPPED | SKIP | ❌ 0 |  | skipped: dependency failed |
| 12.2 | Generate Kitchen Ticket — SKIPPED | SKIP | ❌ 0 |  | skipped: dependency failed |
| 4.1 | Process Payment (Cash) — SKIPPED (no order_id) | SKIP | ❌ 0 |  | skipped: dependency failed |
| 5.1 | Create Order (Takeaway) | POST | ❌ 0 | 1.180903ms | Post "http://localhost:8080/api/v1/orders": dial tcp [::1]:8... |
| 5.2 | Process Payment — Split (Cash + QRIS) — SKIPPED (takeaway creation failed) | SKIP | ❌ 0 |  | skipped: dependency failed |
| 5.3 | Cancel Order (cashier) — SKIPPED (no takeaway_id) | SKIP | ❌ 0 |  | skipped: dependency failed |
| 5.4 | Void Order — SKIPPED (no order_id) | SKIP | ❌ 0 |  | skipped: dependency failed |
| 6.1 | Process Return — SKIPPED (no order_id) | SKIP | ❌ 0 |  | skipped: dependency failed |
| 7.1 | Open Shift | POST | ❌ 0 | 759.647µs | Post "http://localhost:8080/api/v1/shifts/open": dial tcp [:... |
| 7.2 | Get Current Shift | GET | ❌ 0 | 512.584µs | Get "http://localhost:8080/api/v1/shifts/current": dial tcp ... |
| 7.3 | Close Shift | POST | ❌ 0 | 548.224µs | Post "http://localhost:8080/api/v1/shifts/close": dial tcp [... |
| 8.1 | Create Raw Material | POST | ❌ 0 | 2.522297ms | Post "http://localhost:8080/api/v1/master/raw-materials": di... |
| 8.2 | Get Raw Material | GET | ❌ 0 | 2.276154ms | Get "http://localhost:8080/api/v1/master/raw-materials/": di... |
| 8.3 | Update Raw Material | PUT | ❌ 0 | 1.677909ms | Put "http://localhost:8080/api/v1/master/raw-materials/": di... |
| 8.4a | Create Raw Material for Deletion | POST | ❌ 0 | 1.116385ms | Post "http://localhost:8080/api/v1/master/raw-materials": di... |
| 8.4 | Delete Raw Material | DELETE | ❌ 0 | 926.882µs | Delete "http://localhost:8080/api/v1/master/raw-materials/":... |
| 9.0 | Create Recipe (BOM) | POST | ❌ 0 | 908.074µs | Post "http://localhost:8080/api/v1/master/recipes": dial tcp... |
| 9.1 | Get Recipe by Menu ID — SKIPPED (recipe creation failed) | SKIP | ❌ 0 |  | skipped: dependency failed |
| 9.2 | Update Recipe — SKIPPED (no recipe_id) | SKIP | ❌ 0 |  | skipped: dependency failed |
| 9.3 | Delete Recipe — SKIPPED (no recipe_id) | SKIP | ❌ 0 |  | skipped: dependency failed |
| 10.1 | Get All Employees | GET | ❌ 0 | 924.052µs | Get "http://localhost:8080/api/v1/employees": dial tcp [::1]... |
| 10.2 | Create Employee (cashier) | POST | ❌ 0 | 1.060703ms | Post "http://localhost:8080/api/v1/employees": dial tcp [::1... |
| 10.3 | Update Employee | PUT | ❌ 0 | 1.109431ms | Put "http://localhost:8080/api/v1/employees/": dial tcp [::1... |
| 10.4 | Delete Employee | DELETE | ❌ 0 | 884.473µs | Delete "http://localhost:8080/api/v1/employees/": dial tcp [... |
| 11.1 | Sales Summary (this month) | GET | ❌ 0 | 653.672µs | Get "http://localhost:8080/api/v1/analytics/sales-summary?st... |
| 11.2 | Best Sellers (top 5) | GET | ❌ 0 | 555.971µs | Get "http://localhost:8080/api/v1/analytics/best-sellers?lim... |
| 11.3 | Return Impact | GET | ❌ 0 | 437.397µs | Get "http://localhost:8080/api/v1/analytics/return-impact": ... |
| 12.1 | Generate Receipt — SKIPPED (no order_id) | SKIP | ❌ 0 |  | skipped: dependency failed |
| 12.2 | Generate Kitchen Ticket — SKIPPED (no order_id) | SKIP | ❌ 0 |  | skipped: dependency failed |

## Variables Captured

| Key | Value |
|-----|-------|
