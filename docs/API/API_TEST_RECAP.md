# 📊 POS KopiTiam — API Test Recap

**Date:** 2026-05-20 21:56:18

## Summary

| Metric | Value |
|--------|-------|
| Total | 43 |
| ✅ Passed | 43 |
| ❌ Failed | 0 |
| Pass Rate | 100.0% |

## Results

| Phase | Name | Method | Status | Duration | Error |
|-------|------|--------|--------|----------|-------|
| 1.1 | Health Check | GET | ✅ 200 | 2.576637ms | - |
| 1.2 | Login Manager | POST | ✅ 200 | 2.059148466s | - |
| 1.3 | Login Cashier (PIN) | POST | ✅ 200 | 131.047205ms | - |
| 1.4 | Login Kitchen (PIN) | POST | ✅ 200 | 255.781613ms | - |
| 2.1 | Get Active Menus (M2M categories) | GET | ✅ 200 | 7.750986ms | - |
| 3.1 | Create Order (Dine In) | POST | ✅ 201 | 20.534423ms | - |
| 3.2 | Get All Orders | GET | ✅ 200 | 41.433624ms | - |
| 3.3 | Get Order by ID | GET | ✅ 200 | 6.645751ms | - |
| 3.4 | Get Order by Number | GET | ✅ 200 | 7.4543ms | - |
| 3.5 | Confirm Order (stock deduction) | PATCH | ✅ 200 | 19.694204ms | - |
| 3.6 | Update Order Notes | PATCH | ✅ 200 | 10.891863ms | - |
| 4.1 | Process Payment — Cash (single split) | POST | ✅ 200 | 41.83078ms | - |
| 5.1 | Create Order (Takeaway) | POST | ✅ 201 | 12.818016ms | - |
| 5.2 | Process Payment — Split (Cash + QRIS) | POST | ✅ 200 | 810.714983ms | - |
| 5.3a | Create Order for Cancel | POST | ✅ 201 | 14.852221ms | - |
| 5.3 | Cancel Order (cashier) | PATCH | ✅ 200 | 11.431363ms | - |
| 5.4a | Create Order for Void | POST | ✅ 201 | 12.342503ms | - |
| 5.4 | Void Order (manager only) | PATCH | ✅ 200 | 14.362544ms | - |
| 6.1 | Process Return (Cashier) | POST | ✅ 200 | 22.078314ms | - |
| 7.1 | Open Shift | POST | ✅ 201 | 8.506919ms | - |
| 7.2 | Get Current Shift | GET | ✅ 200 | 2.128724ms | - |
| 7.3 | Close Shift | PATCH | ✅ 200 | 10.048961ms | - |
| 8.1 | Create Raw Material | POST | ✅ 201 | 7.533643ms | - |
| 8.2 | Get Raw Material | GET | ✅ 200 | 2.625446ms | - |
| 8.3 | Update Raw Material | PUT | ✅ 200 | 7.783932ms | - |
| 8.4a | Create Raw Material for Deletion | POST | ✅ 201 | 7.65516ms | - |
| 8.4 | Delete Raw Material | DELETE | ✅ 200 | 6.460315ms | - |
| 9.0 | Create Recipe (BOM) | POST | ✅ 200 | 13.703735ms | - |
| 9.1 | Get Recipe by Menu ID | GET | ✅ 200 | 3.760623ms | - |
| 9.2 | Update Recipe | PUT | ✅ 200 | 28.488052ms | - |
| 9.3 | Delete Recipe | DELETE | ✅ 200 | 8.782867ms | - |
| 10.1 | Get All Employees | GET | ✅ 200 | 2.263493ms | - |
| 10.2 | Create Employee (cashier) | POST | ✅ 201 | 2.152001732s | - |
| 10.3 | Update Employee | PUT | ✅ 200 | 6.616646ms | - |
| 10.4 | Delete Employee | DELETE | ✅ 200 | 5.107686ms | - |
| 11.1 | Sales Summary (this month) | GET | ✅ 200 | 30.9153ms | - |
| 11.2 | Best Sellers (top 5) | GET | ✅ 200 | 16.721982ms | - |
| 11.3 | Return Impact | GET | ✅ 200 | 7.521129ms | - |
| 11.4 | COGS | GET | ✅ 200 | 17.429682ms | - |
| 12.1 | Generate Receipt (splits-aware) | GET | ✅ 200 | 5.93523ms | - |
| 12.2 | Generate Kitchen Ticket | GET | ✅ 200 | 7.160991ms | - |
| 13.1 | Get Catalog Menus (Basic Auth) | GET | ✅ 200 | 22.041378ms | - |
| 13.2 | Catalog Self-Checkout (QRIS) | POST | ✅ 200 | 136.330627ms | - |

## Variables Captured

| Key | Value |
|-----|-------|
| manager_token | eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ... |
| cashier_branch_id | e2c8c166-bbf5-4290-a630-044c563bfc03 |
| menu1_name | Ice Lemon Tea |
| menu3_price | 20000 |
| order_id | af653a93-c844-4333-a7cd-9fd960521d51 |
| order_status | pending |
| takeaway_id | eb9321d8-7c3a-4ef8-94f9-595d6e43ef7d |
| takeaway_total | 11240.00 |
| cashier_token | eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ... |
| menu0_price | 12000 |
| menu3_name | Chendol |
| menu4_id | 71121dba-3151-487b-b371-9634f68908b7 |
| order_number | ORD-20260520-0460 |
| shift_id | dbad94fa-80bd-401b-be42-8fc6b45dc32b |
| emp_id | b58f8765-722d-450e-ad8e-d723a0b62ab7 |
| kitchen_token | eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ... |
| rm_name | Arabica Coffee Beans |
| rm_delete_id | cee9cb0b-f20a-402e-85d1-752866283125 |
| menu4_name | Teh Tarik |
| rm_id | 19cd3252-bf60-4dd1-922c-a0057fc544ee |
| catalog_menu_id | 02ff2914-8d43-41f1-b3b8-eb5ce36d43ad |
| confirmed_status | confirmed |
| manager_id | dc92be2c-f2f3-46cc-9688-fc3eaaf46ff8 |
| menu1_price | 14000 |
| menu2_id | 44dc06db-1be7-4d2c-bbd9-504af8db280e |
| menu2_price | 22000 |
| menu3_id | 57e9f605-bac2-4408-9f22-d5c8bedb91b5 |
| menu4_price | 16000 |
| order_total | 46400.00 |
| manager_branch_id | e2c8c166-bbf5-4290-a630-044c563bfc03 |
| cashier_id | 41bfc869-52f0-4a53-9bb7-b742f0638d3d |
| menu1_id | 24fb54cb-a262-4951-bcec-42abed063d93 |
| shift_status | open |
| recipe_id | 87a8e4c4-9f98-40e1-9d31-82b14041cb4c |
| menu0_id | 02ff2914-8d43-41f1-b3b8-eb5ce36d43ad |
| menu0_name | Kopi O |
| menu2_name | Matcha Latte |
