package config

// import "gorm.io/gorm"

// func AddForeignKeys(db *gorm.DB) error {
// 	// Hapus dulu kalau ada (untuk idempoten)
// 	db.Exec("ALTER TABLE menus DROP FOREIGN KEY IF EXISTS fk_menus_category")
// 	db.Exec("ALTER TABLE orders DROP FOREIGN KEY IF EXISTS fk_orders_employee")
// 	db.Exec("ALTER TABLE order_details DROP FOREIGN KEY IF EXISTS fk_order_details_order")
// 	db.Exec("ALTER TABLE order_details DROP FOREIGN KEY IF EXISTS fk_order_details_menu")

// 	// Tambah FK constraints
// 	sqls := []string{
// 		// Menu -> Category
// 		`ALTER TABLE menus
// 		 ADD CONSTRAINT fk_menus_category
// 		 FOREIGN KEY (category_id) REFERENCES categories(category_id)
// 		 ON DELETE RESTRICT ON UPDATE CASCADE`,

// 		// Order -> Employee (nullable)
// 		`ALTER TABLE orders
// 		 ADD CONSTRAINT fk_orders_employee
// 		 FOREIGN KEY (employee_id) REFERENCES employees(employee_id)
// 		 ON DELETE SET NULL ON UPDATE CASCADE`,

// 		// OrderDetail -> Order
// 		`ALTER TABLE order_details
// 		 ADD CONSTRAINT fk_order_details_order
// 		 FOREIGN KEY (order_id) REFERENCES orders(order_id)
// 		 ON DELETE CASCADE ON UPDATE CASCADE`,

// 		// OrderDetail -> Menu
// 		`ALTER TABLE order_details
// 		 ADD CONSTRAINT fk_order_details_menu
// 		 FOREIGN KEY (menu_id) REFERENCES menus(menu_id)
// 		 ON DELETE RESTRICT ON UPDATE CASCADE`,
// 	}

// 	for _, sql := range sqls {
// 		if err := db.Exec(sql).Error; err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }
