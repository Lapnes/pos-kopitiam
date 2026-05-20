-- Migration: Add indexes and foreign keys for POS KopiTiam v3NF
-- Run this AFTER AutoMigrate creates tables
-- Date: 2026-05-16

-- ============================================================
-- 1. FOREIGN KEY CONSTRAINTS (AutoMigrate creates some, but we ensure completeness)
-- ============================================================

-- Orders table FKs
ALTER TABLE orders
    ADD CONSTRAINT fk_orders_branch FOREIGN KEY (branch_id) REFERENCES branches(id) ON DELETE RESTRICT,
    ADD CONSTRAINT fk_orders_employee FOREIGN KEY (employee_id) REFERENCES users(id) ON DELETE RESTRICT,
    ADD CONSTRAINT fk_orders_customer FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE SET NULL,
    ADD CONSTRAINT fk_orders_shift FOREIGN KEY (shift_id) REFERENCES shifts(id) ON DELETE SET NULL;

-- Order details FKs
ALTER TABLE order_details
    -- ADD CONSTRAINT fk_details_order FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
    ADD CONSTRAINT fk_details_menu FOREIGN KEY (menu_id) REFERENCES menus(id) ON DELETE RESTRICT;

-- Payments FKs
-- ALTER TABLE payments
--     ADD CONSTRAINT fk_payments_order FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE;

-- Payment splits FKs
-- ALTER TABLE payment_splits
--     ADD CONSTRAINT fk_splits_payment FOREIGN KEY (payment_id) REFERENCES payments(id) ON DELETE CASCADE;

-- Reservations FKs
ALTER TABLE reservations
    ADD CONSTRAINT fk_reservations_branch FOREIGN KEY (branch_id) REFERENCES branches(id) ON DELETE CASCADE,
    -- ADD CONSTRAINT fk_reservations_table FOREIGN KEY (table_id) REFERENCES tables(id) ON DELETE CASCADE,
    ADD CONSTRAINT fk_reservations_customer FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE SET NULL;

-- Returns FKs
ALTER TABLE order_returns
    ADD CONSTRAINT fk_returns_order FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
    ADD CONSTRAINT fk_returns_processor FOREIGN KEY (processed_by) REFERENCES users(id) ON DELETE RESTRICT;

-- Stock adjustments FKs
ALTER TABLE stock_adjustments
    ADD CONSTRAINT fk_adjustments_material FOREIGN KEY (raw_material_id) REFERENCES raw_materials(id) ON DELETE RESTRICT,
    ADD CONSTRAINT fk_adjustments_user FOREIGN KEY (adjusted_by) REFERENCES users(id) ON DELETE RESTRICT;

-- Shifts FKs
ALTER TABLE shifts
    ADD CONSTRAINT fk_shifts_branch FOREIGN KEY (branch_id) REFERENCES branches(id) ON DELETE CASCADE,
    ADD CONSTRAINT fk_shifts_opener FOREIGN KEY (opened_by) REFERENCES users(id) ON DELETE RESTRICT,
    ADD CONSTRAINT fk_shifts_closer FOREIGN KEY (closed_by) REFERENCES users(id) ON DELETE SET NULL;

-- Cash movements FKs
ALTER TABLE cash_movements
    ADD CONSTRAINT fk_cashmovements_shift FOREIGN KEY (shift_id) REFERENCES shifts(id) ON DELETE CASCADE,
    ADD CONSTRAINT fk_cashmovements_user FOREIGN KEY (recorded_by) REFERENCES users(id) ON DELETE RESTRICT;

-- Recipes FKs
ALTER TABLE recipes
    ADD CONSTRAINT fk_recipes_menu FOREIGN KEY (menu_id) REFERENCES menus(id) ON DELETE CASCADE;

-- Recipe ingredients FKs
-- ALTER TABLE recipe_ingredients
--     ADD CONSTRAINT fk_ingredients_recipe FOREIGN KEY (recipe_id) REFERENCES recipes(id) ON DELETE CASCADE,
--     ADD CONSTRAINT fk_ingredients_material FOREIGN KEY (raw_material_id) REFERENCES raw_materials(id) ON DELETE RESTRICT;

-- Discount applications FKs
ALTER TABLE discount_applications
    ADD CONSTRAINT fk_discountapp_order FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
    ADD CONSTRAINT fk_discountapp_discount FOREIGN KEY (discount_id) REFERENCES discounts(id) ON DELETE RESTRICT;

-- Tax configs FKs
ALTER TABLE tax_configs
    ADD CONSTRAINT fk_tax_branch FOREIGN KEY (branch_id) REFERENCES branches(id) ON DELETE CASCADE;

-- Audit logs FKs
ALTER TABLE audit_logs
    ADD CONSTRAINT fk_audit_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT;

