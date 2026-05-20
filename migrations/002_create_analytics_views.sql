-- Materialized Views for POS KopiTiam v3NF Analytics Dashboard
-- Run this AFTER migration 001 is applied
-- Refresh schedule: every 5 minutes for real-time-ish dashboards
-- Use: REFRESH OR REPLACE VIEW CONCURRENTLY view_name;

-- ============================================================
-- 1. v_daily_sales_summary — Grafik dasbor harian manajer
-- ============================================================
CREATE OR REPLACE VIEW v_daily_sales_summary AS
SELECT
    DATE(o.created_at) AS sale_date,
    o.branch_id,
    COUNT(DISTINCT o.id) AS total_orders,
    SUM(CASE WHEN od.is_voided = false THEN od.subtotal ELSE 0 END) AS gross_revenue,
    SUM(CASE WHEN od.is_voided = false THEN od.cost_price * od.quantity ELSE 0 END) AS total_cogs,
    SUM(CASE WHEN od.is_voided = false THEN od.subtotal ELSE 0 END) -
        SUM(CASE WHEN od.is_voided = false THEN od.cost_price * od.quantity ELSE 0 END) AS net_profit,
    AVG(CASE WHEN od.is_voided = false THEN od.subtotal ELSE NULL END) AS avg_order_value,
    COUNT(CASE WHEN o.status = 'cancelled' THEN 1 END) AS cancelled_orders
FROM orders o
LEFT JOIN order_details od ON od.order_id = o.id AND od.deleted_at IS NULL
WHERE o.deleted_at IS NULL
    AND o.status IN ('confirmed', 'paid', 'served')
GROUP BY DATE(o.created_at), o.branch_id
ORDER BY sale_date DESC;





-- ============================================================
-- 2. v_weekly_sales_summary — Grafik dasbor mingguan manajer
-- ============================================================
CREATE OR REPLACE VIEW v_weekly_sales_summary AS
SELECT
    DATE_SUB(DATE(o.created_at), INTERVAL WEEKDAY(o.created_at) DAY) AS week_start,
    o.branch_id,
    COUNT(DISTINCT o.id) AS total_orders,
    SUM(CASE WHEN od.is_voided = false THEN od.subtotal ELSE 0 END) AS gross_revenue,
    SUM(CASE WHEN od.is_voided = false THEN od.cost_price * od.quantity ELSE 0 END) AS total_cogs,
    SUM(CASE WHEN od.is_voided = false THEN od.subtotal ELSE 0 END) -
        SUM(CASE WHEN od.is_voided = false THEN od.cost_price * od.quantity ELSE 0 END) AS net_profit,
    COUNT(DISTINCT DATE(o.created_at)) AS active_days,
    MAX(o.created_at) AS last_order_at
FROM orders o
LEFT JOIN order_details od ON od.order_id = o.id AND od.deleted_at IS NULL
WHERE o.deleted_at IS NULL
    AND o.status IN ('confirmed', 'paid', 'served')
GROUP BY DATE_SUB(DATE(o.created_at), INTERVAL WEEKDAY(o.created_at) DAY), o.branch_id
ORDER BY week_start DESC;




-- ============================================================
-- 3. v_payment_method_summary — Rekonsiliasi QRIS/Transfer vs uang fisik
-- ============================================================
CREATE OR REPLACE VIEW v_payment_method_summary AS
SELECT
    DATE(p.created_at) AS payment_date,
    o.branch_id,
    ps.payment_method,
    COUNT(DISTINCT p.id) AS total_payments,
    SUM(ps.amount) AS total_amount,
    COUNT(DISTINCT o.id) AS order_count,
    AVG(ps.amount) AS avg_amount
FROM payments p
JOIN payment_splits ps ON ps.payment_id = p.id AND ps.deleted_at IS NULL
JOIN orders o ON o.id = p.order_id
WHERE p.deleted_at IS NULL
    AND p.status = 'success'
    AND o.deleted_at IS NULL
GROUP BY DATE(p.created_at), o.branch_id, ps.payment_method
ORDER BY payment_date DESC, total_amount DESC;





