# 📊 POS KopiTiam — API Test Recap

**Date:** 2026-05-20 21:37:01

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
| 1.1 | Health Check | GET | ✅ 200 | 6.086629ms | - |
| 1.2 | Login Manager | POST | ✅ 200 | 2.741708582s | - |
| 1.3 | Login Cashier (PIN) | POST | ✅ 200 | 385.802366ms | - |
| 1.4 | Login Kitchen (PIN) | POST | ✅ 200 | 213.615316ms | - |
| 2.1 | Get Active Menus (M2M categories) | GET | ✅ 200 | 9.409443ms | - |
| 3.1 | Create Order (Dine In) | POST | ✅ 201 | 19.572739ms | - |
| 3.2 | Get All Orders | GET | ✅ 200 | 43.087293ms | - |
| 3.3 | Get Order by ID | GET | ✅ 200 | 10.229135ms | - |
| 3.4 | Get Order by Number | GET | ✅ 200 | 7.523524ms | - |
| 3.5 | Confirm Order (stock deduction) | PATCH | ✅ 200 | 42.242373ms | - |
| 3.6 | Update Order Notes | PATCH | ✅ 200 | 12.831319ms | - |
| 4.1 | Process Payment — Cash (single split) | POST | ✅ 200 | 12.115966ms | - |
| 5.1 | Create Order (Takeaway) | POST | ✅ 201 | 9.728407ms | - |
| 5.2 | Process Payment — Split (Cash + QRIS) | POST | ✅ 200 | 25.188076ms | - |
| 5.3a | Create Order for Cancel | POST | ✅ 201 | 9.555693ms | - |
| 5.3 | Cancel Order (cashier) | PATCH | ✅ 200 | 10.101033ms | - |
| 5.4a | Create Order for Void | POST | ✅ 201 | 8.246691ms | - |
| 5.4 | Void Order (manager only) | PATCH | ✅ 200 | 12.615401ms | - |
| 6.1 | Process Return (Cashier) | POST | ✅ 200 | 17.366761ms | - |
| 7.1 | Open Shift | POST | ✅ 201 | 3.759321ms | - |
| 7.2 | Get Current Shift | GET | ✅ 200 | 2.145055ms | - |
| 7.3 | Close Shift | PATCH | ✅ 200 | 8.472361ms | - |
| 8.1 | Create Raw Material | POST | ✅ 201 | 3.803172ms | - |
| 8.2 | Get Raw Material | GET | ✅ 200 | 2.603231ms | - |
| 8.3 | Update Raw Material | PUT | ✅ 200 | 3.25187ms | - |
| 8.4a | Create Raw Material for Deletion | POST | ✅ 201 | 5.090517ms | - |
| 8.4 | Delete Raw Material | DELETE | ✅ 200 | 3.349821ms | - |
| 9.0 | Create Recipe (BOM) | POST | ✅ 200 | 10.300134ms | - |
| 9.1 | Get Recipe by Menu ID | GET | ✅ 200 | 4.338664ms | - |
| 9.2 | Update Recipe | PUT | ✅ 200 | 9.613512ms | - |
| 9.3 | Delete Recipe | DELETE | ✅ 200 | 2.476645ms | - |
| 10.1 | Get All Employees | GET | ✅ 200 | 4.176326ms | - |
| 10.2 | Create Employee (cashier) | POST | ✅ 201 | 2.836849927s | - |
| 10.3 | Update Employee | PUT | ✅ 200 | 3.703874ms | - |
| 10.4 | Delete Employee | DELETE | ✅ 200 | 5.547188ms | - |
| 11.1 | Sales Summary (this month) | GET | ✅ 200 | 27.319193ms | - |
| 11.2 | Best Sellers (top 5) | GET | ✅ 200 | 27.120521ms | - |
| 11.3 | Return Impact | GET | ✅ 200 | 5.862161ms | - |
| 11.4 | COGS | GET | ✅ 200 | 15.940747ms | - |
| 12.1 | Generate Receipt (splits-aware) | GET | ✅ 200 | 9.067796ms | - |
| 12.2 | Generate Kitchen Ticket | GET | ✅ 200 | 6.185164ms | - |
| 13.1 | Get Catalog Menus (Basic Auth) | GET | ✅ 200 | 34.083021ms | - |
| 13.2 | Catalog Self-Checkout (QRIS) | POST | ✅ 200 | 57.725309ms | - |

## Variables Captured

| Key | Value |
|-----|-------|
| order_number | ORD-20260520-4915 |
| recipe_id | 29d66cd9-9e70-4e57-a72d-6c549c0477d1 |
| catalog_menu_id | 161b5a19-0740-4897-a864-97f3fde0b37f |
| menu0_id | 13a023fd-38a5-4d23-8960-2645faee7bdf |
| order_total | 106720.00 |
| menu1_price | 12000 |
| menu2_id | 16e9a990-13ed-4f9a-a34b-e37f0c13f36d |
| confirmed_status | confirmed |
| emp_id | 8a8bbdc1-08ff-430a-bb3f-296e63e083fa |
| manager_branch_id | 2d54cb21-beb2-40ac-9850-275f179888b4 |
| cashier_token | eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ... |
| menu2_name | Chendol |
| menu3_name | Teh Tarik |
| order_status | pending |
| rm_name | Arabica Coffee Beans |
| cashier_id | 5dd7f438-7ddf-41a4-9a36-79ee6cbe076c |
| menu1_name | Kopi O |
| takeaway_id | cf60a915-9071-403e-8247-97f792049b8e |
| shift_id | ebd32d2d-ad6c-4976-9c47-94de60b43ecc |
| kitchen_token | eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ... |
| menu2_price | 20000 |
| takeaway_total | 8920.00 |
| rm_id | 77871c87-0751-49ca-a955-b5097b8b2c96 |
| manager_token | eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ... |
| manager_id | 7bf095d0-ccb5-4c6e-aaec-f59baaf437c6 |
| menu4_name | Matcha Latte |
| menu4_price | 22000 |
| shift_status | open |
| cashier_branch_id | 2d54cb21-beb2-40ac-9850-275f179888b4 |
| menu0_name | Nasi Lemak |
| menu0_price | 35000 |
| menu3_price | 16000 |
| menu4_id | 416d57b1-7396-4e1e-a48b-073cfdc57730 |
| order_id | 9a57acbc-40bd-4cd9-83df-ead45b20e136 |
| rm_delete_id | 1111b7f5-2284-4a15-8e19-68f0bd13e15b |
| menu1_id | 161b5a19-0740-4897-a864-97f3fde0b37f |
| menu3_id | 238f2491-52a3-4075-91f9-a3baacfe4480 |
