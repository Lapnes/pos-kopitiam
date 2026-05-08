package config

import "gorm.io/gorm"

func SetupDatabaseLogic(db *gorm.DB) {

	// ======================
	// DROP (SAFE RE-RUN)
	// ======================
	db.Exec(`DROP TRIGGER IF EXISTS before_insert_order_details`)
	db.Exec(`DROP TRIGGER IF EXISTS before_update_order_details`)
	db.Exec(`DROP TRIGGER IF EXISTS after_insert_order_details`)
	db.Exec(`DROP TRIGGER IF EXISTS after_update_order_details`)
	db.Exec(`DROP TRIGGER IF EXISTS after_delete_order_details`)

	db.Exec(`DROP PROCEDURE IF EXISTS create_simple_order`)
	db.Exec(`DROP PROCEDURE IF EXISTS get_total_sales`)

	db.Exec(`DROP FUNCTION IF EXISTS get_order_total`)

	// ======================
	// TRIGGER: VALIDATION + SUBTOTAL
	// ======================
	db.Exec(`
	CREATE TRIGGER before_insert_order_details
	BEFORE INSERT ON order_details
	FOR EACH ROW
	BEGIN
		IF NEW.quantity <= 0 THEN
			SIGNAL SQLSTATE '45000'
			SET MESSAGE_TEXT = 'Quantity must be greater than 0';
		END IF;

		IF NEW.unit_price < 0 THEN
			SIGNAL SQLSTATE '45000'
			SET MESSAGE_TEXT = 'Price cannot be negative';
		END IF;

		SET NEW.subtotal = NEW.quantity * NEW.unit_price;
	END;
	`)

	db.Exec(`
	CREATE TRIGGER before_update_order_details
	BEFORE UPDATE ON order_details
	FOR EACH ROW
	BEGIN
		IF NEW.quantity <= 0 THEN
			SIGNAL SQLSTATE '45000'
			SET MESSAGE_TEXT = 'Quantity must be greater than 0';
		END IF;

		IF NEW.unit_price < 0 THEN
			SIGNAL SQLSTATE '45000'
			SET MESSAGE_TEXT = 'Price cannot be negative';
		END IF;

		SET NEW.subtotal = NEW.quantity * NEW.unit_price;
	END;
	`)

	// ======================
	// TRIGGER: UPDATE TOTAL
	// ======================
	db.Exec(`
	CREATE TRIGGER after_insert_order_details
	AFTER INSERT ON order_details
	FOR EACH ROW
	BEGIN
		UPDATE orders
		SET total_price = (
			SELECT IFNULL(SUM(subtotal),0)
			FROM order_details
			WHERE order_id = NEW.order_id
		)
		WHERE order_id = NEW.order_id;
	END;
	`)

	db.Exec(`
	CREATE TRIGGER after_update_order_details
	AFTER UPDATE ON order_details
	FOR EACH ROW
	BEGIN
		UPDATE orders
		SET total_price = (
			SELECT IFNULL(SUM(subtotal),0)
			FROM order_details
			WHERE order_id = NEW.order_id
		)
		WHERE order_id = NEW.order_id;
	END;
	`)

	db.Exec(`
	CREATE TRIGGER after_delete_order_details
	AFTER DELETE ON order_details
	FOR EACH ROW
	BEGIN
		UPDATE orders
		SET total_price = (
			SELECT IFNULL(SUM(subtotal),0)
			FROM order_details
			WHERE order_id = OLD.order_id
		)
		WHERE order_id = OLD.order_id;
	END;
	`)

	// ======================
	// FUNCTION: GET ORDER TOTAL
	// ======================
	db.Exec(`
	CREATE FUNCTION get_order_total(p_order_id CHAR(36))
	RETURNS BIGINT
	READS SQL DATA
	DETERMINISTIC
	BEGIN
		DECLARE total BIGINT;

		SELECT IFNULL(SUM(subtotal),0)
		INTO total
		FROM order_details
		WHERE order_id = p_order_id;

		RETURN total;
	END;
	`)

	// ======================
	// PROCEDURE: CREATE ORDER
	// ======================
	db.Exec(`
	CREATE PROCEDURE create_simple_order (
		IN p_employee_id BIGINT,
		IN p_menu_id BIGINT,
		IN p_qty INT
	)
	BEGIN
		DECLARE new_order_id CHAR(36);
		DECLARE menu_price BIGINT;

		IF p_qty <= 0 THEN
			SIGNAL SQLSTATE '45000'
			SET MESSAGE_TEXT = 'Quantity must be greater than 0';
		END IF;

		SELECT price INTO menu_price 
		FROM menus 
		WHERE menu_id = p_menu_id;

		IF menu_price IS NULL THEN
			SIGNAL SQLSTATE '45000'
			SET MESSAGE_TEXT = 'Menu not found';
		END IF;

		SET new_order_id = UUID();

		INSERT INTO orders(order_id, employee_id, total_price)
		VALUES (new_order_id, p_employee_id, 0);

		INSERT INTO order_details(order_detail_id, order_id, menu_id, quantity, unit_price)
		VALUES (UUID(), new_order_id, p_menu_id, p_qty, menu_price);
	END;
	`)

	// ======================
	// PROCEDURE: TOTAL SALES
	// ======================
	db.Exec(`
	CREATE PROCEDURE get_total_sales()
	BEGIN
		SELECT IFNULL(SUM(total_price),0) AS total_sales FROM orders;
	END;
	`)
}