-- ============================================================
-- 4. v_menu_profitability_live — Margin keuntungan real-time per menu
-- ============================================================
CREATE OR REPLACE VIEW v_menu_profitability_live AS
SELECT
    od.menu_id,
    od.menu_name,
    m.station,
    SUM(CASE WHEN od.is_voided = false THEN od.quantity ELSE 0 END) AS total_qty_sold,
    SUM(CASE WHEN od.is_voided = false THEN od.subtotal ELSE 0 END) AS total_revenue,
    SUM(CASE WHEN od.is_voided = false THEN od.cost_price * od.quantity ELSE 0 END) AS total_cogs,
    SUM(CASE WHEN od.is_voided = false THEN od.subtotal ELSE 0 END) -
        SUM(CASE WHEN od.is_voided = false THEN od.cost_price * od.quantity ELSE 0 END) AS net_profit,
    CASE WHEN SUM(CASE WHEN od.is_voided = false THEN od.subtotal ELSE 0 END) > 0
        THEN (SUM(CASE WHEN od.is_voided = false THEN od.subtotal ELSE 0 END) -
              SUM(CASE WHEN od.is_voided = false THEN od.cost_price * od.quantity ELSE 0 END)) /
             SUM(CASE WHEN od.is_voided = false THEN od.subtotal ELSE 0 END) * 100
        ELSE 0
    END AS profit_margin_pct,
    AVG(od.price) AS current_selling_price,
    -- Recipe cost if exists
    (SELECT SUM(rm.cost_per_unit * ri.quantity)
     FROM recipes r
     JOIN recipe_ingredients ri ON ri.recipe_id = r.id
     JOIN raw_materials rm ON rm.id = ri.raw_material_id
     WHERE r.menu_id = od.menu_id
     LIMIT 1) AS current_recipe_cost,
    NOW() AS refreshed_at
FROM order_details od
JOIN orders o ON o.id = od.order_id AND o.deleted_at IS NULL
LEFT JOIN menus m ON m.id = od.menu_id
WHERE od.deleted_at IS NULL
    AND o.status IN ('confirmed', 'paid', 'served')
GROUP BY od.menu_id, od.menu_name, m.station
ORDER BY net_profit DESC;





-- ============================================================
-- 5. v_best_selling_items — Langsung ditembak endpoint best-sellers
-- ============================================================
CREATE OR REPLACE VIEW v_best_selling_items AS
SELECT
    od.menu_id,
    od.menu_name,
    m.station,
    SUM(CASE WHEN od.is_voided = false THEN od.quantity ELSE 0 END) AS total_quantity,
    SUM(CASE WHEN od.is_voided = false THEN od.subtotal ELSE 0 END) AS total_sales,
    COUNT(DISTINCT o.id) AS order_count,
    AVG(od.price) AS avg_price,
    MAX(o.created_at) AS last_sold_at,
    -- Trend: compare last 7 days vs previous 7 days
    SUM(CASE WHEN od.is_voided = false AND o.created_at >= NOW() - INTERVAL 7 DAY
        THEN od.quantity ELSE 0 END) AS qty_last_7d,
    SUM(CASE WHEN od.is_voided = false AND o.created_at >= NOW() - INTERVAL 14 DAY
        AND o.created_at < NOW() - INTERVAL 7 DAY
        THEN od.quantity ELSE 0 END) AS qty_prev_7d
FROM order_details od
JOIN orders o ON o.id = od.order_id AND o.deleted_at IS NULL
LEFT JOIN menus m ON m.id = od.menu_id
WHERE od.deleted_at IS NULL
    AND o.status IN ('confirmed', 'paid', 'served')
GROUP BY od.menu_id, od.menu_name, m.station
ORDER BY total_quantity DESC;





-- ============================================================
-- 6. v_shift_reconciliation — Rapor keandalan kasir
-- ============================================================
CREATE OR REPLACE VIEW v_shift_reconciliation AS
SELECT
    s.id AS shift_id,
    s.branch_id,
    s.opened_by AS cashier_id,
    u.name AS cashier_name,
    s.opening_time,
    s.closing_time,
    s.opening_cash,
    s.actual_closing_cash,
    COALESCE(SUM(CASE WHEN cm.type = 'in' THEN cm.amount ELSE 0 END), 0) AS total_cash_in,
    COALESCE(SUM(CASE WHEN cm.type = 'out' THEN cm.amount ELSE 0 END), 0) AS total_cash_out,
    s.opening_cash + COALESCE(SUM(CASE WHEN cm.type = 'in' THEN cm.amount ELSE -cm.amount END), 0) AS expected_closing_cash,
    s.actual_closing_cash - (s.opening_cash + COALESCE(SUM(CASE WHEN cm.type = 'in' THEN cm.amount ELSE -cm.amount END), 0)) AS cash_difference,
    COUNT(DISTINCT o.id) AS total_orders,
    COALESCE(SUM(p.total_paid), 0) AS total_sales,
    CASE WHEN s.actual_closing_cash > 0
        THEN s.actual_closing_cash - (s.opening_cash + COALESCE(SUM(CASE WHEN cm.type = 'in' THEN cm.amount ELSE -cm.amount END), 0))
        ELSE 0
    END AS discrepancy,
    s.status
