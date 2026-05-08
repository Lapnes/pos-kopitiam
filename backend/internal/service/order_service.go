package service

import (
	"fmt"

	"github.com/Lapnes/pos-kopitiam/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func CreateOrder(db *gorm.DB, employeeID *uint64, items []struct {
	MenuID uint64
	Qty    int
}) (*models.Order, error) {

	var createdOrder models.Order

	err := db.Transaction(func(tx *gorm.DB) error {

		createdOrder = models.Order{
			EmployeeID: employeeID,
		}

		if err := tx.Create(&createdOrder).Error; err != nil {
			return err
		}

		for _, item := range items {
			var menu models.Menu

			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
				First(&menu, item.MenuID).Error; err != nil {
				return fmt.Errorf("menu %d not found", item.MenuID)
			}

			if menu.DailyStock < int64(item.Qty) {
				return fmt.Errorf("stock not enough for %s", menu.MenuName)
			}

			// Update stock
			if err := tx.Model(&menu).
				Update("daily_stock", menu.DailyStock-int64(item.Qty)).Error; err != nil {
				return err
			}

			// Tidak perlu subtotal (trigger yang handle)
			detail := models.OrderDetail{
				OrderID:   createdOrder.OrderID,
				MenuID:    menu.MenuID,
				Quantity:  item.Qty,
				UnitPrice: menu.Price,
			}

			if err := tx.Create(&detail).Error; err != nil {
				return err
			}
		}

		return nil
	})

	return &createdOrder, err
}
