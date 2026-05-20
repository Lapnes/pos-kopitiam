package repository

import (
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/models"
	"gorm.io/gorm"
)

type OrderRepository interface {
	Create(order *models.Order) error
	FindByID(id string) (*models.Order, error)
	FindByOrderNumber(orderNumber string) (*models.Order, error)
	FindAllWithDetails() ([]models.Order, error)
	Update(order *models.Order) error
}

type orderRepo struct {
	db *gorm.DB
}

func NewOrderRepo(db *gorm.DB) OrderRepository {
	return &orderRepo{db: db}
}

func (r *orderRepo) Create(order *models.Order) error {
	return r.db.Create(order).Error
}

// FindByID — single order, Preload acceptable (1 order = 1 query + preloads)
func (r *orderRepo) FindByID(id string) (*models.Order, error) {
	var order models.Order
	err := r.db.Preload("OrderDetails").
		Preload("Payments.Splits").
		Preload("Payments").
		First(&order, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// FindByOrderNumber — single order, Preload acceptable
func (r *orderRepo) FindByOrderNumber(orderNumber string) (*models.Order, error) {
	var order models.Order
	err := r.db.Preload("OrderDetails").
		Preload("Payments.Splits").
		Preload("Payments").
		First(&order, "order_number = ?", orderNumber).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// FindAllWithDetails — LIST query: Preload for hasMany (avoids row multiplication)
// GORM Preload uses 2 queries: 1 for orders, 1 batched IN for details.
// This is optimal for collections (hasMany) vs JOINs which duplicate parent rows.
func (r *orderRepo) FindAllWithDetails() ([]models.Order, error) {
	var orders []models.Order

	// Step 1: Fetch all orders with derived aggregates via JOIN
	// We use a subquery approach to get totals without N+1
	err := r.db.Model(&models.Order{}).
		Select(`
			orders.*,
			COALESCE(od_totals.subtotal, 0) as subtotal,
			COALESCE(p_totals.total_paid, 0) as total_paid
		`).
		Joins(`LEFT JOIN (
			SELECT order_id, SUM(CASE WHEN is_voided = 0 THEN subtotal ELSE 0 END) as subtotal
			FROM order_details
			WHERE deleted_at IS NULL
			GROUP BY order_id
		) od_totals ON od_totals.order_id = orders.id`).
		Joins(`LEFT JOIN (
			SELECT p.order_id, SUM(ps.amount) as total_paid
			FROM payments p
			LEFT JOIN payment_splits ps ON ps.payment_id = p.id AND ps.deleted_at IS NULL
			WHERE p.deleted_at IS NULL AND p.status = 'success'
			GROUP BY p.order_id
		) p_totals ON p_totals.order_id = orders.id`).
		Where("orders.deleted_at IS NULL").
		Group("orders.id").
		Order("orders.created_at DESC").
		Find(&orders).Error

	if err != nil {
		return nil, err
	}

	// Step 2: If we need OrderDetails for each order, preload them in a second batched query
	// This avoids the N+1 problem — GORM batches all detail lookups into ONE IN query
	if len(orders) > 0 {
		orderIDs := make([]string, len(orders))
		for i, o := range orders {
			orderIDs[i] = o.ID.String()
		}

		var allDetails []models.OrderDetail
		if err := r.db.Where("order_id IN ? AND deleted_at IS NULL", orderIDs).
			Find(&allDetails).Error; err != nil {
			return nil, err
		}

		// Map details to orders in-memory (O(n) instead of O(n²))
		detailsMap := make(map[string][]models.OrderDetail)
		for _, d := range allDetails {
			key := d.OrderID.String()
			detailsMap[key] = append(detailsMap[key], d)
		}
		for i := range orders {
			orders[i].OrderDetails = detailsMap[orders[i].ID.String()]
		}
	}

	return orders, nil
}

func (r *orderRepo) Update(order *models.Order) error {
	return r.db.Save(order).Error
}