-- Users FKs
ALTER TABLE users
    ADD CONSTRAINT fk_users_branch FOREIGN KEY (branch_id) REFERENCES branches(id) ON DELETE RESTRICT;

-- ============================================================
-- 2. INDEXES (Critical for performance - PostgreSQL does NOT auto-index FKs)
-- ============================================================

-- Order indexes
-- CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_created_at ON orders(created_at DESC);
CREATE INDEX idx_orders_branch_status ON orders(branch_id, status);
CREATE INDEX idx_orders_branch_created ON orders(branch_id, created_at DESC);
CREATE INDEX idx_orders_customer ON orders(customer_id);
CREATE INDEX idx_orders_shift ON orders(shift_id);
CREATE INDEX idx_orders_order_number ON orders(order_number);

-- Order detail indexes
CREATE INDEX idx_details_order ON order_details(order_id);
CREATE INDEX idx_details_menu ON order_details(menu_id);
CREATE INDEX idx_details_voided ON order_details(is_voided);

-- Payment indexes
CREATE INDEX idx_payments_order ON payments(order_id);
CREATE INDEX idx_payments_status ON payments(status);

-- Payment split indexes
CREATE INDEX idx_splits_payment ON payment_splits(payment_id);
CREATE INDEX idx_splits_method ON payment_splits(payment_method);

-- Menu indexes
CREATE INDEX idx_menus_active ON menus(is_active);
CREATE INDEX idx_menus_station ON menus(station);
CREATE INDEX idx_menus_name ON menus(name);

-- Customer indexes
CREATE INDEX idx_customers_phone ON customers(phone);
CREATE INDEX idx_customers_name ON customers(name);

-- Reservation indexes
CREATE INDEX idx_reservations_branch ON reservations(branch_id);
CREATE INDEX idx_reservations_table ON reservations(table_id);
CREATE INDEX idx_reservations_date ON reservations(reservation_date);
CREATE INDEX idx_reservations_status ON reservations(status);

-- Return indexes
CREATE INDEX idx_returns_order ON order_returns(order_id);
CREATE INDEX idx_returns_date ON order_returns(processed_at);

-- Raw material indexes
CREATE INDEX idx_rawmaterials_name ON raw_materials(name);
CREATE INDEX idx_rawmaterials_stock ON raw_materials(current_stock);

-- Recipe indexes
CREATE INDEX idx_recipes_menu ON recipes(menu_id);

-- Recipe ingredient indexes
CREATE INDEX idx_ingredients_recipe ON recipe_ingredients(recipe_id);
CREATE INDEX idx_ingredients_material ON recipe_ingredients(raw_material_id);

-- Shift indexes
CREATE INDEX idx_shifts_branch ON shifts(branch_id);
CREATE INDEX idx_shifts_status ON shifts(status);
CREATE INDEX idx_shifts_opener ON shifts(opened_by);
CREATE INDEX idx_shifts_opening ON shifts(opening_time);

-- Cash movement indexes
CREATE INDEX idx_cashmovements_shift ON cash_movements(shift_id);
CREATE INDEX idx_cashmovements_type ON cash_movements(type);
CREATE INDEX idx_cashmovements_date ON cash_movements(recorded_at);

-- Discount indexes
CREATE INDEX idx_discounts_active ON discounts(is_active);
CREATE INDEX idx_discounts_type ON discounts(type);

-- Tax indexes
CREATE INDEX idx_tax_active ON tax_configs(is_active);
CREATE INDEX idx_tax_type ON tax_configs(tax_type);
CREATE INDEX idx_tax_branch ON tax_configs(branch_id);

-- User indexes
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_active ON users(is_active);
CREATE INDEX idx_users_branch ON users(branch_id);

-- Audit indexes
CREATE INDEX idx_audit_user ON audit_logs(user_id);
CREATE INDEX idx_audit_action ON audit_logs(action);
CREATE INDEX idx_audit_entity ON audit_logs(entity);
CREATE INDEX idx_audit_timestamp ON audit_logs(timestamp);

-- Branch indexes
CREATE INDEX idx_branches_name ON branches(name);

-- ============================================================
-- 3. COMPOSITE INDEXES (For common query patterns)
-- ============================================================

-- Common dashboard query: orders by branch + status + date
CREATE INDEX idx_orders_dashboard ON orders(branch_id, status, created_at DESC);

-- Analytics: order details by order + void status
CREATE INDEX idx_details_analytics ON order_details(order_id, is_voided, subtotal);

-- Shift reconciliation: cash movements by shift + type
CREATE INDEX idx_cashmovements_recon ON cash_movements(shift_id, type, amount);

-- Menu profitability: ingredients by recipe
CREATE INDEX idx_ingredients_profit ON recipe_ingredients(recipe_id, raw_material_id, quantity);