FROM shifts s
JOIN users u ON u.id = s.opened_by
LEFT JOIN cash_movements cm ON cm.shift_id = s.id AND cm.deleted_at IS NULL
LEFT JOIN orders o ON o.shift_id = s.id AND o.deleted_at IS NULL AND o.status IN ('confirmed', 'paid', 'served')
LEFT JOIN (
    SELECT p.order_id, SUM(ps.amount) AS total_paid
    FROM payments p
    JOIN payment_splits ps ON ps.payment_id = p.id AND ps.deleted_at IS NULL
    WHERE p.deleted_at IS NULL AND p.status = 'success'
    GROUP BY p.order_id
) p ON p.order_id = o.id
WHERE s.deleted_at IS NULL
GROUP BY s.id, s.branch_id, s.opened_by, u.name, s.opening_time, s.closing_time,
         s.opening_cash, s.actual_closing_cash, s.status
ORDER BY s.opening_time DESC;






-- ============================================================
-- 7. v_void_and_return_logs — Pantau anomali void/return
-- ============================================================
CREATE OR REPLACE VIEW v_void_and_return_logs AS
SELECT
    DATE(o.updated_at) AS event_date,
    o.branch_id,
    o.employee_id,
    u.name AS employee_name,
    'void_order' AS event_type,
    o.id AS order_id,
    o.order_number,
    o.total AS amount,
    o.notes AS reason,
    COUNT(od.id) AS item_count,
    SUM(od.quantity) AS total_qty_voided
FROM orders o
JOIN users u ON u.id = o.employee_id
JOIN order_details od ON od.order_id = o.id AND od.deleted_at IS NULL AND od.is_voided = true
WHERE o.deleted_at IS NULL
    AND o.status = 'cancelled'
    AND o.updated_at >= NOW() - INTERVAL 30 DAY
GROUP BY DATE(o.updated_at), o.branch_id, o.employee_id, u.name, o.id, o.order_number, o.total, o.notes

UNION ALL

SELECT
    DATE(r.processed_at) AS event_date,
    o.branch_id,
    r.processed_by AS employee_id,
    u.name AS employee_name,
    'return' AS event_type,
    o.id AS order_id,
    o.order_number,
    r.return_amount AS amount,
    r.reason,
    0 AS item_count,
    0 AS total_qty_voided
FROM order_returns r
JOIN orders o ON o.id = r.order_id
JOIN users u ON u.id = r.processed_by
WHERE r.deleted_at IS NULL
    AND r.processed_at >= NOW() - INTERVAL 30 DAY

ORDER BY event_date DESC, amount DESC;






-- ============================================================
-- 8. v_cashier_performance — Karyawan Terbaik / Bonus system
-- ============================================================
CREATE OR REPLACE VIEW v_cashier_performance AS
SELECT
    o.employee_id AS cashier_id,
    u.name AS cashier_name,
    u.branch_id,
    COUNT(DISTINCT o.id) AS total_orders,
    SUM(CASE WHEN od.is_voided = false THEN od.subtotal ELSE 0 END) AS total_sales,
    AVG(CASE WHEN od.is_voided = false THEN od.subtotal ELSE NULL END) AS avg_order_value,
    COUNT(DISTINCT DATE(o.created_at)) AS active_days,
    MAX(o.created_at) AS last_order_at,
    -- Performance by time period
    SUM(CASE WHEN od.is_voided = false AND o.created_at >= NOW() - INTERVAL 7 DAY
        THEN od.subtotal ELSE 0 END) AS sales_last_7d,
    SUM(CASE WHEN od.is_voided = false AND o.created_at >= NOW() - INTERVAL 30 DAY
        THEN od.subtotal ELSE 0 END) AS sales_last_30d,
    -- Void rate (lower is better)
    CASE WHEN COUNT(od.id) > 0
        THEN COUNT(CASE WHEN od.is_voided = true THEN 1 END)* 1.0 / COUNT(od.id) * 100
        ELSE 0
    END AS void_rate_pct,
    -- Customer satisfaction proxy: repeat customers
    COUNT(DISTINCT CASE WHEN o.customer_id IS NOT NULL THEN o.customer_id END) AS unique_customers,
    COUNT(DISTINCT CASE WHEN o.order_type = 'dine_in' THEN o.id END) AS dine_in_orders,
    COUNT(DISTINCT CASE WHEN o.order_type = 'take_away' THEN o.id END) AS take_away_orders
FROM orders o
JOIN users u ON u.id = o.employee_id
LEFT JOIN order_details od ON od.order_id = o.id AND od.deleted_at IS NULL
WHERE o.deleted_at IS NULL
    AND o.status IN ('confirmed', 'paid', 'served')
    AND o.created_at >= NOW() - INTERVAL 90 DAY
GROUP BY o.employee_id, u.name, u.branch_id
ORDER BY total_sales DESC;





