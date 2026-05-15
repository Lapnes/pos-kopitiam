# 📊 POS KopiTiam — API Test Recap

**Date:** 2026-05-15 12:21:04

## Summary

| Metric | Value |
|--------|-------|
| Total | 38 |
| ✅ Passed | 34 |
| ❌ Failed | 4 |
| Pass Rate | 89.5% |

## Results

| Phase | Name | Method | Status | Duration | Error |
|-------|------|--------|--------|----------|-------|
| 1.1 | Health Check | GET | ✅ 200 | 3.031132ms | - |
| 1.2 | Login Manager | POST | ✅ 200 | 2.068456193s | - |
| 1.3 | Login Cashier (PIN) | POST | ✅ 200 | 2.309402ms | - |
| 1.4 | Login Kitchen (PIN) | POST | ✅ 200 | 1.826662ms | - |
| 2.1 | Get Active Menus (M2M categories) | GET | ✅ 200 | 8.959664ms | - |
| 3.1 | Create Order (Dine In) | POST | ✅ 201 | 20.19075ms | - |
| 3.2 | Get All Orders | GET | ✅ 200 | 7.779975ms | - |
| 3.3 | Get Order by ID | GET | ✅ 200 | 4.842558ms | - |
| 3.4 | Get Order by Number | GET | ✅ 200 | 4.137682ms | - |
| 3.5 | Confirm Order (stock deduction) | PUT | ✅ 200 | 16.328249ms | - |
| 3.6 | Update Order Notes | PUT | ✅ 200 | 11.143849ms | - |
| 4.1 | Process Payment — Cash (single split) | POST | ✅ 201 | 7.352072ms | - |
| 5.1 | Create Order (Takeaway) | POST | ✅ 201 | 12.414225ms | - |
| 5.2 | Process Payment — Split (Cash + QRIS) | POST | ✅ 201 | 9.39709ms | - |
| 5.3 | Cancel Order (cashier) | PUT | ✅ 200 | 14.669165ms | - |
| 5.4 | Void Order (manager only) | POST | ✅ 200 | 21.947386ms | - |
| 6.1 | Process Return (Cashier) | POST | ✅ 201 | 40.768449ms | - |
| 7.1 | Open Shift | POST | ✅ 201 | 6.641773ms | - |
| 7.2 | Get Current Shift | GET | ✅ 200 | 1.803468ms | - |
| 7.3 | Close Shift | POST | ✅ 200 | 7.452359ms | - |
| 8.1 | Create Raw Material | POST | ✅ 201 | 6.098205ms | - |
| 8.2 | Get Raw Material | GET | ✅ 200 | 1.752067ms | - |
| 8.3 | Update Raw Material | PUT | ✅ 200 | 5.771516ms | - |
| 8.4a | Create Raw Material for Deletion | POST | ✅ 201 | 6.058765ms | - |
| 8.4 | Delete Raw Material | DELETE | ✅ 200 | 4.789026ms | - |
| 9.0 | Create Recipe (BOM) | POST | ❌ 400 | 7.581723ms | Error 1062 (23000): Duplicate entry '330cd968-18f7-40be-ad14... |
| 9.1 | Get Recipe by Menu ID — SKIPPED (recipe creation failed) | SKIP | ❌ 0 |  | skipped: dependency failed |
| 9.2 | Update Recipe — SKIPPED (no recipe_id) | SKIP | ❌ 0 |  | skipped: dependency failed |
| 9.3 | Delete Recipe — SKIPPED (no recipe_id) | SKIP | ❌ 0 |  | skipped: dependency failed |
| 10.1 | Get All Employees | GET | ✅ 200 | 1.697609ms | - |
| 10.2 | Create Employee (cashier) | POST | ✅ 201 | 2.108637171s | - |
| 10.3 | Update Employee | PUT | ✅ 200 | 7.529463ms | - |
| 10.4 | Delete Employee | DELETE | ✅ 200 | 5.8916ms | - |
| 11.1 | Sales Summary (this month) | GET | ✅ 200 | 7.089719ms | - |
| 11.2 | Best Sellers (top 5) | GET | ✅ 200 | 2.373439ms | - |
| 11.3 | Return Impact | GET | ✅ 200 | 4.516201ms | - |
| 12.1 | Generate Receipt (splits-aware) | GET | ✅ 200 | 7.310543ms | - |
| 12.2 | Generate Kitchen Ticket | GET | ✅ 200 | 6.812705ms | - |

## Variables Captured

| Key | Value |
|-----|-------|
| manager_branch_id | 34678d5c-ddee-4529-8cf1-a42cb9ec00b6 |
| menu0_price | 22000 |
| menu2_name | Milk Coffee |
| menu3_name | Chocolate Pudding |
| menu3_price | 15000 |
| order_status | pending |
| takeaway_id | 0e605a45-c6bc-4e7f-a21f-cc8c395cf5ea |
| rm_id | 2bf6b68b-82ea-447e-ba69-29a823cd5a3d |
| cashier_token | eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ... |
| menu0_name | Avocado Juice |
| menu4_id | 66cdc671-e2a1-4840-b7eb-453078297210 |
| emp_id | 30163d39-e119-40fe-b3de-8e21bc4a09ed |
| manager_id | 718e6b22-b6d9-4b39-ab64-f6ae8abd17f0 |
| cashier_id | 45ff9825-7869-47f3-8aee-8a74df681684 |
| menu4_price | 18000 |
| shift_id | 2934bb05-affd-4d7a-b950-4471e17558f4 |
| manager_token | eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ... |
| cashier_branch_id | 34678d5c-ddee-4529-8cf1-a42cb9ec00b6 |
| order_id | 7038b324-922c-4678-99e7-726fb68bb0ac |
| confirmed_status | confirmed |
| rm_name | Arabica Coffee Beans |
| kitchen_token | eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ... |
| menu1_name | Plain Tea |
| menu1_price | 8000 |
| order_total | 71920.00 |
| menu0_id | 330cd968-18f7-40be-ad14-549c3a856cc7 |
| menu3_id | 4e92db66-631b-45e9-b5bf-4b3753af6053 |
| menu1_id | 331f3bb9-245a-4182-afb5-2699604d883b |
| order_number | ORD-20260515-0314 |
| takeaway_total | 4280.00 |
| shift_status | open |
| shift_variance | 25000.00 |
| shift_expected | 500000.00 |
| shift_sales | 0.00 |
| menu2_id | 4df7184c-59f6-4acc-9e93-c4702cb22f8a |
| menu2_price | 18000 |
| menu4_name | French Fries |
| rm_delete_id | 1d0b06c2-9762-4168-8285-472e7b21060d |
